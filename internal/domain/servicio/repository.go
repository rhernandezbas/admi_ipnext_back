package servicio

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Servicio, error)
	FindAll(ctx context.Context, tipo *Tipo) ([]*Servicio, error)
	Save(ctx context.Context, s *Servicio) error
	Update(ctx context.Context, s *Servicio) error
	Delete(ctx context.Context, id string) error
	GetKPIs(ctx context.Context) (*KPIs, error)
}
