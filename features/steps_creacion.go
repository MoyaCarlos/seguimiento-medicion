package features

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type scenarioContext struct {
	db         *sql.DB
	servicio   *service.CrearHistoriaBacklog
	proyectoID int64
	creada     domain.BacklogItem
	err        error
}

func (s *scenarioContext) reset() {
	if s.db != nil {
		_ = s.db.Close()
		s.db = nil
	}
	s.proyectoID = 1
	s.creada = domain.BacklogItem{}
	s.err = nil
}

func (s *scenarioContext) iniciarBD() error {
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(1)
	if err := repository.Migrar(context.Background(), db); err != nil {
		return err
	}
	s.db = db
	s.servicio = service.NewCrearHistoriaBacklog(repository.NewSQLiteBacklogRepository(db))
	return nil
}

func (s *scenarioContext) existeProyecto(identificador int) error {
	p, err := domain.NewProject("Proyecto de prueba", "", nil, nil)
	if err != nil {
		return err
	}
	guardado, err := repository.NewSQLiteProjectRepository(s.db).Guardar(context.Background(), p)
	if err != nil {
		return err
	}
	if guardado.ID != int64(identificador) {
		return fmt.Errorf("el escenario asume una base vacía: el proyecto quedó con id %d, no %d", guardado.ID, identificador)
	}
	s.proyectoID = guardado.ID
	return nil
}

func (s *scenarioContext) enElPanel() error { return nil }

func (s *scenarioContext) intentoCrea() error { return nil }

func (s *scenarioContext) ingresoValido() error {
	s.creada, s.err = s.servicio.Ejecutar(context.Background(), service.CrearHistoriaInput{
		ProyectoID:  s.proyectoID,
		Titulo:      "Como usuario quiero iniciar sesión",
		Descripcion: "Autenticación con email y contraseña",
		Prioridad:   domain.PrioridadMust,
	})
	return nil
}

func (s *scenarioContext) tituloEnBlanco() error {
	s.creada, s.err = s.servicio.Ejecutar(context.Background(), service.CrearHistoriaInput{
		ProyectoID:  s.proyectoID,
		Titulo:      "",
		Descripcion: "Descripción válida",
		Prioridad:   domain.PrioridadMust,
	})
	return nil
}

func (s *scenarioContext) prioridadInvalida() error {
	s.creada, s.err = s.servicio.Ejecutar(context.Background(), service.CrearHistoriaInput{
		ProyectoID:  s.proyectoID,
		Titulo:      "Historia con prioridad inválida",
		Descripcion: "Descripción válida",
		Prioridad:   domain.Prioridad("Urgente"),
	})
	return nil
}

func (s *scenarioContext) seGuardaEnSQLite() error {
	if s.err != nil {
		return fmt.Errorf("se esperaba creación exitosa, se obtuvo: %w", s.err)
	}
	if s.creada.ID <= 0 {
		return errors.New("se esperaba un identificador asignado")
	}
	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM backlog_items WHERE id = ?", s.creada.ID).Scan(&total); err != nil {
		return err
	}
	if total != 1 {
		return errors.New("la historia no quedó persistida en SQLite")
	}
	return nil
}

func (s *scenarioContext) alFinalDelBacklog() error {
	var maxID int64
	if err := s.db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM backlog_items").Scan(&maxID); err != nil {
		return err
	}
	if s.creada.ID != maxID {
		return fmt.Errorf("se esperaba la historia al final (id %d), máximo actual %d", s.creada.ID, maxID)
	}
	return nil
}

func (s *scenarioContext) estadoNueva() error {
	if s.creada.Estado != domain.EstadoNueva {
		return fmt.Errorf("se esperaba estado %q, se obtuvo %q", domain.EstadoNueva, s.creada.Estado)
	}
	return nil
}

func (s *scenarioContext) sinEstimacion() error {
	if s.creada.EstimacionSP != nil {
		return errors.New("se esperaba la estimación sin asignar")
	}
	return nil
}

func (s *scenarioContext) muestraAdvertenciaDeValidacion() error {
	return s.verificarErrorDeValidacion()
}

func (s *scenarioContext) muestraErrorDeValidacion() error {
	return s.verificarErrorDeValidacion()
}

func (s *scenarioContext) verificarErrorDeValidacion() error {
	if s.err == nil {
		return errors.New("se esperaba un error de validación")
	}
	var verr domain.ValidationError
	if !errors.As(s.err, &verr) {
		return fmt.Errorf("se esperaba ValidationError, se obtuvo %T", s.err)
	}
	return nil
}

func (s *scenarioContext) noRegistraHistoria() error {
	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM backlog_items").Scan(&total); err != nil {
		return err
	}
	if total != 0 {
		return fmt.Errorf("se esperaba 0 historias registradas, se encontraron %d", total)
	}
	return nil
}

// InitializeScenario registra los pasos de los escenarios de HU-01.
func InitializeScenario(ctx *godog.ScenarioContext) {
	sc := &scenarioContext{}

	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		sc.reset()
		return goctx, sc.iniciarBD()
	})
	ctx.After(func(goctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if sc.db != nil {
			_ = sc.db.Close()
			sc.db = nil
		}
		return goctx, nil
	})

	ctx.Step(`^que existe un proyecto con identificador (\d+)$`, sc.existeProyecto)
	ctx.Step(`^que estoy en el panel del Product Backlog del proyecto$`, sc.enElPanel)
	ctx.Step(`^que intento crear una historia$`, sc.intentoCrea)
	ctx.Step(`^ingreso un título, una descripción y una prioridad "M" válidos y presiono "Guardar"$`, sc.ingresoValido)
	ctx.Step(`^dejo el campo "Título" en blanco y presiono "Guardar"$`, sc.tituloEnBlanco)
	ctx.Step(`^ingreso la prioridad "Urgente", que no pertenece al enum MoSCoW \(Must have / Should have / Could have / Won't have\), y presiono "Guardar"$`, sc.prioridadInvalida)
	ctx.Step(`^la historia se guarda de forma persistente en SQLite$`, sc.seGuardaEnSQLite)
	ctx.Step(`^queda al final del Product Backlog según el orden de creación \(identificador incremental posterior\)$`, sc.alFinalDelBacklog)
	ctx.Step(`^su estado es "Nueva"$`, sc.estadoNueva)
	ctx.Step(`^no tiene estimación en Story Points asignada$`, sc.sinEstimacion)
	ctx.Step(`^el sistema muestra una advertencia de validación$`, sc.muestraAdvertenciaDeValidacion)
	ctx.Step(`^el sistema muestra un error de validación$`, sc.muestraErrorDeValidacion)
	ctx.Step(`^no registra la historia en el Product Backlog$`, sc.noRegistraHistoria)

	inicializarPasosSprint(ctx, sc)
}
