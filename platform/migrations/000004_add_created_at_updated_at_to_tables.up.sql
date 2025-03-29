-- Add timestamp columns to account_balances
ALTER TABLE `account_balances`
ADD COLUMN IF NOT EXISTS `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `deleted_at` timestamp NULL DEFAULT NULL;

-- Add timestamp columns to account_details
ALTER TABLE `account_details`
ADD COLUMN IF NOT EXISTS `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `deleted_at` timestamp NULL DEFAULT NULL;

-- account_flags already has created_at and updated_at, just add deleted_at
ALTER TABLE `account_flags`
ADD COLUMN IF NOT EXISTS `deleted_at` timestamp NULL DEFAULT NULL;

-- Add timestamp columns to accounts
ALTER TABLE `accounts`
ADD COLUMN IF NOT EXISTS `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `deleted_at` timestamp NULL DEFAULT NULL;

-- Add timestamp columns to banners
ALTER TABLE `banners`
ADD COLUMN IF NOT EXISTS `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `deleted_at` timestamp NULL DEFAULT NULL;

-- Add timestamp columns to debit_card_design
ALTER TABLE `debit_card_design`
ADD COLUMN IF NOT EXISTS `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `deleted_at` timestamp NULL DEFAULT NULL;

-- Add timestamp columns to debit_card_details
ALTER TABLE `debit_card_details`
ADD COLUMN IF NOT EXISTS `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `deleted_at` timestamp NULL DEFAULT NULL;

-- Add timestamp columns to debit_card_status
ALTER TABLE `debit_card_status`
ADD COLUMN IF NOT EXISTS `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `deleted_at` timestamp NULL DEFAULT NULL;

-- Add timestamp columns to debit_cards
ALTER TABLE `debit_cards`
ADD COLUMN IF NOT EXISTS `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `deleted_at` timestamp NULL DEFAULT NULL;

-- Add timestamp columns to transactions
ALTER TABLE `transactions`
ADD COLUMN IF NOT EXISTS `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `deleted_at` timestamp NULL DEFAULT NULL;

-- Add timestamp columns to user_greetings
ALTER TABLE `user_greetings`
ADD COLUMN IF NOT EXISTS `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `deleted_at` timestamp NULL DEFAULT NULL;

-- Add timestamp columns to users
ALTER TABLE `users`
ADD COLUMN IF NOT EXISTS `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS `deleted_at` timestamp NULL DEFAULT NULL;