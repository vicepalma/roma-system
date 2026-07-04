# ROMA-021 - Filtros/paginacion de check-ins

Estado: Done
Tipo: feature
Fecha: 2026-07-04
Autor: Codex

## Objetivo

Agregar filtros simples por fecha y limite/paginacion basica a check-ins.

## Alcance

- Definir filtros minimos para check-ins.
- Aplicarlos en backend respetando ownership.
- Exponer controles simples en UI si corresponde.
- Mantener vista coach solo para disciples vinculados.

## Fuera de alcance

- Edicion o borrado de check-ins.
- Fotos, adjuntos o mediciones avanzadas.
- Graficos grandes.

## Validaciones esperadas

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- E2E/API con `ROMA_E2E_DB_URL` si se toca contrato o permisos.
- `cd frontend/roma-web && npm run build` si se toca UI.

## Resultado

- `GET /api/checkins` y `GET /api/coach/disciples/:id/checkins` aceptan filtros `from`, `to`, `limit` y `offset`.
- Los filtros por fecha se aplican en backend respetando ownership.
- La UI de `Check-ins` permite filtrar por fecha, elegir cantidad por pagina y navegar anterior/siguiente.
- E2E cubre filtros de check-ins propios y vista coach vinculada.

## Validado

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- `cd backend && ROMA_E2E_DB_URL='postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable' GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`
- `cd frontend/roma-web && npm run build`
