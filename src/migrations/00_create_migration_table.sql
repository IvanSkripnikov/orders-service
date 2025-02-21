CREATE TABLE IF NOT EXISTS migration (
    version VARCHAR(255) PRIMARY KEY NOT NULL,
    apply_time BIGINT
);