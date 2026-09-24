# ROMA-037 — Responsive flujos coach

Estado: Done
Tipo: feature
Fecha: 2026-09-24
Autor: Codex

## Objetivo

Hacer usables en mobile los flujos principales de Maestro.

## Contexto

Depende de ROMA-035 y adapta las tareas de gestión y seguimiento coach sin ampliar el producto.

## Alcance

- Dashboard.
- Discípulos.
- Detalle/progreso de discípulo.
- Asignaciones.
- Programas.
- Ejercicios.
- Biblioteca.
- Entrenamiento personal del coach si corresponde.

## Viewports

- 360, 390 y 430 px + desktop.

## Fuera de alcance

- Nuevas features.
- Admin.
- Rediseño amplio.

## Archivos probables

- `frontend/roma-web/src/pages/Dashboard.tsx`
- `frontend/roma-web/src/pages/DiscipleDetail.tsx`
- `frontend/roma-web/src/pages/Assignments.tsx`
- `frontend/roma-web/src/pages/Programs.tsx`
- `frontend/roma-web/src/pages/Exercises.tsx`
- `frontend/roma-web/src/pages/Templates.tsx`
- Componentes de dashboard, tablas, formularios y layout relacionados.

## Validaciones requeridas

- `cd frontend/roma-web && npm run build`
- Prueba manual coach en 360, 390, 430 px y desktop.

## Criterios de aceptación

- El coach puede navegar, revisar discípulos, gestionar programas/asignaciones y consultar progreso en mobile.
- Tablas, cards y formularios no cortan acciones críticas ni producen overflow accidental.
- Desktop conserva el comportamiento existente.

## Riesgos / notas

- Mantener ownership y permisos sin cambios; este WI solo adapta presentación/interacción responsive.

## Resultado esperado del agente

Implementar solo responsive de flujos coach, validar build/prueba manual, actualizar tracking y no hacer commit automáticamente.

## Resultado

Se ajustaron headers, formularios, acciones táctiles y tablas de los flujos coach para 360/390/430 px, preservando desktop y la lógica existente.

## Validación

- `cd frontend/roma-web && npm run build`
- `git diff --check`
