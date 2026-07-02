package service

import (
	"strings"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// PerfilVeterinarioService contiene la lógica de negocio del módulo de John Erick.
type PerfilVeterinarioService struct {
	repo storage.PerfilVeterinarioRepository
}

func NuevoPerfilVeterinarioService(repo storage.PerfilVeterinarioRepository) *PerfilVeterinarioService {
	return &PerfilVeterinarioService{repo: repo}
}

func (s *PerfilVeterinarioService) Listar() []models.PerfilVeterinario {
	return s.repo.ListarPerfiles()
}

func (s *PerfilVeterinarioService) Obtener(id int) (models.PerfilVeterinario, error) {
	p, ok := s.repo.BuscarPerfilPorID(id)
	if !ok {
		return models.PerfilVeterinario{}, ErrNoEncontrado
	}
	return p, nil
}

func (s *PerfilVeterinarioService) Crear(p models.PerfilVeterinario) (models.PerfilVeterinario, error) {
	if err := validarPerfil(p); err != nil {
		return models.PerfilVeterinario{}, err
	}
	return s.repo.CrearPerfil(p), nil
}

func (s *PerfilVeterinarioService) Actualizar(id int, datos models.PerfilVeterinario) (models.PerfilVeterinario, error) {
	if err := validarPerfil(datos); err != nil {
		return models.PerfilVeterinario{}, err
	}
	actualizado, ok := s.repo.ActualizarPerfil(id, datos)
	if !ok {
		return models.PerfilVeterinario{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *PerfilVeterinarioService) Borrar(id int) error {
	if !s.repo.BorrarPerfil(id) {
		return ErrNoEncontrado
	}
	return nil
}

// validarPerfil aplica las reglas de negocio de las clínicas veterinarias.
// Regla 1: el nombre de la clínica no puede estar vacío.
// Regla 2: el teléfono de contacto es obligatorio para emergencias.
func validarPerfil(p models.PerfilVeterinario) error {
	if strings.TrimSpace(p.Nombre) == "" {
		return ErrNombreVacio
	}
	if strings.TrimSpace(p.Telefono) == "" {
		return ErrTelefonoVacio
	}
	return nil
}
