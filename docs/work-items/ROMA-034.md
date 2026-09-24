# ROMA-034 — Freeze v1 + criterios mobile

Estado: Done
Tipo: docs
Fecha: 2026-09-24
Autor: Codex

## Objetivo

Congelar oficialmente el alcance de Roma v1 y definir criterios verificables de aceptación mobile.

## Alcance

- Fijar viewports objetivo: 360, 390, 430 px y desktop.
- Listar flujos críticos de disciple/independent y coach.
- Definir criterios de navegación, formularios, cards, overflow y acciones principales.
- Establecer que desde ROMA-034 no entran nuevas features.
- Permitir únicamente bugs críticos o bloqueantes durante el cierre de v1.
- Dejar fuera app nativa, PWA avanzada y nuevas features.

## Flujos críticos

### Disciple / independent

- Login y navegación autenticada.
- Ver/activar rutina y abrir Biblioteca.
- Iniciar sesión de entrenamiento, registrar sets y cerrar sesión.
- Revisar Historial.
- Crear, editar y revisar Check-ins.
- Importar una plantilla y verla en Mis rutinas.

### Coach

- Login y navegación autenticada.
- Dashboard y lista de disciples.
- Crear/editar programas y asignaciones.
- Revisar progreso, sesiones, historial y check-ins de un disciple.
- Crear/publicar e importar plantillas.
- Usar entrenamiento personal sin mezclarlo con datos de disciples.

## Criterios mobile verificables

- La aplicación es usable sin zoom horizontal en 360, 390 y 430 px de ancho, además de desktop.
- La navegación permite llegar a todos los flujos críticos sin depender de un sidebar oculto.
- Formularios: campos legibles, controles utilizables con toque, botones visibles y estados loading/error/success comprensibles.
- Cards y listas: contenido ajusta o apila; títulos, acciones y metadatos no se cortan de forma crítica.
- Tablas o contenido ancho: se adaptan a tarjetas/listas o tienen overflow horizontal controlado y explícito.
- Acciones principales quedan visibles, alcanzables y no se solapan; acciones destructivas mantienen confirmación.
- No aparecen scrollbars horizontales accidentales en el viewport principal.
- Los flujos críticos se pueden completar en cada viewport objetivo con teclado móvil y toque.

## Regla De Freeze

Desde ROMA-034 no entra ninguna feature nueva. Solo pueden entrar bugs críticos o bloqueantes que impidan seguridad, navegación o completar un flujo crítico. Toda idea no crítica pasa al backlog post-v1.

## Fuera de alcance

- App nativa.
- PWA avanzada.
- Nuevas features no incluidas en ROMA-032 a ROMA-038.
- Comunidad, reporting avanzado, adjuntos o cambios de modelo no necesarios para cerrar v1.

## Validaciones requeridas

- `git status --short`
- `git diff --check`

No corresponde ejecutar tests de backend/frontend porque este WI solo actualiza documentación.

## Resultado

- Se congeló el alcance de Roma v1 y se definieron viewports, flujos críticos y criterios mobile verificables.
- Se registró la regla de cambios posteriores limitada a bugs críticos o bloqueantes.
- No se modificaron backend ni frontend.
