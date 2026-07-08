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

// recursoRepoMock registra llamadas a RecursoRepository.
type recursoRepoMock struct {
	mock.Mock
}

func (m *recursoRepoMock) ListarRecursos() []models.RecursoClinico {
	return m.Called().Get(0).([]models.RecursoClinico)
}
func (m *recursoRepoMock) BuscarRecursoPorID(id int) (models.RecursoClinico, bool) {
	args := m.Called(id)
	return args.Get(0).(models.RecursoClinico), args.Bool(1)
}
func (m *recursoRepoMock) CrearRecurso(r models.RecursoClinico) models.RecursoClinico {
	return m.Called(r).Get(0).(models.RecursoClinico)
}
func (m *recursoRepoMock) ActualizarRecurso(id int, datos models.RecursoClinico) (models.RecursoClinico, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.RecursoClinico), args.Bool(1)
}
func (m *recursoRepoMock) BorrarRecurso(id int) bool {
	return m.Called(id).Bool(0)
}
func (m *recursoRepoMock) ListarRecursosPorPerfil(perfilID int) []models.RecursoClinico {
	return m.Called(perfilID).Get(0).([]models.RecursoClinico)
}

var _ storage.RecursoRepository = (*recursoRepoMock)(nil)

// TestRecursoService_Crear verifica que validarRecurso rechaza tipo_maquina vacío
// sin llamar al repositorio, y que un recurso válido sí se persiste.
func TestRecursoService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.RecursoClinico
		errEsperado   error
		debePersistir bool
	}{
		{
			"tipo_maquina vacio -> ErrTipoMaquinaVacio y NO llega al repo",
			models.RecursoClinico{PerfilID: 1, TipoMaquina: "", EstaDisponible: true},
			service.ErrTipoMaquinaVacio,
			false,
		},
		{
			"tipo_maquina solo espacios -> ErrTipoMaquinaVacio y NO llega al repo",
			models.RecursoClinico{PerfilID: 2, TipoMaquina: "   ", EstaDisponible: true},
			service.ErrTipoMaquinaVacio,
			false,
		},
		{
			"recurso valido -> sin error y se persiste",
			models.RecursoClinico{PerfilID: 3, TipoMaquina: "Ecógrafo", EstaDisponible: true},
			nil,
			true,
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := new(recursoRepoMock)

			if tc.debePersistir {
				esperado := tc.entrada
				esperado.ID = 20
				repo.On("CrearRecurso", tc.entrada).Return(esperado)
			}

			svc := service.NuevoRecursoService(repo)
			creado, err := svc.Crear(tc.entrada)

			if tc.errEsperado != nil {
				require.ErrorIs(t, err, tc.errEsperado)
				repo.AssertNotCalled(t, "CrearRecurso")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 20, creado.ID)
				repo.AssertCalled(t, "CrearRecurso", tc.entrada)
			}

			repo.AssertExpectations(t)
		})
	}
}

// TestRecursoService_Listar verifica la delegación al repositorio.
func TestRecursoService_Listar(t *testing.T) {
	repo := new(recursoRepoMock)
	esperado := []models.RecursoClinico{
		{ID: 1, PerfilID: 2, TipoMaquina: "Ecógrafo", EstaDisponible: true},
	}
	repo.On("ListarRecursos").Return(esperado)

	svc := service.NuevoRecursoService(repo)
	resultado := svc.Listar()

	assert.Equal(t, 1, len(resultado))
	repo.AssertExpectations(t)
}

// TestRecursoService_Obtener cubre recurso existente e inexistente.
func TestRecursoService_Obtener(t *testing.T) {
	t.Run("recurso existente -> devuelve sin error", func(t *testing.T) {
		repo := new(recursoRepoMock)
		recurso := models.RecursoClinico{ID: 1, PerfilID: 2, TipoMaquina: "Ecógrafo", EstaDisponible: true}
		repo.On("BuscarRecursoPorID", 1).Return(recurso, true)

		svc := service.NuevoRecursoService(repo)
		resultado, err := svc.Obtener(1)

		require.NoError(t, err)
		assert.Equal(t, "Ecógrafo", resultado.TipoMaquina)
		repo.AssertExpectations(t)
	})
	t.Run("recurso inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(recursoRepoMock)
		repo.On("BuscarRecursoPorID", 99).Return(models.RecursoClinico{}, false)

		svc := service.NuevoRecursoService(repo)
		_, err := svc.Obtener(99)

		require.ErrorIs(t, err, service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}

// TestRecursoService_Actualizar cubre actualización exitosa e inexistente.
func TestRecursoService_Actualizar(t *testing.T) {
	t.Run("datos validos y recurso existente -> devuelve actualizado", func(t *testing.T) {
		repo := new(recursoRepoMock)
		datos := models.RecursoClinico{PerfilID: 1, TipoMaquina: "Ventilador", EstaDisponible: false}
		actualizado := datos
		actualizado.ID = 1
		repo.On("ActualizarRecurso", 1, datos).Return(actualizado, true)

		svc := service.NuevoRecursoService(repo)
		resultado, err := svc.Actualizar(1, datos)

		require.NoError(t, err)
		assert.Equal(t, "Ventilador", resultado.TipoMaquina)
		repo.AssertExpectations(t)
	})
	t.Run("recurso inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(recursoRepoMock)
		datos := models.RecursoClinico{PerfilID: 1, TipoMaquina: "Ecógrafo", EstaDisponible: true}
		repo.On("ActualizarRecurso", 99, datos).Return(models.RecursoClinico{}, false)

		svc := service.NuevoRecursoService(repo)
		_, err := svc.Actualizar(99, datos)

		require.ErrorIs(t, err, service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}

// TestRecursoService_Borrar cubre eliminación exitosa e inexistente.
func TestRecursoService_Borrar(t *testing.T) {
	t.Run("recurso existente -> nil", func(t *testing.T) {
		repo := new(recursoRepoMock)
		repo.On("BorrarRecurso", 1).Return(true)

		svc := service.NuevoRecursoService(repo)
		require.NoError(t, svc.Borrar(1))
		repo.AssertExpectations(t)
	})
	t.Run("recurso inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(recursoRepoMock)
		repo.On("BorrarRecurso", 99).Return(false)

		svc := service.NuevoRecursoService(repo)
		require.ErrorIs(t, svc.Borrar(99), service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}
