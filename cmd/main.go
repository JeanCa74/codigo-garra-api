package main

import (
	"fmt"
	"net/http"

	"github.com/JeanCa74/codigo-garra-api/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Módulo 1: Alertas (Jean Carlos)
	r.Route("/api/v1/alertas", func(router chi.Router) {
		router.Post("/", handlers.CreateAlerta)
		router.Get("/", handlers.GetAlertas)
		router.Get("/{id}", handlers.GetAlerta)
		router.Put("/{id}", handlers.UpdateAlerta)
		router.Delete("/{id}", handlers.DeleteAlerta)
	})

	// Módulo 3: Asignaciones (María José)
	r.Route("/api/v1/asignaciones", func(router chi.Router) {
		router.Post("/", handlers.CreateAsignacion)
		router.Get("/", handlers.GetAsignaciones)
		router.Get("/{id}", handlers.GetAsignacion)
		router.Put("/{id}", handlers.UpdateAsignacion)
		router.Delete("/{id}", handlers.DeleteAsignacion)
	})

	fmt.Println("Servidor de Código Garra API corriendo en http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Error al levantar el servidor: %v\n", err)
	}
}
