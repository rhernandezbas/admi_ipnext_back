-- +goose Up

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS usuarios (
    id CHAR(36) PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    rol VARCHAR(20) NOT NULL DEFAULT 'sub-usuario',
    permisos JSON NOT NULL,
    avatar VARCHAR(500),
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transferencias (
    id CHAR(36) PRIMARY KEY,
    beneficiario VARCHAR(255) NOT NULL,
    cbu VARCHAR(22),
    alias VARCHAR(100),
    categoria VARCHAR(100) NOT NULL,
    monto DECIMAL(15,2) NOT NULL,
    moneda VARCHAR(10) NOT NULL DEFAULT 'ARS',
    fecha_pago DATETIME NOT NULL,
    fecha_vencimiento DATETIME,
    frecuencia VARCHAR(20) NOT NULL DEFAULT 'manual',
    estado VARCHAR(20) NOT NULL DEFAULT 'pendiente',
    metodo_pago VARCHAR(20) NOT NULL DEFAULT 'transferencia',
    notas TEXT,
    proveedor_id CHAR(36),
    creado_por CHAR(36) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_transferencias_estado ON transferencias (estado);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_transferencias_fecha_pago ON transferencias (fecha_pago);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_transferencias_categoria ON transferencias (categoria);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS empleados (
    id CHAR(36) PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    puesto VARCHAR(100) NOT NULL,
    area VARCHAR(100) NOT NULL,
    rol VARCHAR(50) NOT NULL,
    sueldo_bruto DECIMAL(15,2) NOT NULL,
    obra_social VARCHAR(100) NOT NULL DEFAULT '',
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    fecha_ingreso DATE NOT NULL,
    avatar VARCHAR(500),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS liquidaciones (
    id CHAR(36) PRIMARY KEY,
    empleado_id CHAR(36) NOT NULL,
    periodo CHAR(7) NOT NULL,
    sueldo_bruto DECIMAL(15,2) NOT NULL,
    deducciones DECIMAL(15,2) NOT NULL DEFAULT 0,
    neto_a_pagar DECIMAL(15,2) NOT NULL,
    estado VARCHAR(20) NOT NULL DEFAULT 'borrador',
    aprobado_por CHAR(36),
    fecha_pago DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS guardias (
    id CHAR(36) PRIMARY KEY,
    empleado_id CHAR(36) NOT NULL,
    fecha DATE NOT NULL,
    horas DECIMAL(5,2) NOT NULL,
    monto DECIMAL(15,2) NOT NULL,
    descripcion TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS compensaciones (
    id CHAR(36) PRIMARY KEY,
    empleado_id CHAR(36) NOT NULL,
    tipo VARCHAR(50) NOT NULL,
    monto DECIMAL(15,2) NOT NULL,
    descripcion TEXT,
    fecha DATE NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS proveedores (
    id CHAR(36) PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    cuit VARCHAR(20) NOT NULL DEFAULT '',
    cbu VARCHAR(22),
    alias VARCHAR(100),
    email VARCHAR(255),
    categoria VARCHAR(100) NOT NULL,
    sitio_web VARCHAR(500),
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS contratos_proveedor (
    id CHAR(36) PRIMARY KEY,
    proveedor_id CHAR(36) NOT NULL,
    codigo VARCHAR(50) NOT NULL UNIQUE,
    descripcion TEXT NOT NULL,
    monto_mensual DECIMAL(15,2) NOT NULL,
    vigencia_desde DATE NOT NULL,
    vigencia_hasta DATE NOT NULL,
    estado VARCHAR(20) NOT NULL DEFAULT 'activo',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS servicios (
    id CHAR(36) PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    tipo VARCHAR(30) NOT NULL,
    proveedor VARCHAR(255) NOT NULL,
    costo_mensual DECIMAL(15,2) NOT NULL,
    vto_factura DATETIME,
    renovacion DATETIME,
    estado VARCHAR(30) NOT NULL DEFAULT 'activo',
    metadata JSON,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS inmuebles (
    id CHAR(36) PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    direccion VARCHAR(500) NOT NULL,
    propietario VARCHAR(255) NOT NULL,
    uso VARCHAR(20) NOT NULL,
    alquiler_mensual DECIMAL(15,2) NOT NULL,
    cbu VARCHAR(22),
    alias VARCHAR(100),
    estado VARCHAR(20) NOT NULL DEFAULT 'pendiente',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS contratos_alquiler (
    id CHAR(36) PRIMARY KEY,
    inmueble_id CHAR(36) NOT NULL,
    vigencia_desde DATE NOT NULL,
    vigencia_hasta DATE NOT NULL,
    ajuste_frecuencia VARCHAR(50) NOT NULL,
    monto_mensual DECIMAL(15,2) NOT NULL,
    estado VARCHAR(20) NOT NULL DEFAULT 'vigente',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pagos_alquiler (
    id CHAR(36) PRIMARY KEY,
    inmueble_id CHAR(36) NOT NULL,
    periodo CHAR(7) NOT NULL,
    monto DECIMAL(15,2) NOT NULL,
    fecha_pago DATETIME,
    estado VARCHAR(20) NOT NULL DEFAULT 'pendiente',
    comprobante VARCHAR(500),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cuentas_bancarias (
    id CHAR(36) PRIMARY KEY,
    banco VARCHAR(255) NOT NULL,
    tipo_cuenta VARCHAR(30) NOT NULL,
    nro_cuenta VARCHAR(100) NOT NULL,
    cbu VARCHAR(22),
    cci VARCHAR(30),
    saldo_actual DECIMAL(15,2) NOT NULL DEFAULT 0,
    moneda VARCHAR(10) NOT NULL DEFAULT 'ARS',
    activa BOOLEAN NOT NULL DEFAULT TRUE,
    ultima_actualizacion DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS movimientos_bancarios (
    id CHAR(36) PRIMARY KEY,
    cuenta_id CHAR(36) NOT NULL,
    tipo VARCHAR(20) NOT NULL,
    monto DECIMAL(15,2) NOT NULL,
    descripcion TEXT NOT NULL,
    fecha DATETIME NOT NULL,
    conciliado BOOLEAN NOT NULL DEFAULT FALSE,
    referencia VARCHAR(255),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT IGNORE INTO usuarios (id, nombre, email, password, rol, permisos, activo, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Administrador',
    'admin@ipnext.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'admin',
    '{"dashboard":true,"transferencias":"escritura","nominas":"escritura","proveedores":"escritura","servicios":"escritura","alquileres":"escritura","tesoreria":"escritura","reportes":"escritura"}',
    TRUE,
    NOW(),
    NOW()
);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS movimientos_bancarios;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS cuentas_bancarias;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS pagos_alquiler;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS contratos_alquiler;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS inmuebles;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS servicios;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS contratos_proveedor;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS proveedores;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS compensaciones;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS guardias;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS liquidaciones;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS empleados;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS transferencias;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS usuarios;
-- +goose StatementEnd
