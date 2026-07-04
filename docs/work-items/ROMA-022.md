# ROMA-022 - Editar check-ins propios

Estado: Planned
Tipo: feature

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
