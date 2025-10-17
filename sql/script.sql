CREATE DATABASE auth_system;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE USER auth_user WITH PASSWORD 'auth_user_2025';


GRANT CONNECT ON DATABASE auth_system TO auth_user;

GRANT USAGE ON SCHEMA public TO auth_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO auth_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO auth_user;

--Para secuencias 
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO auth_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO auth_user;


--Audit User


CREATE USER auth_admin WITH PASSWORD 'audit_secure_2025';
GRANT ALL PRIVILEGES ON DATABASE auth TO auth_admin;
\c auth
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO auth_admin;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO auth_admin;
CREATE ROLE audit_admin;
GRANT audit_admin TO auth_admin;



CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_number_encrypted BYTEA,                -- Encrypted ID/Cédula
    full_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone_number TEXT UNIQUE,
    date_of_birth DATE,
    address TEXT,
    create_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash BYTEA NOT NULL,
    salt BYTEA NOT NULL,

    is_verified BOOLEAN DEFAULT FALSE,
    last_login TIMESTAMP,
    account_status VARCHAR(20) DEFAULT 'active' CHECK (account_status IN ('active', 'suspended', 'locked'))
);

CREATE TABLE auth_methods (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    method_type VARCHAR(50) NOT NULL,          -- 'password', 'biometric', 'otp', 'passkey', etc.
    secret_data BYTEA,                         -- Encrypted key, device hash, biometric template, etc.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used TIMESTAMP
);
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,               -- 'login', 'logout', 'failed_login', 'lock', etc.
    ip_address INET,
    user_agent TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB                             -- optional, flexible for extra info (e.g. device ID)
);
CREATE TABLE recovery_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash BYTEA NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
