package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
)

// NuevoAlmacenPostgres abre una conexión a PostgreSQL y ejecuta AutoMigrate.
// DSN ejemplo: "host=db user=garra password=garra dbname=garra port=5432 sslmode=disable"
func NuevoAlmacenPostgres(dsn string) (*AlmacenGORM, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		&models.AlertaEmergencia{},
		&models.AsignacionTriage{},
		&models.PerfilVeterinario{},
		&models.RecursoClinico{},
		&models.Mascota{},
		&models.HistorialMedico{},
		&models.Usuario{},
	); err != nil {
		return nil, err
	}
	return &AlmacenGORM{db: db}, nil
}
