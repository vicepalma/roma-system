BEGIN;

DELETE FROM assignments
WHERE program_id IN (
  SELECT p.id
  FROM programs p
  JOIN users u ON u.id = p.owner_id
  WHERE u.email = 'fst7.coach@roma.demo'
    AND p.title = '[DEMO] FST-7 Hipertrofia - 5 dias'
);

DELETE FROM program_versions
WHERE program_id IN (
  SELECT p.id
  FROM programs p
  JOIN users u ON u.id = p.owner_id
  WHERE u.email = 'fst7.coach@roma.demo'
    AND p.title = '[DEMO] FST-7 Hipertrofia - 5 dias'
);

DELETE FROM programs
WHERE owner_id IN (SELECT id FROM users WHERE email = 'fst7.coach@roma.demo')
  AND title = '[DEMO] FST-7 Hipertrofia - 5 dias';

DELETE FROM coach_links
WHERE coach_id IN (SELECT id FROM users WHERE email = 'fst7.coach@roma.demo')
   OR disciple_id IN (SELECT id FROM users WHERE email = 'fst7.disciple@roma.demo');

DELETE FROM master_disciple
WHERE master_id IN (SELECT id FROM users WHERE email = 'fst7.coach@roma.demo')
   OR disciple_id IN (SELECT id FROM users WHERE email = 'fst7.disciple@roma.demo');

DELETE FROM users
WHERE email IN ('fst7.coach@roma.demo', 'fst7.disciple@roma.demo');

DELETE FROM exercises
WHERE lower(name) IN (
  'press inclinado con mancuernas',
  'press pecho en maquina convergente',
  'aperturas en polea alta',
  'jalon al pecho agarre neutro',
  'remo pecho apoyado',
  'pullover en polea',
  'sentadilla hack',
  'prensa 45',
  'extension de cuadriceps',
  'press hombro en maquina',
  'elevacion lateral en polea',
  'face pull en polea',
  'curl predicador maquina',
  'curl martillo con cuerda',
  'press cerrado en maquina',
  'extension triceps cuerda'
);

COMMIT;
