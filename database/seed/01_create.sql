BEGIN;

-- Extensiones
CREATE EXTENSION IF NOT EXISTS pgcrypto; -- gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS citext;   -- email case-insensitive

-- 1) Usuarios y relaciones
CREATE TABLE IF NOT EXISTS users (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email          CITEXT NOT NULL UNIQUE,
  password_hash  TEXT   NOT NULL,
  name           TEXT   NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS master_disciple (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  master_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  disciple_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  status       TEXT NOT NULL DEFAULT 'active', -- active|paused|ended
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (master_id, disciple_id)
);
CREATE INDEX IF NOT EXISTS idx_md_master   ON master_disciple(master_id);
CREATE INDEX IF NOT EXISTS idx_md_disciple ON master_disciple(disciple_id);

-- 2) Biblioteca
CREATE TABLE IF NOT EXISTS exercises (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name            TEXT NOT NULL,
  primary_muscle  TEXT NOT NULL,
  equipment       TEXT NULL,
  tags            TEXT[] NOT NULL DEFAULT '{}',
  notes           TEXT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_exercise_name ON exercises (lower(name));

CREATE TABLE IF NOT EXISTS methods (
  id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  key     TEXT NOT NULL UNIQUE,
  name    TEXT NOT NULL,
  params  JSONB NOT NULL DEFAULT '{}'
);

-- 3) Programación
CREATE TABLE IF NOT EXISTS programs (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title       TEXT NOT NULL,
  notes       TEXT NULL,
  visibility  TEXT NOT NULL DEFAULT 'private',
  version     INT  NOT NULL DEFAULT 1,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_programs_owner ON programs(owner_id);

CREATE TABLE IF NOT EXISTS program_weeks (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  program_id  UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
  week_index  INT  NOT NULL CHECK (week_index >= 1)
);
CREATE INDEX IF NOT EXISTS idx_weeks_program ON program_weeks(program_id);

CREATE TABLE IF NOT EXISTS program_days (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  week_id    UUID NOT NULL REFERENCES program_weeks(id) ON DELETE CASCADE,
  day_index  INT  NOT NULL CHECK (day_index >= 1),
  notes      TEXT NULL
);
CREATE INDEX IF NOT EXISTS idx_days_week ON program_days(week_id);

CREATE TABLE IF NOT EXISTS prescriptions (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  day_id        UUID NOT NULL REFERENCES program_days(id) ON DELETE CASCADE,
  exercise_id   UUID NOT NULL REFERENCES exercises(id) ON DELETE RESTRICT,
  series        INT  NOT NULL CHECK (series >= 1),
  reps          TEXT NOT NULL,
  rest_sec      INT  NULL CHECK (rest_sec IS NULL OR rest_sec >= 0),
  to_failure    BOOLEAN NOT NULL DEFAULT FALSE,
  tempo         TEXT NULL,
  rir           INT  NULL CHECK (rir BETWEEN 0 AND 5),
  rpe           NUMERIC(3,1) NULL CHECK (rpe BETWEEN 1 AND 10),
  method_id     UUID NULL REFERENCES methods(id) ON DELETE SET NULL,
  notes         TEXT NULL,
  position      INT  NOT NULL DEFAULT 1
);
CREATE INDEX IF NOT EXISTS idx_presc_day      ON prescriptions(day_id);
CREATE INDEX IF NOT EXISTS idx_presc_exercise ON prescriptions(exercise_id);

-- 4) Asignaciones
CREATE TABLE IF NOT EXISTS assignments (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  program_id       UUID NOT NULL REFERENCES programs(id) ON DELETE RESTRICT,
  program_version  INT  NOT NULL,
  disciple_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  assigned_by      UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  start_date       DATE NOT NULL,
  end_date         DATE NULL,
  is_active        BOOLEAN NOT NULL DEFAULT TRUE,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (start_date <= COALESCE(end_date, start_date))
);
CREATE INDEX IF NOT EXISTS idx_assign_disciple ON assignments(disciple_id);
CREATE INDEX IF NOT EXISTS idx_assign_program  ON assignments(program_id);

-- 5) Ejecución y logs
CREATE TABLE IF NOT EXISTS session_logs (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  assignment_id  UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
  disciple_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  day_id         UUID NOT NULL REFERENCES program_days(id) ON DELETE RESTRICT,
  performed_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  notes          TEXT NULL,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_sess_assign   ON session_logs(assignment_id);
CREATE INDEX IF NOT EXISTS idx_sess_disciple ON session_logs(disciple_id);
CREATE INDEX IF NOT EXISTS idx_session_by_assign_day_created
  ON session_logs (assignment_id, day_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_session_by_assign_day_performed
  ON session_logs (assignment_id, day_id, performed_at DESC);

CREATE TABLE IF NOT EXISTS set_logs (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id      UUID NOT NULL REFERENCES session_logs(id) ON DELETE CASCADE,
  prescription_id UUID NOT NULL REFERENCES prescriptions(id) ON DELETE RESTRICT,
  set_index       INT NOT NULL CHECK (set_index >= 1),
  weight          NUMERIC(8,2) NULL CHECK (weight IS NULL OR weight >= 0),
  reps            INT NOT NULL CHECK (reps >= 0),
  rpe             NUMERIC(3,1) NULL CHECK (rpe BETWEEN 1 AND 10),
  to_failure      BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_sets_session ON set_logs(session_id);

-- 6) Cardio
CREATE TABLE IF NOT EXISTS cardio_segments (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  day_id        UUID NULL REFERENCES program_days(id) ON DELETE CASCADE,
  session_id    UUID NULL REFERENCES session_logs(id) ON DELETE CASCADE,
  modality      TEXT NOT NULL,
  minutes       INT  NOT NULL CHECK (minutes > 0),
  target_hr_min INT  NULL,
  target_hr_max INT  NULL,
  notes         TEXT NULL,
  CHECK ( (day_id IS NOT NULL) <> (session_id IS NOT NULL) )
);
CREATE INDEX IF NOT EXISTS idx_cardio_day     ON cardio_segments(day_id);
CREATE INDEX IF NOT EXISTS idx_cardio_session ON cardio_segments(session_id);

-- 7) Check-ins
CREATE TABLE IF NOT EXISTS checkins (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  disciple_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  weight_kg   NUMERIC(5,2) NULL CHECK (weight_kg IS NULL OR weight_kg > 0),
  notes       TEXT NULL,
  attachments JSONB NOT NULL DEFAULT '[]'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_checkins_disciple ON checkins(disciple_id);

-- 8) Flags
CREATE TABLE IF NOT EXISTS user_flags (
  id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id   UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  injuries  JSONB NOT NULL DEFAULT '{}'::jsonb,
  equipment JSONB NOT NULL DEFAULT '{}'::jsonb,
  level     TEXT NULL
);

-- 9) Enlaces Coach
CREATE TABLE IF NOT EXISTS coach_links (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  coach_id UUID NOT NULL REFERENCES users(id),
  disciple_id UUID NOT NULL REFERENCES users(id),
  status TEXT NOT NULL CHECK (status IN ('pending','accepted','rejected')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (coach_id, disciple_id)
);
CREATE INDEX IF NOT EXISTS idx_coach_links_coach    ON coach_links(coach_id);
CREATE INDEX IF NOT EXISTS idx_coach_links_disciple ON coach_links(disciple_id);

-- 10) Función updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS invitations (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code         TEXT  NOT NULL UNIQUE,                -- código corto compartible
  coach_id     UUID  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  email        CITEXT NOT NULL,                      -- email del invitado
  name         TEXT  NULL,                           -- nombre tentativo
  status       TEXT  NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','accepted','revoked','expired')),
  expires_at   TIMESTAMPTZ NOT NULL,
  accepted_by  UUID  NULL REFERENCES users(id) ON DELETE SET NULL,
  accepted_at  TIMESTAMPTZ NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_invitations_coach  ON invitations(coach_id);
CREATE INDEX IF NOT EXISTS idx_invitations_email  ON invitations(email);
CREATE INDEX IF NOT EXISTS idx_invitations_status ON invitations(status);

-- 11) Triggers (sin IF NOT EXISTS; usamos DO con pg_trigger)

-- programs.updated_at
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger
    WHERE tgname = 'trg_programs_updated'
  ) THEN
    CREATE TRIGGER trg_programs_updated
    BEFORE UPDATE ON programs
    FOR EACH ROW EXECUTE PROCEDURE set_updated_at();
  END IF;
END;
$$;

-- users.updated_at
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger
    WHERE tgname = 'trg_users_updated'
  ) THEN
    CREATE TRIGGER trg_users_updated
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE PROCEDURE set_updated_at();
  END IF;
END;
$$;

-- coach_links.updated_at
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

COMMIT;