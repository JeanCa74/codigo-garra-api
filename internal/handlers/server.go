package handlers

import "github.com/JeanCa74/codigo-garra-api/internal/service"

// Server agrupa los servicios de los que dependen los handlers.
type Server struct {
	Alertas      *service.AlertaService
	Recursos     *service.RecursoService
	Asignaciones *service.AsignacionService
	Auth         *service.AuthService
}

func NewServer(
	alertas *service.AlertaService,
	recursos *service.RecursoService,
	asignaciones *service.AsignacionService,
	auth *service.AuthService,
) *Server {
	return &Server{
		Alertas:      alertas,
		Recursos:     recursos,
		Asignaciones: asignaciones,
		Auth:         auth,
	}
}
