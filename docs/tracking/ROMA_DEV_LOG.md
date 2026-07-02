# ROMA Dev Log

## Estado actual
- Fase actual: Fase 1 self-training endurecido.
- Objetivo actual: una sola rutina propia activa por disciple sin afectar assignments de coach.
- Ultimo checkpoint completado: CHK-008 - Endurecer activacion de self-training.
- Proximo checkpoint: CHK-009 - cerrar gaps de ejecucion y UX de sesiones.

## Decisiones activas
- [2026-06-24] Decision: no partir desde cero; rescatar repo con estabilizacion previa.
- [2026-06-24] Decision: alinear migraciones oficiales con backend actual creando `coach_links` y `program_versions` en Fase 0.
- [2026-06-24] Decision: activar `migrator` en Compose para DB limpia reproducible antes de backend.
- [2026-06-25] Decision: `coach_links` queda como fuente operativa; `master_disciple` queda legacy/compatibilidad.
- [2026-06-29] Decision: self-training usa `programs.kind='self_training'` y self-assignment propio; no usa self-link en `coach_links`.

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

## Pendientes importantes
- Consolidar/eliminar `master_disciple` cuando sea seguro.
- Ampliar E2E cuando aparezcan endpoints de editar sets/check-ins.
- Revisar endpoints no cubiertos: invitaciones, check-ins futuros.
