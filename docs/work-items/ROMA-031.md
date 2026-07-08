# ROMA-031 - Mis rutinas: asignadas, propias e importadas

Estado: Done
Tipo: feature
Fecha: 2026-07-04
Autor: Codex

## Objetivo

Ordenar Mis rutinas para mostrar rutinas asignadas por maestro, rutinas propias y rutinas importadas, con permisos claros.

## Contexto

ROMA-028 detecto que el disciple debe poder ver rutinas asignadas por maestro en Mis rutinas, pero sin eliminarlas desde su seccion personal. Tambien debe distinguir rutinas propias o importadas, que si puede administrar dentro de su espacio personal.

## Alcance

- Definir categorias visibles en Mis rutinas: asignadas por maestro, propias e importadas.
- Mostrar rutinas asignadas por maestro sin accion de eliminacion personal.
- Mantener acciones de edicion/eliminacion solo para rutinas personales cuando corresponda.
- Preparar el modelo visual y de permisos para rutinas importadas.
- Mantener guards backend como fuente de verdad.

## Fuera de alcance

- Implementar biblioteca de rutinas.
- Crear importacion desde biblioteca.
- Permitir que disciple elimine o modifique rutinas asignadas por maestro.
- Cambios grandes de UI fuera de Mis rutinas.

## Archivos probables

- Backend de programs/assignments si el contrato actual no expone los datos necesarios.
- Frontend de Mis rutinas / Programas.
- Tests backend/E2E si se toca contrato o permisos.
- `docs/tracking/ROMA_DEV_LOG.md`

## Validaciones requeridas

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...` si se toca backend.
- E2E/API con `ROMA_E2E_DB_URL` si se toca contrato o permisos.
- `cd frontend/roma-web && npm run build` si se toca frontend.

## Criterios de aceptacion

- Mis rutinas distingue claramente asignadas, propias e importadas.
- Rutinas asignadas por maestro son visibles para el disciple.
- Rutinas asignadas por maestro no se pueden eliminar desde la seccion personal.
- Rutinas propias/importadas conservan acciones permitidas segun ownership.
- UI y backend quedan alineados en permisos.

## Riesgos / notas

- No confiar solo en ocultar botones: backend debe proteger acciones.
- Evitar mezclar este WI con biblioteca o importacion real.

## Resultado esperado del agente

Implementar solo este WI, validar segun archivos tocados, actualizar tracking y no hacer commit automaticamente.

## Resultado

- `/api/programs` ahora devuelve rutinas/programas propios y rutinas asignadas al usuario autenticado.
- El contrato agrega metadata de origen y permisos: `source`, `assignment_id`, `assigned_by`, `is_active`, `can_edit`, `can_delete` y `can_activate`.
- Mis rutinas separa visualmente rutinas asignadas por maestro, propias/programas e importadas.
- Las rutinas asignadas por maestro son visibles, pero no muestran acciones de edicion, eliminacion ni cambios de estructura.
- Las rutinas propias y futuras importadas conservan acciones segun los permisos entregados por backend.
- E2E cubre que una rutina asignada aparece en Mis rutinas con `source=assigned` y no puede eliminarse desde el disciple.

## Validado

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- `cd backend && ROMA_E2E_DB_URL='postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable' GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`
- `cd frontend/roma-web && npm run build`
