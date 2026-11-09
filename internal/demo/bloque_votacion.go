package demo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type votoDemo struct {
	usuario string
	valor   int
}

type rondaDemo struct {
	historia string
	votos    []votoDemo
	revelada bool
}

// rondasDemo: una ronda revelada con todos los votos y una abierta con votos
// parciales, ambas en historias sin estimar del Product Backlog.
var rondasDemo = []rondaDemo{
	{"aviso", []votoDemo{{"carlos.moya", 3}, {"ana.iriarte", 5}, {"candela.vulcano", 5}, {"jimena.martinez", 8}}, true},
	{"digital", []votoDemo{{"ana.iriarte", 3}, {"jimena.martinez", 5}}, false},
}

// bloqueVotacion crea las rondas de Planning Poker. No escribe la estimación
// en las historias.
func bloqueVotacion(ctx context.Context, tx *sql.Tx, c *carga) error {
	for _, definicion := range rondasDemo {
		if err := crearRonda(ctx, tx, c, definicion); err != nil {
			return err
		}
	}
	return nil
}

func crearRonda(ctx context.Context, tx *sql.Tx, c *carga, definicion rondaDemo) error {
	ronda, err := domain.NuevaRonda(c.historias[definicion.historia])
	if err != nil {
		return err
	}
	integrantes := make([]int64, 0, len(definicion.votos))
	for _, voto := range definicion.votos {
		id := c.usuarios[voto.usuario]
		if err := ronda.Votar(id, voto.valor); err != nil {
			return err
		}
		integrantes = append(integrantes, id)
	}
	var revelada any
	if definicion.revelada {
		if err := ronda.Revelar(integrantes, c.ahora); err != nil {
			return err
		}
		revelada = ronda.ReveladaEn.Format(time.RFC3339)
	}

	resultado, err := tx.ExecContext(ctx,
		`INSERT INTO rondas_votacion (historia_id, numero, revelada_en) VALUES (?, ?, ?)`,
		ronda.HistoriaID, ronda.Numero, revelada)
	if err != nil {
		return fmt.Errorf("crear ronda de votación de %s: %w", definicion.historia, err)
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
