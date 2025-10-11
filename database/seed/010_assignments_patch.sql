-- Crea assignments si faltan, resolviendo IDs por email y tomando
-- el programa MÁS RECIENTE de cada owner (Roshi, Baki, Ikki).

-- ROSHI -> Krillin
INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
SELECT
  (SELECT p.id FROM programs p
   JOIN users o ON o.id = p.owner_id
   WHERE o.email = 'roshi@kamehouse.example'
   ORDER BY p.created_at DESC LIMIT 1)                          AS program_id,
  1                                                             AS program_version,
  (SELECT id FROM users WHERE email='krillin@kamehouse.example') AS disciple_id,
  (SELECT id FROM users WHERE email='roshi@kamehouse.example')   AS assigned_by,
  CURRENT_DATE, TRUE
WHERE
  (SELECT COUNT(*) FROM assignments a
     WHERE a.program_id = (SELECT p.id FROM programs p JOIN users o ON o.id=p.owner_id
                           WHERE o.email='roshi@kamehouse.example'
                           ORDER BY p.created_at DESC LIMIT 1)
       AND a.disciple_id = (SELECT id FROM users WHERE email='krillin@kamehouse.example')
  ) = 0
  AND EXISTS (SELECT 1 FROM users WHERE email='roshi@kamehouse.example')
  AND EXISTS (SELECT 1 FROM users WHERE email='krillin@kamehouse.example')
  AND EXISTS (SELECT 1 FROM programs p JOIN users o ON o.id=p.owner_id WHERE o.email='roshi@kamehouse.example');

-- ROSHI -> Yamcha
INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
SELECT
  (SELECT p.id FROM programs p
   JOIN users o ON o.id = p.owner_id
   WHERE o.email = 'roshi@kamehouse.example'
   ORDER BY p.created_at DESC LIMIT 1),
  1,
  (SELECT id FROM users WHERE email='yamcha@capsule.example'),
  (SELECT id FROM users WHERE email='roshi@kamehouse.example'),
  CURRENT_DATE, TRUE
WHERE
  (SELECT COUNT(*) FROM assignments a
     WHERE a.program_id = (SELECT p.id FROM programs p JOIN users o ON o.id=p.owner_id
                           WHERE o.email='roshi@kamehouse.example'
                           ORDER BY p.created_at DESC LIMIT 1)
       AND a.disciple_id = (SELECT id FROM users WHERE email='yamcha@capsule.example')
  ) = 0
  AND EXISTS (SELECT 1 FROM users WHERE email='yamcha@capsule.example')
  AND EXISTS (SELECT 1 FROM programs p JOIN users o ON o.id=p.owner_id WHERE o.email='roshi@kamehouse.example');

-- ROSHI -> Goku
INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
SELECT
  (SELECT p.id FROM programs p
   JOIN users o ON o.id = p.owner_id
   WHERE o.email = 'roshi@kamehouse.example'
   ORDER BY p.created_at DESC LIMIT 1),
  1,
  (SELECT id FROM users WHERE email='goku@capsule.example'),
  (SELECT id FROM users WHERE email='roshi@kamehouse.example'),
  CURRENT_DATE, TRUE
WHERE
  (SELECT COUNT(*) FROM assignments a
     WHERE a.program_id = (SELECT p.id FROM programs p JOIN users o ON o.id=p.owner_id
                           WHERE o.email='roshi@kamehouse.example'
                           ORDER BY p.created_at DESC LIMIT 1)
       AND a.disciple_id = (SELECT id FROM users WHERE email='goku@capsule.example')
  ) = 0
  AND EXISTS (SELECT 1 FROM users WHERE email='goku@capsule.example')
  AND EXISTS (SELECT 1 FROM programs p JOIN users o ON o.id=p.owner_id WHERE o.email='roshi@kamehouse.example');

-- BAKI -> Retsu
INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
SELECT
  (SELECT p.id FROM programs p
   JOIN users o ON o.id = p.owner_id
   WHERE o.email = 'baki.hanma@example.example'
   ORDER BY p.created_at DESC LIMIT 1),
  1,
  (SELECT id FROM users WHERE email='retsu@shinshinkai.example'),
  (SELECT id FROM users WHERE email='baki.hanma@example.example'),
  CURRENT_DATE, TRUE
WHERE
  (SELECT COUNT(*) FROM assignments a
     WHERE a.program_id = (SELECT p.id FROM programs p JOIN users o ON o.id=p.owner_id
                           WHERE o.email='baki.hanma@example.example'
                           ORDER BY p.created_at DESC LIMIT 1)
       AND a.disciple_id = (SELECT id FROM users WHERE email='retsu@shinshinkai.example')
  ) = 0
  AND EXISTS (SELECT 1 FROM users WHERE email='retsu@shinshinkai.example')
  AND EXISTS (SELECT 1 FROM programs p JOIN users o ON o.id=p.owner_id WHERE o.email='baki.hanma@example.example');

-- BAKI -> Katsumi
INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
SELECT
  (SELECT p.id FROM programs p
   JOIN users o ON o.id = p.owner_id
   WHERE o.email = 'baki.hanma@example.example'
   ORDER BY p.created_at DESC LIMIT 1),
  1,
  (SELECT id FROM users WHERE email='katsumi@shinshinkai.example'),
  (SELECT id FROM users WHERE email='baki.hanma@example.example'),
  CURRENT_DATE, TRUE
WHERE
  (SELECT COUNT(*) FROM assignments a
     WHERE a.program_id = (SELECT p.id FROM programs p JOIN users o ON o.id=p.owner_id
                           WHERE o.email='baki.hanma@example.example'
                           ORDER BY p.created_at DESC LIMIT 1)
       AND a.disciple_id = (SELECT id FROM users WHERE email='katsumi@shinshinkai.example')
  ) = 0
  AND EXISTS (SELECT 1 FROM users WHERE email='katsumi@shinshinkai.example')
  AND EXISTS (SELECT 1 FROM programs p JOIN users o ON o.id=p.owner_id WHERE o.email='baki.hanma@example.example');

-- BAKI -> Jack
INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
SELECT
  (SELECT p.id FROM programs p
   JOIN users o ON o.id = p.owner_id
   WHERE o.email = 'baki.hanma@example.example'
   ORDER BY p.created_at DESC LIMIT 1),
  1,
  (SELECT id FROM users WHERE email='jack.hanma@example.example'),
  (SELECT id FROM users WHERE email='baki.hanma@example.example'),
  CURRENT_DATE, TRUE
WHERE
  (SELECT COUNT(*) FROM assignments a
     WHERE a.program_id = (SELECT p.id FROM programs p JOIN users o ON o.id=p.owner_id
                           WHERE o.email='baki.hanma@example.example'
                           ORDER BY p.created_at DESC LIMIT 1)
       AND a.disciple_id = (SELECT id FROM users WHERE email='jack.hanma@example.example')
  ) = 0
  AND EXISTS (SELECT 1 FROM users WHERE email='jack.hanma@example.example')
  AND EXISTS (SELECT 1 FROM programs p JOIN users o ON o.id=p.owner_id WHERE o.email='baki.hanma@example.example');

-- IKKI -> Ikki (auto-asignación)
INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
SELECT
  (SELECT p.id FROM programs p
   JOIN users o ON o.id = p.owner_id
   WHERE o.email = 'ikki@saint.example'
   ORDER BY p.created_at DESC LIMIT 1),
  1,
  (SELECT id FROM users WHERE email='ikki@saint.example'),
  (SELECT id FROM users WHERE email='ikki@saint.example'),
  CURRENT_DATE, TRUE
WHERE
  (SELECT COUNT(*) FROM assignments a
     WHERE a.program_id = (SELECT p.id FROM programs p JOIN users o ON o.id=p.owner_id
                           WHERE o.email='ikki@saint.example'
                           ORDER BY p.created_at DESC LIMIT 1)
       AND a.disciple_id = (SELECT id FROM users WHERE email='ikki@saint.example')
  ) = 0
  AND EXISTS (SELECT 1 FROM users WHERE email='ikki@saint.example')
  AND EXISTS (SELECT 1 FROM programs p JOIN users o ON o.id=p.owner_id WHERE o.email='ikki@saint.example');
