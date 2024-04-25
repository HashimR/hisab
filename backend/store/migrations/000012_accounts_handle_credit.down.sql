ALTER TABLE accounts
    DROP COLUMN is_credit,
    DROP COLUMN credit_card_number_last_four,
    DROP COLUMN credit_card_limit,
    DROP COLUMN credit_card_next_payment_date,
    DROP COLUMN credit_card_next_payment_amount;