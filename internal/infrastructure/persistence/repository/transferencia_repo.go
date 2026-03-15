package repository

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/ipnext/admin-backend/internal/domain/transferencia"
	"gorm.io/gorm"
)

type transferenciaModel struct {
	ID               string     `gorm:"primaryKey;type:char(36)"`
	Beneficiario     string     `gorm:"not null"`
	CBU              *string
	Alias            *string
	Categoria        string     `gorm:"not null"`
	Monto            float64    `gorm:"type:decimal(15,2);not null"`
	Moneda           string     `gorm:"type:enum('ARS','USD');default:ARS"`
	FechaPago        time.Time  `gorm:"not null"`
	FechaVencimiento *time.Time
	Frecuencia       string     `gorm:"type:varchar(20);not null"`
	Estado           string     `gorm:"type:varchar(20);not null"`
	MetodoPago       string     `gorm:"type:varchar(20);not null"`
	Notas            *string
	ProveedorID      *string    `gorm:"type:char(36)"`
	CreadoPor        string     `gorm:"type:char(36);not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (transferenciaModel) TableName() string { return "transferencias" }

type MySQLTransferenciaRepository struct {
	db *gorm.DB
}

func NewMySQLTransferenciaRepository(db *gorm.DB) transferencia.Repository {
	return &MySQLTransferenciaRepository{db: db}
}

func (r *MySQLTransferenciaRepository) FindByID(ctx context.Context, id string) (*transferencia.Transferencia, error) {
	var m transferenciaModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toTransferenciaDomain(&m), nil
}

func (r *MySQLTransferenciaRepository) FindAll(ctx context.Context, f transferencia.Filtros) (*transferencia.ListResult, error) {
	query := r.db.WithContext(ctx).Model(&transferenciaModel{})

	if f.Estado != "" {
		query = query.Where("estado = ?", f.Estado)
	}
	if f.Categoria != "" {
		query = query.Where("categoria = ?", f.Categoria)
	}
	if f.Desde != nil {
		query = query.Where("fecha_pago >= ?", f.Desde)
	}
	if f.Hasta != nil {
		query = query.Where("fecha_pago < ?", f.Hasta)
	}
	if f.Query != "" {
		query = query.Where("beneficiario LIKE ?", "%"+f.Query+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	orden := "fecha_pago"
	if f.Orden != "" {
		orden = f.Orden
	}
	dir := "asc"
	if f.Dir == "desc" {
		dir = "desc"
	}
	query = query.Order(orden + " " + dir)

	page := f.Page
	if page < 1 {
		page = 1
	}
	perPage := f.PerPage
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	query = query.Offset(offset).Limit(perPage)

	var models []transferenciaModel
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*transferencia.Transferencia, 0, len(models))
	for _, m := range models {
		items = append(items, toTransferenciaDomain(&m))
	}

	return &transferencia.ListResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: int(math.Ceil(float64(total) / float64(perPage))),
	}, nil
}

func (r *MySQLTransferenciaRepository) FindRecurrentes(ctx context.Context) ([]*transferencia.Transferencia, error) {
	var models []transferenciaModel
	if err := r.db.WithContext(ctx).Where("frecuencia != ?", "manual").Find(&models).Error; err != nil {
		return nil, err
	}
	return toTransferenciaSlice(models), nil
}

func (r *MySQLTransferenciaRepository) FindByFecha(ctx context.Context, desde, hasta time.Time) ([]*transferencia.Transferencia, error) {
	var models []transferenciaModel
	if err := r.db.WithContext(ctx).Where("fecha_pago >= ? AND fecha_pago < ?", desde, hasta).Order("fecha_pago asc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toTransferenciaSlice(models), nil
}

func (r *MySQLTransferenciaRepository) FindProximasAVencer(ctx context.Context, dias int) ([]*transferencia.Transferencia, error) {
	hasta := time.Now().AddDate(0, 0, dias)
	var models []transferenciaModel
	err := r.db.WithContext(ctx).
		Where("estado = ? AND fecha_pago <= ? AND fecha_pago >= ?", "pendiente", hasta, time.Now()).
		Order("fecha_pago asc").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toTransferenciaSlice(models), nil
}

func (r *MySQLTransferenciaRepository) Save(ctx context.Context, t *transferencia.Transferencia) error {
	return r.db.WithContext(ctx).Create(toTransferenciaModel(t)).Error
}

func (r *MySQLTransferenciaRepository) Update(ctx context.Context, t *transferencia.Transferencia) error {
	return r.db.WithContext(ctx).Save(toTransferenciaModel(t)).Error
}

func (r *MySQLTransferenciaRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&transferenciaModel{}, "id = ?", id).Error
}

func (r *MySQLTransferenciaRepository) SumByCategoria(ctx context.Context, desde, hasta time.Time) (map[string]float64, error) {
	type result struct {
		Categoria string
		Total     float64
	}
	var rows []result
	err := r.db.WithContext(ctx).
		Model(&transferenciaModel{}).
		Select("categoria, SUM(monto) as total").
		Where("fecha_pago >= ? AND fecha_pago < ? AND estado = ?", desde, hasta, "pagado").
		Group("categoria").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	m := make(map[string]float64, len(rows))
	for _, row := range rows {
		m[row.Categoria] = row.Total
	}
	return m, nil
}

func toTransferenciaDomain(m *transferenciaModel) *transferencia.Transferencia {
	return &transferencia.Transferencia{
		ID:               m.ID,
		Beneficiario:     m.Beneficiario,
		CBU:              m.CBU,
		Alias:            m.Alias,
		Categoria:        m.Categoria,
		Monto:            m.Monto,
		Moneda:           m.Moneda,
		FechaPago:        m.FechaPago,
		FechaVencimiento: m.FechaVencimiento,
		Frecuencia:       transferencia.Frecuencia(m.Frecuencia),
		Estado:           transferencia.Estado(m.Estado),
		MetodoPago:       transferencia.MetodoPago(m.MetodoPago),
		Notas:            m.Notas,
		ProveedorID:      m.ProveedorID,
		CreadoPor:        m.CreadoPor,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func toTransferenciaModel(t *transferencia.Transferencia) *transferenciaModel {
	return &transferenciaModel{
		ID:               t.ID,
		Beneficiario:     t.Beneficiario,
		CBU:              t.CBU,
		Alias:            t.Alias,
		Categoria:        t.Categoria,
		Monto:            t.Monto,
		Moneda:           t.Moneda,
		FechaPago:        t.FechaPago,
		FechaVencimiento: t.FechaVencimiento,
		Frecuencia:       string(t.Frecuencia),
		Estado:           string(t.Estado),
		MetodoPago:       string(t.MetodoPago),
		Notas:            t.Notas,
		ProveedorID:      t.ProveedorID,
		CreadoPor:        t.CreadoPor,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}
}

func toTransferenciaSlice(models []transferenciaModel) []*transferencia.Transferencia {
	result := make([]*transferencia.Transferencia, 0, len(models))
	for _, m := range models {
		result = append(result, toTransferenciaDomain(&m))
	}
	return result
}
