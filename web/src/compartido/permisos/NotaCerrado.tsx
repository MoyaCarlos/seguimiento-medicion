/** Nota visible que explica que una acción está deshabilitada porque el proyecto está cerrado (se asocia con aria-describedby). */
export function NotaCerrado({ id }: { id: string }) {
  return (
    <p id={id} className="nota-gestion">
      El proyecto está cerrado
    </p>
  )
}
