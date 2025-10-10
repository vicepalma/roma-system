DELETE FROM exercises WHERE lower(name) IN (
  'pec fly','press plano','posterior en poleas',
  'press militar mancuernas','vuelos laterales','extensión tríceps en polea'
);
DELETE FROM methods WHERE key = 'fst7';
