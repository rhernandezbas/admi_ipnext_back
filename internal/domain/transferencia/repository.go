package transferencia

import (
	"context"
	"time"
)

type Filtros struct {
	Estado    string
	Categoria string
	Desde     *time.Time
	Hasta     *time.Time
	Query     string
	Orden     string
	Dir       string
	Page      int
	PerPage   int
}

type ListResult struct {
	Items      []*Transferencia
	Total      int64
	Page       int
	PerPage    int
	TotalPages int
}

type Repository interface {
	FindByID(ctx context.Context, id string) (*Transferencia, error)
	FindAll(ctx context.Context, filtros Filtros) (*ListResult, error)
	FindRecurrentes(ctx context.Context) ([]*Transferencia, error)
	FindByFecha(ctx context.Context, desde, hasta time.Time) ([]*Transferencia, error)
	FindProximasAVencer(ctx context.Context, dias int) ([]*Transferencia, error)
	Save(ctx context.Context, t *Transferencia) error
	Update(ctx context.Context, t *Transferencia) error
	Delete(ctx context.Context, id string) error
	SumByCategoria(ctx context.Context, desde, hasta time.Time) (map[string]float64, error)
}
