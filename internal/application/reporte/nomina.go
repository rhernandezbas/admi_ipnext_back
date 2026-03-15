package reporte

import (
	"context"

	domemp "github.com/ipnext/admin-backend/internal/domain/empleado"
)

type ReporteNominaItem struct {
	EmpleadoID    string
	Nombre        string
	Puesto        string
	MontoTotal    float64
	Liquidaciones int
}

type ReporteNominaUseCase struct {
	liquidacionRepo domemp.LiquidacionRepository
	empleadoRepo    domemp.Repository
}

func NewReporteNominaUseCase(liquidacionRepo domemp.LiquidacionRepository, empleadoRepo domemp.Repository) *ReporteNominaUseCase {
	return &ReporteNominaUseCase{liquidacionRepo: liquidacionRepo, empleadoRepo: empleadoRepo}
}

func (uc *ReporteNominaUseCase) Execute(ctx context.Context, periodo string) ([]*ReporteNominaItem, error) {
	liqs, err := uc.liquidacionRepo.FindByPeriodo(ctx, periodo)
	if err != nil {
		return nil, err
	}
	byEmp := make(map[string]*ReporteNominaItem)
	for _, l := range liqs {
		if _, ok := byEmp[l.EmpleadoID]; !ok {
			byEmp[l.EmpleadoID] = &ReporteNominaItem{EmpleadoID: l.EmpleadoID}
		}
		byEmp[l.EmpleadoID].MontoTotal += l.NetoAPagar
		byEmp[l.EmpleadoID].Liquidaciones++
	}
	for empID, item := range byEmp {
		emp, err := uc.empleadoRepo.FindByID(ctx, empID)
		if err == nil && emp != nil {
			item.Nombre = emp.Nombre
			item.Puesto = emp.Puesto
		}
	}
	result := make([]*ReporteNominaItem, 0, len(byEmp))
	for _, item := range byEmp {
		result = append(result, item)
	}
	return result, nil
}
