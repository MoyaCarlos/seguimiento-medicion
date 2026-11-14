import { useParams } from 'react-router'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { useEsScrumMaster } from '../../compartido/permisos/useEsScrumMaster'
import { AccionesScrumMaster } from './AccionesScrumMaster'
import { Consenso } from './Consenso'
import { Encabezado } from './Encabezado'
import { Resultado } from './Resultado'
import { ResumenHistoria } from './ResumenHistoria'
import { Ronda } from './Ronda'
import { SinVotacion } from './SinVotacion'
import { useVotacion } from './useVotacion'
import { estaRevelada } from './votacion'
import './votacion.css'

export function Component() {
  const { id = '', historiaId = '' } = useParams()

  if (!/^[1-9]\d*$/.test(historiaId)) return <ErrorMessage message={`La historia "${historiaId}" no es válida`} />

  return <PantallaVotacion key={`${id}/${historiaId}`} proyectoId={id} historiaId={Number(historiaId)} />
}

type Props = { proyectoId: string; historiaId: number }

function PantallaVotacion({ proyectoId, historiaId }: Props) {
  const estado = useVotacion(proyectoId, historiaId)
  const puedeGestionar = useEsScrumMaster(Number(proyectoId))
  const { integrantes, votacion, enviando } = estado

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

      {estado.sinVotacion && (
        <SinVotacion
          deshabilitado={enviando}
          alIniciar={estado.iniciar}
          puedeGestionar={puedeGestionar}
        />
      )}

      {votacion && (
        <>
          <Ronda
            votacion={votacion}
            deshabilitado={enviando}
            alVotar={estado.votar}
          />
          {estaRevelada(votacion) && (
            <>
              <Resultado promedio={votacion.promedio} diferencias={votacion.diferencias} />
              <Consenso
                key={votacion.ronda}
                votacion={votacion}
                deshabilitado={enviando}
                alRegistrar={estado.registrarEstimacion}
                alAbrirRonda={estado.abrirNuevaRonda}
                puedeGestionar={puedeGestionar}
              />
              {estado.aviso && (
                <p className="message votacion-aviso" role="status">
                  {estado.aviso}
                </p>
              )}
            </>
          )}
          {!estaRevelada(votacion) && (
            <AccionesScrumMaster
              votacion={votacion}
              deshabilitado={enviando}
              alRevelar={estado.revelar}
              puedeGestionar={puedeGestionar}
            />
          )}
        </>
      )}
    </section>
  )
}
