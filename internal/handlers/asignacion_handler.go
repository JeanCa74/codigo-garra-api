package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
)

func (s *Server) ListarAsignaciones(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, s.deps.Asignaciones.Listar())
}

func (s *Server) ObtenerAsignacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	asig, err := s.deps.Asignaciones.Obtener(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, asig)
}

func (s *Server) CrearAsignacion(w http.ResponseWriter, r *http.Request) {
	var a models.AsignacionTriage
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		RespondError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	creada, err := s.deps.Asignaciones.Crear(a)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}

func (s *Server) ActualizarAsignacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	var datos models.AsignacionTriage
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	actualizada, err := s.deps.Asignaciones.Actualizar(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizada)
}

func (s *Server) BorrarAsignacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	if err := s.deps.Asignaciones.Borrar(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
