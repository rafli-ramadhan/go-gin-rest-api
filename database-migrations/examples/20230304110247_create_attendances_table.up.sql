CREATE TYPE attendance_status AS ENUM ('check-in', 'check-out', 'none');

CREATE TABLE IF NOT EXISTS attendances (
  account_id INT NOT NULL REFERENCES "accounts" ON UPDATE CASCADE ON DELETE CASCADE,
  location_id INT NOT NULL REFERENCES "locations" ON UPDATE CASCADE ON DELETE CASCADE,
  status attendance_status NOT NULL DEFAULT 'none',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);
