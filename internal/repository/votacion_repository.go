package repository

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// VotacionRepository es el puerto de persistencia de Planning Poker. La ronda
// actual de una historia es la de mayor número (MODELO-DE-DATOS.md §4.7).
//   - CrearRonda devuelve domain.ErrVotacionYaIniciada si la historia ya tiene
//     una ronda con ese número.
//   - ObtenerRondaActual devuelve domain.ErrVotacionNoIniciada si la historia
//     no tiene rondas.
//   - GuardarVoto y MarcarRevelada escriben solo si la ronda sigue sin revelar;
//     si no, devuelven domain.ErrRondaRevelada sin cambiar nada.
type VotacionRepository interface {
	CrearRonda(ctx context.Context, r domain.RondaVotacion) (domain.RondaVotacion, error)
	ObtenerRondaActual(ctx context.Context, historiaID int64) (domain.RondaVotacion, error)
	GuardarVoto(ctx context.Context, rondaID int64, v domain.Voto) error
	MarcarRevelada(ctx context.Context, rondaID int64, en time.Time) error
}
