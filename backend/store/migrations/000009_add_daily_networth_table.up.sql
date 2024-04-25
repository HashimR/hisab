CREATE TABLE daily_net_worth (
     user_id INT NOT NULL,
     date DATE NOT NULL,
     net_worth DECIMAL(10, 3) NOT NULL,
     last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     PRIMARY KEY (user_id, date)
);

ALTER TABLE balances
    ADD COLUMN date DATE DEFAULT (CURDATE());

CREATE INDEX idx_balances_date ON balances (date);

ALTER TABLE accounts
    ADD COLUMN is_active BOOLEAN DEFAULT true;
