CREATE INDEX IF NOT EXISTS idx_notifications_pending_scheduled
    ON notifications(scheduled_at)
    WHERE status = 'pending';

CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
CREATE INDEX IF NOT EXISTS idx_notifications_channel ON notifications(channel);
