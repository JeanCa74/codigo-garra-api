package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
)

func (s *Server) ListarPerfiles(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, s.Perfiles.Listar())
}

func (s *Server) ObtenerPerfil(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	perfil, err := s.Perfiles.Obtener(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, perfil)
}

func (s *Server) CrearPerfil(w http.ResponseWriter, r *http.Request) {
	var p models.PerfilVeterinario
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		RespondError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	creado, err := s.Perfiles.Crear(p)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creado)
}

func (s *Server) ActualizarPerfil(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	var datos models.PerfilVeterinario
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	actualizado, err := s.Perfiles.Actualizar(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizado)
}

func (s *Server) BorrarPerfil(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	if err := s.Perfiles.Borrar(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListarRecursosDePerfil lista los recursos asociados a un perfil veterinario concreto.
func (s *Server) ListarRecursosDePerfil(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	RespondJSON(w, http.StatusOK, s.Recursos.ListarPorPerfil(id))
}
