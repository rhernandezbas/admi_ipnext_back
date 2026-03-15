package reporte

import (
	"context"

	domalq "github.com/ipnext/admin-backend/internal/domain/alquiler"
)

type ReporteInmueblesItem struct {
	InmuebleID  string
	Nombre      string
	Propietario string
	Uso         string
	Estado      string
	AlquilerMensual float64
	TotalPagado float64
	PagosPendientes int
}

type ReporteInmueblesUseCase struct {
	inmuebleRepo domalq.Repository
	pagoRepo     domalq.PagoRepository
}

func NewReporteInmueblesUseCase(inmuebleRepo domalq.Repository, pagoRepo domalq.PagoRepository) *ReporteInmueblesUseCase {
	return &ReporteInmueblesUseCase{inmuebleRepo: inmuebleRepo, pagoRepo: pagoRepo}
}

func (uc *ReporteInmueblesUseCase) Execute(ctx context.Context) ([]*ReporteInmueblesItem, error) {
	inmuebles, err := uc.inmuebleRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*ReporteInmueblesItem, 0, len(inmuebles))
	for _, i := range inmuebles {
		item := &ReporteInmueblesItem{
			InmuebleID:      i.ID,
			Nombre:          i.Nombre,
			Propietario:     i.Propietario,
			Uso:             string(i.Uso),
			Estado:          string(i.Estado),
			AlquilerMensual: i.AlquilerMensual,
		}
		pagos, err := uc.pagoRepo.FindAll(ctx, &i.ID)
		if err == nil {
			for _, p := range pagos {
				if p.Estado == domalq.EstadoPagado {
					item.TotalPagado += p.Monto
				} else if p.Estado == domalq.EstadoPendiente {
					item.PagosPendientes++
				}
			}
		}
		result = append(result, item)
	}
	return result, nil
}
