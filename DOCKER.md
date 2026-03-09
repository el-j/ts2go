# Docker Compose Setup for TS2Go

This directory contains the Docker Compose configuration for running the TS2Go SaaS platform locally.

## Quick Start

1. **Copy the environment template:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` and set secure passwords:**
   ```bash
   # Update these critical values in .env:
   POSTGRES_PASSWORD=<your-secure-postgres-password>
   REDIS_PASSWORD=<your-secure-redis-password>
   MINIO_ROOT_PASSWORD=<your-secure-minio-password>
   JWT_SECRET=<your-64-character-random-string>
   ```

3. **Start the services:**
   ```bash
   docker-compose up -d
   ```

4. **Check the status:**
   ```bash
   docker-compose ps
   ```

## Services

The stack includes the following services:

- **PostgreSQL** (port 5432): Main database
- **Redis** (port 6379): Session store and job queue
- **MinIO** (ports 9000, 9001): S3-compatible object storage
- **API** (port 8080): Go backend API server
- **Worker**: Background job processor for transpilation
- **Frontend** (port 3000): Next.js web interface

## Environment Variables

All sensitive configuration is managed through environment variables. **Never commit your `.env` file to version control.**

### Required Variables

These must be set in your `.env` file:
- `POSTGRES_PASSWORD` - PostgreSQL database password
- `REDIS_PASSWORD` - Redis password
- `MINIO_ROOT_PASSWORD` - MinIO root password
- `JWT_SECRET` - Secret key for JWT token signing (minimum 64 characters)

### Optional Variables

These have defaults but can be customized:
- `POSTGRES_USER` - PostgreSQL username (default: `ts2go`)
- `POSTGRES_DB` - Database name (default: `ts2go_saas`)
- `MINIO_ROOT_USER` - MinIO username (default: `ts2go`)
- `API_PORT` - API server port (default: `8080`)
- `ENV` - Environment mode (default: `development`)
- `LOG_LEVEL` - Logging level (default: `debug`)

## Security Considerations

### Development

For development, you can use simple passwords, but ensure:
1. Your `.env` file is in `.gitignore`
2. You don't commit credentials to version control

### Production

For production deployments:
1. Generate strong random passwords for all password fields (at least 32 characters)
2. Use a cryptographically secure random string for `JWT_SECRET` (at least 64 characters)
3. Set `ENV=production`
4. Set `LOG_LEVEL=info` or `warn`
5. Use Docker secrets or a secret management service instead of `.env` files
6. Enable SSL/TLS for all services
7. Restrict network access using firewall rules

## Generating Secure Passwords

You can generate secure passwords using:

```bash
# For passwords
openssl rand -base64 32

# For JWT secret (64 characters)
openssl rand -base64 64
```

## Common Commands

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes data)
docker-compose down -v

# Rebuild services
docker-compose build

# Restart a specific service
docker-compose restart api
```

## Troubleshooting

### Services won't start

Check if required environment variables are set:
```bash
docker-compose config
```

### Database connection errors

Ensure PostgreSQL is healthy:
```bash
docker-compose ps postgres
docker-compose logs postgres
```

### Storage issues

Check MinIO status:
```bash
docker-compose logs minio
```

Access MinIO console at http://localhost:9001 with credentials from `.env`

## Development Tips

1. **Hot Reload**: Uncomment the `command: air -c .air.toml` line in `docker-compose.yml` for Go API hot reload
2. **Volume Mounts**: Frontend and API source code is mounted for live development
3. **Database Migrations**: Run migrations manually or through the API startup sequence

## License

See the main repository LICENSE file.
