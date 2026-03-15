package transferencia

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ipnext/admin-backend/internal/domain/transferencia"
)

var ErrNoEncontrada = errors.New("transferencia no encontrada")

// --- List ---

type ListUseCase struct{ repo transferencia.Repository }

func NewListUseCase(repo transferencia.Repository) *ListUseCase {
	return &ListUseCase{repo: repo}
}

func (uc *ListUseCase) Execute(ctx context.Context, f transferencia.Filtros) (*transferencia.ListResult, error) {
	return uc.repo.FindAll(ctx, f)
}

// --- Get ---

type GetUseCase struct{ repo transferencia.Repository }

func NewGetUseCase(repo transferencia.Repository) *GetUseCase {
	return &GetUseCase{repo: repo}
}

func (uc *GetUseCase) Execute(ctx context.Context, id string) (*transferencia.Transferencia, error) {
	t, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrNoEncontrada
	}
	return t, nil
}

// --- Create ---

type CreateRequest struct {
	Beneficiario     string
	CBU              *string
	Alias            *string
	Categoria        string
	Monto            float64
	Moneda           string
	FechaPago        time.Time
	FechaVencimiento *time.Time
	Frecuencia       transferencia.Frecuencia
	Estado           transferencia.Estado
	MetodoPago       transferencia.MetodoPago
	Notas            *string
	ProveedorID      *string
	CreadoPor        string
}

type CreateUseCase struct{ repo transferencia.Repository }

func NewCreateUseCase(repo transferencia.Repository) *CreateUseCase {
	return &CreateUseCase{repo: repo}
}

func (uc *CreateUseCase) Execute(ctx context.Context, req CreateRequest) (*transferencia.Transferencia, error) {
	t := &transferencia.Transferencia{
		ID:               uuid.NewString(),
		Beneficiario:     req.Beneficiario,
		CBU:              req.CBU,
		Alias:            req.Alias,
		Categoria:        req.Categoria,
		Monto:            req.Monto,
		Moneda:           req.Moneda,
		FechaPago:        req.FechaPago,
		FechaVencimiento: req.FechaVencimiento,
		Frecuencia:       req.Frecuencia,
		Estado:           transferencia.EstadoPendiente,
		MetodoPago:       req.MetodoPago,
		Notas:            req.Notas,
		ProveedorID:      req.ProveedorID,
		CreadoPor:        req.CreadoPor,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := uc.repo.Save(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

// --- Update ---

type UpdateRequest struct {
	Beneficiario     *string
	CBU              *string
	Alias            *string
	Categoria        *string
	Monto            *float64
	Moneda           *string
	FechaPago        *time.Time
	FechaVencimiento *time.Time
	Frecuencia       *string
	MetodoPago       *string
	Notas            *string
	ProveedorID      *string
}

type UpdateUseCase struct{ repo transferencia.Repository }

func NewUpdateUseCase(repo transferencia.Repository) *UpdateUseCase {
	return &UpdateUseCase{repo: repo}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, id string, req UpdateRequest) (*transferencia.Transferencia, error) {
	t, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrNoEncontrada
	}
	if req.Beneficiario != nil {
		t.Beneficiario = *req.Beneficiario
	}
	if req.CBU != nil {
		t.CBU = req.CBU
	}
	if req.Alias != nil {
		t.Alias = req.Alias
	}
	if req.Categoria != nil {
		t.Categoria = *req.Categoria
	}
	if req.Monto != nil {
		t.Monto = *req.Monto
	}
	if req.Moneda != nil {
		t.Moneda = *req.Moneda
	}
	if req.FechaPago != nil {
		t.FechaPago = *req.FechaPago
	}
	if req.FechaVencimiento != nil {
		t.FechaVencimiento = req.FechaVencimiento
	}
	if req.Frecuencia != nil {
		t.Frecuencia = transferencia.Frecuencia(*req.Frecuencia)
	}
	if req.MetodoPago != nil {
		t.MetodoPago = transferencia.MetodoPago(*req.MetodoPago)
	}
	if req.Notas != nil {
		t.Notas = req.Notas
	}
	if req.ProveedorID != nil {
		t.ProveedorID = req.ProveedorID
	}
	t.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

// --- Delete ---

type DeleteUseCase struct{ repo transferencia.Repository }

func NewDeleteUseCase(repo transferencia.Repository) *DeleteUseCase {
	return &DeleteUseCase{repo: repo}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, id string) error {
	t, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if t == nil {
		return ErrNoEncontrada
	}
	return uc.repo.Delete(ctx, id)
}

// --- CambiarEstado ---

type CambiarEstadoUseCase struct{ repo transferencia.Repository }

func NewCambiarEstadoUseCase(repo transferencia.Repository) *CambiarEstadoUseCase {
	return &CambiarEstadoUseCase{repo: repo}
}

func (uc *CambiarEstadoUseCase) Execute(ctx context.Context, id string, estado transferencia.Estado) (*transferencia.Transferencia, error) {
	t, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrNoEncontrada
	}
	t.Estado = estado
	t.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

// --- Calendario ---

type CalendarioUseCase struct{ repo transferencia.Repository }

func NewCalendarioUseCase(repo transferencia.Repository) *CalendarioUseCase {
	return &CalendarioUseCase{repo: repo}
}

func (uc *CalendarioUseCase) Execute(ctx context.Context, desde, hasta time.Time) ([]*transferencia.Transferencia, error) {
	return uc.repo.FindByFecha(ctx, desde, hasta)
}

// --- Recurrentes ---

type RecurrentesUseCase struct{ repo transferencia.Repository }

func NewRecurrentesUseCase(repo transferencia.Repository) *RecurrentesUseCase {
	return &RecurrentesUseCase{repo: repo}
}

func (uc *RecurrentesUseCase) Execute(ctx context.Context) ([]*transferencia.Transferencia, error) {
	return uc.repo.FindRecurrentes(ctx)
}
