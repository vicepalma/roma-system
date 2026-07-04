# ROMA-027 - Preparacion demo local

Estado: Done
Tipo: docs
Fecha: 2026-07-04
Autor: Codex

## Objetivo

Dejar una demo local facil de levantar y probar.

## Contexto

El MVP inicial llego al cierre de ROMA-026. Antes de seguir con nuevas funcionalidades, conviene documentar una ruta de demo local que permita mostrar el flujo coach/disciple con datos seed opcionales.

## Alcance

- Documentar pasos para levantar backend/frontend.
- Documentar como aplicar seed FST-7 opcional.
- Documentar usuarios demo.
- Documentar flujo demo sugerido para coach y disciple.
- Agregar checklist manual de demo.

## Fuera de alcance

- No cambiar backend.
- No cambiar frontend.
- No crear features.
- No automatizar deploy.
- No tocar migraciones.

## Archivos probables

- `README.md`
- `docs/engineering/DEVELOPMENT_SETUP.md`
- `docs/engineering/LOCAL_DEMO.md`
- `docs/tracking/ROMA_DEV_LOG.md`
- `docs/work-items/ROMA-027.md`

## Validaciones requeridas

- `git status --short`

No es obligatorio correr tests si solo se toca documentacion.

## Criterios de aceptacion

- Existe una guia clara para levantar demo local.
- La guia explica cuando y como aplicar el seed FST-7 opcional.
- Los usuarios demo quedan documentados.
- Hay un flujo manual sugerido para coach y disciple.
- Hay una checklist manual breve para validar la demo.

## Riesgos / notas

- Mantener el seed demo fuera de migraciones oficiales.
- No mezclar mejoras de producto con preparacion de demo.

## Resultado esperado del agente

- Actualizar documentacion operativa de demo local.
- Actualizar este WI con resultado y validaciones.
- Actualizar `docs/tracking/ROMA_DEV_LOG.md`.
- No hacer commit automaticamente.

## Resultado

- Se creo `docs/engineering/LOCAL_DEMO.md` con pasos para levantar backend/frontend.
- Se documento como aplicar y retirar el seed FST-7 opcional.
- Se documentaron usuarios demo y password comun.
- Se agrego flujo demo sugerido para coach y disciple.
- Se agrego checklist manual de demo.
- Se enlazo la guia desde `README.md` y `docs/engineering/DEVELOPMENT_SETUP.md`.
- No se tocaron backend, frontend ni migraciones.

## Validado

- `git status --short`
- `git diff --check`
