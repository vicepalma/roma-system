ALTER TABLE programs
ADD COLUMN IF NOT EXISTS kind TEXT NOT NULL DEFAULT 'coach_program';

ALTER TABLE programs
DROP CONSTRAINT IF EXISTS chk_programs_kind;

ALTER TABLE programs
ADD CONSTRAINT chk_programs_kind CHECK (kind IN ('coach_program', 'self_training'));

UPDATE programs
SET kind = 'coach_program'
WHERE kind IS NULL OR kind NOT IN ('coach_program', 'self_training');

CREATE INDEX IF NOT EXISTS idx_programs_owner_kind ON programs(owner_id, kind);
