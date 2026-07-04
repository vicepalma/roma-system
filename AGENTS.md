# Guia Para Agentes IA

## Producto

Roma System es software para entrenadores y discipulos. El coach construye y administra entrenamiento para otros. El disciple ejecuta entrenamiento y, si no tiene coach, puede crear rutinas propias usando ejercicios existentes.

## Reglas De Oro

- Trabajar una tarea por checkpoint.
- Hacer cambios chicos, verificables y alineados con el estado actual del repo.
- Seguridad en backend primero; la UI solo oculta acciones como apoyo.
- Mantener ownership y relacion coach-disciple.
- Actualizar `docs/tracking/ROMA_DEV_LOG.md` al final.
- No hacer commits automaticamente.
- No implementar ideas sueltas fuera de un Work Item.

## No Construir Ahora

- Comunidad, red social o simpatizantes.
- Biblioteca global de rutinas/plantillas sin checkpoint explicito.
- Fotos, adjuntos o mediciones avanzadas de check-ins.
- Graficos grandes o reporting avanzado.
- Features grandes mezcladas con bugfixes.

## Leer Antes De Trabajar

- `README.md`
- `ROADMAP.md`
- `docs/product/ROMA_CURRENT_STATE.md`
- `docs/product/ROMA_PRODUCT_RULES.md`
- `docs/engineering/DEVELOPMENT_SETUP.md`
- `docs/engineering/TESTING_GUIDE.md`
- `docs/tracking/ROMA_DEV_LOG.md`

## Modo De Trabajo Por Work Item

Todo cambio funcional debe venir desde un Work Item `ROMA-XXX`.

Para implementar un WI, el agente debe leer en este orden:

1. `AGENTS.md`
2. `ROADMAP.md`
3. `docs/tracking/ROMA_DEV_LOG.md`
4. `docs/work-items/ROMA-XXX.md`

El agente debe implementar solo lo pedido por el WI. Si encuentra algo necesario pero fuera de alcance, debe documentarlo como pendiente o sugerir un nuevo WI. Al final del checkpoint debe actualizar `docs/tracking/ROMA_DEV_LOG.md`, validar segun el tipo de cambio y responder con el formato esperado. No debe hacer commits automaticamente.

## Cambios Sin Commit

Los cambios sin commit no bloquean automaticamente un nuevo checkpoint.

Al inicio, el agente debe revisar `git status --short` y advertir si existen cambios locales. Si parecen pertenecer al WI anterior, debe recomendar hacer commit antes de continuar. Si el usuario insiste en continuar, puede seguir trabajando con cuidado. Solo debe bloquearse y pedir decision cuando los cambios sin commit choquen directamente con el nuevo WI o hagan riesgoso distinguir responsabilidades.

## Modo Recuperacion De Contexto

Si el usuario vuelve despues de tiempo o no recuerda el estado, por ejemplo:

```txt
Lee AGENTS.md.
No recuerdo donde quedamos.
Analiza el proyecto y recomiendame el proximo Work Item.
No implementes todavia.
```

El agente debe:

1. leer `AGENTS.md`;
2. leer `ROADMAP.md`;
3. leer `docs/tracking/ROMA_DEV_LOG.md`;
4. revisar `docs/work-items/` si existe;
5. resumir estado actual;
6. identificar ultimo WI completado;
7. recomendar el proximo WI;
8. no implementar nada.

## Si El Usuario Entrega Un Work Item

Ejemplo:

```txt
Lee AGENTS.md.

Work item:
ROMA-020

Implementa.
```

El agente debe abrir `docs/work-items/ROMA-020.md`, confirmar objetivo y alcance desde el WI, implementar solo lo pedido, validar segun el tipo de cambio, actualizar tracking y responder con resumen final.

## Si El Usuario No Entrega Un Work Item

El agente debe leer `ROADMAP.md` y `docs/tracking/ROMA_DEV_LOG.md`, resumir el estado actual, recomendar el proximo WI y no implementar todavia.

## Si El Work Item No Existe

El agente debe:

- no inventar implementacion;
- proponer crear el archivo;
- sugerir contenido minimo;
- esperar confirmacion, salvo que el usuario pida explicitamente crear el WI;
- si crea solo documentacion, no implementar feature.

## Validaciones Esperadas

Backend:

```bash
cd backend
GOCACHE=/tmp/roma-go-cache go test ./...
```

E2E/API:

```bash
cd backend
ROMA_E2E_DB_URL='postgres://roma:roma@localhost:5432/roma_e2e?sslmode=disable' GOCACHE=/tmp/roma-go-cache go test ./... -run E2E -count=1
```

Frontend:

```bash
cd frontend/roma-web
npm run build
```

Reglas de validacion:

- Si solo se toca documentacion, no es obligatorio correr tests; basta `git status --short`.
- Si se toca backend, correr backend tests.
- Si se toca contrato API, correr E2E.
- Si se toca frontend, correr build.
- Si se toca migracion, validar migraciones limpias si es razonable.

## Respuesta Final Esperada

Incluir cambios realizados, archivos modificados, validaciones ejecutadas, riesgos pendientes y proximo checkpoint recomendado.
