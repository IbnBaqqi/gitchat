-- name: CreateUser :one
INSERT INTO users (
    github_id,
    username,
    display_name,
    email,
    avatar_url,
    access_token,
    has_account
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByGithubID :one
SELECT * FROM users
WHERE github_id = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1
LIMIT 1;

-- name: UpdateUserProfile :one
UPDATE users
SET
    display_name = $2,
    email = $3,
    avatar_url = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateLastSeen :exec
UPDATE users
SET last_seen_at = NOW()
WHERE id = $1;

-- name: SetOnlinePrivacy :exec
UPDATE users
SET hide_online_status = $2
WHERE id = $1;

-- name: UpdateAccessToken :exec
UPDATE users
SET access_token = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: SetHasAccount :exec
UPDATE users
SET has_account = TRUE
WHERE id = $1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;