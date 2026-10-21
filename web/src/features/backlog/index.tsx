import { useParams } from 'react-router'

export function Component() {
  const { id } = useParams()

  return (
    <section className="placeholder placeholder--feature" aria-labelledby="placeholder-title">
      <div className="placeholder-icon" aria-hidden="true">✦</div>
      <div>
        <p className="eyebrow">Base del frontend · Sprint 2</p>
        <h1 id="placeholder-title">Product Backlog</h1>
        <p className="project-context">Proyecto actual: <strong>{id}</strong></p>
      </div>
      <p>Esta pantalla será implementada por su historia de usuario.</p>
    </section>
  )
}
