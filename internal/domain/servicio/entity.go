package servicio

import "time"

type Tipo string

const (
	TipoInternet  Tipo = "internet"
	TipoEnergia   Tipo = "energia"
	TipoSeguridad Tipo = "seguridad"
	TipoSoftware  Tipo = "software"
	TipoObraSocial Tipo = "obra_social"
	TipoSeguro    Tipo = "seguro"
)

type Estado string

const (
	EstadoActivo         Estado = "activo"
	EstadoProximoAVencer Estado = "proximo_a_vencer"
	EstadoInactivo       Estado = "inactivo"
)

type Servicio struct {
	ID           string
	Nombre       string
	Tipo         Tipo
	Proveedor    string
	CostoMensual float64
	VtoFactura   *time.Time
	Renovacion   *time.Time
	Estado       Estado
	Metadata     map[string]interface{}
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (s *Servicio) DiasParaRenovar() *int {
	if s.Renovacion == nil {
		return nil
	}
	dias := int(time.Until(*s.Renovacion).Hours() / 24)
	return &dias
}

type KPIs struct {
	CostoTotalMensual float64
	CantidadActivos   int64
	ProximosAVencer   int64
	PorTipo           map[Tipo]float64
}
