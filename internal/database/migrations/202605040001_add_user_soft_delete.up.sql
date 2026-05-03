ALTER TABLE users ADD COLUMN deleted_at TIMESTAMPTZ;

ALTER TABLE users DROP CONSTRAINT users_email_key;
ALTER TABLE users DROP CONSTRAINT users_username_key;

CREATE UNIQUE INDEX users_email_active_key ON users(email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX users_username_active_key ON users(username) WHERE deleted_at IS NULL;
