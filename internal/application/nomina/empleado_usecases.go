package nomina

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ipnext/admin-backend/internal/domain/empleado"
)

var ErrEmpleadoNoEncontrado = errors.New("empleado no encontrado")

type EmpleadoKPIs struct {
	TotalActivos  int64   `json:"totalActivos"`
	CostoNomina   float64 `json:"costoNomina"`
}

// --- List ---

type ListEmpleadosUseCase struct{ repo empleado.Repository }

func NewListEmpleadosUseCase(repo empleado.Repository) *ListEmpleadosUseCase {
	return &ListEmpleadosUseCase{repo: repo}
}
func (uc *ListEmpleadosUseCase) Execute(ctx context.Context) ([]*empleado.Empleado, error) {
	return uc.repo.FindAll(ctx, true)
}

// --- Get ---

type GetEmpleadoUseCase struct{ repo empleado.Repository }

func NewGetEmpleadoUseCase(repo empleado.Repository) *GetEmpleadoUseCase {
	return &GetEmpleadoUseCase{repo: repo}
}
func (uc *GetEmpleadoUseCase) Execute(ctx context.Context, id string) (*empleado.Empleado, error) {
	e, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, ErrEmpleadoNoEncontrado
	}
	return e, nil
}

// --- Create ---

type CreateEmpleadoRequest struct {
	Nombre       string
	Puesto       string
	Area         string
	Rol          string
	SueldoBruto  float64
	ObraSocial   string
	FechaIngreso time.Time
	Avatar       *string
}

type CreateEmpleadoUseCase struct{ repo empleado.Repository }

func NewCreateEmpleadoUseCase(repo empleado.Repository) *CreateEmpleadoUseCase {
	return &CreateEmpleadoUseCase{repo: repo}
}
func (uc *CreateEmpleadoUseCase) Execute(ctx context.Context, req CreateEmpleadoRequest) (*empleado.Empleado, error) {
	e := &empleado.Empleado{
		ID:           uuid.NewString(),
		Nombre:       req.Nombre,
		Puesto:       req.Puesto,
		Area:         req.Area,
		Rol:          req.Rol,
		SueldoBruto:  req.SueldoBruto,
		ObraSocial:   req.ObraSocial,
		Activo:       true,
		FechaIngreso: req.FechaIngreso,
		Avatar:       req.Avatar,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := uc.repo.Save(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

// --- Update ---

type UpdateEmpleadoRequest struct {
	Nombre      *string
	Puesto      *string
	Area        *string
	Rol         *string
	SueldoBruto *float64
	ObraSocial  *string
	Avatar      *string
}

type UpdateEmpleadoUseCase struct{ repo empleado.Repository }

func NewUpdateEmpleadoUseCase(repo empleado.Repository) *UpdateEmpleadoUseCase {
	return &UpdateEmpleadoUseCase{repo: repo}
}
func (uc *UpdateEmpleadoUseCase) Execute(ctx context.Context, id string, req UpdateEmpleadoRequest) (*empleado.Empleado, error) {
	e, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, ErrEmpleadoNoEncontrado
	}
	if req.Nombre != nil {
		e.Nombre = *req.Nombre
	}
	if req.Puesto != nil {
		e.Puesto = *req.Puesto
	}
	if req.Area != nil {
		e.Area = *req.Area
	}
	if req.Rol != nil {
		e.Rol = *req.Rol
	}
	if req.SueldoBruto != nil {
		e.SueldoBruto = *req.SueldoBruto
	}
	if req.ObraSocial != nil {
		e.ObraSocial = *req.ObraSocial
	}
	if req.Avatar != nil {
		e.Avatar = req.Avatar
	}
	e.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

// --- Delete ---

type DeleteEmpleadoUseCase struct{ repo empleado.Repository }

func NewDeleteEmpleadoUseCase(repo empleado.Repository) *DeleteEmpleadoUseCase {
	return &DeleteEmpleadoUseCase{repo: repo}
}
func (uc *DeleteEmpleadoUseCase) Execute(ctx context.Context, id string) error {
	e, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if e == nil {
		return ErrEmpleadoNoEncontrado
	}
	return uc.repo.Delete(ctx, id)
}

// --- KPIs ---

type GetEmpleadoKPIsUseCase struct{ repo empleado.Repository }

func NewGetEmpleadoKPIsUseCase(repo empleado.Repository) *GetEmpleadoKPIsUseCase {
	return &GetEmpleadoKPIsUseCase{repo: repo}
}
func (uc *GetEmpleadoKPIsUseCase) Execute(ctx context.Context) (*EmpleadoKPIs, error) {
	total, err := uc.repo.CountActivos(ctx)
	if err != nil {
		return nil, err
	}
	costo, err := uc.repo.SumSueldos(ctx)
	if err != nil {
		return nil, err
	}
	return &EmpleadoKPIs{TotalActivos: total, CostoNomina: costo}, nil
}
