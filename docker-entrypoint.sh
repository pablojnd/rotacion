#!/bin/sh
set -e

echo "Iniciando API de Rotación..."

# Verificar variables de entorno necesarias
if [ -z "$DB_SERVER" ] || [ -z "$MYSQL_HOST" ]; then
    echo "ADVERTENCIA: Faltan variables de entorno importantes. Verifique la configuración."
fi

# Mostrar información de configuración
echo "Configuración de servidor:"
echo "- Puerto: $SERVER_PORT"
echo "- SQL Server: $DB_SERVER:$DB_PORT"
echo "- MySQL: $MYSQL_HOST:$MYSQL_PORT"

# Iniciar la aplicación
exec /app/rotacion-api
