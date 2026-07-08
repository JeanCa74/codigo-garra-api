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

	// Clínicas con ubicaciones reales de Manta, Ecuador
	p1 := almacen.CrearPerfil(models.PerfilVeterinario{
		Nombre: "Clínica Veterinaria Los Ceibos", Telefono: "+593992345678",
		Direccion: "Av. 4 de Noviembre y Calle 14, Manta", Activo: true,
	})
	p2 := almacen.CrearPerfil(models.PerfilVeterinario{
		Nombre: "VetCare Manta", Telefono: "+593987654321",
		Direccion: "Av. Flavio Alfaro y Av. 24, Manta", Activo: true,
	})

	almacen.CrearRecurso(models.RecursoClinico{
		PerfilID: p1.ID, TipoMaquina: "Ecógrafo portátil", EstaDisponible: true,
	})
	r2 := almacen.CrearRecurso(models.RecursoClinico{
		PerfilID: p1.ID, TipoMaquina: "Ventilador mecánico", EstaDisponible: true,
	})
	almacen.CrearRecurso(models.RecursoClinico{
		PerfilID: p2.ID, TipoMaquina: "Desfibrilador", EstaDisponible: false,
	})

	m1 := almacen.CrearMascota(models.Mascota{
		Nombre: "Rex", Especie: "Perro", Edad: 4, Dueno: "Gabriela Vera",
	})
	m2 := almacen.CrearMascota(models.Mascota{
		Nombre: "Michi", Especie: "Gato", Edad: 2, Dueno: "Juan Cedeño",
	})

	almacen.CrearHistorial(models.HistorialMedico{
		MascotaID: m1.ID, Diagnostico: "Trauma abdominal por atropellamiento",
		Tratamiento: "Cirugía de urgencia + analgésicos", Fecha: "2026-07-01", Veterinario: "Dra. Vélez",
	})
	almacen.CrearHistorial(models.HistorialMedico{
		MascotaID: m2.ID, Diagnostico: "Insuficiencia respiratoria aguda",
		Tratamiento: "Nebulización y antibióticos 7 días", Fecha: "2026-07-05", Veterinario: "Dr. Rivadeneira",
	})

	// Alerta activa: Rex atropellado — vinculada al paciente registrado
	a1 := almacen.CrearAlerta(models.AlertaEmergencia{
		MascotaID: m1.ID, Gravedad: 5,
		Requerimiento: "Ventilador mecánico urgente — perro atropellado", Estado: "Buscando",
	})
	almacen.CrearAlerta(models.AlertaEmergencia{
		MascotaID: m2.ID, Gravedad: 3,
		Requerimiento: "Ecógrafo para diagnóstico respiratorio", Estado: "Atendido",
	})

	// Recurso ya asignado a la emergencia crítica
	almacen.CrearAsignacion(models.AsignacionTriage{
		AlertaID: a1.ID, RecursoID: r2.ID, EstadoConfirmacion: "Confirmado",
	})

	log.Println("seeder: datos iniciales sembrados correctamente")
	return nil
}
