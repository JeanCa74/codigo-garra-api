package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/go-chi/chi/v5"
)

var alertasDB []models.AlertaEmergencia
var alertaID int = 0

func CreateAlerta(w http.ResponseWriter, r *http.Request) {
	var nueva models.AlertaEmergencia
	if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	if nueva.Requerimiento == "" || nueva.Gravedad == 0 {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}
	alertaID++
	nueva.ID = alertaID
	nueva.Estado = "Buscando"
	alertasDB = append(alertasDB, nueva)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nueva)
}
func GetAlertas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(alertasDB)
}
func GetAlerta(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	for _, a := range alertasDB {
		if a.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(a)
			return
		}
	}
	http.Error(w, "No encontrado", http.StatusNotFound)
}