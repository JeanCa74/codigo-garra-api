package middleware

import (
	"net/http"
	"strings"

	"github.com/JeanCa74/codigo-garra-api/internal/service"
)

// Auth valida el Bearer token JWT. Si falta o es inválido responde 401 y corta.
func Auth(auth *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				responderNoAutorizado(w)
				return
			}
			token := strings.TrimPrefix(header, "Bearer ")
			if !auth.ValidarToken(token) {
				responderNoAutorizado(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func responderNoAutorizado(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"token ausente o inválido"}`))
}
