# ROMA Dev Log

## Estado actual
- Fase actual: cierre de Roma v1 / V1 Freeze Plan definido.
- Objetivo actual: ejecutar el cierre mobile y hardening final bajo alcance congelado.
- Ultimo checkpoint completado: ROMA-035 / CHK-035 - Responsive base + navegacion.
- Proximo checkpoint sugerido: ROMA-036 - Responsive flujos disciple/independent.

## Decisiones activas
- [2026-06-24] Decision: no partir desde cero; rescatar repo con estabilizacion previa.
- [2026-06-24] Decision: alinear migraciones oficiales con backend actual creando `coach_links` y `program_versions` en Fase 0.
- [2026-06-24] Decision: activar `migrator` en Compose para DB limpia reproducible antes de backend.
- [2026-06-25] Decision: `coach_links` queda como fuente operativa; `master_disciple` queda legacy/compatibilidad.
- [2026-06-29] Decision: self-training usa `programs.kind='self_training'` y self-assignment propio; no usa self-link en `coach_links`.
- [2026-09-24] Decision: ROMA-032 y ROMA-033 forman parte de Roma v1.
- [2026-09-24] Decision: Roma v1 debe ser usable en telefonos mediante responsive mobile basico; app nativa y PWA avanzada quedan fuera.
- [2026-09-24] Decision: el V1 Freeze Plan queda limitado a ROMA-032 a ROMA-038.
- [2026-09-24] Decision: desde ROMA-034 no entra ninguna feature nueva salvo bug critico o bloqueante; el resto va a backlog post-v1.
- [2026-09-24] Decision: los viewports de aceptación mobile son 360, 390 y 430 px, además de desktop.
- [2026-09-24] Decision: ROMA-035 a ROMA-038 son los WIs finales de cierre de Roma v1.

## V1 Freeze Plan
- ROMA-032 - Diseno biblioteca de rutinas.
- ROMA-033 - Biblioteca de rutinas MVP.
- ROMA-034 - Freeze v1 y criterios mobile.
- ROMA-035 - Responsive base y navegacion.
- ROMA-036 - Responsive flujos disciple/independent.
- ROMA-037 - Responsive flujos coach.
- ROMA-038 - Hardening, regresion y cierre v1.

## Checkpoints
### CHK-001 - Auditoria inicial
Estado: Completado
Resultado: Auditoria, plan de rescate y keep/drop list creados.
Archivos: `docs/audit/ROMA_REPO_AUDIT.md`, `docs/audit/ROMA_RESCUE_PLAN.md`, `docs/audit/ROMA_KEEP_DROP_LIST.md`.

### CHK-002 - Fase 0 build verde
Estado: Completado
Objetivo:
- Backend `go test ./...` verde.
- Frontend `npm run build` verde.
- README raiz minimo.
- Notas Fase 0.
Resultado: Backend y frontend compilan; se agrego README raiz y migracion `0003_phase0_schema_alignment`.
Validado: `cd backend && go test ./...`; `cd frontend/roma-web && npm run build`.

### CHK-003 - Migraciones y arranque local
Estado: Completado
Objetivo: validar DB limpia, migraciones oficiales, seed minimo, backend y endpoint basico.
Resultado: Compose levanta Postgres, corre migrator, deja backend healthy y permite login seed.
Validado: `docker compose -p roma_chk003 -f database/docker-compose.yml up -d backend`; `/healthz`, `/readyz`, `/auth/login`.

### CHK-004 - Roles, ownership y guards minimos
Estado: Completado
Objetivo: persistir `users.role` y proteger recursos criticos desde backend.
Resultado: ejercicios/programas/asignaciones/sesiones/historial tienen guards minimos; frontend oculta acciones segun rol.
Validado: `GOCACHE=/tmp/roma-go-cache go test ./...`; `npm run build`; `docker compose -p roma_chk004 -f database/docker-compose.yml up -d backend`.
Validado manual: coach login OK, disciple login OK, disciple crear ejercicio 403, coach crear ejercicio 201, coach ajeno overview 403.

### CHK-005 - Tests automatizados de permisos
Estado: Completado
Objetivo: cubrir roles, `/me`, ejercicios y guards criticos de ownership/acceso.
Resultado: tests Go agregados para auth, permisos de ejercicios y helpers de guards sobre coach-disciple, programas, assignments, sesiones, sets y consistencia dia-prescripcion.
Validado: `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`; `cd frontend/roma-web && npm run build`.
Pendiente: agregar pruebas E2E/API con DB real para assignments, sesiones e historial cuando el setup de endpoints quede estable.

### CHK-006 - Tests E2E/API backend con DB limpia
Estado: Completado
Objetivo: validar permisos reales con router, handlers, migraciones y seed controlado.
Resultado: E2E con `ROMA_E2E_DB_URL` cubre auth, exercises, programs, assignments, sessions, sets e history; test normal salta si no hay DB.
Validado: `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL=postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable GOCACHE=/tmp/roma-go-cache go test ./... -run E2E`; `npm run build`.
Pendiente: endpoint de editar set no existe; solo queda validado delete ajeno. Mantener DB E2E local, nunca productiva.

### CHK-007 - Self-training minimo
Estado: Completado
Objetivo: disciple sin maestro crea rutina propia, agrega estructura, se auto-asigna y ejecuta sesion.
Resultado: `programs.kind` distingue `coach_program`/`self_training`; guards permiten solo mutacion propia; E2E cubre rutina, self-assignment, sesion, set e historial.
Validado: migraciones 0001-0005 en `roma_e2e`; `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL=postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`; `npm run build`.
Pendiente: manejar multiples self-assignments activos y mejorar UI de inicio/activacion.

### CHK-008 - Endurecer activacion self-training
Estado: Completado
Objetivo: dejar solo una self-assignment activa por disciple sin tocar assignments de coach.
Resultado: migracion `0006` reemplaza indice global por indice parcial de self-assignments; activacion reactiva/crea en transaccion y desactiva solo self-training previo.
Validado: migraciones 0001-0006 en `roma_e2e`; `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL=postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`; `npm run build`.
Pendiente: impedir inicio de sesiones sobre assignments inactivos si se decide aplicar esa regla globalmente.

### CHK-008b - Bugfix dias self-training
Estado: Completado
Objetivo: listar dias agregados al iniciar sesion desde self-training activo.
Resultado: ruta `/api/assignments/:id/days` corregida; scan de `text[]` arreglado; errores inesperados ya no devuelven lista vacia silenciosa.
Validado: `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL=postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`; `npm run build`.

### CHK-009 - Bloqueo de sesiones en assignments inactivos
Estado: Completado.
Objetivo: impedir nuevas sesiones sobre assignments `is_active=false`.
Resultado: `POST /api/sessions` responde `409 assignment_inactive` tras validar ownership; frontend muestra mensaje claro.
Validado: migraciones 0001-0006 en `roma_e2e`; `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL=postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`; `npm run build`.

### CHK-010 - Pulido UX Entrenar
Estado: Completado.
Objetivo: clarificar rutina activa, dias disponibles, ejercicios y accion de continuar/iniciar sesion.
Resultado: Entrenar muestra sesion activa, rutina activa, dias visibles y detalle de ejercicios sin depender de modal.
Validado: `npm run build`.

### CHK-011 - Editar rutina propia
Estado: Completado.
Objetivo: permitir editar titulo y descripcion/notas de rutinas propias y programas propios.
Resultado: Mis rutinas/Programas agrega accion Editar usando `PUT /api/programs/:id`; payload limitado a `title`/`notes`.
Validado: `npm run build`; `GOCACHE=/tmp/roma-go-cache go test ./...`.

### CHK-012 - Descripcion Entrenar
Estado: Completado.
Objetivo: mostrar rutina activa y dias con datos reales y legibles.
Resultado: Entrenar usa `assigned_by`/`disciple_id`, `program.kind`, notas, fecha formateada, conteo y resumen de ejercicios.
Validado: `npm run build`.

### BUGFIX - Cambio de self-training activo
Estado: Completado.
Objetivo: permitir activar rutina B sin violar `assignments_check`.
Resultado: desactivar self-assignments anteriores solo cambia `is_active=false`; E2E cubre activaciones futuras/mismo dia.
Validado: `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL=postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`; `npm run build`.

### CHK-013 - Cierre de sesion y resumen
Estado: Completado.
Objetivo: cerrar sesion, bloquear mutaciones posteriores y mostrar resumen basico.
Resultado: `PATCH /api/sessions/:id` cierra con `status=closed`; sets/cardio/delete quedan bloqueados con `409`; frontend muestra resumen.
Validado: `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL=postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`; `npm run build`.

### CHK-014 - Historial de sesiones cerradas
Estado: Completado.
Objetivo: listar sesiones realizadas con metadata y link a resumen.
Resultado: `/api/history` y `/api/sessions/:id` incluyen programa/semana/dia; Historial lista sesiones con sets, ejercicios y volumen.
Validado: `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL=postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`; `npm run build`.

### CHK-015 - Filtros backend de historial
Estado: Completado.
Objetivo: filtrar sesiones desde backend por `from`, `to`, `status` y `program_id`.
Resultado: `/api/history?group=session` y `/api/history/disciples/:id/sessions` validan filtros y los aplican en SQL; Historial envia filtros reales sin recortar localmente los ultimos 50 items.
Validado: `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL=postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`; `npm run build`.
Pendiente: selector de rutina depende de programas visibles por `/api/programs`; revisar si debe incluir historicos asignados por coach ya no listados.

### CHK-016 - Check-ins basicos
Estado: Completado.
Objetivo: registrar fecha, peso opcional y notas de seguimiento del disciple.
Resultado: `checkins.checked_at` agregado; disciple crea/lista/ve propios; coach vinculado lista/ve; terceros bloqueados; UI disciple simple.
Validado: migraciones 0001-0007 en DB temporal limpia; `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL=postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`; `npm run build`.
Pendiente: UI coach para check-ins en detalle de disciple; edicion/borrado quedan fuera.

### CHK-017 - Check-ins en vista coach/disciple
Estado: Completado.
Objetivo: permitir al coach revisar check-ins recientes del disciple vinculado.
Resultado: detalle de disciple muestra ultimos 5 check-ins con fecha, peso y notas usando `GET /api/coach/disciples/:id/checkins`.
Validado: `npm run build`.
Pendiente: filtros por fecha y detalle historico de check-ins si crece el volumen.

### CHK-018 - Vista coach de progreso del disciple
Estado: Completado.
Objetivo: compactar el detalle del disciple para revisar estado, sesiones recientes y seguimiento.
Resultado: detalle muestra resumen, rutina activa, ultimos entrenamientos, check-ins y accesos a historial/asignaciones; Historial acepta `disciple_id`.
Validado: `npm run build`.
Pendiente: mejorar vista completa de historial coach con nombre del disciple y opciones de contexto.

### ROMA-019 / CHK-019 - Documentacion operativa y flujo por Work Items
Estado: Done
Objetivo: crear una forma segura y simple de desarrollo orientado a agentes.
Resultado: `ROADMAP.md` creado; `AGENTS.md` actualizado para operar por Work Items; creada carpeta `docs/work-items/` con README, template y ROMA-019; creados stubs minimos ROMA-020 a ROMA-025.
Validado: `git status --short`.
Proximo sugerido: ROMA-020.

### ROMA-020 / CHK-020 - Pulir historial coach con contexto del alumno
Estado: Done
Objetivo: mostrar claramente el disciple y los filtros activos cuando el coach revisa historial.
Resultado: `/history?disciple_id=...` muestra nombre, email, ID del disciple, enlace de vuelta y resumen visible de filtros/resultados; los accesos desde detalle de disciple pasan contexto de navegacion.
Validado: `cd frontend/roma-web && npm run build`.
Proximo sugerido: ROMA-021.

### ROMA-021 / CHK-021 - Filtros/paginacion de check-ins
Estado: Done
Objetivo: agregar filtros simples por fecha y paginacion basica a check-ins.
Resultado: endpoints de check-ins propios y coach aceptan `from`, `to`, `limit`, `offset`; UI disciple permite filtrar por fecha, seleccionar pagina y navegar resultados; E2E cubre filtros con ownership.
Validado: `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL=postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`; `npm run build`.
Proximo sugerido: ROMA-022.

### Nota de alcance ROMA-022
ROMA-022 fue acotado a edicion de check-ins propios. Borrado/archivado queda fuera de alcance por ahora y requerira decision/WI aparte si se prioriza.

### ROMA-022 / CHK-022 - Editar check-ins propios
Estado: Done
Objetivo: permitir que el disciple corrija sus propios check-ins sin habilitar borrado.
Resultado: `PATCH /api/checkins/:id` edita `checked_at`, `weight_kg` y `notes` solo para el disciple dueño; UI de Check-ins permite editar inline; coach/otros disciples quedan bloqueados.
Validado: `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL=postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`; `npm run build`.
Proximo sugerido: ROMA-023.

### ROMA-023 / CHK-023 - Biblioteca de rutinas / plantillas Roma, diseno primero
Estado: Done
Objetivo: definir reglas de producto y flujo antes de implementar plantillas.
Resultado: creado `docs/product/ROMA_TEMPLATE_LIBRARY_DESIGN.md` con tipos de plantillas, ownership, visibilidad, relacion con programas existentes y fuera de alcance.
Validado: `git status --short`.
Proximo sugerido: ROMA-024.

### ROMA-024 / CHK-024 - Demo FST-7 seed opcional
Estado: Done
Objetivo: cargar una pauta demo FST-7 como seed opcional para probar Roma con datos reales sin mezclar datos demo en migraciones productivas.
Resultado: `0008_demo_fst7_seed` fue retirado de migraciones oficiales; el seed quedo en `database/seed/demo_fst7.up.sql` y `database/seed/demo_fst7.down.sql`, con cuentas demo `@roma.demo`, vinculo coach-disciple, programa privado `[DEMO] FST-7 Hipertrofia - 5 dias`, 5 dias, 16 prescripciones, 5 cierres FST-7 y assignment activo para ejecucion local.
Validado: `GOCACHE=/tmp/roma-go-cache go test ./...`; migraciones oficiales `0001` a `0007` aplicadas en Postgres temporal limpio sin seed demo; `database/seed/demo_fst7.up.sql` aplicado manualmente; smoke SQL creo una sesion y un set sobre el assignment demo.
Proximo sugerido: ROMA-025.

### ROMA-025 / CHK-025 - Hardening final MVP
Estado: Done
Objetivo: revisar seguridad, validaciones, UX critica, README y flujo demo antes de considerar MVP.
Resultado: rutas `/api/coach/assignments/:id/calendar` y `/api/coach/assignments/:id/activate` ahora requieren rol `coach`; se agrego E2E para bloquear disciples en esas rutas; se actualizo README, estado actual, testing y se creo `docs/engineering/ROMA_MVP_HARDENING.md` con riesgos pendientes.
Validado: `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL='postgres://roma:roma@localhost:55432/roma_e2e?sslmode=disable' GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1` contra Postgres temporal; `npm run build`.
Proximo sugerido: ROMA-026.

### ROMA-026 / CHK-026 - Pulido UI general MVP
Estado: Done
Objetivo: mejorar claridad y consistencia UI del MVP sin redisenar toda la app ni agregar features.
Resultado: se agrego `QueryState` compartido para estados vacios/loading/error y se aplico en Entrenar, Historial, Check-ins, Dashboard coach y Asignaciones; se ajusto una inconsistencia visual menor en la tabla de asignaciones.
Validado: `cd frontend/roma-web && npm run build`; `git diff --check`.
Proximo sugerido: ROMA-027.

### Flujo agentes / ROMA-027 preparado
Estado: Done
Objetivo: ajustar reglas de trabajo con cambios sin commit, cerrar el roadmap inicial y crear el WI de preparacion demo local.
Resultado: `AGENTS.md` aclara que cambios sin commit no bloquean automaticamente; `ROADMAP.md` agrega proximos despues del MVP inicial; `docs/work-items/ROMA-027.md` creado para documentar demo local sin tocar backend/frontend.
Validado: `git status --short`.
Proximo sugerido: ROMA-027.

### ROMA-027 / CHK-027 - Preparacion demo local
Estado: Done
Objetivo: dejar una demo local facil de levantar y probar.
Resultado: se creo `docs/engineering/LOCAL_DEMO.md` con pasos para levantar backend/frontend, aplicar o retirar seed FST-7 opcional, usuarios demo, flujo sugerido para coach/disciple y checklist manual; se enlazo desde README y setup. No se tocaron backend, frontend ni migraciones.
Validado: `git status --short`; `git diff --check`.
Proximo sugerido: ROMA-028.

### ROMA-028 / CHK-028 - Demo local: prueba manual y lista de ajustes menores
Estado: Done
Objetivo: definir una prueba manual completa del MVP local usando `LOCAL_DEMO.md` y el seed FST-7 opcional, registrando bugs, fricciones y ajustes menores sin implementar cambios todavia.
Resultado: se registraron hallazgos de demo: coach tambien necesita self-training; "Mis rutinas" debe distinguir rutinas propias/importadas de rutinas asignadas por maestro; las rutinas asignadas no deben eliminarse desde la seccion personal; self-training independiente funciona correctamente; falta biblioteca de rutinas visible para Discipulos e Independientes; importar desde biblioteca debe copiar a "Mis rutinas" sin borrar la plantilla global al eliminar la copia.
Validado: `git status --short`.
Proximo sugerido: ROMA-029.

### Proximos WIs sugeridos por ROMA-028
- ROMA-029 - Registrar hallazgos demo y ajustar roadmap.
- ROMA-030 - Coach puede usar entrenamiento personal.
- ROMA-031 - Mis rutinas: asignadas, propias e importadas.
- ROMA-032 - Diseno biblioteca de rutinas.
- ROMA-033 - Implementar biblioteca de rutinas MVP.

### ROMA-029 / CHK-029 - Registrar hallazgos demo y ajustar roadmap
Estado: Done
Objetivo: documentar los hallazgos de demo y ordenar la siguiente fase despues del MVP inicial.
Resultado: se consolido la fase post MVP inicial: ROMA-030 coach puede usar entrenamiento personal; ROMA-031 ordena "Mis rutinas" entre asignadas, propias e importadas; ROMA-032 disena biblioteca de rutinas; ROMA-033 implementa el MVP de biblioteca. No se cambiaron backend, frontend ni migraciones.
Validado: `git status --short`; `git diff --check`.
Proximo sugerido: ROMA-030.

### Fase post MVP preparada - ROMA-030 a ROMA-033
Estado: Planned
Objetivo: dejar creados los Work Items ejecutables de la fase posterior al MVP inicial.
Resultado: creados `docs/work-items/ROMA-030.md`, `docs/work-items/ROMA-031.md`, `docs/work-items/ROMA-032.md` y `docs/work-items/ROMA-033.md` usando la plantilla de Work Items. No se cambiaron backend, frontend ni migraciones.
Validado: `git status --short`; `git diff --check`.
Proximo sugerido: ROMA-030.

### ROMA-030 / CHK-030 - Coach puede usar entrenamiento personal
Estado: Done
Objetivo: permitir que un coach use Entrenar, Mis rutinas, Historial personal y Check-ins personales, sin dejar de ser coach.
Resultado: backend permite self-training y check-ins personales para coach; UI muestra Entrenar, Mis rutinas, Historial y Check-ins para coach; Mis rutinas del coach permite crear rutina personal o programa para alumnos; E2E cubre check-in personal y flujo self-training personal del coach.
Validado: `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL='postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable' GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`; `npm run build`.
Proximo sugerido: ROMA-031.

### ROMA-031 / CHK-031 - Mis rutinas: asignadas, propias e importadas
Estado: Done
Objetivo: ordenar Mis rutinas para mostrar rutinas asignadas por maestro, rutinas propias y rutinas importadas, con permisos claros.
Resultado: `/api/programs` incluye programas propios y rutinas asignadas al usuario con metadata de origen/permisos; Mis rutinas agrupa asignadas por maestro, propias/programas e importadas; rutinas asignadas son visibles pero no editables/eliminables desde la seccion personal; E2E cubre visibilidad y bloqueo de eliminacion.
Validado: `GOCACHE=/tmp/roma-go-cache go test ./...`; `ROMA_E2E_DB_URL='postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable' GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`; `npm run build`.
Proximo sugerido: ROMA-033.

### ROMA-032 / CHK-032 - Diseno biblioteca de rutinas
Estado: Done
Objetivo: cerrar reglas de producto y modelo tecnico minimo antes de implementar la biblioteca v1.
Resultado: se definieron plantillas globales `draft`/`published`, permisos de coach/admin, visibilidad autenticada, importacion como copia profunda independiente y separacion entre plantillas, programas y assignments. ROMA-033 queda acotado a listado, previsualizacion, publicacion controlada, importacion y pruebas de ownership. No se tocaron backend, frontend ni migraciones.
Validado: `git diff --check`.
Proximo sugerido: ROMA-033.

## Pendientes importantes
- Consolidar/eliminar `master_disciple` cuando sea seguro.
- Ampliar E2E cuando aparezcan endpoints de editar sets/check-ins.
- Revisar endpoints no cubiertos: invitaciones, check-ins futuros.


### ROMA-033 / CHK-033 - Biblioteca de rutinas MVP
Estado: Done
Objetivo: permitir descubrir, publicar e importar rutinas como copias personales independientes.
Resultado: se agregaron tablas de plantillas y estructura propia mediante migracion `0008`, endpoints autenticados para listar/previsualizar/crear/publicar/importar, copia transaccional a `programs` y UI minima de Biblioteca. Coach puede crear/publicar y disciple puede importar plantillas publicadas.
Validado: `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`; E2E/API con `ROMA_E2E_DB_URL` y `-run E2E -count=1`; `cd frontend/roma-web && npm run build`; `git diff --check`.
Proximo sugerido: ROMA-034.


### ROMA-034 / CHK-034 - Freeze v1 + criterios mobile
Estado: Done
Objetivo: congelar el alcance de Roma v1 y definir criterios verificables de aceptación mobile.
Resultado: se fijaron viewports 360/390/430 px y desktop, flujos críticos de disciple/independent y coach, criterios de navegación, formularios, cards, overflow y acciones principales. Desde este checkpoint solo entran bugs críticos o bloqueantes; app nativa, PWA avanzada y nuevas features quedan fuera. No se modificaron backend ni frontend.
Validado: `git status --short`; `git diff --check`.
Proximo sugerido: ROMA-035.


### WIs finales de Roma v1 preparados
Estado: Planned
Resultado: creados ROMA-035 Responsive base + navegación, ROMA-036 Responsive flujos disciple/independent, ROMA-037 Responsive flujos coach y ROMA-038 Hardening, regresión y cierre Roma v1. No se implementó código.
Proximo sugerido: ROMA-035.


### ROMA-035 / CHK-035 - Responsive base + navegacion
Estado: Done
Objetivo: hacer usable la estructura principal en mobile preservando desktop.
Resultado: se agregó navegación móvil con overlay, header y contenedor responsive, controles táctiles para theme/logout y prevención de overflow horizontal accidental. El sidebar desktop se mantiene.
Validado: `cd frontend/roma-web && npm run build`; `git diff --check`.
Proximo sugerido: ROMA-036.
