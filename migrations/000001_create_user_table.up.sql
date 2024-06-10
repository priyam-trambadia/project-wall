CREATE TABLE IF NOT EXISTS users (
  _id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255)NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  refreshToken VARCHAR(255),
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);