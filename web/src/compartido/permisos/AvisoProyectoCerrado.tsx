export function AvisoProyectoCerrado({ fechaCierre }: { fechaCierre: string | null }) {
  if (fechaCierre === null) return null
  return (
    <div role="status" className="aviso-cerrado">
      Proyecto cerrado el {fechaCierre}
    </div>
  )
}
