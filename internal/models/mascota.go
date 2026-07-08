package models

// Mascota representa el paciente animal registrado en la plataforma.
// Has-Many: HistorialMedico (todas las entradas clínicas de la mascota).
type Mascota struct {
	ID        int               `json:"id" gorm:"primaryKey"`
	Nombre    string            `json:"nombre" gorm:"not null"`
	Especie   string            `json:"especie" gorm:"not null"`
	Edad      int               `json:"edad"`
	Dueno     string            `json:"dueno" gorm:"not null"`
	Historial []HistorialMedico `json:"historial,omitempty" gorm:"foreignKey:MascotaID;constraint:OnDelete:CASCADE"`
}
