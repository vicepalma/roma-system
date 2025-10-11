BEGIN;

-- Mantiene catálogo (exercises/methods) y usuarios
TRUNCATE TABLE
  user_flags,
  checkins,
  cardio_segments,
  set_logs,
  session_logs,
  assignments,
  prescriptions,
  program_days,
  program_weeks,
  programs,
  master_disciple
RESTART IDENTITY CASCADE;

COMMIT;
