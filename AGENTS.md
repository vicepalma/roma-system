# Guia Para Agentes IA

## Producto

Roma System es software para entrenadores y discipulos. El coach construye y administra entrenamiento para otros. El disciple ejecuta entrenamiento y, si no tiene coach, puede crear rutinas propias usando ejercicios existentes.

## Reglas de Oro

- Trabajar una tarea por checkpoint.
- Hacer cambios chicos, verificables y alineados con el estado actual del repo.
- Seguridad en backend primero; la UI solo oculta acciones como apoyo.
- Mantener ownership y relacion coach-disciple.
- Actualizar `docs/tracking/ROMA_DEV_LOG.md` al final.
- No hacer commits automaticamente.

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

## Respuesta Final Esperada

Incluir cambios realizados, archivos modificados, validaciones ejecutadas, riesgos pendientes y proximo checkpoint recomendado.
