package storage

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
)

// AlmacenGORM implementa Almacen + UserRepository usando GORM sobre SQLite.
type AlmacenGORM struct {
	db *gorm.DB
}

// NuevoAlmacenGORM abre la base, ejecuta AutoMigrate y devuelve el almacén.
// Usar DSN ":memory:" para tests aislados; "codigogarra.db" para producción.
func NuevoAlmacenGORM(dsn string) (*AlmacenGORM, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		&models.AlertaEmergencia{},
		&models.RecursoClinico{},
		&models.AsignacionTriage{},
		&models.Usuario{},
	); err != nil {
		return nil, err
	}
	return &AlmacenGORM{db: db}, nil
}

// =========================================================
// ALERTAS
// =========================================================

func (a *AlmacenGORM) ListarAlertas() []models.AlertaEmergencia {
	var alertas []models.AlertaEmergencia
	a.db.Find(&alertas)
	return alertas
}

func (a *AlmacenGORM) BuscarAlertaPorID(id int) (models.AlertaEmergencia, bool) {
	var alerta models.AlertaEmergencia
	if err := a.db.First(&alerta, id).Error; err != nil {
		return models.AlertaEmergencia{}, false
	}
	return alerta, true
}

func (a *AlmacenGORM) CrearAlerta(al models.AlertaEmergencia) models.AlertaEmergencia {
	if al.Estado == "" {
		al.Estado = "Buscando"
	}
	a.db.Create(&al)
	return al
}

func (a *AlmacenGORM) ActualizarAlerta(id int, datos models.AlertaEmergencia) (models.AlertaEmergencia, bool) {
	var existente models.AlertaEmergencia
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.AlertaEmergencia{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenGORM) BorrarAlerta(id int) bool {
	res := a.db.Delete(&models.AlertaEmergencia{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// RECURSOS
// =========================================================

func (a *AlmacenGORM) ListarRecursos() []models.RecursoClinico {
	var recursos []models.RecursoClinico
	a.db.Find(&recursos)
	return recursos
}

func (a *AlmacenGORM) BuscarRecursoPorID(id int) (models.RecursoClinico, bool) {
	var r models.RecursoClinico
	if err := a.db.First(&r, id).Error; err != nil {
		return models.RecursoClinico{}, false
	}
	return r, true
}

func (a *AlmacenGORM) CrearRecurso(r models.RecursoClinico) models.RecursoClinico {
	a.db.Create(&r)
	return r
}

func (a *AlmacenGORM) ActualizarRecurso(id int, datos models.RecursoClinico) (models.RecursoClinico, bool) {
	var existente models.RecursoClinico
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.RecursoClinico{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenGORM) BorrarRecurso(id int) bool {
	res := a.db.Delete(&models.RecursoClinico{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// ASIGNACIONES
// =========================================================

func (a *AlmacenGORM) ListarAsignaciones() []models.AsignacionTriage {
	var asignaciones []models.AsignacionTriage
	a.db.Find(&asignaciones)
	return asignaciones
}

func (a *AlmacenGORM) BuscarAsignacionPorID(id int) (models.AsignacionTriage, bool) {
	var asig models.AsignacionTriage
	if err := a.db.First(&asig, id).Error; err != nil {
		return models.AsignacionTriage{}, false
	}
	return asig, true
}

func (a *AlmacenGORM) CrearAsignacion(asig models.AsignacionTriage) models.AsignacionTriage {
	if asig.EstadoConfirmacion == "" {
		asig.EstadoConfirmacion = "Pendiente"
	}
	a.db.Create(&asig)
	return asig
}

func (a *AlmacenGORM) ActualizarAsignacion(id int, datos models.AsignacionTriage) (models.AsignacionTriage, bool) {
	var existente models.AsignacionTriage
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.AsignacionTriage{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenGORM) BorrarAsignacion(id int) bool {
	res := a.db.Delete(&models.AsignacionTriage{}, id)
	return res.RowsAffected > 0
}

// =========================================================
// USUARIOS
// =========================================================

func (a *AlmacenGORM) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	if err := a.db.Create(&u).Error; err != nil {
		return models.Usuario{}, err
	}
	return u, nil
}

func (a *AlmacenGORM) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	var u models.Usuario
	if err := a.db.Where("email = ?", email).First(&u).Error; err != nil {
		return models.Usuario{}, false
	}
	return u, true
}

// Chequeos en tiempo de compilación.
var _ Almacen = (*AlmacenGORM)(nil)
var _ UserRepository = (*AlmacenGORM)(nil)
