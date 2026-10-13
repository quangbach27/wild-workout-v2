-- name: InsertUser :exec
INSERT INTO users.users (id, display_name, balance, role)
VALUES (@id, @display_name, @balance, @role);

-- name: GetUserByUUID :one
SELECT id, display_name, balance, role
FROM users.users
WHERE id = @id;

-- name: UpdateBalance :one
UPDATE users.users
SET balance = balance + @amount_change::integer
WHERE id = @id AND balance + @amount_change::integer >= 0
RETURNING balance;

-- name: UpdateUser :execrows
UPDATE users.users
SET display_name = @display_name, role = @role
WHERE id = @id;

-- name: UpdateDisplayName :execrows
UPDATE users.users
SET display_name = @display_name
WHERE id = @id;
