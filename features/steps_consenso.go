package features

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	apphttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// pasosConsenso agrega los escenarios de HU-14 a los de la votación de HU-02,
// que crean el equipo, la historia y los votos.
type pasosConsenso struct {
	votacion  *pasosVotacion
	registrar *service.RegistrarEstimacion
	abrir     *service.AbrirNuevaRonda
	historial *service.ConsultarHistorialRondas
	vista     service.VistaVotacion
	consulta  service.HistorialRondas
	mux       *http.ServeMux
	rec       *httptest.ResponseRecorder
}

func (c *pasosConsenso) preparar() {
	db := c.votacion.base.db
	historias := repository.NewSQLiteBacklogRepository(db)
	proyectos := repository.NewSQLiteProjectRepository(db)
	votaciones := repository.NewSQLiteVotacionRepository(db)
	c.registrar = service.NewRegistrarEstimacion(historias, proyectos, votaciones, repository.NewSQLiteSprintRepository(db))
	c.abrir = service.NewAbrirNuevaRonda(historias, proyectos, votaciones)
	c.historial = service.NewConsultarHistorialRondas(historias, proyectos, votaciones)
	c.vista, c.consulta = service.VistaVotacion{}, service.HistorialRondas{}
	handler := apphttp.NewVotacionHandler(service.NewObtenerHistoria(historias), nil, nil, nil, nil, c.registrar, nil, nil)
	c.mux = http.NewServeMux()
	c.mux.HandleFunc("POST /backlog/{id}/votacion/estimacion", handler.RegistrarEstimacion)
	c.rec = nil
}

func (c *pasosConsenso) tieneStoryPointsCargados(sp int) error {
	_, err := c.votacion.base.db.Exec("UPDATE historias SET estimacion_sp = ? WHERE id = ?", sp, c.votacion.historiaID)
	return err
}

func (c *pasosConsenso) asignadaASprint(estado string) error {
	sprintID, err := c.votacion.sembrarSprint(estado)
	if err != nil {
		return err
	}
	_, err = c.votacion.base.db.Exec("UPDATE historias SET sprint_id = ? WHERE id = ?", sprintID, c.votacion.historiaID)
	return err
}

func (c *pasosConsenso) pasoAEstado(estado string) error {
	if estado != string(domain.EstadoCompletada) {
		return fmt.Errorf("el paso solo siembra historias %q", domain.EstadoCompletada)
	}
	return c.votacion.completarHistoria()
}

func (c *pasosConsenso) consultaLaVotacion() (err error) {
	c.vista, err = c.votacion.vistaDe("")
	return err
}

func (c *pasosConsenso) hayDiferencias(minimo, maximo int) error {
	if err := c.consultaLaVotacion(); err != nil {
		return err
	}
	esperadas := domain.Diferencias{Hay: true, Minimo: minimo, Maximo: maximo}
	if c.vista.Diferencias == nil || *c.vista.Diferencias != esperadas {
		return fmt.Errorf("se esperaba %+v, se obtuvo %+v", esperadas, c.vista.Diferencias)
	}
	return nil
}

func (c *pasosConsenso) noHayDiferencias() error {
	if c.vista.Diferencias == nil || c.vista.Diferencias.Hay {
		return fmt.Errorf("se esperaba que no hubiera diferencias, se obtuvo %+v", c.vista.Diferencias)
	}
	return nil
}

func (c *pasosConsenso) noInformaDiferencias() error {
	if c.vista.Diferencias != nil {
		return fmt.Errorf("una ronda oculta no debía informar diferencias, se obtuvo %+v", c.vista.Diferencias)
	}
	return nil
}

func (c *pasosConsenso) registraElValor(nombre string, valor int) error {
	return c.registraElValorEn(nombre, valor, c.votacion.historiaID)
}

func (c *pasosConsenso) registraElValorEn(nombre string, valor int, historiaID int64) error {
	id, err := c.votacion.usuario(nombre)
	if err != nil {
		return err
	}
	_, c.votacion.err = c.registrar.Ejecutar(context.Background(), historiaID, id, valor)
	return nil
}

func (c *pasosConsenso) registraEnInexistente(nombre string, valor, historiaID int) error {
	return c.registraElValorEn(nombre, valor, int64(historiaID))
}

func (c *pasosConsenso) registraQuienNoIntegra(nombre string, valor int) error {
	if _, err := c.votacion.crearUsuario(nombre); err != nil {
		return err
	}
	return c.registraElValor(nombre, valor)
}

// enviaALaAPI pasa por el handler: el service recibe un int, así que un valor que
// no es un entero solo se puede expresar en el JSON del pedido.
func (c *pasosConsenso) enviaALaAPI(nombre, valor string) error {
	id, err := c.votacion.usuario(nombre)
	if err != nil {
		return err
	}
	ruta := fmt.Sprintf("/backlog/%d/votacion/estimacion", c.votacion.historiaID)
	cuerpo := fmt.Sprintf(`{"usuario_id":%d,"estimacion_sp":%s}`, id, valor)
	c.rec = httptest.NewRecorder()
	c.mux.ServeHTTP(c.rec, httptest.NewRequest(http.MethodPost, ruta, strings.NewReader(cuerpo)))
	return nil
}

func (c *pasosConsenso) apiRechaza(estado int, campo string) error {
	if c.rec.Code != estado {
		return fmt.Errorf("se esperaba el estado %d, se obtuvo %d (%s)", estado, c.rec.Code, c.rec.Body.String())
	}
	var resp struct {
		Campo string `json:"campo"`
	}
	if err := json.Unmarshal(c.rec.Body.Bytes(), &resp); err != nil {
		return err
	}
	if resp.Campo != campo {
		return fmt.Errorf("se esperaba una advertencia sobre %q, se obtuvo %q", campo, resp.Campo)
	}
	return nil
}

func (c *pasosConsenso) storyPoints() (*int, error) {
	historia, err := repository.NewSQLiteBacklogRepository(c.votacion.base.db).ObtenerPorID(context.Background(), c.votacion.historiaID)
	return historia.EstimacionSP, err
}

func (c *pasosConsenso) quedaConStoryPoints(esperados int) error {
	if c.votacion.err != nil {
		return fmt.Errorf("no se esperaba error: %v", c.votacion.err)
	}
	return c.tieneStoryPoints(esperados)
}

func (c *pasosConsenso) tieneStoryPoints(esperados int) error {
	sp, err := c.storyPoints()
	if err != nil {
		return err
	}
	if sp == nil || *sp != esperados {
		return fmt.Errorf("se esperaban %d Story Points, la historia tiene %v", esperados, sp)
	}
	return nil
}

// iniciaNuevaRonda fotografía la ronda si la hay: sin votación no hay qué comparar.
func (c *pasosConsenso) iniciaNuevaRonda(nombre string) error {
	id, err := c.votacion.usuario(nombre)
	if err != nil {
		return err
	}
	c.votacion.antes, _ = c.votacion.rondaActual()
	_, c.votacion.err = c.abrir.Ejecutar(context.Background(), c.votacion.historiaID, id)
	return nil
}

func (c *pasosConsenso) inicioUnaNuevaRonda(nombre string) error {
	if err := c.iniciaNuevaRonda(nombre); err != nil {
		return err
	}
	return c.votacion.err
}

func (c *pasosConsenso) rondaActualAbierta(numero int) error {
	if c.votacion.err != nil {
		return fmt.Errorf("no se esperaba error: %v", c.votacion.err)
	}
	r, err := c.votacion.rondaActual()
	if err != nil {
		return err
	}
	if r.Numero != numero || len(r.Votos) != 0 || r.Revelada() {
		return fmt.Errorf("se esperaba la ronda %d sin votos ni revelar, se obtuvo %+v", numero, r)
	}
	return nil
}

func (c *pasosConsenso) votacionNoIniciada() error {
	return c.votacion.esperarError(domain.ErrVotacionNoIniciada)
}

func (c *pasosConsenso) consultaElHistorial() (err error) {
	c.consulta, err = c.historial.Ejecutar(context.Background(), c.votacion.historiaID)
	return err
}

func (c *pasosConsenso) historialSoloConLaRonda(numero int) error {
	if len(c.consulta.Rondas) != 1 || c.consulta.Rondas[0].Numero != numero {
		return fmt.Errorf("se esperaba solo la ronda %d, se obtuvo %+v", numero, c.consulta.Rondas)
	}
	return nil
}

func (c *pasosConsenso) rondaConservaSusVotos(numero int, tabla *godog.Table) error {
	if err := c.consultaElHistorial(); err != nil {
		return err
	}
	for _, r := range c.consulta.Rondas {
		if r.Numero != numero {
			continue
		}
		if len(r.Votos) != len(tabla.Rows)-1 {
			return fmt.Errorf("se esperaban %d votos en la ronda %d, hay %+v", len(tabla.Rows)-1, numero, r.Votos)
		}
		for i, fila := range tabla.Rows[1:] {
			carta, err := strconv.Atoi(fila.Cells[1].Value)
			if err != nil {
				return err
			}
			if r.Votos[i].Nombre != fila.Cells[0].Value || r.Votos[i].Valor != carta {
				return fmt.Errorf("se esperaba la carta %d de %q, se obtuvo %+v", carta, fila.Cells[0].Value, r.Votos[i])
			}
		}
		return nil
	}
	return fmt.Errorf("la ronda %d no está en el historial: %+v", numero, c.consulta.Rondas)
}

// InitializeScenarioConsenso registra los pasos de HU-14 sobre los de HU-02.
func InitializeScenarioConsenso(ctx *godog.ScenarioContext, votacion *pasosVotacion) {
	c := &pasosConsenso{votacion: votacion}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		c.preparar()
		return goctx, nil
	})

	ctx.Step(`^que la historia tiene (\d+) Story Points cargados$`, c.tieneStoryPointsCargados)
	ctx.Step(`^que la historia está asignada a un Sprint "([^"]*)"$`, c.asignadaASprint)
	ctx.Step(`^que la historia pasó a estado "([^"]*)"$`, c.pasoAEstado)
	ctx.Step(`^que "([^"]*)" inició una nueva ronda$`, c.inicioUnaNuevaRonda)

	ctx.Step(`^se consulta la votación de la historia$`, c.consultaLaVotacion)
	ctx.Step(`^se consulta el historial de rondas de la historia$`, c.consultaElHistorial)
	ctx.Step(`^"([^"]*)" registra el valor acordado (\d+)$`, c.registraElValor)
	ctx.Step(`^"([^"]*)" registra el valor acordado (\d+) en la historia con identificador (\d+), que no existe$`, c.registraEnInexistente)
	ctx.Step(`^"([^"]*)", que no es integrante del proyecto, registra el valor acordado (\d+)$`, c.registraQuienNoIntegra)
	ctx.Step(`^"([^"]*)" envía a la API el valor acordado (.+)$`, c.enviaALaAPI)
	ctx.Step(`^"([^"]*)" inicia una nueva ronda$`, c.iniciaNuevaRonda)

	ctx.Step(`^el sistema indica que hay diferencias, con mínimo (\d+) y máximo (\d+)$`, c.hayDiferencias)
	ctx.Step(`^el sistema indica que no hay diferencias$`, c.noHayDiferencias)
	ctx.Step(`^el sistema no informa diferencias, mínimo ni máximo$`, c.noInformaDiferencias)
	ctx.Step(`^la historia queda con (\d+) Story Points$`, c.quedaConStoryPoints)
	ctx.Step(`^la historia mantiene sus (\d+) Story Points$`, c.tieneStoryPoints)
	ctx.Step(`^la estimación se rechaza con una advertencia sobre el campo "([^"]*)"$`, votacion.rechazadoConAdvertenciaSobre)
	ctx.Step(`^la estimación se rechaza porque "([^"]*)" no es Scrum Master del proyecto$`, votacion.noEsScrumMaster)
	ctx.Step(`^la estimación se rechaza porque "([^"]*)" no es integrante del proyecto$`, votacion.noEsIntegrante)
	ctx.Step(`^la estimación se rechaza porque la votación no está iniciada$`, c.votacionNoIniciada)
	ctx.Step(`^la estimación se rechaza con error "historia no encontrada"$`, votacion.historiaNoEncontrada)
	ctx.Step(`^la API rechaza la estimación con el estado (\d+) y una advertencia sobre el campo "([^"]*)"$`, c.apiRechaza)
	ctx.Step(`^la ronda actual es la (\d+), abierta y sin votos$`, c.rondaActualAbierta)
	ctx.Step(`^la ronda (\d+) conserva sus votos en el historial:$`, c.rondaConservaSusVotos)
	ctx.Step(`^la nueva ronda se rechaza con una advertencia sobre el campo "([^"]*)"$`, votacion.rechazadoConAdvertenciaSobre)
	ctx.Step(`^la nueva ronda se rechaza porque "([^"]*)" no es Scrum Master del proyecto$`, votacion.noEsScrumMaster)
	ctx.Step(`^la nueva ronda se rechaza porque la votación no está iniciada$`, c.votacionNoIniciada)
	ctx.Step(`^el historial tiene solo la ronda (\d+)$`, c.historialSoloConLaRonda)
}
