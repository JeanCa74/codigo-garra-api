package service

import (
	"strings"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// AlertaService contiene la lógica de negocio del módulo de Jean Carlos.
type AlertaService struct {
	repo storage.AlertaRepository
}

func NuevoAlertaService(repo storage.AlertaRepository) *AlertaService {
	return &AlertaService{repo: repo}
}

func (s *AlertaService) Listar() []models.AlertaEmergencia {
	return s.repo.ListarAlertas()
}

func (s *AlertaService) Obtener(id int) (models.AlertaEmergencia, error) {
	a, ok := s.repo.BuscarAlertaPorID(id)
	if !ok {
		return models.AlertaEmergencia{}, ErrNoEncontrado
	}
	return a, nil
}

func (s *AlertaService) Crear(a models.AlertaEmergencia) (models.AlertaEmergencia, error) {
	if err := validarAlerta(a); err != nil {
		return models.AlertaEmergencia{}, err
	}
	if a.Estado == "" {
		a.Estado = "Buscando"
	}
	return s.repo.CrearAlerta(a), nil
}

func (s *AlertaService) Actualizar(id int, datos models.AlertaEmergencia) (models.AlertaEmergencia, error) {
	if err := validarAlerta(datos); err != nil {
		return models.AlertaEmergencia{}, err
	}
	actualizada, ok := s.repo.ActualizarAlerta(id, datos)
	if !ok {
		return models.AlertaEmergencia{}, ErrNoEncontrado
	}
	return actualizada, nil
}

func (s *AlertaService) Borrar(id int) error {
	if !s.repo.BorrarAlerta(id) {
		return ErrNoEncontrado
	}
	return nil
}

// validarAlerta aplica las reglas de negocio de triage veterinario.
// Regla 1: requerimiento no puede estar vacío.
// Regla 2: gravedad debe estar en la escala [1,5] (1=leve, 5=crítico).
func validarAlerta(a models.AlertaEmergencia) error {
	if strings.TrimSpace(a.Requerimiento) == "" {
		return ErrRequerimientoVacio
	}
	if a.Gravedad < 1 || a.Gravedad > 5 {
		return ErrGravedadInvalida
	}
	return nil
}
