package background

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/speakeasy-api/gram/internal/background/activities"
	slack_client "github.com/speakeasy-api/gram/internal/thirdparty/slack/client"
	"github.com/speakeasy-api/gram/internal/thirdparty/slack/types"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type ProcessSlackWorkflowParams struct {
	Event types.SlackEvent
}

type ProcessSlackEventResult struct {
	Status string
}

func ExecuteProcessSlackEventWorkflow(ctx context.Context, temporalClient client.Client, params ProcessSlackWorkflowParams) (client.WorkflowRun, error) {
	id := params.Event.EventID
	return temporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:                       fmt.Sprintf("v1:slack-event:%s", id),
		TaskQueue:                string(TaskQueueMain),
		WorkflowIDConflictPolicy: enums.WORKFLOW_ID_CONFLICT_POLICY_FAIL,
		WorkflowIDReusePolicy:    enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
		WorkflowRunTimeout:       time.Minute * 5,
	}, SlackEventWorkflow, params)
}

func SlackEventWorkflow(ctx workflow.Context, params ProcessSlackWorkflowParams) (*ProcessSlackEventResult, error) {
	var a *Activities

	logger := workflow.GetLogger(ctx)
	logger.Info("received slack event", slog.Any("event", params.Event))

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	})

	if params.Event.Event.Text == "" {
		return postSlackErrorMessage(ctx, a, params.Event, fmt.Errorf("no content found in prompt"))
	}

	var toolsResponse activities.SlackProjectContextResponse
	err := workflow.ExecuteActivity(
		ctx,
		a.GetSlackProjectContext,
		params.Event,
	).Get(ctx, &toolsResponse)
	if err != nil {
		return postSlackErrorMessage(ctx, a, params.Event, fmt.Errorf("failed to get slack project context: %w", err))
	}

	// Remove bot tag from the prompt
	words := strings.Fields(params.Event.Event.Text)
	if strings.HasPrefix(words[0], fmt.Sprintf("<@%s>", params.Event.Authorizations[0].UserID)) {
		words = words[1:]
		if len(words) == 0 {
			return postSlackErrorMessage(ctx, a, params.Event, fmt.Errorf("no content found in prompt"))
		}
	}

	if params.Event.Event.Type == "message" && params.Event.Event.ChannelType == "channel" {
		// If we are in a thread only go to chat completions if 'gram' is in the first two words (case-insensitive)
		maxCheck := min(2, len(words))
		gramIdx := -1
		for i := 0; i < maxCheck; i++ {
			if strings.ToLower(words[i]) == "gram" {
				gramIdx = i
				break
			}
		}
		if gramIdx == -1 {
			return &ProcessSlackEventResult{Status: "success"}, nil
		}
		// Remove all words up to and including 'gram'
		words = words[gramIdx+1:]
		if len(words) == 0 {
			return postSlackErrorMessage(ctx, a, params.Event, fmt.Errorf("no content found in prompt after 'gram'"))
		}
	}

	if words[0] == "list" && (len(words) == 1 || words[1] == "tools" || words[1] == "toolsets") {
		// List all toolsets
		if err := workflow.ExecuteActivity(
			ctx,
			a.PostSlackMessage,
			activities.PostSlackMessageInput{
				Event: params.Event,
				PostInput: slack_client.SlackPostMessageInput{
					ChannelID: params.Event.Event.Channel,
					Message:   formatListToolsSlackMessage(toolsResponse),
					ThreadTS:  &params.Event.Event.Ts,
				},
			},
		).Get(ctx, nil); err != nil {
			return nil, fmt.Errorf("failed to post slack response: %w", err)
		}
		return &ProcessSlackEventResult{
			Status: "success",
		}, nil
	}

	// Toolset selection: look for [toolset] or (toolset) as the first word
	potentialSelectedToolset := ""
	if (strings.HasPrefix(words[0], "[") && strings.HasSuffix(words[0], "]")) ||
		(strings.HasPrefix(words[0], "(") && strings.HasSuffix(words[0], ")")) {
		potentialSelectedToolset = words[0][1 : len(words[0])-1]
	}

	chosenToolsetSlug := ""
	if toolsResponse.DefaultToolsetSlug != nil {
		chosenToolsetSlug = *toolsResponse.DefaultToolsetSlug
	}
	if potentialSelectedToolset != "" {
		for _, toolset := range toolsResponse.Toolsets {
			if toolset.Slug == potentialSelectedToolset {
				words = words[1:] // take out the toolset slug from prompt
				chosenToolsetSlug = toolset.Slug
				break
			}
		}
	}

	sanitizedPrompt := strings.Join(words, " ")
	var chatResponse string
	err = workflow.ExecuteActivity(
		ctx,
		a.SlackChatCompletion,
		activities.SlackChatCompletionInput{
			Event:       params.Event,
			Prompt:      sanitizedPrompt,
			ToolsetSlug: chosenToolsetSlug,
		},
	).Get(ctx, &chatResponse)
	if err != nil {
		return postSlackErrorMessage(ctx, a, params.Event, err)
	}

	if err := workflow.ExecuteActivity(
		ctx,
		a.PostSlackMessage,
		activities.PostSlackMessageInput{
			Event: params.Event,
			PostInput: slack_client.SlackPostMessageInput{
				ChannelID: params.Event.Event.Channel,
				Message:   chatResponse,
				ThreadTS:  &params.Event.Event.Ts,
			},
		},
	).Get(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to post slack response: %w", err)
	}

	return &ProcessSlackEventResult{
		Status: "success",
	}, nil
}

func formatListToolsSlackMessage(input activities.SlackProjectContextResponse) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("*Project:* `%s`\n", input.ProjectSlug))
	if input.DefaultToolsetSlug != nil {
		sb.WriteString(fmt.Sprintf("*Default Toolset:* `%s`\n", *input.DefaultToolsetSlug))
	}
	sb.WriteString("\n*Toolsets:*\n")

	for _, ts := range input.Toolsets {
		sb.WriteString(fmt.Sprintf("• *`%s`* (%d tools)\n", ts.Slug, ts.NumberOfTools))
		if ts.Description != nil && *ts.Description != "" {
			sb.WriteString(fmt.Sprintf("  _%s_\n", *ts.Description))
		}
		sb.WriteString(fmt.Sprintf("  created at: `%s`\n", ts.CreatedAt))
		sb.WriteString(fmt.Sprintf("  updated at: `%s`\n\n", ts.UpdatedAt))
	}

	return sb.String()
}

func postSlackErrorMessage(ctx workflow.Context, a *Activities, event types.SlackEvent, err error) (*ProcessSlackEventResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Error("error in slack event workflow", slog.String("error", err.Error()))
	msg := "*Error:* \n Apologies I am unable to complete your request."
	activityErr := workflow.ExecuteActivity(
		ctx,
		a.PostSlackMessage,
		activities.PostSlackMessageInput{
			Event: event,
			PostInput: slack_client.SlackPostMessageInput{
				ChannelID: event.Event.Channel,
				Message:   msg,
				ThreadTS:  &event.Event.Ts,
			},
		},
	).Get(ctx, nil)
	if activityErr != nil {
		return nil, activityErr
	}
	return &ProcessSlackEventResult{
		Status: "failed",
	}, nil
}
