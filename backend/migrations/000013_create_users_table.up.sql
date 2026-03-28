CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TYPE role_enum AS ENUM ('admin', 'viewer');
CREATE TYPE provider_enum AS ENUM ('local', 'google', 'github');
CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT,
    role role_enum DEFAULT 'viewer',
    provider provider_enum,
    provider_id VARCHAR(255),
    avatar_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
)