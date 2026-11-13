export function AvisoSinPermiso({ mensaje }: { mensaje: string }) {
  return (
    <div className="message message--error" role="alert">
      <strong>No tenés permiso</strong>
      <p>{mensaje}</p>
    </div>
  )
}
