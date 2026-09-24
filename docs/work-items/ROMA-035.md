# ROMA-035 — Responsive base + navegación

Estado: Planned
Tipo: feature
Fecha: 2026-09-24
Autor: Codex

## Objetivo

Hacer usable la estructura principal en mobile.

## Contexto

ROMA-034 congeló los criterios mobile de v1. Este WI adapta el shell compartido antes de trabajar los flujos específicos de disciple y coach.

## Alcance

- Navegación visible en 360/390/430 px.
- Reemplazo mobile del sidebar oculto.
- Header y contenedor responsive.
- Theme/logout utilizables con toque.
- Eliminar overflow horizontal accidental.
- Preservar desktop.

## Fuera de alcance

- Nuevas features.
- Rediseño visual completo.
- App nativa o PWA avanzada.

## Archivos probables

- `frontend/roma-web/src/App.tsx`
- `frontend/roma-web/src/components/layout/Sidebar.tsx`
- `frontend/roma-web/src/index.css`
- Componentes compartidos de layout/UI.

## Validaciones requeridas

- `cd frontend/roma-web && npm run build`
- Prueba manual en 360, 390, 430 px y desktop.

## Criterios de aceptación

- Todos los enlaces principales son accesibles en mobile.
- Header, navegación, theme y logout no se solapan ni salen del viewport.
- No hay overflow horizontal accidental en el shell.
- Desktop conserva su navegación y layout actuales.

## Riesgos / notas

- Mantener el cambio limitado al shell; los flujos internos corresponden a ROMA-036/037.

## Resultado esperado del agente

Implementar solo el shell responsive, validar build/prueba manual, actualizar tracking y no hacer commit automáticamente.
