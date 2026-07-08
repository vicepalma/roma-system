# ROMA-030 - Coach puede usar entrenamiento personal

Estado: Done
Tipo: feature
Fecha: 2026-07-04
Autor: Codex

## Objetivo

Permitir que un coach use Entrenar, Mis rutinas, Historial personal y Check-ins personales, sin dejar de ser coach.

## Contexto

ROMA-028 detecto durante la demo local que el maestro tambien puede entrenarse a si mismo. El rol `coach` debe conservar sus capacidades de entrenador, pero tambien acceder a un flujo personal equivalente al self-training de un disciple independiente.

## Alcance

- Revisar reglas actuales de self-training y permisos por rol.
- Permitir que un coach cree/active rutinas personales para si mismo.
- Permitir que un coach use Entrenar para su propio entrenamiento.
- Permitir historial personal del coach sin mezclarlo con historial de disciples.
- Permitir check-ins personales del coach.
- Mantener sin cambios la relacion coach-disciple y los guards de ownership.

## Fuera de alcance

- Biblioteca de rutinas.
- Cambiar modelo de relacion coach-disciple.
- Rutinas asignadas por maestro en Mis rutinas.
- Redisenar UI completa.
- Comunidad, red social o reporting avanzado.

## Archivos probables

- Backend de auth/guards/programs/assignments/sessions/history/checkins.
- Frontend de navegacion, Entrenar, Mis rutinas, Historial y Check-ins.
- Tests backend y E2E/API.
- `docs/tracking/ROMA_DEV_LOG.md`

## Validaciones requeridas

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- E2E/API con `ROMA_E2E_DB_URL` si se toca contrato o permisos.
- `cd frontend/roma-web && npm run build`

## Criterios de aceptacion

- Un coach puede acceder a flujo personal sin perder funcionalidades de coach.
- Un coach puede crear/activar rutina personal y entrenarla.
- Historial personal del coach no se mezcla con historial de sus disciples.
- Check-ins personales del coach quedan separados de los check-ins de disciples.
- Guards de ownership siguen protegiendo recursos ajenos.

## Riesgos / notas

- Evitar tratar al coach como disciple vinculado a si mismo mediante `coach_links`.
- Mantener `programs.kind='self_training'` como flujo personal si aplica.

## Resultado esperado del agente

Implementar solo este WI, validar backend/frontend/E2E segun cambios, actualizar tracking y no hacer commit automaticamente.

## Resultado

- Backend permite que usuarios `coach` creen y muten programas `self_training` propios, manteniendo `coach_program` para programas de alumnos.
- Backend permite check-ins personales para coach usando los endpoints personales `/api/checkins`.
- Entrenar, Mis rutinas, Historial personal y Check-ins quedan accesibles para coach desde la UI.
- Mis rutinas del coach permite crear una rutina personal o un programa para alumnos sin perder capacidades de coach.
- E2E cubre check-in personal del coach y flujo self-training personal del coach con activacion, sesion, set e historial.

## Validado

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- `cd backend && ROMA_E2E_DB_URL='postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable' GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`
- `cd frontend/roma-web && npm run build`
