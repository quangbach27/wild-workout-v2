BEGIN;

CREATE SCHEMA IF NOT EXISTS users;

CREATE TYPE users.user_role AS ENUM (
    'trainer',
    'attendee'
);

CREATE TABLE users.users (
    id            varchar(255)    PRIMARY KEY,
    display_name  text            NOT NULL,
    balance       integer         NOT NULL DEFAULT 0,
    role          users.user_role NOT NULL
);

CREATE INDEX users_role_display_name_idx ON users.users (role, display_name, id);

COMMIT;
