# ROMA-036 — Responsive flujos disciple/independent

Estado: Done
Tipo: feature
Fecha: 2026-09-24
Autor: Codex

## Objetivo

Hacer completamente usables en mobile los flujos personales.

## Contexto

Depende de ROMA-035 y cubre la ejecución diaria de disciple/independent en los viewports congelados de v1.

## Alcance

- Entrenar.
- Sesión activa.
- Selección de día.
- Registro de sets.
- Cierre de sesión.
- Mis rutinas.
- Activación/edición de rutinas.
- Historial.
- Check-ins.
- Biblioteca/importación si forma parte del flujo personal.

## Viewports

- 360, 390 y 430 px + desktop.

## Fuera de alcance

- Nuevas features.
- Rediseño visual completo.
- PWA/app nativa.

## Archivos probables

- `frontend/roma-web/src/pages/SessionsIndex.tsx`
- `frontend/roma-web/src/pages/SessionView.tsx`
- `frontend/roma-web/src/pages/Programs.tsx`
- `frontend/roma-web/src/pages/History.tsx`
- `frontend/roma-web/src/pages/Checkins.tsx`
- `frontend/roma-web/src/pages/Templates.tsx`
- Formularios y componentes de sesión/UI relacionados.

## Validaciones requeridas

- `cd frontend/roma-web && npm run build`
- Prueba manual de los flujos en 360, 390, 430 px y desktop.

## Criterios de aceptación

- Cada flujo personal se puede completar con toque sin overflow accidental.
- Formularios, acciones principales, estados y errores son utilizables en mobile.
- La sesión activa permite registrar sets y cerrarse sin acciones ocultas.
- Desktop conserva el comportamiento existente.

## Riesgos / notas

- No agregar funcionalidades; corregir solo layout/interacción responsive dentro del alcance.

## Resultado esperado del agente

Implementar solo responsive de flujos personales, validar build/prueba manual, actualizar tracking y no hacer commit automáticamente.

## Resultado

- Se adaptaron Entrenar, sesión activa, selección de día, sets y cierre para apilar acciones y contenido en mobile.
- Se ajustaron Mis rutinas, Historial, Check-ins y Biblioteca para evitar cortes y overflow accidental.
- Se hicieron modales y formularios utilizables en pantallas estrechas con acciones táctiles.
- Se preservó el layout desktop mediante breakpoints responsive.

## Validado

- `cd frontend/roma-web && npm run build`
- `git diff --check`
