BEGIN;

CREATE SCHEMA IF NOT EXISTS trainers;

CREATE TYPE trainers.hour_status AS ENUM (
    'availability',
    'not-availability',
    'training-scheduled'
);

CREATE TABLE trainers.hours (
    trainer_uuid varchar(255)         NOT NULL,
    hour         timestamptz          NOT NULL,
    status       trainers.hour_status NOT NULL,

    PRIMARY KEY (trainer_uuid, hour)
);

COMMIT;
