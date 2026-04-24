CREATE TYPE IF NOT EXISTS notification_status AS ENUM ('pending', 'processing', 'sent', 'failed', 'cancelled');
CREATE TYPE IF NOT EXISTS notification_channel AS ENUM ('email', 'telegram');
