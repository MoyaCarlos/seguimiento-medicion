package http

import "github.com/MoyaCarlos/seguimiento-medicion/internal/domain"

type avanceResponse struct{}

func aAvanceResponse(avance domain.AvanceHistorias) avanceResponse {
	return avanceResponse{}
}
