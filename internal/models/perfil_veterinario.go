package models

// PerfilVeterinario representa una clínica veterinaria registrada en el sistema.
type PerfilVeterinario struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Nombre    string `json:"nombre" gorm:"not null"`
	Telefono  string `json:"telefono" gorm:"not null"`
	Direccion string `json:"direccion"`
	Activo    bool   `json:"activo" gorm:"default:true"`
}
