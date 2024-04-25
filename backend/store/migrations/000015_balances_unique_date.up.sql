ALTER TABLE balances
    ADD CONSTRAINT unique_account_date UNIQUE (account_id, date);