UPDATE users
SET password_hash = crypt('secret123', gen_salt('bf'))
WHERE email IN (
  'roshi@kamehouse.example',
  'baki.hanma@example.example',
  'ikki@saint.example',
  'krillin@kamehouse.example',
  'yamcha@capsule.example',
  'goku@capsule.example',
  'retsu@shinshinkai.example',
  'katsumi@shinshinkai.example',
  'jack.hanma@example.example'
);
