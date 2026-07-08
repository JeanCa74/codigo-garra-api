package storage

import "github.com/JeanCa74/codigo-garra-api/internal/models"

// ── Módulo Jean Carlos: Emergencia médica ──────────────────────────────────

type AlertaRepository interface {
	ListarAlertas() []models.AlertaEmergencia
	BuscarAlertaPorID(id int) (models.AlertaEmergencia, bool)
	CrearAlerta(a models.AlertaEmergencia) models.AlertaEmergencia
	ActualizarAlerta(id int, datos models.AlertaEmergencia) (models.AlertaEmergencia, bool)
	BorrarAlerta(id int) bool
}

type AsignacionRepository interface {
	ListarAsignaciones() []models.AsignacionTriage
	BuscarAsignacionPorID(id int) (models.AsignacionTriage, bool)
	CrearAsignacion(a models.AsignacionTriage) models.AsignacionTriage
	ActualizarAsignacion(id int, datos models.AsignacionTriage) (models.AsignacionTriage, bool)
	BorrarAsignacion(id int) bool
	ListarAsignacionesPorAlerta(alertaID int) []models.AsignacionTriage
}

// ── Módulo John Erick: Perfil veterinario ─────────────────────────────────

type PerfilVeterinarioRepository interface {
	ListarPerfiles() []models.PerfilVeterinario
	BuscarPerfilPorID(id int) (models.PerfilVeterinario, bool)
	CrearPerfil(p models.PerfilVeterinario) models.PerfilVeterinario
	ActualizarPerfil(id int, datos models.PerfilVeterinario) (models.PerfilVeterinario, bool)
	BorrarPerfil(id int) bool
}

type RecursoRepository interface {
	ListarRecursos() []models.RecursoClinico
	BuscarRecursoPorID(id int) (models.RecursoClinico, bool)
	CrearRecurso(r models.RecursoClinico) models.RecursoClinico
	ActualizarRecurso(id int, datos models.RecursoClinico) (models.RecursoClinico, bool)
	BorrarRecurso(id int) bool
	ListarRecursosPorPerfil(perfilID int) []models.RecursoClinico
}

// ── Módulo María José: Historial médico ───────────────────────────────────

type MascotaRepository interface {
	ListarMascotas() []models.Mascota
	BuscarMascotaPorID(id int) (models.Mascota, bool)
	CrearMascota(m models.Mascota) models.Mascota
	ActualizarMascota(id int, datos models.Mascota) (models.Mascota, bool)
	BorrarMascota(id int) bool
}

type HistorialMedicoRepository interface {
	ListarHistorial() []models.HistorialMedico
	BuscarHistorialPorID(id int) (models.HistorialMedico, bool)
	CrearHistorial(h models.HistorialMedico) models.HistorialMedico
	ListarHistorialPorMascota(mascotaID int) []models.HistorialMedico
}

// ── Almacen unificado ──────────────────────────────────────────────────────

// Almacen combina los seis repositorios del dominio Código Garra.
type Almacen interface {
	AlertaRepository
	AsignacionRepository
	PerfilVeterinarioRepository
	RecursoRepository
	MascotaRepository
	HistorialMedicoRepository
}

// ── Auth ───────────────────────────────────────────────────────────────────

type UserRepository interface {
	CrearUsuario(u models.Usuario) (models.Usuario, error)
	BuscarUsuarioPorEmail(email string) (models.Usuario, bool)
}
