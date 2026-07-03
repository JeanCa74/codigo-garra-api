package handlers

import (
	"encoding/json"
	"net/http"
)

func (s *Server) Registrar(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	u, err := s.deps.Auth.Registrar(body.Email, body.Password)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, u)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	token, err := s.deps.Auth.Login(body.Email, body.Password)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"token": token})
}
