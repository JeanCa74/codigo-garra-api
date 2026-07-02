package storage

import (
	"sync"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
)

// Memoria es el almacén unificado en RAM; se usa como fake en los tests de handler.
type Memoria struct {
	// Jean Carlos — Emergencia médica
	alertas          []models.AlertaEmergencia
	nextAlertaID     int
	asignaciones     []models.AsignacionTriage
	nextAsignacionID int

	// John Erick — Perfil veterinario
	perfiles      []models.PerfilVeterinario
	nextPerfilID  int
	recursos      []models.RecursoClinico
	nextRecursoID int

	// María José — Historial médico
	mascotas         []models.Mascota
	nextMascotaID    int
	historial        []models.HistorialMedico
	nextHistorialID  int

	mu sync.Mutex
}

func NuevaMemoria() *Memoria {
	return &Memoria{
		alertas:          []models.AlertaEmergencia{},
		nextAlertaID:     1,
		asignaciones:     []models.AsignacionTriage{},
		nextAsignacionID: 1,
		perfiles:         []models.PerfilVeterinario{},
		nextPerfilID:     1,
		recursos:         []models.RecursoClinico{},
		nextRecursoID:    1,
		mascotas:         []models.Mascota{},
		nextMascotaID:    1,
		historial:        []models.HistorialMedico{},
		nextHistorialID:  1,
	}
}

// SeedAlertas carga alertas de prueba para los tests de Jean Carlos.
func (m *Memoria) SeedAlertas() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alertas = []models.AlertaEmergencia{
		{ID: 1, Gravedad: 3, Requerimiento: "Ecógrafo", Estado: "Buscando"},
		{ID: 2, Gravedad: 5, Requerimiento: "Ventilador", Estado: "Asignada"},
	}
	m.nextAlertaID = 3
}

// =========================================================
// ALERTAS (Jean Carlos)
// =========================================================

func (m *Memoria) ListarAlertas() []models.AlertaEmergencia {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]models.AlertaEmergencia, len(m.alertas))
	copy(copia, m.alertas)
	return copia
}

func (m *Memoria) BuscarAlertaPorID(id int) (models.AlertaEmergencia, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, a := range m.alertas {
		if a.ID == id {
			return a, true
		}
	}
	return models.AlertaEmergencia{}, false
}

func (m *Memoria) CrearAlerta(a models.AlertaEmergencia) models.AlertaEmergencia {
	m.mu.Lock()
	defer m.mu.Unlock()
	a.ID = m.nextAlertaID
	m.nextAlertaID++
	if a.Estado == "" {
		a.Estado = "Buscando"
	}
	m.alertas = append(m.alertas, a)
	return a
}

func (m *Memoria) ActualizarAlerta(id int, datos models.AlertaEmergencia) (models.AlertaEmergencia, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, a := range m.alertas {
		if a.ID == id {
			datos.ID = id
			m.alertas[i] = datos
			return datos, true
		}
	}
	return models.AlertaEmergencia{}, false
}

func (m *Memoria) BorrarAlerta(id int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, a := range m.alertas {
		if a.ID == id {
			m.alertas = append(m.alertas[:i], m.alertas[i+1:]...)
			return true
		}
	}
	return false
}

// =========================================================
// ASIGNACIONES (Jean Carlos)
// =========================================================

func (m *Memoria) ListarAsignaciones() []models.AsignacionTriage {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]models.AsignacionTriage, len(m.asignaciones))
	copy(copia, m.asignaciones)
	return copia
}

func (m *Memoria) BuscarAsignacionPorID(id int) (models.AsignacionTriage, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, a := range m.asignaciones {
		if a.ID == id {
			return a, true
		}
	}
	return models.AsignacionTriage{}, false
}

func (m *Memoria) CrearAsignacion(a models.AsignacionTriage) models.AsignacionTriage {
	m.mu.Lock()
	defer m.mu.Unlock()
	a.ID = m.nextAsignacionID
	m.nextAsignacionID++
	if a.EstadoConfirmacion == "" {
		a.EstadoConfirmacion = "Pendiente"
	}
	m.asignaciones = append(m.asignaciones, a)
	return a
}

func (m *Memoria) ActualizarAsignacion(id int, datos models.AsignacionTriage) (models.AsignacionTriage, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, a := range m.asignaciones {
		if a.ID == id {
			datos.ID = id
			m.asignaciones[i] = datos
			return datos, true
		}
	}
	return models.AsignacionTriage{}, false
}

func (m *Memoria) BorrarAsignacion(id int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, a := range m.asignaciones {
		if a.ID == id {
			m.asignaciones = append(m.asignaciones[:i], m.asignaciones[i+1:]...)
			return true
		}
	}
	return false
}

// =========================================================
// PERFILES VETERINARIOS (John Erick)
// =========================================================

func (m *Memoria) ListarPerfiles() []models.PerfilVeterinario {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]models.PerfilVeterinario, len(m.perfiles))
	copy(copia, m.perfiles)
	return copia
}

func (m *Memoria) BuscarPerfilPorID(id int) (models.PerfilVeterinario, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, p := range m.perfiles {
		if p.ID == id {
			return p, true
		}
	}
	return models.PerfilVeterinario{}, false
}

func (m *Memoria) CrearPerfil(p models.PerfilVeterinario) models.PerfilVeterinario {
	m.mu.Lock()
	defer m.mu.Unlock()
	p.ID = m.nextPerfilID
	m.nextPerfilID++
	m.perfiles = append(m.perfiles, p)
	return p
}

func (m *Memoria) ActualizarPerfil(id int, datos models.PerfilVeterinario) (models.PerfilVeterinario, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, p := range m.perfiles {
		if p.ID == id {
			datos.ID = id
			m.perfiles[i] = datos
			return datos, true
		}
	}
	return models.PerfilVeterinario{}, false
}

func (m *Memoria) BorrarPerfil(id int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, p := range m.perfiles {
		if p.ID == id {
			m.perfiles = append(m.perfiles[:i], m.perfiles[i+1:]...)
			return true
		}
	}
	return false
}

// =========================================================
// RECURSOS CLÍNICOS (John Erick)
// =========================================================

func (m *Memoria) ListarRecursos() []models.RecursoClinico {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]models.RecursoClinico, len(m.recursos))
	copy(copia, m.recursos)
	return copia
}

func (m *Memoria) BuscarRecursoPorID(id int) (models.RecursoClinico, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, r := range m.recursos {
		if r.ID == id {
			return r, true
		}
	}
	return models.RecursoClinico{}, false
}

func (m *Memoria) CrearRecurso(r models.RecursoClinico) models.RecursoClinico {
	m.mu.Lock()
	defer m.mu.Unlock()
	r.ID = m.nextRecursoID
	m.nextRecursoID++
	m.recursos = append(m.recursos, r)
	return r
}

func (m *Memoria) ActualizarRecurso(id int, datos models.RecursoClinico) (models.RecursoClinico, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, r := range m.recursos {
		if r.ID == id {
			datos.ID = id
			m.recursos[i] = datos
			return datos, true
		}
	}
	return models.RecursoClinico{}, false
}

func (m *Memoria) BorrarRecurso(id int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, r := range m.recursos {
		if r.ID == id {
			m.recursos = append(m.recursos[:i], m.recursos[i+1:]...)
			return true
		}
	}
	return false
}

func (m *Memoria) ListarRecursosPorPerfil(perfilID int) []models.RecursoClinico {
	m.mu.Lock()
	defer m.mu.Unlock()
	var resultado []models.RecursoClinico
	for _, r := range m.recursos {
		if r.PerfilID == perfilID {
			resultado = append(resultado, r)
		}
	}
	return resultado
}

// =========================================================
// MASCOTAS (María José)
// =========================================================

func (m *Memoria) ListarMascotas() []models.Mascota {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]models.Mascota, len(m.mascotas))
	copy(copia, m.mascotas)
	return copia
}

func (m *Memoria) BuscarMascotaPorID(id int) (models.Mascota, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, ma := range m.mascotas {
		if ma.ID == id {
			return ma, true
		}
	}
	return models.Mascota{}, false
}

func (m *Memoria) CrearMascota(ma models.Mascota) models.Mascota {
	m.mu.Lock()
	defer m.mu.Unlock()
	ma.ID = m.nextMascotaID
	m.nextMascotaID++
	m.mascotas = append(m.mascotas, ma)
	return ma
}

func (m *Memoria) ActualizarMascota(id int, datos models.Mascota) (models.Mascota, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, ma := range m.mascotas {
		if ma.ID == id {
			datos.ID = id
			m.mascotas[i] = datos
			return datos, true
		}
	}
	return models.Mascota{}, false
}

func (m *Memoria) BorrarMascota(id int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, ma := range m.mascotas {
		if ma.ID == id {
			m.mascotas = append(m.mascotas[:i], m.mascotas[i+1:]...)
			return true
		}
	}
	return false
}

// =========================================================
// HISTORIAL MÉDICO (María José)
// =========================================================

func (m *Memoria) ListarHistorial() []models.HistorialMedico {
	m.mu.Lock()
	defer m.mu.Unlock()
	copia := make([]models.HistorialMedico, len(m.historial))
	copy(copia, m.historial)
	return copia
}

func (m *Memoria) BuscarHistorialPorID(id int) (models.HistorialMedico, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, h := range m.historial {
		if h.ID == id {
			return h, true
		}
	}
	return models.HistorialMedico{}, false
}

func (m *Memoria) CrearHistorial(h models.HistorialMedico) models.HistorialMedico {
	m.mu.Lock()
	defer m.mu.Unlock()
	h.ID = m.nextHistorialID
	m.nextHistorialID++
	m.historial = append(m.historial, h)
	return h
}

func (m *Memoria) ListarHistorialPorMascota(mascotaID int) []models.HistorialMedico {
	m.mu.Lock()
	defer m.mu.Unlock()
	var resultado []models.HistorialMedico
	for _, h := range m.historial {
		if h.MascotaID == mascotaID {
			resultado = append(resultado, h)
		}
	}
	return resultado
}

// Chequeo en tiempo de compilación: Memoria debe cumplir Almacen.
var _ Almacen = (*Memoria)(nil)
