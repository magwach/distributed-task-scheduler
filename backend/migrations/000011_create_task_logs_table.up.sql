CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TYPE level_status AS ENUM ('info', 'warning', 'error');
CREATE TABLE task_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    execution_id UUID NOT NULL,
    level level_status NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (execution_id) REFERENCES task_excecutions (id) ON DELETE CASCADE
)