package models

// RecursoClinico representa un equipo médico disponible en una clínica veterinaria.
// Belongs-To: PerfilVeterinario (el perfil dueño del recurso).
type RecursoClinico struct {
	ID             int                `json:"id" gorm:"primaryKey"`
	PerfilID       int                `json:"perfil_id" gorm:"not null"`
	Perfil         *PerfilVeterinario `json:"-" gorm:"foreignKey:PerfilID"`
	TipoMaquina    string             `json:"tipo_maquina" gorm:"not null"`
	EstaDisponible bool               `json:"esta_disponible" gorm:"default:true"`
}
