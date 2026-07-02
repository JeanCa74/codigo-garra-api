package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// TestGORM_CrearYBuscarPerfil verifica que el repositorio real (GORM + SQLite en
// memoria) persiste un PerfilVeterinario y lo recupera fielmente por ID.
func TestGORM_CrearYBuscarPerfil(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err, "debe abrir la base en memoria sin error")

	entrada := models.PerfilVeterinario{
		Nombre:    "Clínica Garra Norte",
		Telefono:  "0987654321",
		Direccion: "Av. Principal 123",
		Activo:    true,
	}

	creado := almacen.CrearPerfil(entrada)

	require.NotZero(t, creado.ID, "GORM debe asignar un ID autoincremental")

	encontrado, ok := almacen.BuscarPerfilPorID(creado.ID)
	require.True(t, ok, "el perfil recién creado debe encontrarse por su ID")
	assert.Equal(t, "Clínica Garra Norte", encontrado.Nombre)
	assert.Equal(t, "0987654321", encontrado.Telefono)
	assert.True(t, encontrado.Activo)
}

// TestGORM_ListarPerfilesReflejaCreacion verifica que ListarPerfiles devuelve los registros insertados.
func TestGORM_ListarPerfilesReflejaCreacion(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err)

	almacen.CrearPerfil(models.PerfilVeterinario{Nombre: "Clínica A", Telefono: "111"})
	almacen.CrearPerfil(models.PerfilVeterinario{Nombre: "Clínica B", Telefono: "222"})

	lista := almacen.ListarPerfiles()
	assert.Len(t, lista, 2, "deben listarse los dos perfiles creados")
}
