# ROMA-029 - Registrar hallazgos demo y ajustar roadmap

Estado: Planned
Tipo: docs
Fecha: 2026-07-04
Autor: Codex

## Objetivo

Documentar los hallazgos de la demo local y ordenar los proximos pasos despues del MVP inicial.

## Contexto

ROMA-028 registro hallazgos de demo sin implementar cambios. Antes de tocar backend o frontend, este WI deja una secuencia clara para transformar esos hallazgos en checkpoints pequenos y verificables.

## Alcance

- Consolidar hallazgos de demo.
- Ajustar `ROADMAP.md` con la siguiente fase post MVP.
- Actualizar `docs/tracking/ROMA_DEV_LOG.md`.
- Dejar ordenados los proximos WIs:
  - ROMA-030 - Coach puede usar entrenamiento personal.
  - ROMA-031 - Mis rutinas: asignadas, propias e importadas.
  - ROMA-032 - Diseno biblioteca de rutinas.
  - ROMA-033 - Implementar biblioteca de rutinas MVP.

## Fuera de alcance

- No cambiar backend.
- No cambiar frontend.
- No tocar migraciones.
- No implementar fixes.
- No agregar features.
- No redisenar UI.

## Archivos probables

- `ROADMAP.md`
- `docs/tracking/ROMA_DEV_LOG.md`
- `docs/work-items/ROMA-029.md`

## Validaciones requeridas

- `git status --short`
- `git diff --check`

No es obligatorio correr tests automatizados porque solo se toca documentacion.

## Criterios de aceptacion

- ROMA-029 documenta hallazgos de demo y ordena proximos pasos.
- ROMA-030 queda definido como el WI donde coach puede usar flujo personal de entrenamiento.
- ROMA-031 queda definido como el WI que ordena "Mis rutinas" distinguiendo asignadas por maestro, propias e importadas.
- ROMA-032 queda definido como el WI de diseno de biblioteca de rutinas antes de implementarla.
- ROMA-033 queda definido como el WI de implementacion MVP de biblioteca de rutinas.

## Resultado esperado del agente

Responder con archivos modificados, validaciones ejecutadas y commit sugerido. No hacer commit automaticamente.
