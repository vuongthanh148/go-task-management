CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the enum type for task status
CREATE TYPE task_status AS ENUM ('Pending', 'In Progress', 'Completed');

CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    assignee_id UUID,
    status task_status NOT NULL DEFAULT 'Pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    due_date TIMESTAMP
);