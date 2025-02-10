-- migrate:up
ALTER TABLE bookings ALTER COLUMN updated_at SET DEFAULT now();
ALTER TABLE booking_count ALTER COLUMN updated_at SET DEFAULT now();

-- migrate:down
ALTER TABLE bookings ALTER COLUMN updated_at DROP DEFAULT;
ALTER TABLE booking_count ALTER COLUMN updated_at DROP DEFAULT;
