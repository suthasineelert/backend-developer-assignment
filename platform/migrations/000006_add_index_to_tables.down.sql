-- Remove indexes from account_balances table
ALTER TABLE `account_balances`
DROP INDEX `idx_account_balances_user_id`;

-- Remove indexes from account_details table
ALTER TABLE `account_details`
DROP INDEX `idx_account_details_user_id`;

-- Remove indexes from account_flags table
ALTER TABLE `account_flags`
DROP INDEX `idx_account_flags_account_id`;

ALTER TABLE `account_flags`
DROP INDEX `idx_account_flags_user_id`;

ALTER TABLE `account_flags`
DROP INDEX `idx_account_flags_account_user`;

-- Remove indexes from accounts table
ALTER TABLE `accounts`
DROP INDEX `idx_accounts_user_id`;

-- Remove indexes from banners table
ALTER TABLE `banners`
DROP INDEX `idx_banners_user_id`;

-- Remove indexes from debit_card_design table
ALTER TABLE `debit_card_design`
DROP INDEX `idx_debit_card_design_user_id`;

-- Remove indexes from debit_card_details table
ALTER TABLE `debit_card_details`
DROP INDEX `idx_debit_card_details_user_id`;

-- Remove indexes from debit_card_status table
ALTER TABLE `debit_card_status`
DROP INDEX `idx_debit_card_status_user_id`;

-- Remove indexes from debit_cards table
ALTER TABLE `debit_cards`
DROP INDEX `idx_debit_cards_user_id`;

-- Remove indexes from transactions table
ALTER TABLE `transactions`
DROP INDEX `idx_transactions_user_id`;