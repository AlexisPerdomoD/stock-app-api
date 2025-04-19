# Stock App Api

Este es un API REST siguiendo de patrones DDD con un enfoque minimalista sobre una app para consultar y sugerir Stocks.

## Features

- Manejo de session básico a través de JWT

- Registro de usuarios

- Registro de stocks a usuarios

- Detalles de recomendaciones de Stocks por `Brokerages`

- Información detallada de los stocks incluyendo `Markets`, `Companies`, y `Brokerages` etc.

- Paginación de resultados tanto para consultas de `Stocks` como de `Recomendaciones`.

- Registro automático de `Stocks` a través de una API de fuentes de datos cada 24 horas.

- Posibilidad de agregar diferentes fuentes de datos de `Stocks` a futuro.

## Dependencies

Este proyecto usa Go 1.23 y depende de:

- `golang.org/x/crypto`: utilidades criptográficas.
- `gorm.io/gorm` + `gorm.io/driver/postgres`: ORM para CockroachDB.
- `github.com/robfig/cron/v3`: tareas programadas (cron).
- `github.com/joho/godotenv`: carga de variables de entorno desde `.env`.
- `github.com/stretchr/testify` : utilidades para tests.

Dependencias adicionales incluyen utilidades para fechas, sincronización, y soporte para PostgreSQL (`pgx`, `puddle`, etc.).

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
CR_RUN_MIGRATE=TRUE

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

Hecho con mucho ❤️ y Go 🐹
