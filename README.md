# Código Garra API

Plataforma REST de emergencias veterinarias desarrollada en Go.  
Proyecto de curso **TDI-601 – Aplicaciones Web II**, Semana 11, Actividad C1.

---

## Integrantes y módulos

| Estudiante | Módulo | Entidades |
|---|---|---|
| **Jean Carlos Arauz** | Emergencia médica | `AlertaEmergencia`, `AsignacionTriage` |
| **John Erick Bello** | Perfil veterinario | `PerfilVeterinario`, `RecursoClinico` |
| **María José Vinces** | Historial médico | `Mascota`, `HistorialMedico` |

---

## Stack técnico

- **Lenguaje**: Go 1.22+
- **ORM**: GORM + SQLite pure-Go (`github.com/glebarez/sqlite`) — sin CGO
- **Router**: chi v5
- **Autenticación**: JWT (`golang-jwt/jwt/v5`) + bcrypt
- **Tests**: `testing` + `testify` (mock, assert, require) + `httptest`

---

## Estructura del proyecto

```
codigo-garra-api/
├── cmd/
│   └── main.go                 — punto de entrada, wiring DI
├── internal/
│   ├── models/                 — entidades GORM (6 modelos + Usuario)
│   │   ├── alerta.go
│   │   ├── asignacion.go
│   │   ├── perfil_veterinario.go
│   │   ├── recursos.go
│   │   ├── mascota.go
│   │   ├── historial_medico.go
│   │   └── usuario.go
│   ├── storage/
│   │   ├── repositorio.go      — interfaces (ISP: 6 repos + UserRepository)
│   │   ├── gorm.go             — implementación GORM + AutoMigrate
│   │   ├── memoria.go          — fake RAM (usado en tests de handler)
│   │   ├── gorm_alerta_test.go
│   │   ├── gorm_asignacion_test.go
│   │   ├── gorm_recurso_test.go
│   │   ├── gorm_perfil_test.go
│   │   └── gorm_historial_test.go
│   ├── service/
│   │   ├── errores.go          — errores de dominio tipados
│   │   ├── alerta.go / alerta_test.go
│   │   ├── asignacion.go / asignacion_test.go
│   │   ├── perfil_veterinario.go / perfil_veterinario_test.go
│   │   ├── recurso.go / recurso_test.go
│   │   ├── mascota.go
│   │   ├── historial_medico.go / historial_medico_test.go
│   │   └── auth.go
│   ├── handlers/
│   │   ├── server.go           — struct Server + NewServer (7 servicios)
│   │   ├── respond.go          — helpers JSON + traducción de errores
│   │   ├── auth_handler.go
│   │   ├── alerta_handler.go / alerta_handler_test.go
│   │   ├── asignacion_handler.go / asignacion_handler_test.go
│   │   ├── perfil_veterinario_handler.go / perfil_veterinario_handler_test.go
│   │   ├── recurso_handler.go / recurso_handler_test.go
│   │   ├── mascota_handler.go
│   │   └── historial_medico_handler.go / historial_medico_handler_test.go
│   └── middleware/
│       ├── auth.go             — JWT Bearer → 401
│       └── cors.go
└── go.mod
```

---

## Ejecutar la API

```bash
go run ./cmd/main.go
# Servidor en http://localhost:8080
```

## Ejecutar todos los tests

```bash
go test ./...
```

---

## Endpoints

Todas las rutas protegidas requieren `Authorization: Bearer <token>`.

### Auth (público)

| Método | Ruta | Descripción |
|---|---|---|
| POST | `/api/v1/auth/register` | Registrar usuario |
| POST | `/api/v1/auth/login` | Login → devuelve JWT |

### Jean Carlos — Emergencia médica

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/api/v1/alertas` | Listar alertas |
| POST | `/api/v1/alertas` | Crear alerta (gravedad 1-5, requerimiento obligatorio) |
| GET | `/api/v1/alertas/{id}` | Obtener alerta |
| PUT | `/api/v1/alertas/{id}` | Actualizar alerta |
| DELETE | `/api/v1/alertas/{id}` | Borrar alerta |
| GET | `/api/v1/asignaciones` | Listar asignaciones de triage |
| POST | `/api/v1/asignaciones` | Crear asignación (alerta_id y recurso_id > 0) |
| GET | `/api/v1/asignaciones/{id}` | Obtener asignación |
| PUT | `/api/v1/asignaciones/{id}` | Actualizar asignación |
| DELETE | `/api/v1/asignaciones/{id}` | Borrar asignación |

### John Erick — Perfil veterinario

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/api/v1/veterinarios` | Listar perfiles veterinarios |
| POST | `/api/v1/veterinarios` | Crear perfil (nombre y telefono obligatorios) |
| GET | `/api/v1/veterinarios/{id}` | Obtener perfil |
| PUT | `/api/v1/veterinarios/{id}` | Actualizar perfil |
| DELETE | `/api/v1/veterinarios/{id}` | Borrar perfil |
| GET | `/api/v1/veterinarios/{id}/recursos` | Recursos del perfil |
| GET | `/api/v1/recursos` | Listar recursos clínicos |
| POST | `/api/v1/recursos` | Crear recurso (tipo_maquina obligatorio) |
| GET | `/api/v1/recursos/{id}` | Obtener recurso |
| PUT | `/api/v1/recursos/{id}` | Actualizar recurso |
| DELETE | `/api/v1/recursos/{id}` | Borrar recurso |

### María José — Historial médico

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/api/v1/mascotas` | Listar mascotas |
| POST | `/api/v1/mascotas` | Crear mascota (nombre obligatorio) |
| GET | `/api/v1/mascotas/{id}` | Obtener mascota |
| PUT | `/api/v1/mascotas/{id}` | Actualizar mascota |
| DELETE | `/api/v1/mascotas/{id}` | Borrar mascota |
| GET | `/api/v1/mascotas/{id}/historial` | Historial de la mascota |
| GET | `/api/v1/historial` | Listar todo el historial |
| POST | `/api/v1/historial` | Crear entrada (mascota_id > 0, diagnostico obligatorio) |
| GET | `/api/v1/historial/{id}` | Obtener entrada de historial |

---

## Tests por estudiante

### Jean Carlos Arauz — Módulo Emergencia médica

| Tipo | Archivo | Qué verifica |
|---|---|---|
| Service mock | `service/alerta_test.go` | `validarAlerta`: gravedad 0/6 → `ErrGravedadInvalida`, `AssertNotCalled("CrearAlerta")` |
| Service mock | `service/asignacion_test.go` | `validarAsignacion`: IDs ≤ 0 → `ErrIDsInvalidos`, `AssertNotCalled("CrearAsignacion")` |
| Handler httptest | `handlers/alerta_handler_test.go` | POST válido → 201; gravedad inválida → 400; sin token → **401** |
| Handler httptest | `handlers/asignacion_handler_test.go` | POST válido → 201; IDs inválidos → 400; sin token → **401** |
| Repository GORM | `storage/gorm_alerta_test.go` | `CrearAlerta` + `BuscarAlertaPorID` con SQLite `:memory:` |
| Repository GORM | `storage/gorm_asignacion_test.go` | `CrearAsignacion` + `BuscarAsignacionPorID` |

### John Erick Bello — Módulo Perfil veterinario

| Tipo | Archivo | Qué verifica |
|---|---|---|
| **Service mock** | `service/perfil_veterinario_test.go` | `validarPerfil`: nombre/teléfono vacíos → `ErrNombreVacio`/`ErrTelefonoVacio`, `AssertNotCalled("CrearPerfil")` |
| Service mock | `service/recurso_test.go` | `validarRecurso`: tipo_maquina vacío → `ErrTipoMaquinaVacio`, `AssertNotCalled` |
| **Handler httptest** | `handlers/perfil_veterinario_handler_test.go` | POST válido → 201; teléfono vacío → 400; sin token → **401** |
| Handler httptest | `handlers/recurso_handler_test.go` | POST válido → 201; tipo_maquina vacío → 400; sin token → 401 |
| **Repository GORM** | `storage/gorm_perfil_test.go` | `CrearPerfil` + `BuscarPerfilPorID` con SQLite `:memory:` |
| Repository GORM | `storage/gorm_recurso_test.go` | `CrearRecurso` + `BuscarRecursoPorID` |

### María José Vinces — Módulo Historial médico

| Tipo | Archivo | Qué verifica |
|---|---|---|
| **Service mock** | `service/historial_medico_test.go` | `validarHistorial`: mascota_id ≤ 0 / diagnóstico vacío → error, `AssertNotCalled("CrearHistorial")` |
| **Handler httptest** | `handlers/historial_medico_handler_test.go` | POST válido → 201; diagnóstico vacío → 400; sin token → **401** |
| **Repository GORM** | `storage/gorm_historial_test.go` | `CrearHistorial` + `BuscarHistorialPorID` + `ListarHistorialPorMascota` |

---

## Reglas de negocio

| Módulo | Entidad | Regla |
|---|---|---|
| Jean Carlos | `AlertaEmergencia` | `gravedad` ∈ [1, 5]; `requerimiento` no vacío |
| Jean Carlos | `AsignacionTriage` | `alerta_id` > 0 y `recurso_id` > 0 |
| John Erick | `PerfilVeterinario` | `nombre` no vacío; `telefono` no vacío |
| John Erick | `RecursoClinico` | `tipo_maquina` no vacío |
| María José | `Mascota` | `nombre` no vacío |
| María José | `HistorialMedico` | `mascota_id` > 0; `diagnostico` no vacío |

---

## Conceptos clave

| Término | Definición |
|---|---|
| **Mock** | Doble de prueba que registra y verifica llamadas. `AssertNotCalled` demuestra que la validación actúa ANTES del repositorio. |
| **Fake** | Implementación alternativa que guarda datos en RAM. No verifica llamadas. Se usa en tests de handler. |
| **httptest** | Paquete estándar de Go para probar handlers sin levantar un servidor real. |
| **:memory:** | DSN de SQLite que abre la base en RAM. Cada test tiene su base aislada. |
| **401** | Lo genera el middleware `Auth`, nunca el handler. El 401 demuestra que el JWT protege las rutas. |
| **AutoMigrate** | GORM crea/actualiza las tablas a partir de los structs. Sin ella el test de repositorio falla en `NotZero(ID)`. |

---

## Ejemplo de uso rápido

```bash
# 1. Registro y login
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"vet@codigogarra.vet","password":"secreto123"}'

TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"vet@codigogarra.vet","password":"secreto123"}' | jq -r .token)

# 2. Crear alerta de emergencia (Jean Carlos)
curl -X POST http://localhost:8080/api/v1/alertas \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"gravedad":4,"requerimiento":"Ventilador mecánico","estado":"Buscando"}'

# 3. Crear perfil veterinario (John Erick)
curl -X POST http://localhost:8080/api/v1/veterinarios \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Clínica Garra Norte","telefono":"0987654321","activo":true}'

# 4. Crear mascota y registrar historial (María José)
curl -X POST http://localhost:8080/api/v1/mascotas \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Simba","especie":"Gato","edad":3,"dueno":"Ana Pérez"}'

curl -X POST http://localhost:8080/api/v1/historial \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"mascota_id":1,"diagnostico":"Gastritis leve","tratamiento":"Dieta blanda","fecha":"2026-07-01"}'
```
