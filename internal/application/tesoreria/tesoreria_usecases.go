package tesoreria

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	domtesoreria "github.com/ipnext/admin-backend/internal/domain/tesoreria"
)

var ErrCuentaNoEncontrada = errors.New("cuenta no encontrada")

// --- Flujo de caja ---

type GetFlujoCajaUseCase struct{ repo domtesoreria.MovimientoRepository }

func NewGetFlujoCajaUseCase(repo domtesoreria.MovimientoRepository) *GetFlujoCajaUseCase {
	return &GetFlujoCajaUseCase{repo: repo}
}

func (uc *GetFlujoCajaUseCase) Execute(ctx context.Context, desde, hasta time.Time) ([]*domtesoreria.FlujoCajaItem, error) {
	movs, err := uc.repo.FindByRango(ctx, desde, hasta)
	if err != nil {
		return nil, err
	}
	byDay := make(map[string]*domtesoreria.FlujoCajaItem)
	for _, m := range movs {
		key := m.Fecha.Format("2006-01-02")
		if _, ok := byDay[key]; !ok {
			byDay[key] = &domtesoreria.FlujoCajaItem{Fecha: m.Fecha}
		}
		if m.Tipo == domtesoreria.TipoIngreso {
			byDay[key].Ingresos += m.Monto
		} else {
			byDay[key].Egresos += m.Monto
		}
	}
	result := make([]*domtesoreria.FlujoCajaItem, 0, len(byDay))
	for _, item := range byDay {
		item.Saldo = item.Ingresos - item.Egresos
		result = append(result, item)
	}
	return result, nil
}

// --- Proyecciones ---

type GetProyeccionesUseCase struct {
	movRepo    domtesoreria.MovimientoRepository
	cuentaRepo domtesoreria.CuentaRepository
}

func NewGetProyeccionesUseCase(movRepo domtesoreria.MovimientoRepository, cuentaRepo domtesoreria.CuentaRepository) *GetProyeccionesUseCase {
	return &GetProyeccionesUseCase{movRepo: movRepo, cuentaRepo: cuentaRepo}
}

func (uc *GetProyeccionesUseCase) Execute(ctx context.Context, meses int) ([]*domtesoreria.ProyeccionItem, error) {
	// Calcular egreso promedio de los últimos 3 meses
	hasta := time.Now()
	desde := hasta.AddDate(0, -3, 0)
	movs, err := uc.movRepo.FindByRango(ctx, desde, hasta)
	if err != nil {
		return nil, err
	}
	var totalEgresos float64
	for _, m := range movs {
		if m.Tipo == domtesoreria.TipoEgreso {
			totalEgresos += m.Monto
		}
	}
	egresoPromedio := totalEgresos / 3

	cuentas, err := uc.cuentaRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var saldoTotal float64
	for _, c := range cuentas {
		if c.Activa {
			saldoTotal += c.SaldoActual
		}
	}

	result := make([]*domtesoreria.ProyeccionItem, 0, meses)
	saldo := saldoTotal
	for i := 1; i <= meses; i++ {
		mes := time.Now().AddDate(0, i, 0)
		saldo -= egresoPromedio
		result = append(result, &domtesoreria.ProyeccionItem{
			Mes:              fmt.Sprintf("%d-%02d", mes.Year(), mes.Month()),
			EgresosPrevistos: egresoPromedio,
			SaldoProyectado:  saldo,
		})
	}
	return result, nil
}

// --- Cuentas ---

type ListCuentasUseCase struct{ repo domtesoreria.CuentaRepository }

func NewListCuentasUseCase(repo domtesoreria.CuentaRepository) *ListCuentasUseCase {
	return &ListCuentasUseCase{repo: repo}
}
func (uc *ListCuentasUseCase) Execute(ctx context.Context) ([]*domtesoreria.CuentaBancaria, error) {
	return uc.repo.FindAll(ctx)
}

type CreateCuentaRequest struct {
	Banco      string
	TipoCuenta string
	NroCuenta  string
	CBU        *string
	CCI        *string
	Moneda     string
}

type CreateCuentaUseCase struct{ repo domtesoreria.CuentaRepository }

func NewCreateCuentaUseCase(repo domtesoreria.CuentaRepository) *CreateCuentaUseCase {
	return &CreateCuentaUseCase{repo: repo}
}
func (uc *CreateCuentaUseCase) Execute(ctx context.Context, req CreateCuentaRequest) (*domtesoreria.CuentaBancaria, error) {
	moneda := req.Moneda
	if moneda == "" {
		moneda = "ARS"
	}
	c := &domtesoreria.CuentaBancaria{
		ID: uuid.NewString(), Banco: req.Banco, TipoCuenta: req.TipoCuenta,
		NroCuenta: req.NroCuenta, CBU: req.CBU, CCI: req.CCI,
		SaldoActual: 0, Moneda: moneda, Activa: true,
		UltimaActualizacion: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := uc.repo.Save(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

type UpdateCuentaRequest struct {
	Banco      *string
	TipoCuenta *string
	NroCuenta  *string
	CBU        *string
	CCI        *string
	SaldoActual *float64
	Activa     *bool
}

type UpdateCuentaUseCase struct{ repo domtesoreria.CuentaRepository }

func NewUpdateCuentaUseCase(repo domtesoreria.CuentaRepository) *UpdateCuentaUseCase {
	return &UpdateCuentaUseCase{repo: repo}
}
func (uc *UpdateCuentaUseCase) Execute(ctx context.Context, id string, req UpdateCuentaRequest) (*domtesoreria.CuentaBancaria, error) {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCuentaNoEncontrada
	}
	if req.Banco != nil {
		c.Banco = *req.Banco
	}
	if req.TipoCuenta != nil {
		c.TipoCuenta = *req.TipoCuenta
	}
	if req.NroCuenta != nil {
		c.NroCuenta = *req.NroCuenta
	}
	if req.CBU != nil {
		c.CBU = req.CBU
	}
	if req.CCI != nil {
		c.CCI = req.CCI
	}
	if req.SaldoActual != nil {
		c.SaldoActual = *req.SaldoActual
	}
	if req.Activa != nil {
		c.Activa = *req.Activa
	}
	c.UltimaActualizacion = time.Now()
	c.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// --- Movimientos ---

type GetConciliacionUseCase struct{ repo domtesoreria.MovimientoRepository }

func NewGetConciliacionUseCase(repo domtesoreria.MovimientoRepository) *GetConciliacionUseCase {
	return &GetConciliacionUseCase{repo: repo}
}
func (uc *GetConciliacionUseCase) Execute(ctx context.Context, cuentaID *string) ([]*domtesoreria.MovimientoBancario, error) {
	return uc.repo.FindAll(ctx, cuentaID, true)
}

type CreateMovimientoRequest struct {
	CuentaID    string
	Tipo        domtesoreria.TipoMovimiento
	Monto       float64
	Descripcion string
	Fecha       time.Time
	Referencia  *string
}

type CreateMovimientoUseCase struct{ repo domtesoreria.MovimientoRepository }

func NewCreateMovimientoUseCase(repo domtesoreria.MovimientoRepository) *CreateMovimientoUseCase {
	return &CreateMovimientoUseCase{repo: repo}
}
func (uc *CreateMovimientoUseCase) Execute(ctx context.Context, req CreateMovimientoRequest) (*domtesoreria.MovimientoBancario, error) {
	m := &domtesoreria.MovimientoBancario{
		ID: uuid.NewString(), CuentaID: req.CuentaID, Tipo: req.Tipo,
		Monto: req.Monto, Descripcion: req.Descripcion, Fecha: req.Fecha,
		Conciliado: false, Referencia: req.Referencia, CreatedAt: time.Now(),
	}
	if err := uc.repo.Save(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

type ConciliarMovimientoUseCase struct{ repo domtesoreria.MovimientoRepository }

func NewConciliarMovimientoUseCase(repo domtesoreria.MovimientoRepository) *ConciliarMovimientoUseCase {
	return &ConciliarMovimientoUseCase{repo: repo}
}
func (uc *ConciliarMovimientoUseCase) Execute(ctx context.Context, id string) error {
	return uc.repo.Conciliar(ctx, id)
}
