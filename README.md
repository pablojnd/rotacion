# API de Rotación

API para consulta de rotación de inventario y ventas.

## Ejecución con Docker

### Requisitos previos
- Docker
- Docker Compose

### Configuración
1. Copie el archivo `.env.example` a `.env`
2. Edite el archivo `.env` y configure las variables de conexión a las bases de datos
3. Puede cambiar el puerto de la aplicación modificando la variable `SERVER_PORT` en el archivo `.env`

### Construcción y ejecución
```bash
# Construir la imagen Docker
docker-compose build

# Iniciar el contenedor
docker-compose up -d

# Ver logs
docker-compose logs -f
```

### Acceso
La API estará disponible en http://localhost:${SERVER_PORT} (por defecto: puerto 8080)
La documentación está disponible en http://localhost:${SERVER_PORT}/docs

## Desarrollo local

### Requisitos previos
- Go 1.21 o superior
- Base de datos SQL Server
- Base de datos MySQL

### Configuración
1. Copie el archivo `.env.example` a `.env`
2. Edite el archivo `.env` y configure las variables de conexión a las bases de datos

### Ejecución
```bash
# Instalar dependencias
go mod download

# Compilar y ejecutar
go run .
```

## Conexiones externas a bases de datos

Para conectarse a bases de datos externas, configure las siguientes variables en el archivo `.env`:

- **SQL Server**: Configurar `DB_SERVER`, `DB_PORT`, `DB_USER`, `DB_PASSWORD` y `DB_NAME`
- **MySQL**: Configurar `MYSQL_HOST`, `MYSQL_PORT`, `MYSQL_USER`, `MYSQL_PASSWORD` y `MYSQL_DATABASE`

## API Endpoints

Consulte la documentación en http://localhost:${SERVER_PORT}/docs para ver todos los endpoints disponibles.
