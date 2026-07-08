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

// construirEntornoPerfiles arma el router con rutas de PerfilVeterinario y middleware Auth real.
func construirEntornoPerfiles(t *testing.T) (http.Handler, string) {
	t.Helper()
	almacen := storage.NuevaMemoria()
	usuarios := nuevoUsuarioRepoFake()

	perfilSvc := service.NuevoPerfilVeterinarioService(almacen)
	authSvc := service.NuevoAuthService(usuarios)

	// --- CORRECCIÓN: Inyección mediante ServerDeps ---
	deps := handlers.ServerDeps{
		Perfiles: perfilSvc,
		Auth:     authSvc,
	}
	srv := handlers.NewServer(deps)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Get("/veterinarios", srv.ListarPerfiles)
			r.Post("/veterinarios", srv.CrearPerfil)
			r.Get("/veterinarios/{id}", srv.ObtenerPerfil)
			r.Put("/veterinarios/{id}", srv.ActualizarPerfil)
			r.Delete("/veterinarios/{id}", srv.BorrarPerfil)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

// TestCrearPerfil_Exitoso: POST con token, nombre y teléfono válidos -> 201 Created.
func TestCrearPerfil_Exitoso(t *testing.T) {
	h, token := construirEntornoPerfiles(t)
	body := `{"nombre":"Clínica Garra Norte","telefono":"0987654321","direccion":"Av. Principal 123","activo":true}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/veterinarios", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var creado models.PerfilVeterinario
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&creado))
	assert.NotZero(t, creado.ID)
	assert.Equal(t, "Clínica Garra Norte", creado.Nombre)
	assert.Equal(t, "0987654321", creado.Telefono)
}

// TestCrearPerfil_TelefonoVacio: teléfono vacío viola la regla de negocio -> 400.
func TestCrearPerfil_TelefonoVacio(t *testing.T) {
	h, token := construirEntornoPerfiles(t)
	body := `{"nombre":"Clínica Garra","telefono":"","activo":true}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/veterinarios", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestRutaPerfiles_SinToken: sin Authorization el middleware devuelve 401.
func TestRutaPerfiles_SinToken(t *testing.T) {
	h, _ := construirEntornoPerfiles(t)
	body := `{"nombre":"Clínica Garra","telefono":"0987654321","activo":true}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/veterinarios", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// TestListarPerfiles_Exitoso: GET /veterinarios devuelve 200 con el array de perfiles.
func TestListarPerfiles_Exitoso(t *testing.T) {
	h, token := construirEntornoPerfiles(t)
	// Crear un perfil primero para que el listado no esté vacío.
	bodyCreate := `{"nombre":"Clínica Lista","telefono":"111222333","activo":true}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/veterinarios", strings.NewReader(bodyCreate))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(httptest.NewRecorder(), reqCreate)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/veterinarios", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var perfiles []models.PerfilVeterinario
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&perfiles))
	assert.GreaterOrEqual(t, len(perfiles), 1)
}

// TestObtenerPerfil_PorID: crea un perfil y lo obtiene por ID → 200.
func TestObtenerPerfil_PorID(t *testing.T) {
	h, token := construirEntornoPerfiles(t)
	bodyCreate := `{"nombre":"Clínica Obtener","telefono":"444555666","activo":true}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/veterinarios", strings.NewReader(bodyCreate))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	recCreate := httptest.NewRecorder()
	h.ServeHTTP(recCreate, reqCreate)
	require.Equal(t, http.StatusCreated, recCreate.Code)

	// El primer perfil creado en memoria siempre tiene ID 1.
	req := httptest.NewRequest(http.MethodGet, "/api/v1/veterinarios/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var perfil models.PerfilVeterinario
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&perfil))
	assert.Equal(t, "Clínica Obtener", perfil.Nombre)
}

// TestActualizarPerfil_Exitoso: crea un perfil y lo actualiza → 200.
func TestActualizarPerfil_Exitoso(t *testing.T) {
	h, token := construirEntornoPerfiles(t)
	bodyCreate := `{"nombre":"Clínica Original","telefono":"777888999","activo":true}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/veterinarios", strings.NewReader(bodyCreate))
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(httptest.NewRecorder(), reqCreate)

	bodyUpdate := `{"nombre":"Clínica Actualizada","telefono":"777888999","activo":false}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/veterinarios/1", strings.NewReader(bodyUpdate))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actualizado models.PerfilVeterinario
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&actualizado))
	assert.Equal(t, "Clínica Actualizada", actualizado.Nombre)
}
