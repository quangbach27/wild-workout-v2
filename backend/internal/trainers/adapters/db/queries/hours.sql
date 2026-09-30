-- name: GetHours :many
SELECT *
FROM trainers.hours
WHERE trainer_uuid = @trainer_uuid
    AND hour = ANY(@hours::timestamptz[]);

-- name: UpsertHours :exec
INSERT INTO trainers.hours (trainer_uuid, hour, status)
SELECT
    @trainer_uuid::text,
    unnest(@hours::timestamptz[]),
    unnest(@status::text[])::trainers.hour_status
ON CONFLICT (trainer_uuid, hour)
DO UPDATE
SET status = EXCLUDED.status;
