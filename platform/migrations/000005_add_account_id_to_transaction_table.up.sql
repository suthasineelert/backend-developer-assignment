ALTER TABLE transactions
ADD COLUMN account_id varchar(50) NOT NULL REFERENCES account (id);