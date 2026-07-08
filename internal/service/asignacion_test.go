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

// asignacionRepoMock registra llamadas a AsignacionRepository.
type asignacionRepoMock struct {
	mock.Mock
}

func (m *asignacionRepoMock) ListarAsignaciones() []models.AsignacionTriage {
	return m.Called().Get(0).([]models.AsignacionTriage)
}
func (m *asignacionRepoMock) BuscarAsignacionPorID(id int) (models.AsignacionTriage, bool) {
	args := m.Called(id)
	return args.Get(0).(models.AsignacionTriage), args.Bool(1)
}
func (m *asignacionRepoMock) CrearAsignacion(a models.AsignacionTriage) models.AsignacionTriage {
	return m.Called(a).Get(0).(models.AsignacionTriage)
}
func (m *asignacionRepoMock) ActualizarAsignacion(id int, datos models.AsignacionTriage) (models.AsignacionTriage, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.AsignacionTriage), args.Bool(1)
}
func (m *asignacionRepoMock) BorrarAsignacion(id int) bool {
	return m.Called(id).Bool(0)
}
func (m *asignacionRepoMock) ListarAsignacionesPorAlerta(alertaID int) []models.AsignacionTriage {
	return m.Called(alertaID).Get(0).([]models.AsignacionTriage)
}

var _ storage.AsignacionRepository = (*asignacionRepoMock)(nil)

// TestAsignacionService_Crear verifica que validarAsignacion rechaza IDs inválidos
// sin llamar al repositorio.
func TestAsignacionService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.AsignacionTriage
		errEsperado   error
		debePersistir bool
	}{
		{
			"alerta_id=0 -> ErrIDsInvalidos y NO llega al repo",
			models.AsignacionTriage{AlertaID: 0, RecursoID: 5},
			service.ErrIDsInvalidos,
			false,
		},
		{
			"recurso_id=0 -> ErrIDsInvalidos y NO llega al repo",
			models.AsignacionTriage{AlertaID: 3, RecursoID: 0},
			service.ErrIDsInvalidos,
			false,
		},
		{
			"ambos IDs validos -> sin error y se persiste",
			models.AsignacionTriage{AlertaID: 2, RecursoID: 4, EstadoConfirmacion: "Pendiente"},
			nil,
			true,
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := new(asignacionRepoMock)

			if tc.debePersistir {
				esperado := tc.entrada
				esperado.ID = 30
				repo.On("CrearAsignacion", tc.entrada).Return(esperado)
			}

			svc := service.NuevoAsignacionService(repo)
			creada, err := svc.Crear(tc.entrada)

			if tc.errEsperado != nil {
				require.ErrorIs(t, err, tc.errEsperado)
				repo.AssertNotCalled(t, "CrearAsignacion")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 30, creada.ID)
				repo.AssertCalled(t, "CrearAsignacion", tc.entrada)
			}

			repo.AssertExpectations(t)
		})
	}
}

// TestAsignacionService_Listar verifica la delegación al repositorio.
func TestAsignacionService_Listar(t *testing.T) {
	repo := new(asignacionRepoMock)
	esperado := []models.AsignacionTriage{
		{ID: 1, AlertaID: 1, RecursoID: 2, EstadoConfirmacion: "Pendiente"},
	}
	repo.On("ListarAsignaciones").Return(esperado)

	svc := service.NuevoAsignacionService(repo)
	resultado := svc.Listar()

	assert.Equal(t, 1, len(resultado))
	repo.AssertExpectations(t)
}

// TestAsignacionService_Obtener cubre asignación existente e inexistente.
func TestAsignacionService_Obtener(t *testing.T) {
	t.Run("asignación existente -> devuelve sin error", func(t *testing.T) {
		repo := new(asignacionRepoMock)
		asig := models.AsignacionTriage{ID: 1, AlertaID: 1, RecursoID: 2, EstadoConfirmacion: "Pendiente"}
		repo.On("BuscarAsignacionPorID", 1).Return(asig, true)

		svc := service.NuevoAsignacionService(repo)
		resultado, err := svc.Obtener(1)

		require.NoError(t, err)
		assert.Equal(t, 1, resultado.AlertaID)
		repo.AssertExpectations(t)
	})
	t.Run("asignación inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(asignacionRepoMock)
		repo.On("BuscarAsignacionPorID", 99).Return(models.AsignacionTriage{}, false)

		svc := service.NuevoAsignacionService(repo)
		_, err := svc.Obtener(99)

		require.ErrorIs(t, err, service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}

// TestAsignacionService_Actualizar cubre actualización exitosa e inexistente.
func TestAsignacionService_Actualizar(t *testing.T) {
	t.Run("datos validos y asignación existente -> devuelve actualizada", func(t *testing.T) {
		repo := new(asignacionRepoMock)
		datos := models.AsignacionTriage{AlertaID: 1, RecursoID: 2, EstadoConfirmacion: "Confirmado"}
		actualizada := datos
		actualizada.ID = 1
		repo.On("ActualizarAsignacion", 1, datos).Return(actualizada, true)

		svc := service.NuevoAsignacionService(repo)
		resultado, err := svc.Actualizar(1, datos)

		require.NoError(t, err)
		assert.Equal(t, "Confirmado", resultado.EstadoConfirmacion)
		repo.AssertExpectations(t)
	})
	t.Run("asignación inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(asignacionRepoMock)
		datos := models.AsignacionTriage{AlertaID: 1, RecursoID: 1, EstadoConfirmacion: "Pendiente"}
		repo.On("ActualizarAsignacion", 99, datos).Return(models.AsignacionTriage{}, false)

		svc := service.NuevoAsignacionService(repo)
		_, err := svc.Actualizar(99, datos)

		require.ErrorIs(t, err, service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}

// TestAsignacionService_Borrar cubre eliminación exitosa e inexistente.
func TestAsignacionService_Borrar(t *testing.T) {
	t.Run("asignación existente -> nil", func(t *testing.T) {
		repo := new(asignacionRepoMock)
		repo.On("BorrarAsignacion", 1).Return(true)

		svc := service.NuevoAsignacionService(repo)
		require.NoError(t, svc.Borrar(1))
		repo.AssertExpectations(t)
	})
	t.Run("asignación inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(asignacionRepoMock)
		repo.On("BorrarAsignacion", 99).Return(false)

		svc := service.NuevoAsignacionService(repo)
		require.ErrorIs(t, svc.Borrar(99), service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}
