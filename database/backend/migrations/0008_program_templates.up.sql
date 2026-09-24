CREATE TABLE program_templates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  notes TEXT NULL,
  status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','published')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  published_at TIMESTAMPTZ NULL
);
CREATE INDEX idx_program_templates_status ON program_templates(status);
CREATE INDEX idx_program_templates_owner ON program_templates(owner_id);

CREATE TABLE program_template_weeks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  template_id UUID NOT NULL REFERENCES program_templates(id) ON DELETE CASCADE,
  week_index INT NOT NULL CHECK (week_index >= 1),
  UNIQUE(template_id, week_index)
);
CREATE TABLE program_template_days (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  week_id UUID NOT NULL REFERENCES program_template_weeks(id) ON DELETE CASCADE,
  day_index INT NOT NULL CHECK (day_index >= 1),
  title TEXT NULL, notes TEXT NULL,
  UNIQUE(week_id, day_index)
);
CREATE TABLE program_template_prescriptions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  day_id UUID NOT NULL REFERENCES program_template_days(id) ON DELETE CASCADE,
  exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE RESTRICT,
  series INT NOT NULL CHECK (series >= 1), reps TEXT NOT NULL,
  rest_sec INT NULL CHECK (rest_sec IS NULL OR rest_sec >= 0),
  to_failure BOOLEAN NOT NULL DEFAULT false, tempo TEXT NULL, rir INT NULL,
  rpe NUMERIC(3,1) NULL, method_id UUID NULL REFERENCES methods(id) ON DELETE SET NULL,
  notes TEXT NULL, position INT NOT NULL DEFAULT 1
);
CREATE INDEX idx_template_days_week ON program_template_days(week_id);
CREATE INDEX idx_template_presc_day ON program_template_prescriptions(day_id);
