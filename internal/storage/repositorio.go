package storage

import "github.com/JeanCa74/codigo-garra-api/internal/models"

// AlertaRepository es el contrato de persistencia del módulo de Jean Carlos.
type AlertaRepository interface {
	ListarAlertas() []models.AlertaEmergencia
	BuscarAlertaPorID(id int) (models.AlertaEmergencia, bool)
	CrearAlerta(a models.AlertaEmergencia) models.AlertaEmergencia
	ActualizarAlerta(id int, datos models.AlertaEmergencia) (models.AlertaEmergencia, bool)
	BorrarAlerta(id int) bool
}

// RecursoRepository es el contrato de persistencia del módulo de John Erick.
type RecursoRepository interface {
	ListarRecursos() []models.RecursoClinico
	BuscarRecursoPorID(id int) (models.RecursoClinico, bool)
	CrearRecurso(r models.RecursoClinico) models.RecursoClinico
	ActualizarRecurso(id int, datos models.RecursoClinico) (models.RecursoClinico, bool)
	BorrarRecurso(id int) bool
}

// AsignacionRepository es el contrato de persistencia del módulo de María José.
type AsignacionRepository interface {
	ListarAsignaciones() []models.AsignacionTriage
	BuscarAsignacionPorID(id int) (models.AsignacionTriage, bool)
	CrearAsignacion(a models.AsignacionTriage) models.AsignacionTriage
	ActualizarAsignacion(id int, datos models.AsignacionTriage) (models.AsignacionTriage, bool)
	BorrarAsignacion(id int) bool
}

// Almacen combina los tres repositorios del dominio Código Garra.
type Almacen interface {
	AlertaRepository
	RecursoRepository
	AsignacionRepository
}

// UserRepository es el contrato de persistencia de usuarios (auth).
type UserRepository interface {
	CrearUsuario(u models.Usuario) (models.Usuario, error)
	BuscarUsuarioPorEmail(email string) (models.Usuario, bool)
}
