package reporte

import (
	"context"
	"time"

	domtrans "github.com/ipnext/admin-backend/internal/domain/transferencia"
)

type ReporteFinancieroItem struct {
	Periodo  string
	Egresos  float64
	Cantidad int
}

type ReporteFinancieroUseCase struct {
	transRepo domtrans.Repository
}

func NewReporteFinancieroUseCase(transRepo domtrans.Repository) *ReporteFinancieroUseCase {
	return &ReporteFinancieroUseCase{transRepo: transRepo}
}

func (uc *ReporteFinancieroUseCase) Execute(ctx context.Context, desde, hasta time.Time) ([]*ReporteFinancieroItem, error) {
	filtros := domtrans.Filtros{
		Desde:   &desde,
		Hasta:   &hasta,
		PerPage: 10000,
		Page:    1,
	}
	result, err := uc.transRepo.FindAll(ctx, filtros)
	if err != nil {
		return nil, err
	}

	byMes := make(map[string]*ReporteFinancieroItem)
	for _, t := range result.Items {
		mes := t.FechaPago.Format("2006-01")
		if _, ok := byMes[mes]; !ok {
			byMes[mes] = &ReporteFinancieroItem{Periodo: mes}
		}
		if t.Estado == domtrans.EstadoPagado {
			byMes[mes].Egresos += t.Monto
			byMes[mes].Cantidad++
		}
	}
	items := make([]*ReporteFinancieroItem, 0, len(byMes))
	for _, item := range byMes {
		items = append(items, item)
	}
	return items, nil
}
