# ROMA-026 - Pulido UI general MVP

Estado: Done
Tipo: chore
Fecha: 2026-07-04
Autor: Codex

## Objetivo

Definir y ejecutar un pulido UI general para el MVP, mejorando claridad y consistencia sin redisenar toda la app ni agregar features nuevas.

## Alcance

- Mejorar textos vacios, loading y error en pantallas principales.
- Revisar navegacion disciple/coach para que los flujos principales sean claros.
- Revisar pantallas principales del MVP:
  - Entrenar.
  - Mis rutinas / Programas.
  - Historial.
  - Check-ins.
  - Vista coach de discipulos/progreso/asignaciones.
- Limpiar duplicaciones visuales evidentes.
- Mejorar consistencia de cards, botones, estados y jerarquia visual.
- Mantener cambios chicos, verificables y alineados con la UI existente.

## Fuera de alcance

- Redisenar toda la app.
- Cambiar backend.
- Cambiar contratos API.
- Agregar features nuevas.
- Crear comunidad, red social, reporting avanzado o biblioteca global de plantillas.
- Agregar fotos, adjuntos o mediciones avanzadas.

## Validaciones esperadas

- `cd frontend/roma-web && npm run build`
- `git status --short`

Si durante el pulido aparece una necesidad backend o de contrato API, documentarla como pendiente y no implementarla en este WI.

## Resultado

- Se agrego `QueryState` como componente UI compartido para estados vacios/loading/error.
- Se aplicaron estados mas consistentes en Entrenar, Historial, Check-ins, Dashboard coach y Asignaciones.
- Se mejoraron textos de estados vacios/loading/error sin agregar flujos nuevos.
- Se corrigio una inconsistencia visual menor en la tabla de asignaciones.
- No se tocaron backend ni contratos API.

## Validado

- `cd frontend/roma-web && npm run build`
- `git diff --check`
