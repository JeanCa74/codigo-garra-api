package models

// AsignacionTriage vincula una AlertaEmergencia con el RecursoClinico asignado.
type AsignacionTriage struct {
	ID                 int    `json:"id" gorm:"primaryKey"`
	AlertaID           int    `json:"alerta_id" gorm:"not null"`
	RecursoID          int    `json:"recurso_id" gorm:"not null"`
	EstadoConfirmacion string `json:"estado_confirmacion" gorm:"default:'Pendiente'"`
}
