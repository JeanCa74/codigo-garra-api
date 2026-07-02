# Código Garra — API Backend

**Aplicaciones Web II (TDI-601) · Semana 11 · Actividad C1**

API REST para conectar dueños de mascotas en emergencias veterinarias nocturnas con las clínicas que disponen del equipo específico requerido en ese momento.

---

## Integrantes — Grupo L (Bloque 3)

| Integrante | Módulo | Entidad | Regla de negocio |
|------------|--------|---------|-----------------|
| **Jean Carlos Arauz** | Triage y Alertas | `AlertaEmergencia` | Gravedad fuera de [1,5] → `ErrGravedadInvalida` (400) |
| **John Erick Bello** | Capacidad Dinámica | `RecursoClinico` | `TipoMaquina` vacío → `ErrTipoMaquinaVacio` (400) |
| **María José Vinces** | Enrutamiento y Matching | `AsignacionTriage` | `AlertaID` o `RecursoID` = 0 → `ErrIDsInvalidos` (400) |

---

## Stack tecnológico

| Componente | Tecnología |
|-----------|-----------|
| Lenguaje | Go 1.22 |
| Router | Chi v5 |
| ORM | GORM |
| Base de datos | SQLite (dev / :memory: en tests) |
| Autenticación | JWT (golang-jwt/jwt/v5) + bcrypt |
| Testing | Testify (mock + assert + require) |

---

## Arquitectura: model → repository → service → handler

```
internal/
├── models/
│   ├── alerta.go          — AlertaEmergencia (GORM tags)
│   ├── recursos.go        — RecursoClinico   (GORM tags)
│   ├── asignacion.go      — AsignacionTriage (GORM tags)
│   └── usuario.go         — Usuario (auth)
├── storage/
│   ├── repositorio.go     — interfaces: AlertaRepository, RecursoRepository,
│   │                         AsignacionRepository, Almacen, UserRepository
│   ├── memoria.go         — fake en RAM (usado en tests de handler)
│   ├── gorm.go            — implementación GORM + SQLite + AutoMigrate
│   ├── gorm_alerta_test.go    ✅ Jean Carlos — repositorio contra SQLite :memory:
│   ├── gorm_recurso_test.go   ✅ John Erick  — repositorio contra SQLite :memory:
│   └── gorm_asignacion_test.go ✅ María José  — repositorio contra SQLite :memory:
├── service/
│   ├── errores.go         — errores de dominio compartidos
│   ├── alerta.go          — validarAlerta (escala de gravedad [1,5])
│   ├── alerta_test.go     ✅ Jean Carlos — service con mock
│   ├── recurso.go         — validarRecurso (tipo_maquina obligatorio)
│   ├── recurso_test.go    ✅ John Erick  — service con mock
│   ├── asignacion.go      — validarAsignacion (IDs > 0)
│   ├── asignacion_test.go ✅ María José  — service con mock
│   └── auth.go            — bcrypt + JWT
├── handlers/
│   ├── server.go          — Server con los tres servicios inyectados
│   ├── respond.go         — RespondJSON / RespondError / statusDeError
│   ├── alerta_handler.go          — CRUD de alertas
│   ├── alerta_handler_test.go     ✅ Jean Carlos — httptest + 401
│   ├── recurso_handler.go         — CRUD de recursos
│   ├── recurso_handler_test.go    ✅ John Erick  — httptest + 401
│   ├── asignacion_handler.go      — CRUD de asignaciones
│   ├── asignacion_handler_test.go ✅ María José  — httptest + 401
│   └── auth_handler.go    — register + login
└── middleware/
    ├── auth.go            — JWT middleware Bearer → produce el 401
    └── cors.go            — CORS para desarrollo
```

---

## Tests — Suite en verde (18 tests)

```bash
go test ./... -cover
```

```
ok  github.com/JeanCa74/codigo-garra-api/internal/handlers   coverage: 25.7%
ok  github.com/JeanCa74/codigo-garra-api/internal/service    coverage: 29.1%
ok  github.com/JeanCa74/codigo-garra-api/internal/storage    coverage: 17.3%
```

---

## Módulo Jean Carlos — Triage y Alertas

### Reglas de negocio

| Regla | Error | HTTP |
|-------|-------|------|
| `Requerimiento` vacío | `ErrRequerimientoVacio` | 400 |
| `Gravedad` fuera de [1,5] | `ErrGravedadInvalida` | 400 |
| No encontrada | `ErrNoEncontrado` | 404 |

### Test 1 — Service con mock (`service/alerta_test.go`)

**Qué comprueba:** Que `validarAlerta` rechaza gravedad fuera de la escala [1,5] **antes** de llamar al repositorio.

**Cómo funciona:**
- Se crea un `alertaRepoMock` (testify/mock) que registra llamadas.
- Se llama a `svc.Crear(...)` con `Gravedad: 0` → se espera `ErrGravedadInvalida`.
- `AssertNotCalled` verifica que el repo **nunca** recibió `CrearAlerta`.

**Qué se rompería:** Si se elimina la validación de gravedad, el mock recibe una llamada inesperada y el test falla.

### Test 2 — Handler con httptest (`handlers/alerta_handler_test.go`)

**Test 2a — `TestCrearAlerta_Exitosa`:** POST con token y gravedad válida → 201 Created.

**Test 2b — `TestRutaAlertas_SinToken` (el 401):** POST sin `Authorization` → 401 Unauthorized. Lo genera el middleware, nunca el handler.

**Qué se rompería:** Si se elimina `r.Use(middleware.Auth(...))`, la petición llega al handler y responde 201.

### Test 3 — Repositorio GORM :memory: (`storage/gorm_alerta_test.go`)

**Qué comprueba:** Que `CrearAlerta` persiste en SQLite y `BuscarAlertaPorID` lo refleja.

**Qué se rompería:** Si `AutoMigrate` no crea la tabla, el ID queda en 0 → falla `require.NotZero`.

---

## Módulo John Erick — Capacidad Dinámica (Recursos)

### Reglas de negocio

| Regla | Error | HTTP |
|-------|-------|------|
| `TipoMaquina` vacío | `ErrTipoMaquinaVacio` | 400 |
| No encontrado | `ErrNoEncontrado` | 404 |

### Test 1 — Service con mock (`service/recurso_test.go`)

**Qué comprueba:** Que `validarRecurso` rechaza `TipoMaquina` vacío **antes** de llamar al repositorio.

**Cómo funciona:**
- Se crea un `recursoRepoMock`.
- Se llama a `svc.Crear(...)` con `TipoMaquina: ""` → se espera `ErrTipoMaquinaVacio`.
- `AssertNotCalled` verifica que el repo **nunca** recibió `CrearRecurso`.

**Qué se rompería:** Sin la validación, el mock recibe la llamada y el test falla.

### Test 2 — Handler con httptest (`handlers/recurso_handler_test.go`)

**Test 2a — `TestCrearRecurso_Exitoso`:** POST con token y tipo_maquina válido → 201 Created.

**Test 2b — `TestRutaRecursos_SinToken` (el 401):** POST sin `Authorization` → 401 Unauthorized.

### Test 3 — Repositorio GORM :memory: (`storage/gorm_recurso_test.go`)

**Qué comprueba:** Que `CrearRecurso` persiste en SQLite y `BuscarRecursoPorID` lo refleja.

---

## Módulo María José — Enrutamiento y Matching (Asignaciones)

### Reglas de negocio

| Regla | Error | HTTP |
|-------|-------|------|
| `AlertaID` ≤ 0 o `RecursoID` ≤ 0 | `ErrIDsInvalidos` | 400 |
| No encontrada | `ErrNoEncontrado` | 404 |

### Test 1 — Service con mock (`service/asignacion_test.go`)

**Qué comprueba:** Que `validarAsignacion` rechaza IDs inválidos **antes** de llamar al repositorio.

**Cómo funciona:**
- Se crea un `asignacionRepoMock`.
- Se llama a `svc.Crear(...)` con `AlertaID: 0` → se espera `ErrIDsInvalidos`.
- `AssertNotCalled` verifica que el repo **nunca** recibió `CrearAsignacion`.

**Qué se rompería:** Sin la validación de IDs, el repo intentaría crear una asignación huérfana.

### Test 2 — Handler con httptest (`handlers/asignacion_handler_test.go`)

**Test 2a — `TestCrearAsignacion_Exitosa`:** POST con IDs válidos y token → 201 Created.

**Test 2b — `TestRutaAsignaciones_SinToken` (el 401):** POST sin `Authorization` → 401 Unauthorized.

### Test 3 — Repositorio GORM :memory: (`storage/gorm_asignacion_test.go`)

**Qué comprueba:** Que `CrearAsignacion` persiste en SQLite y `BuscarAsignacionPorID` lo refleja.

---

## Conceptos clave (glosario para la presentación)

| Término | Definición |
|---------|-----------|
| **Mock** | Doble de prueba que verifica que ciertas llamadas ocurrieron (o no). Usa testify/mock. |
| **Fake** | Implementación alternativa funcional (guarda datos en RAM). No verifica llamadas. |
| **httptest** | Paquete estándar de Go para probar handlers HTTP sin levantar un servidor real. |
| **:memory:** | DSN especial de SQLite que abre la base en RAM; se destruye al cerrar la conexión. |
| **401** | Lo genera el middleware Auth cuando el header Authorization está ausente o el token es inválido. El handler nunca lo produce. |
| **AutoMigrate** | GORM crea/actualiza las tablas automáticamente a partir de los structs Go. |

---

## Ejecutar el servidor

```bash
go run ./cmd/main.go
# → http://localhost:8080
# → crea codigogarra.db automáticamente
```

## Endpoints disponibles

| Método | Ruta | Auth | Módulo |
|--------|------|------|--------|
| POST | `/api/v1/auth/register` | No | — |
| POST | `/api/v1/auth/login` | No | — |
| GET | `/api/v1/alertas` | Bearer JWT | Jean Carlos |
| POST | `/api/v1/alertas` | Bearer JWT | Jean Carlos |
| GET | `/api/v1/alertas/{id}` | Bearer JWT | Jean Carlos |
| PUT | `/api/v1/alertas/{id}` | Bearer JWT | Jean Carlos |
| DELETE | `/api/v1/alertas/{id}` | Bearer JWT | Jean Carlos |
| GET | `/api/v1/recursos` | Bearer JWT | John Erick |
| POST | `/api/v1/recursos` | Bearer JWT | John Erick |
| GET | `/api/v1/recursos/{id}` | Bearer JWT | John Erick |
| PUT | `/api/v1/recursos/{id}` | Bearer JWT | John Erick |
| DELETE | `/api/v1/recursos/{id}` | Bearer JWT | John Erick |
| GET | `/api/v1/asignaciones` | Bearer JWT | María José |
| POST | `/api/v1/asignaciones` | Bearer JWT | María José |
| GET | `/api/v1/asignaciones/{id}` | Bearer JWT | María José |
| PUT | `/api/v1/asignaciones/{id}` | Bearer JWT | María José |
| DELETE | `/api/v1/asignaciones/{id}` | Bearer JWT | María José |
