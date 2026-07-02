package storage

import (
	"sync"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
)

// Memoria es el almacén unificado en RAM; se usa como fake en los tests de handler.
type Memoria struct {
	alertas      []models.AlertaEmergencia
	nextAlertaID int

	recursos      []models.RecursoClinico
	nextRecursoID int

	asignaciones      []models.AsignacionTriage
	nextAsignacionID  int

	mu sync.Mutex
}

func NuevaMemoria() *Memoria {
	return &Memoria{
		alertas:          []models.AlertaEmergencia{},
		nextAlertaID:     1,
		recursos:         []models.RecursoClinico{},
		nextRecursoID:    1,
		asignaciones:     []models.AsignacionTriage{},
		nextAsignacionID: 1,
	}
}

// SeedAlertas carga alertas de prueba para los tests de handler.
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
// ALERTAS
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
// RECURSOS
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

// =========================================================
// ASIGNACIONES
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

// Chequeo en tiempo de compilación.
var _ Almacen = (*Memoria)(nil)
