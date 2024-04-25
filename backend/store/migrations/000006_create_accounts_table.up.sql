CREATE TABLE accounts (
  id INT AUTO_INCREMENT PRIMARY KEY,
  dapi_id VARCHAR(255) NOT NULL,
  user_id INT NOT NULL,
  account_type VARCHAR(50) NOT NULL,
  iban VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  number VARCHAR(50) NOT NULL,
  currency_code VARCHAR(3) NOT NULL,
  last_updated_transactions TIMESTAMP NULL DEFAULT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE balances (
  id INT AUTO_INCREMENT PRIMARY KEY,
  account_id INT NOT NULL,
  user_id INT NOT NULL,
  amount DECIMAL(10, 3) NOT NULL,
  currency VARCHAR(3) NOT NULL,
  last_updated TIMESTAMP NULL DEFAULT NULL
);
