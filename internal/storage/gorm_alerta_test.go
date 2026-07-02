package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// TestGORM_CrearYBuscarAlerta verifica que el repositorio real (GORM + SQLite en
// memoria) refleja la creación al buscar por ID.
func TestGORM_CrearYBuscarAlerta(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err, "debe abrir la base en memoria sin error")

	entrada := models.AlertaEmergencia{
		Gravedad:      4,
		Requerimiento: "Ventilador mecánico",
		Estado:        "Buscando",
	}

	creada := almacen.CrearAlerta(entrada)

	require.NotZero(t, creada.ID, "GORM debe asignar un ID autoincremental")

	encontrada, ok := almacen.BuscarAlertaPorID(creada.ID)
	require.True(t, ok, "la alerta recién creada debe encontrarse por su ID")
	assert.Equal(t, 4, encontrada.Gravedad)
	assert.Equal(t, "Ventilador mecánico", encontrada.Requerimiento)
	assert.Equal(t, "Buscando", encontrada.Estado)
}

// TestGORM_ListarAlertasReflejaCreacion verifica que ListarAlertas devuelve los registros insertados.
func TestGORM_ListarAlertasReflejaCreacion(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err)

	almacen.CrearAlerta(models.AlertaEmergencia{Gravedad: 3, Requerimiento: "Ecógrafo", Estado: "Buscando"})
	almacen.CrearAlerta(models.AlertaEmergencia{Gravedad: 5, Requerimiento: "Desfibrilador", Estado: "Buscando"})

	lista := almacen.ListarAlertas()
	assert.Len(t, lista, 2, "deben listarse las dos alertas creadas")
}
