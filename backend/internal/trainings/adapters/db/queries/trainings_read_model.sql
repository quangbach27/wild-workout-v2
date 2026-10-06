-- name: GetUpcomingTrainingsByUser :many
SELECT *
FROM trainings.trainings
WHERE (attendee_id = @user_id OR trainer_id = @user_id)
    AND hour > @after::timestamptz
ORDER BY hour, id
LIMIT @page_limit::bigint
OFFSET @page_offset::bigint;

-- name: CountUpcomingTrainingsByUser :one
SELECT count(*)
FROM trainings.trainings
WHERE (attendee_id = @user_id OR trainer_id = @user_id)
    AND hour > @after::timestamptz;
