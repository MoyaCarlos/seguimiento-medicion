/** Nota visible que explica por qué una acción de gestión está deshabilitada (se asocia con aria-describedby). */
export function NotaGestion({ id, texto }: { id: string; texto: string }) {
  return (
    <p id={id} className="nota-gestion">
      {texto}
    </p>
  )
}
