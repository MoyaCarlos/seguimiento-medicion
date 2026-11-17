package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// lecturaDe es un servicio de lectura con alcance de proyecto, armado sobre el
// proyecto 1 (equipo: Scrum Master y Product Builder) con una historia (3) y un
// Sprint (1) suyos. recurso es el id que existe; inexistente, uno que no.
type lecturaDe struct {
	nombre string
	leer   func(consultante int64) error
	// inexistente lee un recurso que no existe como consultante.
	inexistente func(consultante int64) error
}

func lecturasDelProyecto(t *testing.T, proyectos *fakeProjectRepository) []lecturaDe {
	t.Helper()
	ctx := context.Background()
	sprints := repository.NewSprintRepositoryEnMemoria()
	finalizado := sprintDeReporte(t, sprints, proyectoPermisos, "Objetivo", domain.SprintFinalizado)
	sprintsFinalizados := &fakeListadorSprints{sprints: []domain.Sprint{finalizado}}
	obtenedorSprint := &fakeObtenedorSprint{sprints: map[int64]domain.Sprint{finalizado.ID: finalizado}}
	metricas := NewCalcularMetricas(proyectos, &fakeDatosMetricas{datos: []domain.DatosMetricasSprint{{SprintID: finalizado.ID, Numero: 1}}})
	reporte := nuevoGenerarReporteDePrueba(proyectos, sprints, &fakeProductBacklogRepository{}, &fakeHistoriasDeSprintRepository{},
		&fakeArrastradasRepository{}, &fakeDatosMetricasRepository{}, &fakeDefectoRepository{})
	historiasDelProyecto := historias(historiaDelBacklog())
	votaciones := votacionAbierta()

	soloError := func(_ any, err error) error { return err }
	return []lecturaDe{
		{"ObtenerEstadoProyecto",
			func(u int64) error {
				return soloError(NewObtenerEstadoProyecto(proyectos, sprints).Ejecutar(ctx, proyectoPermisos, u, time.Now()))
			},
			func(u int64) error {
				return soloError(NewObtenerEstadoProyecto(proyectos, sprints).Ejecutar(ctx, 999, u, time.Now()))
			}},
		{"ListarIntegrantes",
			func(u int64) error {
				return soloError(NewListarIntegrantes(proyectos).Ejecutar(ctx, proyectoPermisos, u))
			},
			func(u int64) error { return soloError(NewListarIntegrantes(proyectos).Ejecutar(ctx, 999, u)) }},
		{"ConsultarProductBacklog",
			func(u int64) error {
				return soloError(NewConsultarProductBacklog(proyectos, &fakeProductBacklogRepository{}).Ejecutar(ctx, proyectoPermisos, u))
			},
			func(u int64) error {
				return soloError(NewConsultarProductBacklog(proyectos, &fakeProductBacklogRepository{}).Ejecutar(ctx, 999, u))
			}},
		{"ObtenerHistoria.ComoIntegrante",
			func(u int64) error {
				return soloError(NewObtenerHistoria(historiasDelProyecto, proyectos).ComoIntegrante(ctx, 3, u))
			},
			func(u int64) error {
				return soloError(NewObtenerHistoria(historiasDelProyecto, proyectos).ComoIntegrante(ctx, 999, u))
			}},
		{"ConsultarEsfuerzoReal",
			func(u int64) error {
				return soloError(NewConsultarEsfuerzoReal(historiasDelProyecto, &fakeRegistroHorasRepository{}, proyectos).Ejecutar(ctx, 3, u))
			},
			func(u int64) error {
				return soloError(NewConsultarEsfuerzoReal(historiasDelProyecto, &fakeRegistroHorasRepository{}, proyectos).Ejecutar(ctx, 999, u))
			}},
		{"ConsultarHistorialRondas",
			func(u int64) error {
				return soloError(NewConsultarHistorialRondas(historiasDelProyecto, proyectos, votaciones).Ejecutar(ctx, 3, u))
			},
			func(u int64) error {
				return soloError(NewConsultarHistorialRondas(historiasDelProyecto, proyectos, votaciones).Ejecutar(ctx, 999, u))
			}},
		{"ConsultarVotacion",
			func(u int64) error {
				return soloError(NewConsultarVotacion(historiasDelProyecto, proyectos, votaciones, domain.CoincidenciaExacta{}).Ejecutar(ctx, 3, u))
			},
			func(u int64) error {
				return soloError(NewConsultarVotacion(historiasDelProyecto, proyectos, votaciones, domain.CoincidenciaExacta{}).Ejecutar(ctx, 999, u))
			}},
		{"ConsultarSprintBacklog",
			func(u int64) error {
				return soloError(NewConsultarSprintBacklog(sprints, &fakeHistoriasDeSprintRepository{}, proyectos).Ejecutar(ctx, finalizado.ID, u))
			},
			func(u int64) error {
				return soloError(NewConsultarSprintBacklog(sprints, &fakeHistoriasDeSprintRepository{}, proyectos).Ejecutar(ctx, 999, u))
			}},
		{"ConsultarSprints",
			func(u int64) error {
				return soloError(NewConsultarSprints(proyectos, sprintsFinalizados).Ejecutar(ctx, proyectoPermisos, u))
			},
			func(u int64) error {
				return soloError(NewConsultarSprints(proyectos, sprintsFinalizados).Ejecutar(ctx, 999, u))
			}},
		{"ListarDefectos",
			func(u int64) error {
				return soloError(NewListarDefectos(proyectos, &fakeDefectoRepository{}).Ejecutar(ctx, proyectoPermisos, u))
			},
			func(u int64) error {
				return soloError(NewListarDefectos(proyectos, &fakeDefectoRepository{}).Ejecutar(ctx, 999, u))
			}},
		{"CalcularMetricas",
			func(u int64) error { return soloError(metricas.Ejecutar(ctx, proyectoPermisos, u)) },
			func(u int64) error { return soloError(metricas.Ejecutar(ctx, 999, u)) }},
		{"ConsultarHistorialSprints",
			func(u int64) error {
				return soloError(NewConsultarHistorialSprints(metricas, sprintsFinalizados).Ejecutar(ctx, proyectoPermisos, u))
			},
			func(u int64) error {
				return soloError(NewConsultarHistorialSprints(metricas, sprintsFinalizados).Ejecutar(ctx, 999, u))
			}},
		{"ConsultarDetalleSprint",
			func(u int64) error {
				return soloError(NewConsultarDetalleSprint(proyectos, obtenedorSprint, &fakeCalculadorMetricas{resultado: domain.MetricasProyecto{
					Sprints: []domain.MetricasDeSprint{metricasDeSprint(finalizado.ID, 1, 5)}}}, &fakeHistorialSprintRepository{}).Ejecutar(ctx, proyectoPermisos, finalizado.ID, u))
			},
			func(u int64) error {
				return soloError(NewConsultarDetalleSprint(proyectos, obtenedorSprint, &fakeCalculadorMetricas{}, &fakeHistorialSprintRepository{}).Ejecutar(ctx, 999, finalizado.ID, u))
			}},
		{"ConsultarAvanceHistorias",
			func(u int64) error {
				return soloError(NewConsultarAvanceHistorias(proyectos, &fakeAvanceHistoriasRepository{}, &fakeListadorSprints{}).Ejecutar(ctx, proyectoPermisos, u))
			},
			func(u int64) error {
				return soloError(NewConsultarAvanceHistorias(proyectos, &fakeAvanceHistoriasRepository{}, &fakeListadorSprints{}).Ejecutar(ctx, 999, u))
			}},
		{"GenerarReporte.EjecutarSprint",
			func(u int64) error {
				return soloError(reporte.EjecutarSprint(ctx, finalizado.ID, u, fechaDeReporte(20)))
			},
			func(u int64) error { return soloError(reporte.EjecutarSprint(ctx, 999, u, fechaDeReporte(20))) }},
		{"GenerarReporte.EjecutarProyecto",
			func(u int64) error {
				return soloError(reporte.EjecutarProyecto(ctx, proyectoPermisos, u, fechaDeReporte(20)))
			},
			func(u int64) error { return soloError(reporte.EjecutarProyecto(ctx, 999, u, fechaDeReporte(20))) }},
	}
}

// Toda lectura con alcance de proyecto exige integrarlo: quien no lo integra
// (o no tiene sesión, id 0) recibe ErrNoEsIntegrante; el Scrum Master y el
// Product Builder leen, también con el proyecto cerrado; el recurso inexistente
// se informa antes (404 antes que 403).
func TestLectura_SoloLosIntegrantesLeen(t *testing.T) {
	escenarios := map[string]func() *fakeProjectRepository{"abierto": proyectoConEquipoGuardado, "cerrado": func() *fakeProjectRepository {
		p := proyectoConEquipoGuardado()
		p.cerrados = map[int64]bool{proyectoPermisos: true}
		return p
	}}
	for estado, proyecto := range escenarios {
		for _, l := range lecturasDelProyecto(t, proyecto()) {
			t.Run(estado+"/"+l.nombre, func(t *testing.T) {
				for _, ajeno := range []int64{ajenoAlProyectoID, 0} {
					if err := l.leer(ajeno); !errors.Is(err, domain.ErrNoEsIntegrante) {
						t.Errorf("consultante %d: se esperaba ErrNoEsIntegrante, se obtuvo %v", ajeno, err)
					}
				}
				for _, integrante := range []int64{scrumMasterID, productBuilderID} {
					if err := l.leer(integrante); err != nil {
						t.Errorf("integrante %d: no se esperaba error, se obtuvo %v", integrante, err)
					}
				}
				err := l.inexistente(ajenoAlProyectoID)
				if err == nil || errors.Is(err, domain.ErrNoEsIntegrante) {
					t.Errorf("recurso inexistente: se esperaba un 404 antes que el 403, se obtuvo %v", err)
				}
			})
		}
	}
}
