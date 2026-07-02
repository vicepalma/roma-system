ALTER TABLE checkins
  ADD COLUMN IF NOT EXISTS checked_at DATE NOT NULL DEFAULT CURRENT_DATE;

CREATE INDEX IF NOT EXISTS idx_checkins_disciple_checked_at
  ON checkins(disciple_id, checked_at DESC);
