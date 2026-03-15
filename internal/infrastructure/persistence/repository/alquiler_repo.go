package repository

import (
	"context"
	"errors"
	"time"

	domalquiler "github.com/ipnext/admin-backend/internal/domain/alquiler"
	"gorm.io/gorm"
)

type inmuebleModel struct {
	ID              string    `gorm:"primaryKey;type:char(36)"`
	Nombre          string    `gorm:"not null"`
	Direccion       string    `gorm:"not null"`
	Propietario     string    `gorm:"not null"`
	Uso             string    `gorm:"type:varchar(20);not null"`
	AlquilerMensual float64   `gorm:"type:decimal(15,2);not null"`
	CBU             *string
	Alias           *string
	Estado          string    `gorm:"type:varchar(20);not null;default:pendiente"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (inmuebleModel) TableName() string { return "inmuebles" }

type contratoAlquilerModel struct {
	ID               string    `gorm:"primaryKey;type:char(36)"`
	InmuebleID       string    `gorm:"type:char(36);not null"`
	VigenciaDesde    time.Time `gorm:"not null"`
	VigenciaHasta    time.Time `gorm:"not null"`
	AjusteFrecuencia string    `gorm:"not null"`
	MontoMensual     float64   `gorm:"type:decimal(15,2);not null"`
	Estado           string    `gorm:"type:varchar(20);not null;default:vigente"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (contratoAlquilerModel) TableName() string { return "contratos_alquiler" }

type pagoAlquilerModel struct {
	ID         string     `gorm:"primaryKey;type:char(36)"`
	InmuebleID string     `gorm:"type:char(36);not null"`
	Periodo    string     `gorm:"type:char(7);not null"`
	Monto      float64    `gorm:"type:decimal(15,2);not null"`
	FechaPago  *time.Time
	Estado     string     `gorm:"type:varchar(20);not null;default:pendiente"`
	Comprobante *string
	CreatedAt  time.Time
}

func (pagoAlquilerModel) TableName() string { return "pagos_alquiler" }

// --- Inmueble ---

type MySQLAlquilerRepository struct{ db *gorm.DB }

func NewMySQLAlquilerRepository(db *gorm.DB) domalquiler.Repository {
	return &MySQLAlquilerRepository{db: db}
}

func (r *MySQLAlquilerRepository) FindByID(ctx context.Context, id string) (*domalquiler.Inmueble, error) {
	var m inmuebleModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toInmuebleDomain(&m), nil
}

func (r *MySQLAlquilerRepository) FindAll(ctx context.Context) ([]*domalquiler.Inmueble, error) {
	var models []inmuebleModel
	if err := r.db.WithContext(ctx).Order("nombre asc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*domalquiler.Inmueble, 0, len(models))
	for _, m := range models {
		result = append(result, toInmuebleDomain(&m))
	}
	return result, nil
}

func (r *MySQLAlquilerRepository) Save(ctx context.Context, i *domalquiler.Inmueble) error {
	return r.db.WithContext(ctx).Create(toInmuebleModel(i)).Error
}

func (r *MySQLAlquilerRepository) Update(ctx context.Context, i *domalquiler.Inmueble) error {
	return r.db.WithContext(ctx).Save(toInmuebleModel(i)).Error
}

func (r *MySQLAlquilerRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&inmuebleModel{}, "id = ?", id).Error
}

// --- Contrato ---

type MySQLContratoAlquilerRepository struct{ db *gorm.DB }

func NewMySQLContratoAlquilerRepository(db *gorm.DB) domalquiler.ContratoRepository {
	return &MySQLContratoAlquilerRepository{db: db}
}

func (r *MySQLContratoAlquilerRepository) FindByID(ctx context.Context, id string) (*domalquiler.ContratoAlquiler, error) {
	var m contratoAlquilerModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toContratoAlquilerDomain(&m), nil
}

func (r *MySQLContratoAlquilerRepository) FindAll(ctx context.Context, inmuebleID *string) ([]*domalquiler.ContratoAlquiler, error) {
	q := r.db.WithContext(ctx)
	if inmuebleID != nil {
		q = q.Where("inmueble_id = ?", *inmuebleID)
	}
	var models []contratoAlquilerModel
	if err := q.Order("vigencia_hasta desc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*domalquiler.ContratoAlquiler, 0, len(models))
	for _, m := range models {
		result = append(result, toContratoAlquilerDomain(&m))
	}
	return result, nil
}

func (r *MySQLContratoAlquilerRepository) FindProximosAVencer(ctx context.Context, dias int) ([]*domalquiler.ContratoAlquiler, error) {
	hasta := time.Now().AddDate(0, 0, dias)
	var models []contratoAlquilerModel
	err := r.db.WithContext(ctx).
		Where("estado = ? AND vigencia_hasta <= ? AND vigencia_hasta >= ?", "vigente", hasta, time.Now()).
		Order("vigencia_hasta asc").Find(&models).Error
	if err != nil {
		return nil, err
	}
	result := make([]*domalquiler.ContratoAlquiler, 0, len(models))
	for _, m := range models {
		result = append(result, toContratoAlquilerDomain(&m))
	}
	return result, nil
}

func (r *MySQLContratoAlquilerRepository) Save(ctx context.Context, c *domalquiler.ContratoAlquiler) error {
	return r.db.WithContext(ctx).Create(toContratoAlquilerModel(c)).Error
}

func (r *MySQLContratoAlquilerRepository) Update(ctx context.Context, c *domalquiler.ContratoAlquiler) error {
	return r.db.WithContext(ctx).Save(toContratoAlquilerModel(c)).Error
}

// --- Pago ---

type MySQLPagoAlquilerRepository struct{ db *gorm.DB }

func NewMySQLPagoAlquilerRepository(db *gorm.DB) domalquiler.PagoRepository {
	return &MySQLPagoAlquilerRepository{db: db}
}

func (r *MySQLPagoAlquilerRepository) FindAll(ctx context.Context, inmuebleID *string) ([]*domalquiler.PagoAlquiler, error) {
	q := r.db.WithContext(ctx)
	if inmuebleID != nil {
		q = q.Where("inmueble_id = ?", *inmuebleID)
	}
	var models []pagoAlquilerModel
	if err := q.Order("periodo desc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*domalquiler.PagoAlquiler, 0, len(models))
	for _, m := range models {
		result = append(result, toPagoAlquilerDomain(&m))
	}
	return result, nil
}

func (r *MySQLPagoAlquilerRepository) Save(ctx context.Context, p *domalquiler.PagoAlquiler) error {
	return r.db.WithContext(ctx).Create(toPagoAlquilerModel(p)).Error
}

// --- Mappers ---

func toInmuebleDomain(m *inmuebleModel) *domalquiler.Inmueble {
	return &domalquiler.Inmueble{ID: m.ID, Nombre: m.Nombre, Direccion: m.Direccion,
		Propietario: m.Propietario, Uso: domalquiler.UsoInmueble(m.Uso),
		AlquilerMensual: m.AlquilerMensual, CBU: m.CBU, Alias: m.Alias,
		Estado: domalquiler.EstadoInmueble(m.Estado), CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}
func toInmuebleModel(i *domalquiler.Inmueble) *inmuebleModel {
	return &inmuebleModel{ID: i.ID, Nombre: i.Nombre, Direccion: i.Direccion,
		Propietario: i.Propietario, Uso: string(i.Uso),
		AlquilerMensual: i.AlquilerMensual, CBU: i.CBU, Alias: i.Alias,
		Estado: string(i.Estado), CreatedAt: i.CreatedAt, UpdatedAt: i.UpdatedAt}
}
func toContratoAlquilerDomain(m *contratoAlquilerModel) *domalquiler.ContratoAlquiler {
	return &domalquiler.ContratoAlquiler{ID: m.ID, InmuebleID: m.InmuebleID,
		VigenciaDesde: m.VigenciaDesde, VigenciaHasta: m.VigenciaHasta,
		AjusteFrecuencia: m.AjusteFrecuencia, MontoMensual: m.MontoMensual,
		Estado: domalquiler.EstadoContrato(m.Estado), CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}
func toContratoAlquilerModel(c *domalquiler.ContratoAlquiler) *contratoAlquilerModel {
	return &contratoAlquilerModel{ID: c.ID, InmuebleID: c.InmuebleID,
		VigenciaDesde: c.VigenciaDesde, VigenciaHasta: c.VigenciaHasta,
		AjusteFrecuencia: c.AjusteFrecuencia, MontoMensual: c.MontoMensual,
		Estado: string(c.Estado), CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
}
func toPagoAlquilerDomain(m *pagoAlquilerModel) *domalquiler.PagoAlquiler {
	return &domalquiler.PagoAlquiler{ID: m.ID, InmuebleID: m.InmuebleID, Periodo: m.Periodo,
		Monto: m.Monto, FechaPago: m.FechaPago, Estado: domalquiler.EstadoInmueble(m.Estado),
		Comprobante: m.Comprobante, CreatedAt: m.CreatedAt}
}
func toPagoAlquilerModel(p *domalquiler.PagoAlquiler) *pagoAlquilerModel {
	return &pagoAlquilerModel{ID: p.ID, InmuebleID: p.InmuebleID, Periodo: p.Periodo,
		Monto: p.Monto, FechaPago: p.FechaPago, Estado: string(p.Estado),
		Comprobante: p.Comprobante, CreatedAt: p.CreatedAt}
}
