# Demo Users

Password comun seed: `secret123`.

## Coach

- `roshi@kamehouse.example` / `secret123` - Maestro Roshi.
- `baki.hanma@example.example` / `secret123` - Baki Hanma.
- `ikki@saint.example` / `secret123` - Ikki (Fenix).

## Disciple

- `krillin@kamehouse.example` / `secret123` - Krillin.
- `yamcha@capsule.example` / `secret123` - Yamcha.
- `goku@capsule.example` / `secret123` - Goku.
- `retsu@shinshinkai.example` / `secret123` - Retsu Kaioh.
- `katsumi@shinshinkai.example` / `secret123` - Katsumi Orochi.
- `jack.hanma@example.example` / `secret123` - Jack Hanma.

## Demo FST-7

- `fst7.coach@roma.demo` / `secret123` - Coach Demo FST-7.
- `fst7.disciple@roma.demo` / `secret123` - Disciple Demo FST-7 con assignment activo `[DEMO] FST-7 Hipertrofia - 5 dias`.
- Estas cuentas no se crean con las migraciones oficiales. Se crean solo al aplicar manualmente `database/seed/demo_fst7.up.sql` sobre una DB local ya migrada.

Aplicar seed demo:

```bash
psql 'postgres://roma:roma@localhost:5432/roma?sslmode=disable' -v ON_ERROR_STOP=1 -f database/seed/demo_fst7.up.sql
```

Retirar seed demo:

```bash
psql 'postgres://roma:roma@localhost:5432/roma?sslmode=disable' -v ON_ERROR_STOP=1 -f database/seed/demo_fst7.down.sql
```

## Notas

- `coach_links` es la fuente operativa actual de relacion coach-disciple.
- `master_disciple` queda como legacy/compatibilidad.
- Los usuarios E2E en tests usan correos `*.e2e@example.test` y no son usuarios demo para uso manual.
- Las cuentas `@roma.demo` son seed local/demo opcional y no representan datos productivos.
