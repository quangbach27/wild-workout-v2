-- name: UpsertTraining :exec
INSERT INTO trainings.trainings (
    id, hour, notes,
    attendee_id, attendee_username,
    trainer_id, trainer_username,
    proposed_new_time, move_proposed_by_id,
    canceled
)
VALUES (
    @id, @hour, @notes,
    @attendee_id, @attendee_username,
    @trainer_id, @trainer_username,
    @proposed_new_time, @move_proposed_by_id,
    @canceled
)
ON CONFLICT (id)
DO UPDATE
SET hour = EXCLUDED.hour,
    notes = EXCLUDED.notes,
    proposed_new_time = EXCLUDED.proposed_new_time,
    move_proposed_by_id = EXCLUDED.move_proposed_by_id,
    canceled = EXCLUDED.canceled;

-- name: GetTraining :one
SELECT *
FROM trainings.trainings
WHERE id = @id;
