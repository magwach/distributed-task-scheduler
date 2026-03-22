CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TYPE task_execution_status AS ENUM ('running', 'success', 'failed');
CREATE TABLE task_excecutions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL,
    status task_execution_status DEFAULT 'running',
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    error_message VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
)