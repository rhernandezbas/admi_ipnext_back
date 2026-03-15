package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ipnext/admin-backend/internal/domain/empleado"
	"gorm.io/gorm"
)

// --- Models ---

type empleadoModel struct {
	ID           string    `gorm:"primaryKey;type:char(36)"`
	Nombre       string    `gorm:"not null"`
	Puesto       string    `gorm:"not null"`
	Area         string    `gorm:"not null"`
	Rol          string    `gorm:"not null"`
	SueldoBruto  float64   `gorm:"type:decimal(15,2);not null"`
	ObraSocial   string    `gorm:"not null"`
	Activo       bool      `gorm:"default:true"`
	FechaIngreso time.Time `gorm:"not null"`
	Avatar       *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (empleadoModel) TableName() string { return "empleados" }

type liquidacionModel struct {
	ID          string    `gorm:"primaryKey;type:char(36)"`
	EmpleadoID  string    `gorm:"type:char(36);not null"`
	Periodo     string    `gorm:"type:char(7);not null"`
	SueldoBruto float64   `gorm:"type:decimal(15,2);not null"`
	Deducciones float64   `gorm:"type:decimal(15,2);not null;default:0"`
	NetoAPagar  float64   `gorm:"type:decimal(15,2);not null"`
	Estado      string    `gorm:"type:varchar(20);not null;default:borrador"`
	AprobadoPor *string   `gorm:"type:char(36)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (liquidacionModel) TableName() string { return "liquidaciones" }

type guardiaModel struct {
	ID         string    `gorm:"primaryKey;type:char(36)"`
	EmpleadoID string    `gorm:"type:char(36);not null"`
	Fecha      time.Time `gorm:"not null"`
	Horas      float64   `gorm:"type:decimal(5,2);not null"`
	Monto      float64   `gorm:"type:decimal(15,2);not null"`
	Notas      *string   `gorm:"column:descripcion"`
	CreatedAt  time.Time
}

func (guardiaModel) TableName() string { return "guardias" }

type compensacionModel struct {
	ID          string    `gorm:"primaryKey;type:char(36)"`
	EmpleadoID  string    `gorm:"type:char(36);not null"`
	Tipo        string    `gorm:"type:varchar(20);not null"`
	Monto       float64   `gorm:"type:decimal(15,2);not null"`
	Fecha       time.Time `gorm:"not null"`
	Descripcion *string
	CreatedAt   time.Time
}

func (compensacionModel) TableName() string { return "compensaciones" }

// --- Empleado Repository ---

type MySQLEmpleadoRepository struct{ db *gorm.DB }

func NewMySQLEmpleadoRepository(db *gorm.DB) empleado.Repository {
	return &MySQLEmpleadoRepository{db: db}
}

func (r *MySQLEmpleadoRepository) FindByID(ctx context.Context, id string) (*empleado.Empleado, error) {
	var m empleadoModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toEmpleadoDomain(&m), nil
}

func (r *MySQLEmpleadoRepository) FindAll(ctx context.Context, soloActivos bool) ([]*empleado.Empleado, error) {
	q := r.db.WithContext(ctx)
	if soloActivos {
		q = q.Where("activo = ?", true)
	}
	var models []empleadoModel
	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*empleado.Empleado, 0, len(models))
	for _, m := range models {
		result = append(result, toEmpleadoDomain(&m))
	}
	return result, nil
}

func (r *MySQLEmpleadoRepository) Save(ctx context.Context, e *empleado.Empleado) error {
	return r.db.WithContext(ctx).Create(toEmpleadoModel(e)).Error
}

func (r *MySQLEmpleadoRepository) Update(ctx context.Context, e *empleado.Empleado) error {
	return r.db.WithContext(ctx).Save(toEmpleadoModel(e)).Error
}

func (r *MySQLEmpleadoRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&empleadoModel{}).Where("id = ?", id).Update("activo", false).Error
}

func (r *MySQLEmpleadoRepository) CountActivos(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&empleadoModel{}).Where("activo = ?", true).Count(&count).Error
	return count, err
}

func (r *MySQLEmpleadoRepository) SumSueldos(ctx context.Context) (float64, error) {
	type res struct{ Total float64 }
	var r2 res
	err := r.db.WithContext(ctx).Model(&empleadoModel{}).
		Select("COALESCE(SUM(sueldo_bruto), 0) as total").
		Where("activo = ?", true).Scan(&r2).Error
	return r2.Total, err
}

// --- Liquidacion Repository ---

type MySQLLiquidacionRepository struct{ db *gorm.DB }

func NewMySQLLiquidacionRepository(db *gorm.DB) empleado.LiquidacionRepository {
	return &MySQLLiquidacionRepository{db: db}
}

func (r *MySQLLiquidacionRepository) FindByID(ctx context.Context, id string) (*empleado.Liquidacion, error) {
	var m liquidacionModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toLiquidacionDomain(&m), nil
}

func (r *MySQLLiquidacionRepository) FindByPeriodo(ctx context.Context, periodo string) ([]*empleado.Liquidacion, error) {
	var models []liquidacionModel
	if err := r.db.WithContext(ctx).Where("periodo = ?", periodo).Find(&models).Error; err != nil {
		return nil, err
	}
	return toLiquidacionSlice(models), nil
}

func (r *MySQLLiquidacionRepository) FindByEmpleado(ctx context.Context, empleadoID string) ([]*empleado.Liquidacion, error) {
	var models []liquidacionModel
	if err := r.db.WithContext(ctx).Where("empleado_id = ?", empleadoID).Order("periodo desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toLiquidacionSlice(models), nil
}

func (r *MySQLLiquidacionRepository) Save(ctx context.Context, l *empleado.Liquidacion) error {
	return r.db.WithContext(ctx).Create(toLiquidacionModel(l)).Error
}

func (r *MySQLLiquidacionRepository) Update(ctx context.Context, l *empleado.Liquidacion) error {
	return r.db.WithContext(ctx).Save(toLiquidacionModel(l)).Error
}

// --- Guardia Repository ---

type MySQLGuardiaRepository struct{ db *gorm.DB }

func NewMySQLGuardiaRepository(db *gorm.DB) empleado.GuardiaRepository {
	return &MySQLGuardiaRepository{db: db}
}

func (r *MySQLGuardiaRepository) FindAll(ctx context.Context, empleadoID *string) ([]*empleado.Guardia, error) {
	q := r.db.WithContext(ctx)
	if empleadoID != nil {
		q = q.Where("empleado_id = ?", *empleadoID)
	}
	var models []guardiaModel
	if err := q.Order("fecha desc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*empleado.Guardia, 0, len(models))
	for _, m := range models {
		result = append(result, toGuardiaDomain(&m))
	}
	return result, nil
}

func (r *MySQLGuardiaRepository) Save(ctx context.Context, g *empleado.Guardia) error {
	return r.db.WithContext(ctx).Create(toGuardiaModel(g)).Error
}

// --- Compensacion Repository ---

type MySQLCompensacionRepository struct{ db *gorm.DB }

func NewMySQLCompensacionRepository(db *gorm.DB) empleado.CompensacionRepository {
	return &MySQLCompensacionRepository{db: db}
}

func (r *MySQLCompensacionRepository) FindAll(ctx context.Context, empleadoID *string) ([]*empleado.Compensacion, error) {
	q := r.db.WithContext(ctx)
	if empleadoID != nil {
		q = q.Where("empleado_id = ?", *empleadoID)
	}
	var models []compensacionModel
	if err := q.Order("fecha desc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*empleado.Compensacion, 0, len(models))
	for _, m := range models {
		result = append(result, toCompensacionDomain(&m))
	}
	return result, nil
}

func (r *MySQLCompensacionRepository) Save(ctx context.Context, c *empleado.Compensacion) error {
	return r.db.WithContext(ctx).Create(toCompensacionModel(c)).Error
}

// --- Mappers ---

func toEmpleadoDomain(m *empleadoModel) *empleado.Empleado {
	return &empleado.Empleado{ID: m.ID, Nombre: m.Nombre, Puesto: m.Puesto, Area: m.Area, Rol: m.Rol,
		SueldoBruto: m.SueldoBruto, ObraSocial: m.ObraSocial, Activo: m.Activo,
		FechaIngreso: m.FechaIngreso, Avatar: m.Avatar, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}
func toEmpleadoModel(e *empleado.Empleado) *empleadoModel {
	return &empleadoModel{ID: e.ID, Nombre: e.Nombre, Puesto: e.Puesto, Area: e.Area, Rol: e.Rol,
		SueldoBruto: e.SueldoBruto, ObraSocial: e.ObraSocial, Activo: e.Activo,
		FechaIngreso: e.FechaIngreso, Avatar: e.Avatar, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt}
}
func toLiquidacionDomain(m *liquidacionModel) *empleado.Liquidacion {
	return &empleado.Liquidacion{ID: m.ID, EmpleadoID: m.EmpleadoID, Periodo: m.Periodo,
		SueldoBruto: m.SueldoBruto, Deducciones: m.Deducciones, NetoAPagar: m.NetoAPagar,
		Estado: empleado.EstadoLiquidacion(m.Estado), AprobadoPor: m.AprobadoPor,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}
func toLiquidacionModel(l *empleado.Liquidacion) *liquidacionModel {
	return &liquidacionModel{ID: l.ID, EmpleadoID: l.EmpleadoID, Periodo: l.Periodo,
		SueldoBruto: l.SueldoBruto, Deducciones: l.Deducciones, NetoAPagar: l.NetoAPagar,
		Estado: string(l.Estado), AprobadoPor: l.AprobadoPor,
		CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt}
}
func toLiquidacionSlice(models []liquidacionModel) []*empleado.Liquidacion {
	r := make([]*empleado.Liquidacion, 0, len(models))
	for _, m := range models {
		r = append(r, toLiquidacionDomain(&m))
	}
	return r
}
func toGuardiaDomain(m *guardiaModel) *empleado.Guardia {
	return &empleado.Guardia{ID: m.ID, EmpleadoID: m.EmpleadoID, Fecha: m.Fecha,
		Horas: m.Horas, Monto: m.Monto, Notas: m.Notas, CreatedAt: m.CreatedAt}
}
func toGuardiaModel(g *empleado.Guardia) *guardiaModel {
	return &guardiaModel{ID: g.ID, EmpleadoID: g.EmpleadoID, Fecha: g.Fecha,
		Horas: g.Horas, Monto: g.Monto, Notas: g.Notas, CreatedAt: g.CreatedAt}
}
func toCompensacionDomain(m *compensacionModel) *empleado.Compensacion {
	return &empleado.Compensacion{ID: m.ID, EmpleadoID: m.EmpleadoID, Tipo: empleado.TipoCompensacion(m.Tipo),
		Monto: m.Monto, Fecha: m.Fecha, Descripcion: m.Descripcion, CreatedAt: m.CreatedAt}
}
func toCompensacionModel(c *empleado.Compensacion) *compensacionModel {
	return &compensacionModel{ID: c.ID, EmpleadoID: c.EmpleadoID, Tipo: string(c.Tipo),
		Monto: c.Monto, Fecha: c.Fecha, Descripcion: c.Descripcion, CreatedAt: c.CreatedAt}
}
