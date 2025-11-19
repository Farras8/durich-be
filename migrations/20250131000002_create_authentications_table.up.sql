CREATE TABLE authentications (
    id VARCHAR(27) PRIMARY KEY,
    user_email TEXT NOT NULL,
    password TEXT NOT NULL,
    encrypted_password TEXT NOT NULL,
    refresh_token_id VARCHAR(27) NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT authentications_user_email_fkey FOREIGN KEY (user_email) REFERENCES users(email) ON DELETE CASCADE
);

