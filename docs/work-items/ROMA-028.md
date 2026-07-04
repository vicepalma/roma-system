# ROMA-028 - Demo local: prueba manual y lista de ajustes menores

Estado: Planned
Tipo: chore
Fecha: 2026-07-04
Autor: Codex

## Objetivo

Definir una prueba manual completa del MVP local usando `LOCAL_DEMO.md` y el seed FST-7 opcional, registrando bugs, fricciones y ajustes menores sin implementar cambios todavia.

## Alcance

- Levantar backend/frontend local.
- Aplicar seed demo FST-7 opcional.
- Probar flujo coach.
- Probar flujo disciple.
- Probar entrenamiento, cierre de sesion, historial y check-ins.
- Registrar problemas encontrados.
- Proponer candidatos para ROMA-029.

## Fuera de alcance

- No cambiar backend.
- No cambiar frontend.
- No tocar migraciones.
- No implementar fixes.
- No agregar features.
- No redisenar UI.

## Archivos probables

- `docs/work-items/ROMA-028.md`
- `docs/tracking/ROMA_DEV_LOG.md`

Si durante la prueba aparecen problemas, registrarlos en este WI o en el tracking sin implementar cambios.

## Validaciones requeridas

- Seguir `docs/engineering/LOCAL_DEMO.md`.
- Registrar resultado de la prueba manual.
- `git status --short`

No es obligatorio correr tests automatizados porque este WI no debe cambiar backend, frontend ni migraciones.

## Criterios de aceptacion

- La demo local fue levantada siguiendo la guia vigente.
- El seed FST-7 opcional fue aplicado o se documento por que no se aplico.
- Se probo el flujo coach.
- Se probo el flujo disciple.
- Se probo entrenamiento, cierre de sesion, historial y check-ins.
- Bugs, fricciones y ajustes menores quedaron registrados.
- Hay candidatos claros para ROMA-029.

## Riesgos / notas

- No convertir hallazgos de demo en fixes dentro de este WI.
- Mantener cambios chicos y documentales; cualquier fix requiere WI posterior.

## Resultado esperado del agente

Responder con resumen de la prueba manual, problemas encontrados, candidatos para ROMA-029, validaciones ejecutadas y recomendacion de siguiente checkpoint.
