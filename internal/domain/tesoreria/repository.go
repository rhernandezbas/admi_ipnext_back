package tesoreria

import (
	"context"
	"time"
)

type CuentaRepository interface {
	FindByID(ctx context.Context, id string) (*CuentaBancaria, error)
	FindAll(ctx context.Context) ([]*CuentaBancaria, error)
	Save(ctx context.Context, c *CuentaBancaria) error
	Update(ctx context.Context, c *CuentaBancaria) error
}

type MovimientoRepository interface {
	FindAll(ctx context.Context, cuentaID *string, soloNoConciliados bool) ([]*MovimientoBancario, error)
	FindByRango(ctx context.Context, desde, hasta time.Time) ([]*MovimientoBancario, error)
	Save(ctx context.Context, m *MovimientoBancario) error
	Conciliar(ctx context.Context, id string) error
}
