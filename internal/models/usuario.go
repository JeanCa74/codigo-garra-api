package models

import "time"

// Usuario representa una cuenta de acceso. Rol puede ser "veterinario" (defecto) o "admin".
type Usuario struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Rol          string    `json:"rol" gorm:"default:'veterinario';not null"`
	CreadoEn     time.Time `json:"creado_en"`
}
