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

// construirEntornoAsignaciones arma el router con rutas de Asignaciones y middleware Auth real.
func construirEntornoAsignaciones(t *testing.T) (http.Handler, string) {
	t.Helper()
	almacen := storage.NuevaMemoria()
	usuarios := nuevoUsuarioRepoFake()

	asignacionSvc := service.NuevoAsignacionService(almacen)
	authSvc := service.NuevoAuthService(usuarios)
	srv := handlers.NewServer(nil, asignacionSvc, nil, nil, nil, nil, authSvc)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Get("/asignaciones", srv.ListarAsignaciones)
			r.Post("/asignaciones", srv.CrearAsignacion)
			r.Get("/asignaciones/{id}", srv.ObtenerAsignacion)
			r.Put("/asignaciones/{id}", srv.ActualizarAsignacion)
			r.Delete("/asignaciones/{id}", srv.BorrarAsignacion)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

// TestCrearAsignacion_Exitosa: POST con IDs válidos y token -> 201 Created.
func TestCrearAsignacion_Exitosa(t *testing.T) {
	h, token := construirEntornoAsignaciones(t)
	body := `{"alerta_id":1,"recurso_id":2,"estado_confirmacion":"Pendiente"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/asignaciones", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var creada models.AsignacionTriage
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&creada))
	assert.NotZero(t, creada.ID)
	assert.Equal(t, 1, creada.AlertaID)
	assert.Equal(t, 2, creada.RecursoID)
}

// TestCrearAsignacion_IDsInvalidos: alerta_id=0 viola la regla de matching -> 400.
func TestCrearAsignacion_IDsInvalidos(t *testing.T) {
	h, token := construirEntornoAsignaciones(t)
	body := `{"alerta_id":0,"recurso_id":2}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/asignaciones", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestRutaAsignaciones_SinToken: sin Authorization el middleware devuelve 401.
//
// Qué se rompería: si se elimina r.Use(middleware.Auth(...)), la petición
// llegaría al handler y respondería 201 en lugar de 401.
func TestRutaAsignaciones_SinToken(t *testing.T) {
	h, _ := construirEntornoAsignaciones(t)
	body := `{"alerta_id":1,"recurso_id":3}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/asignaciones", strings.NewReader(body))
	// A propósito: NO establecemos el header Authorization.
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
