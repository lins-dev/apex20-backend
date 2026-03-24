-- 005_create_campaign_members.sql
-- Source of truth for a user's role within a specific campaign (ADR-002).
-- role is SMALLINT mapped via Protobuf Role enum: 1=GM, 2=Player, 3=Trusted
-- Campaign creator is automatically inserted as GM at the application layer.
CREATE TABLE campaign_members (
    id          UUID        PRIMARY KEY,
    campaign_id UUID        NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role         SMALLINT     NOT NULL,
    display_name VARCHAR(255) NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (campaign_id, user_id)
);

CREATE INDEX idx_campaign_members_campaign_id ON campaign_members(campaign_id);
CREATE INDEX idx_campaign_members_user_id ON campaign_members(user_id);

---- create above / drop below ----

DROP INDEX idx_campaign_members_user_id;
DROP INDEX idx_campaign_members_campaign_id;
DROP TABLE campaign_members;
