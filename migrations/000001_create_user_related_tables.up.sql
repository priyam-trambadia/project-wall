CREATE TABLE IF NOT EXISTS organizations (
  id BIGSERIAL PRIMARY KEY,
  hostname VARCHAR(255)NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255)NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  avatar TEXT NOT NULL DEFAULT '/static/img/default_avatar.jpg',
  bio TEXT NOT NULL DEFAULT '',
  social_links TEXT[] NOT NULL DEFAULT '{}',
  organization_id BIGINT NOT NULL REFERENCES organizations(id), 
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  refresh_token VARCHAR(255),
  is_activated BOOLEAN NOT NULL DEFAULT false
);