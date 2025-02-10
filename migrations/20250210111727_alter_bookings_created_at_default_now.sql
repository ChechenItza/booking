-- migrate:up
ALTER TABLE bookings ALTER COLUMN created_at SET DEFAULT now();

-- migrate:down
ALTER TABLE bookings ALTER COLUMN created_at DROP DEFAULT;