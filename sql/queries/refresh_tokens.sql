-- name: CreateRefreshToken :exec
INSERT INTO
  refresh_tokens (
    token,
    created_at,
    updated_at,
    user_id,
    expires_at
  )
VALUES
  ($1, NOW(), NOW(), $2, $3);

-- name: GetUserFromRefreshToken :one
SELECT
  refresh_tokens.token,
  refresh_tokens.created_at,
  refresh_tokens.updated_at,
  refresh_tokens.user_id,
  refresh_tokens.expires_at,
  refresh_tokens.revoked_at
FROM
  users
  INNER JOIN refresh_tokens ON refresh_tokens.user_id = users.id
WHERE
  refresh_tokens.token = $1
  AND revoked_at IS NULL
  AND expires_at > NOW();

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET
  revoked_at = NOW(),
  updated_at = NOW()
WHERE
  token = $1;
