package models

// RecursoClinico representa un equipo médico de una clínica veterinaria.
type RecursoClinico struct {
	ID             int    `json:"id" gorm:"primaryKey"`
	ClinicaID      int    `json:"clinica_id" gorm:"not null"`
	TipoMaquina    string `json:"tipo_maquina" gorm:"not null"`
	EstaDisponible bool   `json:"esta_disponible" gorm:"default:true"`
}
