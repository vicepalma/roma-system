# ROMA Repo Audit

Fecha de auditoria: 2026-06-24

## Resumen ejecutivo

El repositorio contiene una base rescatable para Roma System, pero no esta en estado sano para continuar sin una fase de rescate tecnica. La arquitectura general coincide con la propuesta original: backend Go, frontend React/TypeScript, PostgreSQL, REST, JWT y Docker para desarrollo. El dominio de entrenamiento esta parcialmente modelado y varias pantallas consumen API real.

El problema principal no es ausencia de codigo, sino inconsistencia entre capas:

- La migracion oficial crea `master_disciple`, pero el backend consulta `coach_links`.
- El backend deriva roles, pero no tiene un modelo persistente de roles ni autorizacion consistente por rol.
- Las rutas de ejercicios, programas, sesiones, historial y asignaciones tienen proteccion desigual de ownership.
- El frontend muestra rutas de maestro a cualquier usuario autenticado.
- Backend y frontend no compilan actualmente.
- Hay seeds y esquemas duplicados que representan versiones distintas del producto.

Recomendacion corta: continuar sobre este repo solo si primero se hace una fase de estabilizacion. No conviene partir desde cero, pero tampoco conviene seguir agregando funcionalidades encima del estado actual.

## Reconocimiento del proyecto

### Lenguajes y frameworks

Backend:

- Go 1.24.
- Gin como framework HTTP.
- GORM para persistencia.
- PostgreSQL.
- JWT con `github.com/golang-jwt/jwt/v5`.
- Hash de password con Argon2id.
- Swagger generado con swaggo.

Frontend:

- React 19.
- TypeScript.
- Vite.
- React Router 7.
- TanStack Query.
- Axios.
- Zustand.
- Tailwind CSS.
- Recharts.
- React Hook Form y Zod instalados, con uso parcial.

Base de datos:

- PostgreSQL 15 en Docker Compose.
- Migraciones SQL en `database/backend/migrations`.
- Seeds alternativos en `database/seed`.

Infra:

- `database/docker-compose.yml` levanta Postgres y backend.
- Migrator esta comentado.
- No hay frontend en compose.
- Carpetas `infra/helm`, `infra/k8s`, `infra/terraform` existen, pero no contienen implementacion util visible.

Tests:

- No se encontraron tests Go ni frontend.
- `go test ./...` falla por Swagger generado incompatible.
- `npm run build` falla por errores TypeScript.

Documentacion:

- No hay README raiz.
- `frontend/roma-web/README.md` es el template de Vite.
- `infra/README.md` esta vacio.
- Swagger existe en `backend/docs`, pero el codigo generado no compila con las dependencias actuales.

Variables de entorno:

- Existen `backend/.env`, `backend/.env.docker`, `frontend/roma-web/.env` y `database/.env.example`.
- El backend real lee `DB_URL`, `PORT`, `ENV`, `JWT_SECRET`, `ACCESS_TTL_MIN`, `REFRESH_TTL_H`, `DEFAULT_TZ`.
- `backend/internal/config/config.go` define `DATABASE_DSN`, pero el servidor actual usa `DB_URL`. Hay dos convenciones activas.

## Estructura de carpetas

```txt
backend/
  cmd/server/main.go
  internal/config
  internal/domain
  internal/middleware
  internal/repository
  internal/security
  internal/service
  internal/transport/http
  docs/swagger.*

frontend/roma-web/
  src/components
  src/lib
  src/pages
  src/router
  src/services
  src/store
  src/types

database/
  backend/migrations
  seed
  scripts

infra/
  helm
  k8s
  terraform
```

La estructura modular es razonable. El mayor problema es que las fronteras no estan aplicadas de forma consistente: hay DTOs duplicados entre `domain`, `repository` y `service`, SQL crudo junto a GORM, y contratos frontend que normalizan respuestas variables porque el backend no entrega formas estables.

## Estado de build

Comandos ejecutados:

```bash
cd backend && go test ./...
cd frontend/roma-web && npm run build
```

Resultado backend:

- Falla en `backend/docs/docs.go`.
- Error: `unknown field LeftDelim/RightDelim in swag.Spec`.
- No hay tests; los paquetes restantes reportan `[no test files]`.

Resultado frontend:

- Falla TypeScript.
- `Exercises.tsx`: mutacion espera `void | Promise<void>` pero recibe `Promise<Exercise>`.
- `ProgramDetail.tsx`: payload usa `index` en vez de `week_index`.
- `ProgramDetail.tsx`: falta `method_id` requerido por el tipo de `addPrescription`.

Conclusion: antes de desarrollar producto, hay que recuperar compilacion reproducible.

## Comparacion contra producto actual

### Roles y permisos

Estado actual:

- Existe tabla `users`.
- No existe columna persistente `role`.
- El rol se deriva: coach si tiene `coach_links` aceptados o programas propios; si no, disciple.
- El signup recibe `role`, pero el comentario indica que no se persiste.
- Hay relacion maestro-discipulo en dos nombres incompatibles:
  - migracion oficial: `master_disciple`;
  - backend y seed nuevo: `coach_links`.
- Auto-entrenamiento aparece como self-link en seeds (`coach_id = disciple_id`) para Ikki.

Brecha:

- No hay regla robusta para "maestro puede construir para otros" vs "discipulo puede ejecutar y, si no tiene maestro, crear rutinas propias".
- Cualquier usuario autenticado puede crear ejercicios globales.
- Cualquier usuario autenticado puede crear programas.
- No hay distincion clara entre programa de maestro, rutina propia simple y catalogo global.
- La proteccion frontend no distingue roles; solo verifica login.

Riesgo:

- Un discipulo puede acceder a pantallas de maestro y llamar endpoints de maestro si el backend no lo bloquea.
- Un usuario sin maestro puede terminar convertido en "coach" solo por crear un programa.
- La derivacion de rol desde links/programas no expresa permisos del producto.

### Dominio de entrenamiento

Cubierto en esquema:

- ejercicios;
- metodos/protocolos;
- programas;
- semanas;
- dias;
- prescripciones;
- asignaciones;
- sesiones;
- set logs;
- cardio;
- check-ins;
- user flags;
- invitaciones.

Cubierto en backend/API:

- ejercicios CRUD;
- programas CRUD parcial;
- semanas y dias;
- prescripciones;
- asignaciones;
- sesiones;
- sets;
- cardio de sesion;
- historial y resumen/pivot;
- coach dashboard parcial;
- invitaciones parcial.

Incompleto o riesgoso:

- `methods` existe en DB, pero no hay API/pantalla de gestion de metodos.
- FST-7 existe solo como seed minimal, no como caso demo completo.
- Check-ins existen solo en DB; no hay handlers/pantallas.
- Cardio prescrito existe por tabla, pero el frontend solo muestra/crea cardio realizado de forma parcial.
- No hay modelo explicito para dias elegidos por discipulo sin maestro.
- No hay separacion entre biblioteca global y recursos creados por maestro.
- Versionado de programas esta incompleto: el servicio usa `program_versions`, pero la migracion oficial no crea esa tabla.
- `program_days.title` es `NOT NULL` en migracion oficial, pero seeds insertan dias sin `title` en varias rutas.

### Seguridad

Existe:

- Login con JWT access/refresh.
- Middleware `AuthRequired`.
- Refresh token.
- Password hashing Argon2id para signup/login actual.
- Rate limit en `/auth`.
- Guards parciales para coach-disciple.

Problemas:

- `JWT_SECRET` vacio no se valida en `GenerateTokens`; si falta, firmaria con secreto vacio.
- CORS esta hardcodeado a localhost y no usa `CORS_ORIGINS`.
- Ownership de programas no se valida en muchas operaciones: obtener, actualizar, borrar, agregar semanas/dias/prescripciones.
- `sessions/:id` GET valida que la sesion sea del usuario; `GET /sets` permite coach; pero `DELETE /sets/:setId` no valida ownership.
- `PATCH /sessions/:id` no valida ownership.
- `POST /sessions` acepta `assignment_id` y `day_id` sin validar que pertenezcan al usuario y al mismo programa.
- `POST /exercises`, `PUT /exercises`, `DELETE /exercises` solo requieren auth.
- `GET /history?disciple_id=` no tiene guard visible por coach o self.
- `PATCH /coach/assignments/:id` y calendario no validan ownership del coach sobre esa asignacion.
- `CoachGuard` y `RequireCoachOf` existen, pero su uso es parcial.

### Frontend

Pantallas existentes:

- Login.
- Dashboard maestro con discipulos.
- Detalle de discipulo con hoy, sets, overview.
- Asignaciones.
- Ejercicios.
- Programas.
- Detalle de programa.
- Inicio de sesiones.
- Vista de sesion.
- Historial pivot.

Lo rescatable:

- Hay integracion real con API usando Axios y TanStack Query.
- La vista de sesion y el flujo "programa activo -> elegir dia -> registrar sets" estan cerca del objetivo.
- Program builder tiene piezas utiles.
- Dashboard maestro ya apunta a discipulos, adherencia y ultimos datos.

Problemas:

- `ProtectedRoute` solo valida autenticacion.
- Sidebar muestra rutas de maestro a todos.
- No existe vista clara para discipulo sin maestro que cree rutinas simples usando ejercicios existentes.
- No existe gating para impedir crear ejercicios desde cuenta discipulo.
- Hay pagina `/profile` enlazada pero no registrada.
- Build TypeScript roto.
- Hay normalizadores para tolerar respuestas inconsistentes del backend.
- README frontend sigue siendo template de Vite.

## Desvios de producto y basura tecnica

No se detecto foco fuerte en comunidad, red social o simpatizantes. El desvio principal es otro: una mezcla de versiones del mismo dominio.

Cosas a sacar o corregir:

- Duplicidad `master_disciple` vs `coach_links`.
- Duplicidad `invitations` vs `invite_codes`.
- Seeds viejos que crean un esquema distinto al de migraciones.
- Carpetas infra vacias si no se usaran en el MVP.
- README de Vite.
- Swagger generado incompatible o dependencia `swag` desalineada.
- Logs `console.log` en servicios frontend.
- Rutas duplicadas raras en `program_handler.go`: `/api/programs/programs/:id/...`.
- `visibility public/unlisted` mientras el producto no es publico ni social.
- `roleHint` en signup si no tiene efecto real.
- Program versioning parcial si no se cierra en DB y API.

## Evaluacion de rescate

Conviene continuar sobre este repo, pero haciendo primero rescate tecnico. No recomiendo repo nuevo salvo que se quiera redisenar UI/UX completa y abandonar el backend actual. Hay suficiente dominio implementado para justificar rescate.

Estimacion de reutilizacion:

- Backend estructura y handlers: 45-55% reutilizable.
- Backend seguridad/authorization: 25-35% reutilizable, requiere endurecimiento.
- Esquema DB: 55-65% reutilizable, pero necesita migracion correctiva.
- Frontend paginas y servicios: 50-60% reutilizable.
- Tests/docs/infra: 10-20% reutilizable.
- Reutilizacion global estimada: 45-55%.

La parte mas valiosa es el esqueleto de entrenamiento: programas, dias, prescripciones, sesiones, sets, historial y pantallas conectadas. La parte que mas conviene reescribir conceptualmente es roles/permisos/ownership.

## Recomendacion final

No partir desde cero. Hacer una Fase 0 obligatoria de limpieza y build verde. Luego fijar el modelo de identidad y permisos antes de tocar features. El producto no puede avanzar confiablemente mientras cualquier usuario autenticado pueda crear recursos globales o manipular entidades ajenas por IDs.
