create table if not exists notify (
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

Create Index idx_notify_status_scheduled on notify(status, scheduled_at)
where status = 0;
