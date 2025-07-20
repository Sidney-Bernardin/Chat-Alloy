CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,

    created_at timestamp,
    updated_at timestamp,

    username text,
    pw_hash bytea,
    pw_salt bytea
);
