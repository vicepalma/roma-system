DROP INDEX IF EXISTS uq_assignments_active_per_disciple;

CREATE UNIQUE INDEX IF NOT EXISTS uq_self_assignments_active_per_disciple
ON assignments (disciple_id)
WHERE is_active = true AND assigned_by = disciple_id;
