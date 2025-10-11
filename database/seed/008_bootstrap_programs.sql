-- ===========================================
-- 008_bootstrap_programs.sql  (SIN CTEs)
-- Crea/asegura programas, semana/día, y
-- prescriptions mínimas para:
--  - Baki: Pec fly + Extensión tríceps en polea
--  - Roshi: Press plano
--  - Ikki: Press militar mancuernas
-- Idempotente (no duplica).
-- ===========================================

-- ====== Ejercicios base (por si faltan) ======
INSERT INTO exercises (name, primary_muscle, equipment, tags)
SELECT 'Pec fly','chest','machine', ARRAY['hypertrophy']
WHERE NOT EXISTS (SELECT 1 FROM exercises WHERE lower(name)=lower('Pec fly'));

INSERT INTO exercises (name, primary_muscle, equipment, tags)
SELECT 'Extensión tríceps en polea','triceps','cable', ARRAY['isolation']
WHERE NOT EXISTS (SELECT 1 FROM exercises WHERE lower(name)=lower('Extensión tríceps en polea'));

INSERT INTO exercises (name, primary_muscle, equipment, tags)
SELECT 'Press militar mancuernas','shoulders','dumbbell', ARRAY['compound']
WHERE NOT EXISTS (SELECT 1 FROM exercises WHERE lower(name)=lower('Press militar mancuernas'));

INSERT INTO exercises (name, primary_muscle, equipment, tags)
SELECT 'Press plano','chest','barbell', ARRAY['compound']
WHERE NOT EXISTS (SELECT 1 FROM exercises WHERE lower(name)=lower('Press plano'));

-- ====== Asegurar 1 programa por owner (si no existe) ======
-- Baki
INSERT INTO programs (owner_id, title, notes, visibility, version)
SELECT u.id, 'Baki Split', 'Seed auto', 'private', 1
FROM users u
WHERE u.email='baki.hanma@example.example'
  AND NOT EXISTS (SELECT 1 FROM programs p WHERE p.owner_id=u.id);

-- Roshi
INSERT INTO programs (owner_id, title, notes, visibility, version)
SELECT u.id, 'Roshi Classic', 'Seed auto', 'private', 1
FROM users u
WHERE u.email='roshi@kamehouse.example'
  AND NOT EXISTS (SELECT 1 FROM programs p WHERE p.owner_id=u.id);

-- Ikki
INSERT INTO programs (owner_id, title, notes, visibility, version)
SELECT u.id, 'Ikki Dawn', 'Seed auto', 'private', 1
FROM users u
WHERE u.email='ikki@saint.example'
  AND NOT EXISTS (SELECT 1 FROM programs p WHERE p.owner_id=u.id);

-- ====== Asegurar semana (week_index=1) por programa vigente ======
-- Baki
INSERT INTO program_weeks (program_id, week_index)
SELECT
  (SELECT p.id FROM programs p
     JOIN users u ON u.id=p.owner_id
    WHERE u.email='baki.hanma@example.example'
    ORDER BY p.created_at DESC
    LIMIT 1),
  1
WHERE NOT EXISTS (
  SELECT 1 FROM program_weeks w
  WHERE w.program_id = (SELECT p.id FROM programs p
                          JOIN users u ON u.id=p.owner_id
                         WHERE u.email='baki.hanma@example.example'
                         ORDER BY p.created_at DESC LIMIT 1)
);

-- Roshi
INSERT INTO program_weeks (program_id, week_index)
SELECT
  (SELECT p.id FROM programs p
     JOIN users u ON u.id=p.owner_id
    WHERE u.email='roshi@kamehouse.example'
    ORDER BY p.created_at DESC
    LIMIT 1),
  1
WHERE NOT EXISTS (
  SELECT 1 FROM program_weeks w
  WHERE w.program_id = (SELECT p.id FROM programs p
                          JOIN users u ON u.id=p.owner_id
                         WHERE u.email='roshi@kamehouse.example'
                         ORDER BY p.created_at DESC LIMIT 1)
);

-- Ikki
INSERT INTO program_weeks (program_id, week_index)
SELECT
  (SELECT p.id FROM programs p
     JOIN users u ON u.id=p.owner_id
    WHERE u.email='ikki@saint.example'
    ORDER BY p.created_at DESC
    LIMIT 1),
  1
WHERE NOT EXISTS (
  SELECT 1 FROM program_weeks w
  WHERE w.program_id = (SELECT p.id FROM programs p
                          JOIN users u ON u.id=p.owner_id
                         WHERE u.email='ikki@saint.example'
                         ORDER BY p.created_at DESC LIMIT 1)
);

-- ====== Asegurar al menos 1 día (preferir day_index=1) ======
-- Baki
INSERT INTO program_days (week_id, day_index, notes)
SELECT
  (SELECT w.id FROM program_weeks w
   WHERE w.program_id = (SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                         WHERE u.email='baki.hanma@example.example'
                         ORDER BY p.created_at DESC LIMIT 1)
   ORDER BY w.week_index ASC LIMIT 1),
  1, 'Auto day'
WHERE NOT EXISTS (
  SELECT 1 FROM program_days d
  JOIN program_weeks w ON w.id=d.week_id
  WHERE w.program_id = (SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                        WHERE u.email='baki.hanma@example.example'
                        ORDER BY p.created_at DESC LIMIT 1)
);

-- Roshi
INSERT INTO program_days (week_id, day_index, notes)
SELECT
  (SELECT w.id FROM program_weeks w
   WHERE w.program_id = (SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                         WHERE u.email='roshi@kamehouse.example'
                         ORDER BY p.created_at DESC LIMIT 1)
   ORDER BY w.week_index ASC LIMIT 1),
  1, 'Auto day'
WHERE NOT EXISTS (
  SELECT 1 FROM program_days d
  JOIN program_weeks w ON w.id=d.week_id
  WHERE w.program_id = (SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                        WHERE u.email='roshi@kamehouse.example'
                        ORDER BY p.created_at DESC LIMIT 1)
);

-- Ikki
INSERT INTO program_days (week_id, day_index, notes)
SELECT
  (SELECT w.id FROM program_weeks w
   WHERE w.program_id = (SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                         WHERE u.email='ikki@saint.example'
                         ORDER BY p.created_at DESC LIMIT 1)
   ORDER BY w.week_index ASC LIMIT 1),
  1, 'Auto day'
WHERE NOT EXISTS (
  SELECT 1 FROM program_days d
  JOIN program_weeks w ON w.id=d.week_id
  WHERE w.program_id = (SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                        WHERE u.email='ikki@saint.example'
                        ORDER BY p.created_at DESC LIMIT 1)
);

-- ====== Prescriptions mínimas por owner/ejercicio ======
-- Baki: Pec fly
INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, position)
SELECT
  -- day preferente: day_index=1; si no, primer day
  COALESCE(
    (SELECT d.id FROM program_days d
       JOIN program_weeks w ON w.id=d.week_id
      WHERE w.program_id=(SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                           WHERE u.email='baki.hanma@example.example'
                           ORDER BY p.created_at DESC LIMIT 1)
        AND d.day_index=1
      ORDER BY d.id LIMIT 1),
    (SELECT d2.id FROM program_days d2
       JOIN program_weeks w2 ON w2.id=d2.week_id
      WHERE w2.program_id=(SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                           WHERE u.email='baki.hanma@example.example'
                           ORDER BY p.created_at DESC LIMIT 1)
      ORDER BY d2.day_index, d2.id LIMIT 1)
  ),
  (SELECT id FROM exercises WHERE lower(name)=lower('Pec fly') LIMIT 1),
  3, '10-12', 90, 1
WHERE EXISTS (SELECT 1 FROM programs p JOIN users u ON u.id=p.owner_id WHERE u.email='baki.hanma@example.example')
  AND NOT EXISTS (
    SELECT 1 FROM prescriptions p
    JOIN program_days d ON d.id=p.day_id
    JOIN program_weeks w ON w.id=d.week_id
    WHERE w.program_id=(SELECT p2.id FROM programs p2 JOIN users u2 ON u2.id=p2.owner_id
                        WHERE u2.email='baki.hanma@example.example'
                        ORDER BY p2.created_at DESC LIMIT 1)
      AND p.exercise_id=(SELECT id FROM exercises WHERE lower(name)=lower('Pec fly') LIMIT 1)
  );

-- Baki: Extensión tríceps en polea
INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, position)
SELECT
  COALESCE(
    (SELECT d.id FROM program_days d
       JOIN program_weeks w ON w.id=d.week_id
      WHERE w.program_id=(SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                           WHERE u.email='baki.hanma@example.example'
                           ORDER BY p.created_at DESC LIMIT 1)
        AND d.day_index=1
      ORDER BY d.id LIMIT 1),
    (SELECT d2.id FROM program_days d2
       JOIN program_weeks w2 ON w2.id=d2.week_id
      WHERE w2.program_id=(SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                           WHERE u.email='baki.hanma@example.example'
                           ORDER BY p.created_at DESC LIMIT 1)
      ORDER BY d2.day_index, d2.id LIMIT 1)
  ),
  (SELECT id FROM exercises WHERE lower(name)=lower('Extensión tríceps en polea') LIMIT 1),
  3, '10-12', 90, 2
WHERE EXISTS (SELECT 1 FROM programs p JOIN users u ON u.id=p.owner_id WHERE u.email='baki.hanma@example.example')
  AND NOT EXISTS (
    SELECT 1 FROM prescriptions p
    JOIN program_days d ON d.id=p.day_id
    JOIN program_weeks w ON w.id=d.week_id
    WHERE w.program_id=(SELECT p2.id FROM programs p2 JOIN users u2 ON u2.id=p2.owner_id
                        WHERE u2.email='baki.hanma@example.example'
                        ORDER BY p2.created_at DESC LIMIT 1)
      AND p.exercise_id=(SELECT id FROM exercises WHERE lower(name)=lower('Extensión tríceps en polea') LIMIT 1)
  );

-- Roshi: Press plano
INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, position)
SELECT
  COALESCE(
    (SELECT d.id FROM program_days d
       JOIN program_weeks w ON w.id=d.week_id
      WHERE w.program_id=(SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                           WHERE u.email='roshi@kamehouse.example'
                           ORDER BY p.created_at DESC LIMIT 1)
        AND d.day_index=1
      ORDER BY d.id LIMIT 1),
    (SELECT d2.id FROM program_days d2
       JOIN program_weeks w2 ON w2.id=d2.week_id
      WHERE w2.program_id=(SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                           WHERE u.email='roshi@kamehouse.example'
                           ORDER BY p.created_at DESC LIMIT 1)
      ORDER BY d2.day_index, d2.id LIMIT 1)
  ),
  (SELECT id FROM exercises WHERE lower(name)=lower('Press plano') LIMIT 1),
  3, '8-10', 120, 1
WHERE EXISTS (SELECT 1 FROM programs p JOIN users u ON u.id=p.owner_id WHERE u.email='roshi@kamehouse.example')
  AND NOT EXISTS (
    SELECT 1 FROM prescriptions p
    JOIN program_days d ON d.id=p.day_id
    JOIN program_weeks w ON w.id=d.week_id
    WHERE w.program_id=(SELECT p2.id FROM programs p2 JOIN users u2 ON u2.id=p2.owner_id
                        WHERE u2.email='roshi@kamehouse.example'
                        ORDER BY p2.created_at DESC LIMIT 1)
      AND p.exercise_id=(SELECT id FROM exercises WHERE lower(name)=lower('Press plano') LIMIT 1)
  );

-- Ikki: Press militar mancuernas
INSERT INTO prescriptions (day_id, exercise_id, series, reps, rest_sec, position)
SELECT
  COALESCE(
    (SELECT d.id FROM program_days d
       JOIN program_weeks w ON w.id=d.week_id
      WHERE w.program_id=(SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                           WHERE u.email='ikki@saint.example'
                           ORDER BY p.created_at DESC LIMIT 1)
        AND d.day_index=1
      ORDER BY d.id LIMIT 1),
    (SELECT d2.id FROM program_days d2
       JOIN program_weeks w2 ON w2.id=d2.week_id
      WHERE w2.program_id=(SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id
                           WHERE u.email='ikki@saint.example'
                           ORDER BY p.created_at DESC LIMIT 1)
      ORDER BY d2.day_index, d2.id LIMIT 1)
  ),
  (SELECT id FROM exercises WHERE lower(name)=lower('Press militar mancuernas') LIMIT 1),
  3, '8-10', 120, 1
WHERE EXISTS (SELECT 1 FROM programs p JOIN users u ON u.id=p.owner_id WHERE u.email='ikki@saint.example')
  AND NOT EXISTS (
    SELECT 1 FROM prescriptions p
    JOIN program_days d ON d.id=p.day_id
    JOIN program_weeks w ON w.id=d.week_id
    WHERE w.program_id=(SELECT p2.id FROM programs p2 JOIN users u2 ON u2.id=p2.owner_id
                        WHERE u2.email='ikki@saint.example'
                        ORDER BY p2.created_at DESC LIMIT 1)
      AND p.exercise_id=(SELECT id FROM exercises WHERE lower(name)=lower('Press militar mancuernas') LIMIT 1)
  );
