package alquiler

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	domalquiler "github.com/ipnext/admin-backend/internal/domain/alquiler"
)

var (
	ErrInmuebleNoEncontrado  = errors.New("inmueble no encontrado")
	ErrContratoNoEncontrado  = errors.New("contrato no encontrado")
)

// --- Inmueble ---

type ListInmueblesUseCase struct{ repo domalquiler.Repository }

func NewListInmueblesUseCase(repo domalquiler.Repository) *ListInmueblesUseCase {
	return &ListInmueblesUseCase{repo: repo}
}
func (uc *ListInmueblesUseCase) Execute(ctx context.Context) ([]*domalquiler.Inmueble, error) {
	return uc.repo.FindAll(ctx)
}

type GetInmuebleUseCase struct{ repo domalquiler.Repository }

func NewGetInmuebleUseCase(repo domalquiler.Repository) *GetInmuebleUseCase {
	return &GetInmuebleUseCase{repo: repo}
}
func (uc *GetInmuebleUseCase) Execute(ctx context.Context, id string) (*domalquiler.Inmueble, error) {
	i, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if i == nil {
		return nil, ErrInmuebleNoEncontrado
	}
	return i, nil
}

type CreateInmuebleRequest struct {
	Nombre          string
	Direccion       string
	Propietario     string
	Uso             domalquiler.UsoInmueble
	AlquilerMensual float64
	CBU             *string
	Alias           *string
}

type CreateInmuebleUseCase struct{ repo domalquiler.Repository }

func NewCreateInmuebleUseCase(repo domalquiler.Repository) *CreateInmuebleUseCase {
	return &CreateInmuebleUseCase{repo: repo}
}
func (uc *CreateInmuebleUseCase) Execute(ctx context.Context, req CreateInmuebleRequest) (*domalquiler.Inmueble, error) {
	i := &domalquiler.Inmueble{
		ID: uuid.NewString(), Nombre: req.Nombre, Direccion: req.Direccion,
		Propietario: req.Propietario, Uso: req.Uso, AlquilerMensual: req.AlquilerMensual,
		CBU: req.CBU, Alias: req.Alias, Estado: domalquiler.EstadoPendiente,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := uc.repo.Save(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}

type UpdateInmuebleRequest struct {
	Nombre          *string
	Direccion       *string
	Propietario     *string
	AlquilerMensual *float64
	CBU             *string
	Alias           *string
	Estado          *domalquiler.EstadoInmueble
}

type UpdateInmuebleUseCase struct{ repo domalquiler.Repository }

func NewUpdateInmuebleUseCase(repo domalquiler.Repository) *UpdateInmuebleUseCase {
	return &UpdateInmuebleUseCase{repo: repo}
}
func (uc *UpdateInmuebleUseCase) Execute(ctx context.Context, id string, req UpdateInmuebleRequest) (*domalquiler.Inmueble, error) {
	i, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if i == nil {
		return nil, ErrInmuebleNoEncontrado
	}
	if req.Nombre != nil {
		i.Nombre = *req.Nombre
	}
	if req.Direccion != nil {
		i.Direccion = *req.Direccion
	}
	if req.Propietario != nil {
		i.Propietario = *req.Propietario
	}
	if req.AlquilerMensual != nil {
		i.AlquilerMensual = *req.AlquilerMensual
	}
	if req.CBU != nil {
		i.CBU = req.CBU
	}
	if req.Alias != nil {
		i.Alias = req.Alias
	}
	if req.Estado != nil {
		i.Estado = *req.Estado
	}
	i.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}

type DeleteInmuebleUseCase struct{ repo domalquiler.Repository }

func NewDeleteInmuebleUseCase(repo domalquiler.Repository) *DeleteInmuebleUseCase {
	return &DeleteInmuebleUseCase{repo: repo}
}
func (uc *DeleteInmuebleUseCase) Execute(ctx context.Context, id string) error {
	i, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if i == nil {
		return ErrInmuebleNoEncontrado
	}
	return uc.repo.Delete(ctx, id)
}

// --- Contrato ---

type ListContratosUseCase struct{ repo domalquiler.ContratoRepository }

func NewListContratosUseCase(repo domalquiler.ContratoRepository) *ListContratosUseCase {
	return &ListContratosUseCase{repo: repo}
}
func (uc *ListContratosUseCase) Execute(ctx context.Context, inmuebleID *string) ([]*domalquiler.ContratoAlquiler, error) {
	return uc.repo.FindAll(ctx, inmuebleID)
}

type CreateContratoRequest struct {
	InmuebleID       string
	VigenciaDesde    time.Time
	VigenciaHasta    time.Time
	AjusteFrecuencia string
	MontoMensual     float64
}

type CreateContratoUseCase struct{ repo domalquiler.ContratoRepository }

func NewCreateContratoUseCase(repo domalquiler.ContratoRepository) *CreateContratoUseCase {
	return &CreateContratoUseCase{repo: repo}
}
func (uc *CreateContratoUseCase) Execute(ctx context.Context, req CreateContratoRequest) (*domalquiler.ContratoAlquiler, error) {
	c := &domalquiler.ContratoAlquiler{
		ID: uuid.NewString(), InmuebleID: req.InmuebleID,
		VigenciaDesde: req.VigenciaDesde, VigenciaHasta: req.VigenciaHasta,
		AjusteFrecuencia: req.AjusteFrecuencia, MontoMensual: req.MontoMensual,
		Estado: domalquiler.EstadoVigente, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := uc.repo.Save(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

type VencimientosUseCase struct{ repo domalquiler.ContratoRepository }

func NewVencimientosUseCase(repo domalquiler.ContratoRepository) *VencimientosUseCase {
	return &VencimientosUseCase{repo: repo}
}
func (uc *VencimientosUseCase) Execute(ctx context.Context, dias int) ([]*domalquiler.ContratoAlquiler, error) {
	return uc.repo.FindProximosAVencer(ctx, dias)
}

// --- Pago ---

type ListPagosUseCase struct{ repo domalquiler.PagoRepository }

func NewListPagosUseCase(repo domalquiler.PagoRepository) *ListPagosUseCase {
	return &ListPagosUseCase{repo: repo}
}
func (uc *ListPagosUseCase) Execute(ctx context.Context, inmuebleID *string) ([]*domalquiler.PagoAlquiler, error) {
	return uc.repo.FindAll(ctx, inmuebleID)
}

type CreatePagoRequest struct {
	InmuebleID  string
	Periodo     string
	Monto       float64
	FechaPago   *time.Time
	Comprobante *string
}

type CreatePagoUseCase struct{ repo domalquiler.PagoRepository }

func NewCreatePagoUseCase(repo domalquiler.PagoRepository) *CreatePagoUseCase {
	return &CreatePagoUseCase{repo: repo}
}
func (uc *CreatePagoUseCase) Execute(ctx context.Context, req CreatePagoRequest) (*domalquiler.PagoAlquiler, error) {
	estado := domalquiler.EstadoPendiente
	if req.FechaPago != nil {
		estado = domalquiler.EstadoPagado
	}
	p := &domalquiler.PagoAlquiler{
		ID: uuid.NewString(), InmuebleID: req.InmuebleID, Periodo: req.Periodo,
		Monto: req.Monto, FechaPago: req.FechaPago, Estado: estado,
		Comprobante: req.Comprobante, CreatedAt: time.Now(),
	}
	if err := uc.repo.Save(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
