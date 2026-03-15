package tesoreria

import "time"

type CuentaBancaria struct {
	ID                  string
	Banco               string
	TipoCuenta          string
	NroCuenta           string
	CBU                 *string
	CCI                 *string
	SaldoActual         float64
	Moneda              string
	Activa              bool
	UltimaActualizacion time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type TipoMovimiento string

const (
	TipoIngreso TipoMovimiento = "ingreso"
	TipoEgreso  TipoMovimiento = "egreso"
)

type MovimientoBancario struct {
	ID          string
	CuentaID    string
	Tipo        TipoMovimiento
	Monto       float64
	Descripcion string
	Fecha       time.Time
	Conciliado  bool
	Referencia  *string
	CreatedAt   time.Time
}

type FlujoCajaItem struct {
	Fecha    time.Time
	Ingresos float64
	Egresos  float64
	Saldo    float64
}

type ProyeccionItem struct {
	Mes           string // YYYY-MM
	EgresosPrevistos float64
	SaldoProyectado  float64
}
