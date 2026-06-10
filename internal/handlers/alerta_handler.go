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