package handlers

import "github.com/JeanCa74/codigo-garra-api/internal/service"

// ServerDeps agrupa las dependencias para evitar el "Constructor Hell".
type ServerDeps struct {
	// Jean Carlos — Emergencia médica
	Alertas      *service.AlertaService
	Asignaciones *service.AsignacionService
	// John Erick — Perfil veterinario
	Perfiles *service.PerfilVeterinarioService
	Recursos *service.RecursoService
	// María José — Historial médico
	Mascotas     *service.MascotaService
	HistorialMed *service.HistorialMedicoService
	// Auth compartido
	Auth *service.AuthService
}

// Server agrupa los servicios de los tres módulos del proyecto Código Garra.
type Server struct {
	deps ServerDeps
}

// NewServer recibe el struct de dependencias para una inicialización limpia.
func NewServer(deps ServerDeps) *Server {
	return &Server{
		deps: deps,
	}
}
