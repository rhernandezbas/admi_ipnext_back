package empleado

import "time"

type Empleado struct {
	ID           string
	Nombre       string
	Puesto       string
	Area         string
	Rol          string
	SueldoBruto  float64
	ObraSocial   string
	Activo       bool
	FechaIngreso time.Time
	Avatar       *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type EstadoLiquidacion string

const (
	EstadoBorrador  EstadoLiquidacion = "borrador"
	EstadoAprobada  EstadoLiquidacion = "aprobada"
	EstadoPagada    EstadoLiquidacion = "pagada"
)

type Liquidacion struct {
	ID          string
	EmpleadoID  string
	Periodo     string // YYYY-MM
	SueldoBruto float64
	Deducciones float64
	NetoAPagar  float64
	Estado      EstadoLiquidacion
	AprobadoPor *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (l *Liquidacion) Calcular() {
	l.NetoAPagar = l.SueldoBruto - l.Deducciones
}

type Guardia struct {
	ID         string
	EmpleadoID string
	Fecha      time.Time
	Horas      float64
	Monto      float64
	Notas      *string
	CreatedAt  time.Time
}

type TipoCompensacion string

const (
	TipoBono     TipoCompensacion = "bono"
	TipoAdelanto TipoCompensacion = "adelanto"
	TipoExtra    TipoCompensacion = "extra"
	TipoOtro     TipoCompensacion = "otro"
)

type Compensacion struct {
	ID          string
	EmpleadoID  string
	Tipo        TipoCompensacion
	Monto       float64
	Fecha       time.Time
	Descripcion *string
	CreatedAt   time.Time
}
