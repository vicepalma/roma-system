# Demo Local

Esta guia deja una demo local rapida para mostrar Roma System con backend, frontend y datos opcionales FST-7. No requiere cambios de codigo ni deploy.

## Preparar entorno

Desde la raiz del repo, levantar Postgres, migraciones oficiales y backend:

```bash
docker compose -f database/docker-compose.yml up --build
```

El backend queda en `http://localhost:8080`. Esperar a que el servicio responda healthy antes de abrir el frontend.

En otra terminal, levantar frontend:

```bash
cd frontend/roma-web
npm run dev
```

El frontend queda en `http://localhost:5173` y usa `VITE_API_BASE=http://localhost:8080` si se necesita configurar el endpoint explicitamente.

## Seed FST-7 opcional

Las migraciones oficiales no cargan datos demo FST-7. Para una demo con rutina ya asignada, aplicar el seed opcional sobre una DB local ya migrada:

```bash
psql 'postgres://roma:roma@localhost:5432/roma?sslmode=disable' -v ON_ERROR_STOP=1 -f database/seed/demo_fst7.up.sql
```

Para retirar esos datos demo:

```bash
psql 'postgres://roma:roma@localhost:5432/roma?sslmode=disable' -v ON_ERROR_STOP=1 -f database/seed/demo_fst7.down.sql
```

Mantener estos datos fuera de migraciones oficiales. Las cuentas `@roma.demo` son solo para uso local/demo.

## Usuarios demo

Password comun: `secret123`.

Demo FST-7 opcional:

- `fst7.coach@roma.demo` - Coach Demo FST-7.
- `fst7.disciple@roma.demo` - Disciple Demo FST-7 con assignment activo `[DEMO] FST-7 Hipertrofia - 5 dias`.

Seed minimo oficial:

- Coach: `roshi@kamehouse.example`, `baki.hanma@example.example`, `ikki@saint.example`.
- Disciple: `krillin@kamehouse.example`, `yamcha@capsule.example`, `goku@capsule.example`, `retsu@shinshinkai.example`, `katsumi@shinshinkai.example`, `jack.hanma@example.example`.

Ver tambien `docs/engineering/DEMO_USERS.md`.

## Flujo demo sugerido

Coach:

1. Entrar como `fst7.coach@roma.demo`.
2. Revisar dashboard de discipulos.
3. Abrir el detalle de `fst7.disciple@roma.demo`.
4. Revisar rutina activa, sesiones recientes, historial y check-ins si existen.
5. Abrir asignaciones para confirmar el programa `[DEMO] FST-7 Hipertrofia - 5 dias`.

Disciple:

1. Entrar como `fst7.disciple@roma.demo`.
2. Abrir Entrenar y revisar la rutina activa.
3. Seleccionar un dia FST-7 e iniciar sesion.
4. Registrar al menos un set.
5. Cerrar la sesion y revisar el resumen.
6. Abrir Historial para confirmar que la sesion cerrada aparece.
7. Crear o editar un check-in basico desde Check-ins.

## Checklist manual de cierre v1

Ejecutar también en 360, 390 y 430 px y en desktop para coach, disciple e independent. La declaración final de v1 requiere además backend unitario, E2E/API y migraciones desde DB limpia verdes.

## Checklist manual

- Backend responde en `http://localhost:8080/healthz`.
- Frontend abre en `http://localhost:5173`.
- Login coach demo funciona.
- Login disciple demo funciona.
- El coach ve solo discipulos vinculados.
- El disciple ve assignment activo FST-7.
- El disciple puede iniciar y cerrar una sesion.
- La sesion cerrada aparece en Historial.
- El check-in basico se crea o edita correctamente.
- El seed demo puede retirarse con `database/seed/demo_fst7.down.sql`.
