package models

import "time"

// AlertaEmergencia representa una llamada de urgencia veterinaria.
// Gravedad sigue la escala de triage: 1 (leve) a 5 (crítico).
// MascotaID = 0 significa que la mascota aún no está registrada en el sistema.
type AlertaEmergencia struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	MascotaID     int       `json:"mascota_id"`
	Mascota       *Mascota  `json:"-" gorm:"foreignKey:MascotaID"`
	Gravedad      int       `json:"gravedad" gorm:"not null"`
	Requerimiento string    `json:"requerimiento" gorm:"not null"`
	Estado        string    `json:"estado" gorm:"default:'Buscando'"`
	CreadoEn      time.Time `json:"creado_en" gorm:"autoCreateTime"`
}
