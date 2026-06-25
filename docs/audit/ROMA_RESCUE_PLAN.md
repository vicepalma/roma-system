# ROMA Rescue Plan

## Principios de rescate

- No agregar features nuevas hasta que backend y frontend compilen.
- Resolver primero identidad, permisos y ownership.
- Mantener el dominio de entrenamiento ya modelado, pero corregir nombres y contratos.
- Eliminar duplicidades antes de ampliar el modelo.
- Priorizar flujos del producto actual: maestro gestiona, discipulo ejecuta, discipulo sin maestro crea rutina propia con piezas existentes.

## Fase 0 - Estabilizacion y limpieza

Objetivo: dejar el repo ejecutable y entendible.

Tareas:

- Crear README raiz con comandos reales de desarrollo.
- Corregir build backend:
  - regenerar Swagger con version compatible o retirar campos `LeftDelim`/`RightDelim`;
  - confirmar `go test ./...` verde.
- Corregir build frontend:
  - ajustar tipos en `Exercises.tsx`;
  - corregir `ProgramDetail.tsx` para usar `week_index`;
  - corregir payload de prescripcion o hacer `method_id` opcional.
- Definir una sola convencion env:
  - preferir `DB_URL` o `DATABASE_DSN`, no ambas;
  - documentar `JWT_SECRET`, TTLs, CORS y `VITE_API_BASE`.
- Decidir una sola tabla de relacion maestro-discipulo:
  - recomendacion: usar `coach_links` o renombrar todo a `master_disciple`, pero no ambos;
  - crear migracion correctiva.
- Eliminar o archivar seeds viejos que no representen el esquema oficial.
- Decidir si `invitations` o `invite_codes` sobrevive.
- Quitar rutas duplicadas de programas.
- Remover carpetas infra vacias o marcarlas explicitamente como futuras.

Criterio de salida:

- `go test ./...` verde.
- `npm run build` verde.
- DB limpia creada por migraciones oficiales y seed minimo.
- Login funcional contra compose local.

## Fase 1 - Core tecnico: identidad, roles y permisos

Objetivo: que el producto tenga reglas enforceables, no solo pantallas.

Modelo recomendado:

- `users`
- `user_roles` o columna `role` si sera simple.
- Roles iniciales:
  - `coach`
  - `disciple`
- Estado derivado adicional:
  - discipulo con maestro;
  - discipulo sin maestro.

Reglas:

- Maestro puede crear ejercicios, metodos y programas.
- Discipulo con maestro solo ejecuta y registra.
- Discipulo sin maestro puede crear rutinas propias simples usando ejercicios existentes.
- Discipulo no puede crear ejercicios globales.
- Discipulo no puede crear metodos globales.
- Discipulo no puede administrar otros usuarios.

Tareas backend:

- Persistir rol real o una politica equivalente auditable.
- Crear middleware/guards:
  - `RequireRole(coach)`;
  - `RequireSelfOrCoachOf(discipleID)`;
  - `RequireProgramOwner(programID)`;
  - `RequireAssignmentOwnerOrCoach(assignmentID)`;
  - `RequireSessionOwnerOrCoach(sessionID)`.
- Aplicar guards a ejercicios, programas, semanas, dias, prescripciones, asignaciones, sesiones, sets e historial.
- Validar que `assignment_id`, `day_id` y `prescription_id` correspondan entre si.
- Validar `JWT_SECRET` obligatorio al arrancar.
- Unificar CORS por config.

Tareas frontend:

- Cargar `/me` despues de login y guardar usuario/rol.
- Sidebar por rol.
- Rutas protegidas por rol.
- Ocultar acciones no permitidas.
- No confiar solo en frontend; backend manda.

Criterio de salida:

- Casos negativos manuales o automatizados:
  - discipulo no crea ejercicio;
  - discipulo con maestro no crea programa para otros;
  - coach no ve discipulo ajeno;
  - usuario no edita sesion ajena;
  - self-training solo toca recursos propios.

## Fase 2 - Training Core

Objetivo: consolidar el constructor de entrenamiento.

Tareas:

- Mantener y endurecer ejercicios:
  - catalogo global;
  - CRUD solo coach/admin si se mantiene admin futuro;
  - filtros por musculo/equipo/tags.
- Decidir alcance de metodos/protocolos:
  - para MVP, usar lectura de metodos existentes;
  - crear FST-7 como seed/demo oficial;
  - dejar CRUD de metodos fuera salvo que sea imprescindible.
- Consolidar programas:
  - owner obligatorio;
  - semanas;
  - dias;
  - prescripciones;
  - cardio prescrito;
  - notas;
  - orden de ejercicios.
- Definir versionado:
  - o se implementa bien `program_versions`;
  - o se elimina del MVP y se usa `program.version` simple sin snapshots.
- Crear flujo self-training:
  - rutina propia simple;
  - dias elegidos por el discipulo;
  - prescripciones usando ejercicios existentes;
  - sin crear ejercicios ni metodos.

Criterio de salida:

- Coach crea programa completo y lo asigna.
- Discipulo sin maestro crea rutina propia simple.
- La pauta FST-7 demo puede cargarse desde seed oficial.

## Fase 3 - Ejecucion

Objetivo: que el discipulo pueda entrenar sin friccion y con datos confiables.

Tareas:

- Consolidar "hoy entreno":
  - programa activo;
  - dia calculado o elegido;
  - prescripciones;
  - cardio prescrito;
  - sesion abierta actual.
- Sesiones:
  - iniciar;
  - registrar sets;
  - editar/borrar set propio;
  - registrar cardio;
  - cerrar sesion;
  - notas.
- Validar integridad:
  - sets solo sobre prescripciones del dia/sesion;
  - una sesion abierta por discipulo, salvo decision contraria;
  - sesiones cerradas con reglas claras de edicion.
- Historial:
  - sesiones recientes;
  - volumen por ejercicio/musculo;
  - PRs basicos;
  - plan vs realizado.

Criterio de salida:

- Discipulo con maestro ejecuta rutina asignada.
- Discipulo sin maestro ejecuta rutina propia.
- Historial refleja los logs sin exponer datos ajenos.

## Fase 4 - Maestro Dashboard

Objetivo: que el coach tenga gestion de alumnos y seguimiento basico.

Tareas:

- Lista de discipulos.
- Invitaciones o links maestro-discipulo, con una sola tabla.
- Perfil/detalle de discipulo.
- Asignaciones activas e historicas.
- Ultimos entrenamientos.
- Adherencia simple.
- Check-ins:
  - crear check-in desde discipulo;
  - coach ve check-ins de sus discipulos;
  - peso, notas, attachments si se mantiene.

Criterio de salida:

- Coach ve solo sus discipulos.
- Coach asigna programa y revisa progreso.
- Check-ins basicos funcionando end to end.

## Fase 5 - Calidad, demo y hardening

Objetivo: preparar el MVP para uso controlado.

Tareas:

- Tests backend:
  - auth;
  - guards;
  - ownership;
  - sesiones/sets;
  - asignaciones.
- Tests frontend minimos:
  - login;
  - rutas por rol;
  - crear programa;
  - registrar sesion.
- Seed demo FST-7:
  - semanas;
  - dias;
  - grupos musculares;
  - ejercicios;
  - series/reps/descansos;
  - cardio;
  - notas;
  - progresion.
- Documentacion:
  - arquitectura;
  - comandos;
  - datos demo;
  - matriz de permisos.
- Docker Compose completo:
  - Postgres;
  - migrator;
  - backend;
  - frontend opcional.

Criterio de salida:

- Un desarrollador nuevo puede levantar el sistema con README.
- Demo FST-7 recorre maestro -> asignacion -> ejecucion -> historial.
- Tests criticos pasan en CI local.

## Orden recomendado de trabajo

1. Build verde.
2. Migracion correctiva de relaciones maestro-discipulo.
3. Modelo de rol/permisos.
4. Ownership en backend.
5. Sidebar y rutas por rol.
6. Self-training.
7. Demo FST-7.
8. Check-ins y reportes.

## Riesgos si se ignora el orden

- Seguir agregando UI aumentara la deuda porque las reglas de producto no estan enforced.
- Agregar FST-7 sin corregir permisos solo creara mas datos inseguros.
- Mantener seeds divergentes hara que bugs aparezcan segun quien haya inicializado la DB.
- Seguir con roles derivados impedira expresar correctamente "discipulo sin maestro".
