# ROMA Keep / Drop List

## Tabla de decision

| Area | Categoria | Decision | Motivo |
|---|---|---:|---|
| Estructura monolito modular backend | Mantener, pero refactorizar | Mantener | La separacion `domain/repository/service/transport` es util, pero hay DTOs duplicados y SQL/GORM mezclados sin criterio claro. |
| Go + Gin + GORM | Mantener tal cual | Mantener | Stack suficiente para MVP. No hay razon para cambiar framework. |
| React + Vite + TypeScript | Mantener tal cual | Mantener | Stack correcto. El problema es consistencia y tipos, no tecnologia. |
| PostgreSQL | Mantener tal cual | Mantener | Encaja con el dominio relacional. |
| `users` | Mantener, pero refactorizar | Mantener | Base valida, pero falta rol/politica persistente. |
| Roles derivados desde links/programas | Reescribir desde cero | Reescribir | No expresa permisos reales y convierte acciones en identidad. |
| `master_disciple` | Eliminar o migrar | Drop/Merge | La app usa `coach_links`; no pueden convivir dos fuentes de verdad. |
| `coach_links` | Mantener, pero refactorizar | Mantener | Encaja mejor con invitaciones/status, pero debe estar en migraciones oficiales y tener nombre definitivo. |
| Self-link para auto-entrenamiento | Reutilizar parcialmente | Refactor | Sirve como demo, pero conceptualmente no basta para "discipulo sin maestro". |
| Ejercicios | Mantener, pero refactorizar | Mantener | CRUD y filtros son utiles; faltan permisos para impedir creacion por discipulos. |
| Metodos/protocolos | Reutilizar parcialmente | Mantener limitado | Tabla existe y FST-7 seed; falta API/UI. Para MVP conviene lectura/seed, no CRUD amplio. |
| Programas/semanas/dias/prescripciones | Mantener, pero refactorizar | Mantener | Es el nucleo mas rescatable; requiere ownership, constraints y contratos estables. |
| Versionado de programas | Reutilizar parcialmente | Decidir | Hay servicio para versiones, pero falta tabla oficial. O se completa o se elimina del MVP. |
| Asignaciones | Mantener, pero refactorizar | Mantener | Concepto correcto; necesita ownership, activacion segura y consistencia de `is_active`. |
| Sesiones | Mantener, pero refactorizar | Mantener | Flujo base existe; faltan validaciones de assignment/day/prescription y ownership en patch/delete. |
| Set logs | Mantener, pero refactorizar | Mantener | Modelo correcto; borrar/editar debe autorizarse. |
| Cardio | Reutilizar parcialmente | Mantener | Tabla flexible; API/frontend solo cubren parte realizada. Falta prescripcion y contrato de payload. |
| Check-ins | Reutilizar parcialmente | Mantener despues | Tabla existe; sin API/UI. Implementar en fase dashboard/progreso. |
| Historial/pivots | Mantener, pero refactorizar | Mantener | Util para progreso; endpoints requieren guards por self/coach. |
| Invitaciones | Reutilizar parcialmente | Consolidar | Hay `invitations` e `invite_codes`; elegir una. |
| Auth JWT | Mantener, pero refactorizar | Mantener | Base util; validar secreto, unificar config, revisar refresh. |
| Middleware de auth | Mantener, pero refactorizar | Mantener | Funciona como base, pero faltan politicas por recurso. |
| Guards de coach | Reutilizar parcialmente | Refactor | Hay piezas utiles, uso parcial e inconsistencia de tabla. |
| Rate limit auth | Mantener tal cual | Mantener | Correcto para MVP. |
| Swagger generado | Reutilizar parcialmente | Regenerar | Actualmente rompe build. Mantener solo si se automatiza. |
| Frontend login | Mantener, pero refactorizar | Mantener | Funciona como base; falta cargar usuario/rol y manejar errores consistentes. |
| Frontend dashboard maestro | Mantener, pero refactorizar | Mantener | Alineado con producto; depende de permisos correctos. |
| Frontend detalle discipulo | Mantener, pero refactorizar | Mantener | Muy cerca del flujo coach revisa/alumno ejecuta, pero mezcla responsabilidades. |
| Frontend ejercicios | Mantener, pero refactorizar | Mantener | Pantalla util para maestro; ocultar/bloquear para discipulo. |
| Frontend programas | Mantener, pero refactorizar | Mantener | Constructor rescatable; necesita simplificacion y build verde. |
| Frontend sesiones | Mantener, pero refactorizar | Mantener | Flujo principal de ejecucion esta presente. |
| Frontend historial | Reutilizar parcialmente | Mantener | Basico; requiere contratos estables y permisos. |
| Sidebar actual | Reescribir parcialmente | Refactor fuerte | Muestra todo a todos; debe ser por rol/estado. |
| ProtectedRoute actual | Mantener, pero refactorizar | Refactor | Solo auth; agregar permisos/roles. |
| README frontend template | Eliminar | Drop | No documenta Roma System. |
| README infra vacio | Eliminar o reemplazar | Drop/Reemplazar | No aporta. |
| Carpetas `infra/helm`, `infra/k8s`, `infra/terraform` vacias | Eliminar | Drop | No sirven al MVP si estan vacias. |
| Seeds `database/seed` | Reutilizar parcialmente | Consolidar | Tienen datos utiles, pero esquema distinto al de migraciones. |
| Migraciones oficiales | Mantener, pero refactorizar | Mantener | Buen punto de partida; corregir coach_links, program_versions, title, constraints. |
| Docker Compose | Mantener, pero refactorizar | Mantener | Base util; migrator comentado y falta frontend. |

## Mantener tal cual

- Stack general: Go, Gin, React, Vite, PostgreSQL.
- Layout modular del repo.
- Rate limit en auth como idea.
- Uso de TanStack Query/Axios/Zustand.

## Mantener, pero refactorizar

- Auth JWT.
- Programas/semanas/dias/prescripciones.
- Ejercicios.
- Asignaciones.
- Sesiones y sets.
- Historial.
- Dashboard maestro.
- Constructor frontend de programas.
- Vista de sesiones.
- Docker Compose.

## Reutilizar parcialmente

- Seeds demo.
- FST-7 seed.
- Metodos/protocolos.
- Check-ins.
- Invitaciones.
- Guards de coach.
- Swagger.
- Cardio.
- Self-link de Ikki como idea de demo, no como modelo final.

## Eliminar

- README de Vite.
- README infra vacio.
- Carpetas infra vacias si no se implementaran ahora.
- Seeds que creen esquemas divergentes.
- Una de las dos tablas de invitaciones.
- Una de las dos relaciones maestro-discipulo.
- Rutas duplicadas `/programs/programs/...`.
- `console.log` en servicios frontend.
- `visibility public/unlisted` si no hay uso real en MVP privado.

## Reescribir desde cero

- Modelo de roles/permisos.
- Matriz de autorizacion backend.
- Navegacion frontend por rol.
- Flujo explicito de discipulo sin maestro.
- Contrato oficial de inicializacion DB.

## Lista priorizada de drops

1. Eliminar la doble fuente `master_disciple` vs `coach_links`.
2. Eliminar la doble fuente `invitations` vs `invite_codes`.
3. Eliminar seeds antiguos que no pasen por migraciones oficiales.
4. Eliminar docs template/vacias.
5. Eliminar rutas duplicadas y endpoints muertos.
6. Eliminar o postergar `public/unlisted` hasta que exista necesidad real.
7. Eliminar referencias a `roleHint` si no se persistira.

## Lista priorizada de rescate

1. Recuperar build backend/frontend.
2. Corregir migraciones oficiales.
3. Blindar permisos y ownership.
4. Rescatar builder de programas.
5. Rescatar ejecucion de sesiones y sets.
6. Rescatar historial/pivot.
7. Cerrar dashboard maestro.
8. Agregar self-training con permisos correctos.
