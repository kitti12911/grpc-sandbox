CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    email TEXT NOT NULL,
    username TEXT NOT NULL,
    display_name TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT users_email_key UNIQUE (email),
    CONSTRAINT users_username_key UNIQUE (username),
    CONSTRAINT users_status_check CHECK (status IN ('active', 'disabled', 'pending'))
);

CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL UNIQUE,
    first_name TEXT,
    last_name TEXT,
    phone_number TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX user_profiles_user_id_idx ON user_profiles(user_id);

CREATE TABLE user_addresses (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_profile_id UUID NOT NULL UNIQUE,
    line1 TEXT,
    line2 TEXT,
    city TEXT,
    state TEXT,
    postal_code TEXT,
    country_code TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX user_addresses_user_profile_id_idx ON user_addresses(user_profile_id);
