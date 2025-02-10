-- migrate:up
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE bookings (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id int NOT NULL,
    resource_id int NOT NULL,
    start_at timestamptz NOT NULL,
    end_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    EXCLUDE USING gist (
        resource_id WITH =,
        tstzrange(start_at, end_at, '[)') WITH &&
    )
);

-- migrate:down
DROP TABLE IF EXISTS bookings;
DROP EXTENSION IF EXISTS btree_gist;

