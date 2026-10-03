package service

import (
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

type fakeProjectRepository struct {
	proyectos     map[int64]domain.Project
	miembros      []domain.Membership
	listaMiembros []domain.Member
	err           error
	creates       int
	updates       int
}

func (f *fakeProjectRepository) Create(p *domain.Project) error {
	if f.err != nil {
		return f.err
	}
	f.creates++
	if f.proyectos == nil {
		f.proyectos = map[int64]domain.Project{}
	}
	if p.ID == 0 {
		p.ID = int64(f.creates)
	}
	p.CreatedAt = time.Now().UTC()
	f.proyectos[p.ID] = *p
	return nil
}

func (f *fakeProjectRepository) CreateWithScrumMaster(p *domain.Project, creatorID int64) error {
	if err := f.Create(p); err != nil {
		return err
	}
	return f.AddMember(&domain.Membership{ProjectID: p.ID, UserID: creatorID, Role: domain.RolScrumMaster})
}

func (f *fakeProjectRepository) GetByID(id int64) (*domain.Project, error) {
	if f.err != nil {
		return nil, f.err
	}
	p, ok := f.proyectos[id]
	if !ok {
		return nil, domain.ErrProyectoNoEncontrado
	}
	return &p, nil
}

func (f *fakeProjectRepository) Update(p *domain.Project) error {
	if f.err != nil {
		return f.err
	}
	if _, ok := f.proyectos[p.ID]; !ok {
		return domain.ErrProyectoNoEncontrado
	}
	f.updates++
	f.proyectos[p.ID] = *p
	return nil
}

func (f *fakeProjectRepository) AddMember(m *domain.Membership) error {
	if f.err != nil {
		return f.err
	}
	for _, existente := range f.miembros {
		if existente.ProjectID == m.ProjectID && existente.UserID == m.UserID {
			return repository.ErrMiembroDuplicado
		}
	}
	f.miembros = append(f.miembros, *m)
	return nil
}

func (f *fakeProjectRepository) ListMembers(int64) ([]domain.Member, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.listaMiembros, nil
}

type fakeUserRepository struct {
	usuarios map[string]domain.User
	creates  int
	err      error
}

func (f *fakeUserRepository) FindByNormalizedName(normalized string) (*domain.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	u, ok := f.usuarios[normalized]
	if !ok {
		return nil, domain.ErrNoEncontrado
	}
	return &u, nil
}

func (f *fakeUserRepository) Create(u *domain.User) error {
	if f.err != nil {
		return f.err
	}
	if f.usuarios == nil {
		f.usuarios = map[string]domain.User{}
	}
	if _, ok := f.usuarios[u.NormalizedName]; ok {
		return repository.ErrUsuarioDuplicado
	}
	f.creates++
	u.ID = int64(f.creates)
	u.CreatedAt = time.Now().UTC()
	f.usuarios[u.NormalizedName] = *u
	return nil
}
