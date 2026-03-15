package servicio

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	domservicio "github.com/ipnext/admin-backend/internal/domain/servicio"
)

var ErrServicioNoEncontrado = errors.New("servicio no encontrado")

type ListUseCase struct{ repo domservicio.Repository }

func NewListUseCase(repo domservicio.Repository) *ListUseCase { return &ListUseCase{repo: repo} }
func (uc *ListUseCase) Execute(ctx context.Context, tipo *domservicio.Tipo) ([]*domservicio.Servicio, error) {
	return uc.repo.FindAll(ctx, tipo)
}

type GetKPIsUseCase struct{ repo domservicio.Repository }

func NewGetKPIsUseCase(repo domservicio.Repository) *GetKPIsUseCase {
	return &GetKPIsUseCase{repo: repo}
}
func (uc *GetKPIsUseCase) Execute(ctx context.Context) (*domservicio.KPIs, error) {
	return uc.repo.GetKPIs(ctx)
}

type GetUseCase struct{ repo domservicio.Repository }

func NewGetUseCase(repo domservicio.Repository) *GetUseCase { return &GetUseCase{repo: repo} }
func (uc *GetUseCase) Execute(ctx context.Context, id string) (*domservicio.Servicio, error) {
	s, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, ErrServicioNoEncontrado
	}
	return s, nil
}

type CreateRequest struct {
	Nombre       string
	Tipo         domservicio.Tipo
	Proveedor    string
	CostoMensual float64
	VtoFactura   *time.Time
	Renovacion   *time.Time
	Metadata     map[string]interface{}
}

type CreateUseCase struct{ repo domservicio.Repository }

func NewCreateUseCase(repo domservicio.Repository) *CreateUseCase { return &CreateUseCase{repo: repo} }
func (uc *CreateUseCase) Execute(ctx context.Context, req CreateRequest) (*domservicio.Servicio, error) {
	s := &domservicio.Servicio{
		ID: uuid.NewString(), Nombre: req.Nombre, Tipo: req.Tipo, Proveedor: req.Proveedor,
		CostoMensual: req.CostoMensual, VtoFactura: req.VtoFactura, Renovacion: req.Renovacion,
		Estado: domservicio.EstadoActivo, Metadata: req.Metadata,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := uc.repo.Save(ctx, s); err != nil {
		return nil, err
	}
	return s, nil
}

type UpdateRequest struct {
	Nombre       *string
	Proveedor    *string
	CostoMensual *float64
	VtoFactura   *time.Time
	Renovacion   *time.Time
	Estado       *domservicio.Estado
}

type UpdateUseCase struct{ repo domservicio.Repository }

func NewUpdateUseCase(repo domservicio.Repository) *UpdateUseCase { return &UpdateUseCase{repo: repo} }
func (uc *UpdateUseCase) Execute(ctx context.Context, id string, req UpdateRequest) (*domservicio.Servicio, error) {
	s, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, ErrServicioNoEncontrado
	}
	if req.Nombre != nil {
		s.Nombre = *req.Nombre
	}
	if req.Proveedor != nil {
		s.Proveedor = *req.Proveedor
	}
	if req.CostoMensual != nil {
		s.CostoMensual = *req.CostoMensual
	}
	if req.VtoFactura != nil {
		s.VtoFactura = req.VtoFactura
	}
	if req.Renovacion != nil {
		s.Renovacion = req.Renovacion
	}
	if req.Estado != nil {
		s.Estado = *req.Estado
	}
	s.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, s); err != nil {
		return nil, err
	}
	return s, nil
}

type DeleteUseCase struct{ repo domservicio.Repository }

func NewDeleteUseCase(repo domservicio.Repository) *DeleteUseCase { return &DeleteUseCase{repo: repo} }
func (uc *DeleteUseCase) Execute(ctx context.Context, id string) error {
	s, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if s == nil {
		return ErrServicioNoEncontrado
	}
	return uc.repo.Delete(ctx, id)
}
