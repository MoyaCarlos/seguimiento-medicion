# Contrato de UI: tarjetas plegables y confirmación (#219)

## Tarjeta de historia (`SprintBacklog`)

```html
<li>
  <details class="historia-tarjeta">
    <summary>
      <h4>{titulo}</h4>
      <span class="historia-marca" data-estado="{estado}">{estado}</span>
      <span class="historia-tarjeta__sp">{estimacion_sp} SP | Sin estimar</span>
    </summary>
    <div class="historia-tarjeta__detalle">
      <p>{etiqueta de prioridad} ({prioridad})</p>
      <p>Completada el {fecha_finalizacion}</p> <!-- solo si está completada -->
      <EstadoHistoria />   <!-- solo si el Sprint está activo -->
      <EsfuerzoHistoria />
    </div>
  </details>
</li>
```

- El resumen no incluye la prioridad ni la fecha de finalización.
- El detalle conserva el orden de los datos y el esfuerzo existentes.
- La tarjeta es plegable con ratón y teclado mediante `<summary>`.

## `EstadoHistoria`

Estados internos de UI:

| Fase | Qué se muestra |
|---|---|
| `reposo` | Botón «Iniciar» o «Completar» (o nada si la historia está Completada). |
| `confirmando` | Pregunta y botones «Confirmar» y «Cancelar». |
| `enviando` | «Confirmar» y «Cancelar» deshabilitados; no se envían más pedidos. |

Reglas:

- Pulsar la acción pasa a `confirmando` y **no** envía nada.
- Pulsar «Confirmar» envía un único `PUT /backlog/{id}/estado`.
- «Cancelar» y Escape vuelven a `reposo` sin enviar.
- Un clic fuera de la confirmación no envía nada.
- El foco inicial al confirmar queda en «Cancelar».
- Con `puedeEscribir = false`, la acción queda deshabilitada con `NotaCerrado` y no se abre la
  confirmación.
- Ante un error de la API se cierra la confirmación, se muestra el mensaje tal cual y la acción
  vuelve a quedar disponible.

## Estilos

En `sprints.css`, sin estilos globales nuevos:

- `.sprint-backlog`: lista sin viñetas, en grilla con separación.
- `.historia-tarjeta`: borde, radio y `overflow: hidden`.
- `.historia-tarjeta > summary`: fila con `cursor: pointer`, separación y foco visible.
- `.historia-marca[data-estado="..."]`: colores para «Nueva», «En progreso» y «Completada».
- `.estado-historia__confirmacion` y `.estado-historia__acciones`: pregunta y botones en línea.
