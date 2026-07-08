package models

// HistorialMedico es una entrada clínica en el historial de una mascota.
// Belongs-To: Mascota (el paciente al que pertenece el registro).
type HistorialMedico struct {
	ID          int      `json:"id" gorm:"primaryKey"`
	MascotaID   int      `json:"mascota_id" gorm:"not null"`
	Mascota     *Mascota `json:"-" gorm:"foreignKey:MascotaID"`
	Diagnostico string   `json:"diagnostico" gorm:"not null"`
	Tratamiento string   `json:"tratamiento"`
	Fecha       string   `json:"fecha" gorm:"not null"` // YYYY-MM-DD
	Veterinario string   `json:"veterinario"`
}
