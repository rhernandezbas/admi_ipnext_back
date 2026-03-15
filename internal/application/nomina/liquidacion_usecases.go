package nomina

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ipnext/admin-backend/internal/domain/empleado"
)

var (
	ErrLiquidacionNoEncontrada   = errors.New("liquidación no encontrada")
	ErrLiquidacionYaAprobada     = errors.New("la liquidación ya fue aprobada o pagada")
)

// --- List ---

type ListLiquidacionesUseCase struct{ repo empleado.LiquidacionRepository }

func NewListLiquidacionesUseCase(repo empleado.LiquidacionRepository) *ListLiquidacionesUseCase {
	return &ListLiquidacionesUseCase{repo: repo}
}
func (uc *ListLiquidacionesUseCase) Execute(ctx context.Context, periodo string) ([]*empleado.Liquidacion, error) {
	return uc.repo.FindByPeriodo(ctx, periodo)
}

// --- Create ---

type CreateLiquidacionRequest struct {
	EmpleadoID  string
	Periodo     string
	SueldoBruto float64
	Deducciones float64
}

type CreateLiquidacionUseCase struct {
	repo         empleado.LiquidacionRepository
	empleadoRepo empleado.Repository
}

func NewCreateLiquidacionUseCase(repo empleado.LiquidacionRepository, empleadoRepo empleado.Repository) *CreateLiquidacionUseCase {
	return &CreateLiquidacionUseCase{repo: repo, empleadoRepo: empleadoRepo}
}

func (uc *CreateLiquidacionUseCase) Execute(ctx context.Context, req CreateLiquidacionRequest) (*empleado.Liquidacion, error) {
	e, err := uc.empleadoRepo.FindByID(ctx, req.EmpleadoID)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, ErrEmpleadoNoEncontrado
	}

	l := &empleado.Liquidacion{
		ID:          uuid.NewString(),
		EmpleadoID:  req.EmpleadoID,
		Periodo:     req.Periodo,
		SueldoBruto: req.SueldoBruto,
		Deducciones: req.Deducciones,
		Estado:      empleado.EstadoBorrador,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	l.Calcular()

	if err := uc.repo.Save(ctx, l); err != nil {
		return nil, err
	}
	return l, nil
}

// --- Aprobar ---

type AprobarLiquidacionUseCase struct {
	repo empleado.LiquidacionRepository
}

func NewAprobarLiquidacionUseCase(repo empleado.LiquidacionRepository) *AprobarLiquidacionUseCase {
	return &AprobarLiquidacionUseCase{repo: repo}
}

func (uc *AprobarLiquidacionUseCase) Execute(ctx context.Context, id, aprobadoPor string) (*empleado.Liquidacion, error) {
	l, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, ErrLiquidacionNoEncontrada
	}
	if l.Estado != empleado.EstadoBorrador {
		return nil, ErrLiquidacionYaAprobada
	}
	l.Estado = empleado.EstadoAprobada
	l.AprobadoPor = &aprobadoPor
	l.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, l); err != nil {
		return nil, err
	}
	return l, nil
}
