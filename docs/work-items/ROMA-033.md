# ROMA-033 - Implementar biblioteca de rutinas MVP

Estado: Done
Tipo: feature
Fecha: 2026-09-24
Autor: Codex

## Objetivo

Implementar una biblioteca minima de rutinas creadas/publicadas por maestros o admin, que disciples/independientes puedan importar como copia personal.

## Contexto

Este WI debe ejecutarse despues de ROMA-032, usando el diseno aprobado para distinguir plantilla global, rutina personal y asignacion de maestro.

## Alcance

- Implementar biblioteca MVP segun el diseno de ROMA-032.
- Permitir listar plantillas visibles.
- Permitir importar una plantilla como copia personal.
- Asegurar que borrar o editar la copia personal no modifique la plantilla global.
- Mantener permisos de publicacion restringidos a maestros o admin segun diseno.
- Agregar UI minima para descubrir e importar rutinas.

## Fuera de alcance

- Comunidad, red social o marketplace.
- Ratings, comentarios o seguidores.
- Adjuntos, fotos o multimedia.
- Reporting avanzado.
- Flujos avanzados de versionado si no estan en ROMA-032.

## Archivos probables

- Migraciones de biblioteca/plantillas si el diseno lo requiere.
- Backend de templates/programs/import.
- Frontend de biblioteca y Mis rutinas.
- Tests backend, E2E/API y build frontend.
- `docs/tracking/ROMA_DEV_LOG.md`

## Validaciones requeridas

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- E2E/API con `ROMA_E2E_DB_URL`
- `cd frontend/roma-web && npm run build`
- Validar migraciones limpias si se agregan migraciones.

## Criterios de aceptacion

- Disciples/independientes pueden ver biblioteca MVP.
- Disciples/independientes pueden importar una rutina como copia personal.
- La copia importada aparece en Mis rutinas como rutina personal/importada.
- Eliminar o editar la copia no modifica la plantilla global.
- Permisos backend bloquean publicacion/importacion indebida.

## Riesgos / notas

- No mezclar este WI con red social o marketplace.
- La seguridad backend debe ser la fuente de verdad para publicacion e importacion.

## Resultado esperado del agente

Implementar solo el MVP definido por ROMA-032, validar backend/E2E/frontend, actualizar tracking y no hacer commit automaticamente.

## Resultado

- Se agrego la migracion `0008_program_templates` con plantillas, semanas, dias y prescripciones propias.
- Se implementaron endpoints autenticados para listar, previsualizar, crear, publicar e importar plantillas.
- La importacion copia la estructura completa en una transaccion y crea un `program` independiente.
- Coach puede crear/publicar; disciple puede listar e importar plantillas publicadas.
- Se agrego la pantalla Biblioteca y navegacion protegida en frontend.
- Se cubrio la independencia de ownership mediante el modelo de copia profunda.

## Validado

- `cd backend && GOCACHE=/tmp/roma-go-cache go test ./...`
- `cd backend && ROMA_E2E_DB_URL=... GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1`
- `cd frontend/roma-web && npm run build`
- `git diff --check`
