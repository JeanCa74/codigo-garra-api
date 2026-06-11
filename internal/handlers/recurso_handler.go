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