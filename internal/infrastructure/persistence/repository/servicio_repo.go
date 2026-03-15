package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	domservicio "github.com/ipnext/admin-backend/internal/domain/servicio"
	"gorm.io/gorm"
)

type servicioModel struct {
	ID           string    `gorm:"primaryKey;type:char(36)"`
	Nombre       string    `gorm:"not null"`
	Tipo         string    `gorm:"type:varchar(30);not null"`
	Proveedor    string    `gorm:"not null"`
	CostoMensual float64   `gorm:"type:decimal(15,2);not null"`
	VtoFactura   *time.Time
	Renovacion   *time.Time
	Estado       string    `gorm:"type:varchar(30);not null;default:activo"`
	Metadata     *string   `gorm:"type:json"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (servicioModel) TableName() string { return "servicios" }

type MySQLServicioRepository struct{ db *gorm.DB }

func NewMySQLServicioRepository(db *gorm.DB) domservicio.Repository {
	return &MySQLServicioRepository{db: db}
}

func (r *MySQLServicioRepository) FindByID(ctx context.Context, id string) (*domservicio.Servicio, error) {
	var m servicioModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toServicioDomain(&m)
}

func (r *MySQLServicioRepository) FindAll(ctx context.Context, tipo *domservicio.Tipo) ([]*domservicio.Servicio, error) {
	q := r.db.WithContext(ctx)
	if tipo != nil {
		q = q.Where("tipo = ?", string(*tipo))
	}
	var models []servicioModel
	if err := q.Order("nombre asc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*domservicio.Servicio, 0, len(models))
	for _, m := range models {
		s, err := toServicioDomain(&m)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (r *MySQLServicioRepository) Save(ctx context.Context, s *domservicio.Servicio) error {
	m, err := toServicioModel(s)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *MySQLServicioRepository) Update(ctx context.Context, s *domservicio.Servicio) error {
	m, err := toServicioModel(s)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *MySQLServicioRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&servicioModel{}).Where("id = ?", id).Update("estado", "inactivo").Error
}

func (r *MySQLServicioRepository) GetKPIs(ctx context.Context) (*domservicio.KPIs, error) {
	type row struct {
		Tipo  string
		Total float64
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&servicioModel{}).
		Select("tipo, SUM(costo_mensual) as total").
		Where("estado = ?", "activo").
		Group("tipo").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	var total float64
	var count int64
	porTipo := make(map[domservicio.Tipo]float64)
	for _, row := range rows {
		porTipo[domservicio.Tipo(row.Tipo)] = row.Total
		total += row.Total
		count++
	}
	var proximos int64
	r.db.WithContext(ctx).Model(&servicioModel{}).Where("estado = ?", "proximo_a_vencer").Count(&proximos)
	return &domservicio.KPIs{
		CostoTotalMensual: total,
		CantidadActivos:   count,
		ProximosAVencer:   proximos,
		PorTipo:           porTipo,
	}, nil
}

func toServicioDomain(m *servicioModel) (*domservicio.Servicio, error) {
	s := &domservicio.Servicio{
		ID: m.ID, Nombre: m.Nombre, Tipo: domservicio.Tipo(m.Tipo),
		Proveedor: m.Proveedor, CostoMensual: m.CostoMensual,
		VtoFactura: m.VtoFactura, Renovacion: m.Renovacion,
		Estado: domservicio.Estado(m.Estado), CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
	if m.Metadata != nil {
		var meta map[string]interface{}
		if err := json.Unmarshal([]byte(*m.Metadata), &meta); err != nil {
			return nil, err
		}
		s.Metadata = meta
	}
	return s, nil
}

func toServicioModel(s *domservicio.Servicio) (*servicioModel, error) {
	m := &servicioModel{
		ID: s.ID, Nombre: s.Nombre, Tipo: string(s.Tipo),
		Proveedor: s.Proveedor, CostoMensual: s.CostoMensual,
		VtoFactura: s.VtoFactura, Renovacion: s.Renovacion,
		Estado: string(s.Estado), CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
	if s.Metadata != nil {
		b, err := json.Marshal(s.Metadata)
		if err != nil {
			return nil, err
		}
		str := string(b)
		m.Metadata = &str
	}
	return m, nil
}
