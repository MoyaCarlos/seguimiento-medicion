package features

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// pasosVotacion guarda el estado de los escenarios de HU-02 sobre la base del
// contexto de HU-01. El equipo y la historia los crea su propio Antecedentes.
type pasosVotacion struct {
	base       *scenarioContext
	votaciones *repository.SQLiteVotacionRepository
	iniciar    *service.IniciarVotacion
	votar      *service.Votar
	revelar    *service.RevelarVotacion
	consultar  *service.ConsultarVotacion
	usuarios   map[string]int64
	proyectoID int64
	historiaID int64
	antes      domain.RondaVotacion
	err        error
}

func (p *pasosVotacion) preparar() {
	historias := repository.NewSQLiteBacklogRepository(p.base.db)
	proyectos := repository.NewSQLiteProjectRepository(p.base.db)
	p.votaciones = repository.NewSQLiteVotacionRepository(p.base.db)
	p.iniciar = service.NewIniciarVotacion(historias, proyectos, p.votaciones)
	p.votar = service.NewVotar(historias, proyectos, p.votaciones)
	p.revelar = service.NewRevelarVotacion(historias, proyectos, p.votaciones)
	p.consultar = service.NewConsultarVotacion(historias, proyectos, p.votaciones)
	p.usuarios = map[string]int64{}
	p.proyectoID, p.historiaID = 0, 0
	p.antes, p.err = domain.RondaVotacion{}, nil
}

func (p *pasosVotacion) crearUsuario(nombre string) (int64, error) {
	u, err := domain.NewUser(nombre)
	if err != nil {
		return 0, err
	}
	guardado, err := repository.NewSQLiteUserRepository(p.base.db).Guardar(context.Background(), u)
	if err != nil {
		return 0, err
	}
	p.usuarios[nombre] = guardado.ID
	return guardado.ID, nil
}

func (p *pasosVotacion) usuario(nombre string) (int64, error) {
	id, ok := p.usuarios[nombre]
	if !ok {
		return 0, fmt.Errorf("el escenario no creó a %q", nombre)
	}
	return id, nil
}

func (p *pasosVotacion) equipo(scrumMaster, builder1, builder2 string) error {
	ctx := context.Background()
	carlosID, err := p.crearUsuario(scrumMaster)
	if err != nil {
		return err
	}
	proyecto, err := domain.NewProject("Software Metrics & Estimation", "", nil, nil)
	if err != nil {
		return err
	}
	proyectos := repository.NewSQLiteProjectRepository(p.base.db)
	guardado, err := proyectos.GuardarConScrumMaster(ctx, proyecto, carlosID)
	if err != nil {
		return err
	}
	p.proyectoID = guardado.ID

	for _, nombre := range []string{builder1, builder2} {
		id, err := p.crearUsuario(nombre)
		if err != nil {
			return err
		}
		m, err := domain.NewMembership(p.proyectoID, id, domain.RolProductBuilder)
		if err != nil {
			return err
		}
		if _, err := proyectos.AgregarIntegrante(ctx, m); err != nil {
			return err
		}
	}
	return nil
}

func (p *pasosVotacion) historiaEnEstado(estado string) error {
	if estado != string(domain.EstadoNueva) {
		return fmt.Errorf("toda historia nace %q, no %q", domain.EstadoNueva, estado)
	}
	item, err := domain.NewBacklogItem(p.proyectoID, "Iniciar sesión", "Como integrante quiero iniciar sesión", domain.PrioridadMust, nil)
	if err != nil {
		return err
	}
	guardada, err := repository.NewSQLiteBacklogRepository(p.base.db).Guardar(context.Background(), item)
	if err != nil {
		return err
	}
	p.historiaID = guardada.ID
	return nil
}

// historiaCompletada siembra el estado por SQL (la transición es de HU-11),
// con el Sprint y la fecha de finalización que exigen los CHECK de §4.4.
func (p *pasosVotacion) historiaCompletada(estado string) error {
	if estado != string(domain.EstadoCompletada) {
		return fmt.Errorf("el paso solo siembra historias %q", domain.EstadoCompletada)
	}
	resultado, err := p.base.db.Exec(
		`INSERT INTO sprints (proyecto_id, estado, sprint_goal, fecha_inicio, fecha_fin) VALUES (?, 'Activo', 'Objetivo', '2026-10-01', '2026-10-15')`,
		p.proyectoID)
	if err != nil {
		return err
	}
	sprintID, err := resultado.LastInsertId()
	if err != nil {
		return err
	}
	if _, err := p.base.db.Exec(
		`UPDATE historias SET estado = 'Completada', sprint_id = ?, fecha_finalizacion = '2026-10-10' WHERE id = ?`,
		sprintID, p.historiaID); err != nil {
		return err
	}
	return p.sinVotacion()
}

func (p *pasosVotacion) inicio(nombre string) error {
	return p.inicioDe(nombre, p.historiaID)
}

func (p *pasosVotacion) inicioDe(nombre string, historiaID int64) error {
	id, err := p.usuario(nombre)
	if err != nil {
		return err
	}
	_, p.err = p.iniciar.Ejecutar(context.Background(), historiaID, id)
	return nil
}

func (p *pasosVotacion) inicioDeInexistente(nombre string, historiaID int) error {
	return p.inicioDe(nombre, int64(historiaID))
}

func (p *pasosVotacion) inicioLaVotacion(nombre string) error {
	if err := p.inicio(nombre); err != nil {
		return err
	}
	return p.err
}

func (p *pasosVotacion) rondaActual() (domain.RondaVotacion, error) {
	return p.votaciones.ObtenerRondaActual(context.Background(), p.historiaID)
}

// fotografiar guarda la ronda antes de la acción, para comprobar que un
// rechazo no la cambia.
func (p *pasosVotacion) fotografiar() error {
	r, err := p.rondaActual()
	p.antes = r
	return err
}

func (p *pasosVotacion) vota(nombre string, carta int) error {
	if err := p.fotografiar(); err != nil {
		return err
	}
	id, err := p.usuario(nombre)
	if err != nil {
		return err
	}
	_, p.err = p.votar.Ejecutar(context.Background(), service.VotarInput{HistoriaID: p.historiaID, UsuarioID: id, Valor: carta})
	return nil
}

func (p *pasosVotacion) votaQuienNoIntegra(nombre string, carta int) error {
	if _, err := p.crearUsuario(nombre); err != nil {
		return err
	}
	return p.vota(nombre, carta)
}

func (p *pasosVotacion) votoLaCarta(nombre string, carta int) error {
	if err := p.vota(nombre, carta); err != nil {
		return err
	}
	return p.err
}

func (p *pasosVotacion) votaronEstasCartas(tabla *godog.Table) error {
	for _, fila := range tabla.Rows[1:] {
		carta, err := strconv.Atoi(fila.Cells[1].Value)
		if err != nil {
			return err
		}
		if err := p.votoLaCarta(fila.Cells[0].Value, carta); err != nil {
			return err
		}
	}
	return nil
}

func (p *pasosVotacion) presionaRevelar(nombre string) error {
	if err := p.fotografiar(); err != nil {
		return err
	}
	id, err := p.usuario(nombre)
	if err != nil {
		return err
	}
	_, p.err = p.revelar.Ejecutar(context.Background(), p.historiaID, id, time.Now())
	return nil
}

func (p *pasosVotacion) revelo(nombre string) error {
	if err := p.presionaRevelar(nombre); err != nil {
		return err
	}
	return p.err
}

func (p *pasosVotacion) vistaDe(nombre string) (service.VistaVotacion, error) {
	var id int64
	if nombre != "" {
		var err error
		if id, err = p.usuario(nombre); err != nil {
			return service.VistaVotacion{}, err
		}
	}
	return p.consultar.Ejecutar(context.Background(), p.historiaID, id)
}

func (p *pasosVotacion) votoDe(vista service.VistaVotacion, nombre string) (service.VotoDeIntegrante, error) {
	id, err := p.usuario(nombre)
	if err != nil {
		return service.VotoDeIntegrante{}, err
	}
	for _, v := range vista.Integrantes {
		if v.UsuarioID == id {
			return v, nil
		}
	}
	return service.VotoDeIntegrante{}, fmt.Errorf("%q no figura en la votación", nombre)
}

func (p *pasosVotacion) votoRegistradoEnRonda(nombre string, numero int) error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error: %v", p.err)
	}
	return p.tieneUnSoloVoto(nombre, -1, numero)
}

func (p *pasosVotacion) unSoloVotoConCarta(nombre string, carta int) error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error: %v", p.err)
	}
	return p.tieneUnSoloVoto(nombre, carta, 1)
}

// tieneUnSoloVoto comprueba en la base el voto del integrante; carta -1 no
// mira el valor.
func (p *pasosVotacion) tieneUnSoloVoto(nombre string, carta, numero int) error {
	id, err := p.usuario(nombre)
	if err != nil {
		return err
	}
	r, err := p.rondaActual()
	if err != nil {
		return err
	}
	if r.Numero != numero {
		return fmt.Errorf("se esperaba la ronda %d, la actual es la %d", numero, r.Numero)
	}
	var suyos []domain.Voto
	for _, v := range r.Votos {
		if v.UsuarioID == id {
			suyos = append(suyos, v)
		}
	}
	if len(suyos) != 1 || (carta != -1 && suyos[0].Valor != carta) {
		return fmt.Errorf("se esperaba un voto de %q con la carta %d, hay %+v", nombre, carta, suyos)
	}
	return nil
}

func (p *pasosVotacion) veVotoOculto(quienConsulta, nombre string) error {
	vista, err := p.vistaDe(quienConsulta)
	if err != nil {
		return err
	}
	v, err := p.votoDe(vista, nombre)
	if err != nil {
		return err
	}
	if !v.YaVoto || v.Valor != nil {
		return fmt.Errorf("%q debía ver que %q votó con el valor oculto, ve %+v", quienConsulta, nombre, v)
	}
	return nil
}

func (p *pasosVotacion) veSuPropiaCarta(nombre string, carta int) error {
	vista, err := p.vistaDe(nombre)
	if err != nil {
		return err
	}
	v, err := p.votoDe(vista, nombre)
	if err != nil {
		return err
	}
	if v.Valor == nil || *v.Valor != carta {
		return fmt.Errorf("%q debía ver su carta %d, ve %+v", nombre, carta, v)
	}
	return nil
}

func (p *pasosVotacion) seVeCadaCarta(tabla *godog.Table) error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error al revelar: %v", p.err)
	}
	vista, err := p.vistaDe("")
	if err != nil {
		return err
	}
	for _, fila := range tabla.Rows[1:] {
		nombre := fila.Cells[0].Value
		carta, err := strconv.Atoi(fila.Cells[1].Value)
		if err != nil {
			return err
		}
		v, err := p.votoDe(vista, nombre)
		if err != nil {
			return err
		}
		if v.Valor == nil || *v.Valor != carta {
			return fmt.Errorf("se esperaba ver la carta %d de %q, se ve %+v", carta, nombre, v)
		}
	}
	return nil
}

func (p *pasosVotacion) promedioSugerido(enteros, decimales string, carta int) error {
	exacto, err := strconv.ParseFloat(enteros+"."+decimales, 64)
	if err != nil {
		return err
	}
	vista, err := p.vistaDe("")
	if err != nil {
		return err
	}
	esperado := domain.Promedio{Exacto: exacto, CartaSugerida: carta}
	if vista.Promedio == nil || *vista.Promedio != esperado {
		return fmt.Errorf("se esperaba el promedio %+v, se obtuvo %+v", esperado, vista.Promedio)
	}
	return nil
}

func (p *pasosVotacion) rechazadoConAdvertenciaSobre(campo string) error {
	var verr domain.ValidationError
	if !errors.As(p.err, &verr) {
		return fmt.Errorf("se esperaba una advertencia de validación, se obtuvo %v", p.err)
	}
	if verr.Campo != campo {
		return fmt.Errorf("se esperaba una advertencia sobre %q, se obtuvo %q", campo, verr.Campo)
	}
	return nil
}

func (p *pasosVotacion) esperarError(esperado error) error {
	if !errors.Is(p.err, esperado) {
		return fmt.Errorf("se esperaba %q, se obtuvo %v", esperado, p.err)
	}
	return nil
}

func (p *pasosVotacion) noEsIntegrante(string) error {
	return p.esperarError(domain.ErrNoEsIntegrante)
}

func (p *pasosVotacion) noEsScrumMaster(string) error {
	return p.esperarError(domain.ErrNoEsScrumMaster)
}

func (p *pasosVotacion) yaTieneVotacion() error {
	return p.esperarError(domain.ErrVotacionYaIniciada)
}

func (p *pasosVotacion) historiaNoEncontrada() error {
	return p.esperarError(domain.ErrHistoriaNoEncontrada)
}

func (p *pasosVotacion) rondaSinCambios() error {
	r, err := p.rondaActual()
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(r, p.antes) {
		return fmt.Errorf("la ronda no debía cambiar: antes %+v, ahora %+v", p.antes, r)
	}
	return nil
}

func (p *pasosVotacion) votosReveladosSinCambios() error {
	if err := p.rondaSinCambios(); err != nil {
		return err
	}
	if !p.antes.Revelada() {
		return errors.New("la ronda debía seguir revelada")
	}
	return nil
}

func (p *pasosVotacion) siguenOcultos() error {
	r, err := p.rondaActual()
	if err != nil {
		return err
	}
	if r.Revelada() {
		return fmt.Errorf("los votos debían seguir ocultos, la ronda se reveló en %v", r.ReveladaEn)
	}
	return nil
}

func (p *pasosVotacion) sinVotacion() error {
	return p.cantidadDeRondas(0)
}

func (p *pasosVotacion) unaSolaRonda() error {
	return p.cantidadDeRondas(1)
}

func (p *pasosVotacion) cantidadDeRondas(esperadas int) error {
	var total int
	if err := p.base.db.QueryRow("SELECT COUNT(*) FROM rondas_votacion WHERE historia_id = ?", p.historiaID).Scan(&total); err != nil {
		return err
	}
	if total != esperadas {
		return fmt.Errorf("se esperaban %d rondas de votación, hay %d", esperadas, total)
	}
	return nil
}

func (p *pasosVotacion) rondaUnoAbierta() error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error: %v", p.err)
	}
	r, err := p.rondaActual()
	if err != nil {
		return err
	}
	if r.Numero != 1 || len(r.Votos) != 0 || r.Revelada() {
		return fmt.Errorf("se esperaba la ronda 1 sin votos ni revelar, se obtuvo %+v", r)
	}
	return nil
}

// InitializeScenarioVotacion registra los pasos de HU-02 sobre la base de
// HU-01.
func InitializeScenarioVotacion(ctx *godog.ScenarioContext, base *scenarioContext) {
	p := &pasosVotacion{base: base}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		p.preparar()
		return goctx, nil
	})

	ctx.Step(`^que el proyecto de la votación tiene a "([^"]*)" como Scrum Master y a "([^"]*)" y "([^"]*)" como Product Builders$`, p.equipo)
	ctx.Step(`^que el proyecto de la votación tiene una historia en estado "([^"]*)"$`, p.historiaEnEstado)
	ctx.Step(`^que la historia está en estado "([^"]*)" y todavía no tiene votación$`, p.historiaCompletada)
	ctx.Step(`^que la historia todavía no tiene votación$`, p.sinVotacion)
	ctx.Step(`^que "([^"]*)" inició la votación de la historia$`, p.inicioLaVotacion)
	ctx.Step(`^que "([^"]*)" votó la carta (\d+)$`, p.votoLaCarta)
	ctx.Step(`^que votaron estas cartas:$`, p.votaronEstasCartas)
	ctx.Step(`^que "([^"]*)" reveló los votos$`, p.revelo)

	ctx.Step(`^"([^"]*)" selecciona la carta (\d+) y confirma su voto$`, p.vota)
	ctx.Step(`^"([^"]*)", que no es integrante del proyecto, selecciona la carta (\d+) y confirma su voto$`, p.votaQuienNoIntegra)
	ctx.Step(`^"([^"]*)" presiona "Revelar"$`, p.presionaRevelar)
	ctx.Step(`^"([^"]*)" inicia la votación de la historia$`, p.inicio)
	ctx.Step(`^"([^"]*)" inicia la votación de la historia con identificador (\d+), que no existe$`, p.inicioDeInexistente)

	ctx.Step(`^el voto de "([^"]*)" queda registrado en la ronda (\d+)$`, p.votoRegistradoEnRonda)
	ctx.Step(`^"([^"]*)" ve que "([^"]*)" votó, con el valor "Oculto"$`, p.veVotoOculto)
	ctx.Step(`^"([^"]*)" ve su propia carta (\d+)$`, p.veSuPropiaCarta)
	ctx.Step(`^"([^"]*)" tiene un solo voto en la ronda, con la carta (\d+)$`, p.unSoloVotoConCarta)
	ctx.Step(`^el voto se rechaza con una advertencia sobre el campo "([^"]*)"$`, p.rechazadoConAdvertenciaSobre)
	ctx.Step(`^la revelación se rechaza con una advertencia sobre el campo "([^"]*)"$`, p.rechazadoConAdvertenciaSobre)
	ctx.Step(`^el inicio de la votación se rechaza con una advertencia sobre el campo "([^"]*)"$`, p.rechazadoConAdvertenciaSobre)
	ctx.Step(`^el voto se rechaza porque "([^"]*)" no es integrante del proyecto$`, p.noEsIntegrante)
	ctx.Step(`^la revelación se rechaza porque "([^"]*)" no es Scrum Master del proyecto$`, p.noEsScrumMaster)
	ctx.Step(`^el inicio de la votación se rechaza porque "([^"]*)" no es Scrum Master del proyecto$`, p.noEsScrumMaster)
	ctx.Step(`^el inicio de la votación se rechaza porque la historia ya tiene una votación$`, p.yaTieneVotacion)
	ctx.Step(`^el inicio de la votación se rechaza con error "historia no encontrada"$`, p.historiaNoEncontrada)
	ctx.Step(`^la ronda de votación no cambia$`, p.rondaSinCambios)
	ctx.Step(`^se ve la carta elegida por cada integrante:$`, p.seVeCadaCarta)
	ctx.Step(`^el promedio sugerido es (\d+),(\d+), con la carta sugerida (\d+)$`, p.promedioSugerido)
	ctx.Step(`^los votos revelados no cambian$`, p.votosReveladosSinCambios)
	ctx.Step(`^los votos siguen ocultos$`, p.siguenOcultos)
	ctx.Step(`^la historia tiene abierta la ronda 1, sin votos y sin revelar$`, p.rondaUnoAbierta)
	ctx.Step(`^la historia sigue con una sola ronda de votación$`, p.unaSolaRonda)
	ctx.Step(`^la historia sigue sin rondas de votación$`, p.sinVotacion)
}
