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

// NuevoAlmacenGORM abre la base, ejecuta AutoMigrate para todas las entidades y devuelve el almacén.
// DSN ":memory:" para tests aislados; "codigogarra.db" para producción.
func NuevoAlmacenGORM(dsn string) (*AlmacenGORM, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		// Jean Carlos — Emergencia médica
		&models.AlertaEmergencia{},
		&models.AsignacionTriage{},
		// John Erick — Perfil veterinario
		&models.PerfilVeterinario{},
		&models.RecursoClinico{},
		// María José — Historial médico
		&models.Mascota{},
		&models.HistorialMedico{},
		// Auth
		&models.Usuario{},
	); err != nil {
		return nil, err
	}
	return &AlmacenGORM{db: db}, nil
}

// =========================================================
// ALERTAS (Jean Carlos)
// =========================================================

func (a *AlmacenGORM) ListarAlertas() []models.AlertaEmergencia {
	var alertas []models.AlertaEmergencia
	a.db.Preload("Mascota").Find(&alertas)
	return alertas
}

func (a *AlmacenGORM) BuscarAlertaPorID(id int) (models.AlertaEmergencia, bool) {
	var alerta models.AlertaEmergencia
	if err := a.db.Preload("Mascota").First(&alerta, id).Error; err != nil {
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
	datos.CreadoEn = existente.CreadoEn // preservar timestamp de creación
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenGORM) ListarAsignacionesPorAlerta(alertaID int) []models.AsignacionTriage {
	var asignaciones []models.AsignacionTriage
	a.db.Where("alerta_id = ?", alertaID).Find(&asignaciones)
	return asignaciones
}

func (a *AlmacenGORM) BorrarAlerta(id int) bool {
	return a.db.Delete(&models.AlertaEmergencia{}, id).RowsAffected > 0
}

// =========================================================
// ASIGNACIONES (Jean Carlos)
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
	return a.db.Delete(&models.AsignacionTriage{}, id).RowsAffected > 0
}

// =========================================================
// PERFILES VETERINARIOS (John Erick)
// =========================================================

func (a *AlmacenGORM) ListarPerfiles() []models.PerfilVeterinario {
	var perfiles []models.PerfilVeterinario
	a.db.Preload("Recursos").Find(&perfiles)
	return perfiles
}

func (a *AlmacenGORM) BuscarPerfilPorID(id int) (models.PerfilVeterinario, bool) {
	var p models.PerfilVeterinario
	if err := a.db.Preload("Recursos").First(&p, id).Error; err != nil {
		return models.PerfilVeterinario{}, false
	}
	return p, true
}

func (a *AlmacenGORM) CrearPerfil(p models.PerfilVeterinario) models.PerfilVeterinario {
	a.db.Create(&p)
	return p
}

func (a *AlmacenGORM) ActualizarPerfil(id int, datos models.PerfilVeterinario) (models.PerfilVeterinario, bool) {
	var existente models.PerfilVeterinario
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.PerfilVeterinario{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenGORM) BorrarPerfil(id int) bool {
	return a.db.Delete(&models.PerfilVeterinario{}, id).RowsAffected > 0
}

// =========================================================
// RECURSOS CLÍNICOS (John Erick)
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
	return a.db.Delete(&models.RecursoClinico{}, id).RowsAffected > 0
}

func (a *AlmacenGORM) ListarRecursosPorPerfil(perfilID int) []models.RecursoClinico {
	var recursos []models.RecursoClinico
	a.db.Where("perfil_id = ?", perfilID).Find(&recursos)
	return recursos
}

// =========================================================
// MASCOTAS (María José)
// =========================================================

func (a *AlmacenGORM) ListarMascotas() []models.Mascota {
	var mascotas []models.Mascota
	a.db.Preload("Historial").Find(&mascotas)
	return mascotas
}

func (a *AlmacenGORM) BuscarMascotaPorID(id int) (models.Mascota, bool) {
	var m models.Mascota
	if err := a.db.Preload("Historial").First(&m, id).Error; err != nil {
		return models.Mascota{}, false
	}
	return m, true
}

func (a *AlmacenGORM) CrearMascota(m models.Mascota) models.Mascota {
	a.db.Create(&m)
	return m
}

func (a *AlmacenGORM) ActualizarMascota(id int, datos models.Mascota) (models.Mascota, bool) {
	var existente models.Mascota
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Mascota{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenGORM) BorrarMascota(id int) bool {
	return a.db.Delete(&models.Mascota{}, id).RowsAffected > 0
}

// =========================================================
// HISTORIAL MÉDICO (María José)
// =========================================================

func (a *AlmacenGORM) ListarHistorial() []models.HistorialMedico {
	var historial []models.HistorialMedico
	a.db.Find(&historial)
	return historial
}

func (a *AlmacenGORM) BuscarHistorialPorID(id int) (models.HistorialMedico, bool) {
	var h models.HistorialMedico
	if err := a.db.First(&h, id).Error; err != nil {
		return models.HistorialMedico{}, false
	}
	return h, true
}

func (a *AlmacenGORM) CrearHistorial(h models.HistorialMedico) models.HistorialMedico {
	a.db.Create(&h)
	return h
}

func (a *AlmacenGORM) ListarHistorialPorMascota(mascotaID int) []models.HistorialMedico {
	var historial []models.HistorialMedico
	a.db.Where("mascota_id = ?", mascotaID).Find(&historial)
	return historial
}

// =========================================================
// USUARIOS (Auth)
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
