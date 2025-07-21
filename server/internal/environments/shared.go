package environments

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/speakeasy-api/gram/server/internal/encryption"
	"github.com/speakeasy-api/gram/server/internal/environments/repo"
	"github.com/speakeasy-api/gram/server/internal/gateway"
)

// EnvironmentEntries should be directly accessed through this interface to handle encryption and redaction.
type EnvironmentEntries struct {
	logger *slog.Logger
	repo   *repo.Queries
	enc    *encryption.Client
}

func NewEnvironmentEntries(logger *slog.Logger, db repo.DBTX, enc *encryption.Client) *EnvironmentEntries {
	return &EnvironmentEntries{
		logger: logger,
		repo:   repo.New(db),
		enc:    enc,
	}
}

func (e *EnvironmentEntries) Load(ctx context.Context, projectID uuid.UUID, envIDOrSlug gateway.SlugOrID) (map[string]string, error) {
	environmentID := envIDOrSlug.ID
	if envIDOrSlug.IsEmpty() {
		return nil, fmt.Errorf("environment id or slug is required")
	}

	if environmentID == uuid.Nil {
		envModel, err := e.repo.GetEnvironmentBySlug(ctx, repo.GetEnvironmentBySlugParams{
			ProjectID: projectID,
			Slug:      strings.ToLower(envIDOrSlug.Slug),
		})
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, gateway.ErrNotFound
		case err != nil:
			return nil, fmt.Errorf("get environment by slug: %w", err)
		}

		environmentID = envModel.ID
	}

	if environmentID == uuid.Nil {
		return nil, fmt.Errorf("environment not found for slug or id: %s", envIDOrSlug)
	}

	entries, err := e.ListEnvironmentEntries(ctx, projectID, environmentID, false)
	if err != nil {
		return nil, fmt.Errorf("list environment entries: %w", err)
	}

	envMap := make(map[string]string, len(entries))
	for _, entry := range entries {
		envMap[entry.Name] = entry.Value
	}
	return envMap, nil
}

func (e *EnvironmentEntries) ListEnvironmentEntries(ctx context.Context, projectID uuid.UUID, environmentID uuid.UUID, redacted bool) ([]repo.EnvironmentEntry, error) {
	entries, err := e.repo.ListEnvironmentEntries(ctx, repo.ListEnvironmentEntriesParams{
		ProjectID:     projectID,
		EnvironmentID: environmentID,
	})
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	decryptedEntries := make([]repo.EnvironmentEntry, len(entries))
	for i, entry := range entries {
		value, err := e.enc.Decrypt(entry.Value)
		if err != nil {
			return nil, fmt.Errorf("decrypt entry %s: %w", entry.Name, err)
		}

		if redacted {
			value = redactedEnvironment(value)
		}

		decryptedEntries[i] = repo.EnvironmentEntry{
			Name:          entry.Name,
			Value:         value,
			EnvironmentID: entry.EnvironmentID,
			CreatedAt:     entry.CreatedAt,
			UpdatedAt:     entry.UpdatedAt,
		}
	}

	return decryptedEntries, nil
}

func (e *EnvironmentEntries) CreateEnvironmentEntries(ctx context.Context, params repo.CreateEnvironmentEntriesParams) ([]repo.EnvironmentEntry, error) {
	encryptedValues := make([]string, len(params.Values))
	originalValues := make(map[string]string, len(params.Values))

	for i, value := range params.Values {
		encryptedValue, err := e.enc.Encrypt([]byte(value))
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt value for entry %s: %w", params.Names[i], err)
		}
		encryptedValues[i] = encryptedValue
		originalValues[params.Names[i]] = value // avoid having to needlessly decrypt the value
	}

	params.Values = encryptedValues
	createdEntries, err := e.repo.CreateEnvironmentEntries(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create environment entries: %w", err)
	}

	decryptedEntries := make([]repo.EnvironmentEntry, len(createdEntries))
	for i, entry := range createdEntries {
		decryptedEntries[i] = repo.EnvironmentEntry{
			Name:          entry.Name,
			Value:         redactedEnvironment(originalValues[entry.Name]), // avoid having to needlessly decrypt the value
			EnvironmentID: entry.EnvironmentID,
			CreatedAt:     entry.CreatedAt,
			UpdatedAt:     entry.UpdatedAt,
		}
	}

	return decryptedEntries, nil
}

func (e *EnvironmentEntries) UpdateEnvironmentEntry(ctx context.Context, params repo.UpsertEnvironmentEntryParams) error {
	encryptedValue, err := e.enc.Encrypt([]byte(params.Value))
	if err != nil {
		return fmt.Errorf("failed to encrypt value: %w", err)
	}

	params.Value = encryptedValue
	_, err = e.repo.UpsertEnvironmentEntry(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update environment entry: %w", err)
	}

	return nil
}

func (e *EnvironmentEntries) DeleteEnvironmentEntry(ctx context.Context, params repo.DeleteEnvironmentEntryParams) error {
	err := e.repo.DeleteEnvironmentEntry(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to delete environment entry: %w", err)
	}

	return nil
}

func redactedEnvironment(val string) string {
	if val == "" {
		return "<EMPTY>"
	}
	if len(val) <= 3 {
		return strings.Repeat("*", 5)
	}
	return val[:3] + strings.Repeat("*", 5)
}
