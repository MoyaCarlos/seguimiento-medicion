# Contrato: componente `Markdown`

```tsx
type Props = {
  texto: string
}

export function Markdown({ texto }: Props): JSX.Element
```

- Renderiza `texto` como CommonMark dentro de un contenedor con la clase `markdown`.
- No interpreta HTML crudo.
- No usa `dangerouslySetInnerHTML`.
- Las URLs se filtran con la transformación por defecto de `react-markdown`.
- No recibe ni emite eventos; es un componente de presentación puro.

## Uso

```tsx
<Markdown texto={historia.descripcion} />
<Markdown texto={historia.criterios_aceptacion} />
```

El contenedor padre decide cuándo se muestra el mensaje de criterios vacíos.
