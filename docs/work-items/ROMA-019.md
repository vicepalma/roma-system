# ROMA-019 - Documentacion operativa y flujo por Work Items

Estado: Done
Tipo: docs
Fecha: 2026-07-04
Autor: Codex

## Objetivo

Crear la documentacion operativa minima para que Roma System se trabaje por Work Items `ROMA-XXX`, con roadmap, reglas para agentes, plantilla y tracking.

## Contexto

El repo ya tenia documentacion base, pero faltaba un flujo explicito para que un agente pueda retomar el proyecto, ejecutar un WI concreto o recomendar el siguiente sin depender de memoria externa.

## Alcance

- Crear o mejorar `ROADMAP.md`.
- Actualizar `AGENTS.md` como manual principal para agentes.
- Crear `docs/work-items/README.md`.
- Crear `docs/work-items/TEMPLATE.md`.
- Crear este Work Item `docs/work-items/ROMA-019.md`.
- Actualizar `docs/tracking/ROMA_DEV_LOG.md`.
- Crear stubs minimos para `ROMA-020` a `ROMA-025`.

## Fuera de alcance

- Cambios funcionales backend.
- Cambios funcionales frontend.
- Migraciones.
- Features nuevas.
- Commits automaticos.

## Archivos probables

- `ROADMAP.md`
- `AGENTS.md`
- `docs/work-items/README.md`
- `docs/work-items/TEMPLATE.md`
- `docs/work-items/ROMA-019.md`
- `docs/work-items/ROMA-020.md`
- `docs/work-items/ROMA-021.md`
- `docs/work-items/ROMA-022.md`
- `docs/work-items/ROMA-023.md`
- `docs/work-items/ROMA-024.md`
- `docs/work-items/ROMA-025.md`
- `docs/tracking/ROMA_DEV_LOG.md`

## Validaciones requeridas

```bash
git status --short
```

No es necesario correr tests porque solo se toca documentacion.

## Criterios de aceptacion

- `ROADMAP.md` resume producto, estado actual, completados y proximos WIs.
- `AGENTS.md` documenta el flujo por Work Item, recuperacion de contexto y validaciones.
- `docs/work-items/` existe con README, template, ROMA-019 y stubs minimos.
- `ROMA_DEV_LOG.md` registra ROMA-019 / CHK-019 y recomienda ROMA-020.

## Riesgos / notas

- Mantener sincronizados roadmap, dev log y WIs en cada checkpoint futuro.

## Resultado esperado del agente

Responder con archivos creados/actualizados, cambios en roadmap y AGENTS, flujo de uso, stubs creados, validacion ejecutada y recomendacion de commit.
