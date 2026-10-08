-- Sample data for local development: the mock users from backend/cmd/main.go. Safe to run repeatedly.
-- Run with `make seed` (root) after `make up`.
BEGIN;

-- Dev-only reset: drop every trainer except the mock one (leftovers from tests
-- and manual onboarding), together with their hours.
DELETE FROM trainers.hours
WHERE trainer_uuid <> '11111111-1111-4111-8111-111111111111';
DELETE FROM users.users
WHERE role = 'trainer' AND id <> '11111111-1111-4111-8111-111111111111';

INSERT INTO users.users (id, display_name, balance, role)
VALUES
    ('11111111-1111-4111-8111-111111111111', 'Mock Trainer', 0, 'trainer'),
    ('22222222-2222-4222-8222-222222222222', 'Mock Attendee', 10, 'attendee')
ON CONFLICT (id) DO NOTHING;

COMMIT;
