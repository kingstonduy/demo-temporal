CREATE TABLE limit_manage (
    account_id VARCHAR(255) PRIMARY KEY,
    amount BIGINT NOT NULL
);

INSERT INTO limit_manage (account_id, amount)
VALUES
('OCB12345',  100000),
('OCB00001', 100000),
('OCB00002', 100000),
('OCB00003', 100000),
('OCB00004', 100000),
