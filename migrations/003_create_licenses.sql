CREATE TABLE licenses (
    id VARCHAR(36) PRIMARY KEY,
    publication_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    passphrase TEXT NOT NULL,
    hint TEXT NOT NULL,
    publication_url TEXT NOT NULL,
    right_print INTEGER,
    right_copy INTEGER,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (publication_id) REFERENCES publications(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);