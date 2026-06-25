ALTER TABLE users
ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'disciple';

ALTER TABLE users
DROP CONSTRAINT IF EXISTS chk_users_role;

ALTER TABLE users
ADD CONSTRAINT chk_users_role CHECK (role IN ('coach', 'disciple'));

UPDATE users
SET role = 'coach'
WHERE id IN (
  SELECT DISTINCT coach_id FROM coach_links WHERE status = 'accepted'
)
OR id IN (
  SELECT DISTINCT owner_id FROM programs
);

UPDATE users
SET role = 'disciple'
WHERE role IS NULL OR role NOT IN ('coach', 'disciple');
