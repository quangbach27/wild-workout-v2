BEGIN;

CREATE SCHEMA IF NOT EXISTS trainings;

CREATE TABLE trainings.trainings (
    id                    uuid                PRIMARY KEY,
    hour                  timestamptz         NOT NULL,
    notes                 text                NOT NULL,

    attendee_id           varchar(255)        NOT NULL,
    attendee_username     text                NOT NULL DEFAULT '',
    trainer_id            varchar(255)        NOT NULL,
    trainer_username      text                NOT NULL DEFAULT '',

    proposed_new_time         timestamptz,
    move_proposed_by_id       varchar(255),

    canceled              boolean             NOT NULL DEFAULT false
);

CREATE INDEX trainings_attendee_id_idx ON trainings.trainings (attendee_id);
CREATE INDEX trainings_trainer_id_idx ON trainings.trainings (trainer_id);

-- a user can't have two active trainings at the same hour; canceled ones free the slot
CREATE UNIQUE INDEX trainings_hour_trainer_id_uniq ON trainings.trainings (hour, trainer_id) WHERE NOT canceled;
CREATE UNIQUE INDEX trainings_hour_attendee_id_uniq ON trainings.trainings (hour, attendee_id) WHERE NOT canceled;

COMMIT;
