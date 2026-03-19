-- 002_create_campaigns.sql
CREATE TABLE campaigns (
    id          UUID         PRIMARY KEY,
    user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ  NULL
);

CREATE INDEX idx_campaigns_user_id ON campaigns(user_id);

---- create above / drop below ----

DROP INDEX idx_campaigns_user_id;
DROP TABLE campaigns;
