BEGIN;

-- Demo local: FST-7 seed data kept isolated with @roma.demo users and [DEMO] titles.
WITH cfg AS (
  SELECT '$argon2id$v=19$m=65536,t=1,p=8$CJcnBbz7AoJRdZslhgpJxg$gTVT+6YGG2M4r87dLvbfIAY1inQx1MqP0v1/qLcKHjk'::text AS ph
)
INSERT INTO users (email, password_hash, name, role)
SELECT email, (SELECT ph FROM cfg), name, role
FROM (
  VALUES
    ('fst7.coach@roma.demo'::citext, 'Coach Demo FST-7'::text, 'coach'::text),
    ('fst7.disciple@roma.demo'::citext, 'Disciple Demo FST-7'::text, 'disciple'::text)
) AS u(email, name, role)
ON CONFLICT (email) DO UPDATE
SET name = EXCLUDED.name,
    role = EXCLUDED.role;

WITH u AS (
  SELECT
    (SELECT id FROM users WHERE email = 'fst7.coach@roma.demo') AS coach_id,
    (SELECT id FROM users WHERE email = 'fst7.disciple@roma.demo') AS disciple_id
)
INSERT INTO master_disciple (master_id, disciple_id, status)
SELECT coach_id, disciple_id, 'active'
FROM u
WHERE coach_id IS NOT NULL
  AND disciple_id IS NOT NULL
ON CONFLICT (master_id, disciple_id) DO UPDATE
SET status = 'active';

WITH u AS (
  SELECT
    (SELECT id FROM users WHERE email = 'fst7.coach@roma.demo') AS coach_id,
    (SELECT id FROM users WHERE email = 'fst7.disciple@roma.demo') AS disciple_id
)
INSERT INTO coach_links (coach_id, disciple_id, status)
SELECT coach_id, disciple_id, 'accepted'
FROM u
WHERE coach_id IS NOT NULL
  AND disciple_id IS NOT NULL
ON CONFLICT (coach_id, disciple_id) DO UPDATE
SET status = 'accepted';

INSERT INTO methods (key, name, params)
VALUES (
  'fst7',
  'FST-7',
  '{"series":7,"rest_sec":30,"to_failure":true,"target_reps":"10-12"}'::jsonb
)
ON CONFLICT (key) DO UPDATE
SET name = EXCLUDED.name,
    params = EXCLUDED.params;

INSERT INTO exercises (name, primary_muscle, equipment, tags, notes)
VALUES
  ('Press inclinado con mancuernas', 'chest', 'dumbbell', ARRAY['demo','fst7','compound'], 'Pecho superior, recorrido controlado.'),
  ('Press pecho en maquina convergente', 'chest', 'machine', ARRAY['demo','fst7','compound'], 'Serie pesada estable antes del bombeo final.'),
  ('Aperturas en polea alta', 'chest', 'cable', ARRAY['demo','fst7','isolation'], 'Cruzar ligeramente al centro y sostener contraccion.'),
  ('Jalon al pecho agarre neutro', 'back', 'machine', ARRAY['demo','fst7','compound'], 'Enfasis en depresion escapular.'),
  ('Remo pecho apoyado', 'back', 'machine', ARRAY['demo','fst7','compound'], 'Evitar impulso lumbar.'),
  ('Pullover en polea', 'back', 'cable', ARRAY['demo','fst7','isolation'], 'Mantener brazos casi extendidos.'),
  ('Sentadilla hack', 'legs', 'machine', ARRAY['demo','fst7','compound'], 'Profundidad consistente sin bloquear rodillas.'),
  ('Prensa 45', 'legs', 'machine', ARRAY['demo','fst7','compound'], 'Pies medios, recorrido completo.'),
  ('Extension de cuadriceps', 'legs', 'machine', ARRAY['demo','fst7','isolation'], 'Pausa arriba en cada repeticion.'),
  ('Press hombro en maquina', 'shoulders', 'machine', ARRAY['demo','fst7','compound'], 'Controlar bajada y no arquear espalda.'),
  ('Elevacion lateral en polea', 'shoulders', 'cable', ARRAY['demo','fst7','isolation'], 'Tension continua para deltoide medio.'),
  ('Face pull en polea', 'shoulders', 'cable', ARRAY['demo','fst7','rear-delt'], 'Codos altos, rotacion externa suave.'),
  ('Curl predicador maquina', 'arms', 'machine', ARRAY['demo','fst7','isolation'], 'No despegar el brazo del apoyo.'),
  ('Curl martillo con cuerda', 'arms', 'cable', ARRAY['demo','fst7','isolation'], 'Separar puntas al final del recorrido.'),
  ('Press cerrado en maquina', 'arms', 'machine', ARRAY['demo','fst7','compound'], 'Triceps con hombro estable.'),
  ('Extension triceps cuerda', 'arms', 'cable', ARRAY['demo','fst7','isolation'], 'Abrir cuerda y bloquear suave.')
ON CONFLICT (lower(name)) DO UPDATE
SET primary_muscle = EXCLUDED.primary_muscle,
    equipment = EXCLUDED.equipment,
    tags = EXCLUDED.tags,
    notes = EXCLUDED.notes;

WITH owner_row AS (
  SELECT id AS owner_id FROM users WHERE email = 'fst7.coach@roma.demo'
), existing AS (
  SELECT p.id
  FROM programs p
  JOIN owner_row o ON o.owner_id = p.owner_id
  WHERE p.title = '[DEMO] FST-7 Hipertrofia - 5 dias'
), inserted AS (
  INSERT INTO programs (owner_id, title, notes, visibility, version, kind)
  SELECT
    owner_id,
    '[DEMO] FST-7 Hipertrofia - 5 dias',
    'Demo local FST-7: bloques de hipertrofia con un cierre de 7 series, 30-45s de descanso y foco en bombeo. Datos seed, no productivos.',
    'private',
    1,
    'coach_program'
  FROM owner_row
  WHERE NOT EXISTS (SELECT 1 FROM existing)
  RETURNING id
), target AS (
  SELECT id FROM inserted
  UNION ALL
  SELECT id FROM existing
), version_row AS (
  INSERT INTO program_versions (program_id, version, title, notes)
  SELECT
    id,
    1,
    '[DEMO] FST-7 Hipertrofia - 5 dias',
    'Demo local FST-7 seed.'
  FROM target
  ON CONFLICT (program_id, version) DO NOTHING
  RETURNING program_id
)
INSERT INTO program_weeks (program_id, week_index)
SELECT id, 1
FROM target
WHERE NOT EXISTS (
  SELECT 1 FROM program_weeks w WHERE w.program_id = target.id AND w.week_index = 1
);

WITH program_row AS (
  SELECT p.id AS program_id
  FROM programs p
  JOIN users u ON u.id = p.owner_id
  WHERE u.email = 'fst7.coach@roma.demo'
    AND p.title = '[DEMO] FST-7 Hipertrofia - 5 dias'
), week_row AS (
  SELECT w.id AS week_id
  FROM program_weeks w
  JOIN program_row p ON p.program_id = w.program_id
  WHERE w.week_index = 1
)
INSERT INTO program_days (week_id, day_index, title, notes)
SELECT week_id, day_index, title, notes
FROM week_row
CROSS JOIN (
  VALUES
    (1, 'Pecho FST-7', 'Dia de pecho con cierre FST-7 en aperturas de polea.'),
    (2, 'Espalda FST-7', 'Tirones pesados y bombeo final en pullover.'),
    (3, 'Piernas FST-7', 'Cuadriceps dominante con extension final de 7 series.'),
    (4, 'Hombros FST-7', 'Deltoide medio como cierre metabolico.'),
    (5, 'Brazos FST-7', 'Biceps y triceps con cierre de alto bombeo.')
) AS d(day_index, title, notes)
WHERE NOT EXISTS (
  SELECT 1 FROM program_days pd
  WHERE pd.week_id = week_row.week_id
    AND pd.day_index = d.day_index
);

WITH method_row AS (
  SELECT id AS method_id FROM methods WHERE key = 'fst7'
), days AS (
  SELECT d.id AS day_id, d.day_index
  FROM program_days d
  JOIN program_weeks w ON w.id = d.week_id
  JOIN programs p ON p.id = w.program_id
  JOIN users u ON u.id = p.owner_id
  WHERE u.email = 'fst7.coach@roma.demo'
    AND p.title = '[DEMO] FST-7 Hipertrofia - 5 dias'
), plan AS (
  SELECT *
  FROM (
    VALUES
      (1, 1, 'Press inclinado con mancuernas', 4, '8-10', 120, false, NULL::text, 2, 'Compuesto principal, dejar 1-2 reps en reserva.'),
      (1, 2, 'Press pecho en maquina convergente', 3, '10-12', 90, false, NULL::text, 1, 'Controlar la excentrica.'),
      (1, 3, 'Aperturas en polea alta', 7, '10-12', 30, true, 'fst7', 0, 'Cierre FST-7: descanso corto, bombeo y tecnica limpia.'),
      (2, 1, 'Jalon al pecho agarre neutro', 4, '8-10', 120, false, NULL::text, 2, 'Pausar abajo sin balanceo.'),
      (2, 2, 'Remo pecho apoyado', 4, '10', 90, false, NULL::text, 1, 'Recorrido completo y pecho fijo.'),
      (2, 3, 'Pullover en polea', 7, '10-12', 30, true, 'fst7', 0, 'Cierre FST-7 para dorsales.'),
      (3, 1, 'Sentadilla hack', 4, '8-10', 150, false, NULL::text, 2, 'Calentar antes de la primera serie efectiva.'),
      (3, 2, 'Prensa 45', 3, '12', 120, false, NULL::text, 1, 'No bloquear rodillas.'),
      (3, 3, 'Extension de cuadriceps', 7, '10-12', 30, true, 'fst7', 0, 'Cierre FST-7 con pausa arriba.'),
      (4, 1, 'Press hombro en maquina', 4, '8-10', 120, false, NULL::text, 2, 'Mantener espalda apoyada.'),
      (4, 2, 'Face pull en polea', 3, '12-15', 60, false, NULL::text, 2, 'Activacion posterior y salud de hombro.'),
      (4, 3, 'Elevacion lateral en polea', 7, '10-12', 30, true, 'fst7', 0, 'Cierre FST-7 para deltoide medio.'),
      (5, 1, 'Curl predicador maquina', 4, '10-12', 75, false, NULL::text, 1, 'Evitar rebote al extender.'),
      (5, 2, 'Press cerrado en maquina', 4, '8-10', 90, false, NULL::text, 1, 'Triceps pesado y estable.'),
      (5, 3, 'Curl martillo con cuerda', 3, '12', 60, false, NULL::text, 1, 'Braquial y antebrazo.'),
      (5, 4, 'Extension triceps cuerda', 7, '10-12', 30, true, 'fst7', 0, 'Cierre FST-7 de triceps.')
  ) AS p(day_index, position, exercise_name, series, reps, rest_sec, to_failure, method_key, rir, notes)
)
INSERT INTO prescriptions (
  day_id,
  exercise_id,
  series,
  reps,
  rest_sec,
  to_failure,
  rir,
  method_id,
  notes,
  position
)
SELECT
  d.day_id,
  e.id,
  p.series,
  p.reps,
  p.rest_sec,
  p.to_failure,
  p.rir,
  CASE WHEN p.method_key = 'fst7' THEN (SELECT method_id FROM method_row) ELSE NULL END,
  p.notes,
  p.position
FROM plan p
JOIN days d ON d.day_index = p.day_index
JOIN exercises e ON e.name = p.exercise_name
WHERE NOT EXISTS (
  SELECT 1
  FROM prescriptions pr
  WHERE pr.day_id = d.day_id
    AND pr.exercise_id = e.id
    AND pr.position = p.position
);

WITH ids AS (
  SELECT
    (SELECT p.id
     FROM programs p
     JOIN users u ON u.id = p.owner_id
     WHERE u.email = 'fst7.coach@roma.demo'
       AND p.title = '[DEMO] FST-7 Hipertrofia - 5 dias'
     LIMIT 1) AS program_id,
    (SELECT id FROM users WHERE email = 'fst7.coach@roma.demo') AS coach_id,
    (SELECT id FROM users WHERE email = 'fst7.disciple@roma.demo') AS disciple_id
)
INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
SELECT program_id, 1, disciple_id, coach_id, CURRENT_DATE, TRUE
FROM ids
WHERE program_id IS NOT NULL
  AND coach_id IS NOT NULL
  AND disciple_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1
    FROM assignments a
    WHERE a.program_id = ids.program_id
      AND a.disciple_id = ids.disciple_id
      AND a.assigned_by = ids.coach_id
  );

COMMIT;
