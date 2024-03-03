CREATE TABLE napas (
    account_id VARCHAR(255) PRIMARY KEY,
    account_name VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL
);

INSERT INTO napas (account_id, account_name, amount)
VALUES
('TMCP23456', 'Tran Thi B', 0),
('BIDV34567', 'Le Chi C', 0),
('STB45678', 'Pham Thi D', 0),
('TCB56789', 'Do Van E', 0),
('MB67890', 'Nguyen Thi G', 0),
('VP78901', 'Tran Van H', 0),
('VCB89012', 'Le Chi I', 0),
('VCB90123', 'Pham Thi J', 0),
('HSBC01234', 'Do Van K', 0),
('TMCP56789', 'Tran Van M', 0),
('BIDV67890', 'Le Chi N', 0),
('STB78901', 'Pham Thi O', 0),
('TCB89012', 'Do Van P', 0),
('MB90123', 'Nguyen Thi Q', 0),
('VP01234', 'Tran Van R', 0),
('VCB12345', 'Le Chi S', 0),
('VCB23456', 'Pham Thi T', 0),
('HSBC34567', 'Do Van U', 0);