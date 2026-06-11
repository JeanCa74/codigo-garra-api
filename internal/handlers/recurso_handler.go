package handlers
import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/go-chi/chi/v5"
	)
var recursosDB []models.RecursoClinico
var recursoID int = 0

// CreateRecurso inicializa la disponibilidad de un equipo medico
func CreateRecurso(w http.ResponseWriter, r *http.Request) {
var nuevo models.RecursoClinico
if err := json.NewDecoder(r.Body).Decode(&nuevo); err != nil {
http.Error(w, "Datos inválidos", http.StatusBadRequest)
return
}
recursoID++
nuevo.ID = recursoID
recursosDB = append(recursosDB, nuevo)
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(nuevo)
}

func GetRecursos(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(recursosDB)
}

func GetRecurso(w http.ResponseWriter, r *http.Request) {
id, _ := strconv.Atoi(chi.URLParam(r, "id"))
for _, rec := range recursosDB {
if rec.ID == id {
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(rec)
return
}
}
http.Error(w, "Recurso clinico no hallado en la base", http.StatusNotFound)
}

func UpdateRecurso(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(chi.URLParam(r, "id"))
    for i, rec := range recursosDB {
        if rec.ID == id {
            var actualizado models.RecursoClinico
            json.NewDecoder(r.Body).Decode(&actualizado)
            actualizado.ID = id
            recursosDB[i] = actualizado
            json.NewEncoder(w).Encode(actualizado)
            return
        }
    }
    http.Error(w, "Recurso no encontrado", http.StatusNotFound)
}