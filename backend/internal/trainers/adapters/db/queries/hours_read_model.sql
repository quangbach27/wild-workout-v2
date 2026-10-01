-- name: GetHoursInRange :many
SELECT *
FROM trainers.hours
WHERE trainer_uuid = @trainer_uuid
    AND hour >= @date_from::timestamptz
    AND hour < @date_to::timestamptz
ORDER BY hour;
