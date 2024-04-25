ALTER TABLE accounts
    ADD COLUMN is_credit BOOLEAN,
    ADD COLUMN credit_card_number_last_four VARCHAR(4),
    ADD COLUMN credit_card_limit DECIMAL(10, 3),
    ADD COLUMN credit_card_next_payment_date TIMESTAMP,
    ADD COLUMN credit_card_next_payment_amount DECIMAL(10, 3);