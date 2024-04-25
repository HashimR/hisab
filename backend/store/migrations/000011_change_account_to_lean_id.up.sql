ALTER TABLE accounts
    CHANGE COLUMN dapi_id lean_account_id VARCHAR(255) NOT NULL UNIQUE;