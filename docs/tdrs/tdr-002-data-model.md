# TDR-002: Modelo de Datos — MySQL

## Estado
Propuesto

## Contexto
Definir el schema de la base de datos MySQL alineado con los tipos TypeScript del frontend (TDR-002 frontend). Las entidades Go en `domain/` mapean directamente a estas tablas.

## Tablas

### usuarios
```sql
CREATE TABLE usuarios (
  id          CHAR(36)     PRIMARY KEY,          -- UUID
  nombre      VARCHAR(100) NOT NULL,
  email       VARCHAR(150) NOT NULL UNIQUE,
  password    VARCHAR(255) NOT NULL,             -- bcrypt hash
  rol         ENUM('admin','sub-usuario') NOT NULL DEFAULT 'sub-usuario',
  permisos    JSON         NOT NULL,             -- PermisosUsuario serializado
  avatar      VARCHAR(500) NULL,
  activo      TINYINT(1)   NOT NULL DEFAULT 1,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### transferencias
```sql
CREATE TABLE transferencias (
  id                CHAR(36)        PRIMARY KEY,
  beneficiario      VARCHAR(150)    NOT NULL,
  cbu               VARCHAR(22)     NULL,
  alias             VARCHAR(100)    NULL,
  categoria         VARCHAR(50)     NOT NULL,
  monto             DECIMAL(15,2)   NOT NULL,
  moneda            ENUM('ARS','USD') NOT NULL DEFAULT 'ARS',
  fecha_pago        DATE            NOT NULL,
  fecha_vencimiento DATE            NULL,
  frecuencia        ENUM('manual','mensual','semanal','quincenal','semestral','anual') NOT NULL DEFAULT 'manual',
  estado            ENUM('pendiente','pagado','vencido','programado','en_proceso') NOT NULL DEFAULT 'pendiente',
  metodo_pago       ENUM('transferencia','debito','efectivo','cheque') NOT NULL,
  notas             TEXT            NULL,
  proveedor_id      CHAR(36)        NULL REFERENCES proveedores(id),
  creado_por        CHAR(36)        NOT NULL REFERENCES usuarios(id),
  created_at        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_estado (estado),
  INDEX idx_fecha_pago (fecha_pago),
  INDEX idx_proveedor (proveedor_id)
);
```

### empleados
```sql
CREATE TABLE empleados (
  id             CHAR(36)      PRIMARY KEY,
  nombre         VARCHAR(100)  NOT NULL,
  puesto         VARCHAR(100)  NOT NULL,
  area           VARCHAR(100)  NOT NULL,
  rol            VARCHAR(100)  NOT NULL,
  sueldo_bruto   DECIMAL(15,2) NOT NULL,
  obra_social    VARCHAR(100)  NOT NULL,
  activo         TINYINT(1)    NOT NULL DEFAULT 1,
  fecha_ingreso  DATE          NOT NULL,
  avatar         VARCHAR(500)  NULL,
  created_at     DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at     DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### liquidaciones
```sql
CREATE TABLE liquidaciones (
  id             CHAR(36)      PRIMARY KEY,
  empleado_id    CHAR(36)      NOT NULL REFERENCES empleados(id),
  periodo        CHAR(7)       NOT NULL,           -- YYYY-MM
  sueldo_bruto   DECIMAL(15,2) NOT NULL,
  deducciones    DECIMAL(15,2) NOT NULL DEFAULT 0,
  neto_a_pagar   DECIMAL(15,2) NOT NULL,
  estado         ENUM('borrador','aprobada','pagada') NOT NULL DEFAULT 'borrador',
  aprobado_por   CHAR(36)      NULL REFERENCES usuarios(id),
  created_at     DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at     DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_empleado_periodo (empleado_id, periodo)
);
```

### guardias
```sql
CREATE TABLE guardias (
  id           CHAR(36)      PRIMARY KEY,
  empleado_id  CHAR(36)      NOT NULL REFERENCES empleados(id),
  fecha        DATE          NOT NULL,
  horas        DECIMAL(5,2)  NOT NULL,
  monto        DECIMAL(15,2) NOT NULL,
  notas        TEXT          NULL,
  created_at   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### compensaciones
```sql
CREATE TABLE compensaciones (
  id           CHAR(36)      PRIMARY KEY,
  empleado_id  CHAR(36)      NOT NULL REFERENCES empleados(id),
  tipo         ENUM('bono','adelanto','extra','otro') NOT NULL,
  monto        DECIMAL(15,2) NOT NULL,
  fecha        DATE          NOT NULL,
  descripcion  VARCHAR(255)  NULL,
  created_at   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### proveedores
```sql
CREATE TABLE proveedores (
  id         CHAR(36)     PRIMARY KEY,
  nombre     VARCHAR(150) NOT NULL,
  cuit       VARCHAR(13)  NOT NULL UNIQUE,
  cbu        VARCHAR(22)  NULL,
  alias      VARCHAR(100) NULL,
  email      VARCHAR(150) NULL,
  categoria  VARCHAR(50)  NOT NULL,
  sitio_web  VARCHAR(500) NULL,
  activo     TINYINT(1)   NOT NULL DEFAULT 1,
  created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### contratos_proveedores
```sql
CREATE TABLE contratos_proveedores (
  id              CHAR(36)      PRIMARY KEY,
  codigo          VARCHAR(30)   NOT NULL UNIQUE,   -- CTR-2024-001
  proveedor_id    CHAR(36)      NOT NULL REFERENCES proveedores(id),
  vigencia_desde  DATE          NOT NULL,
  vigencia_hasta  DATE          NOT NULL,
  monto_anual     DECIMAL(15,2) NOT NULL,
  renovacion      DATE          NULL,
  estado          ENUM('activo','proximo_a_vencer','vencido','en_proceso') NOT NULL DEFAULT 'activo',
  created_at      DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### servicios
```sql
CREATE TABLE servicios (
  id            CHAR(36)      PRIMARY KEY,
  nombre        VARCHAR(150)  NOT NULL,
  tipo          ENUM('internet','energia','seguridad','software','obra_social','seguro') NOT NULL,
  proveedor     VARCHAR(150)  NOT NULL,
  costo_mensual DECIMAL(15,2) NOT NULL,
  vto_factura   DATE          NULL,
  renovacion    DATE          NULL,
  estado        ENUM('activo','proximo_a_vencer','inactivo') NOT NULL DEFAULT 'activo',
  metadata      JSON          NULL,
  created_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### inmuebles
```sql
CREATE TABLE inmuebles (
  id               CHAR(36)      PRIMARY KEY,
  nombre           VARCHAR(150)  NOT NULL,
  direccion        VARCHAR(255)  NOT NULL,
  propietario      VARCHAR(150)  NOT NULL,
  uso              ENUM('nodo','oficina','deposito','otro') NOT NULL,
  alquiler_mensual DECIMAL(15,2) NOT NULL,
  cbu              VARCHAR(22)   NULL,
  alias            VARCHAR(100)  NULL,
  estado           ENUM('pagado','pendiente','vencido') NOT NULL DEFAULT 'pendiente',
  created_at       DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at       DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### contratos_alquiler
```sql
CREATE TABLE contratos_alquiler (
  id               CHAR(36)      PRIMARY KEY,
  inmueble_id      CHAR(36)      NOT NULL REFERENCES inmuebles(id),
  vigencia_desde   DATE          NOT NULL,
  vigencia_hasta   DATE          NOT NULL,
  ajuste_frecuencia VARCHAR(50)  NOT NULL,    -- "6 meses", "anual"
  monto_mensual    DECIMAL(15,2) NOT NULL,
  estado           ENUM('vigente','por_vencer','vencido') NOT NULL DEFAULT 'vigente',
  created_at       DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at       DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### pagos_alquiler
```sql
CREATE TABLE pagos_alquiler (
  id           CHAR(36)      PRIMARY KEY,
  inmueble_id  CHAR(36)      NOT NULL REFERENCES inmuebles(id),
  periodo      CHAR(7)       NOT NULL,      -- YYYY-MM
  monto        DECIMAL(15,2) NOT NULL,
  fecha_pago   DATE          NULL,
  estado       ENUM('pagado','pendiente','vencido') NOT NULL DEFAULT 'pendiente',
  comprobante  VARCHAR(500)  NULL,
  created_at   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### cuentas_bancarias
```sql
CREATE TABLE cuentas_bancarias (
  id                  CHAR(36)      PRIMARY KEY,
  banco               VARCHAR(100)  NOT NULL,
  tipo_cuenta         VARCHAR(50)   NOT NULL,
  nro_cuenta          VARCHAR(30)   NOT NULL,
  cbu                 VARCHAR(22)   NULL,
  cci                 VARCHAR(30)   NULL,
  saldo_actual        DECIMAL(15,2) NOT NULL DEFAULT 0,
  moneda              ENUM('ARS','USD') NOT NULL DEFAULT 'ARS',
  activa              TINYINT(1)    NOT NULL DEFAULT 1,
  ultima_actualizacion DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at          DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at          DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### movimientos_bancarios
```sql
CREATE TABLE movimientos_bancarios (
  id           CHAR(36)      PRIMARY KEY,
  cuenta_id    CHAR(36)      NOT NULL REFERENCES cuentas_bancarias(id),
  tipo         ENUM('ingreso','egreso') NOT NULL,
  monto        DECIMAL(15,2) NOT NULL,
  descripcion  VARCHAR(255)  NOT NULL,
  fecha        DATE          NOT NULL,
  conciliado   TINYINT(1)    NOT NULL DEFAULT 0,
  referencia   VARCHAR(100)  NULL,
  created_at   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_cuenta_fecha (cuenta_id, fecha)
);
```

## Tipos Go (domain entities)

Los structs en `internal/domain/` usan `snake_case` para GORM:

```go
type Transferencia struct {
    ID               string    `gorm:"primaryKey;type:char(36)"`
    Beneficiario     string    `gorm:"not null"`
    CBU              *string
    Alias            *string
    Categoria        string    `gorm:"not null"`
    Monto            float64   `gorm:"type:decimal(15,2);not null"`
    Moneda           string    `gorm:"type:enum('ARS','USD');default:ARS"`
    FechaPago        time.Time `gorm:"not null"`
    FechaVencimiento *time.Time
    Frecuencia       string    `gorm:"type:enum(...)"`
    Estado           string    `gorm:"type:enum(...)"`
    MetodoPago       string    `gorm:"type:enum(...)"`
    Notas            *string
    ProveedorID      *string
    CreadoPor        string
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

## Consecuencias
- Los IDs son UUIDs (CHAR 36) generados en la capa de aplicación Go.
- Los campos JSON (permisos, metadata) son flexibles pero deben validarse en la capa de aplicación.
- Los índices en `estado` y `fecha_pago` son críticos para las queries del dashboard.
- El campo `periodo` (CHAR 7, formato YYYY-MM) se usa en liquidaciones y pagos de alquiler para simplificar queries por mes.
