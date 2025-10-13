BEGIN;

-- ========= Password: 'secret123' (prehashed Argon2id) =========
WITH cfg AS (
  SELECT '$argon2id$v=19$m=65536,t=1,p=8$CJcnBbz7AoJRdZslhgpJxg$gTVT+6YGG2M4r87dLvbfIAY1inQx1MqP0v1/qLcKHjk'::text AS ph
)
-- ========= Usuarios =========
INSERT INTO users (email, password_hash, name)
SELECT * FROM (
  -- Coaches
  SELECT 'roshi@kamehouse.example',        (SELECT ph FROM cfg), 'Maestro Roshi' UNION ALL
  SELECT 'baki.hanma@example.example',     (SELECT ph FROM cfg), 'Baki Hanma'    UNION ALL
  SELECT 'ikki@saint.example',             (SELECT ph FROM cfg), 'Ikki (Fénix)'  UNION ALL
  -- Disciples
  SELECT 'krillin@kamehouse.example',      (SELECT ph FROM cfg), 'Krillin'       UNION ALL
  SELECT 'yamcha@capsule.example',         (SELECT ph FROM cfg), 'Yamcha'        UNION ALL
  SELECT 'goku@capsule.example',           (SELECT ph FROM cfg), 'Goku'          UNION ALL
  SELECT 'retsu@shinshinkai.example',      (SELECT ph FROM cfg), 'Retsu Kaioh'   UNION ALL
  SELECT 'katsumi@shinshinkai.example',    (SELECT ph FROM cfg), 'Katsumi Orochi' UNION ALL
  SELECT 'jack.hanma@example.example',     (SELECT ph FROM cfg), 'Jack Hanma'
) AS u(email,password_hash,name)
ON CONFLICT (email) DO NOTHING;


-- ========= Vínculos coach–disciple =========
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
INSERT INTO coach_links (coach_id, disciple_id, status)
SELECT *
FROM (
  SELECT roshi,  krillin, 'accepted'::text FROM u UNION ALL
  SELECT roshi,  yamcha,  'accepted'       FROM u UNION ALL
  SELECT roshi,  goku,    'accepted'       FROM u UNION ALL
  SELECT baki,   retsu,   'accepted'       FROM u UNION ALL
  SELECT baki,   katsumi, 'accepted'       FROM u UNION ALL
  SELECT baki,   jack,    'accepted'       FROM u UNION ALL
  SELECT ikki,   ikki,    'accepted'       FROM u
) AS x(coach_id, disciple_id, status)
WHERE x.coach_id IS NOT NULL AND x.disciple_id IS NOT NULL
ON CONFLICT (coach_id, disciple_id) DO NOTHING;

-- ========= Catálogo de ejercicios =========
INSERT INTO exercises (name, primary_muscle, equipment, tags)
SELECT v.name, v.primary_muscle, v.equipment, v.tags
FROM (
  VALUES
    ('Pec fly',                   'chest',    'machine',  ARRAY['hypertrophy']::text[]),
    ('Press plano',               'chest',    'barbell',  ARRAY['compound']::text[]),
    ('Posterior en poleas',       'shoulders','cable',    ARRAY['rear-delt']::text[]),
    ('Press militar mancuernas',  'shoulders','dumbbell', ARRAY['compound']::text[]),
    ('Vuelos laterales',          'shoulders','dumbbell', ARRAY['isolation']::text[]),
    ('Extensión tríceps en polea','triceps',  'cable',    ARRAY['isolation']::text[])
) AS v(name, primary_muscle, equipment, tags)
WHERE NOT EXISTS (SELECT 1 FROM exercises e WHERE e.name = v.name);

-- ========= Programas (Roshi / Baki / Ikki) =========
WITH owners AS (
  SELECT
    (SELECT id FROM users WHERE email='roshi@kamehouse.example')        AS roshi,
    (SELECT id FROM users WHERE email='baki.hanma@example.example')     AS baki,
    (SELECT id FROM users WHERE email='ikki@saint.example')             AS ikki
),
cand AS (
  SELECT roshi AS owner_id, 'Kame Style - Pecho/Hombro'::text AS title, 'Programa de Roshi'::text AS notes, 'private'::text AS visibility, 1 AS version FROM owners
  UNION ALL
  SELECT baki , 'Shinshinkai - Full Body',   'Programa de Baki',  'private', 1 FROM owners
  UNION ALL
  SELECT ikki , 'Fénix Solo',                'Auto Ikki',         'private', 1 FROM owners
)
INSERT INTO programs (owner_id, title, notes, visibility, version)
SELECT c.owner_id, c.title, c.notes, c.visibility, c.version
FROM cand c
WHERE c.owner_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM programs pr
    WHERE pr.owner_id = c.owner_id AND pr.title = c.title
  );

-- ========= Garantizar week 1 =========
INSERT INTO program_weeks (program_id, week_index)
SELECT p.id, 1
FROM programs p
WHERE NOT EXISTS (
  SELECT 1 FROM program_weeks w WHERE w.program_id=p.id AND w.week_index=1
);

-- ========= Garantizar day 1 =========
INSERT INTO program_days (week_id, day_index, notes)
SELECT w.id, 1,
       COALESCE(
         (SELECT CASE
             WHEN pr.title ILIKE '%Kame%'        THEN 'Pecho/Hombro'
             WHEN pr.title ILIKE '%Shinshinkai%' THEN 'Full body A'
             ELSE 'Auto día 1'
           END
          FROM programs pr WHERE pr.id = w.program_id),
         'Auto día 1'
       )
FROM program_weeks w
WHERE w.week_index=1
  AND NOT EXISTS (SELECT 1 FROM program_days d WHERE d.week_id=w.id AND d.day_index=1);

-- ========= Garantizar prescripción base en cada day 1 =========
INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, to_failure, position)
SELECT d.id,
       COALESCE(
         (SELECT id FROM exercises WHERE name IN
            ('Press plano','Pec fly','Press militar mancuernas','Extensión tríceps en polea')
          ORDER BY name LIMIT 1),
         (SELECT id FROM exercises ORDER BY id LIMIT 1)
       ),
       3, '10-12', 90, FALSE, 1
FROM program_days d
WHERE d.day_index=1
  AND NOT EXISTS (SELECT 1 FROM prescriptions pp WHERE pp.day_id=d.id);

-- ========= Assignments activos =========
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
),
pairs AS (
  -- Roshi
  SELECT (SELECT id FROM programs WHERE owner_id=(SELECT roshi FROM u) ORDER BY created_at DESC LIMIT 1) AS p_id,
         (SELECT krillin FROM u) AS d_id, (SELECT roshi FROM u) AS m_id
  UNION ALL SELECT (SELECT id FROM programs WHERE owner_id=(SELECT roshi FROM u) ORDER BY created_at DESC LIMIT 1),
                   (SELECT yamcha FROM u), (SELECT roshi FROM u)
  UNION ALL SELECT (SELECT id FROM programs WHERE owner_id=(SELECT roshi FROM u) ORDER BY created_at DESC LIMIT 1),
                   (SELECT goku   FROM u), (SELECT roshi FROM u)
  -- Baki
  UNION ALL SELECT (SELECT id FROM programs WHERE owner_id=(SELECT baki FROM u) ORDER BY created_at DESC LIMIT 1),
                   (SELECT retsu  FROM u), (SELECT baki FROM u)
  UNION ALL SELECT (SELECT id FROM programs WHERE owner_id=(SELECT baki FROM u) ORDER BY created_at DESC LIMIT 1),
                   (SELECT katsumi FROM u), (SELECT baki FROM u)
  UNION ALL SELECT (SELECT id FROM programs WHERE owner_id=(SELECT baki FROM u) ORDER BY created_at DESC LIMIT 1),
                   (SELECT jack   FROM u), (SELECT baki FROM u)
  -- Ikki (auto)
  UNION ALL SELECT (SELECT id FROM programs WHERE owner_id=(SELECT ikki FROM u) ORDER BY created_at DESC LIMIT 1),
                   (SELECT ikki   FROM u), (SELECT ikki FROM u)
)
INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
SELECT x.p_id, 1, x.d_id, x.m_id, CURRENT_DATE, TRUE
FROM pairs x
WHERE x.p_id IS NOT NULL AND x.d_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM assignments a
    WHERE a.disciple_id = x.d_id AND a.program_id = x.p_id AND a.is_active = TRUE
  );

-- ========= Sesiones + sets =========
-- Retsu — Pec fly (hace 3 días, 10:30)
WITH
coach AS (SELECT id FROM users WHERE email='baki.hanma@example.example'),
disc  AS (SELECT id FROM users WHERE email='retsu@shinshinkai.example'),
prog  AS (SELECT p.id FROM programs p WHERE p.owner_id=(SELECT id FROM coach) ORDER BY p.created_at DESC LIMIT 1),
day1  AS (
  SELECT d.id AS day_id
  FROM program_days d JOIN program_weeks w ON w.id=d.week_id
  WHERE w.program_id=(SELECT id FROM prog) AND d.day_index=1
  ORDER BY d.id LIMIT 1
),
ex    AS (SELECT id FROM exercises WHERE lower(name)=lower('Pec fly') LIMIT 1),
asg   AS (
  SELECT a.id AS asg_id FROM assignments a
  WHERE a.disciple_id=(SELECT id FROM disc) AND a.program_id=(SELECT id FROM prog)
  ORDER BY a.created_at DESC LIMIT 1
),
ins_p AS (
  INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, to_failure, position)
  SELECT (SELECT day_id FROM day1), (SELECT id FROM ex), 3, '12', 60, FALSE, 1
  WHERE NOT EXISTS (
    SELECT 1 FROM prescriptions px
    WHERE px.day_id=(SELECT day_id FROM day1) AND px.exercise_id=(SELECT id FROM ex)
  )
  RETURNING id, 1::int AS position
),
presc AS (
  SELECT id FROM (
    SELECT px.id, px.position
    FROM prescriptions px
    WHERE px.day_id=(SELECT day_id FROM day1) AND px.exercise_id=(SELECT id FROM ex)
    UNION ALL
    SELECT id, position FROM ins_p
  ) z
  ORDER BY z.position, z.id
  LIMIT 1
),
s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT (SELECT asg_id FROM asg),(SELECT id FROM disc),(SELECT day_id FROM day1),
         (CURRENT_DATE - INTERVAL '3 days') + TIME '10:30','Retsu - Pec fly'
  WHERE EXISTS (SELECT 1 FROM asg) AND EXISTS (SELECT 1 FROM day1) AND EXISTS (SELECT 1 FROM presc)
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM s), (SELECT id FROM presc),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (VALUES
        (1, 30::numeric, 15, 7.0::numeric, false),
        (2, 35::numeric, 12, 7.5::numeric, false),
        (3, 40::numeric, 10, 8.0::numeric, false)
     ) AS x(set_index, weight, reps, rpe, to_failure)
WHERE EXISTS (SELECT 1 FROM s);

-- Katsumi — Extensión tríceps en polea (hace 2 días, 19:15)
WITH
coach AS (SELECT id FROM users WHERE email='baki.hanma@example.example'),
disc  AS (SELECT id FROM users WHERE email='katsumi@shinshinkai.example'),
prog  AS (SELECT p.id FROM programs p WHERE p.owner_id=(SELECT id FROM coach) ORDER BY p.created_at DESC LIMIT 1),
day1  AS (
  SELECT d.id AS day_id
  FROM program_days d JOIN program_weeks w ON w.id=d.week_id
  WHERE w.program_id=(SELECT id FROM prog) AND d.day_index=1
  ORDER BY d.id LIMIT 1
),
ex    AS (SELECT id FROM exercises WHERE lower(name)=lower('Extensión tríceps en polea') LIMIT 1),
asg   AS (
  SELECT a.id AS asg_id FROM assignments a
  WHERE a.disciple_id=(SELECT id FROM disc) AND a.program_id=(SELECT id FROM prog)
  ORDER BY a.created_at DESC LIMIT 1
),
ins_p AS (
  INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, to_failure, position)
  SELECT (SELECT day_id FROM day1), (SELECT id FROM ex), 3, '12', 60, FALSE, 2
  WHERE NOT EXISTS (
    SELECT 1 FROM prescriptions px
    WHERE px.day_id=(SELECT day_id FROM day1) AND px.exercise_id=(SELECT id FROM ex)
  )
  RETURNING id, 2::int AS position
),
presc AS (
  SELECT id FROM (
    SELECT px.id, px.position
    FROM prescriptions px
    WHERE px.day_id=(SELECT day_id FROM day1) AND px.exercise_id=(SELECT id FROM ex)
    UNION ALL
    SELECT id, position FROM ins_p
  ) z
  ORDER BY z.position, z.id
  LIMIT 1
),
s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT (SELECT asg_id FROM asg),(SELECT id FROM disc),(SELECT day_id FROM day1),
         (CURRENT_DATE - INTERVAL '2 days') + TIME '19:15','Katsumi - Tríceps en polea'
  WHERE EXISTS (SELECT 1 FROM asg) AND EXISTS (SELECT 1 FROM day1) AND EXISTS (SELECT 1 FROM presc)
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM s), (SELECT id FROM presc),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (VALUES
        (1, 25::numeric, 15, 7.0::numeric, false),
        (2, 30::numeric, 12, 7.5::numeric, false),
        (3, 35::numeric, 12, 8.0::numeric, true)
     ) AS x(set_index, weight, reps, rpe, to_failure)
WHERE EXISTS (SELECT 1 FROM s);

-- Jack — (ayer 08:05) Pec fly
WITH
coach AS (SELECT id FROM users WHERE email='baki.hanma@example.example'),
disc  AS (SELECT id FROM users WHERE email='jack.hanma@example.example'),
prog  AS (SELECT p.id FROM programs p WHERE p.owner_id=(SELECT id FROM coach) ORDER BY p.created_at DESC LIMIT 1),
day1  AS (
  SELECT d.id AS day_id
  FROM program_days d JOIN program_weeks w ON w.id=d.week_id
  WHERE w.program_id=(SELECT id FROM prog) AND d.day_index=1
  ORDER BY d.id LIMIT 1
),
ex    AS (SELECT id FROM exercises WHERE lower(name)=lower('Pec fly') LIMIT 1),
asg   AS (
  SELECT a.id AS asg_id FROM assignments a
  WHERE a.disciple_id=(SELECT id FROM disc) AND a.program_id=(SELECT id FROM prog)
  ORDER BY a.created_at DESC LIMIT 1
),
ins_p AS (
  INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, to_failure, position)
  SELECT (SELECT day_id FROM day1), (SELECT id FROM ex), 3, '12', 60, FALSE, 1
  WHERE NOT EXISTS (
    SELECT 1 FROM prescriptions px
    WHERE px.day_id=(SELECT day_id FROM day1) AND px.exercise_id=(SELECT id FROM ex)
  )
  RETURNING id, 1::int AS position
),
presc AS (
  SELECT id FROM (
    SELECT px.id, px.position
    FROM prescriptions px
    WHERE px.day_id=(SELECT day_id FROM day1) AND px.exercise_id=(SELECT id FROM ex)
    UNION ALL
    SELECT id, position FROM ins_p
  ) z
  ORDER BY z.position, z.id
  LIMIT 1
),
s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT (SELECT asg_id FROM asg),(SELECT id FROM disc),(SELECT day_id FROM day1),
         (CURRENT_DATE - INTERVAL '1 day') + TIME '08:05','Jack - Full body A'
  WHERE EXISTS (SELECT 1 FROM asg) AND EXISTS (SELECT 1 FROM day1) AND EXISTS (SELECT 1 FROM presc)
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM s), (SELECT id FROM presc),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (VALUES
        (1, 50::numeric, 10, 7.5::numeric, false),
        (2, 55::numeric,  8, 8.0::numeric, false),
        (3, 60::numeric,  6, 8.5::numeric, true)
     ) AS x(set_index, weight, reps, rpe, to_failure)
WHERE EXISTS (SELECT 1 FROM s);

-- Jack — MISMA sesión, tríceps (esta parte ya era válida)
WITH sess AS (
  SELECT id
  FROM session_logs
  WHERE notes = 'Jack - Full body A'
  ORDER BY performed_at DESC
  LIMIT 1
),
presc AS (
  SELECT px.id
  FROM prescriptions px
  JOIN program_days d ON d.id=px.day_id
  JOIN program_weeks w ON w.id=d.week_id
  JOIN programs pr ON pr.id=w.program_id
  JOIN users cu ON cu.id=pr.owner_id
  JOIN exercises e ON e.id=px.exercise_id
  WHERE cu.email='baki.hanma@example.example'
    AND d.day_index=1
    AND lower(e.name)=lower('Extensión tríceps en polea')
  ORDER BY px.position, px.id
  LIMIT 1
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM sess), (SELECT id FROM presc),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (VALUES
        (1, 35::numeric, 12, 7.0::numeric, false),
        (2, 40::numeric, 10, 7.5::numeric, false),
        (3, 45::numeric,  8, 8.5::numeric, true)
     ) AS x(set_index, weight, reps, rpe, to_failure)
WHERE EXISTS (SELECT 1 FROM sess) AND EXISTS (SELECT 1 FROM presc);

-- Krillin — Press plano (hace 4 días 18:40)
WITH
coach AS (SELECT id FROM users WHERE email='roshi@kamehouse.example'),
disc  AS (SELECT id FROM users WHERE email='krillin@kamehouse.example'),
prog  AS (SELECT p.id FROM programs p WHERE p.owner_id=(SELECT id FROM coach) ORDER BY p.created_at DESC LIMIT 1),
day1  AS (
  SELECT d.id AS day_id
  FROM program_days d JOIN program_weeks w ON w.id=d.week_id
  WHERE w.program_id=(SELECT id FROM prog) AND d.day_index=1
  ORDER BY d.id LIMIT 1
),
ex    AS (SELECT id FROM exercises WHERE lower(name)=lower('Press plano') LIMIT 1),
asg   AS (
  SELECT a.id AS asg_id FROM assignments a
  WHERE a.disciple_id=(SELECT id FROM disc) AND a.program_id=(SELECT id FROM prog)
  ORDER BY a.created_at DESC LIMIT 1
),
ins_p AS (
  INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, to_failure, position)
  SELECT (SELECT day_id FROM day1), (SELECT id FROM ex), 4, '10', 90, FALSE, 1
  WHERE NOT EXISTS (
    SELECT 1 FROM prescriptions px
    WHERE px.day_id=(SELECT day_id FROM day1) AND px.exercise_id=(SELECT id FROM ex)
  )
  RETURNING id, 1::int AS position
),
presc AS (
  SELECT id FROM (
    SELECT px.id, px.position
    FROM prescriptions px
    WHERE px.day_id=(SELECT day_id FROM day1) AND px.exercise_id=(SELECT id FROM ex)
    UNION ALL
    SELECT id, position FROM ins_p
  ) z
  ORDER BY z.position, z.id
  LIMIT 1
),
s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT (SELECT asg_id FROM asg),(SELECT id FROM disc),(SELECT day_id FROM day1),
         (CURRENT_DATE - INTERVAL '4 days') + TIME '18:40','Krillin - extra pecho'
  WHERE EXISTS (SELECT 1 FROM asg) AND EXISTS (SELECT 1 FROM day1) AND EXISTS (SELECT 1 FROM presc)
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM s), (SELECT id FROM presc),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (VALUES
        (1, 62.5::numeric, 12, 7.5::numeric, false),
        (2, 67.5::numeric, 10, 8.0::numeric, false),
        (3, 72.5::numeric,  8, 8.5::numeric, true)
     ) AS x(set_index, weight, reps, rpe, to_failure)
WHERE EXISTS (SELECT 1 FROM s);

-- Ikki — Press militar mancuernas (HOY 06:50)
WITH
coach AS (SELECT id FROM users WHERE email='ikki@saint.example'),
disc  AS (SELECT id FROM users WHERE email='ikki@saint.example'),
prog  AS (SELECT p.id FROM programs p WHERE p.owner_id=(SELECT id FROM coach) ORDER BY p.created_at DESC LIMIT 1),
day1  AS (
  SELECT d.id AS day_id
  FROM program_days d JOIN program_weeks w ON w.id=d.week_id
  WHERE w.program_id=(SELECT id FROM prog) AND d.day_index=1
  ORDER BY d.id LIMIT 1
),
ex    AS (SELECT id FROM exercises WHERE lower(name)=lower('Press militar mancuernas') LIMIT 1),
asg   AS (
  SELECT a.id AS asg_id FROM assignments a
  WHERE a.disciple_id=(SELECT id FROM disc) AND a.program_id=(SELECT id FROM prog)
  ORDER BY a.created_at DESC LIMIT 1
),
ins_p AS (
  INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, to_failure, position)
  SELECT (SELECT day_id FROM day1), (SELECT id FROM ex), 4, '8-10', 90, FALSE, 1
  WHERE NOT EXISTS (
    SELECT 1 FROM prescriptions px
    WHERE px.day_id=(SELECT day_id FROM day1) AND px.exercise_id=(SELECT id FROM ex)
  )
  RETURNING id, 1::int AS position
),
presc AS (
  SELECT id FROM (
    SELECT px.id, px.position
    FROM prescriptions px
    WHERE px.day_id=(SELECT day_id FROM day1) AND px.exercise_id=(SELECT id FROM ex)
    UNION ALL
    SELECT id, position FROM ins_p
  ) z
  ORDER BY z.position, z.id
  LIMIT 1
),
s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT (SELECT asg_id FROM asg),(SELECT id FROM disc),(SELECT day_id FROM day1),
         (CURRENT_DATE + TIME '06:50'),'Ikki - amanecer'
  WHERE EXISTS (SELECT 1 FROM asg) AND EXISTS (SELECT 1 FROM day1) AND EXISTS (SELECT 1 FROM presc)
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM s), (SELECT id FROM presc),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM (VALUES
        (1, 22.5::numeric, 10, 7.0::numeric, false),
        (2, 25.0::numeric,  8, 8.0::numeric, false),
        (3, 27.5::numeric,  6, 8.5::numeric, true)
     ) AS x(set_index, weight, reps, rpe, to_failure)
WHERE EXISTS (SELECT 1 FROM s);

COMMIT;
