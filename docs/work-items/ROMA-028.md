# ROMA-028 - Demo local: prueba manual y lista de ajustes menores

Estado: Done
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

## Resultado

Se registraron hallazgos de la demo local sin implementar cambios funcionales.

Hallazgos:

- El Maestro deberia poder usar funcionalidades de entrenamiento personal/self-training, porque tambien puede entrenarse a si mismo.
- En "Mis rutinas", el Discipulo deberia ver rutinas asignadas por maestro, pero sin opcion de eliminarlas si vienen del maestro.
- El Discipulo/Independiente si puede eliminar rutinas propias o importadas desde biblioteca, pero solo de su seccion personal.
- El Independiente/self-training funciona correctamente.
- Falta Biblioteca de rutinas: rutinas creadas/publicadas por maestros o admin, visibles para Discipulos e Independientes.
- Al agregar una rutina desde biblioteca, debe copiarse/importarse a "Mis rutinas"; eliminarla desde ahi no debe borrar la plantilla global.

## Candidatos para proximos WIs

- ROMA-029 - Reglas de producto para rutinas personales, asignadas e importadas.
- ROMA-030 - Permitir self-training para coach.
- ROMA-031 - Mis rutinas: mostrar rutinas asignadas por maestro sin eliminacion.
- ROMA-032 - Biblioteca de rutinas: diseno tecnico y modelo de copia/importacion.
- ROMA-033 - Biblioteca de rutinas: MVP de lectura e importacion.

## Validado

- `git status --short`

No se ejecutaron tests automatizados porque solo se actualizo documentacion.
