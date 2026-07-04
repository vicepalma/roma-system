# ROMA-021 - Filtros/paginacion de check-ins

Estado: Planned
Tipo: feature

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
