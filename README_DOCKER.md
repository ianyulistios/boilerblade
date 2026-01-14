# Docker Setup Guide

This guide explains how to run Boilerblade using Docker and Docker Compose.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+

## Quick Start

### 1. Build and Run with Docker Compose

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop all services
docker-compose down
```

### 2. Access Services

Once all services are running:

- **Application**: http://localhost:3000
- **Swagger UI**: http://localhost:3000/swagger/index.html
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)

### 3. Database Connection

The application will automatically connect to:
- **PostgreSQL**: `postgres:5432`
- **Redis**: `redis:6379`
- **RabbitMQ**: `rabbitmq:5672`

## Docker Compose Services

### Services Included

1. **app** - Boilerblade application (port 3000)
2. **postgres** - PostgreSQL database (port 5432)
3. **redis** - Redis cache (port 6379)
4. **rabbitmq** - RabbitMQ message queue (ports 5672, 15672)

### Service Health Checks

All services include health checks to ensure they're ready before the application starts.

### Dockerize Integration

The application uses [dockerize](https://github.com/jwilder/dockerize) to wait for dependent services (PostgreSQL, Redis, RabbitMQ) to be ready before starting. This ensures:

- **PostgreSQL** is accepting connections on port 5432
- **Redis** is accepting connections on port 6379
- **RabbitMQ** is accepting connections on port 5672

The app will wait up to 60 seconds for all services to become available. If any service fails to start within the timeout, the container will exit with an error.

**Note:** If you disable certain services (e.g., `ENABLE_REDIS=false`), you may need to modify the Dockerfile CMD to remove the corresponding `-wait` flag, or the container will fail waiting for a service that's not needed.

## Development Mode

For development with hot reload:

```bash
# Install air (hot reload tool)
go install github.com/cosmtrek/air@latest

# Run with development override
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

Or use air directly (without Docker):

```bash
air
```

## Environment Variables

Environment variables are set in `docker-compose.yml`. To customize:

1. Create a `.env` file (or modify docker-compose.yml)
2. Update the environment section in `docker-compose.yml`

### Key Environment Variables

```env
# Application
MODE=development
FIBER_PORT=3000
SERVER_MODE=both

# Database
DB_TYPE=postgres
DB_HOST=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=boilerblade

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# RabbitMQ
AMQP_HOST=rabbitmq
AMQP_USER=guest
AMQP_PASSWORD=guest
```

## Docker Commands

### Build

```bash
# Build the application image
docker-compose build

# Build without cache
docker-compose build --no-cache
```

### Run

```bash
# Start all services in background
docker-compose up -d

# Start and view logs
docker-compose up

# Start specific service
docker-compose up postgres redis
```

### Logs

```bash
# View all logs
docker-compose logs

# View app logs
docker-compose logs app

# Follow logs
docker-compose logs -f app

# View last 100 lines
docker-compose logs --tail=100 app
```

### Stop and Cleanup

```bash
# Stop services (keeps volumes)
docker-compose stop

# Stop and remove containers
docker-compose down

# Stop and remove containers + volumes
docker-compose down -v

# Remove images
docker-compose down --rmi all
```

### Database Operations

```bash
# Access PostgreSQL shell
docker-compose exec postgres psql -U postgres -d boilerblade

# Backup database
docker-compose exec postgres pg_dump -U postgres boilerblade > backup.sql

# Restore database
docker-compose exec -T postgres psql -U postgres boilerblade < backup.sql
```

### Redis Operations

```bash
# Access Redis CLI
docker-compose exec redis redis-cli

# Flush all data
docker-compose exec redis redis-cli FLUSHALL
```

## Production Deployment

For production, consider:

1. **Use production Dockerfile** (multi-stage build)
2. **Set strong passwords** in environment variables
3. **Use secrets management** (Docker secrets, Kubernetes secrets, etc.)
4. **Enable SSL/TLS** for database connections
5. **Use managed services** for databases (RDS, ElastiCache, etc.)
6. **Set resource limits** in docker-compose.yml

### Production docker-compose Example

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
    restart: always
```

## Troubleshooting

### Application won't start

1. Check if all services are healthy:
   ```bash
   docker-compose ps
   ```

2. Check application logs:
   ```bash
   docker-compose logs app
   ```

3. Verify environment variables:
   ```bash
   docker-compose exec app env | grep DB_
   ```

### Database connection errors

1. Ensure PostgreSQL is healthy:
   ```bash
   docker-compose exec postgres pg_isready -U postgres
   ```

2. Check database logs:
   ```bash
   docker-compose logs postgres
   ```

### Port conflicts

If ports are already in use, modify ports in `docker-compose.yml`:

```yaml
ports:
  - "3001:3000"  # Use 3001 instead of 3000
```

## Volumes

Docker Compose creates persistent volumes for:

- `postgres_data` - PostgreSQL data
- `redis_data` - Redis data
- `rabbitmq_data` - RabbitMQ data

These volumes persist even when containers are removed (unless using `docker-compose down -v`).

## Network

All services are connected via the `boilerblade-network` bridge network, allowing them to communicate using service names as hostnames.

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [Redis Docker Image](https://hub.docker.com/_/redis)
- [RabbitMQ Docker Image](https://hub.docker.com/_/rabbitmq)
