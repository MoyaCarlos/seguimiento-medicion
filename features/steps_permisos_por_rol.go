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
)

// permisosContext ejercita los escenarios de HU-21 solo contra la API real
// (app.NuevoMux): cada persona se registra con /auth/registro y actúa con su
// token. No usa atajos de prueba en el código de producción.
type permisosContext struct {
	sesion     *sesionDeEscenario
	creador    string
	proyectos  map[string]int64
	historias  map[string]int64
	ultima     *httptest.ResponseRecorder
	resultados []string
	antes      string
	vista      votacionDeEscenario
}

type votacionDeEscenario struct {
	Integrantes []struct {
		UsuarioID int64 `json:"usuario_id"`
		YaVoto    bool  `json:"ya_voto"`
		Valor     *int  `json:"valor"`
	} `json:"integrantes"`
}

func (c *permisosContext) reiniciar() {
	c.proyectos = map[string]int64{}
	c.historias = map[string]int64{}
	c.ultima, c.resultados, c.antes, c.creador = nil, nil, "", ""
	c.vista = votacionDeEscenario{}
}

// pedir llama a la API como nombre y guarda la respuesta.
func (c *permisosContext) pedir(nombre, metodo, ruta string, cuerpo any) (*httptest.ResponseRecorder, error) {
	rec, err := c.sesion.peticionDe(nombre, metodo, ruta, cuerpo)
	if err != nil {
		return nil, err
	}
	c.ultima = rec
	return rec, nil
}

// esperar llama a la API como nombre y exige el código esperado.
func (c *permisosContext) esperar(nombre, metodo, ruta string, cuerpo any, codigo int) (*httptest.ResponseRecorder, error) {
	rec, err := c.pedir(nombre, metodo, ruta, cuerpo)
	if err != nil {
		return nil, err
	}
	if rec.Code != codigo {
		return nil, fmt.Errorf("%s %s como %q: se esperaba %d, se obtuvo %d %s", metodo, ruta, nombre, codigo, rec.Code, rec.Body.String())
	}
	return rec, nil
}

func (c *permisosContext) proyecto(nombre string) (int64, error) {
	id, ok := c.proyectos[nombre]
	if !ok {
		return 0, fmt.Errorf("el escenario no creó el proyecto %q", nombre)
	}
	return id, nil
}

func (c *permisosContext) historia(titulo string) (int64, error) {
	id, ok := c.historias[titulo]
	if !ok {
		return 0, fmt.Errorf("el escenario no creó la historia %q", titulo)
	}
	return id, nil
}

func leerID(rec *httptest.ResponseRecorder) (int64, error) {
	var resp struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil || resp.ID == 0 {
		return 0, fmt.Errorf("respuesta sin id: %s", rec.Body.String())
	}
	return resp.ID, nil
}

// --- Antecedentes ---

func (c *permisosContext) creoElProyecto(persona, nombre string) error {
	rec, err := c.esperar(persona, http.MethodPost, "/projects", map[string]any{"nombre": nombre}, http.StatusCreated)
	if err != nil {
		return err
	}
	c.creador = persona
	c.proyectos[nombre], err = leerID(rec)
	return err
}

func (c *permisosContext) esProductBuilderDe(persona, proyecto string) error {
	id, err := c.proyecto(proyecto)
	if err != nil {
		return err
	}
	if _, err := c.sesion.token(persona); err != nil {
		return err
	}
	_, err = c.esperar(c.scrumMasterDe(), http.MethodPost, fmt.Sprintf("/projects/%d/members", id),
		map[string]any{"nombre": persona, "rol": "product_builder"}, http.StatusCreated)
	return err
}

// scrumMasterDe es quien creó el proyecto del escenario: su Scrum Master.
func (c *permisosContext) scrumMasterDe() string { return c.creador }

func (c *permisosContext) tieneLaHistoria(proyecto, titulo string) error {
	id, err := c.proyecto(proyecto)
	if err != nil {
		return err
	}
	rec, err := c.esperar(c.scrumMasterDe(), http.MethodPost, "/backlog",
		map[string]any{"proyecto_id": id, "titulo": titulo, "descripcion": "d", "prioridad": "M"}, http.StatusCreated)
	if err != nil {
		return err
	}
	c.historias[titulo], err = leerID(rec)
	return err
}

func (c *permisosContext) tieneCuenta(persona string) error {
	_, err := c.sesion.token(persona)
	return err
}

func (c *permisosContext) noInicieSesion() error {
	c.sesion.nombre = ""
	return nil
}

// --- Acciones de quien tiene la sesión ---

func (c *permisosContext) crearSprint(proyecto string) error {
	id, err := c.proyecto(proyecto)
	if err != nil {
		return err
	}
	return c.accion(http.MethodPost, "/sprints", map[string]any{"proyecto_id": id})
}

func (c *permisosContext) iniciarVotacion(titulo string) error {
	id, err := c.historia(titulo)
	if err != nil {
		return err
	}
	return c.accion(http.MethodPost, fmt.Sprintf("/backlog/%d/votacion", id), nil)
}

func (c *permisosContext) editarProyecto(proyecto string) error {
	id, err := c.proyecto(proyecto)
	if err != nil {
		return err
	}
	return c.accion(http.MethodPut, fmt.Sprintf("/projects/%d", id),
		map[string]any{"nombre": proyecto, "descripcion": "editado", "fecha_inicio": nil, "fecha_fin": nil})
}

func (c *permisosContext) agregarAlEquipo(persona, proyecto string) error {
	return c.agregarAlEquipoComo(persona, proyecto, "")
}

func (c *permisosContext) agregarAlEquipoComo(persona, proyecto, rol string) error {
	id, err := c.proyecto(proyecto)
	if err != nil {
		return err
	}
	cuerpo := map[string]any{"nombre": persona}
	if rol != "" {
		cuerpo["rol"] = rol
	}
	return c.accion(http.MethodPost, fmt.Sprintf("/projects/%d/members", id), cuerpo)
}

// accion pide como quien tiene la sesión y anota el resultado.
func (c *permisosContext) accion(metodo, ruta string, cuerpo any) error {
	rec, err := c.pedir(c.sesion.actor(), metodo, ruta, cuerpo)
	if err != nil {
		return err
	}
	c.resultados = append(c.resultados, fmt.Sprintf("%s %s → %d %s", metodo, ruta, rec.Code, strings.TrimSpace(rec.Body.String())))
	return nil
}

func (c *permisosContext) cadaAccionSeRealiza() error {
	if len(c.resultados) == 0 {
		return fmt.Errorf("el escenario no hizo ninguna acción")
	}
	for _, r := range c.resultados {
		codigo := strings.SplitN(strings.SplitN(r, "→ ", 2)[1], " ", 2)[0]
		if codigo != "200" && codigo != "201" {
			return fmt.Errorf("la acción no se realizó: %s", r)
		}
	}
	return nil
}

// intento guarda el estado del proyecto antes de la acción para comprobar
// después que no cambió nada.
func (c *permisosContext) intento(accion func(string) error) func(string) error {
	return func(arg string) error {
		antes, err := c.estado()
		if err != nil {
			return err
		}
		c.antes = antes
		return accion(arg)
	}
}

func (c *permisosContext) intentoAgregarAlEquipo(persona, proyecto string) error {
	antes, err := c.estado()
	if err != nil {
		return err
	}
	c.antes = antes
	return c.agregarAlEquipo(persona, proyecto)
}

// estado es lo que el Scrum Master ve del proyecto «Tienda»: panel, backlog,
// Sprints, equipo, y votación y horas de «Carrito».
func (c *permisosContext) estado() (string, error) {
	proyecto, err := c.proyecto("Tienda")
	if err != nil {
		return "", err
	}
	historia, err := c.historia("Carrito")
	if err != nil {
		return "", err
	}
	var partes []string
	for _, ruta := range []string{
		fmt.Sprintf("/projects/%d", proyecto),
		fmt.Sprintf("/projects/%d/backlog", proyecto),
		fmt.Sprintf("/projects/%d/sprints", proyecto),
		fmt.Sprintf("/projects/%d/members", proyecto),
		fmt.Sprintf("/backlog/%d/votacion", historia),
		fmt.Sprintf("/backlog/%d/horas", historia),
	} {
		rec, err := c.sesion.peticionDe(c.scrumMasterDe(), http.MethodGet, ruta, nil)
		if err != nil {
			return "", err
		}
		partes = append(partes, fmt.Sprintf("%s %d %s", ruta, rec.Code, rec.Body.String()))
	}
	return strings.Join(partes, "\n"), nil
}

func (c *permisosContext) noCambiaNada() error {
	despues, err := c.estado()
	if err != nil {
		return err
	}
	if despues != c.antes {
		return fmt.Errorf("el proyecto cambió:\nantes:\n%s\ndespués:\n%s", c.antes, despues)
	}
	return nil
}

// --- Respuestas ---

func (c *permisosContext) ultimaRespuesta() (*httptest.ResponseRecorder, error) {
	if c.ultima == nil {
		return nil, fmt.Errorf("el escenario no hizo ninguna petición")
	}
	return c.ultima, nil
}

func (c *permisosContext) mensajeDe(rec *httptest.ResponseRecorder) string {
	var resp struct {
		Mensaje string `json:"mensaje"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	return resp.Mensaje
}

func (c *permisosContext) responde403ConMensajeQueEmpieza(prefijo string) error {
	rec, err := c.ultimaRespuesta()
	if err != nil {
		return err
	}
	if rec.Code != http.StatusForbidden || !strings.HasPrefix(c.mensajeDe(rec), prefijo) {
		return fmt.Errorf("se esperaba 403 con un mensaje que empiece con %q, se obtuvo %d %s", prefijo, rec.Code, rec.Body.String())
	}
	return nil
}

func (c *permisosContext) responde403ConMensaje(mensaje string) error {
	rec, err := c.ultimaRespuesta()
	if err != nil {
		return err
	}
	if rec.Code != http.StatusForbidden || c.mensajeDe(rec) != mensaje {
		return fmt.Errorf("se esperaba 403 %q, se obtuvo %d %s", mensaje, rec.Code, rec.Body.String())
	}
	return nil
}

func (c *permisosContext) rechazaConCodigo(codigo int) error {
	rec, err := c.ultimaRespuesta()
	if err != nil {
		return err
	}
	if rec.Code != codigo {
		return fmt.Errorf("se esperaba %d, se obtuvo %d %s", codigo, rec.Code, rec.Body.String())
	}
	return nil
}

// --- Mis proyectos ---

func (c *permisosContext) consultoMisProyectos() error {
	_, err := c.pedir(c.sesion.actor(), http.MethodGet, "/projects", nil)
	return err
}

func (c *permisosContext) abroElProyecto(proyecto string) error {
	id, err := c.proyecto(proyecto)
	if err != nil {
		return err
	}
	_, err = c.pedir(c.sesion.actor(), http.MethodGet, fmt.Sprintf("/projects/%d", id), nil)
	return err
}

func (c *permisosContext) misProyectos() ([]struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}, error) {
	rec, err := c.ultimaRespuesta()
	if err != nil {
		return nil, err
	}
	if rec.Code != http.StatusOK {
		return nil, fmt.Errorf("se esperaba 200, se obtuvo %d %s", rec.Code, rec.Body.String())
	}
	var lista []struct {
		Nombre string `json:"nombre"`
		Rol    string `json:"rol"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &lista); err != nil {
		return nil, fmt.Errorf("respuesta inválida: %s", rec.Body.String())
	}
	return lista, nil
}

func (c *permisosContext) veoSoloConRol(proyecto, rol string) error {
	lista, err := c.misProyectos()
	if err != nil {
		return err
	}
	if len(lista) != 1 || lista[0].Nombre != proyecto || lista[0].Rol != rol {
		return fmt.Errorf("se esperaba solo %q con el rol %q, se obtuvo %+v", proyecto, rol, lista)
	}
	return nil
}

func (c *permisosContext) reciboListaVacia() error {
	rec, err := c.ultimaRespuesta()
	if err != nil {
		return err
	}
	if body := strings.TrimSpace(rec.Body.String()); rec.Code != http.StatusOK || body != "[]" {
		return fmt.Errorf("se esperaba 200 [], se obtuvo %d %s", rec.Code, body)
	}
	return nil
}

// --- Equipo ---

func (c *permisosContext) integrantes(proyecto string) ([]struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}, error) {
	id, err := c.proyecto(proyecto)
	if err != nil {
		return nil, err
	}
	rec, err := c.sesion.peticionDe(c.scrumMasterDe(), http.MethodGet, fmt.Sprintf("/projects/%d/members", id), nil)
	if err != nil {
		return nil, err
	}
	var lista []struct {
		Nombre string `json:"nombre"`
		Rol    string `json:"rol"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &lista); err != nil {
		return nil, fmt.Errorf("equipo inválido: %d %s", rec.Code, rec.Body.String())
	}
	return lista, nil
}

func (c *permisosContext) unicoScrumMaster(proyecto, persona string) error {
	lista, err := c.integrantes(proyecto)
	if err != nil {
		return err
	}
	var scrumMasters []string
	for _, m := range lista {
		if m.Rol == "scrum_master" {
			scrumMasters = append(scrumMasters, m.Nombre)
		}
	}
	if len(scrumMasters) != 1 || scrumMasters[0] != persona {
		return fmt.Errorf("se esperaba a %q como único Scrum Master, se obtuvo %v", persona, scrumMasters)
	}
	return nil
}

func (c *permisosContext) quedaComo(persona, rol, proyecto string) error {
	if err := c.cadaAccionSeRealiza(); err != nil {
		return err
	}
	lista, err := c.integrantes(proyecto)
	if err != nil {
		return err
	}
	for _, m := range lista {
		if m.Nombre == persona && m.Rol == rol {
			return nil
		}
	}
	return fmt.Errorf("%q no quedó como %q: %+v", persona, rol, lista)
}

func (c *permisosContext) agregarSinCuentaSeRechaza(persona string) error {
	id, err := c.proyecto("Tienda")
	if err != nil {
		return err
	}
	rec, err := c.pedir(c.sesion.actor(), http.MethodPost, fmt.Sprintf("/projects/%d/members", id), map[string]any{"nombre": persona})
	if err != nil {
		return err
	}
	if rec.Code != http.StatusNotFound || c.mensajeDe(rec) != "usuario no encontrado" {
		return fmt.Errorf("se esperaba 404 «usuario no encontrado», se obtuvo %d %s", rec.Code, rec.Body.String())
	}
	lista, err := c.integrantes("Tienda")
	if err != nil {
		return err
	}
	for _, m := range lista {
		if m.Nombre == persona {
			return fmt.Errorf("%q no debía entrar al equipo", persona)
		}
	}
	return nil
}

// --- Identidad de la sesión ---

func (c *permisosContext) inicioLaVotacionComo(persona, titulo string) error {
	id, err := c.historia(titulo)
	if err != nil {
		return err
	}
	_, err = c.esperar(persona, http.MethodPost, fmt.Sprintf("/backlog/%d/votacion", id), nil, http.StatusCreated)
	return err
}

// votoIndicandoOtroUsuario manda en el cuerpo el usuario_id de otra persona:
// la API debe ignorarlo y registrar el voto a nombre de la sesión.
func (c *permisosContext) votoIndicandoOtroUsuario(valor int, titulo, otro string) error {
	id, err := c.historia(titulo)
	if err != nil {
		return err
	}
	otroID, err := c.sesion.idDe(otro)
	if err != nil {
		return err
	}
	_, err = c.esperar(c.sesion.actor(), http.MethodPost, fmt.Sprintf("/backlog/%d/votacion/votos", id),
		map[string]any{"valor": valor, "usuario_id": otroID}, http.StatusOK)
	return err
}

func (c *permisosContext) vistaDeLaVotacion(persona, titulo, query string) (votacionDeEscenario, error) {
	id, err := c.historia(titulo)
	if err != nil {
		return votacionDeEscenario{}, err
	}
	rec, err := c.esperar(persona, http.MethodGet, fmt.Sprintf("/backlog/%d/votacion%s", id, query), nil, http.StatusOK)
	if err != nil {
		return votacionDeEscenario{}, err
	}
	var vista votacionDeEscenario
	if err := json.Unmarshal(rec.Body.Bytes(), &vista); err != nil {
		return votacionDeEscenario{}, fmt.Errorf("votación inválida: %s", rec.Body.String())
	}
	return vista, nil
}

// consultaIndicandoOtroUsuario pide la votación con ?usuario_id= de otra
// persona: la API debe ignorarlo y mostrar la vista de la sesión.
func (c *permisosContext) consultaIndicandoOtroUsuario(persona, titulo, otro string) error {
	otroID, err := c.sesion.idDe(otro)
	if err != nil {
		return err
	}
	c.vista, err = c.vistaDeLaVotacion(persona, titulo, "?usuario_id="+strconv.FormatInt(otroID, 10))
	return err
}

func (c *permisosContext) votoANombreDe(valor int, persona string) error {
	id, err := c.sesion.idDe(persona)
	if err != nil {
		return err
	}
	propia, err := c.vistaDeLaVotacion(persona, "Carrito", "")
	if err != nil {
		return err
	}
	for _, m := range propia.Integrantes {
		if m.UsuarioID == id {
			if !m.YaVoto || m.Valor == nil || *m.Valor != valor {
				return fmt.Errorf("se esperaba el voto %d a nombre de %q, se obtuvo %+v", valor, persona, m)
			}
			return nil
		}
	}
	return fmt.Errorf("%q no figura en la votación", persona)
}

func (c *permisosContext) noVeElVotoOculto(quien, persona string) error {
	quienID, err := c.sesion.idDe(quien)
	if err != nil {
		return err
	}
	personaID, err := c.sesion.idDe(persona)
	if err != nil {
		return err
	}
	for _, m := range c.vista.Integrantes {
		switch m.UsuarioID {
		case personaID:
			if !m.YaVoto || m.Valor != nil {
				return fmt.Errorf("%q debía ver que %q votó, sin ver la carta: %+v", quien, persona, m)
			}
		case quienID:
			if m.YaVoto {
				return fmt.Errorf("el voto no debía quedar a nombre de %q: %+v", quien, m)
			}
		}
	}
	return nil
}

// InitializeScenarioPermisos registra los pasos de permisos_por_rol.feature
// (HU-21) y devuelve su contexto para los escenarios que reutilizan sus pasos.
func InitializeScenarioPermisos(ctx *godog.ScenarioContext, sesion *sesionDeEscenario) *permisosContext {
	c := &permisosContext{sesion: sesion}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		c.reiniciar()
		return goctx, nil
	})

	ctx.Step(`^que "([^"]*)" creó el proyecto "([^"]*)"$`, c.creoElProyecto)
	ctx.Step(`^que "([^"]*)" tiene cuenta y es Product Builder de "([^"]*)"$`, c.esProductBuilderDe)
	ctx.Step(`^que "([^"]*)" tiene la historia "([^"]*)"$`, c.tieneLaHistoria)
	ctx.Step(`^que "([^"]*)" tiene cuenta$`, c.tieneCuenta)
	ctx.Step(`^que "([^"]*)" tiene cuenta y no integra "[^"]*"$`, c.tieneCuenta)
	ctx.Step(`^que "([^"]*)" tiene cuenta y no integra ningún proyecto$`, c.tieneCuenta)
	ctx.Step(`^que no inicié sesión$`, c.noInicieSesion)

	ctx.Step(`^creo un Sprint en "([^"]*)"$`, c.crearSprint)
	ctx.Step(`^inicio la votación de "([^"]*)"$`, c.iniciarVotacion)
	ctx.Step(`^edito el proyecto "([^"]*)"$`, c.editarProyecto)
	ctx.Step(`^cada acción se realiza$`, c.cadaAccionSeRealiza)

	ctx.Step(`^intento crear un Sprint en "([^"]*)"$`, c.intento(c.crearSprint))
	ctx.Step(`^intento iniciar la votación de "([^"]*)"$`, c.intento(c.iniciarVotacion))
	ctx.Step(`^intento editar el proyecto "([^"]*)"$`, c.intento(c.editarProyecto))
	ctx.Step(`^intento agregar a "([^"]*)" al equipo de "([^"]*)"$`, c.intentoAgregarAlEquipo)
	ctx.Step(`^el sistema responde 403 con un mensaje que empieza con "([^"]*)"$`, c.responde403ConMensajeQueEmpieza)
	ctx.Step(`^no cambia nada$`, c.noCambiaNada)

	ctx.Step(`^consulto mis proyectos$`, c.consultoMisProyectos)
	ctx.Step(`^el sistema rechaza la petición con (\d+)$`, c.rechazaConCodigo)
	ctx.Step(`^abro el proyecto "([^"]*)"$`, c.abroElProyecto)
	ctx.Step(`^el sistema responde 403 con el mensaje "([^"]*)"$`, c.responde403ConMensaje)
	ctx.Step(`^veo solo "([^"]*)" con el rol "([^"]*)"$`, c.veoSoloConRol)
	ctx.Step(`^recibo una lista vacía$`, c.reciboListaVacia)

	ctx.Step(`^agrego a "([^"]*)" al equipo de "([^"]*)" como "([^"]*)"$`, c.agregarAlEquipoComo)
	ctx.Step(`^agrego a "([^"]*)" al equipo de "([^"]*)"$`, c.agregarAlEquipo)
	ctx.Step(`^"([^"]*)" sigue teniendo a "([^"]*)" como único Scrum Master$`, c.unicoScrumMaster)
	ctx.Step(`^"([^"]*)" queda como "([^"]*)" de "([^"]*)"$`, c.quedaComo)
	ctx.Step(`^si agrego a "([^"]*)", que no tiene cuenta, el sistema lo rechaza$`, c.agregarSinCuentaSeRechaza)

	ctx.Step(`^que "([^"]*)" inició la votación de "([^"]*)"$`, c.inicioLaVotacionComo)
	ctx.Step(`^voto (\d+) en "([^"]*)" indicando el usuario de "([^"]*)"$`, c.votoIndicandoOtroUsuario)
	ctx.Step(`^"([^"]*)" consulta la votación de "([^"]*)" indicando el usuario de "([^"]*)"$`, c.consultaIndicandoOtroUsuario)
	ctx.Step(`^el voto de (\d+) queda a nombre de "([^"]*)"$`, c.votoANombreDe)
	ctx.Step(`^"([^"]*)" no ve el voto oculto de "([^"]*)"$`, c.noVeElVotoOculto)
	return c
}
