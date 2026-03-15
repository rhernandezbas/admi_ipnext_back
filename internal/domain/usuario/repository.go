package usuario

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Usuario, error)
	FindByEmail(ctx context.Context, email string) (*Usuario, error)
	FindAll(ctx context.Context) ([]*Usuario, error)
	Save(ctx context.Context, u *Usuario) error
	Update(ctx context.Context, u *Usuario) error
	Delete(ctx context.Context, id string) error
}
