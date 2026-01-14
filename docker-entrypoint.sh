#!/bin/sh
set -e

# Build dockerize command based on enabled services
DOCKERIZE_CMD="dockerize"

# Wait for PostgreSQL if database is enabled
if [ "${ENABLE_DB:-true}" = "true" ] && [ -n "${DB_HOST}" ]; then
    DOCKERIZE_CMD="${DOCKERIZE_CMD} -wait tcp://${DB_HOST}:${DB_PORT:-5432}"
fi

# Wait for Redis if Redis is enabled
if [ "${ENABLE_REDIS:-true}" = "true" ] && [ -n "${REDIS_HOST}" ]; then
    DOCKERIZE_CMD="${DOCKERIZE_CMD} -wait tcp://${REDIS_HOST}:${REDIS_PORT:-6379}"
fi

# Wait for RabbitMQ if AMQP is enabled
if [ "${ENABLE_AMQP:-true}" = "true" ] && [ -n "${AMQP_HOST}" ]; then
    DOCKERIZE_CMD="${DOCKERIZE_CMD} -wait tcp://${AMQP_HOST}:${AMQP_PORT:-5672}"
fi

# Determine command to run (default: ./boilerblade, can be overridden)
APP_CMD="${1:-./boilerblade}"

# Add timeout and execute application
exec ${DOCKERIZE_CMD} -timeout 60s ${APP_CMD}
