package models

// AlertaEmergencia representa una llamada de urgencia veterinaria.
// Gravedad sigue la escala de triage: 1 (leve) a 5 (crítico).
type AlertaEmergencia struct {
	ID            int    `json:"id" gorm:"primaryKey"`
	Gravedad      int    `json:"gravedad" gorm:"not null"`
	Requerimiento string `json:"requerimiento" gorm:"not null"`
	Estado        string `json:"estado" gorm:"default:'Buscando'"`
}
