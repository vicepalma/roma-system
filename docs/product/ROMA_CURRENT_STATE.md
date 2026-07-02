# Roma System - Estado Actual

## Roles

- `coach`: crea ejercicios, programas, semanas, dias, prescripciones y asigna rutinas.
- `disciple`: ejecuta entrenamiento, registra sesiones/sets/cardio/notas, ve historial y crea check-ins.
- `disciple` sin coach: puede crear rutinas propias `self_training` usando ejercicios existentes.

## Flujo Disciple

- Ve y activa rutinas propias.
- Inicia entrenamiento desde la pantalla Entrenar.
- Selecciona dias del assignment activo.
- Registra sets.
- Finaliza sesiones.
- Revisa historial y resumen de sesiones.
- Crea y lista check-ins propios.

## Flujo Coach

- Gestiona discipulos vinculados.
- Crea ejercicios y programas.
- Asigna programas a discipulos.
- Revisa sesiones, historial, progreso y check-ins de discipulos vinculados.

## Self-Training

- `programs.kind` distingue `coach_program` y `self_training`.
- El self-training usa self-assignment explicito.
- Solo una self-training activa por disciple.
- No usa self-link en `coach_links`.
- No da permisos de coach al disciple.

## Sesiones e Historial

- No se puede iniciar sesion sobre assignment inactivo.
- La sesion cerrada usa `status=closed` y `ended_at`.
- Sesiones cerradas quedan en modo lectura para sets/cardio.
- Historial incluye programa, semana/dia, sets, ejercicios, volumen, status y filtros backend.

## Check-ins

- El disciple crea check-ins con fecha, peso opcional y notas.
- Coach vinculado puede listar/ver check-ins del disciple.
- No hay edicion/borrado, fotos, adjuntos ni graficos todavia.

## Tests

- Backend unitario con `go test ./...`.
- E2E/API con `ROMA_E2E_DB_URL` contra Postgres real de test.
- Frontend validado con `npm run build`.

## Pendientes

- Contexto de alumno mas claro en historial coach.
- Filtros/paginacion de check-ins.
- Demo FST-7 completa.
- Biblioteca de rutinas/plantillas si se prioriza.
