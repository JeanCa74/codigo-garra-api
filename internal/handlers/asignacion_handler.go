package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/go-chi/chi/v5"
)

var asignacionesDB []models.AsignacionTriage
var asignacionID int = 0

// CreateAsignacion empareja y confirma el bloqueo temporal de un recurso

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

func GetAsignacion(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	for _, asig := range asignacionesDB {
		if asig.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(asig)
			return
		}
	}
	http.Error(w, "Match de emergencia no encontrado", http.StatusNotFound)
}

func UpdateAsignacion(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var actualizacion models.AsignacionTriage
	json.NewDecoder(r.Body).Decode(&actualizacion)
	for i, asig := range asignacionesDB {
		if asig.ID == id {
			asignacionesDB[i].EstadoConfirmacion = actualizacion.EstadoConfirmacion
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(asignacionesDB[i])
			return
		}
	}
	http.Error(w, "Match no encontrado", http.StatusNotFound)
}

func DeleteAsignacion(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	for i, asig := range asignacionesDB {
		if asig.ID == id {
			asignacionesDB = append(asignacionesDB[:i], asignacionesDB[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.Error(w, "Match no encontrado", http.StatusNotFound)
}
