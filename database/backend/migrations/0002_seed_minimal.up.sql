BEGIN;

-- ===== Hash para "1234" (mismo de tus seeds) =====
-- $2a$10$PbvWZsbmMSbPSH5hQLm5peE4AB7p0RhxT6QEBjtkPC3T0UPq2Up5i
WITH cfg AS (
  -- Genera dinámicamente un hash bcrypt para la contraseña 'secret123' (cost = 10)
  SELECT crypt('secret123', gen_salt('bf', 10))::text AS ph
)

-- ===== Usuarios: maestros y discípulos =====
INSERT INTO users (email, password_hash, name)
SELECT * FROM (
  -- Maestros
  SELECT 'roshi@kamehouse.example',        (SELECT ph FROM cfg), 'Maestro Roshi' UNION ALL
  SELECT 'baki.hanma@example.example',     (SELECT ph FROM cfg), 'Baki Hanma'    UNION ALL
  -- Auto-maestro
  SELECT 'ikki@saint.example',             (SELECT ph FROM cfg), 'Ikki (Fénix)'  UNION ALL

  -- Discípulos de Roshi
  SELECT 'krillin@kamehouse.example',      (SELECT ph FROM cfg), 'Krillin'       UNION ALL
  SELECT 'yamcha@capsule.example',         (SELECT ph FROM cfg), 'Yamcha'        UNION ALL
  SELECT 'goku@capsule.example',           (SELECT ph FROM cfg), 'Goku'          UNION ALL

  -- Discípulos de Baki
  SELECT 'retsu@shinshinkai.example',      (SELECT ph FROM cfg), 'Retsu Kaioh'   UNION ALL
  SELECT 'katsumi@shinshinkai.example',    (SELECT ph FROM cfg), 'Katsumi Orochi' UNION ALL
  SELECT 'jack.hanma@example.example',     (SELECT ph FROM cfg), 'Jack Hanma'
) AS u(email,password_hash,name)
ON CONFLICT (email) DO NOTHING;

-- ===== Vínculos maestro–discípulo (fix alias) =====
WITH u AS (
  SELECT
    (SELECT id FROM users WHERE email='roshi@kamehouse.example')            AS roshi,
    (SELECT id FROM users WHERE email='baki.hanma@example.example')         AS baki,
    (SELECT id FROM users WHERE email='ikki@saint.example')                 AS ikki,
    (SELECT id FROM users WHERE email='krillin@kamehouse.example')          AS krillin,
    (SELECT id FROM users WHERE email='yamcha@capsule.example')             AS yamcha,
    (SELECT id FROM users WHERE email='goku@capsule.example')               AS goku,
    (SELECT id FROM users WHERE email='retsu@shinshinkai.example')          AS retsu,
    (SELECT id FROM users WHERE email='katsumi@shinshinkai.example')        AS katsumi,
    (SELECT id FROM users WHERE email='jack.hanma@example.example')         AS jack
)
INSERT INTO master_disciple (master_id, disciple_id, status)
SELECT *
FROM (
  SELECT roshi  AS master_id, krillin AS disciple_id, 'active'::text AS status FROM u UNION ALL
  SELECT roshi, yamcha,  'active' FROM u UNION ALL
  SELECT roshi, goku,    'active' FROM u UNION ALL
  SELECT baki,  retsu,   'active' FROM u UNION ALL
  SELECT baki,  katsumi, 'active' FROM u UNION ALL
  SELECT baki,  jack,    'active' FROM u UNION ALL
  SELECT ikki,  ikki,    'active' FROM u
) AS x(master_id, disciple_id, status)
WHERE x.master_id IS NOT NULL
  AND x.disciple_id IS NOT NULL
ON CONFLICT (master_id, disciple_id) DO NOTHING;

-- ===== Catálogo mínimo (método + ejercicios) =====
INSERT INTO methods (key, name, params)
VALUES ('fst7','FST-7','{"series":7,"rest_sec":30,"to_failure":true,"target_reps":"10-12"}')
ON CONFLICT (key) DO NOTHING;

INSERT INTO exercises (name, primary_muscle, equipment, tags) VALUES
  ('Pec fly', 'chest', 'machine',  ARRAY['hypertrophy']),
  ('Press plano', 'chest', 'barbell', ARRAY['compound']),
  ('Posterior en poleas', 'shoulders', 'cable', ARRAY['rear-delt']),
  ('Press militar mancuernas', 'shoulders', 'dumbbell', ARRAY['compound']),
  ('Vuelos laterales', 'shoulders', 'dumbbell', ARRAY['isolation']),
  ('Extensión tríceps en polea', 'triceps', 'cable', ARRAY['isolation'])
ON CONFLICT (lower(name)) DO NOTHING;

-- ===== Programas por owner (Roshi / Baki / Ikki) + week(1) + day(1) =====
WITH owners AS (
  SELECT
    (SELECT id FROM users WHERE email='roshi@kamehouse.example')        AS roshi,
    (SELECT id FROM users WHERE email='baki.hanma@example.example')     AS baki,
    (SELECT id FROM users WHERE email='ikki@saint.example')             AS ikki
),
prog AS (
  INSERT INTO programs (owner_id, title, notes, visibility, version)
  SELECT roshi, 'Kame Style - Pecho/Hombro', 'Programa de Roshi', 'private', 1 FROM owners WHERE roshi IS NOT NULL
  UNION ALL
  SELECT baki,  'Shinshinkai - Full Body',   'Programa de Baki',  'private', 1 FROM owners WHERE baki  IS NOT NULL
  UNION ALL
  SELECT ikki,  'Fénix Solo',                'Auto Ikki',         'private', 1 FROM owners WHERE ikki  IS NOT NULL
  ON CONFLICT DO NOTHING
  RETURNING id, owner_id
),
weeks AS (
  INSERT INTO program_weeks (program_id, week_index)
  SELECT id, 1 FROM prog
  ON CONFLICT DO NOTHING
  RETURNING id, program_id
),
days AS (
  INSERT INTO program_days (week_id, day_index, notes)
  SELECT w.id, 1,
         CASE
           WHEN p.title ILIKE '%Kame%'        THEN 'Pecho/Hombro'
           WHEN p.title ILIKE '%Shinshinkai%' THEN 'Full body A'
           ELSE 'Auto día 1'
         END
  FROM weeks w
  JOIN programs p ON p.id = w.program_id
  ON CONFLICT DO NOTHING
  RETURNING id, week_id
)
-- ===== Prescripciones mínimas por programa =====
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
FROM days d
JOIN program_weeks w ON w.id = d.week_id
JOIN programs p ON p.id = w.program_id;

-- ===== Safety net: garantizar week 1, day 1 y prescripciones mínimas por PROGRAMA (mover aquí) =====

-- 1) Asegura que todos los programas tengan week 1
INSERT INTO program_weeks (program_id, week_index)
SELECT p.id, 1
FROM programs p
WHERE NOT EXISTS (
  SELECT 1 FROM program_weeks w WHERE w.program_id = p.id AND w.week_index = 1
);

-- 2) Asegura que cada week 1 tenga day 1
INSERT INTO program_days (week_id, day_index, notes)
SELECT w.id, 1, COALESCE(
  (SELECT CASE
           WHEN pr.title ILIKE '%Kame%'        THEN 'Pecho/Hombro'
           WHEN pr.title ILIKE '%Shinshinkai%' THEN 'Full body A'
           ELSE 'Auto día 1'
         END
   FROM programs pr WHERE pr.id = w.program_id),
  'Auto día 1'
)
FROM program_weeks w
WHERE w.week_index = 1
  AND NOT EXISTS (
    SELECT 1 FROM program_days d WHERE d.week_id = w.id AND d.day_index = 1
);

-- 3) Asegura que el day 1 de cada programa tenga al menos 1 prescripción
INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, to_failure, position)
SELECT d.id,
       COALESCE(
         (SELECT id FROM exercises WHERE name IN
            ('Press plano','Pec fly','Press militar mancuernas','Extensión tríceps en polea')
          ORDER BY name LIMIT 1),
         (SELECT id FROM exercises ORDER BY id LIMIT 1)
       ) AS exercise_id,
       3, '10-12', 90, FALSE, 1
FROM program_days d
LEFT JOIN prescriptions p ON p.day_id = d.id
WHERE d.day_index = 1
  AND p.id IS NULL;

-- 4) Prescripciones específicas utilizadas por las sesiones (repetir WITH d en cada INSERT)

-- Pec fly
WITH d AS (
  SELECT d.id AS day_id
  FROM program_days d
  JOIN program_weeks w ON w.id=d.week_id
  WHERE d.day_index=1
)
INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, to_failure, position)
SELECT d.day_id, e.id, 3, '12', 60, FALSE, 1
FROM d
JOIN exercises e ON e.name='Pec fly'
WHERE NOT EXISTS (SELECT 1 FROM prescriptions p WHERE p.day_id=d.day_id AND p.exercise_id=e.id);

-- Extensión tríceps en polea
WITH d AS (
  SELECT d.id AS day_id
  FROM program_days d
  JOIN program_weeks w ON w.id=d.week_id
  WHERE d.day_index=1
)
INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, to_failure, position)
SELECT d.day_id, e.id, 3, '12', 60, FALSE, 2
FROM d
JOIN exercises e ON e.name='Extensión tríceps en polea'
WHERE NOT EXISTS (SELECT 1 FROM prescriptions p WHERE p.day_id=d.day_id AND p.exercise_id=e.id);

-- Press plano
WITH d AS (
  SELECT d.id AS day_id
  FROM program_days d
  JOIN program_weeks w ON w.id=d.week_id
  WHERE d.day_index=1
)
INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, to_failure, position)
SELECT d.day_id, e.id, 4, '10', 90, FALSE, 1
FROM d
JOIN exercises e ON e.name='Press plano'
WHERE NOT EXISTS (SELECT 1 FROM prescriptions p WHERE p.day_id=d.day_id AND p.exercise_id=e.id);

-- Press militar mancuernas
WITH d AS (
  SELECT d.id AS day_id
  FROM program_days d
  JOIN program_weeks w ON w.id=d.week_id
  WHERE d.day_index=1
)
INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, to_failure, position)
SELECT d.day_id, e.id, 4, '8-10', 90, FALSE, 1
FROM d
JOIN exercises e ON e.name='Press militar mancuernas'
WHERE NOT EXISTS (SELECT 1 FROM prescriptions p WHERE p.day_id=d.day_id AND p.exercise_id=e.id);

-- ===== Assignments (Roshi→3, Baki→3, Ikki→auto) =====
WITH u AS (
  SELECT
    (SELECT id FROM users WHERE email='roshi@kamehouse.example')            AS roshi,
    (SELECT id FROM users WHERE email='baki.hanma@example.example')         AS baki,
    (SELECT id FROM users WHERE email='ikki@saint.example')                 AS ikki,
    (SELECT id FROM users WHERE email='krillin@kamehouse.example')          AS krillin,
    (SELECT id FROM users WHERE email='yamcha@capsule.example')             AS yamcha,
    (SELECT id FROM users WHERE email='goku@capsule.example')               AS goku,
    (SELECT id FROM users WHERE email='retsu@shinshinkai.example')          AS retsu,
    (SELECT id FROM users WHERE email='katsumi@shinshinkai.example')        AS katsumi,
    (SELECT id FROM users WHERE email='jack.hanma@example.example')         AS jack
)
INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
SELECT p.id, 1, d_id, m_id, CURRENT_DATE, TRUE
FROM (
  -- Roshi
  SELECT (SELECT id FROM programs WHERE owner_id=(SELECT roshi  FROM u) ORDER BY created_at DESC LIMIT 1) AS p_id,
         (SELECT krillin FROM u) AS d_id,
         (SELECT roshi   FROM u) AS m_id
  UNION ALL
  SELECT (SELECT id FROM programs WHERE owner_id=(SELECT roshi  FROM u) ORDER BY created_at DESC LIMIT 1),
         (SELECT yamcha FROM u), (SELECT roshi FROM u)
  UNION ALL
  SELECT (SELECT id FROM programs WHERE owner_id=(SELECT roshi  FROM u) ORDER BY created_at DESC LIMIT 1),
         (SELECT goku   FROM u), (SELECT roshi FROM u)

  -- Baki
  UNION ALL
  SELECT (SELECT id FROM programs WHERE owner_id=(SELECT baki   FROM u) ORDER BY created_at DESC LIMIT 1),
         (SELECT retsu  FROM u), (SELECT baki  FROM u)
  UNION ALL
  SELECT (SELECT id FROM programs WHERE owner_id=(SELECT baki   FROM u) ORDER BY created_at DESC LIMIT 1),
         (SELECT katsumi FROM u), (SELECT baki FROM u)
  UNION ALL
  SELECT (SELECT id FROM programs WHERE owner_id=(SELECT baki   FROM u) ORDER BY created_at DESC LIMIT 1),
         (SELECT jack   FROM u), (SELECT baki FROM u)

  -- Ikki (auto)
  UNION ALL
  SELECT (SELECT id FROM programs WHERE owner_id=(SELECT ikki   FROM u) ORDER BY created_at DESC LIMIT 1),
         (SELECT ikki   FROM u), (SELECT ikki FROM u)
) x
JOIN programs p ON p.id = x.p_id;

-- ===== Sesiones y sets de ejemplo (con guardias por si faltara day/prescripción) =====

-- Retsu: Pec fly (hace 3 días, 10:30)
WITH emails AS (
  SELECT 'baki.hanma@example.example'::text AS coach, 'retsu@shinshinkai.example'::text AS disciple
), owner_prog AS (
  SELECT (SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
          WHERE u.email=(SELECT coach FROM emails) ORDER BY p.created_at DESC LIMIT 1) AS prog_id
), day1 AS (
  SELECT COALESCE(
           (SELECT d.id FROM program_days d JOIN program_weeks w ON w.id=d.week_id
             WHERE w.program_id=(SELECT prog_id FROM owner_prog) AND d.day_index=1
             ORDER BY d.id LIMIT 1),
           (SELECT d2.id FROM program_days d2 JOIN program_weeks w2 ON w2.id=d2.week_id
             WHERE w2.program_id=(SELECT prog_id FROM owner_prog)
             ORDER BY d2.day_index, d2.id LIMIT 1)
         ) AS day_id
), asg AS (
  SELECT (SELECT a.id FROM assignments a
            WHERE a.disciple_id=(SELECT id FROM users WHERE email=(SELECT disciple FROM emails))
              AND a.program_id =(SELECT prog_id FROM owner_prog)
            ORDER BY a.created_at DESC LIMIT 1) AS asg_id
), new_s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT (SELECT asg_id FROM asg),
         (SELECT id FROM users WHERE email=(SELECT disciple FROM emails)),
         (SELECT day_id FROM day1),
         (CURRENT_DATE - INTERVAL '3 days') + TIME '10:30',
         'Retsu - Pec fly'
  WHERE (SELECT day_id FROM day1) IS NOT NULL
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM new_s LIMIT 1),
       (SELECT p.id FROM prescriptions p WHERE p.day_id=(SELECT day_id FROM day1)
         AND p.exercise_id=(SELECT e.id FROM exercises e WHERE e.name='Pec fly' LIMIT 1)
         ORDER BY p.position, p.id LIMIT 1),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (VALUES (1,30::numeric,15,7.0::numeric,false),
             (2,35::numeric,12,7.5::numeric,false),
             (3,40::numeric,10,8.0::numeric,false)) AS x(set_index,weight,reps,rpe,to_failure)
WHERE EXISTS (SELECT 1 FROM new_s);

-- Katsumi: Tríceps en polea (hace 2 días, 19:15)
WITH emails AS (
  SELECT 'baki.hanma@example.example'::text AS coach, 'katsumi@shinshinkai.example'::text AS disciple
), owner_prog AS (
  SELECT (SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
          WHERE u.email=(SELECT coach FROM emails) ORDER BY p.created_at DESC LIMIT 1) AS prog_id
), day1 AS (
  SELECT COALESCE(
           (SELECT d.id FROM program_days d JOIN program_weeks w ON w.id=d.week_id
             WHERE w.program_id=(SELECT prog_id FROM owner_prog) AND d.day_index=1
             ORDER BY d.id LIMIT 1),
           (SELECT d2.id FROM program_days d2 JOIN program_weeks w2 ON w2.id=d2.week_id
             WHERE w2.program_id=(SELECT prog_id FROM owner_prog)
             ORDER BY d2.day_index, d2.id LIMIT 1)
         ) AS day_id
), asg AS (
  SELECT (SELECT a.id FROM assignments a
            WHERE a.disciple_id=(SELECT id FROM users WHERE email=(SELECT disciple FROM emails))
              AND a.program_id =(SELECT prog_id FROM owner_prog)
            ORDER BY a.created_at DESC LIMIT 1) AS asg_id
), new_s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT (SELECT asg_id FROM asg),
         (SELECT id FROM users WHERE email=(SELECT disciple FROM emails)),
         (SELECT day_id FROM day1),
         (CURRENT_DATE - INTERVAL '2 days') + TIME '19:15',
         'Katsumi - Tríceps en polea'
  WHERE (SELECT day_id FROM day1) IS NOT NULL
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM new_s LIMIT 1),
       (SELECT p.id FROM prescriptions p WHERE p.day_id=(SELECT day_id FROM day1)
         AND p.exercise_id=(SELECT e.id FROM exercises e WHERE e.name='Extensión tríceps en polea' LIMIT 1)
         ORDER BY p.position, p.id LIMIT 1),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (VALUES (1,25::numeric,15,7.0::numeric,false),
             (2,30::numeric,12,7.5::numeric,false),
             (3,35::numeric,12,8.0::numeric,true)) AS x(set_index,weight,reps,rpe,to_failure)
WHERE EXISTS (SELECT 1 FROM new_s);

-- Jack: Full body A (ayer 08:05) – Pec fly + tríceps (dos prescripciones)
WITH emails AS (
  SELECT 'baki.hanma@example.example'::text AS coach, 'jack.hanma@example.example'::text AS disciple
), owner_prog AS (
  SELECT (SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
          WHERE u.email=(SELECT coach FROM emails) ORDER BY p.created_at DESC LIMIT 1) AS prog_id
), day1 AS (
  SELECT COALESCE(
           (SELECT d.id FROM program_days d JOIN program_weeks w ON w.id=d.week_id
             WHERE w.program_id=(SELECT prog_id FROM owner_prog) AND d.day_index=1
             ORDER BY d.id LIMIT 1),
           (SELECT d2.id FROM program_days d2 JOIN program_weeks w2 ON w2.id=d2.week_id
             WHERE w2.program_id=(SELECT prog_id FROM owner_prog)
             ORDER BY d2.day_index, d2.id LIMIT 1)
         ) AS day_id
), asg AS (
  SELECT (SELECT a.id FROM assignments a
            WHERE a.disciple_id=(SELECT id FROM users WHERE email=(SELECT disciple FROM emails))
              AND a.program_id =(SELECT prog_id FROM owner_prog)
            ORDER BY a.created_at DESC LIMIT 1) AS asg_id
), s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT (SELECT asg_id FROM asg),
         (SELECT id FROM users WHERE email=(SELECT disciple FROM emails)),
         (SELECT day_id FROM day1),
         (CURRENT_DATE - INTERVAL '1 day') + TIME '08:05',
         'Jack - Full body A'
  WHERE (SELECT day_id FROM day1) IS NOT NULL
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM s LIMIT 1),
       (SELECT p.id FROM prescriptions p WHERE p.day_id=(SELECT day_id FROM day1)
         AND p.exercise_id=(SELECT e.id FROM exercises e WHERE e.name='Pec fly' LIMIT 1)
         ORDER BY p.position, p.id LIMIT 1),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (VALUES (1,50::numeric,10,7.5::numeric,false),
             (2,55::numeric, 8,8.0::numeric,false),
             (3,60::numeric, 6,8.5::numeric,true)) AS x(set_index,weight,reps,rpe,to_failure)
WHERE EXISTS (SELECT 1 FROM s);

WITH sess AS (
  SELECT id
  FROM session_logs
  WHERE notes = 'Jack - Full body A'
  ORDER BY performed_at DESC
  LIMIT 1
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT
  (SELECT id FROM sess),
  (
    SELECT p.id FROM prescriptions p
    JOIN program_days d ON d.id=p.day_id
    JOIN program_weeks w ON w.id=d.week_id
    JOIN programs pr ON pr.id=w.program_id
    JOIN users cu ON cu.id=pr.owner_id
    JOIN exercises e ON e.id=p.exercise_id
    WHERE cu.email='baki.hanma@example.example'
      AND d.day_index=1
      AND e.name='Extensión tríceps en polea'
    ORDER BY pr.created_at DESC, p.position, p.id
    LIMIT 1
  ) AS prescription_id,
  x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (VALUES (1,35::numeric,12,7.0::numeric,false),
             (2,40::numeric,10,7.5::numeric,false),
             (3,45::numeric, 8,8.5::numeric,true)) AS x(set_index,weight,reps,rpe,to_failure)
WHERE EXISTS (SELECT 1 FROM sess);

-- Krillin: Press plano (hace 4 días 18:40)
WITH emails AS (
  SELECT 'roshi@kamehouse.example'::text AS coach, 'krillin@kamehouse.example'::text AS disciple
), owner_prog AS (
  SELECT (SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
          WHERE u.email=(SELECT coach FROM emails) ORDER BY p.created_at DESC LIMIT 1) AS prog_id
), day1 AS (
  SELECT COALESCE(
           (SELECT d.id FROM program_days d JOIN program_weeks w ON w.id=d.week_id
             WHERE w.program_id=(SELECT prog_id FROM owner_prog) AND d.day_index=1
             ORDER BY d.id LIMIT 1),
           (SELECT d2.id FROM program_days d2 JOIN program_weeks w2 ON w2.id=d2.week_id
             WHERE w2.program_id=(SELECT prog_id FROM owner_prog)
             ORDER BY d2.day_index, d2.id LIMIT 1)
         ) AS day_id
), asg AS (
  SELECT (SELECT a.id FROM assignments a
            WHERE a.disciple_id=(SELECT id FROM users WHERE email=(SELECT disciple FROM emails))
              AND a.program_id =(SELECT prog_id FROM owner_prog)
            ORDER BY a.created_at DESC LIMIT 1) AS asg_id
), new_s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT (SELECT asg_id FROM asg),
         (SELECT id FROM users WHERE email=(SELECT disciple FROM emails)),
         (SELECT day_id FROM day1),
         (CURRENT_DATE - INTERVAL '4 days') + TIME '18:40',
         'Krillin - extra pecho'
  WHERE (SELECT day_id FROM day1) IS NOT NULL
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM new_s LIMIT 1),
       (SELECT p.id FROM prescriptions p WHERE p.day_id=(SELECT day_id FROM day1)
         AND p.exercise_id=(SELECT e.id FROM exercises e WHERE e.name='Press plano' LIMIT 1)
         ORDER BY p.position, p.id LIMIT 1),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (VALUES (1,62.5::numeric,12,7.5::numeric,false),
             (2,67.5::numeric,10,8.0::numeric,false),
             (3,72.5::numeric, 8,8.5::numeric,true)) AS x(set_index,weight,reps,rpe,to_failure)
WHERE EXISTS (SELECT 1 FROM new_s);

-- Ikki: Press militar (hoy 06:50)
WITH emails AS (
  SELECT 'ikki@saint.example'::text AS coach
), owner_prog AS (
  SELECT (SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
          WHERE u.email=(SELECT coach FROM emails) ORDER BY p.created_at DESC LIMIT 1) AS prog_id
), day1 AS (
  SELECT COALESCE(
           (SELECT d.id FROM program_days d JOIN program_weeks w ON w.id=d.week_id
             WHERE w.program_id=(SELECT prog_id FROM owner_prog) AND d.day_index=1
             ORDER BY d.id LIMIT 1),
           (SELECT d2.id FROM program_days d2 JOIN program_weeks w2 ON w2.id=d2.week_id
             WHERE w2.program_id=(SELECT prog_id FROM owner_prog)
             ORDER BY d2.day_index, d2.id LIMIT 1)
         ) AS day_id
), asg AS (
  SELECT (SELECT a.id FROM assignments a
            WHERE a.disciple_id=(SELECT id FROM users WHERE email=(SELECT coach FROM emails))
              AND a.program_id =(SELECT prog_id FROM owner_prog)
            ORDER BY a.created_at DESC LIMIT 1) AS asg_id
), new_s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT (SELECT asg_id FROM asg),
         (SELECT id FROM users WHERE email=(SELECT coach FROM emails)),
         (SELECT day_id FROM day1),
         (CURRENT_DATE + TIME '06:50'),
         'Ikki - amanecer'
  WHERE (SELECT day_id FROM day1) IS NOT NULL
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM new_s LIMIT 1),
       (SELECT p.id FROM prescriptions p WHERE p.day_id=(SELECT day_id FROM day1)
         AND p.exercise_id=(SELECT e.id FROM exercises e WHERE e.name='Press militar mancuernas' LIMIT 1)
         ORDER BY p.position, p.id LIMIT 1),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (VALUES (1,22.5::numeric,10,7.0::numeric,false),
             (2,25.0::numeric, 8,8.0::numeric,false),
             (3,27.5::numeric, 6,8.5::numeric,true)) AS x(set_index,weight,reps,rpe,to_failure)
WHERE EXISTS (SELECT 1 FROM new_s);

COMMIT;
