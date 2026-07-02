package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// TestGORM_CrearYBuscarRecurso verifica que el repositorio real (GORM + SQLite en
// memoria) refleja la creación al buscar por ID.
func TestGORM_CrearYBuscarRecurso(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err, "debe abrir la base en memoria sin error")

	entrada := models.RecursoClinico{
		PerfilID:       7,
		TipoMaquina:    "Ecógrafo portátil",
		EstaDisponible: true,
	}

	creado := almacen.CrearRecurso(entrada)

	require.NotZero(t, creado.ID, "GORM debe asignar un ID autoincremental")

	encontrado, ok := almacen.BuscarRecursoPorID(creado.ID)
	require.True(t, ok, "el recurso recién creado debe encontrarse por su ID")
	assert.Equal(t, 7, encontrado.PerfilID)
	assert.Equal(t, "Ecógrafo portátil", encontrado.TipoMaquina)
	assert.True(t, encontrado.EstaDisponible)
}

// TestGORM_ListarRecursosReflejaCreacion verifica que ListarRecursos devuelve los registros insertados.
func TestGORM_ListarRecursosReflejaCreacion(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err)

	almacen.CrearRecurso(models.RecursoClinico{PerfilID: 1, TipoMaquina: "Rayos X", EstaDisponible: true})
	almacen.CrearRecurso(models.RecursoClinico{PerfilID: 2, TipoMaquina: "Ventilador", EstaDisponible: false})

	lista := almacen.ListarRecursos()
	assert.Len(t, lista, 2, "deben listarse los dos recursos creados")
}
