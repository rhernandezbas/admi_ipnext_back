package reporte

import (
	"context"

	domprov "github.com/ipnext/admin-backend/internal/domain/proveedor"
)

type ReporteProveedoresItem struct {
	ProveedorID     string
	NombreProveedor string
	TotalPagado     float64
	CantidadPagos   int64
	Posicion        int
}

type ReporteProveedoresUseCase struct {
	rankingRepo domprov.RankingRepository
}

func NewReporteProveedoresUseCase(rankingRepo domprov.RankingRepository) *ReporteProveedoresUseCase {
	return &ReporteProveedoresUseCase{rankingRepo: rankingRepo}
}

func (uc *ReporteProveedoresUseCase) Execute(ctx context.Context) ([]*ReporteProveedoresItem, error) {
	ranking, err := uc.rankingRepo.GetRanking(ctx, 100)
	if err != nil {
		return nil, err
	}
	result := make([]*ReporteProveedoresItem, 0, len(ranking))
	for i, r := range ranking {
		result = append(result, &ReporteProveedoresItem{
			ProveedorID:     r.ProveedorID,
			NombreProveedor: r.NombreProveedor,
			TotalPagado:     r.TotalPagado,
			CantidadPagos:   r.CantidadPagos,
			Posicion:        i + 1,
		})
	}
	return result, nil
}
