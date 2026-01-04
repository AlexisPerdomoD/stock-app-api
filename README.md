# Stock App Api

## Description

Este proyecto es una **API REST para la gestión y consulta de información bursátil**, construida en **Go**, siguiendo una arquitectura por capas.

La aplicación permite:

- Autenticación y registro de usuarios.
- Consulta de mercados.
- Consulta de acciones (stocks).
- Gestión de acciones favoritas por usuario.
- Consulta de históricos, tendencias y estadísticas de acciones.
- Obtención de recomendaciones por acción.
- Ingesta periódica de datos desde múltiples fuentes externas.
- Ejecución de tareas programadas en segundo plano (scheduler).

La API utiliza **Gin** como framework HTTP, **CockroachDB** como base de datos, y **JWT** para autenticación.  
Las dependencias se inyectan explícitamente en el `main` para mantener claridad en la inicialización y facilitar la testabilidad.

---

## Arquitectura (alto nivel)

- **Application**
  - Casos de uso (use cases)
  - Servicios de dominio
- **Infrastructure**
  - Persistencia en CockroachDB
  - Handlers HTTP y middlewares
  - Integraciones con servicios externos
  - Scheduler
- **Delivery**
  - API REST (Gin)
  - Documentación Swagger (OpenAPI)

---

## API Documentation

La documentación de la API está generada con **Swagger (OpenAPI)**.

### Acceso a Swagger UI

Una vez que el servidor está corriendo:

```http
host:port/swagger/index.html
```

Ejemplo:

```http
http://localhost:3000/swagger/index.html
```

## Scheduler

La aplicación ejecuta tareas automáticas en segundo plano para la ingesta de datos de acciones:

- Soporte para múltiples fuentes de datos.
- Intervalo de ejecución por defecto: **15 minutos**.
- Timeout por tarea: **3 minutos**.
- El scheduler se inicia automáticamente al arrancar el servidor.

## Entry Point

El punto de entrada de la aplicación es:

```sh
cmd/server/main.go
```

Desde este archivo se inicializan:

- Variables de entorno.
- Conexión a la base de datos.
- Repositorios.
- Casos de uso.
- Handlers HTTP.
- Scheduler.
- Servidor HTTP.

## Dependencies

Este proyecto usa **Go 1.24** y depende principalmente de:

### Core

- `github.com/gin-gonic/gin`  
  Framework HTTP.
- `github.com/gin-contrib/cors`  
  Middleware CORS.
- `github.com/jmoiron/sqlx`  
  Acceso a base de datos SQL.
- `github.com/lib/pq`  
  Driver PostgreSQL.
- `github.com/go-playground/validator/v10`  
  Validación de structs.
- `github.com/golang-jwt/jwt/v5`  
  Autenticación JWT.
- `golang.org/x/crypto`  
  Utilidades criptográficas.

### Configuración y jobs

- `github.com/joho/godotenv`  
  Carga de variables de entorno desde `.env`.
- `github.com/robfig/cron/v3`  
  Tareas programadas (cron).

### Documentación (Swagger)

- `github.com/swaggo/swag`
- `github.com/swaggo/gin-swagger`
- `github.com/swaggo/files`

### Testing

- `github.com/stretchr/testify`  
  Utilidades para tests.

### Dependencias indirectas

Incluyen soporte para:

- Serialización JSON de alto rendimiento
- Parsing OpenAPI / Swagger
- Utilidades de sincronización, sistema y tooling (`golang.org/x/*`)
- YAML, Protobuf y helpers internos

### Dependencias instaladas en el host

- `migrate`: cli escrita en Go para migrar bases de datos.
- `swago`: generador de documentación para APIs REST.

## Environment Variables

Crea un archivo `.env` en la raíz del proyecto siguiendo el formato de `.env.example`. Asegúrate de completar los valores sensibles como contraseñas, llaves y secretos.

```env
# COCKROACHDB SETUP CR
CR_HOST=localhost
CR_PORT=26257
CR_USER=root
CR_PASSWORD=yourpassword
CR_DB=defaultdb
CR_SSL=disable

# DATA SOURCING
MAIN_SOURCE_STOCK_URI=https://your-api-url.com/path
MAIN_SOURCE_STOCK_KEY=your-api-key

# AUTH
SESSION_SECRET=your-session-secret

# SERVER
SERVER_PORT=3000
GIN_MODE=debug
```

Ver `go.mod` para detalles.

## Setup

### Requisitos

Asegúrate de tener instalado:

- **Go 1.24+**
- **Docker** y **Docker Compose**
- **golang-migrate**  
  👉 https://github.com/golang-migrate/migrate

Verificación rápida:

    go version
    docker --version
    docker compose version
    migrate -version

---

### Variables de entorno

Crea un archivo `.env` en la raíz del proyecto con las variables necesarias para CockroachDB:

    CR_USER=
    CR_PASSWORD=
    CR_HOST=
    CR_PORT=
    CR_DB=
    CR_SSL=

Estas variables son cargadas automáticamente por el `Makefile`.

---

### Base de datos local (CockroachDB)

Levantar la base de datos local:

    make local-up-db

Aplicar migraciones:

    make local-migrate-up

---

### Ejecutar la aplicación

Con la base de datos y migraciones listas:

    make local-start

---

### Poblar la base de datos

Ejecuta el seed inicial:

    make local-populate-db

---

### Testing

Levanta la base de datos de test, ejecuta migraciones y corre los tests:

    make test-local

---

### Comandos Docker / Producción

Aplicar migraciones usando las variables del `.env`:

    make migrate-up

Levantar el servidor:

    make start

---

### Notas

- El proyecto usa **CockroachDB** vía Docker para entornos local y test.
- Las migraciones se gestionan con **golang-migrate**.
- Los comandos del `Makefile` están pensados para entornos tipo Unix (Linux / macOS).

## Tests

Para ejecutar los tests, ejecutar `go test ./...` desde la raíz del proyecto.

Hecho con mucho ❤️ y Go 🐹
