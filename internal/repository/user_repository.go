package repository

import "github.com/MoyaCarlos/seguimiento-medicion/internal/domain"

// UserRepository abstrae la persistencia de los integrantes del equipo.
type UserRepository interface {
	FindByNormalizedName(normalizedName string) (*domain.User, error)
	Create(u *domain.User) error
}
