-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    uuid          UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    first_name    VARCHAR(255)        NOT NULL,
    password_hash VARCHAR(255)        NOT NULL
);

CREATE TABLE IF NOT EXISTS posts
(
    id               SERIAL PRIMARY KEY,
    title            VARCHAR(255) NOT NULL,
    content          TEXT         NOT NULL,
    comments_enabled BOOLEAN      NOT NULL       DEFAULT TRUE,
    author_uuid      UUID         NOT NULL,
    created_at       TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (author_uuid) REFERENCES users (uuid)
);

CREATE TABLE comments
(
    id               SERIAL PRIMARY KEY,
    post_id          INTEGER NOT NULL,
    parent_id        INTEGER NOT NULL            DEFAULT 0,
    author_uuid      UUID    NOT NULL,
    content          TEXT    NOT NULL,
    has_sub_comments BOOLEAN NOT NULL            DEFAULT FALSE,
    created_at       TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (author_uuid) REFERENCES users (uuid)
);