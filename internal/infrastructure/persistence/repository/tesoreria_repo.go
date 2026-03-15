package repository

import (
	"context"
	"errors"
	"time"

	domtesoreria "github.com/ipnext/admin-backend/internal/domain/tesoreria"
	"gorm.io/gorm"
)

type cuentaBancariaModel struct {
	ID                  string    `gorm:"primaryKey;type:char(36)"`
	Banco               string    `gorm:"not null"`
	TipoCuenta          string    `gorm:"type:varchar(30);not null"`
	NroCuenta           string    `gorm:"not null"`
	CBU                 *string
	CCI                 *string
	SaldoActual         float64   `gorm:"type:decimal(15,2);not null;default:0"`
	Moneda              string    `gorm:"type:varchar(10);not null;default:ARS"`
	Activa              bool      `gorm:"not null;default:true"`
	UltimaActualizacion time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (cuentaBancariaModel) TableName() string { return "cuentas_bancarias" }

type movimientoBancarioModel struct {
	ID          string    `gorm:"primaryKey;type:char(36)"`
	CuentaID    string    `gorm:"type:char(36);not null"`
	Tipo        string    `gorm:"type:varchar(20);not null"`
	Monto       float64   `gorm:"type:decimal(15,2);not null"`
	Descripcion string    `gorm:"not null"`
	Fecha       time.Time `gorm:"not null"`
	Conciliado  bool      `gorm:"not null;default:false"`
	Referencia  *string
	CreatedAt   time.Time
}

func (movimientoBancarioModel) TableName() string { return "movimientos_bancarios" }

// --- CuentaBancaria ---

type MySQLCuentaRepository struct{ db *gorm.DB }

func NewMySQLCuentaRepository(db *gorm.DB) domtesoreria.CuentaRepository {
	return &MySQLCuentaRepository{db: db}
}

func (r *MySQLCuentaRepository) FindByID(ctx context.Context, id string) (*domtesoreria.CuentaBancaria, error) {
	var m cuentaBancariaModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toCuentaDomain(&m), nil
}

func (r *MySQLCuentaRepository) FindAll(ctx context.Context) ([]*domtesoreria.CuentaBancaria, error) {
	var models []cuentaBancariaModel
	if err := r.db.WithContext(ctx).Order("banco asc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*domtesoreria.CuentaBancaria, 0, len(models))
	for _, m := range models {
		result = append(result, toCuentaDomain(&m))
	}
	return result, nil
}

func (r *MySQLCuentaRepository) Save(ctx context.Context, c *domtesoreria.CuentaBancaria) error {
	return r.db.WithContext(ctx).Create(toCuentaModel(c)).Error
}

func (r *MySQLCuentaRepository) Update(ctx context.Context, c *domtesoreria.CuentaBancaria) error {
	return r.db.WithContext(ctx).Save(toCuentaModel(c)).Error
}

// --- MovimientoBancario ---

type MySQLMovimientoRepository struct{ db *gorm.DB }

func NewMySQLMovimientoRepository(db *gorm.DB) domtesoreria.MovimientoRepository {
	return &MySQLMovimientoRepository{db: db}
}

func (r *MySQLMovimientoRepository) FindAll(ctx context.Context, cuentaID *string, soloNoConciliados bool) ([]*domtesoreria.MovimientoBancario, error) {
	q := r.db.WithContext(ctx)
	if cuentaID != nil {
		q = q.Where("cuenta_id = ?", *cuentaID)
	}
	if soloNoConciliados {
		q = q.Where("conciliado = false")
	}
	var models []movimientoBancarioModel
	if err := q.Order("fecha desc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*domtesoreria.MovimientoBancario, 0, len(models))
	for _, m := range models {
		result = append(result, toMovimientoDomain(&m))
	}
	return result, nil
}

func (r *MySQLMovimientoRepository) FindByRango(ctx context.Context, desde, hasta time.Time) ([]*domtesoreria.MovimientoBancario, error) {
	var models []movimientoBancarioModel
	err := r.db.WithContext(ctx).
		Where("fecha >= ? AND fecha <= ?", desde, hasta).
		Order("fecha asc").Find(&models).Error
	if err != nil {
		return nil, err
	}
	result := make([]*domtesoreria.MovimientoBancario, 0, len(models))
	for _, m := range models {
		result = append(result, toMovimientoDomain(&m))
	}
	return result, nil
}

func (r *MySQLMovimientoRepository) Save(ctx context.Context, m *domtesoreria.MovimientoBancario) error {
	return r.db.WithContext(ctx).Create(toMovimientoModel(m)).Error
}

func (r *MySQLMovimientoRepository) Conciliar(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&movimientoBancarioModel{}).Where("id = ?", id).Update("conciliado", true).Error
}

// --- Mappers ---

func toCuentaDomain(m *cuentaBancariaModel) *domtesoreria.CuentaBancaria {
	return &domtesoreria.CuentaBancaria{
		ID: m.ID, Banco: m.Banco, TipoCuenta: m.TipoCuenta, NroCuenta: m.NroCuenta,
		CBU: m.CBU, CCI: m.CCI, SaldoActual: m.SaldoActual, Moneda: m.Moneda,
		Activa: m.Activa, UltimaActualizacion: m.UltimaActualizacion,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func toCuentaModel(c *domtesoreria.CuentaBancaria) *cuentaBancariaModel {
	return &cuentaBancariaModel{
		ID: c.ID, Banco: c.Banco, TipoCuenta: c.TipoCuenta, NroCuenta: c.NroCuenta,
		CBU: c.CBU, CCI: c.CCI, SaldoActual: c.SaldoActual, Moneda: c.Moneda,
		Activa: c.Activa, UltimaActualizacion: c.UltimaActualizacion,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
}

func toMovimientoDomain(m *movimientoBancarioModel) *domtesoreria.MovimientoBancario {
	return &domtesoreria.MovimientoBancario{
		ID: m.ID, CuentaID: m.CuentaID, Tipo: domtesoreria.TipoMovimiento(m.Tipo),
		Monto: m.Monto, Descripcion: m.Descripcion, Fecha: m.Fecha,
		Conciliado: m.Conciliado, Referencia: m.Referencia, CreatedAt: m.CreatedAt,
	}
}

func toMovimientoModel(m *domtesoreria.MovimientoBancario) *movimientoBancarioModel {
	return &movimientoBancarioModel{
		ID: m.ID, CuentaID: m.CuentaID, Tipo: string(m.Tipo),
		Monto: m.Monto, Descripcion: m.Descripcion, Fecha: m.Fecha,
		Conciliado: m.Conciliado, Referencia: m.Referencia, CreatedAt: m.CreatedAt,
	}
}
