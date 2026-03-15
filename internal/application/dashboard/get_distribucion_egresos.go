package dashboard

import (
	"context"
	"time"

	"github.com/ipnext/admin-backend/internal/domain/transferencia"
)

type DistribucionItem struct {
	Categoria  string  `json:"categoria"`
	Monto      float64 `json:"monto"`
	Porcentaje float64 `json:"porcentaje"`
}

type GetDistribucionEgresosUseCase struct {
	transRepo transferencia.Repository
}

func NewGetDistribucionEgresosUseCase(transRepo transferencia.Repository) *GetDistribucionEgresosUseCase {
	return &GetDistribucionEgresosUseCase{transRepo: transRepo}
}

func (uc *GetDistribucionEgresosUseCase) Execute(ctx context.Context) ([]*DistribucionItem, error) {
	now := time.Now()
	inicioMes := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	finMes := inicioMes.AddDate(0, 1, 0)

	porCategoria, err := uc.transRepo.SumByCategoria(ctx, inicioMes, finMes)
	if err != nil {
		return nil, err
	}

	var total float64
	for _, monto := range porCategoria {
		total += monto
	}

	result := make([]*DistribucionItem, 0, len(porCategoria))
	for cat, monto := range porCategoria {
		porcentaje := 0.0
		if total > 0 {
			porcentaje = (monto / total) * 100
		}
		result = append(result, &DistribucionItem{
			Categoria:  cat,
			Monto:      monto,
			Porcentaje: porcentaje,
		})
	}
	return result, nil
}
