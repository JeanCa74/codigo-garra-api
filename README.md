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
- **ORM**: GORM + SQLite (dev/tests) / PostgreSQL 16 (Docker)
- **Router**: chi v5
- **Autenticación**: JWT (`golang-jwt/jwt/v5`) + bcrypt — roles: `veterinario` | `admin`
- **Tests**: `testing` + `testify` (mock, assert, require) + `httptest`
- **CI**: GitHub Actions (build → vet → test con cobertura)
- **Docker**: multi-stage image (~15 MB) + docker-compose con PostgreSQL

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
│   ├── middleware/
│   │   ├── auth.go             — JWT Bearer → 401
│   │   └── cors.go
│   └── web/
│       ├── web.go              — sirve la interfaz web embebida (go:embed)
│       └── static/             — SPA de demostración (index.html, app.js, styles.css)
└── go.mod
```

---

## Inicio rápido con Docker

```bash
# Levanta PostgreSQL + API + seeder automático en un solo comando
docker-compose up --build

# La API queda disponible en http://localhost:8080
# Cuentas pre-cargadas por el seeder:
#   admin@codigogarra.vet / Admin123!   (rol: admin)
#   vet@codigogarra.vet   / Vet123!    (rol: veterinario)
```

## Ejecutar en local (SQLite)

```bash
go run ./cmd/main.go
# Servidor en http://localhost:8080, base de datos: garra.db
```

Variables de entorno — copiar `.env.example` a `.env` para personalizar.

## Interfaz web de demostración

Al levantar el servidor, la raíz `http://localhost:8080/` sirve una **SPA de
demostración** (HTML/CSS/JS puro, sin dependencias ni build) que consume la
propia API y permite ver cómo podría lucir la aplicación:

- **Login / registro** con las cuentas del seeder (botones de acceso rápido
  «Admin» y «Veterinario» en la pantalla de inicio).
- **Panel** con estadísticas en vivo (alertas activas, mascotas, clínicas,
  recursos disponibles) y lista de emergencias críticas (gravedad 4-5).
- **Pestañas CRUD** para los tres módulos: Alertas, Asignaciones, Clínicas,
  Recursos, Mascotas e Historial — con edición de estados en línea.

Decisiones de seguridad del front end:

| Medida | Detalle |
|---|---|
| Mismo origen | La SPA se embebe en el binario con `go:embed` y se sirve desde la API — sin CORS abierto ni servidor extra. |
| CSP estricta | `default-src 'self'; frame-ancestors 'none'` — sin scripts inline ni recursos externos, y no se puede embeber en iframes. |
| Anti-XSS | Todo dato de la API se pinta con `textContent`, nunca con `innerHTML`. |
| Token | El JWT vive en `sessionStorage` (se borra al cerrar la pestaña) y ante un `401` la sesión se cierra sola. |
| Roles | Los botones de eliminar solo se muestran al rol `admin`; el servidor sigue validando con `RequireRol` (403). La UI acompaña, nunca reemplaza, la autorización del backend. |

## Ejecutar todos los tests

```bash
go test ./...
# o con cobertura:
go test ./... -cover
```

## Roles y autorización

| Rol | Puede hacer |
|-----|------------|
| `veterinario` | GET, POST, PUT en todos los recursos |
| `admin` | Todo lo anterior + DELETE en todos los recursos |

El rol se incluye en el JWT al hacer login. Los DELETE requieren `rol: admin`.

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
| **403** | Lo genera `RequireRol("admin")`. Indica que el token es válido pero el rol es insuficiente. |
| **AutoMigrate** | GORM crea/actualiza las tablas a partir de los structs. Sin ella el test de repositorio falla en `NotZero(ID)`. |
| **Has-Many** | `PerfilVeterinario` → `[]RecursoClinico`; `Mascota` → `[]HistorialMedico`. GORM los precarga con `Preload`. |
| **Belongs-To** | `RecursoClinico` → `PerfilVeterinario`; `HistorialMedico` → `Mascota`; `AsignacionTriage` → `AlertaEmergencia` y `RecursoClinico`. |

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
