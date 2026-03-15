package dashboard

import (
	"context"

	"github.com/ipnext/admin-backend/internal/domain/transferencia"
)

type PagoUrgente struct {
	ID           string  `json:"id"`
	Beneficiario string  `json:"beneficiario"`
	Monto        float64 `json:"monto"`
	Moneda       string  `json:"moneda"`
	FechaPago    string  `json:"fechaPago"`
	Estado       string  `json:"estado"`
	Categoria    string  `json:"categoria"`
}

type GetPagosUrgentesUseCase struct {
	transRepo transferencia.Repository
}

func NewGetPagosUrgentesUseCase(transRepo transferencia.Repository) *GetPagosUrgentesUseCase {
	return &GetPagosUrgentesUseCase{transRepo: transRepo}
}

func (uc *GetPagosUrgentesUseCase) Execute(ctx context.Context) ([]*PagoUrgente, error) {
	items, err := uc.transRepo.FindProximasAVencer(ctx, 7)
	if err != nil {
		return nil, err
	}

	result := make([]*PagoUrgente, 0, len(items))
	for _, t := range items {
		result = append(result, &PagoUrgente{
			ID:           t.ID,
			Beneficiario: t.Beneficiario,
			Monto:        t.Monto,
			Moneda:       t.Moneda,
			FechaPago:    t.FechaPago.Format("2006-01-02"),
			Estado:       string(t.Estado),
			Categoria:    t.Categoria,
		})
	}
	return result, nil
}
