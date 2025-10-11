BEGIN;

-- Borra todo lo de negocio respetando FKs
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
  master_disciple,
  methods,
  exercises,
  users
RESTART IDENTITY CASCADE;

COMMIT;
