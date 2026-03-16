# TDR-005: Cargas Sociales por Empleado

## Estado
Pendiente de implementación

## Contexto
El frontend actualmente calcula las cargas sociales como 30% fijo global.
Se requiere que cada empleado tenga configuración propia: porcentaje y/o monto fijo.
Si se especifica monto fijo, tiene precedencia sobre el porcentaje.

---

## Cambio en schema — tabla `empleados`

```sql
ALTER TABLE empleados
  ADD COLUMN cargas_sociales_pct   DECIMAL(5,2)  NOT NULL DEFAULT 30.00
    COMMENT 'Porcentaje de cargas sociales sobre sueldo bruto (ej: 30.00 = 30%)',
  ADD COLUMN cargas_sociales_monto DECIMAL(15,2) NULL DEFAULT NULL
    COMMENT 'Monto fijo de cargas sociales. Si no es NULL, tiene precedencia sobre el porcentaje.';
```

Agregar migración correspondiente en `migrations/`.

---

## Regla de cálculo

```
si cargas_sociales_monto IS NOT NULL:
    cargas_sociales_calculado = cargas_sociales_monto
sino:
    cargas_sociales_calculado = sueldo_bruto * cargas_sociales_pct / 100
```

---

## Cambios en la API

### GET /nominas/empleados y GET /nominas/empleados/:id

Agregar campos al response:

```json
{
  "ID": "uuid",
  "Nombre": "Juan Perez",
  "SueldoBruto": 150000,
  "CargasSocialesPct": 30.00,
  "CargasSocialesMonto": null,
  "CargasSocialesCalculado": 45000
}
```

`CargasSocialesCalculado` se calcula en el handler antes de serializar, no se persiste.

### POST /nominas/empleados

Agregar campos opcionales al body:

```json
{
  "nombre": "string (required)",
  "puesto": "string (required)",
  "area": "string (required)",
  "rol": "string (required)",
  "sueldoBruto": 150000,
  "obraSocial": "string (required)",
  "fechaIngreso": "YYYY-MM-DD (required)",
  "cargasSocialesPct": 30.0,
  "cargasSocialesMonto": null
}
```

Si no se envían, usar defaults: `cargasSocialesPct=30.00`, `cargasSocialesMonto=NULL`.

### PATCH /nominas/empleados/:id

Agregar campos opcionales al body:

```json
{
  "cargasSocialesPct": 35.5,
  "cargasSocialesMonto": null
}
```

Para borrar el monto fijo y volver a porcentaje, enviar `"cargasSocialesMonto": null` explícitamente.

---

## Notas de implementación

1. Crear migración: `migrations/00X_add_cargas_sociales_to_empleados.sql`
2. Actualizar struct `Empleado` en `internal/models/` con los nuevos campos
3. Actualizar handler GET para calcular y devolver `CargasSocialesCalculado`
4. Actualizar handlers POST y PATCH para aceptar y persistir los nuevos campos
5. El campo `CargasSocialesCalculado` NO se guarda en DB, se calcula on-the-fly
