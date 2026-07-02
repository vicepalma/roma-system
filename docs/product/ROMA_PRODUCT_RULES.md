# Roma System - Reglas De Producto

## Coach

- Puede crear ejercicios.
- Puede crear metodos/protocolos, programas, semanas, dias y prescripciones.
- Puede asignar programas propios a discipulos vinculados.
- Puede revisar progreso, historial, sesiones y check-ins de discipulos vinculados.
- No puede ver datos de discipulos ajenos.
- No crea check-ins por el disciple en el estado actual.

## Disciple

- No puede crear ejercicios.
- No puede crear metodos/protocolos globales.
- No puede administrar usuarios ni asignar rutinas a otros.
- Puede ejecutar assignments propios.
- Puede registrar sesiones, sets, cardio, notas e historial propio.
- Puede crear check-ins propios.

## Disciple Sin Coach

- Puede crear rutinas propias simples con `programs.kind = self_training`.
- Puede agregar semanas/dias/prescripciones usando ejercicios existentes.
- Puede activar una rutina propia para entrenar.
- Solo puede tener una self-training activa a la vez.

## Sesiones

- No se inicia sesion sobre assignment inactivo.
- `day_id` debe pertenecer al programa asignado.
- La prescripcion debe pertenecer al dia de la sesion.
- Una sesion cerrada queda en modo lectura.
- Historial pasado sigue visible aunque el assignment quede inactivo.

## Check-ins

- Los crea el disciple.
- Coach vinculado puede verlos.
- Coach no vinculado y otros disciples no pueden acceder.
- Sin adjuntos/fotos ni mediciones avanzadas por ahora.

## Fuera De Foco

- Comunidad.
- Red social.
- Simpatizantes.
- Features grandes sin checkpoint.
