package models

// RecursoClinico representa un equipo médico disponible en una clínica.
// PerfilID referencia al PerfilVeterinario propietario del equipo.
type RecursoClinico struct {
	ID             int    `json:"id" gorm:"primaryKey"`
	PerfilID       int    `json:"perfil_id" gorm:"not null"`
	TipoMaquina    string `json:"tipo_maquina" gorm:"not null"`
	EstaDisponible bool   `json:"esta_disponible" gorm:"default:true"`
}
