# ROMA-020 - Pulir historial coach con contexto del alumno

Estado: Done
Tipo: feature
Fecha: 2026-07-04
Autor: Codex

## Objetivo

Cuando el coach vea historial de un alumno, mostrar claramente nombre del alumno y filtros aplicados.

## Alcance

- Revisar vista de historial usada por coach.
- Mostrar contexto visible del disciple seleccionado.
- Asegurar que filtros activos sean claros para el coach.
- Mantener guards y ownership en backend.

## Fuera de alcance

- Nuevos graficos o reporting avanzado.
- Cambios al modelo de permisos.
- Biblioteca de rutinas.

## Validaciones esperadas

- `cd frontend/roma-web && npm run build`
- Si se toca backend o contrato API: `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- Si se toca contrato API: E2E/API con `ROMA_E2E_DB_URL`.

## Resultado

- La vista `/history?disciple_id=...` ahora muestra contexto visible del disciple seleccionado con nombre, email e ID.
- Los enlaces desde detalle de disciple pasan nombre/email por estado de navegacion.
- Si se abre la URL directamente, el historial coach intenta resolver el disciple desde `/api/coach/disciples`.
- Los filtros activos de sesiones quedan visibles como resumen junto al total de resultados.

## Validado

- `cd frontend/roma-web && npm run build`
