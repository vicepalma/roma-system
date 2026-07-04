# Testing Guide

## Backend Unitario

```bash
cd backend
GOCACHE=/tmp/roma-go-cache go test ./...
```

Cubre servicios/handlers/guards y debe correrse ante cambios backend.

## E2E/API Con DB Real De Test

Requiere Postgres con migraciones aplicadas y una DB dedicada, por ejemplo `roma_e2e`.

```bash
cd backend
ROMA_E2E_DB_URL='postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable' GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1
```

La suite E2E limpia/trunca datos. Nunca usar con DB real.

Cubre auth, roles, ejercicios, programas, assignments, self-training, sesiones, sets, historial, check-ins y ownership principal.

## Frontend Build

```bash
cd frontend/roma-web
npm run build
```

Debe correrse ante cambios en frontend o contratos consumidos por frontend.

## Seeds Demo Opcionales

Los seeds demo no forman parte de las migraciones oficiales. Para validar la demo FST-7, primero aplica migraciones oficiales sobre una DB local/test y luego:

```bash
psql 'postgres://roma:roma@localhost:5432/roma?sslmode=disable' -v ON_ERROR_STOP=1 -f database/seed/demo_fst7.up.sql
```

## Cuando Correr Cada Validacion

- Solo docs: `git status --short`.
- Backend o permisos: backend unitario y E2E/API.
- Frontend: `npm run build`.
- Migraciones: validar migraciones desde DB limpia y correr E2E/API si afectan contrato.
