package models

// AsignacionTriage vincula una AlertaEmergencia con el RecursoClinico asignado.
// Belongs-To: AlertaEmergencia, RecursoClinico.
type AsignacionTriage struct {
	ID                 int               `json:"id" gorm:"primaryKey"`
	AlertaID           int               `json:"alerta_id" gorm:"not null"`
	Alerta             *AlertaEmergencia `json:"-" gorm:"foreignKey:AlertaID;constraint:OnDelete:RESTRICT"`
	RecursoID          int               `json:"recurso_id" gorm:"not null"`
	Recurso            *RecursoClinico   `json:"-" gorm:"foreignKey:RecursoID;constraint:OnDelete:RESTRICT"`
	EstadoConfirmacion string            `json:"estado_confirmacion" gorm:"default:'Pendiente'"`
}
