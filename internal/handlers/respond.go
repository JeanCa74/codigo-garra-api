package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/JeanCa74/codigo-garra-api/internal/service"
)

func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error codificando JSON: %v", err)
	}
}

func RespondError(w http.ResponseWriter, status int, mensaje string) {
	RespondJSON(w, status, map[string]string{"error": mensaje})
}

// statusDeError traduce errores de dominio a código HTTP.
func statusDeError(err error) int {
	switch {
	case errors.Is(err, service.ErrNoEncontrado):
		return http.StatusNotFound
	case errors.Is(err, service.ErrEmailEnUso):
		return http.StatusConflict
	case errors.Is(err, service.ErrCredencialesInvalidas):
		return http.StatusUnauthorized
	case errors.Is(err, service.ErrRequerimientoVacio),
		errors.Is(err, service.ErrGravedadInvalida),
		errors.Is(err, service.ErrTipoMaquinaVacio),
		errors.Is(err, service.ErrIDsInvalidos):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
