# IPNEXT Backend — Ralph Loop Spec

## Objetivo

Construir el backend Go de IPNEXT Admin: una API REST con Go + Gin + GORM + MySQL, arquitectura hexagonal, que sirve al frontend React de administración empresarial.

## Contexto

- Docs del proyecto: `docs/` (leer ANTES de implementar cualquier módulo)
- Stack: Go 1.22, Gin, GORM, MySQL 8
- Arquitectura: Hexagonal (ports & adapters)
- Convenciones: ver `docs/adrs/` y `docs/tdrs/`

## Instrucciones por iteración

En cada iteración debés:

1. **Leer `TRACKER.md`** — identificar cuál es la próxima tarea con estado `[ ]`
2. **Leer los docs relevantes** en `docs/` para esa tarea
3. **Implementar una sola tarea** del tracker (no más de una por iteración)
4. **Marcar la tarea como completada** en `TRACKER.md` → cambiar `[ ]` por `[x]`
5. **Actualizar la sección "Última iteración"** en `TRACKER.md` con qué se hizo y qué sigue
6. **Correr `go build ./...`** para verificar que compila sin errores
7. Si hay errores de compilación, corregirlos antes de finalizar la iteración
8. Si todas las tareas del tracker están `[x]`, emitir la promise de completitud

## Reglas

- Seguir estrictamente la estructura de carpetas definida en `docs/tdrs/tdr-001-project-structure.md`
- Los tipos Go de dominio deben alinearse con `docs/tdrs/tdr-002-data-model.md`
- Los endpoints deben seguir `docs/tdrs/tdr-003-endpoints.md`
- La autenticación debe seguir `docs/adrs/003-auth-strategy.md`
- Los IDs siempre son UUIDs generados con `github.com/google/uuid`
- Manejo de errores explícito — no panic, siempre retornar `error`
- Cada handler debe tener su propio archivo en `internal/infrastructure/http/handler/`
- Los repositorios implementan las interfaces definidas en `internal/domain/`

## Señal de completitud

Cuando TODAS las tareas del `TRACKER.md` estén marcadas `[x]`, emitir exactamente:

<promise>IPNEXT BACKEND COMPLETE</promise>
