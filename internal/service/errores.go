package service

import "errors"

// Errores de dominio → traducción HTTP en handlers/respond.go
//
//	400 Bad Request  : ErrRequerimientoVacio, ErrGravedadInvalida, ErrIDsInvalidos,
//	                   ErrNombreVacio, ErrTelefonoVacio, ErrTipoMaquinaVacio,
//	                   ErrDiagnosticoVacio, ErrMascotaIDInvalido
//	404 Not Found    : ErrNoEncontrado
//	409 Conflict     : ErrEmailEnUso
//	401 Unauthorized : ErrCredencialesInvalidas (y el middleware JWT)
var (
	// Jean Carlos — Emergencia médica
	ErrRequerimientoVacio = errors.New("el campo requerimiento es obligatorio")
	ErrGravedadInvalida   = errors.New("la gravedad debe estar entre 1 y 5")
	ErrIDsInvalidos       = errors.New("alerta_id y recurso_id deben ser mayores a cero")

	// John Erick — Perfil veterinario
	ErrNombreVacio      = errors.New("el campo nombre es obligatorio")
	ErrTelefonoVacio    = errors.New("el campo telefono es obligatorio")
	ErrTipoMaquinaVacio = errors.New("el campo tipo_maquina es obligatorio")

	// María José — Historial médico
	ErrDiagnosticoVacio  = errors.New("el campo diagnostico es obligatorio")
	ErrMascotaIDInvalido = errors.New("el mascota_id debe ser mayor a cero")

	// Compartidos
	ErrNoEncontrado          = errors.New("recurso no encontrado")
	ErrEmailEnUso            = errors.New("el email ya esta registrado")
	ErrCredencialesInvalidas = errors.New("email o contrasena incorrectos")
)
