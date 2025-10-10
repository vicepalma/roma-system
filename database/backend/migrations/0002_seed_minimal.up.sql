-- Método FST-7
INSERT INTO methods (key, name, params)
VALUES ('fst7','FST-7', '{"series":7,"rest_sec":30,"to_failure":true,"target_reps":"10-12"}');

-- Algunos ejercicios base
INSERT INTO exercises (name, primary_muscle, equipment, tags) VALUES
 ('Pec fly', 'chest', 'machine',  ARRAY['hypertrophy']),
 ('Press plano', 'chest', 'barbell', ARRAY['compound']),
 ('Posterior en poleas', 'shoulders', 'cable', ARRAY['rear-delt']),
 ('Press militar mancuernas', 'shoulders', 'dumbbell', ARRAY['compound']),
 ('Vuelos laterales', 'shoulders', 'dumbbell', ARRAY['isolation']),
 ('Extensión tríceps en polea', 'triceps', 'cable', ARRAY['isolation']);
