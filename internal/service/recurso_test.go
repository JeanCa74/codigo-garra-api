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
			models.RecursoClinico{ClinicaID: 1, TipoMaquina: "", EstaDisponible: true},
			service.ErrTipoMaquinaVacio,
			false,
		},
		{
			"tipo_maquina solo espacios -> ErrTipoMaquinaVacio y NO llega al repo",
			models.RecursoClinico{ClinicaID: 2, TipoMaquina: "   ", EstaDisponible: true},
			service.ErrTipoMaquinaVacio,
			false,
		},
		{
			"recurso valido -> sin error y se persiste",
			models.RecursoClinico{ClinicaID: 3, TipoMaquina: "Ecógrafo", EstaDisponible: true},
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
