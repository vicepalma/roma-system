-- ===========================================
-- 011_extra_sessions_by_email.sql (CTEs por bloque, sin reuse entre sentencias)
-- ===========================================

-- ============ Retsu: Pec Fly (hace 3 días 10:30) ============
WITH emails AS (
  SELECT
    'baki.hanma@example.example'::text AS baki_email,
    'retsu@shinshinkai.example'::text  AS retsu_email
),
owners AS (
  SELECT (SELECT u.id FROM users u WHERE u.email=(SELECT baki_email FROM emails)) AS baki_id
),
programs_by_owner AS (
  SELECT (SELECT p.id FROM programs p JOIN owners o ON o.baki_id=p.owner_id ORDER BY p.created_at DESC LIMIT 1) AS baki_prog
),
days_by_owner AS (
  SELECT COALESCE(
           (SELECT d.id FROM program_days d JOIN program_weeks w ON w.id=d.week_id
             WHERE w.program_id=(SELECT baki_prog FROM programs_by_owner) AND d.day_index=1
             ORDER BY d.id LIMIT 1),
           (SELECT d2.id FROM program_days d2 JOIN program_weeks w2 ON w2.id=d2.week_id
             WHERE w2.program_id=(SELECT baki_prog FROM programs_by_owner)
             ORDER BY d2.day_index, d2.id LIMIT 1)
         ) AS baki_day
),
disciple_ids AS (
  SELECT (SELECT u.id FROM users u WHERE u.email=(SELECT retsu_email FROM emails)) AS retsu_id
),
assignments_resolved AS (
  SELECT (SELECT a.id FROM assignments a
            WHERE a.disciple_id=(SELECT retsu_id FROM disciple_ids)
              AND a.program_id =(SELECT baki_prog FROM programs_by_owner)
            ORDER BY a.created_at DESC LIMIT 1) AS asg_retsu
),
prescriptions_resolved AS (
  SELECT (SELECT p.id FROM prescriptions p
            WHERE p.day_id=(SELECT baki_day FROM days_by_owner)
              AND p.exercise_id=(SELECT e.id FROM exercises e WHERE e.name='Pec fly' LIMIT 1)
            ORDER BY p.position, p.id LIMIT 1) AS presc_baki_pec
),
new_s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT
    (SELECT asg_retsu FROM assignments_resolved),
    (SELECT retsu_id  FROM disciple_ids),
    (SELECT baki_day  FROM days_by_owner),
    (CURRENT_DATE - INTERVAL '3 days') + TIME '10:30',
    'Retsu - Pec fly'
  WHERE (SELECT asg_retsu FROM assignments_resolved) IS NOT NULL
    AND (SELECT baki_day  FROM days_by_owner)       IS NOT NULL
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT
  (SELECT id FROM new_s),
  (SELECT presc_baki_pec FROM prescriptions_resolved),
  x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM ( VALUES
  (1,30::numeric,15,7.0::numeric,false),
  (2,35::numeric,12,7.5::numeric,false),
  (3,40::numeric,10,8.0::numeric,false)
) AS x(set_index,weight,reps,rpe,to_failure)
WHERE (SELECT id FROM new_s) IS NOT NULL
  AND (SELECT presc_baki_pec FROM prescriptions_resolved) IS NOT NULL;

-- ============ Katsumi: Tríceps en polea (hace 2 días 19:15) ============
WITH emails AS (
  SELECT
    'baki.hanma@example.example'::text AS baki_email,
    'katsumi@shinshinkai.example'::text AS katsumi_email
),
owners AS (
  SELECT (SELECT u.id FROM users u WHERE u.email=(SELECT baki_email FROM emails)) AS baki_id
),
programs_by_owner AS (
  SELECT (SELECT p.id FROM programs p JOIN owners o ON o.baki_id=p.owner_id ORDER BY p.created_at DESC LIMIT 1) AS baki_prog
),
days_by_owner AS (
  SELECT COALESCE(
           (SELECT d.id FROM program_days d JOIN program_weeks w ON w.id=d.week_id
             WHERE w.program_id=(SELECT baki_prog FROM programs_by_owner) AND d.day_index=1
             ORDER BY d.id LIMIT 1),
           (SELECT d2.id FROM program_days d2 JOIN program_weeks w2 ON w2.id=d2.week_id
             WHERE w2.program_id=(SELECT baki_prog FROM programs_by_owner)
             ORDER BY d2.day_index, d2.id LIMIT 1)
         ) AS baki_day
),
disciple_ids AS (
  SELECT (SELECT u.id FROM users u WHERE u.email=(SELECT katsumi_email FROM emails)) AS katsumi_id
),
assignments_resolved AS (
  SELECT (SELECT a.id FROM assignments a
            WHERE a.disciple_id=(SELECT katsumi_id FROM disciple_ids)
              AND a.program_id =(SELECT baki_prog FROM programs_by_owner)
            ORDER BY a.created_at DESC LIMIT 1) AS asg_katsumi
),
prescriptions_resolved AS (
  SELECT (SELECT p.id FROM prescriptions p
            WHERE p.day_id=(SELECT baki_day FROM days_by_owner)
              AND p.exercise_id=(SELECT e.id FROM exercises e WHERE e.name='Extensión tríceps en polea' LIMIT 1)
            ORDER BY p.position, p.id LIMIT 1) AS presc_baki_triceps
),
new_s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT
    (SELECT asg_katsumi FROM assignments_resolved),
    (SELECT katsumi_id  FROM disciple_ids),
    (SELECT baki_day    FROM days_by_owner),
    (CURRENT_DATE - INTERVAL '2 days') + TIME '19:15',
    'Katsumi - Tríceps en polea'
  WHERE (SELECT asg_katsumi FROM assignments_resolved) IS NOT NULL
    AND (SELECT baki_day    FROM days_by_owner)        IS NOT NULL
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT
  (SELECT id FROM new_s),
  (SELECT presc_baki_triceps FROM prescriptions_resolved),
  x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM ( VALUES
  (1,25::numeric,15,7.0::numeric,false),
  (2,30::numeric,12,7.5::numeric,false),
  (3,35::numeric,12,8.0::numeric,true)
) AS x(set_index,weight,reps,rpe,to_failure)
WHERE (SELECT id FROM new_s) IS NOT NULL
  AND (SELECT presc_baki_triceps FROM prescriptions_resolved) IS NOT NULL;

-- ============ Jack: Full body A (ayer 08:05) – Pec fly ============
WITH emails AS (
  SELECT
    'baki.hanma@example.example'::text AS baki_email,
    'jack.hanma@example.example'::text AS jack_email
),
owners AS (
  SELECT (SELECT u.id FROM users u WHERE u.email=(SELECT baki_email FROM emails)) AS baki_id
),
programs_by_owner AS (
  SELECT (SELECT p.id FROM programs p JOIN owners o ON o.baki_id=p.owner_id ORDER BY p.created_at DESC LIMIT 1) AS baki_prog
),
days_by_owner AS (
  SELECT COALESCE(
           (SELECT d.id FROM program_days d JOIN program_weeks w ON w.id=d.week_id
             WHERE w.program_id=(SELECT baki_prog FROM programs_by_owner) AND d.day_index=1
             ORDER BY d.id LIMIT 1),
           (SELECT d2.id FROM program_days d2 JOIN program_weeks w2 ON w2.id=d2.week_id
             WHERE w2.program_id=(SELECT baki_prog FROM programs_by_owner)
             ORDER BY d2.day_index, d2.id LIMIT 1)
         ) AS baki_day
),
disciple_ids AS (
  SELECT (SELECT u.id FROM users u WHERE u.email=(SELECT jack_email FROM emails)) AS jack_id
),
assignments_resolved AS (
  SELECT (SELECT a.id FROM assignments a
            WHERE a.disciple_id=(SELECT jack_id FROM disciple_ids)
              AND a.program_id =(SELECT baki_prog FROM programs_by_owner)
            ORDER BY a.created_at DESC LIMIT 1) AS asg_jack
),
prescriptions_resolved AS (
  SELECT
    (SELECT p.id FROM prescriptions p
      WHERE p.day_id=(SELECT baki_day FROM days_by_owner)
        AND p.exercise_id=(SELECT e.id FROM exercises e WHERE e.name='Pec fly' LIMIT 1)
      ORDER BY p.position, p.id LIMIT 1) AS presc_baki_pec
),
new_s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT
    (SELECT asg_jack FROM assignments_resolved),
    (SELECT jack_id  FROM disciple_ids),
    (SELECT baki_day FROM days_by_owner),
    (CURRENT_DATE - INTERVAL '1 day') + TIME '08:05',
    'Jack - Full body A'
  WHERE (SELECT asg_jack FROM assignments_resolved) IS NOT NULL
    AND (SELECT baki_day FROM days_by_owner)        IS NOT NULL
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM new_s), (SELECT presc_baki_pec FROM prescriptions_resolved),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM ( VALUES
  (1,50::numeric,10,7.5::numeric,false),
  (2,55::numeric, 8,8.0::numeric,false),
  (3,60::numeric, 6,8.5::numeric,true)
) AS x(set_index,weight,reps,rpe,to_failure)
WHERE (SELECT id FROM new_s) IS NOT NULL
  AND (SELECT presc_baki_pec FROM prescriptions_resolved) IS NOT NULL;

-- Jack: Extensión tríceps en polea (MISMA SESIÓN) → subselect inline, sin CTE previas
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT
  s.id,
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
FROM session_logs s
CROSS JOIN ( VALUES
  (1,35::numeric,12,7.0::numeric,false),
  (2,40::numeric,10,7.5::numeric,false),
  (3,45::numeric, 8,8.5::numeric,true)
) AS x(set_index,weight,reps,rpe,to_failure)
WHERE s.notes = 'Jack - Full body A'
ORDER BY s.performed_at DESC
LIMIT 3;

-- ============ Ikki: Press militar (hoy 06:50) ============
WITH emails AS (
  SELECT 'ikki@saint.example'::text AS ikki_email
),
owners AS (
  SELECT (SELECT u.id FROM users u WHERE u.email=(SELECT ikki_email FROM emails)) AS ikki_id
),
programs_by_owner AS (
  SELECT (SELECT p.id FROM programs p JOIN owners o ON o.ikki_id=p.owner_id ORDER BY p.created_at DESC LIMIT 1) AS ikki_prog
),
days_by_owner AS (
  SELECT COALESCE(
           (SELECT d.id FROM program_days d JOIN program_weeks w ON w.id=d.week_id
             WHERE w.program_id=(SELECT ikki_prog FROM programs_by_owner) AND d.day_index=1
             ORDER BY d.id LIMIT 1),
           (SELECT d2.id FROM program_days d2 JOIN program_weeks w2 ON w2.id=d2.week_id
             WHERE w2.program_id=(SELECT ikki_prog FROM programs_by_owner)
             ORDER BY d2.day_index, d2.id LIMIT 1)
         ) AS ikki_day
),
disciple_ids AS (
  SELECT (SELECT u.id FROM users u WHERE u.email=(SELECT ikki_email FROM emails)) AS ikki_id
),
assignments_resolved AS (
  SELECT (SELECT a.id FROM assignments a
            WHERE a.disciple_id=(SELECT ikki_id FROM disciple_ids)
              AND a.program_id =(SELECT ikki_prog FROM programs_by_owner)
            ORDER BY a.created_at DESC LIMIT 1) AS asg_ikki
),
prescriptions_resolved AS (
  SELECT (SELECT p.id FROM prescriptions p
            WHERE p.day_id=(SELECT ikki_day FROM days_by_owner)
              AND p.exercise_id=(SELECT e.id FROM exercises e WHERE e.name='Press militar mancuernas' LIMIT 1)
            ORDER BY p.position, p.id LIMIT 1) AS presc_ikki_press
),
new_s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT
    (SELECT asg_ikki FROM assignments_resolved),
    (SELECT ikki_id FROM disciple_ids),
    (SELECT ikki_day FROM days_by_owner),
    (CURRENT_DATE + TIME '06:50'),
    'Ikki - amanecer'
  WHERE (SELECT asg_ikki FROM assignments_resolved) IS NOT NULL
    AND (SELECT ikki_day FROM days_by_owner)       IS NOT NULL
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM new_s), (SELECT presc_ikki_press FROM prescriptions_resolved),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM ( VALUES
  (1,22.5::numeric,10,7.0::numeric,false),
  (2,25.0::numeric, 8,8.0::numeric,false),
  (3,27.5::numeric, 6,8.5::numeric,true)
) AS x(set_index,weight,reps,rpe,to_failure)
WHERE (SELECT id FROM new_s) IS NOT NULL
  AND (SELECT presc_ikki_press FROM prescriptions_resolved) IS NOT NULL;

-- ============ Krillin: Press plano (hace 4 días 18:40) ============
WITH emails AS (
  SELECT
    'roshi@kamehouse.example'::text   AS roshi_email,
    'krillin@kamehouse.example'::text AS krillin_email
),
owners AS (
  SELECT (SELECT u.id FROM users u WHERE u.email=(SELECT roshi_email FROM emails)) AS roshi_id
),
programs_by_owner AS (
  SELECT (SELECT p.id FROM programs p JOIN owners o ON o.roshi_id=p.owner_id ORDER BY p.created_at DESC LIMIT 1) AS roshi_prog
),
days_by_owner AS (
  SELECT COALESCE(
           (SELECT d.id FROM program_days d JOIN program_weeks w ON w.id=d.week_id
             WHERE w.program_id=(SELECT roshi_prog FROM programs_by_owner) AND d.day_index=1
             ORDER BY d.id LIMIT 1),
           (SELECT d2.id FROM program_days d2 JOIN program_weeks w2 ON w2.id=d2.week_id
             WHERE w2.program_id=(SELECT roshi_prog FROM programs_by_owner)
             ORDER BY d2.day_index, d2.id LIMIT 1)
         ) AS roshi_day
),
disciple_ids AS (
  SELECT (SELECT u.id FROM users u WHERE u.email=(SELECT krillin_email FROM emails)) AS krillin_id
),
assignments_resolved AS (
  SELECT (SELECT a.id FROM assignments a
            WHERE a.disciple_id=(SELECT krillin_id FROM disciple_ids)
              AND a.program_id =(SELECT roshi_prog FROM programs_by_owner)
            ORDER BY a.created_at DESC LIMIT 1) AS asg_krillin
),
prescriptions_resolved AS (
  SELECT (SELECT p.id FROM prescriptions p
            WHERE p.day_id=(SELECT roshi_day FROM days_by_owner)
              AND p.exercise_id=(SELECT e.id FROM exercises e WHERE e.name='Press plano' LIMIT 1)
            ORDER BY p.position, p.id LIMIT 1) AS presc_roshi_press
),
new_s AS (
  INSERT INTO session_logs (assignment_id, disciple_id, day_id, performed_at, notes)
  SELECT
    (SELECT asg_krillin FROM assignments_resolved),
    (SELECT krillin_id FROM disciple_ids),
    (SELECT roshi_day  FROM days_by_owner),
    (CURRENT_DATE - INTERVAL '4 days') + TIME '18:40',
    'Krillin - extra pecho'
  WHERE (SELECT asg_krillin FROM assignments_resolved) IS NOT NULL
    AND (SELECT roshi_day  FROM days_by_owner)       IS NOT NULL
  RETURNING id
)
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure)
SELECT (SELECT id FROM new_s), (SELECT presc_roshi_press FROM prescriptions_resolved),
       x.set_index, x.weight, x.reps, x.rpe, x.to_failure
FROM ( VALUES
  (1,62.5::numeric,12,7.5::numeric,false),
  (2,67.5::numeric,10,8.0::numeric,false),
  (3,72.5::numeric, 8,8.5::numeric,true)
) AS x(set_index,weight,reps,rpe,to_failure)
WHERE (SELECT id FROM new_s) IS NOT NULL
  AND (SELECT presc_roshi_press FROM prescriptions_resolved) IS NOT NULL;
