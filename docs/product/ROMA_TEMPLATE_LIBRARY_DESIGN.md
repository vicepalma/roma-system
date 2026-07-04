# Roma System - Diseno Biblioteca De Rutinas / Plantillas

## Objetivo

Disenar una biblioteca de rutinas/plantillas reutilizables sin implementarla todavia.

La biblioteca debe ayudar a crear programas mas rapido, pero no debe convertir Roma System en comunidad, red social, marketplace ni catalogo publico abierto.

## Principios

- Una plantilla no se ejecuta directamente.
- Una plantilla se copia o instancia como `program` antes de asignarse o activarse.
- Las sesiones siempre siguen usando `assignments`, `programs`, semanas, dias y prescripciones.
- La seguridad backend manda: ownership y relacion coach-disciple no dependen de ocultar botones en UI.
- El MVP debe favorecer plantillas privadas o curadas, no publicacion social.

## Tipos Propuestos

### Plantilla privada de coach

- La crea y administra un coach.
- Solo ese coach puede verla, editarla y usarla.
- Puede copiarla como programa propio para asignarla a disciples vinculados.
- No queda visible para otros coaches ni disciples.

### Plantilla propia de disciple

- No es prioritaria para MVP.
- Si se implementa, solo debe permitir crear una rutina propia `self_training` usando ejercicios existentes.
- No debe permitir crear ejercicios ni metodos globales.

### Plantilla Roma curada

- Puede existir como seed/demo o contenido interno controlado.
- Debe ser de solo lectura para usuarios normales.
- Puede copiarse como programa del coach o rutina propia del disciple, segun reglas definidas.
- No implica comunidad ni publicacion por usuarios.

## Relacion Con Programas Existentes

La plantilla debe ser fuente de copia, no reemplazo de `programs`.

Flujo recomendado:

1. Usuario elige una plantilla disponible.
2. Backend valida acceso a la plantilla.
3. Backend crea un `program` nuevo copiando semanas, dias y prescripciones.
4. El nuevo programa queda bajo ownership del usuario que lo instancio.
5. Desde ahi se usa el flujo actual:
   - coach asigna a disciple vinculado;
   - disciple sin coach activa como `self_training`.

No se deben crear `assignments` directamente desde la plantilla sin crear antes un programa concreto.

## Ownership Y Visibilidad

- Coach puede usar plantillas propias y plantillas Roma curadas.
- Coach no puede ver ni usar plantillas privadas de otro coach.
- Disciple puede usar plantillas Roma curadas solo si el producto lo habilita explicitamente.
- Disciple no puede publicar plantillas globales.
- Coach vinculado no obtiene permiso sobre rutinas propias privadas del disciple.

## Reglas De Ejercicios

- Las prescripciones de una plantilla deben referenciar ejercicios existentes.
- Si una plantilla Roma usa ejercicios seed, esos ejercicios deben existir en el catalogo local.
- Disciple no puede crear ejercicios al instanciar una plantilla.
- Si falta un ejercicio requerido, la instancia debe fallar con error claro o exigir una decision de mapeo en un WI posterior.

## Modelo Conceptual Futuro

Opciones a evaluar antes de implementar:

- Tabla nueva `program_templates` con estructura propia.
- Reutilizar `programs` con un nuevo `kind`, por ejemplo `template`, solo si no contamina reglas actuales.
- Guardar metadata de origen en programas copiados, por ejemplo `source_template_id`, si aporta trazabilidad.

La decision debe tomarse en un WI tecnico posterior, no en este diseno.

## Flujo UI Propuesto

Para coach:

1. Abrir seccion "Plantillas".
2. Ver plantillas propias y plantillas Roma curadas.
3. Previsualizar semanas, dias y ejercicios.
4. Crear programa desde plantilla.
5. Editar el programa copiado si necesita ajustes.
6. Asignarlo con el flujo actual.

Para disciple sin coach:

1. Abrir Mis rutinas.
2. Elegir crear rutina desde plantilla Roma si esta habilitado.
3. Copiar como `self_training`.
4. Activarla con la regla actual de una sola rutina propia activa.

## Fuera De Alcance Para MVP Inmediato

- Publicar plantillas entre usuarios.
- Marketplace.
- Likes, comentarios, seguidores o comunidad.
- Versionado avanzado de plantillas.
- Adjuntos, imagenes o media.
- Recomendador automatico.
- Instalar plantillas desde internet.

## Proximos Work Items Posibles

- Definir modelo tecnico de plantillas.
- Crear seed/demo de plantilla Roma FST-7.
- Implementar copia de plantilla a programa.
- Agregar UI minima de seleccion/previsualizacion.
- Agregar E2E de ownership para plantillas.
