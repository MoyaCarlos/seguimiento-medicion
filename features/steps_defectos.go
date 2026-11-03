package features

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// pasosDefectos guarda el estado de los escenarios de HU-08. Reutiliza la base
// y el proyecto reales del Background de HU-01 (no redefine ese paso) y arma
// los Sprints con los services de HU-05.
type pasosDefectos struct {
	base                   *scenarioContext
	historias              *repository.SQLiteBacklogRepository
	defectos               *repository.SQLiteDefectoRepository
	obtenerProyecto        *service.ObtenerProyecto
	registrar              *service.RegistrarDefecto
	resolver               *service.ResolverDefecto
	listar                 *service.ListarDefectos
	crearSprint            *service.CrearSprint
	iniciarSprint          *service.IniciarSprint
	cerrarSprint           *service.CerrarSprint
	hoy                    time.Time
	sprints                map[string]int64
	historiasPorTitulo     map[string]int64
	defectosPorDescripcion map[string]int64
	otroProyecto           int64
	ultimo                 domain.DefectoDetalle
	listado                []domain.DefectoDetalle
	err                    error
}

func (p *pasosDefectos) preparar() {
	p.hoy = hoyCivil()
	p.historias = repository.NewSQLiteBacklogRepository(p.base.db)
	p.defectos = repository.NewSQLiteDefectoRepository(p.base.db)
	sprints := repository.NewSQLiteSprintRepository(p.base.db)
	proyectos := repository.NewSQLiteProjectRepository(p.base.db)
	p.obtenerProyecto = service.NewObtenerProyecto(proyectos)
	p.registrar = service.NewRegistrarDefecto(sprints, p.historias, p.defectos)
	p.resolver = service.NewResolverDefecto(p.defectos, sprints)
	p.listar = service.NewListarDefectos(proyectos, p.defectos)
	p.crearSprint = service.NewCrearSprint(proyectos, sprints)
	p.iniciarSprint = service.NewIniciarSprint(sprints)
	p.cerrarSprint = service.NewCerrarSprint(sprints, p.historias)
	p.sprints, p.historiasPorTitulo, p.defectosPorDescripcion = map[string]int64{}, map[string]int64{}, map[string]int64{}
	p.otroProyecto, p.ultimo, p.listado, p.err = 0, domain.DefectoDetalle{}, nil, nil
}

func (p *pasosDefectos) crearHistoria(proyectoID int64, titulo string) (int64, error) {
	item, err := domain.NewBacklogItem(proyectoID, titulo, "Descripción de "+titulo, domain.PrioridadMust, nil)
	if err != nil {
		return 0, err
	}
	guardada, err := p.historias.Guardar(context.Background(), item)
	return guardada.ID, err
}

func (p *pasosDefectos) tieneHistoria(titulo string) error {
	id, err := p.crearHistoria(p.base.proyectoID, titulo)
	p.historiasPorTitulo[titulo] = id
	return err
}

// sprintEn crea un Sprint del proyecto y lo lleva hasta el estado pedido con
// los casos de uso de HU-05.
func (p *pasosDefectos) sprintEn(proyectoID int64, estado string) (int64, error) {
	ctx := context.Background()
	sprint, err := p.crearSprint.Ejecutar(ctx, proyectoID)
	if err != nil || estado == string(domain.SprintPendiente) {
		return sprint.ID, err
	}
	if _, err := p.iniciarSprint.Ejecutar(ctx, service.IniciarSprintInput{
		SprintID: sprint.ID, SprintGoal: "Objetivo",
		FechaInicio: p.hoy.AddDate(0, 0, -7), FechaFin: p.hoy.AddDate(0, 0, 7),
	}); err != nil {
		return 0, err
	}
	if estado == string(domain.SprintFinalizado) {
		_, err = p.cerrarSprint.Ejecutar(ctx, sprint.ID)
	}
	return sprint.ID, err
}

func (p *pasosDefectos) tieneSprint(nombre, estado string) error {
	id, err := p.sprintEn(p.base.proyectoID, estado)
	p.sprints[nombre] = id
	return err
}

func (p *pasosDefectos) cerroSprint(nombre string) error {
	_, err := p.cerrarSprint.Ejecutar(context.Background(), p.sprints[nombre])
	return err
}

// historiaCompletadaEn deja la historia Completada en un Sprint ya Finalizado:
// la completa con SQL mientras el Sprint está Activo y recién después lo cierra
// (CerrarSprint deja en el Sprint las completadas), sin pasar por HU-15.
func (p *pasosDefectos) historiaCompletadaEn(titulo, sprint string) error {
	id, err := p.sprintEn(p.base.proyectoID, string(domain.SprintActivo))
	if err != nil {
		return err
	}
	p.sprints[sprint] = id
	if _, err := p.base.db.Exec(
		`UPDATE historias SET estado = 'Completada', fecha_finalizacion = ?, sprint_id = ? WHERE id = ?`,
		p.hoy.Format(time.DateOnly), id, p.historiasPorTitulo[titulo],
	); err != nil {
		return err
	}
	return p.cerroSprint(sprint)
}

func (p *pasosDefectos) historiaSigueEnBacklog(titulo string) error {
	h, err := p.historias.ObtenerPorID(context.Background(), p.historiasPorTitulo[titulo])
	if err != nil {
		return err
	}
	if h.SprintID != nil {
		return fmt.Errorf("la historia %q está en el Sprint %d", titulo, *h.SprintID)
	}
	return nil
}

func (p *pasosDefectos) noTieneSprints() error { return nil }

// registrarEn sigue el orden del handler: el proyecto inexistente domina sobre
// el resto de las validaciones.
func (p *pasosDefectos) registrarEn(in service.RegistrarDefectoInput) {
	if _, p.err = p.obtenerProyecto.Ejecutar(context.Background(), in.ProyectoID); p.err != nil {
		return
	}
	p.ultimo, p.err = p.registrar.Ejecutar(context.Background(), in)
}

func (p *pasosDefectos) entrada(descripcion, severidad string, historiaID int64) service.RegistrarDefectoInput {
	return service.RegistrarDefectoInput{
		ProyectoID: p.base.proyectoID, HistoriaID: historiaID, Descripcion: descripcion, Severidad: severidad,
	}
}

func (p *pasosDefectos) registraDefecto(descripcion, severidad, historia string) error {
	p.registrarEn(p.entrada(descripcion, severidad, p.historiasPorTitulo[historia]))
	return nil
}

func (p *pasosDefectos) historiaPrincipal() int64 {
	return p.historiasPorTitulo["Carrito de compras"]
}

func (p *pasosDefectos) intentaConDescripcionVacia() error {
	p.registrarEn(p.entrada("", "Alta", p.historiaPrincipal()))
	return nil
}

func (p *pasosDefectos) intentaConSeveridad(severidad string) error {
	p.registrarEn(p.entrada("El total no se actualiza", severidad, p.historiaPrincipal()))
	return nil
}

func (p *pasosDefectos) intentaSinHistoria() error {
	p.registrarEn(p.entrada("El total no se actualiza", "Alta", 0))
	return nil
}

func (p *pasosDefectos) intentaConHistoriaInexistente(historiaID int) error {
	p.registrarEn(p.entrada("El total no se actualiza", "Alta", int64(historiaID)))
	return nil
}

// asegurarOtroProyecto crea, una sola vez, un segundo proyecto con un Sprint Activo.
func (p *pasosDefectos) asegurarOtroProyecto() error {
	if p.otroProyecto != 0 {
		return nil
	}
	proyecto, err := domain.NewProject("Otro proyecto", "", nil, nil)
	if err != nil {
		return err
	}
	guardado, err := repository.NewSQLiteProjectRepository(p.base.db).Guardar(context.Background(), proyecto)
	if err != nil {
		return err
	}
	p.otroProyecto = guardado.ID
	_, err = p.sprintEn(p.otroProyecto, string(domain.SprintActivo))
	return err
}

func (p *pasosDefectos) intentaConHistoriaAjena(titulo string) error {
	if err := p.asegurarOtroProyecto(); err != nil {
		return err
	}
	id, err := p.crearHistoria(p.otroProyecto, titulo)
	if err != nil {
		return err
	}
	p.registrarEn(p.entrada("El total no se actualiza", "Alta", id))
	return nil
}

func (p *pasosDefectos) intentaEnProyectoInexistente(proyectoID int) error {
	in := p.entrada("El total no se actualiza", "Alta", p.historiaPrincipal())
	in.ProyectoID = int64(proyectoID)
	p.registrarEn(in)
	return nil
}

func (p *pasosDefectos) seRegistro(descripcion, severidad string) error {
	p.registrarEn(p.entrada(descripcion, severidad, p.historiaPrincipal()))
	if p.err != nil {
		return fmt.Errorf("no se pudo sembrar el defecto %q: %w", descripcion, p.err)
	}
	p.defectosPorDescripcion[descripcion] = p.ultimo.ID
	return nil
}

func (p *pasosDefectos) otroProyectoConDefecto() error {
	if err := p.asegurarOtroProyecto(); err != nil {
		return err
	}
	historiaID, err := p.crearHistoria(p.otroProyecto, "Historia ajena")
	if err != nil {
		return err
	}
	_, err = p.registrar.Ejecutar(context.Background(), service.RegistrarDefectoInput{
		ProyectoID: p.otroProyecto, HistoriaID: historiaID, Descripcion: "Defecto ajeno", Severidad: "Baja",
	})
	return err
}

func (p *pasosDefectos) resolverDefecto(id int64) {
	var resuelto domain.DefectoDetalle
	resuelto, p.err = p.resolver.Ejecutar(context.Background(), id)
	if p.err == nil {
		p.ultimo = resuelto
	}
}

func (p *pasosDefectos) marcaResuelto() error {
	p.resolverDefecto(p.ultimo.ID)
	return nil
}

func (p *pasosDefectos) seMarcoResuelto(descripcion string) error {
	p.resolverDefecto(p.defectosPorDescripcion[descripcion])
	if p.err != nil {
		return fmt.Errorf("no se pudo resolver el defecto %q: %w", descripcion, p.err)
	}
	return nil
}

func (p *pasosDefectos) intentaResolverInexistente(id int) error {
	p.resolverDefecto(int64(id))
	return nil
}

// guardado relee el último defecto de la base, para verificar lo persistido.
func (p *pasosDefectos) guardado() (domain.DefectoDetalle, error) {
	if p.ultimo.ID == 0 {
		return domain.DefectoDetalle{}, fmt.Errorf("no hay un defecto registrado (error: %v)", p.err)
	}
	return p.defectos.ObtenerPorID(context.Background(), p.ultimo.ID)
}

func (p *pasosDefectos) quedaRegistrado(descripcion, severidad, historia string) error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error: %w", p.err)
	}
	d, err := p.guardado()
	if err != nil {
		return err
	}
	if d.Descripcion != descripcion || string(d.Severidad) != severidad ||
		d.HistoriaID != p.historiasPorTitulo[historia] || d.TituloHistoria != historia {
		return fmt.Errorf("defecto inesperado: %+v", d)
	}
	return nil
}

func (p *pasosDefectos) quedaPendiente(estado, sprint string) error {
	d, err := p.guardado()
	if err != nil {
		return err
	}
	if string(d.Estado) != estado || d.SprintDeteccionID != p.sprints[sprint] || d.SprintResolucionID != nil {
		return fmt.Errorf("se esperaba %s detectado en %q (%d) sin resolución, se obtuvo %+v", estado, sprint, p.sprints[sprint], d)
	}
	return nil
}

func (p *pasosDefectos) quedaResuelto(sprint string) error {
	d, err := p.guardado()
	if err != nil {
		return err
	}
	if d.Estado != domain.DefectoResuelto || d.SprintResolucionID == nil || *d.SprintResolucionID != p.sprints[sprint] {
		return fmt.Errorf("se esperaba Resuelto en %q (%d), se obtuvo %+v", sprint, p.sprints[sprint], d)
	}
	return nil
}

func (p *pasosDefectos) sigueDetectadoEn(sprint string) error {
	d, err := p.guardado()
	if err != nil {
		return err
	}
	if d.SprintDeteccionID != p.sprints[sprint] {
		return fmt.Errorf("se esperaba detectado en %q (%d), se obtuvo el Sprint %d", sprint, p.sprints[sprint], d.SprintDeteccionID)
	}
	return nil
}

func (p *pasosDefectos) rechazoConMensaje(mensaje string) error {
	var verr domain.ValidationError
	if !errors.As(p.err, &verr) || verr.Mensaje != mensaje {
		return fmt.Errorf("se esperaba el mensaje %q, se obtuvo %v", mensaje, p.err)
	}
	return nil
}

func (p *pasosDefectos) rechazoSobreCampo(campo string) error {
	var verr domain.ValidationError
	if !errors.As(p.err, &verr) || verr.Campo != campo {
		return fmt.Errorf("se esperaba un error sobre %q, se obtuvo %v", campo, p.err)
	}
	return nil
}

func (p *pasosDefectos) rechazoConError(esperado error) func() error {
	return func() error {
		if !errors.Is(p.err, esperado) {
			return fmt.Errorf("se esperaba %v, se obtuvo %v", esperado, p.err)
		}
		return nil
	}
}

func (p *pasosDefectos) sigueSinDefectos() error {
	defectos, err := p.listar.Ejecutar(context.Background(), p.base.proyectoID)
	if err != nil {
		return err
	}
	if len(defectos) != 0 {
		return fmt.Errorf("se esperaba que el proyecto siga sin defectos, tiene %d", len(defectos))
	}
	return nil
}

func (p *pasosDefectos) consultaDe(proyectoID int64) {
	p.listado, p.err = p.listar.Ejecutar(context.Background(), proyectoID)
}

func (p *pasosDefectos) consulta() error {
	p.consultaDe(p.base.proyectoID)
	return nil
}

func (p *pasosDefectos) consultaInexistente(proyectoID int) error {
	p.consultaDe(int64(proyectoID))
	return nil
}

func (p *pasosDefectos) veoDefectos(tabla *godog.Table) error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error: %w", p.err)
	}
	filas := tabla.Rows[1:]
	if len(p.listado) != len(filas) {
		return fmt.Errorf("se esperaban %d defectos, se obtuvieron %d", len(filas), len(p.listado))
	}
	for i, fila := range filas {
		d := p.listado[i]
		descripcion, severidad, estado := fila.Cells[0].Value, fila.Cells[1].Value, fila.Cells[2].Value
		historia, detectado, resuelto := fila.Cells[3].Value, fila.Cells[4].Value, fila.Cells[5].Value
		resolucionOK := d.SprintResolucionID == nil
		if resuelto != "" {
			resolucionOK = d.SprintResolucionID != nil && *d.SprintResolucionID == p.sprints[resuelto]
		}
		if d.Descripcion != descripcion || string(d.Severidad) != severidad || string(d.Estado) != estado ||
			d.TituloHistoria != historia || d.SprintDeteccionID != p.sprints[detectado] || !resolucionOK || d.CreadoEn.IsZero() {
			return fmt.Errorf("defecto %d: se esperaba %s/%s/%s/%s/%s/%s, se obtuvo %+v",
				i+1, descripcion, severidad, estado, historia, detectado, resuelto, d)
		}
	}
	return nil
}

func (p *pasosDefectos) listaVacia() error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error: %w", p.err)
	}
	if p.listado == nil || len(p.listado) != 0 {
		return fmt.Errorf("se esperaba una lista vacía, se obtuvo %+v", p.listado)
	}
	return nil
}

func InitializeScenarioDefectos(ctx *godog.ScenarioContext, base *scenarioContext) {
	p := &pasosDefectos{base: base}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		p.preparar()
		return goctx, nil
	})

	ctx.Step(`^que el proyecto tiene la historia "([^"]*)" en el Product Backlog$`, p.tieneHistoria)
	ctx.Step(`^que el proyecto tiene el Sprint "([^"]*)" en estado "(Pendiente|Activo|Finalizado)"$`, p.tieneSprint)
	ctx.Step(`^que el proyecto no tiene Sprints$`, p.noTieneSprints)
	ctx.Step(`^que la historia "([^"]*)" se completó en el Sprint "([^"]*)", ya finalizado$`, p.historiaCompletadaEn)
	ctx.Step(`^que la historia "([^"]*)" sigue en el Product Backlog$`, p.historiaSigueEnBacklog)
	ctx.Step(`^que se cerró el Sprint "([^"]*)"$`, p.cerroSprint)
	ctx.Step(`^(?:intento registrar|registro) el defecto "([^"]*)" con severidad "([^"]*)" en la historia "([^"]*)"$`, p.registraDefecto)
	ctx.Step(`^intento registrar un defecto con la descripción vacía$`, p.intentaConDescripcionVacia)
	ctx.Step(`^intento registrar un defecto con la severidad "([^"]*)"$`, p.intentaConSeveridad)
	ctx.Step(`^intento registrar un defecto sin historia relacionada$`, p.intentaSinHistoria)
	ctx.Step(`^intento registrar un defecto relacionado con la historia con identificador (\d+)$`, p.intentaConHistoriaInexistente)
	ctx.Step(`^intento registrar un defecto relacionado con la historia "([^"]*)" de otro proyecto$`, p.intentaConHistoriaAjena)
	ctx.Step(`^intento registrar un defecto en el proyecto con identificador (\d+), que no existe$`, p.intentaEnProyectoInexistente)
	ctx.Step(`^que se registró el defecto "([^"]*)" con severidad "([^"]*)"$`, p.seRegistro)
	ctx.Step(`^que otro proyecto tiene un defecto registrado$`, p.otroProyectoConDefecto)
	ctx.Step(`^(?:intento marcar|marco) el defecto como "Resuelto"$`, p.marcaResuelto)
	ctx.Step(`^que se marcó como "Resuelto" el defecto "([^"]*)"$`, p.seMarcoResuelto)
	ctx.Step(`^intento marcar como "Resuelto" el defecto con identificador (\d+)$`, p.intentaResolverInexistente)
	ctx.Step(`^el defecto queda registrado con la descripción "([^"]*)", severidad "([^"]*)" y la historia "([^"]*)"$`, p.quedaRegistrado)
	ctx.Step(`^el defecto queda en estado "(Pendiente|Resuelto)", detectado en el Sprint "([^"]*)" y sin Sprint de resolución$`, p.quedaPendiente)
	ctx.Step(`^el defecto queda en estado "Resuelto", resuelto en el Sprint "([^"]*)"$`, p.quedaResuelto)
	ctx.Step(`^el defecto sigue detectado en el Sprint "([^"]*)"$`, p.sigueDetectadoEn)
	ctx.Step(`^el registro se rechaza con el mensaje "([^"]*)"$`, p.rechazoConMensaje)
	ctx.Step(`^el registro se rechaza con un error sobre el campo "([^"]*)"$`, p.rechazoSobreCampo)
	ctx.Step(`^el registro se rechaza con error "proyecto no encontrado"$`, p.rechazoConError(domain.ErrProyectoNoEncontrado))
	ctx.Step(`^el proyecto sigue sin defectos$`, p.sigueSinDefectos)
	ctx.Step(`^la resolución se rechaza con el mensaje "([^"]*)"$`, p.rechazoConMensaje)
	ctx.Step(`^la resolución se rechaza con error "defecto no encontrado"$`, p.rechazoConError(domain.ErrDefectoNoEncontrado))
	ctx.Step(`^consulto los defectos del proyecto$`, p.consulta)
	ctx.Step(`^consulto los defectos del proyecto con identificador (\d+), que no existe$`, p.consultaInexistente)
	ctx.Step(`^veo estos defectos, en este orden, cada uno con su fecha de registro:$`, p.veoDefectos)
	ctx.Step(`^veo una lista de defectos vacía, sin error$`, p.listaVacia)
	ctx.Step(`^la consulta de defectos se rechaza con error "proyecto no encontrado"$`, p.rechazoConError(domain.ErrProyectoNoEncontrado))
}
