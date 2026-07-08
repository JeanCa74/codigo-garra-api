package models

// PerfilVeterinario representa a una clínica veterinaria registrada en la plataforma.
// Has-Many: RecursoClinico (los equipos que posee la clínica).
type PerfilVeterinario struct {
	ID        int              `json:"id" gorm:"primaryKey"`
	Nombre    string           `json:"nombre" gorm:"not null"`
	Telefono  string           `json:"telefono" gorm:"not null"`
	Direccion string           `json:"direccion"`
	Activo    bool             `json:"activo" gorm:"default:true"`
	Recursos  []RecursoClinico `json:"recursos,omitempty" gorm:"foreignKey:PerfilID;constraint:OnDelete:CASCADE"`
}
