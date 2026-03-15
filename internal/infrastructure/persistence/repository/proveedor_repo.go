package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ipnext/admin-backend/internal/domain/proveedor"
	"gorm.io/gorm"
)

type proveedorModel struct {
	ID        string    `gorm:"primaryKey;type:char(36)"`
	Nombre    string    `gorm:"not null"`
	CUIT      string    `gorm:"column:cuit;not null"`
	CBU       *string
	Alias     *string
	Email     *string
	Categoria string    `gorm:"not null"`
	SitioWeb  *string
	Activo    bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (proveedorModel) TableName() string { return "proveedores" }

type contratoProveedorModel struct {
	ID            string     `gorm:"primaryKey;type:char(36)"`
	Codigo        string     `gorm:"uniqueIndex;not null"`
	ProveedorID   string     `gorm:"type:char(36);not null"`
	Descripcion   string     `gorm:"type:text;not null;default:''"`
	VigenciaDesde time.Time  `gorm:"not null"`
	VigenciaHasta time.Time  `gorm:"not null"`
	MontoAnual    float64    `gorm:"column:monto_mensual;type:decimal(15,2);not null"`
	Estado        string     `gorm:"type:varchar(30);not null;default:activo"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (contratoProveedorModel) TableName() string { return "contratos_proveedor" }

// --- Proveedor Repository ---

type MySQLProveedorRepository struct{ db *gorm.DB }

func NewMySQLProveedorRepository(db *gorm.DB) proveedor.Repository {
	return &MySQLProveedorRepository{db: db}
}

func (r *MySQLProveedorRepository) FindByID(ctx context.Context, id string) (*proveedor.Proveedor, error) {
	var m proveedorModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toProveedorDomain(&m), nil
}

func (r *MySQLProveedorRepository) FindAll(ctx context.Context, query string, soloActivos bool) ([]*proveedor.Proveedor, error) {
	q := r.db.WithContext(ctx)
	if soloActivos {
		q = q.Where("activo = ?", true)
	}
	if query != "" {
		q = q.Where("nombre LIKE ? OR cuit LIKE ?", "%"+query+"%", "%"+query+"%")
	}
	var models []proveedorModel
	if err := q.Order("nombre asc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*proveedor.Proveedor, 0, len(models))
	for _, m := range models {
		result = append(result, toProveedorDomain(&m))
	}
	return result, nil
}

func (r *MySQLProveedorRepository) Save(ctx context.Context, p *proveedor.Proveedor) error {
	return r.db.WithContext(ctx).Create(toProveedorModel(p)).Error
}

func (r *MySQLProveedorRepository) Update(ctx context.Context, p *proveedor.Proveedor) error {
	return r.db.WithContext(ctx).Save(toProveedorModel(p)).Error
}

func (r *MySQLProveedorRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&proveedorModel{}).Where("id = ?", id).Update("activo", false).Error
}

// --- Contrato Repository ---

type MySQLContratoProveedorRepository struct{ db *gorm.DB }

func NewMySQLContratoProveedorRepository(db *gorm.DB) proveedor.ContratoRepository {
	return &MySQLContratoProveedorRepository{db: db}
}

func (r *MySQLContratoProveedorRepository) FindByID(ctx context.Context, id string) (*proveedor.ContratoProveedor, error) {
	var m contratoProveedorModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toContratoDomain(&m), nil
}

func (r *MySQLContratoProveedorRepository) FindAll(ctx context.Context, proveedorID *string) ([]*proveedor.ContratoProveedor, error) {
	q := r.db.WithContext(ctx)
	if proveedorID != nil {
		q = q.Where("proveedor_id = ?", *proveedorID)
	}
	var models []contratoProveedorModel
	if err := q.Order("vigencia_hasta desc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*proveedor.ContratoProveedor, 0, len(models))
	for _, m := range models {
		result = append(result, toContratoDomain(&m))
	}
	return result, nil
}

func (r *MySQLContratoProveedorRepository) Save(ctx context.Context, c *proveedor.ContratoProveedor) error {
	return r.db.WithContext(ctx).Create(toContratoModel(c)).Error
}

func (r *MySQLContratoProveedorRepository) Update(ctx context.Context, c *proveedor.ContratoProveedor) error {
	return r.db.WithContext(ctx).Save(toContratoModel(c)).Error
}

// --- Ranking Repository ---

type MySQLRankingRepository struct{ db *gorm.DB }

func NewMySQLRankingRepository(db *gorm.DB) proveedor.RankingRepository {
	return &MySQLRankingRepository{db: db}
}

func (r *MySQLRankingRepository) GetRanking(ctx context.Context, limit int) ([]*proveedor.RankingItem, error) {
	type row struct {
		ProveedorID     string
		NombreProveedor string
		TotalPagado     float64
		CantidadPagos   int64
	}
	var rows []row
	err := r.db.WithContext(ctx).
		Table("transferencias t").
		Select("t.proveedor_id, p.nombre as nombre_proveedor, SUM(t.monto) as total_pagado, COUNT(t.id) as cantidad_pagos").
		Joins("JOIN proveedores p ON p.id = t.proveedor_id").
		Where("t.proveedor_id IS NOT NULL AND t.estado = ?", "pagado").
		Group("t.proveedor_id, p.nombre").
		Order("total_pagado desc").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make([]*proveedor.RankingItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, &proveedor.RankingItem{
			ProveedorID:     row.ProveedorID,
			NombreProveedor: row.NombreProveedor,
			TotalPagado:     row.TotalPagado,
			CantidadPagos:   row.CantidadPagos,
		})
	}
	return result, nil
}

// --- Mappers ---

func toProveedorDomain(m *proveedorModel) *proveedor.Proveedor {
	return &proveedor.Proveedor{ID: m.ID, Nombre: m.Nombre, CUIT: m.CUIT, CBU: m.CBU,
		Alias: m.Alias, Email: m.Email, Categoria: m.Categoria, SitioWeb: m.SitioWeb,
		Activo: m.Activo, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}
func toProveedorModel(p *proveedor.Proveedor) *proveedorModel {
	return &proveedorModel{ID: p.ID, Nombre: p.Nombre, CUIT: p.CUIT, CBU: p.CBU,
		Alias: p.Alias, Email: p.Email, Categoria: p.Categoria, SitioWeb: p.SitioWeb,
		Activo: p.Activo, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}
func toContratoDomain(m *contratoProveedorModel) *proveedor.ContratoProveedor {
	return &proveedor.ContratoProveedor{ID: m.ID, Codigo: m.Codigo, ProveedorID: m.ProveedorID,
		Descripcion: m.Descripcion, VigenciaDesde: m.VigenciaDesde, VigenciaHasta: m.VigenciaHasta,
		MontoAnual: m.MontoAnual, Estado: proveedor.EstadoContrato(m.Estado),
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}
func toContratoModel(c *proveedor.ContratoProveedor) *contratoProveedorModel {
	return &contratoProveedorModel{ID: c.ID, Codigo: c.Codigo, ProveedorID: c.ProveedorID,
		Descripcion: c.Descripcion, VigenciaDesde: c.VigenciaDesde, VigenciaHasta: c.VigenciaHasta,
		MontoAnual: c.MontoAnual, Estado: string(c.Estado),
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
}
