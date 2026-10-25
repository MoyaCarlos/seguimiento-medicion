package features

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// pasosBacklog guarda el estado de los escenarios de HU-15. Reutiliza la base y
// el proyecto reales del Background de HU-01 (no redefine ese paso).
type pasosBacklog struct {
	base      *scenarioContext
	historias *repository.SQLiteBacklogRepository
	consultar *service.ConsultarProductBacklog
	editar    *service.EditarHistoria
	backlog   []domain.BacklogItem
	historia  domain.BacklogItem
	editada   domain.BacklogItem
	enviado   service.EditarHistoriaInput
	err       error
}

func (p *pasosBacklog) preparar() {
	p.historias = repository.NewSQLiteBacklogRepository(p.base.db)
	proyectos := repository.NewSQLiteProjectRepository(p.base.db)
	p.consultar = service.NewConsultarProductBacklog(proyectos, p.historias)
	p.editar = service.NewEditarHistoria(p.historias, repository.NewSQLiteSprintRepository(p.base.db))
	p.backlog, p.historia, p.editada, p.err = nil, domain.BacklogItem{}, domain.BacklogItem{}, nil
	p.enviado = service.EditarHistoriaInput{}
}

func (p *pasosBacklog) crearHistoria(titulo string, prioridad domain.Prioridad) (domain.BacklogItem, error) {
	item, err := domain.NewBacklogItem(p.base.proyectoID, titulo, "Descripción de "+titulo, prioridad, nil)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	return p.historias.Guardar(context.Background(), item)
}

// crearSprint siembra un Sprint directo en la base: la asignación de historias
// a un Sprint llega con HU-06.
func (p *pasosBacklog) crearSprint(estado string) (int64, error) {
	var (
		resultado interface{ LastInsertId() (int64, error) }
		err       error
	)
	if estado == string(domain.SprintPendiente) {
		resultado, err = p.base.db.Exec(`INSERT INTO sprints (proyecto_id, estado) VALUES (?, ?)`, p.base.proyectoID, estado)
	} else {
		resultado, err = p.base.db.Exec(
			`INSERT INTO sprints (proyecto_id, estado, sprint_goal, fecha_inicio, fecha_fin) VALUES (?, ?, 'Objetivo', '2026-10-01', '2026-10-15')`,
			p.base.proyectoID, estado)
	}
	if err != nil {
		return 0, err
	}
	return resultado.LastInsertId()
}

func (p *pasosBacklog) asignarASprint(historiaID, sprintID int64, estimacion *int) error {
	_, err := p.base.db.Exec(
		`UPDATE historias SET sprint_id = ?, estado = 'En progreso', estimacion_sp = ? WHERE id = ?`,
		sprintID, estimacion, historiaID)
	return err
}

func (p *pasosBacklog) proyectoTieneHistorias(tabla *godog.Table) error {
	var sprintID int64
	for _, fila := range tabla.Rows[1:] {
		titulo, prioridad, enSprint := fila.Cells[0].Value, fila.Cells[1].Value, fila.Cells[2].Value
		h, err := p.crearHistoria(titulo, domain.Prioridad(prioridad))
		if err != nil {
			return err
		}
		if enSprint != "si" {
			continue
		}
		if sprintID == 0 {
			if sprintID, err = p.crearSprint(string(domain.SprintActivo)); err != nil {
				return err
			}
		}
		if err := p.asignarASprint(h.ID, sprintID, nil); err != nil {
			return err
		}
	}
	return nil
}

func (p *pasosBacklog) proyectoSinHistoriasSinSprint() error {
	h, err := p.crearHistoria("En un Sprint", domain.PrioridadMust)
	if err != nil {
		return err
	}
	sprintID, err := p.crearSprint(string(domain.SprintActivo))
	if err != nil {
		return err
	}
	return p.asignarASprint(h.ID, sprintID, nil)
}

func (p *pasosBacklog) consultoBacklog() error {
	p.backlog, p.err = p.consultar.Ejecutar(context.Background(), p.base.proyectoID)
	return nil
}

func (p *pasosBacklog) consultoBacklogDeInexistente(identificador int) error {
	p.backlog, p.err = p.consultar.Ejecutar(context.Background(), int64(identificador))
	return nil
}

func (p *pasosBacklog) veoHistoriasEnOrden(tabla *godog.Table) error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error: %v", p.err)
	}
	esperados := make([]string, 0, len(tabla.Rows)-1)
	for _, fila := range tabla.Rows[1:] {
		esperados = append(esperados, fila.Cells[0].Value)
	}
	obtenidos := make([]string, 0, len(p.backlog))
	for _, h := range p.backlog {
		obtenidos = append(obtenidos, h.Titulo)
	}
	if fmt.Sprint(obtenidos) != fmt.Sprint(esperados) {
		return fmt.Errorf("se esperaba %v, se obtuvo %v", esperados, obtenidos)
	}
	return nil
}

func (p *pasosBacklog) veoBacklogVacio() error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error: %v", p.err)
	}
	if p.backlog == nil || len(p.backlog) != 0 {
		return fmt.Errorf("se esperaba un Product Backlog vacío, se obtuvo %#v", p.backlog)
	}
	return nil
}

func (p *pasosBacklog) consultaRechazadaProyectoNoEncontrado() error {
	if !errors.Is(p.err, domain.ErrProyectoNoEncontrado) {
		return fmt.Errorf("se esperaba proyecto no encontrado, se obtuvo %v", p.err)
	}
	return nil
}

func (p *pasosBacklog) existeHistoriaEnBacklog() error {
	h, err := p.crearHistoria("Iniciar sesión", domain.PrioridadMust)
	p.historia = h
	return err
}

func (p *pasosBacklog) existeHistoriaEnSprint(estimacion *int, estado string) error {
	h, err := p.crearHistoria("Historia del Sprint", domain.PrioridadMust)
	if err != nil {
		return err
	}
	sprintID, err := p.crearSprint(estado)
	if err != nil {
		return err
	}
	if err := p.asignarASprint(h.ID, sprintID, estimacion); err != nil {
		return err
	}
	p.historia, err = p.historias.ObtenerPorID(context.Background(), h.ID)
	return err
}

func (p *pasosBacklog) existeHistoriaConSPEnSprint(sp int, estado string) error {
	return p.existeHistoriaEnSprint(&sp, estado)
}

func (p *pasosBacklog) existeHistoriaSinEstimarEnSprint(estado string) error {
	return p.existeHistoriaEnSprint(nil, estado)
}

// guardar edita la historia del escenario partiendo de sus datos actuales: la
// edición envía siempre todos los campos editables.
func (p *pasosBacklog) guardar(cambios func(*service.EditarHistoriaInput)) error {
	input := service.EditarHistoriaInput{
		ID:                  p.historia.ID,
		Titulo:              p.historia.Titulo,
		Descripcion:         p.historia.Descripcion,
		CriteriosAceptacion: p.historia.CriteriosAceptacion,
		EstimacionSP:        p.historia.EstimacionSP,
	}
	cambios(&input)
	p.enviado = input
	p.editada, p.err = p.editar.Ejecutar(context.Background(), input)
	return nil
}

func (p *pasosBacklog) cambioTodo(titulo string, sp int) error {
	return p.guardar(func(in *service.EditarHistoriaInput) {
		criterios := "Dado que tengo cuenta de Google, cuando ingreso, entonces veo el tablero"
		in.Titulo, in.Descripcion, in.CriteriosAceptacion, in.EstimacionSP = titulo, "Autenticación con OAuth", &criterios, &sp
	})
}

func (p *pasosBacklog) cambioTitulo(titulo string) error {
	return p.guardar(func(in *service.EditarHistoriaInput) { in.Titulo = titulo })
}

func (p *pasosBacklog) cambioStoryPoints(sp int) error {
	return p.guardar(func(in *service.EditarHistoriaInput) { in.EstimacionSP = &sp })
}

func (p *pasosBacklog) intentoEditarInexistente(identificador int) error {
	p.historia = domain.BacklogItem{ID: int64(identificador)}
	return p.guardar(func(in *service.EditarHistoriaInput) { in.Titulo = "" })
}

func (p *pasosBacklog) historiaGuardada() (domain.BacklogItem, error) {
	if p.err != nil {
		return domain.BacklogItem{}, fmt.Errorf("no se esperaba error: %v", p.err)
	}
	return p.historias.ObtenerPorID(context.Background(), p.historia.ID)
}

func (p *pasosBacklog) quedaConDatosNuevos() error {
	h, err := p.historiaGuardada()
	if err != nil {
		return err
	}
	if h.Titulo != p.enviado.Titulo || h.Descripcion != p.enviado.Descripcion ||
		textoOpcional(h.CriteriosAceptacion) != textoOpcional(p.enviado.CriteriosAceptacion) ||
		enteroOpcional(h.EstimacionSP) != enteroOpcional(p.enviado.EstimacionSP) {
		return fmt.Errorf("la historia no quedó con los datos enviados: se envió %+v, quedó %+v", p.enviado, h)
	}
	return nil
}

func (p *pasosBacklog) conservaIdentificadorYEstado(estado string) error {
	h, err := p.historiaGuardada()
	if err != nil {
		return err
	}
	if p.editada.ID != p.historia.ID {
		return fmt.Errorf("la edición devolvió el id %d, se esperaba %d", p.editada.ID, p.historia.ID)
	}
	if p.editada.Estado != domain.Estado(estado) || h.Estado != domain.Estado(estado) {
		return fmt.Errorf("se esperaba el estado %q, la edición devolvió %q y quedó %q", estado, p.editada.Estado, h.Estado)
	}
	return nil
}

func (p *pasosBacklog) quedaConStoryPoints(sp int) error {
	h, err := p.historiaGuardada()
	if err != nil {
		return err
	}
	if h.EstimacionSP == nil || *h.EstimacionSP != sp {
		return fmt.Errorf("se esperaban %d Story Points, se obtuvo %v", sp, h.EstimacionSP)
	}
	return nil
}

func (p *pasosBacklog) rechazadaConAdvertenciaSobre(campo string) error {
	var verr domain.ValidationError
	if !errors.As(p.err, &verr) {
		return fmt.Errorf("se esperaba una advertencia de validación, se obtuvo %v", p.err)
	}
	if verr.Campo != campo {
		return fmt.Errorf("se esperaba una advertencia sobre %q, se obtuvo %q", campo, verr.Campo)
	}
	return nil
}

func (p *pasosBacklog) quedaSinModificaciones() error {
	h, err := p.historias.ObtenerPorID(context.Background(), p.historia.ID)
	if err != nil {
		return err
	}
	if h.Titulo != p.historia.Titulo || h.Descripcion != p.historia.Descripcion ||
		textoOpcional(h.CriteriosAceptacion) != textoOpcional(p.historia.CriteriosAceptacion) ||
		enteroOpcional(h.EstimacionSP) != enteroOpcional(p.historia.EstimacionSP) {
		return fmt.Errorf("la historia fue modificada: antes %+v, ahora %+v", p.historia, h)
	}
	return nil
}

func (p *pasosBacklog) rechazadaHistoriaNoEncontrada() error {
	if !errors.Is(p.err, domain.ErrHistoriaNoEncontrada) {
		return fmt.Errorf("se esperaba historia no encontrada, se obtuvo %v", p.err)
	}
	return nil
}

func textoOpcional(v *string) string {
	if v == nil {
		return "<nil>"
	}
	return *v
}

func enteroOpcional(v *int) string {
	if v == nil {
		return "<nil>"
	}
	return strconv.Itoa(*v)
}

// InitializeScenarioBacklog registra los pasos de HU-15 sobre el contexto de
// HU-01, para compartir el proyecto creado por el Antecedentes.
func InitializeScenarioBacklog(ctx *godog.ScenarioContext, base *scenarioContext) {
	p := &pasosBacklog{base: base}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		p.preparar()
		return goctx, nil
	})

	ctx.Step(`^que el proyecto tiene estas historias, en este orden de creación:$`, p.proyectoTieneHistorias)
	ctx.Step(`^que el proyecto no tiene historias sin Sprint asignado$`, p.proyectoSinHistoriasSinSprint)
	ctx.Step(`^consulto el Product Backlog del proyecto$`, p.consultoBacklog)
	ctx.Step(`^consulto el Product Backlog del proyecto con identificador (\d+), que no existe$`, p.consultoBacklogDeInexistente)
	ctx.Step(`^veo las historias en este orden:$`, p.veoHistoriasEnOrden)
	ctx.Step(`^veo un Product Backlog vacío, sin error$`, p.veoBacklogVacio)
	ctx.Step(`^la consulta se rechaza con error "proyecto no encontrado"$`, p.consultaRechazadaProyectoNoEncontrado)
	ctx.Step(`^que existe una historia en el Product Backlog del proyecto$`, p.existeHistoriaEnBacklog)
	ctx.Step(`^que existe una historia con (\d+) Story Points asignada a un Sprint en estado "([^"]*)"$`, p.existeHistoriaConSPEnSprint)
	ctx.Step(`^que existe una historia sin estimar asignada a un Sprint en estado "([^"]*)"$`, p.existeHistoriaSinEstimarEnSprint)
	ctx.Step(`^cambio su título a "([^"]*)", su descripción, sus criterios de aceptación y sus Story Points a (\d+), y guardo$`, p.cambioTodo)
	ctx.Step(`^cambio su título a "([^"]*)" y guardo$`, p.cambioTitulo)
	ctx.Step(`^cambio sus Story Points a (\d+) y guardo$`, p.cambioStoryPoints)
	ctx.Step(`^intento editar la historia con identificador (\d+), que no existe, dejando su título en blanco$`, p.intentoEditarInexistente)
	ctx.Step(`^la historia queda guardada con los datos nuevos$`, p.quedaConDatosNuevos)
	ctx.Step(`^conserva su identificador y su estado "([^"]*)"$`, p.conservaIdentificadorYEstado)
	ctx.Step(`^la historia queda guardada con (\d+) Story Points$`, p.quedaConStoryPoints)
	ctx.Step(`^la edición se rechaza con una advertencia sobre el campo "([^"]*)"$`, p.rechazadaConAdvertenciaSobre)
	ctx.Step(`^la historia queda sin modificaciones$`, p.quedaSinModificaciones)
	ctx.Step(`^la edición se rechaza con error "historia no encontrada"$`, p.rechazadaHistoriaNoEncontrada)
}
