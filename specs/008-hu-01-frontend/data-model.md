# Data Model: Pantalla de creación de historias (HU-01, frontend)

Esta historia **no crea ni modifica tablas**. La historia que se crea se persiste en la
tabla `historias` (hoy `backlog_items`) documentada en
[`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md) §4.4; el backend ya aplica sus
restricciones. Acá solo se describe lo que la pantalla necesita para armar el formulario
y el cuerpo del `POST`.

## Tipos del frontend

| Tipo | Definición | Uso |
|---|---|---|
| `Prioridad` | `'M' \| 'S' \| 'C' \| 'W'` | Valor del `<select>` de prioridad y del campo `prioridad` del POST |
| `ValorNegocio` | `1 \| 2 \| 3 \| 5 \| 8 \| 13 \| 21` | Valor numérico admitido para `valor_negocio` |
| `FormularioHistoria` | `{ titulo: string; descripcion: string; prioridad: Prioridad; valorNegocio: '' \| `${ValorNegocio}` }` | Estado del formulario; `''` representa "Sin valor" |
| `CuerpoCrearHistoria` | `{ proyecto_id: number; titulo: string; descripcion: string; prioridad: Prioridad; valor_negocio?: ValorNegocio }` | Cuerpo que consume `POST /backlog` |
| `HistoriaCreada` | `{ id; proyecto_id; titulo; descripcion; prioridad; estado; valor_negocio: number \| null; estimacion_sp: number \| null }` | Respuesta 201 que alimenta la confirmación |
| `ErrorDeCampo` | `Partial<Record<'titulo' \| 'descripcion' \| 'prioridad' \| 'valor_negocio', string>>` | Mensajes locales y del servidor por campo |

## Reglas de la pantalla

| Regla | Detalle |
|---|---|
| Proyecto | Se lee de la URL (`:id`). Solo se acepta `/^\d+$/` con valor > 0 y `Number.isSafeInteger`; si no, aviso general y formulario deshabilitado |
| Título | Obligatorio; se evalúa después de `trim`; máximo 200 puntos de código |
| Descripción | Obligatoria; se evalúa después de `trim`; máximo 2000 puntos de código |
| Prioridad | `M`, `S`, `C` o `W`; por defecto `M` |
| Valor de negocio | Opcional; `''` (Sin valor) o `1, 2, 3, 5, 8, 13, 21`; por defecto `''` |
| Envío | `titulo` y `descripcion` recortados; `valor_negocio` omitido si es `''`; una sola llamada por guardado |
| Validación simultánea | Se muestran todos los errores de campo del intento, no solo el primero |
| Mensajes de error | Permanecen hasta el próximo intento; no se borran al escribir |

## Estados de la pantalla

```text
inicial ──(escribir)──► editando ──(Guardar con errores)──► editando + advertencias
   ▲                        │
   │                        ├──(Guardar válido)──► enviando ──► éxito: confirmación + formulario inicial
   │                        │                          └─────► error de campo: mensaje junto al campo
   │                        │                          └─────► error general: aviso arriba del formulario
   └────────────────────────┴──(nuevo intento)──► reemplaza confirmación y errores anteriores
```

- **inicial**: `titulo` y `descripcion` vacíos, prioridad `M`, valor `''`.
- **enviando**: botón "Guardar" deshabilitado e `LoadingIndicator` visible; se ignoran
  nuevos envíos.
- **éxito**: confirmación con título, prioridad legible, estado "Nueva", sin estimación y,
  si corresponde, valor de negocio. El formulario vuelve a **inicial** y conserva el
  proyecto de la URL.
- **error de campo**: el mensaje se asocia al control (`aria-describedby`) y el campo se
  marca inválido.
- **error general**: `ErrorMessage` arriba del formulario; los datos ingresados no se
  borran.

## Fuera de alcance

No hay listado, edición, borrado, estimación ni Story Points en esta pantalla (HU-15 y
HU-02 los cubren). Tampoco se consulta el proyecto antes de guardar: el 404 surge del
`POST`.
