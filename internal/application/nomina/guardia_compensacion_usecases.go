package nomina

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ipnext/admin-backend/internal/domain/empleado"
)

// --- Guardias ---

type ListGuardiasUseCase struct{ repo empleado.GuardiaRepository }

func NewListGuardiasUseCase(repo empleado.GuardiaRepository) *ListGuardiasUseCase {
	return &ListGuardiasUseCase{repo: repo}
}
func (uc *ListGuardiasUseCase) Execute(ctx context.Context, empleadoID *string) ([]*empleado.Guardia, error) {
	return uc.repo.FindAll(ctx, empleadoID)
}

type CreateGuardiaRequest struct {
	EmpleadoID string
	Fecha      time.Time
	Horas      float64
	Monto      float64
	Notas      *string
}

type CreateGuardiaUseCase struct{ repo empleado.GuardiaRepository }

func NewCreateGuardiaUseCase(repo empleado.GuardiaRepository) *CreateGuardiaUseCase {
	return &CreateGuardiaUseCase{repo: repo}
}
func (uc *CreateGuardiaUseCase) Execute(ctx context.Context, req CreateGuardiaRequest) (*empleado.Guardia, error) {
	g := &empleado.Guardia{
		ID:         uuid.NewString(),
		EmpleadoID: req.EmpleadoID,
		Fecha:      req.Fecha,
		Horas:      req.Horas,
		Monto:      req.Monto,
		Notas:      req.Notas,
		CreatedAt:  time.Now(),
	}
	if err := uc.repo.Save(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}

// --- Compensaciones ---

type ListCompensacionesUseCase struct{ repo empleado.CompensacionRepository }

func NewListCompensacionesUseCase(repo empleado.CompensacionRepository) *ListCompensacionesUseCase {
	return &ListCompensacionesUseCase{repo: repo}
}
func (uc *ListCompensacionesUseCase) Execute(ctx context.Context, empleadoID *string) ([]*empleado.Compensacion, error) {
	return uc.repo.FindAll(ctx, empleadoID)
}

type CreateCompensacionRequest struct {
	EmpleadoID  string
	Tipo        empleado.TipoCompensacion
	Monto       float64
	Fecha       time.Time
	Descripcion *string
}

type CreateCompensacionUseCase struct{ repo empleado.CompensacionRepository }

func NewCreateCompensacionUseCase(repo empleado.CompensacionRepository) *CreateCompensacionUseCase {
	return &CreateCompensacionUseCase{repo: repo}
}
func (uc *CreateCompensacionUseCase) Execute(ctx context.Context, req CreateCompensacionRequest) (*empleado.Compensacion, error) {
	c := &empleado.Compensacion{
		ID:          uuid.NewString(),
		EmpleadoID:  req.EmpleadoID,
		Tipo:        req.Tipo,
		Monto:       req.Monto,
		Fecha:       req.Fecha,
		Descripcion: req.Descripcion,
		CreatedAt:   time.Now(),
	}
	if err := uc.repo.Save(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}
