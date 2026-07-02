package service

import "errors"

// Errores de dominio. El handler los traduce a códigos HTTP:
//
//	ErrRequerimientoVacio, ErrGravedadInvalida,
//	ErrTipoMaquinaVacio, ErrIDsInvalidos     -> 400 Bad Request
//	ErrNoEncontrado                           -> 404 Not Found
//	ErrEmailEnUso                             -> 409 Conflict
//	ErrCredencialesInvalidas                  -> 401 Unauthorized
var (
	ErrRequerimientoVacio   = errors.New("el campo requerimiento es obligatorio")
	ErrGravedadInvalida     = errors.New("la gravedad debe estar entre 1 y 5")
	ErrTipoMaquinaVacio     = errors.New("el campo tipo_maquina es obligatorio")
	ErrIDsInvalidos         = errors.New("alerta_id y recurso_id deben ser mayores a cero")
	ErrNoEncontrado         = errors.New("recurso no encontrado")
	ErrEmailEnUso           = errors.New("el email ya esta registrado")
	ErrCredencialesInvalidas = errors.New("email o contrasena incorrectos")
)
