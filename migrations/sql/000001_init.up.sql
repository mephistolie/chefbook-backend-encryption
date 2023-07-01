CREATE TABLE vault_keys
(
    user_id     uuid PRIMARY KEY NOT NULL UNIQUE,
    public_key  bytea            NOT NULL,
    private_key bytea            NOT NULL
);

CREATE TYPE recipe_key_request_status as ENUM ('approved', 'pending', 'declined');

CREATE TABLE recipe_keys
(
    recipe_id uuid PRIMARY KEY                                  NOT NULL,
    user_id   uuid REFERENCES users (user_id) ON DELETE CASCADE NOT NULL,
    key       bytea                                                      DEFAULT NULL,
    status    recipe_key_request_status                         NOT NULL DEFAULT 'pending',
    UNIQUE (recipe_id, user_id)
);

CREATE TABLE vault_deletions
(
    user_id     uuid REFERENCES users (user_id) ON DELETE CASCADE NOT NULL UNIQUE,
    delete_code VARCHAR(6)                                        NOT NULL,
);

CREATE TABLE inbox
(
    message_id uuid PRIMARY KEY         NOT NULL UNIQUE,
    timestamp  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now():: timestamp
);

CREATE TABLE outbox
(
    message_id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    exchange   VARCHAR(64)                      DEFAULT '',
    type       VARCHAR(64)      NOT NULL,
    body       JSONB            NOT NULL        DEFAULT '{}'::jsonb
);

