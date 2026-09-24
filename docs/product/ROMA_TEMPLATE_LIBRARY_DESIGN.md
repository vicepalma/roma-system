# Roma System - Diseno Biblioteca De Rutinas / Plantillas

## Objetivo

Definir la biblioteca de rutinas que se implementara en ROMA-033, sin introducir comunidad ni un catalogo publico abierto.

La biblioteca debe ayudar a crear programas mas rapido, pero no debe convertir Roma System en comunidad, red social, marketplace ni catalogo publico abierto.

## Principios

- Una plantilla no se ejecuta directamente.
- Una plantilla se copia o instancia como `program` antes de asignarse o activarse.
- Las sesiones siempre siguen usando `assignments`, `programs`, semanas, dias y prescripciones.
- La seguridad backend manda: ownership y relacion coach-disciple no dependen de ocultar botones en UI.
- La biblioteca v1 contiene plantillas globales curadas, visibles solo para usuarios autenticados.
- Una copia importada es independiente: editarla o eliminarla nunca modifica la plantilla.

## Tipos Y Estados V1

### Plantilla global

- Es una rutina reutilizable con estructura propia de semanas, dias y prescripciones.
- Puede ser creada por un coach autorizado o administrada por un admin cuando exista ese rol operativo.
- `draft` solo es visible para su propietario; `published` es visible para usuarios autenticados.
- No se ejecuta ni recibe assignments directamente.

### Rutina importada/personal

- Es un `program` nuevo creado mediante copia profunda de una plantilla publicada.
- Queda bajo ownership del usuario que importa.
- Un coach puede editarla y asignarla a un disciple vinculado; un disciple/independent puede usarla como `self_training`.

### Asignacion de maestro

- Sigue siendo la relacion existente entre un `program` concreto y un disciple.
- Nunca se crea directamente desde una plantilla.
- El disciple puede verla en Mis rutinas, pero no editarla ni eliminarla desde su espacio personal.

## Relacion Con Programas Existentes

La plantilla debe ser fuente de copia, no reemplazo de `programs`.

Flujo recomendado:

1. Usuario elige una plantilla disponible.
2. Backend valida acceso a la plantilla.
3. Backend crea un `program` nuevo copiando profundamente semanas, dias y prescripciones.
4. El nuevo programa queda bajo ownership del usuario que lo instancio.
5. Desde ahi se usa el flujo actual:
   - coach asigna a disciple vinculado;
   - disciple sin coach activa como `self_training`.

No se deben crear `assignments` directamente desde la plantilla sin crear antes un programa concreto.

## Ownership, Publicacion Y Visibilidad

- Coach puede crear y administrar sus plantillas; solo el propietario puede editar, publicar, despublicar o eliminar un draft.
- Un admin puede administrar plantillas globales si el rol existe en la instalacion; no se agrega un rol admin en ROMA-033.
- Disciple e independent pueden listar plantillas `published` e importarlas, pero no crearlas ni publicarlas.
- Usuarios no autenticados no pueden listar ni importar plantillas.
- Coach vinculado no obtiene permiso sobre rutinas propias privadas del disciple.
- Retirar o editar una plantilla no cambia programas ya importados.

## Reglas De Ejercicios

- Las prescripciones de una plantilla deben referenciar ejercicios existentes.
- Si una plantilla Roma usa ejercicios seed, esos ejercicios deben existir en el catalogo local.
- Disciple no puede crear ejercicios al instanciar una plantilla.
- Si falta un ejercicio requerido, la instancia debe fallar con error claro o exigir una decision de mapeo en un WI posterior.

## Modelo Tecnico Minimo Para ROMA-033

- Usar una entidad separada `program_templates`; no reutilizar `programs.kind` como `template`.
- Guardar ownership, estado (`draft`/`published`), titulo, notas y fechas de publicacion en la plantilla.
- Mantener estructura propia de plantilla para semanas, dias y prescripciones, referenciando ejercicios existentes.
- Registrar en el `program` importado el origen de plantilla cuando sea necesario para trazabilidad (`source_template_id` o metadata equivalente).
- La importacion debe ejecutarse en una transaccion y copiar toda la estructura; si falla, no deja programa parcial.
- No se requiere versionado avanzado: una plantilla publicada puede actualizarse para futuras importaciones, sin afectar copias existentes.

## Flujo UI Minimo V1

Para coach:

1. Abrir Biblioteca.
2. Ver plantillas publicadas y su informacion basica.
3. Previsualizar estructura y ejercicios.
4. Importar como programa propio.
5. Editar el programa copiado si corresponde.
6. Usar los flujos existentes de asignacion o self-training.

Para disciple/independent:

1. Abrir Biblioteca desde la navegacion disponible.
2. Ver solo plantillas publicadas.
3. Importar una copia personal como `self_training`.
4. Activarla desde Mis rutinas con la regla actual de una sola rutina propia activa.

## Alcance ROMA-033

- Listar plantillas publicadas para usuarios autenticados.
- Crear/publicar plantillas mediante permisos de coach/admin definidos arriba.
- Previsualizar estructura minima.
- Importar una plantilla como copia profunda a un `program` propio.
- Mostrar la copia en Mis rutinas y proteger la independencia entre copia y plantilla.
- Cubrir ownership, publicacion e importacion con tests backend/E2E.

## Fuera De Alcance V1

- Publicacion anonima o catalogo publico sin autenticacion.
- Plantillas propias de disciple como fuente publicable.
- Marketplace.
- Likes, comentarios, seguidores o comunidad.
- Versionado avanzado de plantillas.
- Adjuntos, imagenes o media.
- Recomendador automatico.
- Instalar plantillas desde internet.
- Flujos avanzados de versionado, archivado o migracion de plantillas.
