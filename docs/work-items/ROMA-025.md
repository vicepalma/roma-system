# ROMA-025 - Hardening final MVP

Estado: Planned
Tipo: chore

## Objetivo

Revisar seguridad, validaciones, UX critica, README y flujo demo antes de considerar MVP.

## Alcance

- Revisar permisos y ownership criticos.
- Revisar errores visibles para flujos principales.
- Revisar README/setup/testing para demo local.
- Crear lista corta de bloqueantes MVP si aparecen.

## Fuera de alcance

- Features nuevas grandes.
- Red social/comunidad.
- Reporting avanzado.
- Reescrituras amplias.

## Validaciones esperadas

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- E2E/API con `ROMA_E2E_DB_URL`
- `cd frontend/roma-web && npm run build`
