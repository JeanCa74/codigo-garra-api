package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
)

var asignacionesDB []models.AsignacionTriage
var asignacionID int = 0

func CreateAsignacion(w http.ResponseWriter, r *http.Request) {
	var nueva models.AsignacionTriage
	if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	asignacionID++
	nueva.ID = asignacionID
	nueva.EstadoConfirmacion = "Pendiente"
	asignacionesDB = append(asignacionesDB, nueva)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nueva)
}

func GetAsignaciones(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(asignacionesDB)
}
