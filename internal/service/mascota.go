package service

import (
	"strings"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// MascotaService contiene la lógica de negocio del módulo de María José (entidad Mascota).
type MascotaService struct {
	repo storage.MascotaRepository
}

func NuevoMascotaService(repo storage.MascotaRepository) *MascotaService {
	return &MascotaService{repo: repo}
}

func (s *MascotaService) Listar() []models.Mascota {
	return s.repo.ListarMascotas()
}

func (s *MascotaService) Obtener(id int) (models.Mascota, error) {
	m, ok := s.repo.BuscarMascotaPorID(id)
	if !ok {
		return models.Mascota{}, ErrNoEncontrado
	}
	return m, nil
}

func (s *MascotaService) Crear(m models.Mascota) (models.Mascota, error) {
	if err := validarMascota(m); err != nil {
		return models.Mascota{}, err
	}
	return s.repo.CrearMascota(m), nil
}

func (s *MascotaService) Actualizar(id int, datos models.Mascota) (models.Mascota, error) {
	if err := validarMascota(datos); err != nil {
		return models.Mascota{}, err
	}
	actualizada, ok := s.repo.ActualizarMascota(id, datos)
	if !ok {
		return models.Mascota{}, ErrNoEncontrado
	}
	return actualizada, nil
}

func (s *MascotaService) Borrar(id int) error {
	if !s.repo.BorrarMascota(id) {
		return ErrNoEncontrado
	}
	return nil
}

// validarMascota: el nombre de la mascota es el campo mínimo requerido para registrar al paciente.
func validarMascota(m models.Mascota) error {
	if strings.TrimSpace(m.Nombre) == "" {
		return ErrNombreVacio
	}
	return nil
}
