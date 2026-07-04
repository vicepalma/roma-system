# Roma MVP Hardening

Fecha: 2026-07-04
Work Item: ROMA-025

## Revision

- Seguridad backend: se revisaron rutas criticas de coach, assignments, sesiones, programas, historial y check-ins.
- Validaciones: se priorizo mantener ownership y relacion coach-disciple como regla backend.
- UX critica: se revisaron flujos Entrenar, Check-ins, Historial, vista de coach y demo local desde documentacion y contratos actuales.
- Demo local: FST-7 queda como seed opcional, no como migracion obligatoria.

## Ajuste aplicado

- `/api/coach/assignments/:id/calendar` requiere rol `coach`.
- `/api/coach/assignments/:id/activate` requiere rol `coach`.
- La activacion de assignments de coach ya no acepta al disciple como actor valido; self-training mantiene su ruta separada `/api/programs/:id/self-assignment`.

## Bloqueantes MVP

- No se encontraron bloqueantes criticos nuevos para continuar hacia preparacion demo.

## Riesgos pendientes

- `master_disciple` sigue como legacy/compatibilidad; `coach_links` es la fuente operativa.
- `invitations` e `invite_codes` siguen duplicados conceptualmente y conviene consolidarlos en un WI posterior.
- La biblioteca de plantillas sigue en diseno; no hay implementacion funcional de plantillas.
- Borrado/archivado de check-ins queda fuera de alcance.
- Reporting avanzado, graficos grandes, adjuntos y comunidad siguen fuera de foco MVP.

## Recomendacion

Continuar con ROMA-026 para pulido UI general MVP o ROMA-027 para preparacion demo local, segun si se prioriza polish visual o ensayo de demo.
