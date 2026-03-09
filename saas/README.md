# SaaS Directory Structure

This directory contains all components for the ts2go Software-as-a-Service platform.

## 📁 Directory Layout

```
saas/
├── backend/           # Go backend services
│   ├── api/           # RESTful API server (Gin/Echo)
│   ├── auth/          # Authentication service (JWT, OAuth2)
│   ├── db/            # Database models & migrations (PostgreSQL)
│   ├── storage/       # File storage (S3/MinIO)
│   ├── queue/         # Job queue (Redis + BullMQ)
│   ├── worker/        # Transpilation workers
│   └── middleware/    # API middleware (auth, rate limit, logging)
│
├── frontend/          # Next.js web application
│   ├── pages/         # Next.js pages
│   ├── components/    # React components
│   ├── lib/           # Utilities
│   └── public/        # Static assets
│
├── k8s/               # Kubernetes manifests
│   ├── api/           # API deployment
│   ├── worker/        # Worker deployment
│   ├── frontend/      # Frontend deployment
│   └── ingress/       # Ingress rules
│
├── helm/              # Helm charts
│   └── ts2go-saas/    # Main chart
│
└── docs/              # SaaS-specific documentation
    ├── api/           # API documentation
    ├── architecture/  # Architecture diagrams
    └── guides/        # Development guides
```

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- Node.js 20+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose

### Local Development

```bash
# Start infrastructure
docker-compose up -d postgres redis

# Start backend API
cd backend/api
go run main.go

# Start frontend
cd frontend
npm install
npm run dev

# Start worker
cd backend/worker
go run main.go
```

## 📚 Documentation

- **[SaaS Transformation Plan](../docs/SAAS_TRANSFORMATION_PLAN.md)** - Complete plan
- **[Implementation Tasks](../docs/SAAS_IMPLEMENTATION_TASKS.md)** - Detailed task list
- **[API Reference](docs/api/)** - API documentation (coming soon)
- **[Architecture](docs/architecture/)** - System design (coming soon)

## 🎯 Current Status

**Phase:** Planning & Setup  
**Progress:** 0/20 weeks  
**Next Steps:**
1. Set up PostgreSQL schema
2. Implement authentication service
3. Create basic API endpoints
4. Set up job queue

## 🔗 Quick Links

- API Server: http://localhost:8080
- Frontend: http://localhost:3000
- GraphQL Playground: http://localhost:8080/graphql
- Swagger Docs: http://localhost:8080/swagger

---

**For detailed information, see the main SaaS planning documents in the `docs/` directory.**
