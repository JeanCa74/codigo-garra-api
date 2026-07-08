package seed

import (
	"log"

	"github.com/JeanCa74/codigo-garra-api/internal/models"
	"github.com/JeanCa74/codigo-garra-api/internal/service"
	"github.com/JeanCa74/codigo-garra-api/internal/storage"
)

// Sembrar inserta datos iniciales solo si la base está vacía.
// Cuentas creadas:
//   - admin@codigogarra.vet / Admin123!  (rol: admin)
//   - vet@codigogarra.vet  / Vet123!    (rol: veterinario)
func Sembrar(almacen storage.Almacen, usuarios storage.UserRepository, authSvc *service.AuthService) error {
	if _, existe := usuarios.BuscarUsuarioPorEmail("admin@codigogarra.vet"); existe {
		log.Println("seeder: ya existen datos, omitiendo")
		return nil
	}

	if _, err := authSvc.RegistrarAdmin("admin@codigogarra.vet", "Admin123!"); err != nil {
		return err
	}
	if _, err := authSvc.Registrar("vet@codigogarra.vet", "Vet123!"); err != nil {
		return err
	}

	p1 := almacen.CrearPerfil(models.PerfilVeterinario{
		Nombre: "Clínica Garra Norte", Telefono: "+56912345678",
		Direccion: "Av. Providencia 1234, Santiago", Activo: true,
	})
	p2 := almacen.CrearPerfil(models.PerfilVeterinario{
		Nombre: "Veterinaria El Bosque", Telefono: "+56987654321",
		Direccion: "Av. Las Condes 5678, Santiago", Activo: true,
	})

	almacen.CrearRecurso(models.RecursoClinico{
		PerfilID: p1.ID, TipoMaquina: "Ecógrafo portátil", EstaDisponible: true,
	})
	almacen.CrearRecurso(models.RecursoClinico{
		PerfilID: p1.ID, TipoMaquina: "Ventilador mecánico", EstaDisponible: true,
	})
	almacen.CrearRecurso(models.RecursoClinico{
		PerfilID: p2.ID, TipoMaquina: "Desfibrilador", EstaDisponible: false,
	})

	m1 := almacen.CrearMascota(models.Mascota{
		Nombre: "Max", Especie: "Perro", Edad: 3, Dueno: "Ana García",
	})
	m2 := almacen.CrearMascota(models.Mascota{
		Nombre: "Luna", Especie: "Gato", Edad: 5, Dueno: "Carlos López",
	})

	almacen.CrearHistorial(models.HistorialMedico{
		MascotaID: m1.ID, Diagnostico: "Fractura en pata delantera",
		Tratamiento: "Inmovilización y analgésicos", Fecha: "2026-06-15", Veterinario: "Dra. Pérez",
	})
	almacen.CrearHistorial(models.HistorialMedico{
		MascotaID: m2.ID, Diagnostico: "Infección respiratoria aguda",
		Tratamiento: "Antibióticos 7 días", Fecha: "2026-07-01", Veterinario: "Dr. Soto",
	})

	almacen.CrearAlerta(models.AlertaEmergencia{
		Gravedad: 5, Requerimiento: "Necesita ventilador mecánico urgente", Estado: "Buscando",
	})
	almacen.CrearAlerta(models.AlertaEmergencia{
		Gravedad: 3, Requerimiento: "Requiere ecógrafo para diagnóstico", Estado: "Atendido",
	})

	almacen.CrearAsignacion(models.AsignacionTriage{
		AlertaID: 1, RecursoID: 2, EstadoConfirmacion: "Confirmado",
	})

	log.Println("seeder: datos iniciales sembrados correctamente")
	return nil
}
