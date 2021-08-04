SET statement_timeout = 60000;
set lock_timeout = 60000;

ALTER TABLE bookings
    ADD COLUMN canceled BOOLEAN NOT NULL DEFAULT FALSE;
