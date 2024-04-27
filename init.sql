CREATE TABLE IF NOT EXISTS config (
  personal_deduction DECIMAL,
  max_k_receipt DECIMAL
);

INSERT INTO config (personal_deduction, max_k_receipt) VALUES (60000, 50000);