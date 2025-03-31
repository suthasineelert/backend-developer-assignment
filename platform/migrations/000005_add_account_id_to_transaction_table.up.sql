ALTER TABLE transactions
ADD COLUMN account_id INTEGER NOT NULL REFERENCES account (id);