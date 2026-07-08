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

// construirEntornoRecursos arma el router con rutas de Recursos y middleware Auth real.
func construirEntornoRecursos(t *testing.T) (http.Handler, string) {
	t.Helper()
	almacen := storage.NuevaMemoria()
	usuarios := nuevoUsuarioRepoFake()

	recursoSvc := service.NuevoRecursoService(almacen)
	authSvc := service.NuevoAuthService(usuarios)

	// --- CORRECCIÓN: Inyección mediante ServerDeps ---
	deps := handlers.ServerDeps{
		Recursos: recursoSvc,
		Auth:     authSvc,
	}
	srv := handlers.NewServer(deps)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Get("/recursos", srv.ListarRecursos)
			r.Post("/recursos", srv.CrearRecurso)
			r.Get("/recursos/{id}", srv.ObtenerRecurso)
			r.Put("/recursos/{id}", srv.ActualizarRecurso)
			r.Delete("/recursos/{id}", srv.BorrarRecurso)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

// TestCrearRecurso_Exitoso: POST con token y tipo_maquina válido -> 201 Created.
func TestCrearRecurso_Exitoso(t *testing.T) {
	h, token := construirEntornoRecursos(t)
	body := `{"perfil_id":5,"tipo_maquina":"Ecógrafo","esta_disponible":true}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/recursos", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var creado models.RecursoClinico
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&creado))
	assert.NotZero(t, creado.ID)
	assert.Equal(t, "Ecógrafo", creado.TipoMaquina)
}

// TestCrearRecurso_TipoMaquinaVacio: tipo_maquina vacío viola la regla de negocio -> 400.
func TestCrearRecurso_TipoMaquinaVacio(t *testing.T) {
	h, token := construirEntornoRecursos(t)
	body := `{"perfil_id":2,"tipo_maquina":"","esta_disponible":true}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/recursos", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestRutaRecursos_SinToken: sin Authorization el middleware devuelve 401.
func TestRutaRecursos_SinToken(t *testing.T) {
	h, _ := construirEntornoRecursos(t)
	body := `{"perfil_id":1,"tipo_maquina":"Ventilador","esta_disponible":true}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/recursos", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// TestListarRecursos_Exitoso: GET /recursos con token devuelve 200.
func TestListarRecursos_Exitoso(t *testing.T) {
	h, token := construirEntornoRecursos(t)
	bodyCreate := `{"perfil_id":1,"tipo_maquina":"Ecógrafo","esta_disponible":true}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/recursos", strings.NewReader(bodyCreate))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(httptest.NewRecorder(), reqCreate)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/recursos", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var recursos []models.RecursoClinico
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&recursos))
	assert.GreaterOrEqual(t, len(recursos), 1)
}

// TestObtenerRecurso_PorID: crea un recurso y lo obtiene por ID → 200.
func TestObtenerRecurso_PorID(t *testing.T) {
	h, token := construirEntornoRecursos(t)
	bodyCreate := `{"perfil_id":2,"tipo_maquina":"Desfibrilador","esta_disponible":false}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/recursos", strings.NewReader(bodyCreate))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(httptest.NewRecorder(), reqCreate)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/recursos/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var recurso models.RecursoClinico
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&recurso))
	assert.Equal(t, "Desfibrilador", recurso.TipoMaquina)
}

// TestActualizarRecurso_Exitoso: crea un recurso y lo actualiza → 200.
func TestActualizarRecurso_Exitoso(t *testing.T) {
	h, token := construirEntornoRecursos(t)
	bodyCreate := `{"perfil_id":1,"tipo_maquina":"Ecógrafo","esta_disponible":true}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/recursos", strings.NewReader(bodyCreate))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(httptest.NewRecorder(), reqCreate)

	bodyUpdate := `{"perfil_id":1,"tipo_maquina":"Ecógrafo 4K","esta_disponible":false}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/recursos/1", strings.NewReader(bodyUpdate))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actualizado models.RecursoClinico
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&actualizado))
	assert.Equal(t, "Ecógrafo 4K", actualizado.TipoMaquina)
}
