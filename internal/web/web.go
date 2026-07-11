// Package web sirve la interfaz de demostración de Código Garra.
// Los archivos estáticos se embeben en el binario con go:embed, de modo que
// la imagen Docker no necesita copiar nada extra y no hay rutas de disco
// que un atacante pueda manipular (no existe path traversal posible).
package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var archivos embed.FS

// Handler devuelve el servidor de archivos de la interfaz web.
// index.html responde en "/" y los assets (app.js, styles.css) en sus rutas.
func Handler() http.Handler {
	contenido, err := fs.Sub(archivos, "static")
	if err != nil {
		panic(err) // imposible en runtime: el embed se resuelve en compilación
	}
	return http.FileServer(http.FS(contenido))
}
