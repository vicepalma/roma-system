-- Más sesiones para probar history/pivot con variedad de días/ejercicios.
-- Depende de 002_characters_seed.sql cargado previamente.

-- =========================
-- 1) INSERT de sesiones
-- =========================

-- Retsu: hace "Pec fly" hace 3 días a las 10:30
INSERT INTO session_logs (id, assignment_id, disciple_id, day_id, performed_at, notes)
VALUES (
  gen_random_uuid(),
  (SELECT a.id FROM assignments a
   WHERE a.disciple_id = 'bbb11111-1111-1111-1111-111111111111'
     AND a.program_id  = '38173377-7223-413f-8295-62bd8d2aa978'
   ORDER BY a.created_at DESC LIMIT 1),
  'bbb11111-1111-1111-1111-111111111111',
  (SELECT d.id FROM program_days d
   JOIN program_weeks w ON w.id = d.week_id
   WHERE w.program_id = '38173377-7223-413f-8295-62bd8d2aa978'
     AND d.day_index = 1
   LIMIT 1),
  (CURRENT_DATE - INTERVAL '3 days') + TIME '10:30',
  'Retsu - Pec fly'
);

-- Katsumi: "Extensión tríceps en polea" hace 2 días a las 19:15
INSERT INTO session_logs (id, assignment_id, disciple_id, day_id, performed_at, notes)
VALUES (
  gen_random_uuid(),
  (SELECT a.id FROM assignments a
   WHERE a.disciple_id = 'bbb22222-2222-2222-2222-222222222222'
     AND a.program_id  = '38173377-7223-413f-8295-62bd8d2aa978'
   ORDER BY a.created_at DESC LIMIT 1),
  'bbb22222-2222-2222-2222-222222222222',
  (SELECT d.id FROM program_days d
   JOIN program_weeks w ON w.id = d.week_id
   WHERE w.program_id = '38173377-7223-413f-8295-62bd8d2aa978'
     AND d.day_index = 1
   LIMIT 1),
  (CURRENT_DATE - INTERVAL '2 days') + TIME '19:15',
  'Katsumi - Tríceps en polea'
);

-- Jack: hace ambas (pec + triceps) ayer a las 08:05
INSERT INTO session_logs (id, assignment_id, disciple_id, day_id, performed_at, notes)
VALUES (
  gen_random_uuid(),
  (SELECT a.id FROM assignments a
   WHERE a.disciple_id = 'bbb33333-3333-3333-3333-333333333333'
     AND a.program_id  = '38173377-7223-413f-8295-62bd8d2aa978'
   ORDER BY a.created_at DESC LIMIT 1),
  'bbb33333-3333-3333-3333-333333333333',
  (SELECT d.id FROM program_days d
   JOIN program_weeks w ON w.id = d.week_id
   WHERE w.program_id = '38173377-7223-413f-8295-62bd8d2aa978'
     AND d.day_index = 1
   LIMIT 1),
  (CURRENT_DATE - INTERVAL '1 day') + TIME '08:05',
  'Jack - Full body A'
);

-- Ikki: "Press militar mancuernas" hoy a las 06:50
INSERT INTO session_logs (id, assignment_id, disciple_id, day_id, performed_at, notes)
VALUES (
  gen_random_uuid(),
  (SELECT a.id FROM assignments a
   WHERE a.disciple_id = '33333333-3333-3333-3333-333333333333'
     AND a.program_id  = '49173377-7223-413f-8295-62bd8d2aa978'
   ORDER BY a.created_at DESC LIMIT 1),
  '33333333-3333-3333-3333-333333333333',
  (SELECT d.id FROM program_days d
   JOIN program_weeks w ON w.id = d.week_id
   WHERE w.program_id = '49173377-7223-413f-8295-62bd8d2aa978'
     AND d.day_index = 1
   LIMIT 1),
  (CURRENT_DATE + TIME '06:50'),
  'Ikki - amanecer'
);

-- Krillin: sesión extra (Press plano) hace 4 días a las 18:40
INSERT INTO session_logs (id, assignment_id, disciple_id, day_id, performed_at, notes)
VALUES (
  gen_random_uuid(),
  (SELECT a.id FROM assignments a
   WHERE a.disciple_id = 'aaa11111-1111-1111-1111-111111111111'
     AND a.program_id  = '27173377-7223-413f-8295-62bd8d2aa978'
   ORDER BY a.created_at DESC LIMIT 1),
  'aaa11111-1111-1111-1111-111111111111',
  (SELECT d.id FROM program_days d
   JOIN program_weeks w ON w.id = d.week_id
   WHERE w.program_id = '27173377-7223-413f-8295-62bd8d2aa978'
     AND d.day_index = 1
   LIMIT 1),
  (CURRENT_DATE - INTERVAL '4 days') + TIME '18:40',
  'Krillin - extra pecho'
);

-- =========================
-- 2) INSERT de sets
-- =========================

-- Helper: subconsulta para obtener sesión por disciple + note
--   (tomamos la más reciente con ese note)
-- Helper: subconsulta para obtener prescription por programa/día/ejercicio

-- Retsu (Pec fly)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT
  (SELECT s.id FROM session_logs s
   WHERE s.disciple_id = 'bbb11111-1111-1111-1111-111111111111'
     AND s.notes = 'Retsu - Pec fly'
   ORDER BY s.performed_at DESC LIMIT 1) AS session_id,
  (SELECT p.id FROM prescriptions p
   JOIN program_days d ON d.id = p.day_id
   JOIN program_weeks w ON w.id = d.week_id
   JOIN exercises e ON e.id = p.exercise_id
   WHERE w.program_id = '38173377-7223-413f-8295-62bd8d2aa978'
     AND d.day_index = 1
     AND e.name = 'Pec fly'
   LIMIT 1) AS prescription_id,
  x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (
  SELECT 1 AS set_index, 30::numeric AS weight, 15 AS reps, 7.0::numeric AS rpe, FALSE AS to_failure
  UNION ALL SELECT 2, 35, 12, 7.5, FALSE
  UNION ALL SELECT 3, 40, 10, 8.0, FALSE
) x;

-- Katsumi (Extensión tríceps en polea)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT
  (SELECT s.id FROM session_logs s
   WHERE s.disciple_id = 'bbb22222-2222-2222-2222-222222222222'
     AND s.notes = 'Katsumi - Tríceps en polea'
   ORDER BY s.performed_at DESC LIMIT 1),
  (SELECT p.id FROM prescriptions p
   JOIN program_days d ON d.id = p.day_id
   JOIN program_weeks w ON w.id = d.week_id
   JOIN exercises e ON e.id = p.exercise_id
   WHERE w.program_id = '38173377-7223-413f-8295-62bd8d2aa978'
     AND d.day_index = 1
     AND e.name = 'Extensión tríceps en polea'
   LIMIT 1),
  x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (
  SELECT 1 AS set_index, 25::numeric AS weight, 15 AS reps, 7.0::numeric AS rpe, FALSE AS to_failure
  UNION ALL SELECT 2, 30, 12, 7.5, FALSE
  UNION ALL SELECT 3, 35, 12, 8.0, TRUE
) x;

-- Jack (Pec fly)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT
  (SELECT s.id FROM session_logs s
   WHERE s.disciple_id = 'bbb33333-3333-3333-3333-333333333333'
     AND s.notes = 'Jack - Full body A'
   ORDER BY s.performed_at DESC LIMIT 1),
  (SELECT p.id FROM prescriptions p
   JOIN program_days d ON d.id = p.day_id
   JOIN program_weeks w ON w.id = d.week_id
   JOIN exercises e ON e.id = p.exercise_id
   WHERE w.program_id = '38173377-7223-413f-8295-62bd8d2aa978'
     AND d.day_index = 1
     AND e.name = 'Pec fly'
   LIMIT 1),
  x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (
  SELECT 1 AS set_index, 50::numeric AS weight, 10 AS reps, 7.5::numeric AS rpe, FALSE AS to_failure
  UNION ALL SELECT 2, 55, 8, 8.0, FALSE
  UNION ALL SELECT 3, 60, 6, 8.5, TRUE
) x;

-- Jack (Extensión tríceps en polea)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT
  (SELECT s.id FROM session_logs s
   WHERE s.disciple_id = 'bbb33333-3333-3333-3333-333333333333'
     AND s.notes = 'Jack - Full body A'
   ORDER BY s.performed_at DESC LIMIT 1),
  (SELECT p.id FROM prescriptions p
   JOIN program_days d ON d.id = p.day_id
   JOIN program_weeks w ON w.id = d.week_id
   JOIN exercises e ON e.id = p.exercise_id
   WHERE w.program_id = '38173377-7223-413f-8295-62bd8d2aa978'
     AND d.day_index = 1
     AND e.name = 'Extensión tríceps en polea'
   LIMIT 1),
  x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (
  SELECT 1 AS set_index, 35::numeric AS weight, 12 AS reps, 7.0::numeric AS rpe, FALSE AS to_failure
  UNION ALL SELECT 2, 40, 10, 7.5, FALSE
  UNION ALL SELECT 3, 45, 8, 8.5, TRUE
) x;

-- Ikki (Press militar mancuernas)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT
  (SELECT s.id FROM session_logs s
   WHERE s.disciple_id = '33333333-3333-3333-3333-333333333333'
     AND s.notes = 'Ikki - amanecer'
   ORDER BY s.performed_at DESC LIMIT 1),
  (SELECT p.id FROM prescriptions p
   JOIN program_days d ON d.id = p.day_id
   JOIN program_weeks w ON w.id = d.week_id
   JOIN exercises e ON e.id = p.exercise_id
   WHERE w.program_id = '49173377-7223-413f-8295-62bd8d2aa978'
     AND d.day_index = 1
     AND e.name = 'Press militar mancuernas'
   LIMIT 1),
  x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (
  SELECT 1 AS set_index, 22.5::numeric AS weight, 10 AS reps, 7.0::numeric AS rpe, FALSE AS to_failure
  UNION ALL SELECT 2, 25, 8, 8.0, FALSE
  UNION ALL SELECT 3, 27.5, 6, 8.5, TRUE
) x;

-- Krillin (Press plano) - sesión extra
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT
  (SELECT s.id FROM session_logs s
   WHERE s.disciple_id = 'aaa11111-1111-1111-1111-111111111111'
     AND s.notes = 'Krillin - extra pecho'
   ORDER BY s.performed_at DESC LIMIT 1),
  (SELECT p.id FROM prescriptions p
   JOIN program_days d ON d.id = p.day_id
   JOIN program_weeks w ON w.id = d.week_id
   JOIN exercises e ON e.id = p.exercise_id
   WHERE w.program_id = '27173377-7223-413f-8295-62bd8d2aa978'
     AND d.day_index = 1
     AND e.name = 'Press plano'
   LIMIT 1),
  x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (
  SELECT 1 AS set_index, 62.5::numeric AS weight, 12 AS reps, 7.5::numeric AS rpe, FALSE AS to_failure
  UNION ALL SELECT 2, 67.5, 10, 8.0, FALSE
  UNION ALL SELECT 3, 72.5, 8, 8.5, TRUE
) x;
