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
