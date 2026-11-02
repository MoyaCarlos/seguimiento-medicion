import { useParams } from 'react-router'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { AccionesScrumMaster } from './AccionesScrumMaster'
import { Encabezado } from './Encabezado'
import { Resultado } from './Resultado'
import { ResumenHistoria } from './ResumenHistoria'
import { Ronda } from './Ronda'
import { SelectorActor } from './SelectorActor'
import { SinVotacion } from './SinVotacion'
import { useVotacion } from './useVotacion'
import { estaRevelada } from './votacion'
import './votacion.css'

export function Component() {
  const { id = '', historiaId = '' } = useParams()

  return <PantallaVotacion key={`${id}/${historiaId}`} proyectoId={id} historiaId={Number(historiaId)} />
}

type Props = { proyectoId: string; historiaId: number }

function PantallaVotacion({ proyectoId, historiaId }: Props) {
  const estado = useVotacion(proyectoId, historiaId)
  const { integrantes, actor, scrumMaster, votacion, enviando } = estado

  return (
    <section className="votacion" aria-labelledby="votacion-title">
      <Encabezado
        proyectoId={proyectoId}
        historiaId={historiaId}
        deshabilitado={enviando || estado.cargando}
        alActualizar={estado.recargar}
      />

      <ResumenHistoria historiaId={historiaId} historia={estado.historia} />

      {estado.errorDeIntegrantes && <ErrorMessage message={estado.errorDeIntegrantes} />}
      {estado.errorDeVotacion && <ErrorMessage message={estado.errorDeVotacion} />}
      {estado.errorDeAccion && <ErrorMessage message={estado.errorDeAccion} />}
      {estado.cargando && <LoadingIndicator />}

      {integrantes !== null && integrantes.length === 0 && <p>El proyecto no tiene integrantes</p>}
      {integrantes !== null && integrantes.length > 0 && (
        <SelectorActor
          integrantes={integrantes}
          actorId={actor?.id}
          deshabilitado={enviando}
          alElegir={estado.elegirActor}
        />
      )}

      {estado.sinVotacion && (
        <SinVotacion
          actor={actor}
          puedeIniciar={scrumMaster !== undefined}
          hayIntegrantes={Boolean(integrantes?.length)}
          deshabilitado={enviando}
          alIniciar={estado.iniciar}
        />
      )}

      {votacion && (
        <>
          <Ronda votacion={votacion} actor={actor} deshabilitado={enviando} alVotar={estado.votar} />
          {estaRevelada(votacion) && (
            <>
              <Resultado promedio={votacion.promedio} />
              <p>Registrar la estimación acordada se hace en otra historia (HU-14).</p>
            </>
          )}
          {scrumMaster && !estaRevelada(votacion) && (
            <AccionesScrumMaster
              votacion={votacion}
              scrumMaster={scrumMaster}
              deshabilitado={enviando}
              alRevelar={estado.revelar}
            />
          )}
        </>
      )}
    </section>
  )
}
