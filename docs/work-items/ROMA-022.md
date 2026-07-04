# ROMA-022 - Editar check-ins propios

Estado: Done
Tipo: feature
Fecha: 2026-07-04
Autor: Codex

## Objetivo

Permitir que el disciple edite sus propios check-ins para corregir errores.

## Alcance

- Permitir edicion solo para check-ins propios del disciple autenticado.
- Campos editables: `checked_at`, `weight_kg`, `notes`.
- Mantener bloqueo para coach, otros disciples y usuarios no vinculados.
- Validar fecha, peso opcional y notas con las mismas reglas de creacion.

## Fuera de alcance

- Borrar check-ins.
- Soft delete.
- Hard delete.
- Edicion por coach.
- Check-ins creados por coach.
- Creacion de check-ins por coach.
- Adjuntos, fotos o mediciones avanzadas.
- Auditoria avanzada.

## Validaciones esperadas

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- E2E/API con `ROMA_E2E_DB_URL`.
- `cd frontend/roma-web && npm run build` si se toca UI.

## Resultado

- `PATCH /api/checkins/:id` permite editar check-ins propios del disciple autenticado.
- Campos editables: `checked_at`, `weight_kg`, `notes`.
- Coach, otros disciples y usuarios no vinculados no pueden editar check-ins ajenos.
- La pantalla `Check-ins` permite editar inline un check-in existente y guardar/cancelar cambios.
- No se implemento borrado, soft delete ni hard delete.

## Validado

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- `cd backend && ROMA_E2E_DB_URL='postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable' GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`
- `cd frontend/roma-web && npm run build`
