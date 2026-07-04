# ROMA-022 - Edicion/borrado de check-ins

Estado: Planned
Tipo: feature

## Objetivo

Permitir que el disciple edite o borre sus propios check-ins, si se decide.

## Alcance

- Definir reglas de edicion/borrado para check-ins propios.
- Implementar endpoints y UI solo si se confirma el alcance.
- Mantener bloqueo para coach, otros disciples y usuarios no vinculados.

## Fuera de alcance

- Check-ins creados por coach.
- Adjuntos, fotos o mediciones avanzadas.
- Auditoria compleja de cambios.

## Validaciones esperadas

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- E2E/API con `ROMA_E2E_DB_URL`.
- `cd frontend/roma-web && npm run build` si se toca UI.
