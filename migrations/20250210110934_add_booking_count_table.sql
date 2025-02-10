-- migrate:up
CREATE TABLE booking_count (
    resource_id int PRIMARY KEY,
    count int NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL
);

-- migrate:down

DROP TABLE booking_count;