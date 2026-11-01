package demo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// estadoEnProgreso aún no tiene constante en el dominio: la transición de
// estados es de HU-11.
const estadoEnProgreso domain.Estado = "En progreso"

type historiaDemo struct {
	clave       string
	titulo      string
	descripcion string
	criterios   string
	prioridad   domain.Prioridad
	valor       int
	puntos      int
	sprint      int // 1 a 3; 0 es el Product Backlog
	estado      domain.Estado
	finalizada  int // día relativo a hoy; solo para las Completadas
	arrastrada  bool
}

var historiasDemo = []historiaDemo{
	{clave: "turno", titulo: "Solicitar turno para el examen teórico",
		descripcion: "Como **solicitante** quiero elegir día y hora para rendir el examen teórico para:\n\n- evitar filas\n- llegar con el trámite iniciado",
		criterios:   "- Se muestran solo los turnos con cupo\n- Al confirmar se envía el comprobante",
		prioridad:   domain.PrioridadMust, valor: 13, puntos: 5, sprint: 1, estado: domain.EstadoCompletada, finalizada: -10},
	{clave: "datos", titulo: "Registrar datos personales del solicitante",
		descripcion: "Como solicitante quiero cargar mi DNI, domicilio y fecha de nacimiento para iniciar el trámite.",
		prioridad:   domain.PrioridadMust, valor: 8, puntos: 3, sprint: 1, estado: domain.EstadoCompletada, finalizada: -11},
	{clave: "apto", titulo: "Validar el certificado de aptitud psicofísica",
		descripcion: "Como agente de licencias quiero validar el certificado médico para habilitar el examen.",
		criterios:   "- El certificado no puede tener más de 30 días\n- Se rechaza si falta la firma del médico",
		prioridad:   domain.PrioridadShould, valor: 8, puntos: 8, sprint: 1, estado: domain.EstadoCompletada, finalizada: -9},
	{clave: "tasa", titulo: "Pagar la tasa de la licencia en línea",
		descripcion: "Como solicitante quiero pagar la tasa con tarjeta para no ir a una oficina de cobro.",
		prioridad:   domain.PrioridadMust, valor: 13, puntos: 8, estado: domain.EstadoNueva, sprint: 0, arrastrada: true},
	{clave: "examen", titulo: "Rendir el examen teórico en línea",
		descripcion: "Como **solicitante** quiero responder el examen desde el navegador para:\n\n1. rendir sin viajar\n2. ver el resultado al terminar",
		criterios:   "- El examen tiene 30 preguntas\n- Se aprueba con 24 respuestas correctas\n- El tiempo máximo es de 45 minutos",
		prioridad:   domain.PrioridadMust, valor: 21, puntos: 13, sprint: 2, estado: estadoEnProgreso},
	{clave: "comprobante", titulo: "Emitir el comprobante del turno",
		descripcion: "Como solicitante quiero descargar el comprobante del turno para presentarlo el día del examen.",
		criterios:   "- El comprobante incluye un código QR",
		prioridad:   domain.PrioridadShould, valor: 5, puntos: 3, sprint: 2, estado: estadoEnProgreso},
	{clave: "estado", titulo: "Consultar el estado del trámite",
		descripcion: "Como solicitante quiero ver en qué etapa está mi trámite para saber qué falta.",
		prioridad:   domain.PrioridadCould, valor: 3, puntos: 2, sprint: 2, estado: domain.EstadoNueva},
	{clave: "renovar", titulo: "Renovar la licencia vencida",
		descripcion: "Como **conductor** quiero renovar mi licencia para:\n\n- seguir conduciendo legalmente\n- conservar mi antigüedad",
		prioridad:   domain.PrioridadShould, valor: 8, puntos: 5, sprint: 3, estado: domain.EstadoNueva},
	{clave: "aviso", titulo: "Avisar por correo el vencimiento",
		descripcion: "Como conductor quiero recibir un correo antes de que venza mi licencia para renovarla a tiempo.",
		prioridad:   domain.PrioridadCould, valor: 5, estado: domain.EstadoNueva},
	{clave: "digital", titulo: "Mostrar la licencia en formato digital",
		descripcion: "Como conductor quiero llevar la licencia en el celular para no depender del plástico.",
		prioridad:   domain.PrioridadWont, estado: domain.EstadoNueva},
}

// bloqueHistorias crea las historias con su estado final ya resuelto y
// registra la historia arrastrada del Sprint 1.
func bloqueHistorias(ctx context.Context, tx *sql.Tx, c *carga) error {
	for _, h := range historiasDemo {
		id, err := insertarHistoria(ctx, tx, c, h)
		if err != nil {
			return err
		}
		c.historias[h.clave] = id
		c.resumen.Historias++

		if h.arrastrada {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO historias_arrastradas (sprint_id, historia_id, proyecto_id, estimacion_sp) VALUES (?, ?, ?, ?)`,
				c.sprints[0], id, proyectoID, h.puntos); err != nil {
				return fmt.Errorf("registrar historia arrastrada %s: %w", h.clave, err)
			}
		}
	}
	return nil
}

func insertarHistoria(ctx context.Context, tx *sql.Tx, c *carga, h historiaDemo) (int64, error) {
	valor := opcional(h.valor)
	item, err := domain.NewBacklogItem(proyectoID, h.titulo, h.descripcion, h.prioridad, valor)
	if err != nil {
		return 0, err
	}
	var criterios *string
	if h.criterios != "" {
		criterios = &h.criterios
	}
	item, err = item.ConDatosEditados(item.Titulo, item.Descripcion, criterios, opcional(h.puntos))
	if err != nil {
		return 0, err
	}

	var sprintID, finalizada any
	if h.sprint > 0 {
		sprintID = c.sprints[h.sprint-1]
	}
	if h.estado == domain.EstadoCompletada {
		finalizada = c.dia(h.finalizada).Format(time.DateOnly)
	}
	resultado, err := tx.ExecContext(ctx,
		`INSERT INTO historias (proyecto_id, titulo, descripcion, criterios_aceptacion, prioridad, estado, valor_negocio, estimacion_sp, sprint_id, fecha_finalizacion)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		proyectoID, item.Titulo, item.Descripcion, item.CriteriosAceptacion, string(item.Prioridad), string(h.estado),
		item.ValorNegocio, item.EstimacionSP, sprintID, finalizada)
	if err != nil {
		return 0, fmt.Errorf("crear historia %s: %w", h.clave, err)
	}
	return resultado.LastInsertId()
}

func opcional(n int) *int {
	if n == 0 {
		return nil
	}
	return &n
}
