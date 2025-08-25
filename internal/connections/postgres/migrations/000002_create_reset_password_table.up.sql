CREATE TABLE reset_passwords (
    token UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_token UUID NOT NULL,
    hash TEXT NOT NULL,
    valid_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_users FOREIGN KEY(user_token) REFERENCES users(token) ON DELETE CASCADE
);
