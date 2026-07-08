package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
)

func (s *Server) ListarAlertas(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, s.deps.Alertas.Listar())
}

func (s *Server) ObtenerAlerta(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	alerta, err := s.deps.Alertas.Obtener(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, alerta)
}

func (s *Server) CrearAlerta(w http.ResponseWriter, r *http.Request) {
	var a models.AlertaEmergencia
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		RespondError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	creada, err := s.deps.Alertas.Crear(a)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}

func (s *Server) ActualizarAlerta(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	var datos models.AlertaEmergencia
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	actualizada, err := s.deps.Alertas.Actualizar(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizada)
}

func (s *Server) BorrarAlerta(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	if err := s.deps.Alertas.Borrar(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListarAsignacionesDeAlerta devuelve todos los recursos asignados a una alerta concreta.
// Responde [] (array vacío) cuando la alerta no tiene asignaciones aún.
func (s *Server) ListarAsignacionesDeAlerta(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	asignaciones := s.deps.Asignaciones.ListarPorAlerta(id)
	if asignaciones == nil {
		asignaciones = []models.AsignacionTriage{}
	}
	RespondJSON(w, http.StatusOK, asignaciones)
}
