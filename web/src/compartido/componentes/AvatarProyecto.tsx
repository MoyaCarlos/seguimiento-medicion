type Props = { id: number | string; nombre: string; chico?: boolean }

const PALETAS = 5

export function AvatarProyecto({ id, nombre, chico = false }: Props) {
  return (
    <span
      className={`avatar-proyecto avatar-proyecto--${Number(id) % PALETAS}${chico ? ' avatar-proyecto--chico' : ''}`}
      data-inicial={(nombre.trim().charAt(0) || '?').toUpperCase()}
      aria-hidden="true"
    />
  )
}
