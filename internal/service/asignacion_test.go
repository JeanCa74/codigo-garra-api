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
