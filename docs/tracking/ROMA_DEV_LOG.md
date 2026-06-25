# ROMA Dev Log

## Estado actual
- Fase actual: Fase 1 base minima aplicada.
- Objetivo actual: roles persistentes y guards backend enforceables.
- Ultimo checkpoint completado: CHK-004 - Roles, ownership y guards minimos.
- Proximo checkpoint: CHK-005 - ampliar cobertura de permisos y pruebas automatizadas.

## Decisiones activas
- [2026-06-24] Decision: no partir desde cero; rescatar repo con estabilizacion previa.
- [2026-06-24] Decision: alinear migraciones oficiales con backend actual creando `coach_links` y `program_versions` en Fase 0.
- [2026-06-24] Decision: activar `migrator` en Compose para DB limpia reproducible antes de backend.
- [2026-06-25] Decision: `coach_links` queda como fuente operativa; `master_disciple` queda legacy/compatibilidad.

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

## Pendientes importantes
- Consolidar/eliminar `master_disciple` cuando sea seguro.
- Agregar tests automatizados de permisos.
- Revisar endpoints no cubiertos: invitaciones, check-ins futuros, self-training futuro.
