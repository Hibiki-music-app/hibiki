CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username TEXT NOT NULL UNIQUE,
                       display_name TEXT,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE credentials (
                             id SERIAL PRIMARY KEY,
                             user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                             credential_id BYTEA NOT NULL UNIQUE,
                             public_key BYTEA NOT NULL,
                             sign_count BIGINT NOT NULL,
                             attestation_type TEXT,
                             transports TEXT[],
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
                          id SERIAL PRIMARY KEY,
                          session_key TEXT NOT NULL UNIQUE,
                          user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          data JSONB NOT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          expires_at TIMESTAMP NOT NULL
);
