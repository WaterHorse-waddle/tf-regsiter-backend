-- ==============================================================================
-- Announcements table for user-facing hackathon updates
-- ==============================================================================

CREATE TABLE IF NOT EXISTS announcements (
    announcement_id uuid PRIMARY KEY,
    title text NOT NULL,
    body text NOT NULL,
    created_by_email text,
    created_by_name text,
    is_active boolean DEFAULT true,
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_announcements_is_active_created_at
    ON announcements(is_active, created_at DESC);