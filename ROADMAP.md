# Roadmap Roma System

## Producto

Roma System es software para entrenadores y discipulos. El coach construye y administra entrenamiento para otros; el disciple ejecuta entrenamiento y, si no tiene coach, puede crear rutinas propias usando ejercicios existentes.

No es una comunidad, red social ni sistema de simpatizantes.

## Estado Actual

El proyecto ya tiene backend Go, frontend React/Vite/TypeScript, PostgreSQL con migraciones, Docker Compose, roles `coach` / `disciple`, guards de ownership, tests unitarios backend, E2E/API con DB real, self-training minimo, sesiones activas/cerradas, historial con filtros, check-ins basicos y vista coach de progreso del disciple.

Roma v1 entra en fase de cierre. La biblioteca de rutinas y el responsive mobile basico para los flujos principales forman parte obligatoria de v1.

## Completado

- ROMA-001 - Auditoria inicial.
- ROMA-002 - Fase 0 build verde.
- ROMA-003 - Migraciones y arranque local.
- ROMA-004 - Roles y guards.
- ROMA-005 - Tests de permisos.
- ROMA-006 - E2E/API con DB.
- ROMA-007 - Self-training minimo.
- ROMA-008 - Activacion self-training unica.
- ROMA-009 - Bloquear sesiones sobre assignments inactivos.
- ROMA-010 - UX Entrenar.
- ROMA-011 - Editar rutinas.
- ROMA-012 - Rutina activa y dias descriptivos.
- ROMA-013 - Cierre de sesion y resumen.
- ROMA-014 - Historial de sesiones.
- ROMA-015 - Filtros de historial.
- ROMA-016 - Check-ins basicos.
- ROMA-017 - Check-ins en vista coach.
- ROMA-018 - Vista coach progreso discipulo.
- ROMA-019 - Documentacion operativa y flujo por Work Items.
- ROMA-020 - Pulir historial coach con contexto del alumno.
- ROMA-021 - Filtros/paginacion de check-ins.
- ROMA-022 - Editar check-ins propios.
- ROMA-023 - Biblioteca de rutinas / plantillas Roma, diseno primero.
- ROMA-024 - Demo FST-7 seed opcional.
- ROMA-025 - Hardening final MVP.
- ROMA-026 - Pulido UI general MVP.
- ROMA-027 - Preparacion demo local.
- ROMA-028 - Demo local: prueba manual y lista de ajustes menores.
- ROMA-029 - Registrar hallazgos demo y ajustar roadmap.
- ROMA-030 - Coach puede usar entrenamiento personal.
- ROMA-031 - Mis rutinas: asignadas, propias e importadas.
- ROMA-032 - Diseno biblioteca de rutinas.
- ROMA-033 - Biblioteca de rutinas MVP.
- ROMA-034 - Freeze v1 y criterios mobile.
- ROMA-035 - Responsive base y navegacion.
- ROMA-036 - Responsive flujos disciple/independent.
- ROMA-037 - Responsive flujos coach.

## Modo De Trabajo Actual

- El proyecto se trabaja por Work Items `ROMA-XXX`.
- Cada Work Item debe existir en `docs/work-items/` antes de implementar.
- Si el usuario no entrega un WI, el agente debe leer este roadmap y `docs/tracking/ROMA_DEV_LOG.md`, recomendar el siguiente WI y no implementar todavia.
- No implementar ideas sueltas.
- No mezclar features grandes con bugfixes.
- El tracking historico vive en `docs/tracking/ROMA_DEV_LOG.md`; el alcance ejecutable vive en cada WI.

## Proximo Work Item Recomendado

- ROMA-038 - Hardening, regresion y cierre v1.

WIs finales preparados: `docs/work-items/ROMA-035.md` a `docs/work-items/ROMA-038.md`.

## V1 Freeze Plan

- ROMA-032 - Diseno biblioteca de rutinas.
- ROMA-033 - Biblioteca de rutinas MVP.
- ROMA-034 - Freeze v1 y criterios mobile.
- ROMA-035 - Responsive base y navegacion.
- ROMA-036 - Responsive flujos disciple/independent.
- ROMA-037 - Responsive flujos coach.
- ROMA-038 - Hardening, regresion y cierre v1.

Desde ROMA-034 no entra ninguna feature nueva salvo un bug critico o bloqueante. Cualquier otra idea debe registrarse para backlog post-v1.

## Criterios Mobile Del Freeze

- Viewports objetivo: 360, 390 y 430 px, además de desktop.
- Navegación completa sin depender del sidebar oculto en móvil.
- Formularios y acciones principales utilizables con toque, con estados loading/error/success claros.
- Cards y listas apilables, sin cortes críticos ni solapamiento de acciones.
- Tablas o contenido ancho con adaptación o overflow horizontal controlado.
- Sin overflow horizontal accidental en los flujos críticos.


## Criterio De Cierre V1

- ROMA-032 disena la biblioteca de rutinas antes de implementarla.
- ROMA-033 implementa el MVP de biblioteca de rutinas.
- ROMA-034 congela el alcance y define criterios mobile verificables.
- ROMA-035 deja navegacion y layout base utilizables en telefono.
- ROMA-036 adapta los flujos principales de disciple/independent.
- ROMA-037 adapta los flujos principales de coach.
- ROMA-038 ejecuta hardening, regresion y validacion final de v1.

## Fuera De Alcance Actual

- Comunidad, red social o simpatizantes.
- Biblioteca global de rutinas/plantillas sin WI explicito.
- Fotos, adjuntos o mediciones avanzadas de check-ins.
- Graficos grandes o reporting avanzado.
- Features grandes mezcladas con bugfixes.
- App nativa o PWA avanzada.
- Nuevas features no criticas una vez iniciado ROMA-034.
