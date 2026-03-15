package alquiler

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Inmueble, error)
	FindAll(ctx context.Context) ([]*Inmueble, error)
	Save(ctx context.Context, i *Inmueble) error
	Update(ctx context.Context, i *Inmueble) error
	Delete(ctx context.Context, id string) error
}

type ContratoRepository interface {
	FindByID(ctx context.Context, id string) (*ContratoAlquiler, error)
	FindAll(ctx context.Context, inmuebleID *string) ([]*ContratoAlquiler, error)
	FindProximosAVencer(ctx context.Context, dias int) ([]*ContratoAlquiler, error)
	Save(ctx context.Context, c *ContratoAlquiler) error
	Update(ctx context.Context, c *ContratoAlquiler) error
}

type PagoRepository interface {
	FindAll(ctx context.Context, inmuebleID *string) ([]*PagoAlquiler, error)
	Save(ctx context.Context, p *PagoAlquiler) error
}
