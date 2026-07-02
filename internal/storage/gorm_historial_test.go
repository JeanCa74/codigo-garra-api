package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// TestGORM_CrearYBuscarHistorial verifica que el repositorio real (GORM + SQLite en
// memoria) persiste un HistorialMedico y lo recupera fielmente por ID.
func TestGORM_CrearYBuscarHistorial(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err, "debe abrir la base en memoria sin error")

	entrada := models.HistorialMedico{
		MascotaID:   3,
		Diagnostico: "Otitis externa",
		Tratamiento: "Antibiótico tópico",
		Fecha:       "2026-07-01",
		Veterinario: "Dra. López",
	}

	creado := almacen.CrearHistorial(entrada)

	require.NotZero(t, creado.ID, "GORM debe asignar un ID autoincremental")

	encontrado, ok := almacen.BuscarHistorialPorID(creado.ID)
	require.True(t, ok, "el historial recién creado debe encontrarse por su ID")
	assert.Equal(t, 3, encontrado.MascotaID)
	assert.Equal(t, "Otitis externa", encontrado.Diagnostico)
	assert.Equal(t, "2026-07-01", encontrado.Fecha)
}

// TestGORM_ListarHistorialPorMascota verifica el filtro por mascota_id.
func TestGORM_ListarHistorialPorMascota(t *testing.T) {
	almacen, err := storage.NuevoAlmacenGORM(":memory:")
	require.NoError(t, err)

	almacen.CrearHistorial(models.HistorialMedico{MascotaID: 1, Diagnostico: "Fiebre", Fecha: "2026-06-01"})
	almacen.CrearHistorial(models.HistorialMedico{MascotaID: 2, Diagnostico: "Fractura", Fecha: "2026-06-15"})
	almacen.CrearHistorial(models.HistorialMedico{MascotaID: 1, Diagnostico: "Seguimiento", Fecha: "2026-07-01"})

	deMascota1 := almacen.ListarHistorialPorMascota(1)
	assert.Len(t, deMascota1, 2, "la mascota 1 tiene dos entradas de historial")

	deMascota2 := almacen.ListarHistorialPorMascota(2)
	assert.Len(t, deMascota2, 1, "la mascota 2 tiene una sola entrada")
}
