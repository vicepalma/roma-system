# ROMA-024 - Demo FST-7 seed

Estado: Planned
Tipo: chore

## Objetivo

Cargar una pauta demo FST-7 como programa/rutina seed para probar Roma con datos reales.

## Alcance

- Definir datos seed realistas para una demo local.
- Mantenerlos separados de datos productivos.
- Verificar que la rutina demo se pueda revisar y ejecutar.

## Fuera de alcance

- Biblioteca completa de plantillas.
- Reporting avanzado.
- Cambios grandes en UI.

## Validaciones esperadas

- Validar migraciones/seeds desde DB limpia si se modifica seed oficial.
- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- E2E/API con `ROMA_E2E_DB_URL` si afecta datos usados por tests.
- `cd frontend/roma-web && npm run build` si se toca UI.
