package middleware

import (
	"net/http"
)

// RequireRol rechaza con 403 Forbidden si el rol del usuario autenticado
// no coincide con el requerido. Debe colocarse detrás del middleware Auth.
func RequireRol(rol string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rolUsuario, _ := r.Context().Value(ctxRol).(string)
			if rolUsuario != rol {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error":"permisos insuficientes"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
