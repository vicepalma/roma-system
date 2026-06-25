# Roma System

Roma System es una aplicacion para entrenadores y discipulos. El foco actual es entrenamiento, gestion de alumnos y ejecucion de rutinas; no es una red social ni una comunidad publica.

## Stack

- Backend: Go, Gin, GORM, PostgreSQL, JWT.
- Frontend: React, TypeScript, Vite, TanStack Query, Axios, Zustand, Tailwind CSS.
- Base de datos local: PostgreSQL via Docker Compose.

## Estructura

```txt
backend/              API Go
frontend/roma-web/    App React/Vite
database/             Compose, migraciones y seeds
docs/audit/           Auditoria y plan de rescate
```

## Variables de entorno

Backend local esperado:

```env
PORT=8080
ENV=dev
DB_URL=postgres://roma:roma@localhost:5432/roma?sslmode=disable
JWT_SECRET=change-me
ACCESS_TTL_MIN=15
REFRESH_TTL_H=168
DEFAULT_TZ=America/Santiago
```

Frontend local esperado:

```env
VITE_API_BASE=http://localhost:8080
```

Nota: `backend/internal/config/config.go` todavia define `DATABASE_DSN`, pero el servidor actual usa `DB_URL`. Hasta que se unifique la configuracion, usa `DB_URL` para ejecutar `backend/cmd/server`.

## Build y checks

Backend:

```bash
cd backend
go test ./...
```

Frontend:

```bash
cd frontend/roma-web
npm run build
```

## Desarrollo local

Levantar Postgres y backend:

```bash
cd database
docker compose up --build
```

El migrator de Compose esta comentado. Si inicializas una base limpia, aplica las migraciones de `database/backend/migrations` con la herramienta de migraciones que uses en tu entorno.

Frontend en modo dev:

```bash
cd frontend/roma-web
npm run dev
```

## Estado de Fase 0

La Fase 0 busca dejar el repositorio compilable y entendible, sin agregar features nuevas. Las decisiones de producto y permisos siguen documentadas en:

- `docs/audit/ROMA_REPO_AUDIT.md`
- `docs/audit/ROMA_RESCUE_PLAN.md`
- `docs/audit/ROMA_KEEP_DROP_LIST.md`
