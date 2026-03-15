package alquiler

import "time"

type EstadoInmueble string

const (
	EstadoPagado   EstadoInmueble = "pagado"
	EstadoPendiente EstadoInmueble = "pendiente"
	EstadoVencido  EstadoInmueble = "vencido"
)

type UsoInmueble string

const (
	UsoNodo     UsoInmueble = "nodo"
	UsoOficina  UsoInmueble = "oficina"
	UsoDeposito UsoInmueble = "deposito"
	UsoOtro     UsoInmueble = "otro"
)

type Inmueble struct {
	ID              string
	Nombre          string
	Direccion       string
	Propietario     string
	Uso             UsoInmueble
	AlquilerMensual float64
	CBU             *string
	Alias           *string
	Estado          EstadoInmueble
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type EstadoContrato string

const (
	EstadoVigente    EstadoContrato = "vigente"
	EstadoPorVencer  EstadoContrato = "por_vencer"
	EstadoContratoVencido EstadoContrato = "vencido"
)

type ContratoAlquiler struct {
	ID               string
	InmuebleID       string
	VigenciaDesde    time.Time
	VigenciaHasta    time.Time
	AjusteFrecuencia string
	MontoMensual     float64
	Estado           EstadoContrato
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (c *ContratoAlquiler) DiasParaVencer() int {
	return int(time.Until(c.VigenciaHasta).Hours() / 24)
}

type PagoAlquiler struct {
	ID          string
	InmuebleID  string
	Periodo     string // YYYY-MM
	Monto       float64
	FechaPago   *time.Time
	Estado      EstadoInmueble
	Comprobante *string
	CreatedAt   time.Time
}
