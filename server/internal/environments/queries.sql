-- name: CreateEnvironment :one
INSERT INTO environments (
    organization_id,
    project_id,
    name,
    slug,
    description
) VALUES (
    @organization_id,
    @project_id,
    @name,
    @slug,
    @description
) RETURNING *;

-- name: ListEnvironments :many
SELECT *
FROM environments e
WHERE e.project_id = $1 AND e.deleted IS FALSE
ORDER BY e.created_at DESC;

-- name: GetEnvironmentBySlug :one
-- returns: GetEnvironmentByIDRow
SELECT *
FROM environments e
WHERE e.slug = $1 AND e.project_id = $2 AND e.deleted IS FALSE;

-- name: GetEnvironmentByID :one
SELECT *
FROM environments e
WHERE e.id = $1 AND e.project_id = $2 AND e.deleted IS FALSE;


-- name: UpdateEnvironment :one
UPDATE environments
SET 
    name = COALESCE(@name, name),
    description = COALESCE(@description, description),
    updated_at = now()
WHERE slug = @slug AND project_id = @project_id AND deleted IS FALSE
RETURNING *;

-- name: ListEnvironmentEntries :many
SELECT *
FROM environment_entries ee
WHERE ee.environment_id = $1
ORDER BY ee.name ASC;

-- name: DeleteEnvironment :exec
WITH deleted_env AS (
    UPDATE environments
    SET deleted_at = now()
    WHERE environments.slug = $1 AND environments.project_id = $2 AND environments.deleted IS FALSE
    RETURNING slug, project_id
)
UPDATE toolsets
SET default_environment_slug = NULL
FROM deleted_env
WHERE toolsets.default_environment_slug = deleted_env.slug AND toolsets.project_id = deleted_env.project_id;

-- name: CreateEnvironmentEntries :many
INSERT INTO environment_entries (
    environment_id,
    name,
    value
) 
/*
 Parameters:
 - environment_id: uuid
 - names: text[]
 - values: text[]
*/
VALUES (
    @environment_id::uuid,
    unnest(@names::text[]),
    unnest(@values::text[])
)
RETURNING *;

-- name: UpsertEnvironmentEntry :one
INSERT INTO environment_entries (environment_id, name, value, updated_at)
VALUES ($1, $2, $3, now())
ON CONFLICT (environment_id, name) 
DO UPDATE SET 
    value = EXCLUDED.value,
    updated_at = now()
RETURNING *;

-- name: DeleteEnvironmentEntry :exec
DELETE FROM environment_entries
WHERE environment_id = $1 AND name = $2;
