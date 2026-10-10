-- ==============================================================================
-- Issues table for participant-raised issues during the hackathon
-- ==============================================================================

CREATE TABLE IF NOT EXISTS issues (
    issue_id     uuid PRIMARY KEY,
    team_id      text REFERENCES teams(team_id) ON DELETE SET NULL,
    user_id      uuid REFERENCES users(user_id) ON DELETE SET NULL,
    user_name    text,
    user_email   text,
    issue_type   text NOT NULL,
    title        text NOT NULL,
    description  text NOT NULL,
    status       text DEFAULT 'open',
    created_at   timestamptz DEFAULT now(),
    updated_at   timestamptz DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_issues_status
    ON issues(status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_issues_team_id
    ON issues(team_id);
