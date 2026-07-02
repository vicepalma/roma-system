# Roma System

Roma System es una aplicacion para entrenadores y discipulos. El foco actual es entrenamiento, gestion de alumnos, ejecucion de rutinas, historial y seguimiento basico. No es una red social ni una comunidad publica.

## Stack

- Backend: Go, Gin, GORM, PostgreSQL, JWT.
- Frontend: React, TypeScript, Vite, TanStack Query, Axios, Zustand, Tailwind CSS.
- Infra local: Docker Compose para Postgres, migrator y backend.

## Estado Actual

Ya existen roles `coach` / `disciple`, guards de ownership, self-training minimo, sesiones activas/cerradas, resumen de entrenamiento, historial con filtros backend, check-ins basicos y vista de progreso del disciple para coach.

## Comandos Rapidos

Levantar Postgres, migraciones y backend:

```bash
docker compose -f database/docker-compose.yml up --build
```

Levantar frontend:

```bash
cd frontend/roma-web
npm run dev
```

Tests backend:

```bash
cd backend
GOCACHE=/tmp/roma-go-cache go test ./...
```

E2E/API con DB real de test:

```bash
cd backend
ROMA_E2E_DB_URL='postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable' GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1
```

Build frontend:

```bash
cd frontend/roma-web
npm run build
```

## Documentos Principales

- `AGENTS.md`
- `ROADMAP.md`
- `docs/product/ROMA_CURRENT_STATE.md`
- `docs/product/ROMA_PRODUCT_RULES.md`
- `docs/engineering/DEVELOPMENT_SETUP.md`
- `docs/engineering/TESTING_GUIDE.md`
- `docs/engineering/DEMO_USERS.md`
- `docs/tracking/ROMA_DEV_LOG.md`
