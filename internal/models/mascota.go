package models

// Mascota representa el paciente animal registrado en el sistema.
type Mascota struct {
	ID      int    `json:"id" gorm:"primaryKey"`
	Nombre  string `json:"nombre" gorm:"not null"`
	Especie string `json:"especie" gorm:"not null"` // perro, gato, ave, etc.
	Edad    int    `json:"edad"`
	Dueno   string `json:"dueno" gorm:"not null"`
}
