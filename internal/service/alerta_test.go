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

// alertaRepoMock registra llamadas a AlertaRepository para verificar que
// la validación rechaza datos antes de llegar a la base de datos.
type alertaRepoMock struct {
	mock.Mock
}

func (m *alertaRepoMock) ListarAlertas() []models.AlertaEmergencia {
	return m.Called().Get(0).([]models.AlertaEmergencia)
}
func (m *alertaRepoMock) BuscarAlertaPorID(id int) (models.AlertaEmergencia, bool) {
	args := m.Called(id)
	return args.Get(0).(models.AlertaEmergencia), args.Bool(1)
}
func (m *alertaRepoMock) CrearAlerta(a models.AlertaEmergencia) models.AlertaEmergencia {
	return m.Called(a).Get(0).(models.AlertaEmergencia)
}
func (m *alertaRepoMock) ActualizarAlerta(id int, datos models.AlertaEmergencia) (models.AlertaEmergencia, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.AlertaEmergencia), args.Bool(1)
}
func (m *alertaRepoMock) BorrarAlerta(id int) bool {
	return m.Called(id).Bool(0)
}

var _ storage.AlertaRepository = (*alertaRepoMock)(nil)

// TestAlertaService_Crear verifica que validarAlerta rechaza datos inválidos
// sin tocar el repositorio, y que datos válidos sí persisten.
func TestAlertaService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.AlertaEmergencia
		errEsperado   error
		debePersistir bool
	}{
		{
			"gravedad 0 -> ErrGravedadInvalida y NO llega al repo",
			models.AlertaEmergencia{Gravedad: 0, Requerimiento: "Ecógrafo"},
			service.ErrGravedadInvalida,
			false,
		},
		{
			"gravedad 6 fuera de escala -> ErrGravedadInvalida y NO llega al repo",
			models.AlertaEmergencia{Gravedad: 6, Requerimiento: "Ventilador"},
			service.ErrGravedadInvalida,
			false,
		},
		{
			"requerimiento vacio -> ErrRequerimientoVacio y NO llega al repo",
			models.AlertaEmergencia{Gravedad: 3, Requerimiento: "   "},
			service.ErrRequerimientoVacio,
			false,
		},
		{
			"datos validos -> sin error y se persiste",
			models.AlertaEmergencia{Gravedad: 4, Requerimiento: "Ecógrafo", Estado: "Buscando"},
			nil,
			true,
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := new(alertaRepoMock)

			if tc.debePersistir {
				esperado := tc.entrada
				esperado.ID = 10
				repo.On("CrearAlerta", tc.entrada).Return(esperado)
			}

			svc := service.NuevoAlertaService(repo)
			creada, err := svc.Crear(tc.entrada)

			if tc.errEsperado != nil {
				require.ErrorIs(t, err, tc.errEsperado)
				// La clave: el repositorio NO debe recibir la llamada.
				repo.AssertNotCalled(t, "CrearAlerta")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 10, creada.ID)
				repo.AssertCalled(t, "CrearAlerta", tc.entrada)
			}

			repo.AssertExpectations(t)
		})
	}
}

// TestAlertaService_Listar verifica que el servicio delega correctamente al repositorio.
func TestAlertaService_Listar(t *testing.T) {
	repo := new(alertaRepoMock)
	esperado := []models.AlertaEmergencia{
		{ID: 1, Gravedad: 3, Requerimiento: "Ecógrafo", Estado: "Buscando"},
		{ID: 2, Gravedad: 5, Requerimiento: "Ventilador", Estado: "Asignada"},
	}
	repo.On("ListarAlertas").Return(esperado)

	svc := service.NuevoAlertaService(repo)
	resultado := svc.Listar()

	assert.Equal(t, 2, len(resultado))
	repo.AssertExpectations(t)
}

// TestAlertaService_Obtener cubre el camino feliz (existe) y el triste (no existe).
func TestAlertaService_Obtener(t *testing.T) {
	t.Run("alerta existente -> devuelve alerta sin error", func(t *testing.T) {
		repo := new(alertaRepoMock)
		alerta := models.AlertaEmergencia{ID: 1, Gravedad: 3, Requerimiento: "Ecógrafo", Estado: "Buscando"}
		repo.On("BuscarAlertaPorID", 1).Return(alerta, true)

		svc := service.NuevoAlertaService(repo)
		resultado, err := svc.Obtener(1)

		require.NoError(t, err)
		assert.Equal(t, alerta, resultado)
		repo.AssertExpectations(t)
	})
	t.Run("alerta inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(alertaRepoMock)
		repo.On("BuscarAlertaPorID", 99).Return(models.AlertaEmergencia{}, false)

		svc := service.NuevoAlertaService(repo)
		_, err := svc.Obtener(99)

		require.ErrorIs(t, err, service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}

// TestAlertaService_Actualizar cubre actualización exitosa y registro inexistente.
func TestAlertaService_Actualizar(t *testing.T) {
	t.Run("datos validos y alerta existente -> devuelve actualizada", func(t *testing.T) {
		repo := new(alertaRepoMock)
		datos := models.AlertaEmergencia{Gravedad: 4, Requerimiento: "Ventilador", Estado: "Atendido"}
		actualizada := datos
		actualizada.ID = 1
		repo.On("ActualizarAlerta", 1, datos).Return(actualizada, true)

		svc := service.NuevoAlertaService(repo)
		resultado, err := svc.Actualizar(1, datos)

		require.NoError(t, err)
		assert.Equal(t, 1, resultado.ID)
		repo.AssertExpectations(t)
	})
	t.Run("alerta inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(alertaRepoMock)
		datos := models.AlertaEmergencia{Gravedad: 2, Requerimiento: "Ecógrafo", Estado: "Buscando"}
		repo.On("ActualizarAlerta", 99, datos).Return(models.AlertaEmergencia{}, false)

		svc := service.NuevoAlertaService(repo)
		_, err := svc.Actualizar(99, datos)

		require.ErrorIs(t, err, service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}

// TestAlertaService_Borrar cubre eliminación exitosa y registro inexistente.
func TestAlertaService_Borrar(t *testing.T) {
	t.Run("alerta existente -> nil", func(t *testing.T) {
		repo := new(alertaRepoMock)
		repo.On("BorrarAlerta", 1).Return(true)

		svc := service.NuevoAlertaService(repo)
		require.NoError(t, svc.Borrar(1))
		repo.AssertExpectations(t)
	})
	t.Run("alerta inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(alertaRepoMock)
		repo.On("BorrarAlerta", 99).Return(false)

		svc := service.NuevoAlertaService(repo)
		require.ErrorIs(t, svc.Borrar(99), service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}
