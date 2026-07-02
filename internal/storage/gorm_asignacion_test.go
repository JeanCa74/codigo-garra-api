package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// TestGORM_CrearYBuscarAsignacion verifica que el repositorio real (GORM + SQLite en
// memoria) refleja la creación al buscar por ID.
func TestGORM_CrearYBuscarAsignacion(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err, "debe abrir la base en memoria sin error")

	entrada := models.AsignacionTriage{
		AlertaID:           3,
		RecursoID:          5,
		EstadoConfirmacion: "Pendiente",
	}

	creada := almacen.CrearAsignacion(entrada)

	require.NotZero(t, creada.ID, "GORM debe asignar un ID autoincremental")

	encontrada, ok := almacen.BuscarAsignacionPorID(creada.ID)
	require.True(t, ok, "la asignación recién creada debe encontrarse por su ID")
	assert.Equal(t, 3, encontrada.AlertaID)
	assert.Equal(t, 5, encontrada.RecursoID)
	assert.Equal(t, "Pendiente", encontrada.EstadoConfirmacion)
}

// TestGORM_ListarAsignacionesReflejaCreacion verifica que ListarAsignaciones devuelve los registros insertados.
func TestGORM_ListarAsignacionesReflejaCreacion(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err)

	almacen.CrearAsignacion(models.AsignacionTriage{AlertaID: 1, RecursoID: 2, EstadoConfirmacion: "Pendiente"})
	almacen.CrearAsignacion(models.AsignacionTriage{AlertaID: 2, RecursoID: 3, EstadoConfirmacion: "Confirmada"})

	lista := almacen.ListarAsignaciones()
	assert.Len(t, lista, 2, "deben listarse las dos asignaciones creadas")
}
