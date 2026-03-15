# TDR-002: Modelo de Datos — MySQL/MariaDB

## Estado
Implementado

## Contexto
Schema de la base de datos implementado en `migrations/001_initial_schema.sql`. Las entidades Go en `domain/` mapean a estas tablas a través de los modelos GORM en `infrastructure/persistence/repository/`.

## Notas de implementación

- Cada `CREATE TABLE` va en su propio bloque `-- +goose StatementBegin / StatementEnd` (requerimiento de MariaDB).
- Los IDs son UUIDs en formato `CHAR(36)`, generados en la capa de aplicación Go (`uuid.NewString()`).
- Las columnas de fecha usan `DATETIME` (no `DATE`) para mayor compatibilidad con GORM.

---

## Tablas

### usuarios
```sql
CREATE TABLE usuarios (
  id          CHAR(36)     PRIMARY KEY,
  nombre      VARCHAR(255) NOT NULL,
  email       VARCHAR(255) NOT NULL UNIQUE,
  password    VARCHAR(255) NOT NULL,           -- bcrypt hash
  rol         VARCHAR(20)  NOT NULL DEFAULT 'sub-usuario',
  permisos    JSON         NOT NULL,
  avatar      VARCHAR(500) NULL,
  activo      BOOLEAN      NOT NULL DEFAULT TRUE,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

**Usuario seed (admin):**
- email: `admin@ipnext.com` / password: `password`
- rol: `admin`, todos los permisos en `escritura`

---

### transferencias
```sql
CREATE TABLE transferencias (
  id                CHAR(36)      PRIMARY KEY,
  beneficiario      VARCHAR(255)  NOT NULL,
  cbu               VARCHAR(22)   NULL,
  alias             VARCHAR(100)  NULL,
  categoria         VARCHAR(100)  NOT NULL,
  monto             DECIMAL(15,2) NOT NULL,
  moneda            VARCHAR(10)   NOT NULL DEFAULT 'ARS',
  fecha_pago        DATETIME      NOT NULL,
  fecha_vencimiento DATETIME      NULL,
  frecuencia        VARCHAR(20)   NOT NULL DEFAULT 'manual',
  estado            VARCHAR(20)   NOT NULL DEFAULT 'pendiente',
  metodo_pago       VARCHAR(20)   NOT NULL DEFAULT 'transferencia',
  notas             TEXT          NULL,
  proveedor_id      CHAR(36)      NULL,
  creado_por        CHAR(36)      NOT NULL,
  created_at        DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at        DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_transferencias_estado (estado),
  INDEX idx_transferencias_fecha_pago (fecha_pago),
  INDEX idx_transferencias_categoria (categoria)
);
```

---

### empleados
```sql
CREATE TABLE empleados (
  id            CHAR(36)      PRIMARY KEY,
  nombre        VARCHAR(255)  NOT NULL,
  puesto        VARCHAR(100)  NOT NULL,
  area          VARCHAR(100)  NOT NULL,
  rol           VARCHAR(50)   NOT NULL,
  sueldo_bruto  DECIMAL(15,2) NOT NULL,
  obra_social   VARCHAR(100)  NOT NULL DEFAULT '',
  activo        BOOLEAN       NOT NULL DEFAULT TRUE,
  fecha_ingreso DATE          NOT NULL,
  avatar        VARCHAR(500)  NULL,
  created_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

---

### liquidaciones
```sql
CREATE TABLE liquidaciones (
  id           CHAR(36)      PRIMARY KEY,
  empleado_id  CHAR(36)      NOT NULL,
  periodo      CHAR(7)       NOT NULL,           -- YYYY-MM
  sueldo_bruto DECIMAL(15,2) NOT NULL,
  deducciones  DECIMAL(15,2) NOT NULL DEFAULT 0,
  neto_a_pagar DECIMAL(15,2) NOT NULL,
  estado       VARCHAR(20)   NOT NULL DEFAULT 'borrador',
  aprobado_por CHAR(36)      NULL,               -- ID del usuario que aprobó
  fecha_pago   DATETIME      NULL,
  created_at   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

---

### guardias
```sql
CREATE TABLE guardias (
  id          CHAR(36)      PRIMARY KEY,
  empleado_id CHAR(36)      NOT NULL,
  fecha       DATE          NOT NULL,
  horas       DECIMAL(5,2)  NOT NULL,
  monto       DECIMAL(15,2) NOT NULL,
  descripcion TEXT          NULL,                -- mapeado como "Notas" en el modelo Go
  created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

### compensaciones
```sql
CREATE TABLE compensaciones (
  id          CHAR(36)      PRIMARY KEY,
  empleado_id CHAR(36)      NOT NULL,
  tipo        VARCHAR(50)   NOT NULL,
  monto       DECIMAL(15,2) NOT NULL,
  descripcion TEXT          NULL,
  fecha       DATE          NOT NULL,
  created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

### proveedores
```sql
CREATE TABLE proveedores (
  id         CHAR(36)     PRIMARY KEY,
  nombre     VARCHAR(255) NOT NULL,
  cuit       VARCHAR(20)  NOT NULL DEFAULT '',
  cbu        VARCHAR(22)  NULL,
  alias      VARCHAR(100) NULL,
  email      VARCHAR(255) NULL,
  categoria  VARCHAR(100) NOT NULL,
  sitio_web  VARCHAR(500) NULL,
  activo     BOOLEAN      NOT NULL DEFAULT TRUE,
  created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

---

### contratos_proveedor
```sql
CREATE TABLE contratos_proveedor (
  id             CHAR(36)      PRIMARY KEY,
  proveedor_id   CHAR(36)      NOT NULL,
  codigo         VARCHAR(50)   NOT NULL UNIQUE,  -- auto-generado: CTR-2026-XXXX
  descripcion    TEXT          NOT NULL DEFAULT '',
  monto_mensual  DECIMAL(15,2) NOT NULL,         -- mapeado como "MontoAnual" en el modelo Go
  vigencia_desde DATE          NOT NULL,
  vigencia_hasta DATE          NOT NULL,
  estado         VARCHAR(20)   NOT NULL DEFAULT 'activo',
  created_at     DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at     DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

> **Nota:** El campo Go se llama `MontoAnual` pero mapea a la columna `monto_mensual` vía `gorm:"column:monto_mensual"`. El código auto-genera el `codigo` en formato `CTR-{año}-{seq}`.

---

### servicios
```sql
CREATE TABLE servicios (
  id            CHAR(36)      PRIMARY KEY,
  nombre        VARCHAR(255)  NOT NULL,
  tipo          VARCHAR(30)   NOT NULL,
  proveedor     VARCHAR(255)  NOT NULL,
  costo_mensual DECIMAL(15,2) NOT NULL,
  vto_factura   DATETIME      NULL,
  renovacion    DATETIME      NULL,
  estado        VARCHAR(30)   NOT NULL DEFAULT 'activo',
  metadata      JSON          NULL,
  created_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

---

### inmuebles
```sql
CREATE TABLE inmuebles (
  id               CHAR(36)      PRIMARY KEY,
  nombre           VARCHAR(255)  NOT NULL,
  direccion        VARCHAR(500)  NOT NULL,
  propietario      VARCHAR(255)  NOT NULL,
  uso              VARCHAR(20)   NOT NULL,
  alquiler_mensual DECIMAL(15,2) NOT NULL,
  cbu              VARCHAR(22)   NULL,
  alias            VARCHAR(100)  NULL,
  estado           VARCHAR(20)   NOT NULL DEFAULT 'pendiente',
  created_at       DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at       DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

---

### contratos_alquiler
```sql
CREATE TABLE contratos_alquiler (
  id                CHAR(36)      PRIMARY KEY,
  inmueble_id       CHAR(36)      NOT NULL,
  vigencia_desde    DATE          NOT NULL,
  vigencia_hasta    DATE          NOT NULL,
  ajuste_frecuencia VARCHAR(50)   NOT NULL,
  monto_mensual     DECIMAL(15,2) NOT NULL,
  estado            VARCHAR(20)   NOT NULL DEFAULT 'vigente',
  created_at        DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at        DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

---

### pagos_alquiler
```sql
CREATE TABLE pagos_alquiler (
  id          CHAR(36)      PRIMARY KEY,
  inmueble_id CHAR(36)      NOT NULL,
  periodo     CHAR(7)       NOT NULL,            -- YYYY-MM
  monto       DECIMAL(15,2) NOT NULL,
  fecha_pago  DATETIME      NULL,
  estado      VARCHAR(20)   NOT NULL DEFAULT 'pendiente',
  comprobante VARCHAR(500)  NULL,
  created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

### cuentas_bancarias
```sql
CREATE TABLE cuentas_bancarias (
  id                   CHAR(36)      PRIMARY KEY,
  banco                VARCHAR(255)  NOT NULL,
  tipo_cuenta          VARCHAR(30)   NOT NULL,
  nro_cuenta           VARCHAR(100)  NOT NULL,
  cbu                  VARCHAR(22)   NULL,
  cci                  VARCHAR(30)   NULL,
  saldo_actual         DECIMAL(15,2) NOT NULL DEFAULT 0,
  moneda               VARCHAR(10)   NOT NULL DEFAULT 'ARS',
  activa               BOOLEAN       NOT NULL DEFAULT TRUE,
  ultima_actualizacion DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at           DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at           DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

---

### movimientos_bancarios
```sql
CREATE TABLE movimientos_bancarios (
  id          CHAR(36)      PRIMARY KEY,
  cuenta_id   CHAR(36)      NOT NULL,
  tipo        VARCHAR(20)   NOT NULL,            -- 'ingreso' | 'egreso'
  monto       DECIMAL(15,2) NOT NULL,
  descripcion TEXT          NOT NULL,
  fecha       DATETIME      NOT NULL,
  conciliado  BOOLEAN       NOT NULL DEFAULT FALSE,
  referencia  VARCHAR(255)  NULL,
  created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

## Consecuencias
- Los IDs son UUIDs (`CHAR(36)`) generados en la capa de aplicación Go.
- Los campos JSON (`permisos`, `metadata`) son flexibles pero validados en la capa de aplicación.
- El campo `periodo` (`CHAR(7)`, formato `YYYY-MM`) se usa en liquidaciones y pagos de alquiler.
- Las fechas en request body se reciben como `YYYY-MM-DD`; los handlers parsean a `time.Time`.
- Los campos marcados con `gorm:"column:..."` tienen discrepancia entre nombre Go y columna DB — documentado en cada tabla.
