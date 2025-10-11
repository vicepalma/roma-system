-- === Usuarios base ===
INSERT INTO users (id, name, email, password_hash, created_at) VALUES
('11111111-1111-1111-1111-111111111111', 'Maestro Roshi', 'coach@example.com', '$2a$10$PbvWZsbmMSbPSH5hQLm5peE4AB7p0RhxT6QEBjtkPC3T0UPq2Up5i', NOW()), -- password: 1234
('22222222-2222-2222-2222-222222222222', 'Krilin', 'disciple@example.com', '$2a$10$PbvWZsbmMSbPSH5hQLm5peE4AB7p0RhxT6QEBjtkPC3T0UPq2Up5i', NOW()); -- password: 1234

-- === Vínculo coach-discípulo (ya aceptado) ===
INSERT INTO coach_links (coach_id, disciple_id, status, created_at)
VALUES ('11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', 'accepted', NOW());

-- === Programa ejemplo ===
INSERT INTO programs (id, name, description, created_by, version, created_at)
VALUES ('27173377-7223-413f-8295-62bd8d2aa978', 'Programa FST-7 Pecho/Hombro', 'Programa ejemplo de 4 semanas', 
        '11111111-1111-1111-1111-111111111111', 1, NOW());

-- === Semana y Día ===
INSERT INTO weeks (id, program_id, week_index) 
VALUES ('cc87fb42-2a6d-4285-a54a-99bde61b9b49', '27173377-7223-413f-8295-62bd8d2aa978', 1);

INSERT INTO days (id, week_id, day_index, notes)
VALUES ('e39429d3-7293-4138-b4c4-fc8f313d9827', 'cc87fb42-2a6d-4285-a54a-99bde61b9b49', 1, 'Pecho/Hombro');

-- === Prescripción ejemplo ===
INSERT INTO prescriptions (id, day_id, exercise_id, series, reps, rest_sec, position)
VALUES ('5a1be6da-8927-4169-b6e5-afaf73fe8635', 'e39429d3-7293-4138-b4c4-fc8f313d9827', 
        '76395cbf-74ce-4ff5-9ba2-9a6937476d7e', 4, '10-12', 90, 1);

-- === Asignación del programa al discípulo ===
INSERT INTO assignments (id, program_id, disciple_id, assigned_by, start_date, is_active, created_at)
VALUES ('03ac230c-9bdf-4e24-9294-5832f589ee64', '27173377-7223-413f-8295-62bd8d2aa978', 
        '22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', 
        CURRENT_DATE, TRUE, NOW());
