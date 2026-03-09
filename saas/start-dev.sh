#!/bin/bash
# Quick start script for local development

set -e

echo "🚀 Starting ts2go SaaS local development environment..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Copy .env.example if .env doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
fi

# Start infrastructure with Docker Compose
echo "🐳 Starting PostgreSQL, Redis, and MinIO..."
docker-compose up -d postgres redis minio

# Wait for services to be healthy
echo "⏳ Waiting for services to be ready..."
sleep 5

# Check if services are healthy
echo "✅ Checking service health..."
docker-compose ps

# Build the API server
echo "🔨 Building API server..."
cd backend/api && go build -o ../../bin/ts2go-api .
cd ../..

echo ""
echo "✅ Local development environment ready!"
echo ""
echo "📊 Services running:"
echo "   - PostgreSQL: localhost:5432"
echo "   - Redis: localhost:6379"
echo "   - MinIO API: localhost:9000"
echo "   - MinIO Console: http://localhost:9001"
echo ""
echo "🔑 MinIO credentials:"
echo "   - Username: ts2go"
echo "   - Password: ts2go_minio_password"
echo ""
echo "▶️  To start the API server:"
echo "   ./bin/ts2go-api"
echo ""
echo "🛑 To stop services:"
echo "   docker-compose down"
echo ""
