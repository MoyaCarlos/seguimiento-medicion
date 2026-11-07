package features

import (
	"context"
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	apihttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type historiaAvanceJSON struct {
	ID           int64  `json:"id"`
	Titulo       string `json:"titulo"`
	EstimacionSP *int   `json:"estimacion_sp"`
	Sprint       *struct {
		ID     int64 `json:"id"`
		Numero int   `json:"numero"`
	} `json:"sprint"`
}

type avanceJSON struct {
	Total struct {
		Historias   int `json:"historias"`
		StoryPoints int `json:"story_points"`
	} `json:"total"`
	Grupos []struct {
		Estado      string               `json:"estado"`
		Cantidad    int                  `json:"cantidad"`
		StoryPoints int                  `json:"story_points"`
		Historias   []historiaAvanceJSON `json:"historias"`
	} `json:"grupos"`
}

// pasosAvance guarda el estado de los escenarios de HU-18. Reutiliza la base y
// el proyecto del Background, y los services de Sprint e historias de HU-05,
// HU-06 y HU-11 a través de pasosAsignacion; la consulta pasa por la ruta real.
type pasosAvance struct {
	base    *scenarioContext
	plan    *pasosAsignacion
	cambiar *service.CambiarEstadoHistoria
	rutas   nethttp.Handler

	sprints   map[string]int64
	respuesta *httptest.ResponseRecorder
}

func (p *pasosAvance) preparar() {
	p.plan = &pasosAsignacion{base: p.base}
	p.plan.preparar()
	p.cambiar = service.NewCambiarEstadoHistoria(p.plan.backlog, p.plan.sprints, func() time.Time { return relojEstadoHistoria })
	handler := apihttp.NewAvanceHandler(service.NewConsultarAvanceHistorias(
		repository.NewSQLiteProjectRepository(p.base.db), p.plan.backlog, p.plan.sprints))
	mux := nethttp.NewServeMux()
	mux.HandleFunc("GET /projects/{id}/avance", handler.Consultar)
	p.rutas = mux
	p.sprints = map[string]int64{}
	p.respuesta = nil
}

func (p *pasosAvance) proyecto() error { return p.base.existeProyecto(1) }

// sprintEnEstado crea el Sprint la primera vez que se lo nombra, con los
// services de HU-05, y lo reutiliza en los pasos siguientes.
func (p *pasosAvance) sprintEnEstado(nombre, estado string) (int64, error) {
	if id, existe := p.sprints[nombre]; existe {
		return id, nil
	}
	id, err := p.plan.crearSprintEnEstado(estado)
	if err != nil {
		return 0, err
	}
	p.sprints[nombre] = id
	return id, nil
}

func (p *pasosAvance) historiaNueva(titulo string, puntos int) error {
	_, err := p.plan.nuevaHistoria(titulo, &puntos)
	return err
}

func (p *pasosAvance) historiaNuevaSinEstimar(titulo string) error {
	_, err := p.plan.nuevaHistoria(titulo, nil)
	return err
}

func (p *pasosAvance) historiaNuevaEnSprintPendiente(titulo string, puntos int, sprint string) error {
	return p.historiaEnSprint(titulo, puntos, sprint, "Pendiente", domain.EstadoNueva)
}

func (p *pasosAvance) historiaEnSprintActivo(titulo string, puntos int, estado, sprint string) error {
	return p.historiaEnSprint(titulo, puntos, sprint, "Activo", domain.Estado(estado))
}

// historiaEnSprint asigna la historia al Sprint y la lleva al estado pedido
// con los services de HU-06 y HU-11.
func (p *pasosAvance) historiaEnSprint(titulo string, puntos int, sprint, estadoSprint string, estado domain.Estado) error {
	ctx := context.Background()
	historiaID, err := p.plan.nuevaHistoria(titulo, &puntos)
	if err != nil {
		return err
	}
	sprintID, err := p.sprintEnEstado(sprint, estadoSprint)
	if err != nil {
		return err
	}
	if _, err := p.plan.asignar.Ejecutar(ctx, historiaID, sprintID); err != nil {
		return err
	}
	recorrido := map[domain.Estado][]domain.Estado{
		domain.EstadoNueva:      nil,
		domain.EstadoEnProgreso: {domain.EstadoEnProgreso},
		domain.EstadoCompletada: {domain.EstadoEnProgreso, domain.EstadoCompletada},
	}
	for _, siguiente := range recorrido[estado] {
		if _, err := p.cambiar.Ejecutar(ctx, historiaID, siguiente); err != nil {
			return err
		}
	}
	return nil
}

func (p *pasosAvance) seCierraElSprint(nombre string) error {
	id, ok := p.sprints[nombre]
	if !ok {
		return fmt.Errorf("el Sprint %q no fue creado en el escenario", nombre)
	}
	_, err := p.plan.cerrar.Ejecutar(context.Background(), id)
	return err
}

func (p *pasosAvance) consultarCon(id string) error {
	req := httptest.NewRequest(nethttp.MethodGet, "/projects/"+id+"/avance", nil)
	p.respuesta = httptest.NewRecorder()
	p.rutas.ServeHTTP(p.respuesta, req)
	return nil
}

func (p *pasosAvance) consultar() error {
	return p.consultarCon(strconv.FormatInt(p.base.proyectoID, 10))
}

func (p *pasosAvance) consultarProyecto(id int) error { return p.consultarCon(strconv.Itoa(id)) }

func (p *pasosAvance) consultarTexto(id string) error { return p.consultarCon(id) }

func (p *pasosAvance) avance() (avanceJSON, error) {
	var a avanceJSON
	if p.respuesta == nil {
		return a, fmt.Errorf("todavía no se consultó el avance")
	}
	if p.respuesta.Code != nethttp.StatusOK {
		return a, fmt.Errorf("se esperaba 200, se obtuvo %d (%s)", p.respuesta.Code, p.respuesta.Body.String())
	}
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &a); err != nil {
		return a, fmt.Errorf("la respuesta no es un avance válido: %w", err)
	}
	return a, nil
}

func (p *pasosAvance) informaGrupos(tabla *godog.Table) error {
	a, err := p.avance()
	if err != nil {
		return err
	}
	filas := tabla.Rows[1:]
	if len(a.Grupos) != len(filas) {
		return fmt.Errorf("se esperaban %d grupos, se obtuvieron %d", len(filas), len(a.Grupos))
	}
	for i, fila := range filas {
		esperado := [3]string{fila.Cells[0].Value, fila.Cells[1].Value, fila.Cells[2].Value}
		g := a.Grupos[i]
		obtenido := [3]string{g.Estado, strconv.Itoa(g.Cantidad), strconv.Itoa(g.StoryPoints)}
		if esperado != obtenido {
			return fmt.Errorf("grupo %d: se esperaba %v, se obtuvo %v", i, esperado, obtenido)
		}
		if len(g.Historias) != g.Cantidad {
			return fmt.Errorf("grupo %q: informa %d historias pero lista %d", g.Estado, g.Cantidad, len(g.Historias))
		}
	}
	return nil
}

func (p *pasosAvance) informaTotal(historias, puntos int) error {
	a, err := p.avance()
	if err != nil {
		return err
	}
	if a.Total.Historias != historias || a.Total.StoryPoints != puntos {
		return fmt.Errorf("se esperaban %d historias y %d Story Points en total, se obtuvo %d y %d",
			historias, puntos, a.Total.Historias, a.Total.StoryPoints)
	}
	return nil
}

func (p *pasosAvance) grupoInforma(estado string, historias, puntos int) error {
	a, err := p.avance()
	if err != nil {
		return err
	}
	for _, g := range a.Grupos {
		if g.Estado == estado {
			if g.Cantidad != historias || g.StoryPoints != puntos {
				return fmt.Errorf("grupo %q: se esperaban %d historias y %d Story Points, se obtuvo %d y %d",
					estado, historias, puntos, g.Cantidad, g.StoryPoints)
			}
			return nil
		}
	}
	return fmt.Errorf("el avance no trae el grupo %q", estado)
}

func (p *pasosAvance) historia(titulo string) (*historiaAvanceJSON, error) {
	a, err := p.avance()
	if err != nil {
		return nil, err
	}
	for gi := range a.Grupos {
		for hi := range a.Grupos[gi].Historias {
			if a.Grupos[gi].Historias[hi].Titulo == titulo {
				return &a.Grupos[gi].Historias[hi], nil
			}
		}
	}
	return nil, fmt.Errorf("el avance no trae la historia %q", titulo)
}

func (p *pasosAvance) figuraSinEstimar(titulo string) error {
	h, err := p.historia(titulo)
	if err != nil {
		return err
	}
	if h.EstimacionSP != nil {
		return fmt.Errorf("la historia %q debía figurar sin estimar, tiene %d Story Points", titulo, *h.EstimacionSP)
	}
	return nil
}

func (p *pasosAvance) figuraEnElSprint(titulo, nombre string) error {
	h, err := p.historia(titulo)
	if err != nil {
		return err
	}
	esperado, ok := p.sprints[nombre]
	if !ok {
		return fmt.Errorf("el Sprint %q no fue creado en el escenario", nombre)
	}
	if h.Sprint == nil || h.Sprint.ID != esperado {
		return fmt.Errorf("la historia %q debía figurar en el Sprint %q (id %d), figura en %v", titulo, nombre, esperado, h.Sprint)
	}
	sprint, err := p.plan.sprints.ObtenerPorID(context.Background(), esperado)
	if err != nil {
		return err
	}
	if h.Sprint.Numero != sprint.Numero {
		return fmt.Errorf("la historia %q informa el número de Sprint %d, se esperaba %d", titulo, h.Sprint.Numero, sprint.Numero)
	}
	return nil
}

func (p *pasosAvance) figuraEnElProductBacklog(titulo string) error {
	h, err := p.historia(titulo)
	if err != nil {
		return err
	}
	if h.Sprint != nil {
		return fmt.Errorf("la historia %q debía figurar en el Product Backlog, figura en el Sprint %d", titulo, h.Sprint.ID)
	}
	return nil
}

func (p *pasosAvance) respuestaEs404() error {
	if p.respuesta == nil || p.respuesta.Code != nethttp.StatusNotFound {
		return fmt.Errorf("se esperaba 404, se obtuvo %v", p.respuesta)
	}
	var cuerpo struct {
		Mensaje string `json:"mensaje"`
	}
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &cuerpo); err != nil {
		return err
	}
	if cuerpo.Mensaje != "proyecto no encontrado" {
		return fmt.Errorf("se esperaba %q, se obtuvo %q", "proyecto no encontrado", cuerpo.Mensaje)
	}
	return nil
}

func (p *pasosAvance) respuestaEs400ConCampo(campo string) error {
	if p.respuesta == nil || p.respuesta.Code != nethttp.StatusBadRequest {
		return fmt.Errorf("se esperaba 400, se obtuvo %v", p.respuesta)
	}
	var cuerpo struct {
		Campo string `json:"campo"`
	}
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &cuerpo); err != nil {
		return err
	}
	if cuerpo.Campo != campo {
		return fmt.Errorf("se esperaba el error sobre el campo %q, se obtuvo %q", campo, cuerpo.Campo)
	}
	return nil
}

// InitializeScenarioAvance registra los pasos de HU-18 (@HU-18) sobre el
// contexto de HU-01 (base de datos y proyecto del Background).
func InitializeScenarioAvance(ctx *godog.ScenarioContext, base *scenarioContext) {
	p := &pasosAvance{base: base}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		p.preparar()
		return goctx, nil
	})

	ctx.Step(`^un proyecto para el avance de historias$`, p.proyecto)
	ctx.Step(`^la historia de avance "([^"]*)" de (\d+) Story Points está Nueva$`, p.historiaNueva)
	ctx.Step(`^la historia de avance "([^"]*)" sin estimar está Nueva$`, p.historiaNuevaSinEstimar)
	ctx.Step(`^la historia de avance "([^"]*)" de (\d+) Story Points está Nueva y asignada al Sprint pendiente "([^"]*)"$`, p.historiaNuevaEnSprintPendiente)
	ctx.Step(`^la historia de avance "([^"]*)" de (\d+) Story Points (?:está|estaba) (En progreso|Completada) en el Sprint "([^"]*)"$`, p.historiaEnSprintActivo)
	ctx.Step(`^se cierra el Sprint "([^"]*)" del avance$`, p.seCierraElSprint)

	ctx.Step(`^consulto el avance de las historias del proyecto$`, p.consultar)
	ctx.Step(`^consulto el avance de las historias del proyecto (\d+)$`, p.consultarProyecto)
	ctx.Step(`^consulto el avance de las historias del proyecto "([^"]*)"$`, p.consultarTexto)

	ctx.Step(`^el avance informa estos grupos en este orden:$`, p.informaGrupos)
	ctx.Step(`^el avance informa un total de (\d+) historias y (\d+) Story Points$`, p.informaTotal)
	ctx.Step(`^el grupo "([^"]*)" informa (\d+) historias y (\d+) Story Points$`, p.grupoInforma)
	ctx.Step(`^la historia "([^"]*)" figura como "sin estimar"$`, p.figuraSinEstimar)
	ctx.Step(`^la historia "([^"]*)" figura en el Sprint "([^"]*)"$`, p.figuraEnElSprint)
	ctx.Step(`^la historia "([^"]*)" figura en el "Product Backlog"$`, p.figuraEnElProductBacklog)
	ctx.Step(`^la respuesta del avance es 404 "proyecto no encontrado"$`, p.respuestaEs404)
	ctx.Step(`^la respuesta del avance es 400 con el campo "([^"]*)"$`, p.respuestaEs400ConCampo)
}
