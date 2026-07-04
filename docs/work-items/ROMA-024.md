# ROMA-024 - Demo FST-7 seed

Estado: Done
Tipo: chore
Fecha: 2026-07-04
Autor: Codex

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

## Resultado

- Se retiro el seed FST-7 de migraciones oficiales para no cargar datos demo en cualquier DB migrada.
- Se agrego el seed opcional `database/seed/demo_fst7.up.sql` con usuarios demo dedicados `@roma.demo`.
- Se agrego `database/seed/demo_fst7.down.sql` para retirar esos datos demo.
- Se creo el programa privado `[DEMO] FST-7 Hipertrofia - 5 dias` para `fst7.coach@roma.demo`.
- Se agrego un assignment activo para `fst7.disciple@roma.demo`.
- La rutina demo incluye 5 dias, 16 prescripciones y 5 cierres con metodo FST-7.
- Se documento la aplicacion manual del seed demo en `docs/engineering/DEMO_USERS.md` y `docs/engineering/DEVELOPMENT_SETUP.md`.

## Validado

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- Migraciones oficiales `0001` a `0007` aplicadas en Postgres temporal limpio sin datos demo FST-7.
- Seed opcional `database/seed/demo_fst7.up.sql` aplicado manualmente sobre DB temporal ya migrada.
- Smoke SQL en DB temporal: 1 assignment, 5 dias, 16 prescripciones, 5 finishers FST-7.
- Smoke SQL en DB temporal: se pudo crear una sesion y registrar un set sobre el assignment demo.
