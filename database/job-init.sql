
CREATE TABLE jobs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    job_name TEXT NOT NULL,
    image_name TEXT NOT NULL,
    image_version TEXT NOT NULL,
    adjustment_parameters JSONB NOT NULL,
    creation_zone TEXT NOT NULL,
    worker_id TEXT,
    compute_zone TEXT,
    carbon_intensity INTEGER DEFAULT -1,
    carbon_savings INTEGER DEFAULT -1,
    result TEXT DEFAULT '',
    error_message TEXT DEFAULT '',
    job_status TEXT DEFAULT 'queued'
);
