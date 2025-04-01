-- Add indexes to account_balances table
ALTER TABLE `account_balances` ADD INDEX `idx_account_balances_user_id` (`user_id`);

-- Add indexes to account_details table
ALTER TABLE `account_details` ADD INDEX `idx_account_details_user_id` (`user_id`);

-- Add indexes to account_flags table
ALTER TABLE `account_flags` ADD INDEX `idx_account_flags_account_id` (`account_id`);

ALTER TABLE `account_flags` ADD INDEX `idx_account_flags_user_id` (`user_id`);

ALTER TABLE `account_flags` ADD INDEX `idx_account_flags_account_user` (`account_id`, `user_id`);

-- Add indexes to accounts table
ALTER TABLE `accounts` ADD INDEX `idx_accounts_user_id` (`user_id`);

-- Add indexes to banners table
ALTER TABLE `banners` ADD INDEX `idx_banners_user_id` (`user_id`);

-- Add indexes to debit_card_design table
ALTER TABLE `debit_card_design` ADD INDEX `idx_debit_card_design_user_id` (`user_id`);

-- Add indexes to debit_card_details table
ALTER TABLE `debit_card_details` ADD INDEX `idx_debit_card_details_user_id` (`user_id`);

-- Add indexes to debit_card_status table
ALTER TABLE `debit_card_status` ADD INDEX `idx_debit_card_status_user_id` (`user_id`);

-- Add indexes to debit_cards table
ALTER TABLE `debit_cards` ADD INDEX `idx_debit_cards_user_id` (`user_id`);

-- Add indexes to transactions table
ALTER TABLE `transactions` ADD INDEX `idx_transactions_user_id` (`user_id`);