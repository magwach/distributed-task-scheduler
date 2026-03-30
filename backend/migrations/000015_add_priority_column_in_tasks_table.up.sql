CREATE TYPE priority_enum AS ENUM ('low', 'normal', 'high');
ALTER TABLE tasks ADD COLUMN priority priority_enum NOT NULL DEFAULT 'normal';