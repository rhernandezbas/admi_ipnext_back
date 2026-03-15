package empleado

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Empleado, error)
	FindAll(ctx context.Context, soloActivos bool) ([]*Empleado, error)
	Save(ctx context.Context, e *Empleado) error
	Update(ctx context.Context, e *Empleado) error
	Delete(ctx context.Context, id string) error
	CountActivos(ctx context.Context) (int64, error)
	SumSueldos(ctx context.Context) (float64, error)
}

type LiquidacionRepository interface {
	FindByID(ctx context.Context, id string) (*Liquidacion, error)
	FindByPeriodo(ctx context.Context, periodo string) ([]*Liquidacion, error)
	FindByEmpleado(ctx context.Context, empleadoID string) ([]*Liquidacion, error)
	Save(ctx context.Context, l *Liquidacion) error
	Update(ctx context.Context, l *Liquidacion) error
}

type GuardiaRepository interface {
	FindAll(ctx context.Context, empleadoID *string) ([]*Guardia, error)
	Save(ctx context.Context, g *Guardia) error
}

type CompensacionRepository interface {
	FindAll(ctx context.Context, empleadoID *string) ([]*Compensacion, error)
	Save(ctx context.Context, c *Compensacion) error
}
