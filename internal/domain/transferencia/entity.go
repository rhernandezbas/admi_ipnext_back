package transferencia

import "time"

type Estado string

const (
	EstadoPendiente   Estado = "pendiente"
	EstadoPagado      Estado = "pagado"
	EstadoVencido     Estado = "vencido"
	EstadoProgramado  Estado = "programado"
	EstadoEnProceso   Estado = "en_proceso"
)

type Frecuencia string

const (
	FrecuenciaManual     Frecuencia = "manual"
	FrecuenciaMensual    Frecuencia = "mensual"
	FrecuenciaSemanal    Frecuencia = "semanal"
	FrecuenciaQuincenal  Frecuencia = "quincenal"
	FrecuenciaSemestral  Frecuencia = "semestral"
	FrecuenciaAnual      Frecuencia = "anual"
)

type MetodoPago string

const (
	MetodoTransferencia MetodoPago = "transferencia"
	MetodoDebito        MetodoPago = "debito"
	MetodoEfectivo      MetodoPago = "efectivo"
	MetodoCheque        MetodoPago = "cheque"
)

type Transferencia struct {
	ID               string
	Beneficiario     string
	CBU              *string
	Alias            *string
	Categoria        string
	Monto            float64
	Moneda           string
	FechaPago        time.Time
	FechaVencimiento *time.Time
	Frecuencia       Frecuencia
	Estado           Estado
	MetodoPago       MetodoPago
	Notas            *string
	ProveedorID      *string
	CreadoPor        string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (t *Transferencia) EsRecurrente() bool {
	return t.Frecuencia != FrecuenciaManual
}

func (t *Transferencia) EstaVencida() bool {
	if t.FechaVencimiento == nil {
		return false
	}
	return time.Now().After(*t.FechaVencimiento) && t.Estado == EstadoPendiente
}
