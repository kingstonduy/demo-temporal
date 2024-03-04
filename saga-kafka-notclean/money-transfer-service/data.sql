CREATE TABLE transaction (
    transaction_id VARCHAR(255) PRIMARY KEY,
    from_account_id VARCHAR(255) NOT NULL,
    to_account_id VARCHAR(255) NOT NULL,
    to_account_name VARCHAR(255) NOT NULL,
    message TEXT,
    amount BIGINT NOT NULL,
    timestamp VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL
);