# Development Setup

## Requisitos

- Go instalado para backend.
- Node/npm para frontend.
- Docker y Docker Compose para Postgres, migrator y backend local.

## Variables Importantes

Backend:

```env
PORT=8080
ENV=dev
DB_URL=postgres://roma:roma@localhost:5432/roma?sslmode=disable
JWT_SECRET=change-me
ACCESS_TTL_MIN=15
REFRESH_TTL_H=168
DEFAULT_TZ=America/Santiago
```

Frontend:

```env
VITE_API_BASE=http://localhost:8080
```

## Levantar Entorno Local

Desde la raiz:

```bash
docker compose -f database/docker-compose.yml up --build
```

Esto levanta Postgres, corre migraciones y levanta backend.

Frontend:

```bash
cd frontend/roma-web
npm run dev
```

## Seeds Demo Opcionales

Las migraciones oficiales dejan el esquema y el seed minimo operativo. Los datos demo no se aplican automaticamente.

Para una ruta completa de demo local, ver `docs/engineering/LOCAL_DEMO.md`.

Para cargar la demo FST-7 sobre una DB local ya migrada:

```bash
psql 'postgres://roma:roma@localhost:5432/roma?sslmode=disable' -v ON_ERROR_STOP=1 -f database/seed/demo_fst7.up.sql
```

Esto crea las cuentas `fst7.coach@roma.demo` y `fst7.disciple@roma.demo` con password `secret123`, mas el programa `[DEMO] FST-7 Hipertrofia - 5 dias` asignado al disciple demo.

Para retirar esos datos demo:

```bash
psql 'postgres://roma:roma@localhost:5432/roma?sslmode=disable' -v ON_ERROR_STOP=1 -f database/seed/demo_fst7.down.sql
```

## Reset Local

Si necesitas resetear la DB local de compose:

```bash
docker compose -f database/docker-compose.yml down -v
docker compose -f database/docker-compose.yml up --build
```

No ejecutes resets sobre bases que contengan datos reales.

## Advertencias

- `ROMA_E2E_DB_URL` se usa para E2E y la suite trunca tablas.
- No apuntes `ROMA_E2E_DB_URL` a una DB real o compartida.
- Archivos temporales como `creds.txt` y `sesion.txt` no deben commitearse.
- El backend actual usa `DB_URL` para `cmd/server`.
