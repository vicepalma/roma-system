BEGIN;

-- 1) Resolver usuarios por email
WITH u AS (
  SELECT
    (SELECT id FROM users WHERE email='roshi@kamehouse.example')   AS roshi,
    (SELECT id FROM users WHERE email='krillin@kamehouse.example') AS krillin,
    (SELECT id FROM users WHERE email='yamcha@capsule.example')    AS yamcha,
    (SELECT id FROM users WHERE email='goku@capsule.example')      AS goku,

    (SELECT id FROM users WHERE email='baki.hanma@example.example')        AS baki,
    (SELECT id FROM users WHERE email='retsu@shinshinkai.example') AS retsu,
    (SELECT id FROM users WHERE email='katsumi@shinshinkai.example') AS katsumi,
    (SELECT id FROM users WHERE email='jack.hanma@example.example')        AS jack,

    (SELECT id FROM users WHERE email='ikki@saint.example')        AS ikki
),

-- 2) Vínculos maestro–discípulo (incluye auto-vínculo de Ikki)
md AS (
  INSERT INTO master_disciple (master_id, disciple_id, status)
  SELECT roshi, krillin, 'active' FROM u WHERE roshi IS NOT NULL AND krillin IS NOT NULL
  UNION ALL SELECT roshi, yamcha,  'active' FROM u
  UNION ALL SELECT roshi, goku,    'active' FROM u
  UNION ALL SELECT baki,  retsu,   'active' FROM u
  UNION ALL SELECT baki,  katsumi, 'active' FROM u
  UNION ALL SELECT baki,  jack,    'active' FROM u
  UNION ALL SELECT ikki,  ikki,    'active' FROM u
  ON CONFLICT (master_id,disciple_id) DO NOTHING
  RETURNING 1
),

-- 3) Método y catálogo mínimo (por si faltan)
met AS (
  INSERT INTO methods (key, name, params)
  VALUES ('fst7','FST-7','{"series":7,"rest_sec":30,"to_failure":true,"target_reps":"10-12"}')
  ON CONFLICT (key) DO NOTHING
  RETURNING 1
),
ex AS (
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

-- 4) Programas (uno por dueño)
prog AS (
  INSERT INTO programs (owner_id, title, notes, visibility, version)
  SELECT roshi, 'Kame Style - Pecho/Hombro', 'Programa de Roshi', 'private', 1 FROM u WHERE roshi IS NOT NULL
  UNION ALL
  SELECT baki,  'Shinshinkai - Full Body',   'Programa de Baki',  'private', 1 FROM u WHERE baki  IS NOT NULL
  UNION ALL
  SELECT ikki,  'Fénix Solo',                'Auto Ikki',         'private', 1 FROM u WHERE ikki  IS NOT NULL
  ON CONFLICT DO NOTHING
  RETURNING id, owner_id, title
),

-- 5) Week 1 para cada programa
w AS (
  INSERT INTO program_weeks (program_id, week_index)
  SELECT id, 1 FROM prog
  RETURNING id, program_id
),

-- 6) Day 1 con notas segun el programa
d AS (
  INSERT INTO program_days (week_id, day_index, notes)
  SELECT w.id, 1,
         CASE
           WHEN p.title ILIKE '%Kame%'        THEN 'Pecho/Hombro'
           WHEN p.title ILIKE '%Shinshinkai%' THEN 'Full body A'
           ELSE 'Auto día 1'
         END
  FROM w
  JOIN programs p ON p.id = w.program_id
  RETURNING id, week_id
),

-- 7) 1 prescripción por programa (elige ejercicio por título)
presc AS (
  INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, to_failure, position)
  SELECT
    d.id,
    (SELECT e.id
       FROM exercises e
      WHERE e.name = CASE
                      WHEN p.title ILIKE '%Kame%'        THEN 'Press plano'
                      WHEN p.title ILIKE '%Shinshinkai%' THEN 'Pec fly'
                      ELSE 'Press militar mancuernas'
                    END
      LIMIT 1),
    4, '10-12', 90, FALSE, 1
  FROM d
  JOIN program_weeks w2 ON w2.id = d.week_id
  JOIN programs p ON p.id = w2.program_id
  RETURNING id, day_id
),

-- 8) Assignments: Roshi→(Krillin,Yamcha,Goku), Baki→(Retsu,Katsumi,Jack), Ikki→Ikki
asgn AS (
  INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
  -- Roshi
  SELECT pr.id, 1, u.krillin, u.roshi, CURRENT_DATE, TRUE FROM programs pr CROSS JOIN u
   WHERE pr.owner_id = u.roshi AND u.krillin IS NOT NULL
  UNION ALL
  SELECT pr.id, 1, u.yamcha,  u.roshi, CURRENT_DATE, TRUE FROM programs pr CROSS JOIN u
   WHERE pr.owner_id = u.roshi AND u.yamcha IS NOT NULL
  UNION ALL
  SELECT pr.id, 1, u.goku,    u.roshi, CURRENT_DATE, TRUE FROM programs pr CROSS JOIN u
   WHERE pr.owner_id = u.roshi AND u.goku IS NOT NULL
  -- Baki
  UNION ALL
  SELECT pr.id, 1, u.retsu,   u.baki,  CURRENT_DATE, TRUE FROM programs pr CROSS JOIN u
   WHERE pr.owner_id = u.baki  AND u.retsu IS NOT NULL
  UNION ALL
  SELECT pr.id, 1, u.katsumi, u.baki,  CURRENT_DATE, TRUE FROM programs pr CROSS JOIN u
   WHERE pr.owner_id = u.baki  AND u.katsumi IS NOT NULL
  UNION ALL
  SELECT pr.id, 1, u.jack,    u.baki,  CURRENT_DATE, TRUE FROM programs pr CROSS JOIN u
   WHERE pr.owner_id = u.baki  AND u.jack IS NOT NULL
  -- Ikki (self)
  UNION ALL
  SELECT pr.id, 1, u.ikki,    u.ikki,  CURRENT_DATE, TRUE FROM programs pr CROSS JOIN u
   WHERE pr.owner_id = u.ikki  AND u.ikki IS NOT NULL
  ON CONFLICT DO NOTHING
  RETURNING id, program_id, disciple_id
)

-- 9) Una sesión de ejemplo para Krillin (para que los pivots no queden en 0)
INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
SELECT
  a.id,
  u.krillin,
  d2.id,
  NOW(),
  'Seed sesión Krillin'
FROM asgn a
JOIN assignments a2 ON a2.id = a.id
JOIN programs pr ON pr.id = a2.program_id
JOIN program_weeks w2 ON w2.program_id = pr.id
JOIN program_days d2 ON d2.week_id = w2.id
JOIN u ON TRUE
WHERE u.krillin IS NOT NULL AND pr.owner_id = u.roshi
ORDER BY pr.created_at DESC
LIMIT 1;

-- 10) 3 sets para esa sesión (Press plano)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT
  s.id,
  (SELECT p.id FROM prescriptions p
    JOIN exercises e ON e.id = p.exercise_id
   WHERE p.day_id = s.day_id AND e.name = 'Press plano' LIMIT 1),
  x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM session_logs s
JOIN (VALUES
  (1, 60::numeric, 12, 7.5::numeric, FALSE),
  (2, 65::numeric, 10, 8.0::numeric, FALSE),
  (3, 70::numeric,  8, 8.5::numeric, TRUE )
) AS x(set_index,weight,reps,rpe,to_failure) ON TRUE
WHERE s.notes = 'Seed sesión Krillin'
ORDER BY s.performed_at DESC
LIMIT 3;

COMMIT;
