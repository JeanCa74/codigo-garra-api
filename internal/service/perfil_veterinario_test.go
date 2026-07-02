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
