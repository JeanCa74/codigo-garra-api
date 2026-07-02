package service

import (
	"strings"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// HistorialMedicoService contiene la lógica de negocio del módulo de María José (entidad HistorialMedico).
type HistorialMedicoService struct {
	repo storage.HistorialMedicoRepository
}

func NuevoHistorialMedicoService(repo storage.HistorialMedicoRepository) *HistorialMedicoService {
	return &HistorialMedicoService{repo: repo}
}

func (s *HistorialMedicoService) Listar() []models.HistorialMedico {
	return s.repo.ListarHistorial()
}

func (s *HistorialMedicoService) ObtenerPorID(id int) (models.HistorialMedico, error) {
	h, ok := s.repo.BuscarHistorialPorID(id)
	if !ok {
		return models.HistorialMedico{}, ErrNoEncontrado
	}
	return h, nil
}

func (s *HistorialMedicoService) Crear(h models.HistorialMedico) (models.HistorialMedico, error) {
	if err := validarHistorial(h); err != nil {
		return models.HistorialMedico{}, err
	}
	return s.repo.CrearHistorial(h), nil
}

func (s *HistorialMedicoService) ListarPorMascota(mascotaID int) []models.HistorialMedico {
	return s.repo.ListarHistorialPorMascota(mascotaID)
}

// validarHistorial exige mascota_id > 0 y que el diagnóstico no esté vacío.
// Sin diagnóstico no existe entrada clínica válida; sin mascota_id no se puede asociar al paciente.
func validarHistorial(h models.HistorialMedico) error {
	if h.MascotaID <= 0 {
		return ErrMascotaIDInvalido
	}
	if strings.TrimSpace(h.Diagnostico) == "" {
		return ErrDiagnosticoVacio
	}
	return nil
}
