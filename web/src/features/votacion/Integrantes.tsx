import { ETIQUETA_ROL, textoDeVoto, type VotoDeIntegrante } from './votacion'

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
