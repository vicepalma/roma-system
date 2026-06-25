-- Phase 0 alignment: keep official migrations compatible with current backend code.

CREATE TABLE IF NOT EXISTS coach_links (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  coach_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  disciple_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  status      TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','accepted','rejected')),
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (coach_id, disciple_id)
);

CREATE INDEX IF NOT EXISTS idx_coach_links_coach ON coach_links(coach_id);
CREATE INDEX IF NOT EXISTS idx_coach_links_disciple ON coach_links(disciple_id);

INSERT INTO coach_links (coach_id, disciple_id, status, created_at)
SELECT
  master_id,
  disciple_id,
  CASE
    WHEN status = 'active' THEN 'accepted'
    WHEN status = 'ended' THEN 'rejected'
    ELSE 'pending'
  END,
  created_at
FROM master_disciple
ON CONFLICT (coach_id, disciple_id) DO NOTHING;

CREATE TABLE IF NOT EXISTS program_versions (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  program_id UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
  version    INT NOT NULL CHECK (version >= 1),
  title      TEXT NOT NULL,
  notes      TEXT NULL,
  created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM now())::BIGINT,
  UNIQUE (program_id, version)
);

CREATE INDEX IF NOT EXISTS idx_program_versions_program ON program_versions(program_id);

INSERT INTO program_versions (program_id, version, title, notes, created_at)
SELECT id, version, title, notes, EXTRACT(EPOCH FROM created_at)::BIGINT
FROM programs
ON CONFLICT (program_id, version) DO NOTHING;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger
    WHERE tgname = 'trg_coach_links_updated'
  ) THEN
    CREATE TRIGGER trg_coach_links_updated
    BEFORE UPDATE ON coach_links
    FOR EACH ROW EXECUTE PROCEDURE set_updated_at();
  END IF;
END;
$$;
