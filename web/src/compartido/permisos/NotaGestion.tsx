import { Icono } from '../componentes/Icono'

/** Nota visible que explica por qué una acción de gestión está deshabilitada (se asocia con aria-describedby). */
export function NotaGestion({ id, texto }: { id: string; texto: string }) {
  return (
    <p id={id} className="nota-gestion">
      <Icono nombre="candado" />
      <span>{texto}</span>
    </p>
  )
}
