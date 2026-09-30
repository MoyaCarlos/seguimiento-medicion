package repository

import "github.com/MoyaCarlos/seguimiento-medicion/internal/domain"

// ProjectRepository es lo mínimo que HU-05 necesita para poder crear un
// Sprint referenciando un Project existente. HU-04 puede agregar métodos
// (Update, AddMember, List, etc.) sin romper esto — solo no renombrar ni
// sacar los que ya están acá sin avisar a quien esté en HU-05.
type ProjectRepository interface {
	Create(p *domain.Project) error
	GetByID(id string) (*domain.Project, error)
}
