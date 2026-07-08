# CIERRE — Código Garra API (H3)
**TDI-601 Aplicaciones Web II · Semana 11**

---

## ¿Qué construimos?

Plataforma de emergencias veterinarias con API REST en Go que conecta
mascotas en estado crítico con clínicas y recursos disponibles.

### Módulos (3 × 2 entidades)

| # | Módulo | Entidades | Integrante |
|---|--------|-----------|-----------|
| 1 | Emergencia médica | AlertaEmergencia, AsignacionTriage | Jean Carlos |
| 2 | Perfil veterinario | PerfilVeterinario, RecursoClinico | John Erick |
| 3 | Historial médico | Mascota, HistorialMedico | María José |

---

## Lo que aprendimos

### Arquitectura en capas
Separar `handlers → services → storage` hace que cada capa sea testeable
en forma independiente. Los handlers solo hablan HTTP; los servicios solo
conocen lógica de negocio; el almacén solo habla GORM.

### GORM y relaciones entre entidades
- **Has-Many**: `PerfilVeterinario → []RecursoClinico`,
  `Mascota → []HistorialMedico`
- **Belongs-To**: `RecursoClinico → PerfilVeterinario`,
  `HistorialMedico → Mascota`,
  `AsignacionTriage → AlertaEmergencia` y `→ RecursoClinico`
- `Preload("Recursos")` permite devolver los recursos anidados en la
  respuesta del perfil sin queries adicionales.

### JWT con roles
Incluir `"rol"` en los claims del JWT permite proteger rutas sensibles
(DELETE) sin consultar la base de datos en cada request. El middleware
`RequireRol("admin")` lee el claim del contexto en O(1).

### Inyección de dependencias con `ServerDeps`
Reemplazar 7 parámetros posicionales por un struct elimina el "constructor
hell" y hace los tests más legibles: cada test inicializa solo los campos
que necesita.

### Tests con mocks (testify/mock)
Los tests de servicio usan fakes en memoria. Los de handler levantan un
router chi completo con el middleware real, lo que prueba la integración
HTTP sin levantar una base de datos.

### Docker multi-stage
La imagen final pesa ~15 MB (Alpine + binario estático). El stage builder
compila con `CGO_ENABLED=0` para que el binario sea portable.

---

## ¿Qué haríamos diferente?

1. **Migraciones con versión explícita** — GORM `AutoMigrate` es conveniente
   para desarrollo pero en producción se usan herramientas como `goose` o
   `atlas` para tener historial de cambios reversibles.

2. **Variables de entorno para el secreto JWT** — `secretoJWT` está
   hardcoded. En producción se debe leer de una variable de entorno o de
   un secret manager.

3. **Paginación desde el inicio** — Los endpoints `GET /alertas` devuelven
   toda la tabla. Con muchos registros esto escala mal; habríamos añadido
   `?page=&limit=` desde la primera iteración.

4. **Logging estructurado** — Reemplazar `log.Println` por `slog` (stdlib
   desde Go 1.21) para tener logs en JSON que se puedan filtrar en producción.

5. **Tests de integración con testcontainers** — Levantar un PostgreSQL real
   durante los tests de repositorio habría detectado diferencias entre
   SQLite y PostgreSQL (e.g., comportamiento de `ILIKE` vs `LIKE`).

---

## Próximos pasos

- [ ] Paginación en listados (`?page=1&limit=20`)
- [ ] Endpoint de búsqueda `/alertas?gravedad=5&estado=Buscando`
- [ ] Refresh token (el JWT actual expira en 24 h sin posibilidad de renovar)
- [ ] WebSocket para notificaciones en tiempo real cuando se crea una alerta
- [ ] RBAC más granular con tabla de permisos en la BD
- [ ] Métricas con Prometheus y dashboard en Grafana

---

*Proyecto desarrollado como entrega H3 · julio 2026*
