package proveedor

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Proveedor, error)
	FindAll(ctx context.Context, query string, soloActivos bool) ([]*Proveedor, error)
	Save(ctx context.Context, p *Proveedor) error
	Update(ctx context.Context, p *Proveedor) error
	Delete(ctx context.Context, id string) error
}

type ContratoRepository interface {
	FindByID(ctx context.Context, id string) (*ContratoProveedor, error)
	FindAll(ctx context.Context, proveedorID *string) ([]*ContratoProveedor, error)
	Save(ctx context.Context, c *ContratoProveedor) error
	Update(ctx context.Context, c *ContratoProveedor) error
}

type RankingRepository interface {
	GetRanking(ctx context.Context, limit int) ([]*RankingItem, error)
}
