# ROMA-025 - Hardening final MVP

Estado: Done
Tipo: chore
Fecha: 2026-07-04
Autor: Codex

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

## Resultado

- Se revisaron permisos y ownership criticos de assignments, sesiones, programas, historial, check-ins y flujo coach/disciple.
- Se endurecieron rutas de coach para que calendario y activacion de assignments bajo `/api/coach/assignments/:id` requieran rol `coach`.
- Se agrego cobertura E2E para bloquear a disciples en activacion/calendario de assignments por rutas coach.
- Se actualizo documentacion de estado, testing y README para reflejar check-ins editables y seed demo FST-7 opcional.
- Se creo `docs/engineering/ROMA_MVP_HARDENING.md` con riesgos pendientes y recomendacion siguiente.

## Validado

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- `cd backend && ROMA_E2E_DB_URL='postgres://roma:roma@localhost:55432/roma_e2e?sslmode=disable' GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1` contra Postgres temporal.
- `cd frontend/roma-web && npm run build`
