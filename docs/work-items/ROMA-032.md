# ROMA-032 - Diseno biblioteca de rutinas

Estado: Planned
Tipo: docs
Fecha: 2026-07-04
Autor: Codex

## Objetivo

Disenar la biblioteca global de rutinas antes de implementarla. Debe distinguir plantilla global, rutina personal y asignacion de maestro.

## Contexto

La demo local mostro la necesidad de una biblioteca de rutinas creadas/publicadas por maestros o admin, visibles para disciples e independientes. Antes de tocar backend/frontend, hay que cerrar reglas de producto y modelo tecnico minimo.

## Alcance

- Definir conceptos: plantilla global, rutina personal, rutina importada y asignacion de maestro.
- Definir ownership, visibilidad y permisos de publicacion.
- Definir flujo de importacion como copia personal.
- Definir que eliminar una copia personal no borra la plantilla global.
- Definir alcance MVP para ROMA-033.

## Fuera de alcance

- Implementar backend.
- Implementar frontend.
- Crear migraciones.
- Publicacion social, marketplace o comunidad.
- Adjuntos, fotos o reporting avanzado.

## Archivos probables

- `docs/product/ROMA_TEMPLATE_LIBRARY_DESIGN.md`
- `docs/work-items/ROMA-032.md`
- `docs/tracking/ROMA_DEV_LOG.md`

## Validaciones requeridas

- `git status --short`
- `git diff --check`

No es obligatorio correr tests automatizados porque este WI es de documentacion/diseno.

## Criterios de aceptacion

- Queda definida la diferencia entre plantilla global, rutina personal y asignacion de maestro.
- Queda definido quien puede crear/publicar plantillas.
- Queda definido como disciples/independientes importan una copia personal.
- Queda definido el alcance implementable de ROMA-033.

## Riesgos / notas

- Mantener el diseno acotado al MVP.
- No introducir comunidad, marketplace ni biblioteca social.

## Resultado esperado del agente

Documentar el diseno, actualizar tracking y no implementar cambios funcionales.
