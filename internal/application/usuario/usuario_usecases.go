package usuario

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	domusr "github.com/ipnext/admin-backend/internal/domain/usuario"
	"golang.org/x/crypto/bcrypt"
)

var ErrUsuarioNoEncontrado = errors.New("usuario no encontrado")

type ListUseCase struct{ repo domusr.Repository }

func NewListUseCase(repo domusr.Repository) *ListUseCase { return &ListUseCase{repo: repo} }
func (uc *ListUseCase) Execute(ctx context.Context) ([]*domusr.Usuario, error) {
	return uc.repo.FindAll(ctx)
}

type CreateRequest struct {
	Nombre   string
	Email    string
	Password string
	Permisos domusr.Permisos
}

type CreateUseCase struct{ repo domusr.Repository }

func NewCreateUseCase(repo domusr.Repository) *CreateUseCase { return &CreateUseCase{repo: repo} }
func (uc *CreateUseCase) Execute(ctx context.Context, req CreateRequest) (*domusr.Usuario, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &domusr.Usuario{
		ID: uuid.NewString(), Nombre: req.Nombre, Email: req.Email,
		Password: string(hash), Rol: domusr.RolSubUsuario,
		Permisos: req.Permisos, Activo: true,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := uc.repo.Save(ctx, u); err != nil {
		return nil, err
	}
	u.Password = ""
	return u, nil
}

type UpdateRequest struct {
	Nombre   *string
	Permisos *domusr.Permisos
	Activo   *bool
}

type UpdateUseCase struct{ repo domusr.Repository }

func NewUpdateUseCase(repo domusr.Repository) *UpdateUseCase { return &UpdateUseCase{repo: repo} }
func (uc *UpdateUseCase) Execute(ctx context.Context, id string, req UpdateRequest) (*domusr.Usuario, error) {
	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrUsuarioNoEncontrado
	}
	if req.Nombre != nil {
		u.Nombre = *req.Nombre
	}
	if req.Permisos != nil {
		u.Permisos = *req.Permisos
	}
	if req.Activo != nil {
		u.Activo = *req.Activo
	}
	u.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	u.Password = ""
	return u, nil
}

type DeleteUseCase struct{ repo domusr.Repository }

func NewDeleteUseCase(repo domusr.Repository) *DeleteUseCase { return &DeleteUseCase{repo: repo} }
func (uc *DeleteUseCase) Execute(ctx context.Context, id string) error {
	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if u == nil {
		return ErrUsuarioNoEncontrado
	}
	return uc.repo.Delete(ctx, id)
}
