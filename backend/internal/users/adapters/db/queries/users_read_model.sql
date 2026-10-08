-- name: GetTrainers :many
SELECT id, display_name, balance
FROM users.users
WHERE role = 'trainer'
ORDER BY display_name, id
LIMIT @page_limit::bigint
OFFSET @page_offset::bigint;

-- name: CountTrainers :one
SELECT count(*)
FROM users.users
WHERE role = 'trainer';
