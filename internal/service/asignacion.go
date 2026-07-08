package service

import (
	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// AsignacionService contiene la lógica de negocio del módulo de María José.
type AsignacionService struct {
	repo storage.AsignacionRepository
}

func NuevoAsignacionService(repo storage.AsignacionRepository) *AsignacionService {
	return &AsignacionService{repo: repo}
}

func (s *AsignacionService) Listar() []models.AsignacionTriage {
	return s.repo.ListarAsignaciones()
}

func (s *AsignacionService) Obtener(id int) (models.AsignacionTriage, error) {
	a, ok := s.repo.BuscarAsignacionPorID(id)
	if !ok {
		return models.AsignacionTriage{}, ErrNoEncontrado
	}
	return a, nil
}

func (s *AsignacionService) Crear(a models.AsignacionTriage) (models.AsignacionTriage, error) {
	if err := validarAsignacion(a); err != nil {
		return models.AsignacionTriage{}, err
	}
	if a.EstadoConfirmacion == "" {
		a.EstadoConfirmacion = "Pendiente"
	}
	return s.repo.CrearAsignacion(a), nil
}

func (s *AsignacionService) Actualizar(id int, datos models.AsignacionTriage) (models.AsignacionTriage, error) {
	if err := validarAsignacion(datos); err != nil {
		return models.AsignacionTriage{}, err
	}
	actualizada, ok := s.repo.ActualizarAsignacion(id, datos)
	if !ok {
		return models.AsignacionTriage{}, ErrNoEncontrado
	}
	return actualizada, nil
}

func (s *AsignacionService) Borrar(id int) error {
	if !s.repo.BorrarAsignacion(id) {
		return ErrNoEncontrado
	}
	return nil
}

func (s *AsignacionService) ListarPorAlerta(alertaID int) []models.AsignacionTriage {
	return s.repo.ListarAsignacionesPorAlerta(alertaID)
}

// validarAsignacion aplica las reglas del módulo de enrutamiento y matching.
// Regla: alerta_id y recurso_id deben ser mayores a cero para que el match sea válido.
func validarAsignacion(a models.AsignacionTriage) error {
	if a.AlertaID <= 0 || a.RecursoID <= 0 {
		return ErrIDsInvalidos
	}
	return nil
}
