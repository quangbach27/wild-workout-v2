BEGIN;

CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE users.users (
    id            varchar(255)  PRIMARY KEY,
    display_name  text          NOT NULL,
    balance       integer       NOT NULL DEFAULT 0,
    role          text          NOT NULL
);

CREATE INDEX users_role_display_name_idx ON users.users (role, display_name, id);

COMMIT;
