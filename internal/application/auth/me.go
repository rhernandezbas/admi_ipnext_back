package auth

import (
	"context"
	"errors"

	"github.com/ipnext/admin-backend/internal/domain/usuario"
)

var ErrUsuarioNoEncontrado = errors.New("usuario no encontrado")

type GetMeUseCase struct {
	repo usuario.Repository
}

func NewGetMeUseCase(repo usuario.Repository) *GetMeUseCase {
	return &GetMeUseCase{repo: repo}
}

func (uc *GetMeUseCase) Execute(ctx context.Context, id string) (*usuario.Usuario, error) {
	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrUsuarioNoEncontrado
	}
	return u, nil
}
