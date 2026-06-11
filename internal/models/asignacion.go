package models

type AsignacionTriage struct {
	ID                 int    `json:"id"`
	AlertaID           int    `json:"alerta_id"`
	RecursoID          int    `json:"recurso_id"`
	EstadoConfirmacion string `json:"estado_confirmacion"`
}
