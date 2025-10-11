-- ==== Prerrequisitos (por si faltan) ====
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

-- ==== Usuarios (bcrypt de "1234") ====
-- hash: $2a$10$PbvWZsbmMSbPSH5hQLm5peE4AB7p0RhxT6QEBjtkPC3T0UPq2Up5i
INSERT INTO users (id, email, password_hash, name) VALUES
('11111111-1111-1111-1111-111111111111','ada.coach@example.com','$2a$10$PbvWZsbmMSbPSH5hQLm5peE4AB7p0RhxT6QEBjtkPC3T0UPq2Up5i','Ada Coach'),
('22222222-2222-2222-2222-222222222222','bruno.coach@example.com','$2a$10$PbvWZsbmMSbPSH5hQLm5peE4AB7p0RhxT6QEBjtkPC3T0UPq2Up5i','Bruno Coach'),
('33333333-3333-3333-3333-333333333333','sam.solo@example.com','$2a$10$PbvWZsbmMSbPSH5hQLm5peE4AB7p0RhxT6QEBjtkPC3T0UPq2Up5i','Sam Solo'),

('aaa11111-1111-1111-1111-111111111111','disc1.ada@example.com','$2a$10$PbvWZsbmMSbPSH5hQLm5peE4AB7p0RhxT6QEBjtkPC3T0UPq2Up5i','Disc Ada 1'),
('aaa22222-2222-2222-2222-222222222222','disc2.ada@example.com','$2a$10$PbvWZsbmMSbPSH5hQLm5peE4AB7p0RhxT6QEBjtkPC3T0UPq2Up5i','Disc Ada 2'),
('aaa33333-3333-3333-3333-333333333333','disc3.ada@example.com','$2a$10$PbvWZsbmMSbPSH5hQLm5peE4AB7p0RhxT6QEBjtkPC3T0UPq2Up5i','Disc Ada 3'),

('bbb11111-1111-1111-1111-111111111111','disc1.bruno@example.com','$2a$10$PbvWZsbmMSbPSH5hQLm5peE4AB7p0RhxT6QEBjtkPC3T0UPq2Up5i','Disc Bruno 1'),
('bbb22222-2222-2222-2222-222222222222','disc2.bruno@example.com','$2a$10$PbvWZsbmMSbPSH5hQLm5peE4AB7p0RhxT6QEBjtkPC3T0UPq2Up5i','Disc Bruno 2'),
('bbb33333-3333-3333-3333-333333333333','disc3.bruno@example.com','$2a$10$PbvWZsbmMSbPSH5hQLm5peE4AB7p0RhxT6QEBjtkPC3T0UPq2Up5i','Disc Bruno 3')
ON CONFLICT (email) DO NOTHING;

-- ==== Vínculos master-disciple ====
-- Ada con 3 discípulos
INSERT INTO master_disciple (master_id, disciple_id, status) VALUES
('11111111-1111-1111-1111-111111111111','aaa11111-1111-1111-1111-111111111111','active'),
('11111111-1111-1111-1111-111111111111','aaa22222-2222-2222-2222-222222222222','active'),
('11111111-1111-1111-1111-111111111111','aaa33333-3333-3333-3333-333333333333','active')
ON CONFLICT (master_id,disciple_id) DO NOTHING;

-- Bruno con 3 discípulos
INSERT INTO master_disciple (master_id, disciple_id, status) VALUES
('22222222-2222-2222-2222-222222222222','bbb11111-1111-1111-1111-111111111111','active'),
('22222222-2222-2222-2222-222222222222','bbb22222-2222-2222-2222-222222222222','active'),
('22222222-2222-2222-2222-222222222222','bbb33333-3333-3333-3333-333333333333','active')
ON CONFLICT (master_id,disciple_id) DO NOTHING;

-- Sam es su propio discípulo (auto-vínculo)
INSERT INTO master_disciple (master_id, disciple_id, status) VALUES
('33333333-3333-3333-3333-333333333333','33333333-3333-3333-3333-333333333333','active')
ON CONFLICT (master_id,disciple_id) DO NOTHING;

-- ==== Métodos (por si no existe) ====
INSERT INTO methods (key, name, params)
VALUES ('fst7','FST-7','{"series":7,"rest_sec":30,"to_failure":true,"target_reps":"10-12"}')
ON CONFLICT (key) DO NOTHING;

-- ==== Ejercicios base (por si faltan) ====
INSERT INTO exercises (name, primary_muscle, equipment, tags) VALUES
 ('Pec fly', 'chest', 'machine',  ARRAY['hypertrophy']),
 ('Press plano', 'chest', 'barbell', ARRAY['compound']),
 ('Posterior en poleas', 'shoulders', 'cable', ARRAY['rear-delt']),
 ('Press militar mancuernas', 'shoulders', 'dumbbell', ARRAY['compound']),
 ('Vuelos laterales', 'shoulders', 'dumbbell', ARRAY['isolation']),
 ('Extensión tríceps en polea', 'triceps', 'cable', ARRAY['isolation'])
ON CONFLICT (lower(name)) DO NOTHING;

-- ==== Programas (uno por maestro) ====
INSERT INTO programs (id, owner_id, title, notes, visibility, version) VALUES
('27173377-7223-413f-8295-62bd8d2aa978','11111111-1111-1111-1111-111111111111','Pecho/Hombro 4w','Programa ejemplo Ada','private',1),
('38173377-7223-413f-8295-62bd8d2aa978','22222222-2222-2222-2222-222222222222','Full Body 2w','Programa ejemplo Bruno','private',1),
('49173377-7223-413f-8295-62bd8d2aa978','33333333-3333-3333-3333-333333333333','Auto-Rutina','Programa auto Sam','private',1)
ON CONFLICT (id) DO NOTHING;

-- ==== Weeks / Days (1 semana y 1 día por programa para el seed) ====
INSERT INTO program_weeks (id, program_id, week_index) VALUES
('cc87fb42-2a6d-4285-a54a-99bde61b9b49', '27173377-7223-413f-8295-62bd8d2aa978', 1),
('cc87fb42-2a6d-4285-a54a-99bde61b9b50', '38173377-7223-413f-8295-62bd8d2aa978', 1),
('cc87fb42-2a6d-4285-a54a-99bde61b9b51', '49173377-7223-413f-8295-62bd8d2aa978', 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO program_days (id, week_id, day_index, notes) VALUES
('e39429d3-7293-4138-b4c4-fc8f313d9827','cc87fb42-2a6d-4285-a54a-99bde61b9b49',1,'Pecho/Hombro'),
('e39429d3-7293-4138-b4c4-fc8f313d9828','cc87fb42-2a6d-4285-a54a-99bde61b9b50',1,'Full body A'),
('e39429d3-7293-4138-b4c4-fc8f313d9829','cc87fb42-2a6d-4285-a54a-99bde61b9b51',1,'Auto día 1')
ON CONFLICT (id) DO NOTHING;

-- ==== Prescripciones (referenciando ejercicios por nombre) ====
-- Día de Ada
INSERT INTO prescriptions (id, day_id, exercise_id, series, reps, rest_sec, to_failure, position)
VALUES
('5a1be6da-8927-4169-b6e5-afaf73fe8635',
 'e39429d3-7293-4138-b4c4-fc8f313d9827',
 (SELECT id FROM exercises WHERE name='Press plano' LIMIT 1), 4, '10-12', 90, FALSE, 1),
(gen_random_uuid(),
 'e39429d3-7293-4138-b4c4-fc8f313d9827',
 (SELECT id FROM exercises WHERE name='Vuelos laterales' LIMIT 1), 3, '12-15', 60, FALSE, 2)
ON CONFLICT (id) DO NOTHING;

-- Día de Bruno
INSERT INTO prescriptions (id, day_id, exercise_id, series, reps, rest_sec, to_failure, position)
VALUES
(gen_random_uuid(),
 'e39429d3-7293-4138-b4c4-fc8f313d9828',
 (SELECT id FROM exercises WHERE name='Pec fly' LIMIT 1), 4, '12', 60, FALSE, 1),
(gen_random_uuid(),
 'e39429d3-7293-4138-b4c4-fc8f313d9828',
 (SELECT id FROM exercises WHERE name='Extensión tríceps en polea' LIMIT 1), 4, '10-12', 90, FALSE, 2)
ON CONFLICT (id) DO NOTHING;

-- Día de Sam (auto)
INSERT INTO prescriptions (id, day_id, exercise_id, series, reps, rest_sec, to_failure, position)
VALUES
(gen_random_uuid(),
 'e39429d3-7293-4138-b4c4-fc8f313d9829',
 (SELECT id FROM exercises WHERE name='Press militar mancuernas' LIMIT 1), 4, '8-10', 120, FALSE, 1)
ON CONFLICT (id) DO NOTHING;

-- ==== Assignments (programa → discípulos), versión fijada 1 ====
-- Ada asigna su programa a sus 3 discípulos
INSERT INTO assignments (id, program_id, program_version, disciple_id, assigned_by, start_date, is_active)
VALUES
('03ac230c-9bdf-4e24-9294-5832f589ee64','27173377-7223-413f-8295-62bd8d2aa978',1,'aaa11111-1111-1111-1111-111111111111','11111111-1111-1111-1111-111111111111', CURRENT_DATE, TRUE),
(gen_random_uuid(),'27173377-7223-413f-8295-62bd8d2aa978',1,'aaa22222-2222-2222-2222-222222222222','11111111-1111-1111-1111-111111111111', CURRENT_DATE, TRUE),
(gen_random_uuid(),'27173377-7223-413f-8295-62bd8d2aa978',1,'aaa33333-3333-3333-3333-333333333333','11111111-1111-1111-1111-111111111111', CURRENT_DATE, TRUE)
ON CONFLICT (id) DO NOTHING;

-- Bruno asigna a sus 3 discípulos
INSERT INTO assignments (id, program_id, program_version, disciple_id, assigned_by, start_date, is_active)
VALUES
(gen_random_uuid(),'38173377-7223-413f-8295-62bd8d2aa978',1,'bbb11111-1111-1111-1111-111111111111','22222222-2222-2222-2222-222222222222', CURRENT_DATE, TRUE),
(gen_random_uuid(),'38173377-7223-413f-8295-62bd8d2aa978',1,'bbb22222-2222-2222-2222-222222222222','22222222-2222-2222-2222-222222222222', CURRENT_DATE, TRUE),
(gen_random_uuid(),'38173377-7223-413f-8295-62bd8d2aa978',1,'bbb33333-3333-3333-3333-333333333333','22222222-2222-2222-2222-222222222222', CURRENT_DATE, TRUE)
ON CONFLICT (id) DO NOTHING;

-- Sam se asigna su propio programa
INSERT INTO assignments (id, program_id, program_version, disciple_id, assigned_by, start_date, is_active)
VALUES
(gen_random_uuid(),'49173377-7223-413f-8295-62bd8d2aa978',1,'33333333-3333-3333-3333-333333333333','33333333-3333-3333-3333-333333333333', CURRENT_DATE, TRUE)
ON CONFLICT (id) DO NOTHING;

-- ==== Una sesión de ejemplo + sets (para que history/pivot tenga datos) ====
-- Usamos el primer assignment de Ada (id fijo arriba), su día y la prescripción fija
INSERT INTO session_logs (id, assignment_id, disciple_id, day_id, performed_at, notes)
VALUES ('7f2f5bec-fa4d-4f8a-9331-f4ed509e61d0',
        '03ac230c-9bdf-4e24-9294-5832f589ee64',
        'aaa11111-1111-1111-1111-111111111111',
        'e39429d3-7293-4138-b4c4-fc8f313d9827',
        NOW(), 'Seed session Ada Disc 1')
ON CONFLICT (id) DO NOTHING;

-- 3 sets sobre la prescripción 'Press plano'
INSERT INTO set_logs (session_id, prescription_id, set_index, weight, reps, rpe, to_failure) VALUES
('7f2f5bec-fa4d-4f8a-9331-f4ed509e61d0','5a1be6da-8927-4169-b6e5-afaf73fe8635',1, 60,12,7.5,false),
('7f2f5bec-fa4d-4f8a-9331-f4ed509e61d0','5a1be6da-8927-4169-b6e5-afaf73fe8635',2, 65,10,8.0,false),
('7f2f5bec-fa4d-4f8a-9331-f4ed509e61d0','5a1be6da-8927-4169-b6e5-afaf73fe8635',3, 70, 8,8.5,true);
