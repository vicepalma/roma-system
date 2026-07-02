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
