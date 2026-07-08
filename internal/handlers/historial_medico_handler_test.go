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

// construirEntornoHistorial arma el router con rutas de HistorialMedico y middleware Auth real.
func construirEntornoHistorial(t *testing.T) (http.Handler, string) {
	t.Helper()
	almacen := storage.NuevaMemoria()
	usuarios := nuevoUsuarioRepoFake()

	historialSvc := service.NuevoHistorialMedicoService(almacen)
	authSvc := service.NuevoAuthService(usuarios)

	// --- CORRECCIÓN: Inyección mediante ServerDeps ---
	deps := handlers.ServerDeps{
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
			r.Get("/historial", srv.ListarHistorial)
			r.Post("/historial", srv.CrearHistorial)
			r.Get("/historial/{id}", srv.ObtenerHistorial)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

// TestCrearHistorial_Exitoso: POST con token, mascota_id y diagnóstico válidos -> 201 Created.
func TestCrearHistorial_Exitoso(t *testing.T) {
	h, token := construirEntornoHistorial(t)
	body := `{"mascota_id":1,"diagnostico":"Gastritis leve","tratamiento":"Dieta blanda","fecha":"2026-07-01"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/historial", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var creado models.HistorialMedico
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&creado))
	assert.NotZero(t, creado.ID)
	assert.Equal(t, "Gastritis leve", creado.Diagnostico)
}

// TestCrearHistorial_DiagnosticoVacio: diagnóstico vacío viola la regla de negocio -> 400.
func TestCrearHistorial_DiagnosticoVacio(t *testing.T) {
	h, token := construirEntornoHistorial(t)
	body := `{"mascota_id":2,"diagnostico":"","fecha":"2026-07-01"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/historial", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestRutaHistorial_SinToken: sin Authorization el middleware devuelve 401.
func TestRutaHistorial_SinToken(t *testing.T) {
	h, _ := construirEntornoHistorial(t)
	body := `{"mascota_id":1,"diagnostico":"Infección","fecha":"2026-07-01"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/historial", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// TestListarHistorial_Exitoso: GET /historial devuelve 200 con al menos una entrada.
func TestListarHistorial_Exitoso(t *testing.T) {
	h, token := construirEntornoHistorial(t)
	bodyCreate := `{"mascota_id":1,"diagnostico":"Test listado","fecha":"2026-07-09"}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/historial", strings.NewReader(bodyCreate))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(httptest.NewRecorder(), reqCreate)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/historial", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var historial []models.HistorialMedico
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&historial))
	assert.GreaterOrEqual(t, len(historial), 1)
}

// TestObtenerHistorial_PorID: crea una entrada y la obtiene por ID → 200.
func TestObtenerHistorial_PorID(t *testing.T) {
	h, token := construirEntornoHistorial(t)
	bodyCreate := `{"mascota_id":2,"diagnostico":"Fractura leve","fecha":"2026-07-09"}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/historial", strings.NewReader(bodyCreate))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	recCreate := httptest.NewRecorder()
	h.ServeHTTP(recCreate, reqCreate)
	require.Equal(t, http.StatusCreated, recCreate.Code)

	// La primera entrada en memoria siempre tiene ID 1.
	req := httptest.NewRequest(http.MethodGet, "/api/v1/historial/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var entrada models.HistorialMedico
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&entrada))
	assert.Equal(t, "Fractura leve", entrada.Diagnostico)
}
