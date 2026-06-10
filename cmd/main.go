package main

import (
	"fmt"
	"github.com/JeanCa74/codigo-garra-api/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
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

	fmt.Println("Servidor de Código Garra API corriendo en http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Error al levantar el servidor: %v\n", err)
	}
}
