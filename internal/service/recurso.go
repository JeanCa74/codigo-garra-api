package service

import (
	"strings"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// RecursoService contiene la lógica de negocio del módulo de John Erick.
type RecursoService struct {
	repo storage.RecursoRepository
}

func NuevoRecursoService(repo storage.RecursoRepository) *RecursoService {
	return &RecursoService{repo: repo}
}

func (s *RecursoService) Listar() []models.RecursoClinico {
	return s.repo.ListarRecursos()
}

func (s *RecursoService) Obtener(id int) (models.RecursoClinico, error) {
	r, ok := s.repo.BuscarRecursoPorID(id)
	if !ok {
		return models.RecursoClinico{}, ErrNoEncontrado
	}
	return r, nil
}

func (s *RecursoService) Crear(r models.RecursoClinico) (models.RecursoClinico, error) {
	if err := validarRecurso(r); err != nil {
		return models.RecursoClinico{}, err
	}
	return s.repo.CrearRecurso(r), nil
}

func (s *RecursoService) Actualizar(id int, datos models.RecursoClinico) (models.RecursoClinico, error) {
	if err := validarRecurso(datos); err != nil {
		return models.RecursoClinico{}, err
	}
	actualizado, ok := s.repo.ActualizarRecurso(id, datos)
	if !ok {
		return models.RecursoClinico{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *RecursoService) Borrar(id int) error {
	if !s.repo.BorrarRecurso(id) {
		return ErrNoEncontrado
	}
	return nil
}

// validarRecurso aplica las reglas de negocio de capacidad clínica.
// Regla: tipo_maquina no puede estar vacío (identifica el equipo requerido).
func validarRecurso(r models.RecursoClinico) error {
	if strings.TrimSpace(r.TipoMaquina) == "" {
		return ErrTipoMaquinaVacio
	}
	return nil
}
