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



ALTER TABLE users
ADD COLUMN oauth_provider VARCHAR(50),        -- e.g. 'google', 'facebook'
ADD COLUMN oauth_id TEXT,                     -- unique provider user ID
ADD CONSTRAINT unique_oauth UNIQUE (oauth_provider, oauth_id);
ALTER TABLE users
ALTER COLUMN password_hash DROP NOT NULL,
ALTER COLUMN salt DROP NOT NULL;

ALTER TABLE auth_methods
DROP COLUMN secret_data;

CREATE OR REPLACE FUNCTION create_auth_method_for_user()
RETURNS TRIGGER AS $$
BEGIN
    -- If OAuth ID is not null, set method_type = 'oauth'
    IF NEW.oauth_id IS NOT NULL THEN
        INSERT INTO auth_methods(user_id, method_type, last_used)
        VALUES (NEW.id, 'oauth', NOW());
    ELSE
        -- Otherwise, assume password signup
        INSERT INTO auth_methods(user_id, method_type, last_used)
        VALUES (NEW.id, 'password', NOW());
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_auth_method
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION create_auth_method_for_user();


CREATE OR REPLACE FUNCTION update_auth_method_last_login(p_username TEXT)
RETURNS VOID AS $$
DECLARE
    v_user_id UUID;
BEGIN
    -- Get the user ID from users table
    SELECT id INTO v_user_id
    FROM users
    WHERE username = p_username;

    -- Only update if user exists
    IF v_user_id IS NOT NULL THEN
        UPDATE auth_methods
        SET last_used = NOW()
        WHERE user_id = v_user_id;
    END IF;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION log_user_action(
    p_username TEXT,
    p_action TEXT,
    p_ip_address TEXT,
    p_user_agent TEXT,
    p_metadata TEXT
)
RETURNS VOID AS $$
DECLARE 
    v_user_id UUID;
BEGIN
    -- Find the user's ID
    SELECT id INTO v_user_id
    FROM users
    WHERE username = p_username;

    -- Insert into audit_logs
    INSERT INTO audit_logs(user_id, action, ip_address, user_agent, metadata)
    VALUES (v_user_id, p_action, p_ip_address, p_user_agent, p_metadata::jsonb);
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER trg_log_user_registration
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION log_user_registration();


ALTER TABLE users
    DROP CONSTRAINT users_account_status_check,
    ADD CONSTRAINT users_account_status_check
    CHECK (account_status IN ('active', 'pending', 'disabled'));



CREATE OR REPLACE FUNCTION retrieve_user_audits(p_email TEXT)
RETURNS TABLE(
    email TEXT,
    action TEXT,
    ip_address TEXT,
    user_agent TEXT,
    metadata TEXT,
    "timestamp" TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT u.email AS email,
           a.action,
           a.ip_address,
           a.user_agent,
           a.metadata::TEXT,
           a.timestamp
    FROM audit_logs a
    JOIN users u ON u.id = a.user_id
    WHERE u.email = p_email
    ORDER BY a.timestamp DESC;
END;
$$ LANGUAGE plpgsql;

