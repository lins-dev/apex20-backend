-- 002_create_campaigns.sql
CREATE TABLE campaigns (
    id UUID PRIMARY KEY,
    users_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for foreign key lookups
CREATE INDEX idx_campaigns_users_id ON campaigns(users_id);

---- create above / drop below ----

DROP TABLE campaigns;
