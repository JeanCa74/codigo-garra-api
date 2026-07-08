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

// usuarioRepoFake es un fake en memoria para UserRepository (no un mock).
type usuarioRepoFake struct {
	porEmail map[string]models.Usuario
	nextID   int
}

func nuevoUsuarioRepoFake() *usuarioRepoFake {
	return &usuarioRepoFake{porEmail: map[string]models.Usuario{}, nextID: 1}
}

func (f *usuarioRepoFake) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	u.ID = f.nextID
	f.nextID++
	f.porEmail[u.Email] = u
	return u, nil
}

func (f *usuarioRepoFake) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	u, ok := f.porEmail[email]
	return u, ok
}

// registrarYObtenerToken hace register + login contra el router para obtener un JWT real.
func registrarYObtenerToken(t *testing.T, h http.Handler) string {
	t.Helper()
	cred := `{"email":"clinica@codigogarra.vet","password":"clave123"}`

	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(cred))
	h.ServeHTTP(httptest.NewRecorder(), reqReg)

	reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(cred))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, reqLogin)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Token string `json:"token"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.NotEmpty(t, resp.Token)
	return resp.Token
}

// construirEntornoAlertas arma el router con rutas de Alertas y middleware Auth real.
func construirEntornoAlertas(t *testing.T) (http.Handler, string) {
	t.Helper()
	almacen := storage.NuevaMemoria()
	almacen.SeedAlertas()
	usuarios := nuevoUsuarioRepoFake()

	alertaSvc := service.NuevoAlertaService(almacen)
	asignacionSvc := service.NuevoAsignacionService(almacen)
	authSvc := service.NuevoAuthService(usuarios)

	deps := handlers.ServerDeps{
		Alertas:      alertaSvc,
		Asignaciones: asignacionSvc,
		Auth:         authSvc,
	}
	srv := handlers.NewServer(deps)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Get("/alertas", srv.ListarAlertas)
			r.Post("/alertas", srv.CrearAlerta)
			r.Get("/alertas/{id}", srv.ObtenerAlerta)
			r.Put("/alertas/{id}", srv.ActualizarAlerta)
			r.Delete("/alertas/{id}", srv.BorrarAlerta)
			r.Get("/alertas/{id}/asignaciones", srv.ListarAsignacionesDeAlerta)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

// TestCrearAlerta_Exitosa: POST con token y gravedad válida -> 201 Created.
func TestCrearAlerta_Exitosa(t *testing.T) {
	h, token := construirEntornoAlertas(t)
	body := `{"gravedad":4,"requerimiento":"Ventilador mecánico","estado":"Buscando"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/alertas", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusTeapot, rec.Code) // SABOTAJE CP3 Jean Carlos
	var creada models.AlertaEmergencia
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&creada))
	assert.NotZero(t, creada.ID)
	assert.Equal(t, "Ventilador mecánico", creada.Requerimiento)
}

// TestCrearAlerta_GravedadInvalida: gravedad=0 viola la escala de triage -> 400.
func TestCrearAlerta_GravedadInvalida(t *testing.T) {
	h, token := construirEntornoAlertas(t)
	body := `{"gravedad":0,"requerimiento":"Ecógrafo"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/alertas", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestRutaAlertas_SinToken: sin Authorization el middleware devuelve 401.
func TestRutaAlertas_SinToken(t *testing.T) {
	h, _ := construirEntornoAlertas(t)
	body := `{"gravedad":3,"requerimiento":"Radiografía"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/alertas", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// TestListarAlertas_Exitoso: GET /alertas con token válido devuelve 200 y el array de alertas.
func TestListarAlertas_Exitoso(t *testing.T) {
	h, token := construirEntornoAlertas(t) // SeedAlertas carga 2 alertas

	req := httptest.NewRequest(http.MethodGet, "/api/v1/alertas", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var alertas []models.AlertaEmergencia
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&alertas))
	assert.GreaterOrEqual(t, len(alertas), 2)
}

// TestObtenerAlerta_Exitoso: GET /alertas/1 con token devuelve 200 y la alerta correcta.
func TestObtenerAlerta_Exitoso(t *testing.T) {
	h, token := construirEntornoAlertas(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/alertas/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var alerta models.AlertaEmergencia
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&alerta))
	assert.Equal(t, 1, alerta.ID)
}

// TestObtenerAlerta_NoExiste: GET /alertas/999 devuelve 404.
func TestObtenerAlerta_NoExiste(t *testing.T) {
	h, token := construirEntornoAlertas(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/alertas/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestListarAsignacionesDeAlerta: GET /alertas/1/asignaciones devuelve 200 con array.
func TestListarAsignacionesDeAlerta(t *testing.T) {
	h, token := construirEntornoAlertas(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/alertas/1/asignaciones", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	var asignaciones []models.AsignacionTriage
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&asignaciones))
	assert.NotNil(t, asignaciones)
}

// TestActualizarAlerta_Exitosa: PUT /alertas/1 con datos válidos devuelve 200.
func TestActualizarAlerta_Exitosa(t *testing.T) {
	h, token := construirEntornoAlertas(t)
	body := `{"gravedad":5,"requerimiento":"Ventilador mecánico crítico","estado":"Atendido"}`

	req := httptest.NewRequest(http.MethodPut, "/api/v1/alertas/1", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var actualizada models.AlertaEmergencia
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&actualizada))
	assert.Equal(t, "Atendido", actualizada.Estado)
}
