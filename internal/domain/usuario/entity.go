package usuario

import "time"

type Rol string

const (
	RolAdmin      Rol = "admin"
	RolSubUsuario Rol = "sub-usuario"
)

type NivelPermiso string

const (
	NivelNinguno   NivelPermiso = "ninguno"
	NivelLectura   NivelPermiso = "lectura"
	NivelEscritura NivelPermiso = "escritura"
)

type Permisos struct {
	Dashboard      bool         `json:"dashboard"`
	Transferencias NivelPermiso `json:"transferencias"`
	Nominas        NivelPermiso `json:"nominas"`
	Proveedores    NivelPermiso `json:"proveedores"`
	Servicios      NivelPermiso `json:"servicios"`
	Alquileres     NivelPermiso `json:"alquileres"`
	Tesoreria      NivelPermiso `json:"tesoreria"`
	Reportes       NivelPermiso `json:"reportes"`
}

type Usuario struct {
	ID        string
	Nombre    string
	Email     string
	Password  string
	Rol       Rol
	Permisos  Permisos
	Avatar    *string
	Activo    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *Usuario) TienePermiso(modulo string, nivelMinimo NivelPermiso) bool {
	if u.Rol == RolAdmin {
		return true
	}
	nivel := u.nivelModulo(modulo)
	return nivelSuficiente(nivel, nivelMinimo)
}

func (u *Usuario) nivelModulo(modulo string) NivelPermiso {
	switch modulo {
	case "transferencias":
		return u.Permisos.Transferencias
	case "nominas":
		return u.Permisos.Nominas
	case "proveedores":
		return u.Permisos.Proveedores
	case "servicios":
		return u.Permisos.Servicios
	case "alquileres":
		return u.Permisos.Alquileres
	case "tesoreria":
		return u.Permisos.Tesoreria
	case "reportes":
		return u.Permisos.Reportes
	default:
		return NivelNinguno
	}
}

func nivelSuficiente(nivel, minimo NivelPermiso) bool {
	orden := map[NivelPermiso]int{
		NivelNinguno:   0,
		NivelLectura:   1,
		NivelEscritura: 2,
	}
	return orden[nivel] >= orden[minimo]
}
