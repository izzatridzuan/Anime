CREATE TABLE audit_logs(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    action_type TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)