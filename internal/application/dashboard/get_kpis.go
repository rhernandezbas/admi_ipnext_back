package dashboard

import (
	"context"
	"time"

	"github.com/ipnext/admin-backend/internal/domain/transferencia"
	"github.com/ipnext/admin-backend/internal/domain/empleado"
)

type KPIs struct {
	TotalPagosMes    float64 `json:"totalPagosMes"`
	PagosPendientes  int64   `json:"pagosPendientes"`
	PagosVencidos    int64   `json:"pagosVencidos"`
	FlujoCajaMes     float64 `json:"flujoCajaMes"`
	TotalEmpleados   int64   `json:"totalEmpleados"`
	CostoNominaMes   float64 `json:"costoNominaMes"`
}

type GetKPIsUseCase struct {
	transRepo    transferencia.Repository
	empleadoRepo empleado.Repository
}

func NewGetKPIsUseCase(transRepo transferencia.Repository, empleadoRepo empleado.Repository) *GetKPIsUseCase {
	return &GetKPIsUseCase{transRepo: transRepo, empleadoRepo: empleadoRepo}
}

func (uc *GetKPIsUseCase) Execute(ctx context.Context) (*KPIs, error) {
	now := time.Now()
	inicioMes := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	finMes := inicioMes.AddDate(0, 1, 0)

	// Total pagado en el mes
	pagadas, err := uc.transRepo.FindAll(ctx, transferencia.Filtros{
		Estado:  string(transferencia.EstadoPagado),
		Desde:   &inicioMes,
		Hasta:   &finMes,
		Page:    1,
		PerPage: 9999,
	})
	if err != nil {
		return nil, err
	}
	var totalMes float64
	for _, t := range pagadas.Items {
		totalMes += t.Monto
	}

	// Pendientes
	pendientes, err := uc.transRepo.FindAll(ctx, transferencia.Filtros{
		Estado:  string(transferencia.EstadoPendiente),
		Page:    1,
		PerPage: 9999,
	})
	if err != nil {
		return nil, err
	}

	// Vencidos
	vencidos, err := uc.transRepo.FindAll(ctx, transferencia.Filtros{
		Estado:  string(transferencia.EstadoVencido),
		Page:    1,
		PerPage: 9999,
	})
	if err != nil {
		return nil, err
	}

	// Empleados activos y costo nómina
	totalEmpleados, err := uc.empleadoRepo.CountActivos(ctx)
	if err != nil {
		return nil, err
	}
	costoNomina, err := uc.empleadoRepo.SumSueldos(ctx)
	if err != nil {
		return nil, err
	}

	return &KPIs{
		TotalPagosMes:   totalMes,
		PagosPendientes: pendientes.Total,
		PagosVencidos:   vencidos.Total,
		FlujoCajaMes:    totalMes,
		TotalEmpleados:  totalEmpleados,
		CostoNominaMes:  costoNomina,
	}, nil
}
