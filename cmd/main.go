package main

import (
	"log"
	"net/http"
	"os"
	_ "time/tzdata" // embeds timezone DB in binary; required for Alpine (no tzdata package)

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/JeanCa74/codigo-garra-api/internal/config"
	"github.com/JeanCa74/codigo-garra-api/internal/handlers"
	"github.com/JeanCa74/codigo-garra-api/internal/middleware"
	"github.com/JeanCa74/codigo-garra-api/internal/seed"
	"github.com/JeanCa74/codigo-garra-api/internal/service"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

func main() {
	// 1. Configuración desde variables de entorno / .env
	cfg := config.Load()

	// 2. Almacén: PostgreSQL si DATABASE_URL está definido, SQLite en caso contrario.
	var almacen *storage.AlmacenGORM
	var err error
	if cfg.DatabaseURL != "" {
		almacen, err = storage.NuevoAlmacenPostgres(cfg.DatabaseURL)
	} else {
		almacen, err = storage.NuevoAlmacenGORM(cfg.DBPath)
	}
	if err != nil {
		log.Fatal("no se pudo abrir la base de datos: ", err)
	}

	// 3. Capa de servicio — un servicio por entidad del dominio.
	deps := handlers.ServerDeps{
		Alertas:      service.NuevoAlertaService(almacen),
		Asignaciones: service.NuevoAsignacionService(almacen),
		Perfiles:     service.NuevoPerfilVeterinarioService(almacen),
		Recursos:     service.NuevoRecursoService(almacen),
		Mascotas:     service.NuevoMascotaService(almacen),
		HistorialMed: service.NuevoHistorialMedicoService(almacen),
		Auth:         service.NuevoAuthService(almacen),
	}

	// 4. Seeder automático al primer arranque (o cuando SEED_ON_STARTUP=true).
	if os.Getenv("SEED_ON_STARTUP") == "true" {
		if err := seed.Sembrar(almacen, almacen, deps.Auth); err != nil {
			log.Println("aviso seeder:", err)
		}
	}

	// 5. Server con inyección de dependencias.
	servidor := handlers.NewServer(deps)

	// 6. Router + middleware global.
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// 7. Rutas versionadas /api/v1/.
	r.Route("/api/v1", func(r chi.Router) {
		// Públicas: registro y login.
		r.Post("/auth/register", servidor.Registrar)
		r.Post("/auth/login", servidor.Login)

		// Protegidas con JWT válido.
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(deps.Auth))

			// ── Módulo 1: Emergencia médica (Jean Carlos) ──────────────────
			r.Get("/alertas", servidor.ListarAlertas)
			r.Post("/alertas", servidor.CrearAlerta)
			r.Get("/alertas/{id}", servidor.ObtenerAlerta)
			r.Put("/alertas/{id}", servidor.ActualizarAlerta)

			r.Get("/alertas/{id}/asignaciones", servidor.ListarAsignacionesDeAlerta)

			r.Get("/asignaciones", servidor.ListarAsignaciones)
			r.Post("/asignaciones", servidor.CrearAsignacion)
			r.Get("/asignaciones/{id}", servidor.ObtenerAsignacion)
			r.Put("/asignaciones/{id}", servidor.ActualizarAsignacion)

			// ── Módulo 2: Perfil veterinario (John Erick) ──────────────────
			r.Get("/veterinarios", servidor.ListarPerfiles)
			r.Post("/veterinarios", servidor.CrearPerfil)
			r.Get("/veterinarios/{id}", servidor.ObtenerPerfil)
			r.Put("/veterinarios/{id}", servidor.ActualizarPerfil)
			r.Get("/veterinarios/{id}/recursos", servidor.ListarRecursosDePerfil)

			r.Get("/recursos", servidor.ListarRecursos)
			r.Post("/recursos", servidor.CrearRecurso)
			r.Get("/recursos/{id}", servidor.ObtenerRecurso)
			r.Put("/recursos/{id}", servidor.ActualizarRecurso)

			// ── Módulo 3: Historial médico (María José) ────────────────────
			r.Get("/mascotas", servidor.ListarMascotas)
			r.Post("/mascotas", servidor.CrearMascota)
			r.Get("/mascotas/{id}", servidor.ObtenerMascota)
			r.Put("/mascotas/{id}", servidor.ActualizarMascota)
			r.Get("/mascotas/{id}/historial", servidor.ListarHistorialDeMascota)

			r.Get("/historial", servidor.ListarHistorial)
			r.Post("/historial", servidor.CrearHistorial)
			r.Get("/historial/{id}", servidor.ObtenerHistorial)

			// ── Solo admins: eliminación de registros ──────────────────────
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRol("admin"))
				r.Delete("/alertas/{id}", servidor.BorrarAlerta)
				r.Delete("/asignaciones/{id}", servidor.BorrarAsignacion)
				r.Delete("/veterinarios/{id}", servidor.BorrarPerfil)
				r.Delete("/recursos/{id}", servidor.BorrarRecurso)
				r.Delete("/mascotas/{id}", servidor.BorrarMascota)
			})
		})
	})

	log.Printf("Código Garra API corriendo en http://localhost%s", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, r))
}
