CREATE TABLE IF NOT EXISTS projects (
  id bigint PRIMARY KEY,
  github_link varchar(255) NOT NULL,
  title varchar(255) NOT NULL,
  description text,
  owner_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_synced_at timestamp,
  bookmark_count bigint NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS tags (
  id bigint PRIMARY KEY,
  name varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS project_tags (
  project_id bigint NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  tag_id bigint NOT NULL REFERENCES tags(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS languages (
  id bigint PRIMARY KEY,
  name varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS project_languages (
  project_id bigint NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  language_id bigint NOT NULL REFERENCES languages(id) ON DELETE CASCADE
);