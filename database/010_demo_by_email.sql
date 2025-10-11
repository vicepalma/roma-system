-- Resuelve IDs por email (evitamos hardcodear UUIDs)
WITH u AS (
  SELECT
    (SELECT id FROM users WHERE email='roshi@kamehouse.example')   AS roshi,
    (SELECT id FROM users WHERE email='krillin@kamehouse.example') AS krillin,
    (SELECT id FROM users WHERE email='yamcha@capsule.example')    AS yamcha,
    (SELECT id FROM users WHERE email='goku@capsule.example')      AS goku,

    (SELECT id FROM users WHERE email='baki.hanma@example')        AS baki,
    (SELECT id FROM users WHERE email='retsu@shinshinkai.example') AS retsu,
    (SELECT id FROM users WHERE email='katsumi@shinshinkai.example') AS katsumi,
    (SELECT id FROM users WHERE email='jack.hanma@example')        AS jack,

    (SELECT id FROM users WHERE email='ikki@saint.example')        AS ikki
),

ins_master_disciple AS (
  INSERT INTO master_disciple (master_id, disciple_id, status)
  SELECT roshi, krillin, 'active' FROM u WHERE roshi IS NOT NULL AND krillin IS NOT NULL
  UNION ALL SELECT roshi, yamcha,  'active' FROM u
  UNION ALL SELECT roshi, goku,    'active' FROM u
  UNION ALL SELECT baki,  retsu,   'active' FROM u
  UNION ALL SELECT baki,  katsumi, 'active' FROM u
  UNION ALL SELECT baki,  jack,    'active' FROM u
  ON CONFLICT (master_id,disciple_id) DO NOTHING
  RETURNING 1
),
ins_methods AS (
  INSERT INTO methods (key, name, params)
  VALUES ('fst7','FST-7','{"series":7,"rest_sec":30,"to_failure":true,"target_reps":"10-12"}')
  ON CONFLICT (key) DO NOTHING
  RETURNING 1
),
-- Catálogo mínimo por si faltan
ins_ex AS (
  INSERT INTO exercises (name, primary_muscle, equipment, tags) VALUES
    ('Pec fly', 'chest', 'machine',  ARRAY['hypertrophy']),
    ('Press plano', 'chest', 'barbell', ARRAY['compound']),
    ('Posterior en poleas', 'shoulders', 'cable', ARRAY['rear-delt']),
    ('Press militar mancuernas', 'shoulders', 'dumbbell', ARRAY['compound']),
    ('Vuelos laterales', 'shoulders', 'dumbbell', ARRAY['isolation']),
    ('Extensión tríceps en polea', 'triceps', 'cable', ARRAY['isolation'])
  ON CONFLICT (lower(name)) DO NOTHING
  RETURNING 1
),

-- Programas (uno por maestro, + auto Ikki)
prog AS (
  INSERT INTO programs (owner_id, title, notes, visibility, version)
  SELECT roshi, 'Kame Style - Pecho/Hombro', 'Programa de Roshi', 'private', 1 FROM u
  UNION ALL
  SELECT baki,  'Shinshinkai - Full Body',   'Programa de Baki',  'private', 1 FROM u
  UNION ALL
  SELECT ikki,  'Fénix Solo',                'Auto Ikki',         'private', 1 FROM u
  RETURNING id, owner_id, title
),
w AS (
  INSERT INTO program_weeks (program_id, week_index)
  SELECT id, 1 FROM prog
  RETURNING id, program_id
),
d AS (
  INSERT INTO program_days (week_id, day_index, notes)
  SELECT id, 1, CASE
    WHEN (SELECT title FROM prog WHERE prog.program_id = w.program_id LIMIT 1) IS NOT NULL THEN
      CASE
        WHEN (SELECT title FROM prog WHERE prog.program_id = w.program_id LIMIT 1) ILIKE '%Kame%' THEN 'Pecho/Hombro'
        WHEN (SELECT title FROM prog WHERE prog.program_id = w.program_id LIMIT 1) ILIKE '%Shinshinkai%' THEN 'Full body A'
        ELSE 'Auto día 1'
      END
    ELSE 'Día 1'
  END
  FROM w
  RETURNING id, week_id
),
-- Prescripciones: una por programa (usando nombres de ejercicio)
p1 AS (
  INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, to_failure, position)
  SELECT d.id, e.id, 4, '10-12', 90, FALSE, 1
  FROM d
  JOIN program_weeks w2 ON w2.id = d.week_id
  JOIN programs pr ON pr.id = w2.program_id
  JOIN exercises e ON e.name = CASE
      WHEN pr.title ILIKE '%Kame%'       THEN 'Press plano'
      WHEN pr.title ILIKE '%Shinshinkai%' THEN 'Pec fly'
      ELSE 'Press militar mancuernas'
  END
  RETURNING id, day_id
),

-- Assignments: Roshi→(Krillin,Yamcha,Goku), Baki→(Retsu,Katsumi,Jack), Ikki→Ikki
asgn AS (
  INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
  SELECT pr.id, 1, (SELECT krillin FROM u), (SELECT roshi FROM u), CURRENT_DATE, TRUE
  FROM programs pr JOIN u ON pr.owner_id = (SELECT roshi FROM u)
  UNION ALL
  SELECT pr.id, 1, (SELECT yamcha  FROM u), (SELECT roshi FROM u), CURRENT_DATE, TRUE
  FROM programs pr JOIN u ON pr.owner_id = (SELECT roshi FROM u)
  UNION ALL
  SELECT pr.id, 1, (SELECT goku    FROM u), (SELECT roshi FROM u), CURRENT_DATE, TRUE
  FROM programs pr JOIN u ON pr.owner_id = (SELECT roshi FROM u)

  UNION ALL
  SELECT pr.id, 1, (SELECT retsu   FROM u), (SELECT baki  FROM u), CURRENT_DATE, TRUE
  FROM programs pr JOIN u ON pr.owner_id = (SELECT baki FROM u)
  UNION ALL
  SELECT pr.id, 1, (SELECT katsumi FROM u), (SELECT baki  FROM u), CURRENT_DATE, TRUE
  FROM programs pr JOIN u ON pr.owner_id = (SELECT baki FROM u)
  UNION ALL
  SELECT pr.id, 1, (SELECT jack    FROM u), (SELECT baki  FROM u), CURRENT_DATE, TRUE
  FROM programs pr JOIN u ON pr.owner_id = (SELECT baki FROM u)

  UNION ALL
  SELECT pr.id, 1, (SELECT ikki    FROM u), (SELECT ikki  FROM u), CURRENT_DATE, TRUE
  FROM programs pr JOIN u ON pr.owner_id = (SELECT ikki FROM u)
  RETURNING id, program_id, disciple_id
)

-- Una sesión de ejemplo para Krillin (para que pivot no quede en 0)
INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
SELECT a.id, a.disciple_id, d2.id, NOW(), 'Seed sesión Krillin'
FROM asgn a
JOIN programs pr ON pr.id = a.program_id
JOIN program_weeks w2 ON w2.program_id = pr.id
JOIN program_days d2 ON d2.week_id = w2.id
JOIN users u2 ON u2.id = a.disciple_id
WHERE u2.email = 'krillin@kamehouse.example'
LIMIT 1;

-- 3 sets de ejemplo (Press plano)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT
  s.id,
  p.id,
  x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (
  SELECT 1 AS set_index, 60::numeric AS weight, 12 AS reps, 7.5::numeric AS rpe, FALSE AS to_failure
  UNION ALL SELECT 2, 65, 10, 8.0, FALSE
  UNION ALL SELECT 3, 70,  8, 8.5, TRUE
) x
JOIN session_logs s ON s.notes = 'Seed sesión Krillin'
JOIN prescriptions p ON p.day_id = s.day_id
JOIN exercises e ON e.id = p.exercise_id AND e.name = 'Press plano'
ORDER BY s.performed_at DESC
LIMIT 3;
