# ROMA-023 - Biblioteca de rutinas / plantillas Roma, diseno primero

Estado: Done
Tipo: chore
Fecha: 2026-07-04
Autor: Codex

## Objetivo

Disenar antes de implementar una biblioteca de rutinas/plantillas reutilizables.

## Alcance

- Documentar reglas de producto para plantillas.
- Definir ownership, visibilidad y relacion con programas existentes.
- Proponer flujo antes de tocar backend/frontend funcional.

## Fuera de alcance

- Implementar biblioteca global.
- Publicacion social o comunidad.
- Marketplace o plantillas publicas sin decision explicita.

## Validaciones esperadas

- Si es solo diseno/documentacion: `git status --short`.
- Si se implementa algo despues de una decision explicita, validar backend/frontend segun cambios.

## Resultado

- Se creo `docs/product/ROMA_TEMPLATE_LIBRARY_DESIGN.md`.
- El diseno define tipos de plantillas, ownership, visibilidad y relacion con `programs`.
- Se mantiene fuera de alcance la implementacion de biblioteca global, comunidad, marketplace y publicacion por usuarios.
- Se propone que las plantillas se copien a programas antes de asignarse o activarse.

## Validado

- `git status --short`
