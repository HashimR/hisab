ALTER TABLE balances DROP COLUMN IF EXISTS date;
DROP INDEX IF EXISTS idx_balances_date ON balances;

ALTER TABLE accounts DROP COLUMN IF EXISTS is_active;

DROP TABLE IF EXISTS daily_net_worth;
