CREATE TYPE gender_data_type AS ENUM ('male', 'female', 'none');

CREATE TABLE IF NOT EXISTS accounts (
  id SERIAL PRIMARY KEY,
  username VARCHAR(60) UNIQUE,
  full_name VARCHAR(150),
  email VARCHAR(150) UNIQUE,
  password VARCHAR(64),
  address VARCHAR(100),
  employee_number VARCHAR(50),
  job_position VARCHAR(50),
  ktp_number VARCHAR(50) UNIQUE,
  phone_number VARCHAR(20) UNIQUE,
  photo_url VARCHAR(500),
  gender gender_data_type NOT NULL DEFAULT 'none',
  date_of_birth DATE,
  is_verified BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);
