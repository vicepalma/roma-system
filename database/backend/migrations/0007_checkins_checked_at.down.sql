DROP INDEX IF EXISTS idx_checkins_disciple_checked_at;

ALTER TABLE checkins
  DROP COLUMN IF EXISTS checked_at;
