package proveedor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	domproveedor "github.com/ipnext/admin-backend/internal/domain/proveedor"
)

var (
	ErrProveedorNoEncontrado = errors.New("proveedor no encontrado")
	ErrContratoNoEncontrado  = errors.New("contrato no encontrado")
)

// --- List ---

type ListUseCase struct{ repo domproveedor.Repository }

func NewListUseCase(repo domproveedor.Repository) *ListUseCase {
	return &ListUseCase{repo: repo}
}
func (uc *ListUseCase) Execute(ctx context.Context, query string) ([]*domproveedor.Proveedor, error) {
	return uc.repo.FindAll(ctx, query, true)
}

// --- Get ---

type GetUseCase struct{ repo domproveedor.Repository }

func NewGetUseCase(repo domproveedor.Repository) *GetUseCase { return &GetUseCase{repo: repo} }
func (uc *GetUseCase) Execute(ctx context.Context, id string) (*domproveedor.Proveedor, error) {
	p, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProveedorNoEncontrado
	}
	return p, nil
}

// --- Create ---

type CreateRequest struct {
	Nombre    string
	CUIT      string
	CBU       *string
	Alias     *string
	Email     *string
	Categoria string
	SitioWeb  *string
}

type CreateUseCase struct{ repo domproveedor.Repository }

func NewCreateUseCase(repo domproveedor.Repository) *CreateUseCase {
	return &CreateUseCase{repo: repo}
}
func (uc *CreateUseCase) Execute(ctx context.Context, req CreateRequest) (*domproveedor.Proveedor, error) {
	p := &domproveedor.Proveedor{
		ID:        uuid.NewString(),
		Nombre:    req.Nombre,
		CUIT:      req.CUIT,
		CBU:       req.CBU,
		Alias:     req.Alias,
		Email:     req.Email,
		Categoria: req.Categoria,
		SitioWeb:  req.SitioWeb,
		Activo:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := uc.repo.Save(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// --- Update ---

type UpdateRequest struct {
	Nombre    *string
	CUIT      *string
	CBU       *string
	Alias     *string
	Email     *string
	Categoria *string
	SitioWeb  *string
}

type UpdateUseCase struct{ repo domproveedor.Repository }

func NewUpdateUseCase(repo domproveedor.Repository) *UpdateUseCase {
	return &UpdateUseCase{repo: repo}
}
func (uc *UpdateUseCase) Execute(ctx context.Context, id string, req UpdateRequest) (*domproveedor.Proveedor, error) {
	p, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProveedorNoEncontrado
	}
	if req.Nombre != nil {
		p.Nombre = *req.Nombre
	}
	if req.CUIT != nil {
		p.CUIT = *req.CUIT
	}
	if req.CBU != nil {
		p.CBU = req.CBU
	}
	if req.Alias != nil {
		p.Alias = req.Alias
	}
	if req.Email != nil {
		p.Email = req.Email
	}
	if req.Categoria != nil {
		p.Categoria = *req.Categoria
	}
	if req.SitioWeb != nil {
		p.SitioWeb = req.SitioWeb
	}
	p.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// --- Delete ---

type DeleteUseCase struct{ repo domproveedor.Repository }

func NewDeleteUseCase(repo domproveedor.Repository) *DeleteUseCase {
	return &DeleteUseCase{repo: repo}
}
func (uc *DeleteUseCase) Execute(ctx context.Context, id string) error {
	p, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return ErrProveedorNoEncontrado
	}
	return uc.repo.Delete(ctx, id)
}

// --- List Contratos ---

type ListContratosUseCase struct{ repo domproveedor.ContratoRepository }

func NewListContratosUseCase(repo domproveedor.ContratoRepository) *ListContratosUseCase {
	return &ListContratosUseCase{repo: repo}
}
func (uc *ListContratosUseCase) Execute(ctx context.Context, proveedorID *string) ([]*domproveedor.ContratoProveedor, error) {
	return uc.repo.FindAll(ctx, proveedorID)
}

// --- Create Contrato ---

type CreateContratoRequest struct {
	ProveedorID   string
	Descripcion   string
	VigenciaDesde time.Time
	VigenciaHasta time.Time
	MontoAnual    float64
}

type CreateContratoUseCase struct {
	repo         domproveedor.ContratoRepository
	provRepo     domproveedor.Repository
	contratoCount int
}

func NewCreateContratoUseCase(repo domproveedor.ContratoRepository, provRepo domproveedor.Repository) *CreateContratoUseCase {
	return &CreateContratoUseCase{repo: repo, provRepo: provRepo}
}
func (uc *CreateContratoUseCase) Execute(ctx context.Context, req CreateContratoRequest) (*domproveedor.ContratoProveedor, error) {
	p, err := uc.provRepo.FindByID(ctx, req.ProveedorID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProveedorNoEncontrado
	}
	codigo := fmt.Sprintf("CTR-%d-%04d", time.Now().Year(), time.Now().UnixMilli()%10000)
	c := &domproveedor.ContratoProveedor{
		ID:            uuid.NewString(),
		Codigo:        codigo,
		ProveedorID:   req.ProveedorID,
		Descripcion:   req.Descripcion,
		VigenciaDesde: req.VigenciaDesde,
		VigenciaHasta: req.VigenciaHasta,
		MontoAnual:    req.MontoAnual,
		Estado:        domproveedor.EstadoActivo,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := uc.repo.Save(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// --- Update Contrato ---

type UpdateContratoUseCase struct{ repo domproveedor.ContratoRepository }

func NewUpdateContratoUseCase(repo domproveedor.ContratoRepository) *UpdateContratoUseCase {
	return &UpdateContratoUseCase{repo: repo}
}
func (uc *UpdateContratoUseCase) Execute(ctx context.Context, id string, estado domproveedor.EstadoContrato) (*domproveedor.ContratoProveedor, error) {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrContratoNoEncontrado
	}
	c.Estado = estado
	c.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// --- Ranking ---

type GetRankingUseCase struct{ repo domproveedor.RankingRepository }

func NewGetRankingUseCase(repo domproveedor.RankingRepository) *GetRankingUseCase {
	return &GetRankingUseCase{repo: repo}
}
func (uc *GetRankingUseCase) Execute(ctx context.Context) ([]*domproveedor.RankingItem, error) {
	return uc.repo.GetRanking(ctx, 20)
}
