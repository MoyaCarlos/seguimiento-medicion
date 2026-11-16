package demo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type miembroDemo struct {
	usuario string
	rol     domain.Role
	baja    *int // día relativo a hoy de la baja; nil si está activo
}

type rondaAdicional struct {
	historia string
	numero   int
	votos    []votoDemo
	revelada bool
}

// proyectoDemo describe un proyecto adicional. Reutiliza los tipos de las filas
// del proyecto 25 (sprintDemo, historiaDemo, horasDemo, defectoDemo); los
// índices de Sprint son 1-based dentro del propio proyecto.
type proyectoDemo struct {
	nombre, descripcion string
	inicio, fin         int
	cierre              *int
	horasPorSP          int
	equipo              []miembroDemo
	sprints             []sprintDemo
	historias           []historiaDemo
	arrastreDesde       int
	horas               []horasDemo
	rondas              []rondaAdicional
	defectos            []defectoDemo
}

// proyectoEnCurso es el estado propio de un proyecto adicional mientras se
// inserta: los ids de sus Sprints e historias.
type proyectoEnCurso struct {
	id        int64
	def       proyectoDemo
	sprints   []int64
	historias map[string]int64
}

// bloqueProyectosAdicionales crea, sin ids fijos, los proyectos de
// proyectosAdicionales (spec 052). Reutiliza los usuarios de bloqueUsuarios.
func bloqueProyectosAdicionales(ctx context.Context, tx *sql.Tx, c *carga) error {
	for _, def := range proyectosAdicionales {
		p := &proyectoEnCurso{def: def, historias: map[string]int64{}}
		pasos := []func(context.Context, *sql.Tx, *carga, *proyectoEnCurso) error{
			insertarProyectoAdicional, insertarSprintsAdicionales, insertarHistoriasAdicionales,
			insertarHorasAdicionales, insertarRondasAdicionales, insertarDefectosAdicionales,
		}
		for _, paso := range pasos {
			if err := paso(ctx, tx, c, p); err != nil {
				return fmt.Errorf("%s: %w", def.nombre, err)
			}
		}
		c.resumen.ProyectosAdicionales++
	}
	return nil
}

func insertarProyectoAdicional(ctx context.Context, tx *sql.Tx, c *carga, p *proyectoEnCurso) error {
	inicio, fin := c.dia(p.def.inicio), c.dia(p.def.fin)
	proyecto, err := domain.NewProject(p.def.nombre, p.def.descripcion, &inicio, &fin)
	if err != nil {
		return err
	}
	var horasPorSP, cierre any
	if p.def.horasPorSP > 0 {
		horasPorSP = p.def.horasPorSP
	}
	if p.def.cierre != nil {
		cierre = c.dia(*p.def.cierre).Format(time.DateOnly)
	}
	resultado, err := tx.ExecContext(ctx,
		`INSERT INTO proyectos (nombre, descripcion, fecha_inicio, fecha_fin, horas_por_story_point, fecha_cierre, creado_en)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		proyecto.Nombre, proyecto.Descripcion, inicio.Format(time.DateOnly), fin.Format(time.DateOnly),
		horasPorSP, cierre, c.creadoEn())
	if err != nil {
		return fmt.Errorf("crear proyecto: %w", err)
	}
	if p.id, err = resultado.LastInsertId(); err != nil {
		return err
	}
	for _, m := range p.def.equipo {
		if err := insertarMiembroAdicional(ctx, tx, c, p.id, m); err != nil {
			return err
		}
	}
	return nil
}

func insertarMiembroAdicional(ctx context.Context, tx *sql.Tx, c *carga, proyecto int64, m miembroDemo) error {
	membresia, err := domain.NewMembership(proyecto, c.usuarios[m.usuario], m.rol)
	if err != nil {
		return err
	}
	var baja any
	if m.baja != nil {
		baja = c.dia(*m.baja).Format(time.DateOnly)
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO integrantes (proyecto_id, usuario_id, rol, fecha_baja, creado_en) VALUES (?, ?, ?, ?, ?)`,
		membresia.ProjectID, membresia.UserID, string(membresia.Role), baja, c.creadoEn()); err != nil {
		return fmt.Errorf("agregar integrante %s: %w", m.usuario, err)
	}
	return nil
}

func insertarSprintsAdicionales(ctx context.Context, tx *sql.Tx, c *carga, p *proyectoEnCurso) error {
	for _, definicion := range p.def.sprints {
		sprint, err := domain.NewSprint(p.id)
		if err != nil {
			return err
		}
		if definicion.estado != domain.SprintPendiente {
			if err := sprint.Iniciar(definicion.goal, c.dia(definicion.desde), c.dia(definicion.hasta)); err != nil {
				return err
			}
		}
		if definicion.estado == domain.SprintFinalizado {
			if err := sprint.Cerrar(); err != nil {
				return err
			}
		}
		id, err := insertarSprint(ctx, tx, sprint)
		if err != nil {
			return err
		}
		p.sprints = append(p.sprints, id)
		c.resumen.SprintsAdicionales++
	}
	return nil
}

func insertarHistoriasAdicionales(ctx context.Context, tx *sql.Tx, c *carga, p *proyectoEnCurso) error {
	for _, h := range p.def.historias {
		item, err := domain.NewBacklogItem(p.id, h.titulo, h.descripcion, h.prioridad, opcional(h.valor))
		if err != nil {
			return fmt.Errorf("historia %s: %w", h.clave, err)
		}
		var criterios *string
		if h.criterios != "" {
			criterios = &h.criterios
		}
		if item, err = item.ConDatosEditados(item.Titulo, item.Descripcion, criterios, opcional(h.puntos)); err != nil {
			return fmt.Errorf("historia %s: %w", h.clave, err)
		}
		var sprintID, finalizada any
		if h.sprint > 0 {
			sprintID = p.sprints[h.sprint-1]
		}
		if h.estado == domain.EstadoCompletada {
			finalizada = c.dia(h.finalizada).Format(time.DateOnly)
		}
		resultado, err := tx.ExecContext(ctx,
			`INSERT INTO historias (proyecto_id, titulo, descripcion, criterios_aceptacion, prioridad, estado, valor_negocio, estimacion_sp, sprint_id, fecha_finalizacion)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			p.id, item.Titulo, item.Descripcion, item.CriteriosAceptacion, string(item.Prioridad), string(h.estado),
			item.ValorNegocio, item.EstimacionSP, sprintID, finalizada)
		if err != nil {
			return fmt.Errorf("crear historia %s: %w", h.clave, err)
		}
		id, err := resultado.LastInsertId()
		if err != nil {
			return err
		}
		p.historias[h.clave] = id
		c.resumen.HistoriasAdicionales++
		if h.arrastrada {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO historias_arrastradas (sprint_id, historia_id, proyecto_id, estimacion_sp) VALUES (?, ?, ?, ?)`,
				p.sprints[p.def.arrastreDesde-1], id, p.id, h.puntos); err != nil {
				return fmt.Errorf("registrar historia arrastrada %s: %w", h.clave, err)
			}
		}
	}
	return nil
}

func insertarHorasAdicionales(ctx context.Context, tx *sql.Tx, c *carga, p *proyectoEnCurso) error {
	for _, h := range p.def.horas {
		definicion := p.def.sprints[h.sprint-1]
		sprint := domain.Sprint{ID: p.sprints[h.sprint-1], ProyectoID: p.id,
			FechaInicio: c.dia(definicion.desde), FechaFin: c.dia(definicion.hasta)}
		historia := domain.BacklogItem{ID: p.historias[h.historia], ProyectoID: p.id}
		registro, err := domain.NewRegistroHoras(historia, sprint, c.usuarios[h.usuario],
			h.horas, h.actividad, c.dia(h.dia).Format(time.DateOnly), c.ahora)
		if err != nil {
			return fmt.Errorf("horas de %s: %w", h.historia, err)
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO registros_horas (proyecto_id, historia_id, usuario_id, sprint_id, fecha, actividad, minutos, creado_en)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			registro.ProyectoID, registro.HistoriaID, registro.UsuarioID, registro.SprintID,
			registro.Fecha.Format(time.DateOnly), registro.Actividad, registro.Minutos, c.creadoEn()); err != nil {
			return fmt.Errorf("imputar horas de %s: %w", h.historia, err)
		}
		c.resumen.RegistrosHorasAdicionales++
	}
	return nil
}

func insertarRondasAdicionales(ctx context.Context, tx *sql.Tx, c *carga, p *proyectoEnCurso) error {
	for _, r := range p.def.rondas {
		ronda := domain.RondaVotacion{HistoriaID: p.historias[r.historia], Numero: r.numero}
		votantes := make([]int64, 0, len(r.votos))
		for _, voto := range r.votos {
			id := c.usuarios[voto.usuario]
			if err := ronda.Votar(id, voto.valor); err != nil {
				return err
			}
			votantes = append(votantes, id)
		}
		var revelada any
		if r.revelada {
			if err := ronda.Revelar(votantes, c.ahora); err != nil {
				return err
			}
			revelada = ronda.ReveladaEn.Format(time.RFC3339)
		}
		resultado, err := tx.ExecContext(ctx,
			`INSERT INTO rondas_votacion (historia_id, numero, revelada_en) VALUES (?, ?, ?)`,
			ronda.HistoriaID, ronda.Numero, revelada)
		if err != nil {
			return fmt.Errorf("crear ronda de %s: %w", r.historia, err)
		}
		rondaID, err := resultado.LastInsertId()
		if err != nil {
			return err
		}
		c.resumen.RondasAdicionales++
		for _, voto := range ronda.Votos {
			if _, err := tx.ExecContext(ctx, `INSERT INTO votos (ronda_id, usuario_id, valor) VALUES (?, ?, ?)`,
				rondaID, voto.UsuarioID, voto.Valor); err != nil {
				return fmt.Errorf("registrar voto: %w", err)
			}
			c.resumen.VotosAdicionales++
		}
	}
	return nil
}

func insertarDefectosAdicionales(ctx context.Context, tx *sql.Tx, c *carga, p *proyectoEnCurso) error {
	for _, d := range p.def.defectos {
		defecto, err := domain.NewDefecto(p.id, p.historias[d.historia], domain.Sprint{ID: p.sprints[d.deteccion-1]}, d.descripcion, d.severidad)
		if err != nil {
			return err
		}
		var resolucion any
		if d.resolucion > 0 {
			sprint := domain.Sprint{ID: p.sprints[d.resolucion-1]}
			if err := defecto.Resolver(&sprint); err != nil {
				return err
			}
			resolucion = sprint.ID
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO defectos (proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id, creado_en)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			p.id, defecto.HistoriaID, defecto.Descripcion, string(defecto.Severidad), string(defecto.Estado),
			defecto.SprintDeteccionID, resolucion, c.dia(d.dia).Add(10*time.Hour).Format(time.RFC3339)); err != nil {
			return fmt.Errorf("crear defecto de %s: %w", d.historia, err)
		}
		c.resumen.DefectosAdicionales++
	}
	return nil
}
