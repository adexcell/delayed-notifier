CREATE TABLE IF NOT EXISTS notify (
    notify_id UUID primary key,
    payload JSONB not null,
    target varchar(255) not null,
    channel varchar(100) not null,
    status int not null default 0,
    scheduled_at timestamp with time zone not null,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone,
    retry_count int default 0,
    last_error text
);

CREATE INDEX IF NOT EXISTS idx_notify_status_scheduled ON notify(status, scheduled_at)
where status = 0;

CREATE INDEX IF NOT EXISTS idx_notify_created_at ON notify(created_at DESC);
