package proveedor

import "time"

type Proveedor struct {
	ID        string
	Nombre    string
	CUIT      string
	CBU       *string
	Alias     *string
	Email     *string
	Categoria string
	SitioWeb  *string
	Activo    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EstadoContrato string

const (
	EstadoActivo          EstadoContrato = "activo"
	EstadoProximoAVencer  EstadoContrato = "proximo_a_vencer"
	EstadoVencido         EstadoContrato = "vencido"
	EstadoEnProceso       EstadoContrato = "en_proceso"
)

type ContratoProveedor struct {
	ID            string
	Codigo        string // CTR-2024-001
	ProveedorID   string
	Descripcion   string
	VigenciaDesde time.Time
	VigenciaHasta time.Time
	MontoAnual    float64
	Estado        EstadoContrato
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (c *ContratoProveedor) DiasParaVencer() int {
	diff := time.Until(c.VigenciaHasta)
	return int(diff.Hours() / 24)
}

type RankingItem struct {
	ProveedorID     string
	NombreProveedor string
	TotalPagado     float64
	CantidadPagos   int64
}
