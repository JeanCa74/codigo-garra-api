package handlers

import "github.com/JeanCa74/codigo-garra-api/internal/service"

// Server agrupa los servicios de los tres módulos del proyecto Código Garra.
type Server struct {
	// Jean Carlos — Emergencia médica
	Alertas      *service.AlertaService
	Asignaciones *service.AsignacionService
	// John Erick — Perfil veterinario
	Perfiles *service.PerfilVeterinarioService
	Recursos *service.RecursoService
	// María José — Historial médico
	Mascotas    *service.MascotaService
	HistorialMed *service.HistorialMedicoService
	// Auth compartido
	Auth *service.AuthService
}

func NewServer(
	alertas *service.AlertaService,
	asignaciones *service.AsignacionService,
	perfiles *service.PerfilVeterinarioService,
	recursos *service.RecursoService,
	mascotas *service.MascotaService,
	historialMed *service.HistorialMedicoService,
	auth *service.AuthService,
) *Server {
	return &Server{
		Alertas:      alertas,
		Asignaciones: asignaciones,
		Perfiles:     perfiles,
		Recursos:     recursos,
		Mascotas:     mascotas,
		HistorialMed: historialMed,
		Auth:         auth,
	}
}
