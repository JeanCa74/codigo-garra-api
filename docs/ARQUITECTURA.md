# Arquitectura — Código Garra API

## Diagrama de capas

```
┌─────────────────────────────────────────────────────────────────────┐
│                        CLIENTE (Postman / Frontend)                  │
└────────────────────────────┬────────────────────────────────────────┘
                             │ HTTP / JSON
┌────────────────────────────▼────────────────────────────────────────┐
│                          MIDDLEWARE                                   │
│   CORS ──► Logger ──► Recoverer ──► Auth (JWT) ──► RequireRol(admin)│
└────────────────────────────┬────────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────────┐
│                        CAPA DE HANDLERS                               │
│                    (internal/handlers/*.go)                           │
│                                                                       │
│  AlertaHandler  AsignacionHandler  PerfilHandler  RecursoHandler     │
│  MascotaHandler  HistorialHandler  AuthHandler                       │
│                                                                       │
│            Todos accedidos vía ServerDeps struct                     │
└────────────────────────────┬────────────────────────────────────────┘
                             │ Interfaces ISP
┌────────────────────────────▼────────────────────────────────────────┐
│                        CAPA DE SERVICIOS                              │
│                    (internal/service/*.go)                            │
│                                                                       │
│  AlertaService   AsignacionService   PerfilVeterinarioService        │
│  RecursoService  MascotaService      HistorialMedicoService          │
│  AuthService                                                          │
│                                                                       │
│   Validaciones de negocio: gravedad 1–5, campos obligatorios, etc.  │
└────────────────────────────┬────────────────────────────────────────┘
                             │ Interfaces de repositorio
┌────────────────────────────▼────────────────────────────────────────┐
│                      CAPA DE ALMACENAMIENTO                           │
│                    (internal/storage/*.go)                            │
│                                                                       │
│  AlmacenGORM (SQLite / PostgreSQL)   Memoria (fake para tests)      │
│                                                                       │
│  AutoMigrate → relaciones GORM:                                      │
│    PerfilVeterinario ──has-many──► RecursoClinico                    │
│    Mascota           ──has-many──► HistorialMedico                   │
│    AsignacionTriage  ──belongs-to─ AlertaEmergencia                  │
│    AsignacionTriage  ──belongs-to─ RecursoClinico                    │
│    RecursoClinico    ──belongs-to─ PerfilVeterinario                 │
│    HistorialMedico   ──belongs-to─ Mascota                           │
└────────────────────────────┬────────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────────┐
│              BASE DE DATOS                                            │
│   SQLite (dev / tests ":memory:")   PostgreSQL 16 (Docker)           │
└─────────────────────────────────────────────────────────────────────┘
```

## Módulos del proyecto

| Módulo | Integrante | Entidades | Endpoints |
|--------|-----------|-----------|-----------|
| Emergencia médica | Jean Carlos | AlertaEmergencia, AsignacionTriage | `/alertas`, `/asignaciones` |
| Perfil veterinario | John Erick | PerfilVeterinario, RecursoClinico | `/veterinarios`, `/recursos` |
| Historial médico | María José | Mascota, HistorialMedico | `/mascotas`, `/historial` |
| Auth compartido | Grupo | Usuario | `/auth/register`, `/auth/login` |

## Roles y autorización

| Rol | Puede hacer |
|-----|------------|
| `veterinario` | CRUD completo excepto DELETE |
| `admin` | Todo, incluidos los DELETE |

El rol se incluye en el JWT como claim `"rol"` al hacer login. El middleware
`RequireRol("admin")` protege los endpoints de eliminación.

## Flujo de autenticación

```
POST /auth/register  →  bcrypt hash  →  INSERT usuario (rol=veterinario)
POST /auth/login     →  validar hash  →  JWT con {sub, rol, exp}
GET  /alertas        →  Authorization: Bearer <token>  →  Auth middleware
DELETE /alertas/1    →  Auth middleware + RequireRol("admin")
```

## Infraestructura Docker

```
docker-compose up
       │
       ├── db (postgres:16-alpine)
       │      └── healthcheck: pg_isready
       │
       └── api (multi-stage Dockerfile)
              ├── depends_on: db (service_healthy)
              ├── DATABASE_URL → NuevoAlmacenPostgres()
              └── SEED_ON_STARTUP=true → seed.Sembrar()
```
