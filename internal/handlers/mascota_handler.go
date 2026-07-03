package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
)

func (s *Server) ListarMascotas(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, s.deps.Mascotas.Listar())
}

func (s *Server) ObtenerMascota(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	mascota, err := s.deps.Mascotas.Obtener(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, mascota)
}

func (s *Server) CrearMascota(w http.ResponseWriter, r *http.Request) {
	var m models.Mascota
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		RespondError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	creada, err := s.deps.Mascotas.Crear(m)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}

func (s *Server) ActualizarMascota(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	var datos models.Mascota
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	actualizada, err := s.deps.Mascotas.Actualizar(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizada)
}

func (s *Server) BorrarMascota(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id inválido")
		return
	}
	if err := s.deps.Mascotas.Borrar(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
