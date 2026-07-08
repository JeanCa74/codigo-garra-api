package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/service"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// mascotaRepoMock registra llamadas a MascotaRepository.
type mascotaRepoMock struct {
	mock.Mock
}

func (m *mascotaRepoMock) ListarMascotas() []models.Mascota {
	return m.Called().Get(0).([]models.Mascota)
}
func (m *mascotaRepoMock) BuscarMascotaPorID(id int) (models.Mascota, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Mascota), args.Bool(1)
}
func (m *mascotaRepoMock) CrearMascota(ma models.Mascota) models.Mascota {
	return m.Called(ma).Get(0).(models.Mascota)
}
func (m *mascotaRepoMock) ActualizarMascota(id int, datos models.Mascota) (models.Mascota, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Mascota), args.Bool(1)
}
func (m *mascotaRepoMock) BorrarMascota(id int) bool {
	return m.Called(id).Bool(0)
}

var _ storage.MascotaRepository = (*mascotaRepoMock)(nil)

// TestMascotaService_Crear verifica que validarMascota rechaza nombre vacío
// sin llamar al repositorio, y que una mascota válida sí se persiste.
func TestMascotaService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Mascota
		errEsperado   error
		debePersistir bool
	}{
		{
			"nombre vacio -> ErrNombreVacio y NO llega al repo",
			models.Mascota{Nombre: "", Especie: "Perro", Dueno: "Ana"},
			service.ErrNombreVacio,
			false,
		},
		{
			"nombre solo espacios -> ErrNombreVacio y NO llega al repo",
			models.Mascota{Nombre: "   ", Especie: "Gato", Dueno: "Luis"},
			service.ErrNombreVacio,
			false,
		},
		{
			"mascota valida -> sin error y se persiste",
			models.Mascota{Nombre: "Max", Especie: "Perro", Edad: 3, Dueno: "Ana"},
			nil,
			true,
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := new(mascotaRepoMock)

			if tc.debePersistir {
				esperado := tc.entrada
				esperado.ID = 8
				repo.On("CrearMascota", tc.entrada).Return(esperado)
			}

			svc := service.NuevoMascotaService(repo)
			creada, err := svc.Crear(tc.entrada)

			if tc.errEsperado != nil {
				require.ErrorIs(t, err, tc.errEsperado)
				repo.AssertNotCalled(t, "CrearMascota")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 8, creada.ID)
				repo.AssertCalled(t, "CrearMascota", tc.entrada)
			}

			repo.AssertExpectations(t)
		})
	}
}

// TestMascotaService_Listar verifica la delegación al repositorio.
func TestMascotaService_Listar(t *testing.T) {
	repo := new(mascotaRepoMock)
	esperado := []models.Mascota{
		{ID: 1, Nombre: "Max", Especie: "Perro", Edad: 3, Dueno: "Ana"},
	}
	repo.On("ListarMascotas").Return(esperado)

	svc := service.NuevoMascotaService(repo)
	resultado := svc.Listar()

	assert.Equal(t, 1, len(resultado))
	repo.AssertExpectations(t)
}

// TestMascotaService_Obtener cubre mascota existente e inexistente.
func TestMascotaService_Obtener(t *testing.T) {
	t.Run("mascota existente -> devuelve sin error", func(t *testing.T) {
		repo := new(mascotaRepoMock)
		mascota := models.Mascota{ID: 1, Nombre: "Luna", Especie: "Gato", Edad: 5, Dueno: "Carlos"}
		repo.On("BuscarMascotaPorID", 1).Return(mascota, true)

		svc := service.NuevoMascotaService(repo)
		resultado, err := svc.Obtener(1)

		require.NoError(t, err)
		assert.Equal(t, "Luna", resultado.Nombre)
		repo.AssertExpectations(t)
	})
	t.Run("mascota inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(mascotaRepoMock)
		repo.On("BuscarMascotaPorID", 99).Return(models.Mascota{}, false)

		svc := service.NuevoMascotaService(repo)
		_, err := svc.Obtener(99)

		require.ErrorIs(t, err, service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}

// TestMascotaService_Actualizar cubre actualización exitosa e inexistente.
func TestMascotaService_Actualizar(t *testing.T) {
	t.Run("datos validos y mascota existente -> devuelve actualizada", func(t *testing.T) {
		repo := new(mascotaRepoMock)
		datos := models.Mascota{Nombre: "Max Actualizado", Especie: "Perro", Edad: 4, Dueno: "Ana"}
		actualizada := datos
		actualizada.ID = 1
		repo.On("ActualizarMascota", 1, datos).Return(actualizada, true)

		svc := service.NuevoMascotaService(repo)
		resultado, err := svc.Actualizar(1, datos)

		require.NoError(t, err)
		assert.Equal(t, "Max Actualizado", resultado.Nombre)
		repo.AssertExpectations(t)
	})
	t.Run("mascota inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(mascotaRepoMock)
		datos := models.Mascota{Nombre: "Fantasma", Especie: "Perro", Dueno: "Nadie"}
		repo.On("ActualizarMascota", 99, datos).Return(models.Mascota{}, false)

		svc := service.NuevoMascotaService(repo)
		_, err := svc.Actualizar(99, datos)

		require.ErrorIs(t, err, service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}
