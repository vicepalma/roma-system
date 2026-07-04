# Work Items Roma System

Un Work Item es la unidad minima de trabajo para cambiar Roma System. Debe describir un objetivo concreto, alcance claro, validaciones esperadas y criterios de aceptacion.

## Formato

- Los IDs usan formato `ROMA-XXX`, por ejemplo `ROMA-020`.
- Cada WI vive en `docs/work-items/ROMA-XXX.md`.
- Cada WI debe ser pequeno, verificable y con alcance claro.

## Uso Con Codex / Agentes

El flujo recomendado es:

```txt
Lee AGENTS.md.

Work item:
ROMA-020

Implementa.
```

El agente debe leer `AGENTS.md`, `ROADMAP.md`, `docs/tracking/ROMA_DEV_LOG.md` y el archivo del WI antes de implementar. Ningun agente debe implementar ideas sueltas fuera de un WI.

## Estados

- `Planned`: definido, aun no implementado.
- `In Progress`: en ejecucion.
- `Done`: completado y validado segun su alcance.
- `Blocked`: no puede avanzar sin decision, dato externo o cambio previo.
