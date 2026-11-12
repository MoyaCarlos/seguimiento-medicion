package features

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cucumber/godog"
)

// cierreProyectoContext ejercita cerrar_reabrir_proyecto.feature (#300) contra
// la API real, con el proyecto, las historias y las sesiones de HU-21.
type cierreProyectoContext struct {
	permisos *permisosContext
}

// proyectoIDInexistente no lo crea ningún escenario.
const proyectoIDInexistente = 999

type proyectoConCierre struct {
	FechaCierre *string `json:"fecha_cierre"`
}

// --- Contexto ---

func (c *cierreProyectoContext) cerroElProyecto(persona, proyecto string) error {
	id, err := c.permisos.proyecto(proyecto)
	if err != nil {
		return err
	}
	_, err = c.permisos.esperar(persona, http.MethodPost, fmt.Sprintf("/projects/%d/cerrar", id), nil, http.StatusOK)
	return err
}

func (c *cierreProyectoContext) estaEnSituacion(proyecto, situacion string) error {
	if situacion == "abierto" {
		return nil
	}
	return c.cerroElProyecto(c.permisos.scrumMasterDe(), proyecto)
}

func (c *cierreProyectoContext) tieneUnSprintActivo(proyecto string) error {
	id, err := c.permisos.proyecto(proyecto)
	if err != nil {
		return err
	}
	scrumMaster := c.permisos.scrumMasterDe()
	rec, err := c.permisos.esperar(scrumMaster, http.MethodPost, "/sprints", map[string]any{"proyecto_id": id}, http.StatusCreated)
	if err != nil {
		return err
	}
	sprintID, err := leerID(rec)
	if err != nil {
		return err
	}
	_, err = c.permisos.esperar(scrumMaster, http.MethodPost, fmt.Sprintf("/sprints/%d/iniciar", sprintID),
		map[string]any{"sprint_goal": "g", "fecha_inicio": "2026-11-11", "fecha_fin": "2026-11-20"}, http.StatusOK)
	return err
}

// --- Acciones de quien tiene la sesión ---

func (c *cierreProyectoContext) rutaDe(accion string, proyectoID int64) string {
	return fmt.Sprintf("/projects/%d/%s", proyectoID, accion)
}

func (c *cierreProyectoContext) cierroOReabro(accion, proyecto string) error {
	id, err := c.permisos.proyecto(proyecto)
	if err != nil {
		return err
	}
	_, err = c.permisos.pedir(c.permisos.sesion.actor(), http.MethodPost, c.rutaDe(accion, id), nil)
	return err
}

func (c *cierreProyectoContext) cierro(proyecto string) error {
	return c.cierroOReabro("cerrar", proyecto)
}

func (c *cierreProyectoContext) reabro(proyecto string) error {
	return c.cierroOReabro("reabrir", proyecto)
}

// intentar guarda lo que se ve del proyecto antes de la acción, para que
// «no cambia nada» lo compare después.
func (c *cierreProyectoContext) intentar(accion func() error) error {
	antes, err := c.permisos.estado()
	if err != nil {
		return err
	}
	c.permisos.antes = antes
	return accion()
}

func (c *cierreProyectoContext) intentoCerrarOReabrir(accion, proyecto string) error {
	return c.intentar(func() error { return c.cierroOReabro(accion, proyecto) })
}

func (c *cierreProyectoContext) intentoSobreInexistente(accion string) error {
	_, err := c.permisos.pedir(c.permisos.sesion.actor(), http.MethodPost, c.rutaDe(accion, proyectoIDInexistente), nil)
	return err
}

func (c *cierreProyectoContext) intentoCrearHistoria(titulo, proyecto string) error {
	id, err := c.permisos.proyecto(proyecto)
	if err != nil {
		return err
	}
	return c.intentar(func() error {
		_, err := c.permisos.pedir(c.permisos.sesion.actor(), http.MethodPost, "/backlog",
			map[string]any{"proyecto_id": id, "titulo": titulo, "descripcion": "d", "prioridad": "M"})
		return err
	})
}

func (c *cierreProyectoContext) intentoVotar(valor int, titulo string) error {
	id, err := c.permisos.historia(titulo)
	if err != nil {
		return err
	}
	return c.intentar(func() error {
		_, err := c.permisos.pedir(c.permisos.sesion.actor(), http.MethodPost, fmt.Sprintf("/backlog/%d/votacion/votos", id),
			map[string]any{"valor": valor})
		return err
	})
}

func (c *cierreProyectoContext) intentoRegistrarHoras(horas int, titulo string) error {
	id, err := c.permisos.historia(titulo)
	if err != nil {
		return err
	}
	return c.intentar(func() error {
		_, err := c.permisos.pedir(c.permisos.sesion.actor(), http.MethodPost, fmt.Sprintf("/backlog/%d/horas", id),
			map[string]any{"fecha": "2026-11-11", "actividad": "desarrollo", "horas": horas})
		return err
	})
}

func (c *cierreProyectoContext) consultoElBacklog(proyecto string) error {
	id, err := c.permisos.proyecto(proyecto)
	if err != nil {
		return err
	}
	_, err = c.permisos.pedir(c.permisos.sesion.actor(), http.MethodGet, fmt.Sprintf("/projects/%d/backlog", id), nil)
	return err
}

// --- Respuestas ---

func (c *cierreProyectoContext) rechazaConCodigoYMensaje(codigo int, mensaje string) error {
	rec, err := c.permisos.ultimaRespuesta()
	if err != nil {
		return err
	}
	if rec.Code != codigo || c.permisos.mensajeDe(rec) != mensaje {
		return fmt.Errorf("se esperaba %d %q, se obtuvo %d %s", codigo, mensaje, rec.Code, rec.Body.String())
	}
	return nil
}

func (c *cierreProyectoContext) ultimaFueExitosa() error {
	rec, err := c.permisos.ultimaRespuesta()
	if err != nil {
		return err
	}
	if rec.Code != http.StatusOK {
		return fmt.Errorf("se esperaba 200, se obtuvo %d %s", rec.Code, rec.Body.String())
	}
	return nil
}

// fechaDeCierre lee la fecha de cierre del detalle del proyecto (nil: abierto).
func (c *cierreProyectoContext) fechaDeCierre(proyecto string) (*string, error) {
	id, err := c.permisos.proyecto(proyecto)
	if err != nil {
		return nil, err
	}
	rec, err := c.permisos.sesion.peticionDe(c.permisos.scrumMasterDe(), http.MethodGet, fmt.Sprintf("/projects/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var p proyectoConCierre
	if rec.Code != http.StatusOK || json.Unmarshal(rec.Body.Bytes(), &p) != nil {
		return nil, fmt.Errorf("no se pudo consultar %q: %d %s", proyecto, rec.Code, rec.Body.String())
	}
	return p.FechaCierre, nil
}

func (c *cierreProyectoContext) estadoDe(proyecto string) (string, error) {
	id, err := c.permisos.proyecto(proyecto)
	if err != nil {
		return "", err
	}
	rec, err := c.permisos.sesion.peticionDe(c.permisos.scrumMasterDe(), http.MethodGet, fmt.Sprintf("/projects/%d/status", id), nil)
	if err != nil {
		return "", err
	}
	var resp struct {
		Estado string `json:"estado"`
	}
	if rec.Code != http.StatusOK || json.Unmarshal(rec.Body.Bytes(), &resp) != nil {
		return "", fmt.Errorf("no se pudo consultar el estado de %q: %d %s", proyecto, rec.Code, rec.Body.String())
	}
	return resp.Estado, nil
}

func (c *cierreProyectoContext) quedaConFechaDeCierre(proyecto, estado, fecha string) error {
	if err := c.ultimaFueExitosa(); err != nil {
		return err
	}
	cierre, err := c.fechaDeCierre(proyecto)
	if err != nil {
		return err
	}
	if cierre == nil || *cierre != fecha {
		return fmt.Errorf("se esperaba la fecha de cierre %s, se obtuvo %v", fecha, cierre)
	}
	return c.tieneElEstado(proyecto, estado)
}

func (c *cierreProyectoContext) sigueAbierto(proyecto string) error {
	cierre, err := c.fechaDeCierre(proyecto)
	if err != nil {
		return err
	}
	if cierre != nil {
		return fmt.Errorf("%q debía seguir abierto, tiene fecha de cierre %s", proyecto, *cierre)
	}
	return nil
}

func (c *cierreProyectoContext) quedaAbiertoConEstado(proyecto, estado string) error {
	if err := c.ultimaFueExitosa(); err != nil {
		return err
	}
	if err := c.sigueAbierto(proyecto); err != nil {
		return err
	}
	return c.tieneElEstado(proyecto, estado)
}

func (c *cierreProyectoContext) tieneElEstado(proyecto, estado string) error {
	actual, err := c.estadoDe(proyecto)
	if err != nil {
		return err
	}
	if actual != estado {
		return fmt.Errorf("se esperaba el estado %q, se obtuvo %q", estado, actual)
	}
	return nil
}

func (c *cierreProyectoContext) puedeCrearLaHistoria(persona, titulo, proyecto string) error {
	id, err := c.permisos.proyecto(proyecto)
	if err != nil {
		return err
	}
	_, err = c.permisos.esperar(persona, http.MethodPost, "/backlog",
		map[string]any{"proyecto_id": id, "titulo": titulo, "descripcion": "d", "prioridad": "M"}, http.StatusCreated)
	return err
}

func (c *cierreProyectoContext) veoLaHistoriaEnElBacklog(titulo string) error {
	rec, err := c.permisos.ultimaRespuesta()
	if err != nil {
		return err
	}
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), fmt.Sprintf("%q", titulo)) {
		return fmt.Errorf("se esperaba 200 con la historia %q, se obtuvo %d %s", titulo, rec.Code, rec.Body.String())
	}
	return nil
}

// InitializeScenarioCierreProyecto registra los pasos de
// cerrar_reabrir_proyecto.feature (#300) sobre el contexto de HU-21.
func InitializeScenarioCierreProyecto(ctx *godog.ScenarioContext, permisos *permisosContext) {
	c := &cierreProyectoContext{permisos: permisos}

	ctx.Step(`^que "([^"]*)" cerró el proyecto "([^"]*)"$`, c.cerroElProyecto)
	ctx.Step(`^que "([^"]*)" está (abierto|cerrado)$`, c.estaEnSituacion)
	ctx.Step(`^que "([^"]*)" tiene un Sprint activo$`, c.tieneUnSprintActivo)

	ctx.Step(`^cierro el proyecto "([^"]*)"$`, c.cierro)
	ctx.Step(`^reabro el proyecto "([^"]*)"$`, c.reabro)
	ctx.Step(`^intento (cerrar|reabrir) el proyecto "([^"]*)"$`, c.intentoCerrarOReabrir)
	ctx.Step(`^intento (cerrar|reabrir) un proyecto que no existe$`, c.intentoSobreInexistente)
	ctx.Step(`^intento crear la historia "([^"]*)" en "([^"]*)"$`, c.intentoCrearHistoria)
	ctx.Step(`^intento votar (\d+) en "([^"]*)"$`, c.intentoVotar)
	ctx.Step(`^intento registrar (\d+) horas en "([^"]*)"$`, c.intentoRegistrarHoras)
	ctx.Step(`^consulto el backlog de "([^"]*)"$`, c.consultoElBacklog)

	ctx.Step(`^el sistema rechaza la petición con (\d+) y el mensaje "([^"]*)"$`, c.rechazaConCodigoYMensaje)
	ctx.Step(`^"([^"]*)" queda "([^"]*)" con la fecha de cierre "([^"]*)"$`, c.quedaConFechaDeCierre)
	ctx.Step(`^"([^"]*)" sigue abierto$`, c.sigueAbierto)
	ctx.Step(`^"([^"]*)" queda abierto con el estado "([^"]*)"$`, c.quedaAbiertoConEstado)
	ctx.Step(`^"([^"]*)" puede crear la historia "([^"]*)" en "([^"]*)"$`, c.puedeCrearLaHistoria)
	ctx.Step(`^veo la historia "([^"]*)" en el backlog$`, c.veoLaHistoriaEnElBacklog)
}
