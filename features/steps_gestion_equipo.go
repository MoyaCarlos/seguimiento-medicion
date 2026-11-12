package features

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cucumber/godog"
)

// gestionEquipoContext ejercita gestion_equipo.feature (HU-16) contra la API
// real, con el proyecto, las historias y las sesiones de HU-21.
type gestionEquipoContext struct {
	permisos *permisosContext
	// metricasAntes es el cuerpo de GET /projects/{id}/metricas antes de la baja.
	metricasAntes string
}

func (c *gestionEquipoContext) reiniciar() {
	c.metricasAntes = ""
}

func (c *gestionEquipoContext) rutaDeBaja(proyectoID, usuarioID int64) string {
	return fmt.Sprintf("/projects/%d/members/%d", proyectoID, usuarioID)
}

// --- Contexto ---

// votoDe registra el voto de persona en la ronda abierta de la historia.
func (c *gestionEquipoContext) votoDe(persona string, valor int, titulo string) error {
	id, err := c.permisos.historia(titulo)
	if err != nil {
		return err
	}
	_, err = c.permisos.esperar(persona, http.MethodPost, fmt.Sprintf("/backlog/%d/votacion/votos", id),
		map[string]any{"valor": valor}, http.StatusOK)
	return err
}

func (c *gestionEquipoContext) revelar(persona, titulo string) error {
	id, err := c.permisos.historia(titulo)
	if err != nil {
		return err
	}
	_, err = c.permisos.esperar(persona, http.MethodPost, fmt.Sprintf("/backlog/%d/votacion/revelar", id), nil, http.StatusOK)
	return err
}

// votoEnRondaRevelada: el Scrum Master abre la votación, vota 3, persona vota
// valor y el Scrum Master la revela.
func (c *gestionEquipoContext) votoEnRondaRevelada(persona string, valor int, titulo string) error {
	scrumMaster := c.permisos.scrumMasterDe()
	if err := c.permisos.inicioLaVotacionComo(scrumMaster, titulo); err != nil {
		return err
	}
	if err := c.votoDe(scrumMaster, 3, titulo); err != nil {
		return err
	}
	if err := c.votoDe(persona, valor, titulo); err != nil {
		return err
	}
	return c.revelar(scrumMaster, titulo)
}

// imputoEnSprintFinalizado: el Scrum Master crea e inicia un Sprint con la
// historia, persona imputa las horas y el Sprint se cierra.
func (c *gestionEquipoContext) imputoEnSprintFinalizado(persona string, horas int, titulo string) error {
	proyectoID, err := c.permisos.proyecto("Tienda")
	if err != nil {
		return err
	}
	historiaID, err := c.permisos.historia(titulo)
	if err != nil {
		return err
	}
	scrumMaster := c.permisos.scrumMasterDe()
	rec, err := c.permisos.esperar(scrumMaster, http.MethodPost, "/sprints", map[string]any{"proyecto_id": proyectoID}, http.StatusCreated)
	if err != nil {
		return err
	}
	sprintID, err := leerID(rec)
	if err != nil {
		return err
	}
	pasos := []struct {
		persona, metodo, ruta string
		cuerpo                any
		codigo                int
	}{
		{scrumMaster, http.MethodPost, fmt.Sprintf("/sprints/%d/iniciar", sprintID),
			map[string]any{"sprint_goal": "g", "fecha_inicio": "2026-11-01", "fecha_fin": "2026-11-10"}, http.StatusOK},
		{scrumMaster, http.MethodPut, fmt.Sprintf("/backlog/%d/sprint", historiaID), map[string]any{"sprint_id": sprintID}, http.StatusOK},
		{persona, http.MethodPost, fmt.Sprintf("/backlog/%d/horas", historiaID),
			map[string]any{"fecha": "2026-11-05", "actividad": "desarrollo", "horas": horas}, http.StatusCreated},
		{scrumMaster, http.MethodPost, fmt.Sprintf("/sprints/%d/cerrar", sprintID), nil, http.StatusOK},
	}
	for _, p := range pasos {
		if _, err := c.permisos.esperar(p.persona, p.metodo, p.ruta, p.cuerpo, p.codigo); err != nil {
			return err
		}
	}
	return nil
}

func (c *gestionEquipoContext) metricas(proyecto string) (string, error) {
	id, err := c.permisos.proyecto(proyecto)
	if err != nil {
		return "", err
	}
	rec, err := c.permisos.sesion.peticionDe(c.permisos.scrumMasterDe(), http.MethodGet, fmt.Sprintf("/projects/%d/metricas", id), nil)
	if err != nil {
		return "", err
	}
	if rec.Code != http.StatusOK {
		return "", fmt.Errorf("no se pudieron consultar las métricas: %d %s", rec.Code, rec.Body.String())
	}
	return rec.Body.String(), nil
}

func (c *gestionEquipoContext) anoteLasMetricas(proyecto string) error {
	var err error
	c.metricasAntes, err = c.metricas(proyecto)
	return err
}

// pedirBaja pide como persona la baja de quien en el proyecto.
func (c *gestionEquipoContext) pedirBaja(persona, quien string, proyectoID int64) (*httptest.ResponseRecorder, error) {
	usuarioID, err := c.permisos.sesion.idDe(quien)
	if err != nil {
		return nil, err
	}
	return c.permisos.pedir(persona, http.MethodDelete, c.rutaDeBaja(proyectoID, usuarioID), nil)
}

func (c *gestionEquipoContext) quito(quien, proyecto string) error {
	id, err := c.permisos.proyecto(proyecto)
	if err != nil {
		return err
	}
	_, err = c.pedirBaja(c.permisos.sesion.actor(), quien, id)
	return err
}

func (c *gestionEquipoContext) quitoComo(persona, quien, proyecto string) error {
	id, err := c.permisos.proyecto(proyecto)
	if err != nil {
		return err
	}
	rec, err := c.pedirBaja(persona, quien, id)
	if err != nil {
		return err
	}
	if rec.Code != http.StatusNoContent {
		return fmt.Errorf("no se pudo quitar a %q: %d %s", quien, rec.Code, rec.Body.String())
	}
	return nil
}

// intentoQuitar guarda lo que se ve del proyecto antes, para «no cambia nada».
func (c *gestionEquipoContext) intentoQuitar(quien, proyecto string) error {
	antes, err := c.permisos.estado()
	if err != nil {
		return err
	}
	c.permisos.antes = antes
	return c.quito(quien, proyecto)
}

func (c *gestionEquipoContext) intentoQuitarDeInexistente(quien string) error {
	_, err := c.pedirBaja(c.permisos.sesion.actor(), quien, proyectoIDInexistente)
	return err
}

// --- Respuestas ---

func (c *gestionEquipoContext) bajaSinContenido() error {
	rec, err := c.permisos.ultimaRespuesta()
	if err != nil {
		return err
	}
	if rec.Code != http.StatusNoContent || rec.Body.Len() != 0 {
		return fmt.Errorf("se esperaba 204 sin cuerpo, se obtuvo %d %s", rec.Code, rec.Body.String())
	}
	return nil
}

func (c *gestionEquipoContext) vecesEnElEquipo(persona, proyecto, rol string) (int, error) {
	lista, err := c.permisos.integrantes(proyecto)
	if err != nil {
		return 0, err
	}
	veces := 0
	for _, m := range lista {
		if m.Nombre == persona && (rol == "" || m.Rol == rol) {
			veces++
		}
	}
	return veces, nil
}

func (c *gestionEquipoContext) noFiguraEnElEquipo(persona, proyecto string) error {
	veces, err := c.vecesEnElEquipo(persona, proyecto, "")
	if err != nil {
		return err
	}
	if veces != 0 {
		return fmt.Errorf("%q no debía figurar en el equipo de %q", persona, proyecto)
	}
	return nil
}

func (c *gestionEquipoContext) figuraUnaSolaVez(persona, proyecto, rol string) error {
	if err := c.permisos.cadaAccionSeRealiza(); err != nil {
		return err
	}
	veces, err := c.vecesEnElEquipo(persona, proyecto, rol)
	if err != nil {
		return err
	}
	todas, err := c.vecesEnElEquipo(persona, proyecto, "")
	if err != nil {
		return err
	}
	if veces != 1 || todas != 1 {
		return fmt.Errorf("se esperaba a %q una sola vez como %q, figura %d veces (%d con ese rol)", persona, rol, todas, veces)
	}
	return nil
}

func (c *gestionEquipoContext) veSoloEnSusProyectos(persona, proyecto, rol string) error {
	if _, err := c.permisos.esperar(persona, http.MethodGet, "/projects", nil, http.StatusOK); err != nil {
		return err
	}
	return c.permisos.veoSoloConRol(proyecto, rol)
}

// votacionNuevaSoloConElVotoDe: el Scrum Master abre la votación de la
// historia, vota y la revela sin esperar a nadie más.
func (c *gestionEquipoContext) votacionNuevaSoloConElVotoDe(titulo, persona string) error {
	if err := c.permisos.inicioLaVotacionComo(persona, titulo); err != nil {
		return err
	}
	if err := c.votoDe(persona, 3, titulo); err != nil {
		return err
	}
	return c.revelar(persona, titulo)
}

type votoRevelado struct {
	Nombre string `json:"nombre"`
	Valor  int    `json:"valor"`
}

func (c *gestionEquipoContext) historialMuestraElVoto(titulo string, valor int, persona string) error {
	id, err := c.permisos.historia(titulo)
	if err != nil {
		return err
	}
	rec, err := c.permisos.esperar(c.permisos.scrumMasterDe(), http.MethodGet, fmt.Sprintf("/backlog/%d/votacion/rondas", id), nil, http.StatusOK)
	if err != nil {
		return err
	}
	var historial struct {
		Rondas []struct {
			Votos []votoRevelado `json:"votos"`
		} `json:"rondas"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &historial); err != nil {
		return fmt.Errorf("historial inválido: %s", rec.Body.String())
	}
	for _, r := range historial.Rondas {
		for _, v := range r.Votos {
			if v.Nombre == persona && v.Valor == valor {
				return nil
			}
		}
	}
	return fmt.Errorf("el historial no muestra el voto %d de %q: %s", valor, persona, rec.Body.String())
}

func (c *gestionEquipoContext) horasMuestran(titulo string, horas int, persona string) error {
	id, err := c.permisos.historia(titulo)
	if err != nil {
		return err
	}
	rec, err := c.permisos.esperar(c.permisos.scrumMasterDe(), http.MethodGet, fmt.Sprintf("/backlog/%d/horas", id), nil, http.StatusOK)
	if err != nil {
		return err
	}
	var resp struct {
		Registros []struct {
			Integrante string  `json:"integrante"`
			Horas      float64 `json:"horas"`
		} `json:"registros"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		return fmt.Errorf("horas inválidas: %s", rec.Body.String())
	}
	if len(resp.Registros) != 1 || resp.Registros[0].Integrante != persona || resp.Registros[0].Horas != float64(horas) {
		return fmt.Errorf("se esperaban %d horas de %q, se obtuvo %s", horas, persona, rec.Body.String())
	}
	return nil
}

func (c *gestionEquipoContext) metricasSinCambios(proyecto string) error {
	despues, err := c.metricas(proyecto)
	if err != nil {
		return err
	}
	if c.metricasAntes == "" || despues != c.metricasAntes {
		return fmt.Errorf("las métricas cambiaron:\nantes:   %s\ndespués: %s", c.metricasAntes, despues)
	}
	return nil
}

func (c *gestionEquipoContext) votacionReveladaMuestraSolo(titulo string, valor int, persona string) error {
	vista, err := c.permisos.vistaDeLaVotacion(c.permisos.scrumMasterDe(), titulo, "")
	if err != nil {
		return err
	}
	personaID, err := c.permisos.sesion.idDe(persona)
	if err != nil {
		return err
	}
	var votos []string
	for _, m := range vista.Integrantes {
		if m.YaVoto {
			votos = append(votos, fmt.Sprintf("%d:%v", m.UsuarioID, *m.Valor))
		}
	}
	if len(votos) != 1 || votos[0] != fmt.Sprintf("%d:%d", personaID, valor) {
		return fmt.Errorf("se esperaba solo el voto %d de %q, se obtuvo %v", valor, persona, votos)
	}
	return nil
}

func (c *gestionEquipoContext) puedeVotar(persona string, valor int, titulo string) error {
	if err := c.permisos.cadaAccionSeRealiza(); err != nil {
		return err
	}
	return c.votoDe(persona, valor, titulo)
}

// InitializeScenarioGestionEquipo registra los pasos de gestion_equipo.feature
// (HU-16) sobre el contexto de HU-21.
func InitializeScenarioGestionEquipo(ctx *godog.ScenarioContext, permisos *permisosContext) {
	c := &gestionEquipoContext{permisos: permisos}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		c.reiniciar()
		return goctx, nil
	})

	ctx.Step(`^que "([^"]*)" votó (\d+) en la ronda revelada de "([^"]*)"$`, c.votoEnRondaRevelada)
	ctx.Step(`^que "([^"]*)" imputó (\d+) horas en "([^"]*)" en un Sprint ya finalizado$`, c.imputoEnSprintFinalizado)
	ctx.Step(`^que anoté las métricas de "([^"]*)"$`, c.anoteLasMetricas)
	ctx.Step(`^que "([^"]*)" quitó a "([^"]*)" del equipo de "([^"]*)"$`, c.quitoComo)
	ctx.Step(`^que "([^"]*)" ya votó (\d+) en "([^"]*)"$`, c.votoDe)

	ctx.Step(`^quito a "([^"]*)" del equipo de "([^"]*)"$`, c.quito)
	ctx.Step(`^intento quitar a "([^"]*)" del equipo de "([^"]*)"$`, c.intentoQuitar)
	ctx.Step(`^intento quitar a "([^"]*)" de un proyecto que no existe$`, c.intentoQuitarDeInexistente)
	ctx.Step(`^"([^"]*)" revela la votación de "([^"]*)"$`, c.revelar)

	ctx.Step(`^la baja se confirma sin contenido$`, c.bajaSinContenido)
	ctx.Step(`^"([^"]*)" no figura en el equipo de "([^"]*)"$`, c.noFiguraEnElEquipo)
	ctx.Step(`^"([^"]*)" figura una sola vez en el equipo de "([^"]*)" como "([^"]*)"$`, c.figuraUnaSolaVez)
	ctx.Step(`^"([^"]*)" ve solo "([^"]*)" en sus proyectos con el rol "([^"]*)"$`, c.veSoloEnSusProyectos)
	ctx.Step(`^una votación nueva de "([^"]*)" se revela solo con el voto de "([^"]*)"$`, c.votacionNuevaSoloConElVotoDe)
	ctx.Step(`^el historial de rondas de "([^"]*)" muestra el voto (\d+) de "([^"]*)"$`, c.historialMuestraElVoto)
	ctx.Step(`^las horas de "([^"]*)" muestran (\d+) horas de "([^"]*)"$`, c.horasMuestran)
	ctx.Step(`^las métricas de "([^"]*)" son las mismas que antes de la baja$`, c.metricasSinCambios)
	ctx.Step(`^la votación revelada de "([^"]*)" muestra solo el voto (\d+) de "([^"]*)"$`, c.votacionReveladaMuestraSolo)
	ctx.Step(`^"([^"]*)" puede votar (\d+) en "([^"]*)"$`, c.puedeVotar)
}
