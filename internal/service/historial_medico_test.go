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

// historialRepoMock registra llamadas a HistorialMedicoRepository.
type historialRepoMock struct {
	mock.Mock
}

func (m *historialRepoMock) ListarHistorial() []models.HistorialMedico {
	return m.Called().Get(0).([]models.HistorialMedico)
}
func (m *historialRepoMock) BuscarHistorialPorID(id int) (models.HistorialMedico, bool) {
	args := m.Called(id)
	return args.Get(0).(models.HistorialMedico), args.Bool(1)
}
func (m *historialRepoMock) CrearHistorial(h models.HistorialMedico) models.HistorialMedico {
	return m.Called(h).Get(0).(models.HistorialMedico)
}
func (m *historialRepoMock) ListarHistorialPorMascota(mascotaID int) []models.HistorialMedico {
	return m.Called(mascotaID).Get(0).([]models.HistorialMedico)
}

var _ storage.HistorialMedicoRepository = (*historialRepoMock)(nil)

// TestHistorialMedicoService_Crear verifica que validarHistorial rechaza
// mascota_id inválido y diagnóstico vacío sin tocar el repositorio.
func TestHistorialMedicoService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.HistorialMedico
		errEsperado   error
		debePersistir bool
	}{
		{
			"mascota_id cero -> ErrMascotaIDInvalido y NO llega al repo",
			models.HistorialMedico{MascotaID: 0, Diagnostico: "Infección ocular"},
			service.ErrMascotaIDInvalido,
			false,
		},
		{
			"mascota_id negativo -> ErrMascotaIDInvalido y NO llega al repo",
			models.HistorialMedico{MascotaID: -1, Diagnostico: "Fractura"},
			service.ErrMascotaIDInvalido,
			false,
		},
		{
			"diagnostico vacio -> ErrDiagnosticoVacio y NO llega al repo",
			models.HistorialMedico{MascotaID: 2, Diagnostico: ""},
			service.ErrDiagnosticoVacio,
			false,
		},
		{
			"diagnostico solo espacios -> ErrDiagnosticoVacio y NO llega al repo",
			models.HistorialMedico{MascotaID: 3, Diagnostico: "   "},
			service.ErrDiagnosticoVacio,
			false,
		},
		{
			"historial valido -> sin error y se persiste",
			models.HistorialMedico{MascotaID: 1, Diagnostico: "Gastritis leve", Fecha: "2026-07-01"},
			nil,
			true,
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := new(historialRepoMock)

			if tc.debePersistir {
				esperado := tc.entrada
				esperado.ID = 7
				repo.On("CrearHistorial", tc.entrada).Return(esperado)
			}

			svc := service.NuevoHistorialMedicoService(repo)
			creado, err := svc.Crear(tc.entrada)

			if tc.errEsperado != nil {
				require.ErrorIs(t, err, tc.errEsperado)
				// Clave: la validación rechaza los datos ANTES de llegar a la base de datos.
				repo.AssertNotCalled(t, "CrearHistorial")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 7, creado.ID)
				repo.AssertCalled(t, "CrearHistorial", tc.entrada)
			}

			repo.AssertExpectations(t)
		})
	}
}

// TestHistorialMedicoService_Listar verifica la delegación correcta al repositorio.
func TestHistorialMedicoService_Listar(t *testing.T) {
	repo := new(historialRepoMock)
	esperado := []models.HistorialMedico{
		{ID: 1, MascotaID: 1, Diagnostico: "Gastritis leve", Fecha: "2026-07-01"},
	}
	repo.On("ListarHistorial").Return(esperado)

	svc := service.NuevoHistorialMedicoService(repo)
	resultado := svc.Listar()

	assert.Equal(t, 1, len(resultado))
	repo.AssertExpectations(t)
}

// TestHistorialMedicoService_ObtenerPorID cubre entrada existente e inexistente.
func TestHistorialMedicoService_ObtenerPorID(t *testing.T) {
	t.Run("historial existente -> devuelve entrada sin error", func(t *testing.T) {
		repo := new(historialRepoMock)
		entrada := models.HistorialMedico{ID: 1, MascotaID: 1, Diagnostico: "Gastritis", Fecha: "2026-07-01"}
		repo.On("BuscarHistorialPorID", 1).Return(entrada, true)

		svc := service.NuevoHistorialMedicoService(repo)
		resultado, err := svc.ObtenerPorID(1)

		require.NoError(t, err)
		assert.Equal(t, "Gastritis", resultado.Diagnostico)
		repo.AssertExpectations(t)
	})
	t.Run("historial inexistente -> ErrNoEncontrado", func(t *testing.T) {
		repo := new(historialRepoMock)
		repo.On("BuscarHistorialPorID", 99).Return(models.HistorialMedico{}, false)

		svc := service.NuevoHistorialMedicoService(repo)
		_, err := svc.ObtenerPorID(99)

		require.ErrorIs(t, err, service.ErrNoEncontrado)
		repo.AssertExpectations(t)
	})
}
