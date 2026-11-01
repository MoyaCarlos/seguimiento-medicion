package demo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

const historiaVotada = "aviso"

var votosDemo = []struct {
	usuario string
	valor   int
}{
	{"Beta Uno", 3},
	{"Beta Dos", 5},
	{"Beta Tres", 5},
}

// bloqueVotacion crea una ronda revelada en una historia sin estimar del
// Product Backlog. No escribe la estimación en la historia.
func bloqueVotacion(ctx context.Context, tx *sql.Tx, c *carga) error {
	historiaID := c.historias[historiaVotada]
	ronda, err := domain.NuevaRonda(historiaID)
	if err != nil {
		return err
	}
	integrantes := make([]int64, 0, len(votosDemo))
	for _, voto := range votosDemo {
		id := c.usuarios[voto.usuario]
		if err := ronda.Votar(id, voto.valor); err != nil {
			return err
		}
		integrantes = append(integrantes, id)
	}
	if err := ronda.Revelar(integrantes, c.ahora); err != nil {
		return err
	}

	resultado, err := tx.ExecContext(ctx,
		`INSERT INTO rondas_votacion (historia_id, numero, revelada_en) VALUES (?, ?, ?)`,
		ronda.HistoriaID, ronda.Numero, ronda.ReveladaEn.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("crear ronda de votación: %w", err)
	}
	rondaID, err := resultado.LastInsertId()
	if err != nil {
		return err
	}
	c.resumen.Rondas++

	for _, voto := range ronda.Votos {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO votos (ronda_id, usuario_id, valor) VALUES (?, ?, ?)`,
			rondaID, voto.UsuarioID, voto.Valor); err != nil {
			return fmt.Errorf("registrar voto: %w", err)
		}
		c.resumen.Votos++
	}
	return nil
}
