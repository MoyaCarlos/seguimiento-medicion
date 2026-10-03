package repository

import (
	"sync"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// ProjectRepositoryEnMemoria implementa ProjectRepository en memoria para que
// HU-05 y HU-13 usen un fake compartido en lugar de fakes propios.
type ProjectRepositoryEnMemoria struct {
	mu        sync.Mutex
	proyectos map[int64]domain.Project
	miembros  []domain.Membership
	siguiente int64
}

func NewProjectRepositoryEnMemoria() *ProjectRepositoryEnMemoria {
	return &ProjectRepositoryEnMemoria{proyectos: map[int64]domain.Project{}}
}

func (r *ProjectRepositoryEnMemoria) Create(p *domain.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	r.siguiente++
	p.ID = r.siguiente
	r.proyectos[p.ID] = *p
	return nil
}

func (r *ProjectRepositoryEnMemoria) CreateWithScrumMaster(p *domain.Project, creatorID int64) error {
	if err := r.Create(p); err != nil {
		return err
	}
	return r.AddMember(&domain.Membership{ProjectID: p.ID, UserID: creatorID, Role: domain.RolScrumMaster})
}

func (r *ProjectRepositoryEnMemoria) GetByID(id int64) (*domain.Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.proyectos[id]
	if !ok {
		return nil, domain.ErrProyectoNoEncontrado
	}
	return &p, nil
}

func (r *ProjectRepositoryEnMemoria) Update(p *domain.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.proyectos[p.ID]; !ok {
		return domain.ErrProyectoNoEncontrado
	}
	r.proyectos[p.ID] = *p
	return nil
}

func (r *ProjectRepositoryEnMemoria) AddMember(m *domain.Membership) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existente := range r.miembros {
		if existente.ProjectID == m.ProjectID && existente.UserID == m.UserID {
			return ErrMiembroDuplicado
		}
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	r.miembros = append(r.miembros, *m)
	return nil
}

func (r *ProjectRepositoryEnMemoria) ListMembers(projectID int64) ([]domain.Member, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []domain.Member
	for _, m := range r.miembros {
		if m.ProjectID == projectID {
			out = append(out, domain.Member{
				UserID:    m.UserID,
				ProjectID: m.ProjectID,
				Role:      m.Role,
				CreatedAt: m.CreatedAt,
			})
		}
	}
	return out, nil
}
