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

// perfilRepoMock registra llamadas a PerfilVeterinarioRepository.
type perfilRepoMock struct {
	mock.Mock
}

func (m *perfilRepoMock) ListarPerfiles() []models.PerfilVeterinario {
	return m.Called().Get(0).([]models.PerfilVeterinario)
}
func (m *perfilRepoMock) BuscarPerfilPorID(id int) (models.PerfilVeterinario, bool) {
	args := m.Called(id)
	return args.Get(0).(models.PerfilVeterinario), args.Bool(1)
}
func (m *perfilRepoMock) CrearPerfil(p models.PerfilVeterinario) models.PerfilVeterinario {
	return m.Called(p).Get(0).(models.PerfilVeterinario)
}
func (m *perfilRepoMock) ActualizarPerfil(id int, datos models.PerfilVeterinario) (models.PerfilVeterinario, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.PerfilVeterinario), args.Bool(1)
}
func (m *perfilRepoMock) BorrarPerfil(id int) bool {
	return m.Called(id).Bool(0)
}

var _ storage.PerfilVeterinarioRepository = (*perfilRepoMock)(nil)

// TestPerfilVeterinarioService_Crear verifica que validarPerfil rechaza
// nombre y teléfono vacíos sin llegar al repositorio.
func TestPerfilVeterinarioService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.PerfilVeterinario
		errEsperado   error
		debePersistir bool
	}{
		{
			"nombre vacio -> ErrNombreVacio y NO llega al repo",
			models.PerfilVeterinario{Nombre: "", Telefono: "0987654321"},
			service.ErrNombreVacio,
			false,
		},
		{
			"nombre solo espacios -> ErrNombreVacio y NO llega al repo",
			models.PerfilVeterinario{Nombre: "   ", Telefono: "0987654321"},
			service.ErrNombreVacio,
			false,
		},
		{
			"telefono vacio -> ErrTelefonoVacio y NO llega al repo",
			models.PerfilVeterinario{Nombre: "Clínica Garra", Telefono: ""},
			service.ErrTelefonoVacio,
			false,
		},
		{
			"telefono solo espacios -> ErrTelefonoVacio y NO llega al repo",
			models.PerfilVeterinario{Nombre: "Clínica Garra", Telefono: "   "},
			service.ErrTelefonoVacio,
			false,
		},
		{
			"perfil valido -> sin error y se persiste",
			models.PerfilVeterinario{Nombre: "Clínica Garra", Telefono: "0987654321", Activo: true},
			nil,
			true,
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := new(perfilRepoMock)

			if tc.debePersistir {
				esperado := tc.entrada
				esperado.ID = 5
				repo.On("CrearPerfil", tc.entrada).Return(esperado)
			}

			svc := service.NuevoPerfilVeterinarioService(repo)
			creado, err := svc.Crear(tc.entrada)

			if tc.errEsperado != nil {
				require.ErrorIs(t, err, tc.errEsperado)
				// Clave: la validación rechaza los datos ANTES de tocar la base de datos.
				repo.AssertNotCalled(t, "CrearPerfil")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 5, creado.ID)
				repo.AssertCalled(t, "CrearPerfil", tc.entrada)
			}

			repo.AssertExpectations(t)
		})
	}
}

// TestPerfilVeterinarioService_Listar verifica la delegación correcta al repositorio.
func TestPerfilVeterinarioService_Listar(t *testing.T) {
	repo := new(perfilRepoMock)
	esperado := []models.PerfilVeterinario{
		{ID: 1, Nombre: "Clínica Garra Norte", Telefono: "0987654321", Activo: true},
	}
	repo.On("ListarPerfiles").Return(esperado)

	svc := service.NuevoPerfilVeterinarioService(repo)
	resultado := svc.Listar()

	assert.Equal(t, 1, len(resultado))
	repo.AssertExpectations(t)
}

// TestPerfilVeterinarioService_Obtener cubre perfil existente e inexistente.
func TestPerfilVeterinarioService_Obtener(t *testing.T) {
	t.Run("perfil existente -> devuelve perfil sin error", func(t *testing.T) {
		repo := new(perfilRepoMock)
		perfil := models.PerfilVeterinario{ID: 1, Nombre: "Clínica Garra Norte", Telefono: "0987654321"}
		repo.On("BuscarPerfilPorID", 1).Return(perfil, true)

		svc := service.NuevoPerfilVeterinarioService(repo)
		resultado, err := svc.Obtener(1)

		require.NoError(t, err)
		assert.Equal(t, perfil.Nombre, resultado.Nombre)
		repo.AssertExpectations(t)
	})
	t.Run("perfil inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(perfilRepoMock)
		repo.On("BuscarPerfilPorID", 99).Return(models.PerfilVeterinario{}, false)

		svc := service.NuevoPerfilVeterinarioService(repo)
		_, err := svc.Obtener(99)

		require.ErrorIs(t, err, service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}

// TestPerfilVeterinarioService_Actualizar cubre actualización exitosa e inexistente.
func TestPerfilVeterinarioService_Actualizar(t *testing.T) {
	t.Run("datos validos y perfil existente -> devuelve actualizado", func(t *testing.T) {
		repo := new(perfilRepoMock)
		datos := models.PerfilVeterinario{Nombre: "Clínica Garra Sur", Telefono: "0911111111", Activo: false}
		actualizado := datos
		actualizado.ID = 1
		repo.On("ActualizarPerfil", 1, datos).Return(actualizado, true)

		svc := service.NuevoPerfilVeterinarioService(repo)
		resultado, err := svc.Actualizar(1, datos)

		require.NoError(t, err)
		assert.Equal(t, "Clínica Garra Sur", resultado.Nombre)
		repo.AssertExpectations(t)
	})
	t.Run("perfil inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(perfilRepoMock)
		datos := models.PerfilVeterinario{Nombre: "Clínica X", Telefono: "000"}
		repo.On("ActualizarPerfil", 99, datos).Return(models.PerfilVeterinario{}, false)

		svc := service.NuevoPerfilVeterinarioService(repo)
		_, err := svc.Actualizar(99, datos)

		require.ErrorIs(t, err, service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}

// TestPerfilVeterinarioService_Borrar cubre eliminación exitosa e inexistente.
func TestPerfilVeterinarioService_Borrar(t *testing.T) {
	t.Run("perfil existente -> nil", func(t *testing.T) {
		repo := new(perfilRepoMock)
		repo.On("BorrarPerfil", 1).Return(true)

		svc := service.NuevoPerfilVeterinarioService(repo)
		require.NoError(t, svc.Borrar(1))
		repo.AssertExpectations(t)
	})
	t.Run("perfil inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(perfilRepoMock)
		repo.On("BorrarPerfil", 99).Return(false)

		svc := service.NuevoPerfilVeterinarioService(repo)
		require.ErrorIs(t, svc.Borrar(99), service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}
