package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/JeanCa74/codigo-garra-api/internal/handlers"
	"github.com/JeanCa74/codigo-garra-api/internal/middleware"
	"github.com/JeanCa74/codigo-garra-api/internal/service"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

func main() {
	// 1. GORM abre la base, ejecuta AutoMigrate y devuelve el almacén.
	almacen, err := storage.NuevoAlmacenGORM("codigogarra.db")
	if err != nil {
		log.Fatal("no se pudo abrir la base de datos: ", err)
	}

	// 2. Capa de servicio con inyección de dependencias.
	alertaSvc := service.NuevoAlertaService(almacen)
	recursoSvc := service.NuevoRecursoService(almacen)
	asignacionSvc := service.NuevoAsignacionService(almacen)
	authSvc := service.NuevoAuthService(almacen)

	// 3. Server con los servicios inyectados.
	servidor := handlers.NewServer(alertaSvc, recursoSvc, asignacionSvc, authSvc)

	// 4. Router + middleware global.
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// 5. Rutas versionadas /api/v1/.
	r.Route("/api/v1", func(r chi.Router) {
		// Públicas: registro y login.
		r.Post("/auth/register", servidor.Registrar)
		r.Post("/auth/login", servidor.Login)

		// Protegidas: exigen JWT válido en Authorization: Bearer <token>.
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))

			// Módulo Jean Carlos — Alertas de emergencia (Triage)
			r.Get("/alertas", servidor.ListarAlertas)
			r.Post("/alertas", servidor.CrearAlerta)
			r.Get("/alertas/{id}", servidor.ObtenerAlerta)
			r.Put("/alertas/{id}", servidor.ActualizarAlerta)
			r.Delete("/alertas/{id}", servidor.BorrarAlerta)

			// Módulo John Erick — Recursos clínicos (Capacidad Dinámica)
			r.Get("/recursos", servidor.ListarRecursos)
			r.Post("/recursos", servidor.CrearRecurso)
			r.Get("/recursos/{id}", servidor.ObtenerRecurso)
			r.Put("/recursos/{id}", servidor.ActualizarRecurso)
			r.Delete("/recursos/{id}", servidor.BorrarRecurso)

			// Módulo María José — Asignaciones de triage (Enrutamiento y Matching)
			r.Get("/asignaciones", servidor.ListarAsignaciones)
			r.Post("/asignaciones", servidor.CrearAsignacion)
			r.Get("/asignaciones/{id}", servidor.ObtenerAsignacion)
			r.Put("/asignaciones/{id}", servidor.ActualizarAsignacion)
			r.Delete("/asignaciones/{id}", servidor.BorrarAsignacion)
		})
	})

	log.Println("Código Garra API corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
