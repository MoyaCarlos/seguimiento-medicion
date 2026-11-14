import { Link } from 'react-router'

export function NoIntegrante() {
  return (
    <section className="placeholder placeholder--feature" aria-labelledby="no-integrante-title">
      <h1 id="no-integrante-title">No sos integrante de este proyecto</h1>
      <p>
        <Link to="/">Mis proyectos</Link>
      </p>
    </section>
  )
}
