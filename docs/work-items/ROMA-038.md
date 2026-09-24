# ROMA-038 — Hardening, regresión y cierre Roma v1

Estado: Done
Tipo: chore
Fecha: 2026-09-24
Autor: Codex

## Objetivo

Cerrar formalmente Roma v1 estable.

## Contexto

Es el checkpoint final del V1 Freeze Plan. No debe convertirse en una nueva fase de features.

## Alcance

- Backend tests.
- E2E/API.
- Frontend build.
- Migraciones desde DB limpia.
- Seed demo opcional.
- Prueba manual coach/disciple/independent.
- Prueba mobile y desktop.
- Corregir solo bugs críticos o bloqueantes.
- Actualizar README, ROADMAP, CURRENT_STATE, demo y tracking.
- Declarar v1 cerrada si cumple criterios del freeze.

## Fuera de alcance

- Features nuevas.
- Mejoras cosméticas no bloqueantes.
- Backlog post-v1.

## Archivos probables

- `README.md`
- `ROADMAP.md`
- `docs/product/ROMA_CURRENT_STATE.md`
- `docs/engineering/LOCAL_DEMO.md`
- `docs/tracking/ROMA_DEV_LOG.md`
- Tests, migraciones o código solo si un bug crítico/bloqueante lo exige.

## Validaciones requeridas

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- `cd backend && ROMA_E2E_DB_URL='postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable' GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`
- `cd frontend/roma-web && npm run build`
- Migraciones desde DB limpia.
- Checklist manual mobile/desktop y demo.

## Criterios de aceptación

- Suite backend, E2E/API y build frontend verdes.
- Migraciones limpias y seed demo opcional verificados.
- Flujos críticos coach/disciple/independent completados en mobile y desktop.
- No quedan bugs críticos o bloqueantes abiertos.
- Documentación y tracking declaran v1 cerrada.

## Riesgos / notas

- Cualquier defecto no crítico pasa al backlog post-v1.
- No hacer commit automáticamente.

## Resultado esperado del agente

Validar y cerrar Roma v1, corregir únicamente bloqueantes, actualizar documentación y no agregar features.

## Resultado

Se ejecutó el hardening disponible sin agregar features. Backend unitario y build frontend quedaron verdes. La validación E2E y migraciones limpias no pudieron ejecutarse porque Docker/PostgreSQL no está disponible en el entorno actual; no se declararon aprobadas esas verificaciones.

## Validación

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...` — OK
- `cd frontend/roma-web && npm run build` — OK
- `git diff --check` — OK
- E2E/API y migraciones limpias — pendientes por falta de Docker/PostgreSQL local


Nota de cierre: las validaciones pendientes siguen bloqueadas por la ausencia de Docker/PostgreSQL en el entorno; mantener Roma v1 abierta hasta ejecutarlas.


## Validaciones completadas

- Migraciones 0001→0008 desde DB limpia: OK.
- Migraciones sin datos demo: OK (`@roma.demo` = 0).
- Seed FST-7 opcional aplicado: OK.
- Retiro del seed FST-7: OK (`@roma.demo` y programas `[DEMO]` = 0).
- E2E/API con PostgreSQL real: OK.
- Healthcheck backend de demo: OK.

## Validación manual pendiente

Abrir el frontend y completar el checklist de `docs/engineering/LOCAL_DEMO.md` en desktop y viewports 360, 390 y 430 px para coach, disciple e independent.
