package toolsets

import (
	"github.com/speakeasy-api/gram/server/design/security"
	"github.com/speakeasy-api/gram/server/design/shared"
	. "goa.design/goa/v3/dsl"
)

var _ = Service("chat", func() {
	Description("Managed chats for gram AI consumers.")
	Security(security.Session, security.ProjectSlug)
	shared.DeclareErrorResponses()

	Method("listChats", func() {
		Description("List all chats for a project")

		Payload(func() {
			security.SessionPayload()
			security.ProjectPayload()
		})

		Result(ListChatsResult)

		HTTP(func() {
			GET("/rpc/chat.list")
			security.SessionHeader()
			security.ProjectHeader()
			Response(StatusOK)
		})

		Meta("openapi:operationId", "listChats")
		Meta("openapi:extension:x-speakeasy-name-override", "list")
		Meta("openapi:extension:x-speakeasy-react-hook", `{"name": "ListChats"}`)
	})

	Method("loadChat", func() {
		Description("Load a chat by its ID")

		Payload(func() {
			security.SessionPayload()
			security.ProjectPayload()
			Attribute("id", String, "The ID of the chat")
			Required("id")
		})

		Result(Chat)

		HTTP(func() {
			GET("/rpc/chat.load")
			Param("id")
			security.SessionHeader()
			security.ProjectHeader()
			Response(StatusOK)
		})

		Meta("openapi:operationId", "loadChat")
		Meta("openapi:extension:x-speakeasy-name-override", "load")
		Meta("openapi:extension:x-speakeasy-react-hook", `{"name": "LoadChat"}`)
	})
})

var ListChatsResult = Type("ListChatsResult", func() {
	Attribute("chats", ArrayOf(ChatOverview), "The list of chats")
	Required("chats")
})

var ChatOverview = Type("ChatOverview", func() {
	Attribute("id", String, "The ID of the chat")
	Attribute("title", String, "The title of the chat")
	Attribute("user_id", String, "The ID of the user who created the chat")
	Attribute("num_messages", Int, "The number of messages in the chat")
	Attribute("created_at", String, func() {
		Description("When the chat was created.")
		Format(FormatDateTime)
	})
	Attribute("updated_at", String, func() {
		Description("When the chat was last updated.")
		Format(FormatDateTime)
	})

	Required("id", "title", "user_id", "num_messages", "created_at", "updated_at")
})

var Chat = Type("Chat", func() {
	Extend(ChatOverview)
	Attribute("messages", ArrayOf(ChatMessage), "The list of messages in the chat")

	Required("messages")
})

var ChatMessage = Type("ChatMessage", func() {
	Attribute("id", String, "The ID of the message")
	Attribute("role", String, "The role of the message")
	Attribute("content", String, "The content of the message")
	Attribute("model", String, "The model that generated the message")
	Attribute("tool_call_id", String, "The tool call ID of the message")
	Attribute("tool_calls", String, "The tool calls in the message as a JSON blob")
	Attribute("finish_reason", String, "The finish reason of the message")
	Attribute("user_id", String, "The ID of the user who created the message")
	Attribute("created_at", String, func() {
		Description("When the message was created.")
		Format(FormatDateTime)
	})

	Required("id", "role", "model", "created_at")
})
