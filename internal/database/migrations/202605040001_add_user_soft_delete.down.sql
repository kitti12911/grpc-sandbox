DROP INDEX IF EXISTS users_username_active_key;
DROP INDEX IF EXISTS users_email_active_key;

ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);

ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;
