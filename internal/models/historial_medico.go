package models

// HistorialMedico registra cada consulta o atención que recibió una Mascota.
type HistorialMedico struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	MascotaID   int    `json:"mascota_id" gorm:"not null"`
	Diagnostico string `json:"diagnostico" gorm:"not null"`
	Tratamiento string `json:"tratamiento"`
	Fecha       string `json:"fecha" gorm:"not null"` // formato YYYY-MM-DD
	Veterinario string `json:"veterinario"`
}
