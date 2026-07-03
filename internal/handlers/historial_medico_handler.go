package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
)

func (s *Server) ListarHistorial(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, s.deps.HistorialMed.Listar())
}

func (s *Server) ObtenerHistorial(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	h, err := s.deps.HistorialMed.ObtenerPorID(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, h)
}

func (s *Server) CrearHistorial(w http.ResponseWriter, r *http.Request) {
	var h models.HistorialMedico
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
		RespondError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	creado, err := s.deps.HistorialMed.Crear(h)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

// ListarHistorialDeMascota lista todas las entradas del historial de una mascota concreta.
func (s *Server) ListarHistorialDeMascota(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	RespondJSON(w, http.StatusOK, s.deps.HistorialMed.ListarPorMascota(id))
}
