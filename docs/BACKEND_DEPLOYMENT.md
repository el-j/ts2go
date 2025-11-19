# Backend Deployment Strategy

## Overview

This document outlines deployment options for the TS2Go SaaS backend API. The backend is currently located in `saas/backend/` and includes:

- **API**: Go-based REST API (Fiber framework)
- **Database**: PostgreSQL
- **Caching**: Redis
- **Authentication**: JWT-based auth system

## Deployment Options

### Option 1: Uberspace (Recommended for Testing)

**Pros:**
- Fixed cost: €5/month (flat rate, unlimited usage)
- Full SSH access and control
- Easy PostgreSQL and Redis setup
- Supervisord for process management
- No cold starts or function limitations
- Perfect for staging/testing environments

**Cons:**
- Manual server management
- Single server (no auto-scaling)
- Shared infrastructure

**Setup Steps:**

1. **Create Uberspace Account**
   ```bash
   # Sign up at https://uberspace.de
   # SSH access provided immediately
   ```

2. **Install Dependencies**
   ```bash
   # PostgreSQL
   uberspace tools version use postgresql 15
   createdb ts2go
   
   # Redis
   uberspace tools version use redis 7
   redis-cli ping  # Verify
   ```

3. **Deploy Backend**
   ```bash
   # Build locally
   cd saas/backend
   GOOS=linux GOARCH=amd64 go build -o ts2go-api ./cmd/api
   
   # Upload via SCP
   scp ts2go-api <user>@<host>.uberspace.de:~/ts2go-api
   
   # Upload config
   scp .env.production <user>@<host>.uberspace.de:~/ts2go-api/.env
   ```

4. **Configure Supervisord**
   ```bash
   # Create ~/etc/services.d/ts2go-api.ini
   [program:ts2go-api]
   command=/home/<user>/ts2go-api
   directory=/home/<user>
   autostart=yes
   autorestart=yes
   startsecs=30
   environment=PORT="8080",DATABASE_URL="postgresql://...",REDIS_URL="redis://..."
   ```

5. **Start Service**
   ```bash
   supervisorctl reread
   supervisorctl update
   supervisorctl start ts2go-api
   ```

6. **Configure Web Backend**
   ```bash
   # Configure web backend to proxy to port 8080
   uberspace web backend set / --http --port 8080
   ```

**Environment Variables:**
```bash
PORT=8080
DATABASE_URL=postgresql://<user>:<password>@localhost:5432/ts2go
REDIS_URL=redis://localhost:6379
JWT_SECRET=<generate-random-secret>
FRONTEND_URL=https://<username>.github.io
ENVIRONMENT=production
```

**GitHub Actions Deployment:**
```yaml
name: Deploy Backend (Uberspace)

on:
  push:
    branches:
      - main
    paths:
      - 'saas/backend/**'
  workflow_dispatch:

jobs:
  deploy:
    name: Deploy to Uberspace
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22.5'
      
      - name: Build for Linux
        working-directory: saas/backend
        run: |
          GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ts2go-api ./cmd/api
      
      - name: Deploy via SSH
        uses: appleboy/scp-action@v0.1.7
        with:
          host: ${{ secrets.UBERSPACE_HOST }}
          username: ${{ secrets.UBERSPACE_USER }}
          key: ${{ secrets.UBERSPACE_SSH_KEY }}
          source: "saas/backend/ts2go-api"
          target: "~/ts2go-api"
      
      - name: Restart service
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.UBERSPACE_HOST }}
          username: ${{ secrets.UBERSPACE_USER }}
          key: ${{ secrets.UBERSPACE_SSH_KEY }}
          script: supervisorctl restart ts2go-api
```

### Option 2: Cloud Functions (Future Migration)

**Pros:**
- Auto-scaling
- Pay-per-request (cost-efficient at scale)
- No server management
- Global distribution (CDN)

**Cons:**
- Cold starts (latency)
- More complex setup
- Vendor lock-in risk
- Higher cost at low usage

**Platforms:**
- **Vercel**: Easy setup, integrated with GitHub, good for API routes
- **Cloudflare Workers**: Edge functions, lowest latency, KV storage for caching
- **Google Cloud Functions**: Full GCP integration, expensive for low usage
- **AWS Lambda**: Most flexible, steeper learning curve

**When to Migrate:**
- User base grows beyond single server capacity
- Need global distribution
- Auto-scaling becomes critical
- Budget allows for function costs (typically >$20/month at scale)

### Option 3: Docker + VPS (Alternative)

**Pros:**
- Full control
- Docker containerization
- Easy local development parity
- Can self-host cheaply (DigitalOcean, Hetzner, etc.)

**Cons:**
- More DevOps overhead
- Manual scaling
- Security responsibility

**Cost:**
- DigitalOcean Droplet: $6/month (basic)
- Hetzner VPS: €4.50/month (CX11)

## Recommendation

**Phase 1 (Now - Testing):**
Use **Uberspace** for initial deployment and testing:
- Fixed €5/month cost
- Easy PostgreSQL + Redis setup
- No cold starts
- Perfect for staging environment
- Can handle moderate production load

**Phase 2 (Production at Scale):**
Migrate to **Cloudflare Workers** or **Vercel** when:
- User base exceeds 1,000 active users
- Need global distribution
- Traffic patterns justify function costs
- Budget allows >$20/month for backend

## Configuration

### Required Secrets (GitHub Actions)

Add to repository secrets (Settings → Secrets → Actions):

```bash
UBERSPACE_HOST=<username>.uber.space
UBERSPACE_USER=<username>
UBERSPACE_SSH_KEY=<private-key-content>
```

### Database Migrations

Run migrations on deployment:

```bash
# SSH into Uberspace
ssh <user>@<host>.uberspace.de

# Run migrations
cd ~/ts2go-api
./ts2go-api migrate up
```

### Monitoring

**Logs:**
```bash
# View application logs
supervisorctl tail -f ts2go-api

# PostgreSQL logs
tail -f ~/logs/postgresql.log

# Redis logs
tail -f ~/logs/redis.log
```

**Health Check:**
```bash
# Add health endpoint in backend
GET /health
→ {"status": "ok", "version": "v1.0.0", "database": "connected", "redis": "connected"}
```

### Backup Strategy

**Database Backups:**
```bash
# Cron job for daily backups
0 3 * * * pg_dump ts2go | gzip > ~/backups/ts2go-$(date +\%Y\%m\%d).sql.gz

# Keep 7 days of backups
find ~/backups -name "ts2go-*.sql.gz" -mtime +7 -delete
```

**Restore:**
```bash
gunzip -c ~/backups/ts2go-20240115.sql.gz | psql ts2go
```

## Next Steps

1. ✅ Document deployment strategy
2. ⏳ Create Uberspace account
3. ⏳ Set up PostgreSQL and Redis
4. ⏳ Create deployment workflow (`.github/workflows/deploy-backend.yml`)
5. ⏳ Configure environment variables
6. ⏳ Set up monitoring and health checks
7. ⏳ Test deployment pipeline
8. ⏳ Document manual deployment procedure (fallback)
