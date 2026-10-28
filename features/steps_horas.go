package features

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// pasosHoras guarda el estado de los escenarios de HU-07. Reutiliza la base y
// el proyecto reales del Background de HU-01 (no redefine ese paso). La
// asignación de historias a un Sprint llega con HU-06: acá se siembra con SQL.
type pasosHoras struct {
	base          *scenarioContext
	historias     *repository.SQLiteBacklogRepository
	horas         *repository.SQLiteRegistroHorasRepository
	imputar       *service.ImputarHoras
	consultar     *service.ConsultarEsfuerzoReal
	cerrar        *service.CerrarSprint
	hoy           time.Time
	integrantes   map[string]int64
	sprintActivo  int64
	sprintCerrado int64
	historiaID    int64
	otraID        int64
	ultima        domain.RegistroHorasDetalle
	esfuerzo      service.EsfuerzoReal
	minutosAntes  int
	err           error
}

func (p *pasosHoras) preparar() {
	p.hoy = hoyCivil()
	p.historias = repository.NewSQLiteBacklogRepository(p.base.db)
	p.horas = repository.NewSQLiteRegistroHorasRepository(p.base.db)
	sprints := repository.NewSQLiteSprintRepository(p.base.db)
	proyectos := repository.NewSQLiteProjectRepository(p.base.db)
	p.imputar = service.NewImputarHoras(p.historias, sprints, proyectos, p.horas)
	p.consultar = service.NewConsultarEsfuerzoReal(p.historias, p.horas)
	p.cerrar = service.NewCerrarSprint(sprints, p.historias)
	p.integrantes = map[string]int64{}
	p.sprintActivo, p.sprintCerrado, p.historiaID, p.otraID = 0, 0, 0, 0
	p.ultima, p.esfuerzo, p.minutosAntes, p.err = domain.RegistroHorasDetalle{}, service.EsfuerzoReal{}, 0, nil
}

func (p *pasosHoras) dia(nombre string) (time.Time, error) {
	switch nombre {
	case "hoy":
		return p.hoy, nil
	case "ayer":
		return p.hoy.AddDate(0, 0, -1), nil
	}
	return time.Time{}, fmt.Errorf("fecha %q no soportada: se espera hoy o ayer", nombre)
}

func minutosDe(horas string) (int, error) {
	valor, err := strconv.ParseFloat(horas, 64)
	if err != nil {
		return 0, err
	}
	return int(math.Round(valor * 60)), nil
}

func (p *pasosHoras) crearUsuario(nombre string) (int64, error) {
	var id int64
	err := p.base.db.QueryRow(
		`INSERT INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES (?, ?, ?) RETURNING id`,
		nombre, strings.ToLower(nombre), time.Now().UTC().Format(time.RFC3339),
	).Scan(&id)
	return id, err
}

func (p *pasosHoras) esIntegrante(nombre string) error {
	id, err := p.crearUsuario(nombre)
	if err != nil {
		return err
	}
	if _, err := p.base.db.Exec(
		`INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (?, ?, 'product_builder', ?)`,
		p.base.proyectoID, id, time.Now().UTC().Format(time.RFC3339),
	); err != nil {
		return err
	}
	p.integrantes[nombre] = id
	return nil
}

func (p *pasosHoras) crearSprint(estado domain.EstadoSprint, inicio, fin time.Time) (int64, error) {
	var id int64
	err := p.base.db.QueryRow(
		`INSERT INTO sprints (proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES (?, 'Objetivo', ?, ?, ?) RETURNING id`,
		p.base.proyectoID, inicio.Format(time.DateOnly), fin.Format(time.DateOnly), string(estado),
	).Scan(&id)
	return id, err
}

func (p *pasosHoras) tieneSprintActivo() error {
	id, err := p.crearSprint(domain.SprintActivo, p.hoy.AddDate(0, 0, -7), p.hoy.AddDate(0, 0, 7))
	p.sprintActivo = id
	return err
}

func (p *pasosHoras) crearHistoria(titulo string, sprintID *int64) (int64, error) {
	item, err := domain.NewBacklogItem(p.base.proyectoID, titulo, "Descripción de "+titulo, domain.Prioridad("M"), nil)
	if err != nil {
		return 0, err
	}
	guardada, err := p.historias.Guardar(context.Background(), item)
	if err != nil {
		return 0, err
	}
	if sprintID != nil {
		if _, err := p.base.db.Exec(`UPDATE historias SET sprint_id = ? WHERE id = ?`, *sprintID, guardada.ID); err != nil {
			return 0, err
		}
	}
	return guardada.ID, nil
}

func (p *pasosHoras) historiaEnSprintActivo() error {
	id, err := p.crearHistoria("Historia con horas", &p.sprintActivo)
	p.historiaID = id
	return err
}

func (p *pasosHoras) ejecutarImputacion(historiaID int64, integrante, horas, actividad string, fecha time.Time) {
	p.ultima, p.err = p.imputar.Ejecutar(context.Background(), service.ImputarHorasInput{
		HistoriaID:   historiaID,
		IntegranteID: p.integrantes[integrante],
		Fecha:        fecha.Format(time.DateOnly),
		Actividad:    actividad,
		Horas:        horas,
	}, time.Now())
}

func (p *pasosHoras) imputa(integrante, horas, actividad string) error {
	p.ejecutarImputacion(p.historiaID, integrante, horas, actividad, p.hoy)
	return nil
}

func (p *pasosHoras) imputo(integrante, horas string) error {
	p.ejecutarImputacion(p.historiaID, integrante, horas, "Trabajo previo", p.hoy)
	if p.err != nil {
		return fmt.Errorf("no se pudo sembrar la carga: %w", p.err)
	}
	return nil
}

func (p *pasosHoras) totalActual(historiaID int64) (int, error) {
	esfuerzo, err := p.consultar.Ejecutar(context.Background(), historiaID)
	return esfuerzo.TotalMinutos, err
}

// intentar guarda el acumulado previo para verificar después que no cambió.
func (p *pasosHoras) intentar(historiaID int64, integrante, horas, actividad string) error {
	if historiaID == p.historiaID {
		antes, err := p.totalActual(historiaID)
		if err != nil {
			return err
		}
		p.minutosAntes = antes
	}
	p.ejecutarImputacion(historiaID, integrante, horas, actividad, p.hoy)
	return nil
}

func (p *pasosHoras) intentaImputarHoras(integrante, horas string) error {
	return p.intentar(p.historiaID, integrante, horas, "Trabajo")
}

func (p *pasosHoras) intentaConActividad(integrante, actividad string) error {
	textos := map[string]string{"vacía": "", "solo con espacios": "   "}
	return p.intentar(p.historiaID, integrante, "1", textos[actividad])
}

func (p *pasosHoras) intentaSiendoAjeno(nombre string) error {
	id, err := p.crearUsuario(nombre)
	if err != nil {
		return err
	}
	p.integrantes[nombre] = id
	return p.intentar(p.historiaID, nombre, "1", "Trabajo")
}

func (p *pasosHoras) intentaEnInexistente(integrante string, historiaID int) error {
	return p.intentar(int64(historiaID), integrante, "1", "Trabajo")
}

func (p *pasosHoras) otraEnBacklog() error {
	id, err := p.crearHistoria("Otra historia", nil)
	p.otraID = id
	return err
}

func (p *pasosHoras) otraEnSprint(estado string) error {
	var (
		sprintID int64
		err      error
	)
	if estado == string(domain.SprintPendiente) {
		err = p.base.db.QueryRow(`INSERT INTO sprints (proyecto_id, estado) VALUES (?, 'Pendiente') RETURNING id`, p.base.proyectoID).Scan(&sprintID)
	} else {
		sprintID, err = p.crearSprint(domain.EstadoSprint(estado), p.hoy.AddDate(0, 0, -30), p.hoy.AddDate(0, 0, -16))
	}
	if err != nil {
		return err
	}
	id, err := p.crearHistoria("Otra historia", &sprintID)
	p.otraID = id
	return err
}

func (p *pasosHoras) intentaEnOtra(integrante string) error {
	return p.intentar(p.otraID, integrante, "1", "Trabajo")
}

func (p *pasosHoras) cierraYPlanificaEnOtro() error {
	if _, err := p.cerrar.Ejecutar(context.Background(), p.sprintActivo); err != nil {
		return fmt.Errorf("no se pudo cerrar el Sprint: %w", err)
	}
	nuevo, err := p.crearSprint(domain.SprintActivo, p.hoy, p.hoy.AddDate(0, 0, 14))
	if err != nil {
		return err
	}
	if _, err := p.base.db.Exec(`UPDATE historias SET sprint_id = ? WHERE id = ?`, nuevo, p.historiaID); err != nil {
		return err
	}
	p.sprintCerrado, p.sprintActivo = p.sprintActivo, nuevo
	return nil
}

func (p *pasosHoras) cargaRegistrada(integrante, actividad string, minutos int) error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error: %w", p.err)
	}
	r := p.ultima
	if r.ID == 0 || r.NombreIntegrante != integrante || r.UsuarioID != p.integrantes[integrante] ||
		!r.Fecha.Equal(p.hoy) || r.Actividad != actividad || r.Minutos != minutos ||
		r.SprintID != p.sprintActivo || r.HistoriaID != p.historiaID {
		return fmt.Errorf("carga inesperada: %+v", r)
	}
	guardadas, err := p.horas.ListarPorHistoria(context.Background(), p.historiaID)
	if err != nil {
		return err
	}
	for _, g := range guardadas {
		if g.ID == r.ID {
			return nil
		}
	}
	return fmt.Errorf("la carga %d no quedó persistida", r.ID)
}

func (p *pasosHoras) acumuladoEs(horas string) error {
	esperado, err := minutosDe(horas)
	if err != nil {
		return err
	}
	total, err := p.totalActual(p.historiaID)
	if err != nil {
		return err
	}
	if total != esperado {
		return fmt.Errorf("se esperaban %d minutos acumulados, hay %d", esperado, total)
	}
	return nil
}

func (p *pasosHoras) siguenEnSprintCerrado(horas int) error {
	cargas, err := p.horas.ListarPorHistoria(context.Background(), p.historiaID)
	if err != nil {
		return err
	}
	minutos := 0
	for _, c := range cargas {
		if c.SprintID != p.sprintCerrado {
			return fmt.Errorf("la carga %d quedó en el Sprint %d, no en el cerrado %d", c.ID, c.SprintID, p.sprintCerrado)
		}
		minutos += c.Minutos
	}
	if minutos != horas*60 {
		return fmt.Errorf("se esperaban %d minutos en el Sprint cerrado, hay %d", horas*60, minutos)
	}
	return nil
}

func (p *pasosHoras) rechazoSobreCampo(campo string) error {
	var verr domain.ValidationError
	if !errors.As(p.err, &verr) || verr.Campo != campo {
		return fmt.Errorf("se esperaba un error sobre %q, se obtuvo %v", campo, p.err)
	}
	return nil
}

func (p *pasosHoras) rechazoConMensaje(mensaje string) error {
	var verr domain.ValidationError
	if !errors.As(p.err, &verr) || verr.Mensaje != mensaje {
		return fmt.Errorf("se esperaba el mensaje %q, se obtuvo %v", mensaje, p.err)
	}
	return nil
}

func (p *pasosHoras) acumuladoNoCambia() error {
	total, err := p.totalActual(p.historiaID)
	if err != nil {
		return err
	}
	if total != p.minutosAntes {
		return fmt.Errorf("el acumulado cambió de %d a %d minutos", p.minutosAntes, total)
	}
	return nil
}

func (p *pasosHoras) otraSinHoras() error {
	total, err := p.totalActual(p.otraID)
	if err != nil {
		return err
	}
	if total != 0 {
		return fmt.Errorf("la otra historia tiene %d minutos imputados", total)
	}
	return nil
}

func (p *pasosHoras) imputacionHistoriaNoEncontrada() error {
	if !errors.Is(p.err, domain.ErrHistoriaNoEncontrada) {
		return fmt.Errorf("se esperaba ErrHistoriaNoEncontrada, se obtuvo %v", p.err)
	}
	return nil
}

func (p *pasosHoras) historiaTieneCargas(tabla *godog.Table) error {
	for _, fila := range tabla.Rows[1:] {
		integrante, horas, actividad := fila.Cells[0].Value, fila.Cells[1].Value, fila.Cells[2].Value
		fecha, err := p.dia(fila.Cells[3].Value)
		if err != nil {
			return err
		}
		p.ejecutarImputacion(p.historiaID, integrante, horas, actividad, fecha)
		if p.err != nil {
			return fmt.Errorf("no se pudo sembrar la carga de %s: %w", integrante, p.err)
		}
	}
	return nil
}

func (p *pasosHoras) consultoEsfuerzoDe(historiaID int64) {
	p.esfuerzo, p.err = p.consultar.Ejecutar(context.Background(), historiaID)
}

func (p *pasosHoras) consultoEsfuerzo() error {
	p.consultoEsfuerzoDe(p.historiaID)
	return nil
}

func (p *pasosHoras) consultoEsfuerzoInexistente(historiaID int) error {
	p.consultoEsfuerzoDe(int64(historiaID))
	return nil
}

func (p *pasosHoras) totalEs(horas string) error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error: %w", p.err)
	}
	esperado, err := minutosDe(horas)
	if err != nil {
		return err
	}
	if p.esfuerzo.TotalMinutos != esperado {
		return fmt.Errorf("se esperaba un total de %d minutos, se obtuvo %d", esperado, p.esfuerzo.TotalMinutos)
	}
	return nil
}

func (p *pasosHoras) veoCargasEnOrden(tabla *godog.Table) error {
	filas := tabla.Rows[1:]
	if len(p.esfuerzo.Registros) != len(filas) {
		return fmt.Errorf("se esperaban %d cargas, se obtuvieron %d", len(filas), len(p.esfuerzo.Registros))
	}
	for i, fila := range filas {
		minutos, err := minutosDe(fila.Cells[1].Value)
		if err != nil {
			return err
		}
		fecha, err := p.dia(fila.Cells[3].Value)
		if err != nil {
			return err
		}
		r := p.esfuerzo.Registros[i]
		if r.NombreIntegrante != fila.Cells[0].Value || r.Minutos != minutos || r.Actividad != fila.Cells[2].Value ||
			!r.Fecha.Equal(fecha) || r.SprintID != p.sprintActivo {
			return fmt.Errorf("carga %d: se esperaba %s/%s/%s/%s en el Sprint %d, se obtuvo %+v",
				i+1, fila.Cells[0].Value, fila.Cells[1].Value, fila.Cells[2].Value, fila.Cells[3].Value, p.sprintActivo, r)
		}
	}
	return nil
}

func (p *pasosHoras) noVeoCargas() error {
	if p.esfuerzo.Registros == nil || len(p.esfuerzo.Registros) != 0 {
		return fmt.Errorf("se esperaba un detalle vacío, se obtuvo %+v", p.esfuerzo.Registros)
	}
	return nil
}

func (p *pasosHoras) consultaHistoriaNoEncontrada() error {
	if !errors.Is(p.err, domain.ErrHistoriaNoEncontrada) {
		return fmt.Errorf("se esperaba ErrHistoriaNoEncontrada, se obtuvo %v", p.err)
	}
	return nil
}

func InitializeScenarioHoras(ctx *godog.ScenarioContext, base *scenarioContext) {
	p := &pasosHoras{base: base}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		p.preparar()
		return goctx, nil
	})

	ctx.Step(`^que "([^"]*)" es integrante del proyecto$`, p.esIntegrante)
	ctx.Step(`^que el proyecto tiene un Sprint activo$`, p.tieneSprintActivo)
	ctx.Step(`^que existe una historia asignada al Sprint activo$`, p.historiaEnSprintActivo)
	ctx.Step(`^"([^"]*)" imputa ([\d.]+) horas a la historia con la fecha de hoy y la actividad "([^"]*)"$`, p.imputa)
	ctx.Step(`^la carga queda registrada con integrante "([^"]*)", la fecha de hoy, la actividad "([^"]*)", (\d+) minutos y el Sprint activo$`, p.cargaRegistrada)
	ctx.Step(`^el esfuerzo acumulado de la historia es de ([\d.]+) horas$`, p.acumuladoEs)
	ctx.Step(`^que "([^"]*)" imputó ([\d.]+) horas a la historia$`, p.imputo)
	ctx.Step(`^se cierra el Sprint activo y la historia se planifica en un Sprint nuevo$`, p.cierraYPlanificaEnOtro)
	ctx.Step(`^las (\d+) horas siguen atribuidas al Sprint cerrado$`, p.siguenEnSprintCerrado)
	ctx.Step(`^"([^"]*)" intenta imputar "([^"]*)" horas a la historia$`, p.intentaImputarHoras)
	ctx.Step(`^la carga se rechaza con un error sobre el campo "([^"]*)"$`, p.rechazoSobreCampo)
	ctx.Step(`^el esfuerzo acumulado de la historia no cambia$`, p.acumuladoNoCambia)
	ctx.Step(`^"([^"]*)" intenta imputar horas a la historia con una actividad (vacía|solo con espacios)$`, p.intentaConActividad)
	ctx.Step(`^"([^"]*)", que no es integrante del proyecto, intenta imputar horas a la historia$`, p.intentaSiendoAjeno)
	ctx.Step(`^que existe otra historia en el Product Backlog$`, p.otraEnBacklog)
	ctx.Step(`^que existe otra historia en un Sprint en estado "([^"]*)"$`, p.otraEnSprint)
	ctx.Step(`^"([^"]*)" intenta imputar horas a esa otra historia$`, p.intentaEnOtra)
	ctx.Step(`^la carga se rechaza con el mensaje "([^"]*)"$`, p.rechazoConMensaje)
	ctx.Step(`^esa otra historia no tiene horas imputadas$`, p.otraSinHoras)
	ctx.Step(`^"([^"]*)" intenta imputar horas a la historia con identificador (\d+)$`, p.intentaEnInexistente)
	ctx.Step(`^la carga se rechaza con error "historia no encontrada"$`, p.imputacionHistoriaNoEncontrada)
	ctx.Step(`^que la historia tiene estas cargas, cargadas en este orden:$`, p.historiaTieneCargas)
	ctx.Step(`^consulto el esfuerzo real de la historia$`, p.consultoEsfuerzo)
	ctx.Step(`^consulto el esfuerzo real de la historia con identificador (\d+)$`, p.consultoEsfuerzoInexistente)
	ctx.Step(`^el total es de ([\d.]+) horas$`, p.totalEs)
	ctx.Step(`^veo estas cargas en este orden, por fecha y luego por orden de carga, todas en el Sprint activo:$`, p.veoCargasEnOrden)
	ctx.Step(`^no veo ninguna carga$`, p.noVeoCargas)
	ctx.Step(`^la consulta del esfuerzo se rechaza con error "historia no encontrada"$`, p.consultaHistoriaNoEncontrada)
}
