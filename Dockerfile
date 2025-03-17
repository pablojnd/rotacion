# Etapa de compilación
FROM golang:1.22 as build

WORKDIR /app

# Copiar y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar la aplicación con optimizaciones
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /rotacion-api .

# Etapa de ejecución - usando una imagen mínima en lugar de scratch
# para tener certificados SSL y otros archivos necesarios
FROM alpine:3.19

WORKDIR /app

# Instalar certificados CA y tzdata (necesarios para conexiones seguras y manejo de fechas)
RUN apk --no-cache add ca-certificates tzdata

# Copiar el binario compilado
COPY --from=build /rotacion-api /app/

# Copiar los archivos estáticos necesarios
COPY --from=build /app/static /app/static

# Configuración por defecto en caso de que no se proporcione .env
ENV DB_SERVER=localhost \
    DB_USER=sa \
    DB_PASSWORD="" \
    DB_NAME="" \
    DB_PORT=1433 \
    MYSQL_HOST=localhost \
    MYSQL_USER=root \
    MYSQL_PASSWORD="" \
    MYSQL_DATABASE="" \
    MYSQL_PORT=3306 \
    SERVER_PORT=8080

# Exponer el puerto de la aplicación de manera dinámica
# No usamos EXPOSE $SERVER_PORT porque se evalúa en tiempo de build,
# pero documentamos que el puerto predeterminado es 8080
EXPOSE 8080

# Comando de ejecución con entrypoint que permite configuración dinámica
# Usamos un script shell como punto de entrada para manejar la configuración dinámica
COPY docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh
ENTRYPOINT ["/app/docker-entrypoint.sh"]
