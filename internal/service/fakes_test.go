package service

import (
	"context"
	"sort"
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

func (f *fakeProjectRepository) Guardar(_ context.Context, p domain.Project) (domain.Project, error) {
	if f.err != nil {
		return domain.Project{}, f.err
	}
	f.creates++
	if f.proyectos == nil {
		f.proyectos = map[int64]domain.Project{}
	}
	if p.ID == 0 {
		p.ID = int64(f.creates)
	}
	p.CreadoEn = time.Now().UTC()
	f.proyectos[p.ID] = p
	return p, nil
}

func (f *fakeProjectRepository) GuardarConScrumMaster(ctx context.Context, p domain.Project, creadorID int64) (domain.Project, error) {
	guardado, err := f.Guardar(ctx, p)
	if err != nil {
		return domain.Project{}, err
	}
	if _, err := f.AgregarIntegrante(ctx, domain.Membership{ProjectID: guardado.ID, UserID: creadorID, Role: domain.RolScrumMaster}); err != nil {
		return domain.Project{}, err
	}
	return guardado, nil
}

func (f *fakeProjectRepository) ObtenerPorID(_ context.Context, id int64) (domain.Project, error) {
	if f.err != nil {
		return domain.Project{}, f.err
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Project{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *fakeProjectRepository) Listar(_ context.Context) ([]domain.Project, error) {
	if f.err != nil {
		return nil, f.err
	}
	proyectos := make([]domain.Project, 0, len(f.proyectos))
	for _, p := range f.proyectos {
		proyectos = append(proyectos, p)
	}
	sort.Slice(proyectos, func(i, j int) bool { return proyectos[i].ID < proyectos[j].ID })
	return proyectos, nil
}

func (f *fakeProjectRepository) Actualizar(_ context.Context, p domain.Project) error {
	if f.err != nil {
		return f.err
	}
	if _, ok := f.proyectos[p.ID]; !ok {
		return domain.ErrProyectoNoEncontrado
	}
	f.updates++
	f.proyectos[p.ID] = p
	return nil
}

func (f *fakeProjectRepository) AgregarIntegrante(_ context.Context, m domain.Membership) (domain.Membership, error) {
	if f.err != nil {
		return domain.Membership{}, f.err
	}
	for _, existente := range f.miembros {
		if existente.ProjectID == m.ProjectID && existente.UserID == m.UserID {
			return domain.Membership{}, repository.ErrMiembroDuplicado
		}
	}
	if m.CreadoEn.IsZero() {
		m.CreadoEn = time.Now().UTC()
	}
	f.miembros = append(f.miembros, m)
	return m, nil
}

func (f *fakeProjectRepository) ListarIntegrantes(_ context.Context, _ int64) ([]domain.Member, error) {
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

func (f *fakeUserRepository) ObtenerPorNombreNormalizado(_ context.Context, normalized string) (*domain.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	u, ok := f.usuarios[normalized]
	if !ok {
		return nil, domain.ErrNoEncontrado
	}
	return &u, nil
}

func (f *fakeUserRepository) Guardar(_ context.Context, u domain.User) (domain.User, error) {
	if f.err != nil {
		return domain.User{}, f.err
	}
	if f.usuarios == nil {
		f.usuarios = map[string]domain.User{}
	}
	if _, ok := f.usuarios[u.NombreNormalizado]; ok {
		return domain.User{}, repository.ErrUsuarioDuplicado
	}
	f.creates++
	u.ID = int64(f.creates)
	u.CreadoEn = time.Now().UTC()
	f.usuarios[u.NombreNormalizado] = u
	return u, nil
}

type fakeProductBacklogRepository struct {
	historias []domain.BacklogItem
	err       error
	llamadas  int
}

func (f *fakeProductBacklogRepository) ListarSinSprint(_ context.Context, _ int64) ([]domain.BacklogItem, error) {
	f.llamadas++
	if f.err != nil {
		return nil, f.err
	}
	return f.historias, nil
}

type fakeHistoriaRepository struct {
	historias       map[int64]domain.BacklogItem
	errActualizar   error
	actualizaciones []domain.BacklogItem
}

func (f *fakeHistoriaRepository) ObtenerPorID(_ context.Context, id int64) (domain.BacklogItem, error) {
	h, ok := f.historias[id]
	if !ok {
		return domain.BacklogItem{}, domain.ErrHistoriaNoEncontrada
	}
	return h, nil
}

func (f *fakeHistoriaRepository) Actualizar(_ context.Context, item domain.BacklogItem) error {
	if f.errActualizar != nil {
		return f.errActualizar
	}
	f.actualizaciones = append(f.actualizaciones, item)
	f.historias[item.ID] = item
	return nil
}

type fakeRegistroHorasRepository struct {
	registros    []domain.RegistroHoras
	listados     int
	errRegistrar error
}

func (f *fakeRegistroHorasRepository) Registrar(_ context.Context, r domain.RegistroHoras) (domain.RegistroHoras, error) {
	if f.errRegistrar != nil {
		return domain.RegistroHoras{}, f.errRegistrar
	}
	r.ID = int64(len(f.registros) + 1)
	f.registros = append(f.registros, r)
	return r, nil
}

func (f *fakeRegistroHorasRepository) ListarPorHistoria(_ context.Context, historiaID int64) ([]domain.RegistroHorasDetalle, error) {
	f.listados++
	detalle := make([]domain.RegistroHorasDetalle, 0)
	for _, r := range f.registros {
		if r.HistoriaID == historiaID {
			detalle = append(detalle, domain.RegistroHorasDetalle{RegistroHoras: r})
		}
	}
	return detalle, nil
}

type fakeIntegrantesRepository struct {
	miembros []domain.Member
}

func (f *fakeIntegrantesRepository) ListarIntegrantes(_ context.Context, proyectoID int64) ([]domain.Member, error) {
	var resultado []domain.Member
	for _, m := range f.miembros {
		if m.ProjectID == proyectoID {
			resultado = append(resultado, m)
		}
	}
	return resultado, nil
}

// fakeSprintsQueCuentan envuelve el repositorio en memoria y cuenta las lecturas.
type fakeSprintsQueCuentan struct {
	*repository.SprintRepositoryEnMemoria
	lecturas int
}

func (f *fakeSprintsQueCuentan) ObtenerPorID(ctx context.Context, id int64) (domain.Sprint, error) {
	f.lecturas++
	return f.SprintRepositoryEnMemoria.ObtenerPorID(ctx, id)
}
