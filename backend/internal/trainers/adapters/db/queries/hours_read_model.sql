-- name: GetHoursInRange :many
SELECT *
FROM trainers.hours
WHERE trainer_uuid = @trainer_uuid
    AND hour >= @date_from::timestamptz
    AND hour < @date_to::timestamptz
    AND (sqlc.narg('status')::trainers.hour_status IS NULL OR status = sqlc.narg('status'))
ORDER BY hour;
