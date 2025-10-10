DROP TRIGGER IF EXISTS trg_users_updated ON users;
DROP TRIGGER IF EXISTS trg_programs_updated ON programs;
DROP FUNCTION IF EXISTS set_updated_at;

DROP TABLE IF EXISTS user_flags;
DROP TABLE IF EXISTS checkins;
DROP TABLE IF EXISTS cardio_segments;
DROP TABLE IF EXISTS set_logs;
DROP TABLE IF EXISTS session_logs;
DROP TABLE IF EXISTS assignments;
DROP TABLE IF EXISTS prescriptions;
DROP TABLE IF EXISTS program_days;
DROP TABLE IF EXISTS program_weeks;
DROP TABLE IF EXISTS programs;
DROP TABLE IF EXISTS methods;
DROP TABLE IF EXISTS exercises;
DROP TABLE IF EXISTS master_disciple;
DROP TABLE IF EXISTS users;

-- extensiones (no siempre conviene dropearlas en down)
-- DROP EXTENSION IF EXISTS citext;
-- DROP EXTENSION IF EXISTS pgcrypto;
