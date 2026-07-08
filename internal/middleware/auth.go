package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/JeanCa74/codigo-garra-api/internal/service"
)

type contextKey string

const ctxRol contextKey = "rol"

// Auth valida el Bearer token JWT. Extrae el claim "rol" y lo almacena en el contexto.
// Responde 401 y corta si el token está ausente o es inválido.
func Auth(auth *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				responderNoAutorizado(w)
				return
			}
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, ok := auth.ValidarToken(tokenStr)
			if !ok {
				responderNoAutorizado(w)
				return
			}
			rol, _ := claims["rol"].(string)
			ctx := context.WithValue(r.Context(), ctxRol, rol)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func responderNoAutorizado(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"token ausente o inválido"}`))
}
