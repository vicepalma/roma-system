DROP INDEX IF EXISTS idx_programs_owner_kind;

ALTER TABLE programs
DROP CONSTRAINT IF EXISTS chk_programs_kind;

ALTER TABLE programs
DROP COLUMN IF EXISTS kind;
