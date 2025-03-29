-- Remove timestamp columns from account_balances
ALTER TABLE `account_balances`
DROP COLUMN `created_at`,
DROP COLUMN `updated_at`,
DROP COLUMN `deleted_at`;

-- Remove timestamp columns from account_details
ALTER TABLE `account_details`
DROP COLUMN `created_at`,
DROP COLUMN `updated_at`,
DROP COLUMN `deleted_at`;

-- Remove deleted_at from account_flags (keeping created_at and updated_at)
ALTER TABLE `account_flags`
DROP COLUMN `deleted_at`;

-- Remove timestamp columns from accounts
ALTER TABLE `accounts`
DROP COLUMN `created_at`,
DROP COLUMN `updated_at`,
DROP COLUMN `deleted_at`;

-- Remove timestamp columns from banners
ALTER TABLE `banners`
DROP COLUMN `created_at`,
DROP COLUMN `updated_at`,
DROP COLUMN `deleted_at`;

-- Remove timestamp columns from debit_card_design
ALTER TABLE `debit_card_design`
DROP COLUMN `created_at`,
DROP COLUMN `updated_at`,
DROP COLUMN `deleted_at`;

-- Remove timestamp columns from debit_card_details
ALTER TABLE `debit_card_details`
DROP COLUMN `created_at`,
DROP COLUMN `updated_at`,
DROP COLUMN `deleted_at`;

-- Remove timestamp columns from debit_card_status
ALTER TABLE `debit_card_status`
DROP COLUMN `created_at`,
DROP COLUMN `updated_at`,
DROP COLUMN `deleted_at`;

-- Remove timestamp columns from debit_cards
ALTER TABLE `debit_cards`
DROP COLUMN `created_at`,
DROP COLUMN `updated_at`,
DROP COLUMN `deleted_at`;

-- Remove timestamp columns from transactions
ALTER TABLE `transactions`
DROP COLUMN `created_at`,
DROP COLUMN `updated_at`,
DROP COLUMN `deleted_at`;

-- Remove timestamp columns from user_greetings
ALTER TABLE `user_greetings`
DROP COLUMN `created_at`,
DROP COLUMN `updated_at`,
DROP COLUMN `deleted_at`;

-- Remove timestamp columns from users
ALTER TABLE `users`
DROP COLUMN `created_at`,
DROP COLUMN `updated_at`,
DROP COLUMN `deleted_at`;