package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JeanCa74/codigo-garra-api/internal/handlers"
	"github.com/JeanCa74/codigo-garra-api/internal/middleware"
	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/service"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// construirEntornoMascotas arma el router con rutas de Mascotas y middleware Auth real.
func construirEntornoMascotas(t *testing.T) (http.Handler, string) {
	t.Helper()
	almacen := storage.NuevaMemoria()
	usuarios := nuevoUsuarioRepoFake()

	mascotaSvc := service.NuevoMascotaService(almacen)
	historialSvc := service.NuevoHistorialMedicoService(almacen)
	authSvc := service.NuevoAuthService(usuarios)

	deps := handlers.ServerDeps{
		Mascotas:     mascotaSvc,
		HistorialMed: historialSvc,
		Auth:         authSvc,
	}
	srv := handlers.NewServer(deps)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Get("/mascotas", srv.ListarMascotas)
			r.Post("/mascotas", srv.CrearMascota)
			r.Get("/mascotas/{id}", srv.ObtenerMascota)
			r.Put("/mascotas/{id}", srv.ActualizarMascota)
			r.Delete("/mascotas/{id}", srv.BorrarMascota)
			r.Get("/mascotas/{id}/historial", srv.ListarHistorialDeMascota)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

// TestCrearMascota_Exitosa: POST con token y nombre válido -> 201 Created.
func TestCrearMascota_Exitosa(t *testing.T) {
	h, token := construirEntornoMascotas(t)
	body := `{"nombre":"Max","especie":"Perro","edad":3,"dueno":"Ana García"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/mascotas", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusTeapot, rec.Code) // SABOTAJE CP3 María José
	var creada models.Mascota
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&creada))
	assert.NotZero(t, creada.ID)
	assert.Equal(t, "Max", creada.Nombre)
}

// TestCrearMascota_NombreVacio: nombre vacío viola la regla de negocio -> 400.
func TestCrearMascota_NombreVacio(t *testing.T) {
	h, token := construirEntornoMascotas(t)
	body := `{"nombre":"","especie":"Gato","edad":2,"dueno":"Carlos"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/mascotas", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestRutaMascotas_SinToken: sin Authorization el middleware devuelve 401.
func TestRutaMascotas_SinToken(t *testing.T) {
	h, _ := construirEntornoMascotas(t)
	body := `{"nombre":"Luna","especie":"Gato","dueno":"Pedro"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/mascotas", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// TestListarMascotas_Exitoso: GET /mascotas devuelve 200 con el array de mascotas.
func TestListarMascotas_Exitoso(t *testing.T) {
	h, token := construirEntornoMascotas(t)
	bodyCreate := `{"nombre":"Luna","especie":"Gato","edad":5,"dueno":"Carlos"}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/mascotas", strings.NewReader(bodyCreate))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(httptest.NewRecorder(), reqCreate)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/mascotas", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var mascotas []models.Mascota
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&mascotas))
	assert.GreaterOrEqual(t, len(mascotas), 1)
}

// TestObtenerMascota_PorID: crea una mascota y la obtiene por ID → 200.
func TestObtenerMascota_PorID(t *testing.T) {
	h, token := construirEntornoMascotas(t)
	bodyCreate := `{"nombre":"Simba","especie":"León","edad":1,"dueno":"Rey"}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/mascotas", strings.NewReader(bodyCreate))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(httptest.NewRecorder(), reqCreate)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/mascotas/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var mascota models.Mascota
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&mascota))
	assert.Equal(t, "Simba", mascota.Nombre)
}

// TestActualizarMascota_Exitosa: crea una mascota y la actualiza → 200.
func TestActualizarMascota_Exitosa(t *testing.T) {
	h, token := construirEntornoMascotas(t)
	bodyCreate := `{"nombre":"Rocky","especie":"Perro","edad":2,"dueno":"Pedro"}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/mascotas", strings.NewReader(bodyCreate))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(httptest.NewRecorder(), reqCreate)

	bodyUpdate := `{"nombre":"Rocky Junior","especie":"Perro","edad":3,"dueno":"Pedro"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/mascotas/1", strings.NewReader(bodyUpdate))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actualizada models.Mascota
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&actualizada))
	assert.Equal(t, "Rocky Junior", actualizada.Nombre)
}

// TestListarHistorialDeMascota: crea mascota + historial y usa la ruta anidada → 200.
func TestListarHistorialDeMascota(t *testing.T) {
	h, token := construirEntornoMascotas(t)

	// Crear mascota primero (ID = 1 en memoria).
	bodyMascota := `{"nombre":"Rocky","especie":"Perro","edad":2,"dueno":"Pedro"}`
	reqMascota := httptest.NewRequest(http.MethodPost, "/api/v1/mascotas", strings.NewReader(bodyMascota))
	reqMascota.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(httptest.NewRecorder(), reqMascota)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/mascotas/1/historial", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	// Historial vacío es válido; solo verificamos que la ruta responde con 200.
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
}
