-- +goose Up
CREATE TYPE user_role AS ENUM('admin', 'customer');
CREATE TYPE gender_type AS ENUM('male', 'female');
CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR UNIQUE NOT NULL,
    phone_number VARCHAR UNIQUE NOT NULL,
    password_hash VARCHAR NOT NULL,
    role user_role NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR,
    date_of_birth DATE NOT NULL,
    gender gender_type NOT NULL,
    avatar_url VARCHAR,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS gender_type;
DROP TYPE IF EXISTS user_role;
