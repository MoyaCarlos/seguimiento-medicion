package features

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

var (
	inicioSprint = time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	finSprint    = inicioSprint.AddDate(0, 0, 14)
)

// pasosSprint guarda el estado de los escenarios de HU-05. Comparte la base
// y el proyecto del Background con los pasos de HU-01.
type pasosSprint struct {
	base       *scenarioContext
	sprints    *repository.SQLiteSprintRepository
	backlog    *repository.SQLiteBacklogRepository
	crear      *service.CrearSprint
	iniciar    *service.IniciarSprint
	cerrar     *service.CerrarSprint
	sprint     domain.Sprint
	segundo    domain.Sprint
	completada int64
	pendiente  int64
	estimada   int64
	sinEstimar int64
	err        error
}

func (p *pasosSprint) preparar() {
	p.sprints = repository.NewSQLiteSprintRepository(p.base.db)
	p.backlog = repository.NewSQLiteBacklogRepository(p.base.db)
	p.crear = service.NewCrearSprint(repository.NewSQLiteProjectRepository(p.base.db), p.sprints)
	p.iniciar = service.NewIniciarSprint(p.sprints)
	p.cerrar = service.NewCerrarSprint(p.sprints, p.backlog)
	p.sprint, p.segundo, p.err = domain.Sprint{}, domain.Sprint{}, nil
	p.completada, p.pendiente, p.estimada, p.sinEstimar = 0, 0, 0, 0
}

func (p *pasosSprint) nuevoPendiente() (domain.Sprint, error) {
	s, err := domain.NewSprint(p.base.proyectoID)
	if err != nil {
		return domain.Sprint{}, err
	}
	return p.sprints.Guardar(context.Background(), s)
}

func (p *pasosSprint) iniciarCon(id int64, goal string, fin time.Time) (domain.Sprint, error) {
	return p.iniciar.Ejecutar(context.Background(), service.IniciarSprintInput{
		SprintID: id, SprintGoal: goal, FechaInicio: inicioSprint, FechaFin: fin,
	})
}

func (p *pasosSprint) existePendiente() (err error) {
	p.sprint, err = p.nuevoPendiente()
	return err
}

func (p *pasosSprint) existeActivo() error {
	if err := p.existePendiente(); err != nil {
		return err
	}
	activo, err := p.iniciarCon(p.sprint.ID, "Goal del Sprint", finSprint)
	p.sprint = activo
	return err
}

func (p *pasosSprint) existeOtroPendiente() (err error) {
	p.segundo, err = p.nuevoPendiente()
	return err
}

func (p *pasosSprint) creoUnSprint() error {
	p.sprint, p.err = p.crear.Ejecutar(context.Background(), p.base.proyectoID)
	return nil
}

func (p *pasosSprint) intentoCrearParaProyecto(proyectoID int) error {
	_, p.err = p.crear.Ejecutar(context.Background(), int64(proyectoID))
	return nil
}

func (p *pasosSprint) intentoCrearOtro() error {
	_, p.err = p.crear.Ejecutar(context.Background(), p.base.proyectoID)
	return nil
}

func (p *pasosSprint) unSoloPendiente() error {
	delProyecto, err := p.sprints.ListarPorProyecto(context.Background(), p.base.proyectoID)
	if err != nil {
		return err
	}
	pendientes := 0
	for _, s := range delProyecto {
		if s.Estado == domain.SprintPendiente {
			pendientes++
		}
	}
	if pendientes != 1 {
		return fmt.Errorf("se esperaba un solo Sprint Pendiente, hay %d", pendientes)
	}
	return nil
}

func (p *pasosSprint) inicioConGoal(goal string) error {
	iniciado, err := p.iniciarCon(p.sprint.ID, goal, finSprint)
	if err == nil {
		p.sprint = iniciado
	}
	p.err = err
	return nil
}

func (p *pasosSprint) intentoIniciarSegundo() error {
	_, p.err = p.iniciarCon(p.segundo.ID, "Goal del segundo", finSprint)
	return nil
}

func (p *pasosSprint) intentoIniciarConFechaInvalida() error {
	_, p.err = p.iniciarCon(p.sprint.ID, "Goal", inicioSprint)
	return nil
}

func (p *pasosSprint) cierroElSprint() error {
	cerrado, err := p.cerrar.Ejecutar(context.Background(), p.sprint.ID)
	if err == nil {
		p.sprint = cerrado
	}
	p.err = err
	return nil
}

func (p *pasosSprint) historiasAsignadas() error {
	var ids []int64
	for _, titulo := range []string{"Historia completada", "Historia sin terminar"} {
		item, err := domain.NewBacklogItem(p.base.proyectoID, titulo, "Descripción", domain.PrioridadMust, nil)
		if err != nil {
			return err
		}
		g, err := p.backlog.Guardar(context.Background(), item)
		if err != nil {
			return err
		}
		ids = append(ids, g.ID)
	}
	// La asignación de historias a un Sprint es de HU-06; acá se simula con SQL.
	if _, err := p.base.db.Exec("UPDATE historias SET sprint_id = ? WHERE id IN (?, ?)", p.sprint.ID, ids[0], ids[1]); err != nil {
		return err
	}
	if _, err := p.base.db.Exec("UPDATE historias SET estado = ?, fecha_finalizacion = ? WHERE id = ?", string(domain.EstadoCompletada), "2026-10-20", ids[0]); err != nil {
		return err
	}
	if _, err := p.base.db.Exec("UPDATE historias SET estado = ? WHERE id = ?", "En progreso", ids[1]); err != nil {
		return err
	}
	p.completada, p.pendiente = ids[0], ids[1]
	return nil
}

func (p *pasosSprint) historiaNoCompletadaEnElSprint(titulo string, estimacion *int) (int64, error) {
	item, err := domain.NewBacklogItem(p.base.proyectoID, titulo, "Descripción", domain.PrioridadMust, nil)
	if err != nil {
		return 0, err
	}
	item.EstimacionSP = estimacion
	g, err := p.backlog.Guardar(context.Background(), item)
	if err != nil {
		return 0, err
	}
	// La asignación de historias a un Sprint es de HU-06; acá se simula con SQL.
	_, err = p.base.db.Exec("UPDATE historias SET sprint_id = ?, estado = ? WHERE id = ?", p.sprint.ID, "En progreso", g.ID)
	return g.ID, err
}

func (p *pasosSprint) historiaEstimadaYOtraSinEstimar(puntos int) (err error) {
	if p.estimada, err = p.historiaNoCompletadaEnElSprint("Historia estimada", &puntos); err != nil {
		return err
	}
	p.sinEstimar, err = p.historiaNoCompletadaEnElSprint("Historia sin estimar", nil)
	return err
}

// La edición del Product Backlog es de HU-15; acá se simula con SQL.
func (p *pasosSprint) cambioLaEstimacion(puntos int) error {
	_, err := p.base.db.Exec("UPDATE historias SET estimacion_sp = ? WHERE id = ?", puntos, p.estimada)
	return err
}

func (p *pasosSprint) estimacionArrastrada(historiaID int64) (sql.NullInt64, error) {
	if p.err != nil {
		return sql.NullInt64{}, fmt.Errorf("no se esperaba error al cerrar: %w", p.err)
	}
	var estimacion sql.NullInt64
	err := p.base.db.QueryRow(
		"SELECT estimacion_sp FROM historias_arrastradas WHERE sprint_id = ? AND historia_id = ?", p.sprint.ID, historiaID,
	).Scan(&estimacion)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.NullInt64{}, fmt.Errorf("la historia %d no quedó registrada como arrastrada del Sprint %d", historiaID, p.sprint.ID)
	}
	return estimacion, err
}

func (p *pasosSprint) estimadaArrastradaCon(puntos int) error {
	estimacion, err := p.estimacionArrastrada(p.estimada)
	if err != nil {
		return err
	}
	if !estimacion.Valid || estimacion.Int64 != int64(puntos) {
		return fmt.Errorf("se esperaban %d Story Points registrados, quedó %v", puntos, estimacion)
	}
	return nil
}

func (p *pasosSprint) sinEstimarArrastradaSinPuntos() error {
	estimacion, err := p.estimacionArrastrada(p.sinEstimar)
	if err != nil {
		return err
	}
	if estimacion.Valid {
		return fmt.Errorf("se esperaba la historia registrada sin Story Points, quedó con %d", estimacion.Int64)
	}
	return nil
}

func (p *pasosSprint) estadoPersistido(id int64) (domain.EstadoSprint, error) {
	s, err := p.sprints.ObtenerPorID(context.Background(), id)
	return s.Estado, err
}

func (p *pasosSprint) sprintPasaAEstado(estado string) error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error: %w", p.err)
	}
	actual, err := p.estadoPersistido(p.sprint.ID)
	if err != nil {
		return err
	}
	if string(actual) != estado {
		return fmt.Errorf("se esperaba estado %q, quedó %q", estado, actual)
	}
	return nil
}

func (p *pasosSprint) rechazaLaOperacion() error {
	if p.err == nil {
		return errors.New("se esperaba que el sistema rechace la operación")
	}
	return nil
}

func (p *pasosSprint) rechazaPorFechas() error {
	var verr domain.ValidationError
	if !errors.As(p.err, &verr) || verr.Campo != "fecha_fin" {
		return fmt.Errorf("se esperaba un error de validación de fechas, se obtuvo %v", p.err)
	}
	return nil
}

func (p *pasosSprint) segundoSiguePendiente() error {
	actual, err := p.estadoPersistido(p.segundo.ID)
	if err != nil {
		return err
	}
	if actual != domain.SprintPendiente {
		return fmt.Errorf("el segundo Sprint debía seguir Pendiente, quedó %q", actual)
	}
	return nil
}

func (p *pasosSprint) sprintDeHistoria(id int64) (sql.NullInt64, error) {
	var sprintID sql.NullInt64
	err := p.base.db.QueryRow("SELECT sprint_id FROM historias WHERE id = ?", id).Scan(&sprintID)
	return sprintID, err
}

func (p *pasosSprint) noCompletadaVuelveAlBacklog() error {
	sprintID, err := p.sprintDeHistoria(p.pendiente)
	if err != nil {
		return err
	}
	if sprintID.Valid {
		return fmt.Errorf("la historia no completada sigue en el Sprint %d", sprintID.Int64)
	}
	return nil
}

func (p *pasosSprint) completadaSigueVinculada() error {
	sprintID, err := p.sprintDeHistoria(p.completada)
	if err != nil {
		return err
	}
	if !sprintID.Valid || sprintID.Int64 != p.sprint.ID {
		return errors.New("la historia completada debía seguir vinculada al Sprint")
	}
	return nil
}

func (p *pasosSprint) noCompletadaVuelveANueva() error {
	var estado string
	if err := p.base.db.QueryRow("SELECT estado FROM historias WHERE id = ?", p.pendiente).Scan(&estado); err != nil {
		return err
	}
	if estado != string(domain.EstadoNueva) {
		return fmt.Errorf("la historia no completada debía quedar en estado Nueva, quedó %q", estado)
	}
	return nil
}

// inicializarPasosSprint registra los pasos de HU-05.
func inicializarPasosSprint(ctx *godog.ScenarioContext, base *scenarioContext) {
	p := &pasosSprint{base: base}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		p.preparar()
		return goctx, nil
	})

	ctx.Step(`^creo un Sprint para el proyecto$`, p.creoUnSprint)
	ctx.Step(`^intento crear un Sprint para el proyecto con identificador (\d+), que no existe$`, p.intentoCrearParaProyecto)
	ctx.Step(`^intento crear otro Sprint para el proyecto$`, p.intentoCrearOtro)
	ctx.Step(`^el Sprint queda persistido en estado "([^"]*)"$`, p.sprintPasaAEstado)
	ctx.Step(`^el proyecto sigue teniendo un solo Sprint en estado "Pendiente"$`, p.unSoloPendiente)
	ctx.Step(`^que existe un Sprint en estado "Pendiente" para el proyecto, sin otro Sprint "Activo"$`, p.existePendiente)
	ctx.Step(`^que existe un Sprint en estado "Pendiente" para el proyecto$`, p.existePendiente)
	ctx.Step(`^que el proyecto (?:ya )?tiene un Sprint en estado "Activo"$`, p.existeActivo)
	ctx.Step(`^existe otro Sprint en estado "Pendiente" para el mismo proyecto$`, p.existeOtroPendiente)
	ctx.Step(`^ese Sprint tiene una historia completada y otra no completada asignadas$`, p.historiasAsignadas)
	ctx.Step(`^defino el Sprint Goal "([^"]*)", una fecha de inicio y una fecha de fin posterior, y lo inicio$`, p.inicioConGoal)
	ctx.Step(`^intento iniciar ese segundo Sprint$`, p.intentoIniciarSegundo)
	ctx.Step(`^intento iniciarlo con una fecha de fin igual o anterior a la fecha de inicio$`, p.intentoIniciarConFechaInvalida)
	ctx.Step(`^el Scrum Master cierra el Sprint$`, p.cierroElSprint)
	ctx.Step(`^intento cerrarlo$`, p.cierroElSprint)
	ctx.Step(`^el Sprint pasa al estado "([^"]*)"$`, p.sprintPasaAEstado)
	ctx.Step(`^el sistema rechaza la operación$`, p.rechazaLaOperacion)
	ctx.Step(`^el sistema rechaza la operación por validación de fechas$`, p.rechazaPorFechas)
	ctx.Step(`^el segundo Sprint permanece en estado "Pendiente"$`, p.segundoSiguePendiente)
	ctx.Step(`^la historia no completada queda sin Sprint asignado en el Product Backlog$`, p.noCompletadaVuelveAlBacklog)
	ctx.Step(`^la historia no completada queda en estado "Nueva"$`, p.noCompletadaVuelveANueva)
	ctx.Step(`^la historia completada permanece vinculada a ese Sprint$`, p.completadaSigueVinculada)
	ctx.Step(`^ese Sprint tiene una historia no completada estimada en (\d+) Story Points y otra no completada sin estimar$`, p.historiaEstimadaYOtraSinEstimar)
	ctx.Step(`^se cambia la estimación de la historia estimada a (\d+) Story Points en el Product Backlog$`, p.cambioLaEstimacion)
	ctx.Step(`^la historia estimada queda registrada como arrastrada de ese Sprint con (\d+) Story Points$`, p.estimadaArrastradaCon)
	ctx.Step(`^la historia sin estimar queda registrada como arrastrada de ese Sprint sin Story Points$`, p.sinEstimarArrastradaSinPuntos)
}
