import { ETIQUETA_ROL, textoDeVoto, type VotoDeIntegrante } from './votacion'

type Props = { integrantes: VotoDeIntegrante[]; usuarioActualId?: number | null }

export function Integrantes({ integrantes, usuarioActualId = null }: Props) {
  return (
    <ul className="votacion-integrantes">
      {integrantes.map((integrante) => (
        <li
          key={integrante.usuario_id}
          className={`votacion-integrante${integrante.usuario_id === usuarioActualId ? ' votacion-integrante--yo' : ''}`}
          aria-current={integrante.usuario_id === usuarioActualId ? 'true' : undefined}
        >
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
