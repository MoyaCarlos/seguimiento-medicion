import { textoDeVoto, type Rol, type VotoDeIntegrante } from './votacion'

const ETIQUETA_ROL: Record<Rol, string> = {
  scrum_master: 'Scrum Master',
  product_builder: 'Product Builder',
}

export function Integrantes({ integrantes }: { integrantes: VotoDeIntegrante[] }) {
  return (
    <ul className="votacion-integrantes">
      {integrantes.map((integrante) => (
        <li key={integrante.usuario_id} className="votacion-integrante">
          <span className="votacion-integrante__nombre">{integrante.nombre}</span>
          <span className="votacion-integrante__rol">{ETIQUETA_ROL[integrante.rol]}</span>
          {integrante.ya_voto && <span className="votacion-integrante__estado">Ya votó</span>}
          <span className="votacion-solo-lectores">Voto: </span>
          <strong className="votacion-integrante__voto">{textoDeVoto(integrante)}</strong>
        </li>
      ))}
    </ul>
  )
}
