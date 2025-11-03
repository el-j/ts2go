# Build stage
FROM node:20-alpine AS builder-node

# Install TypeScript parser dependencies
WORKDIR /app/parser
COPY internal/transpiler/parser/package*.json ./
RUN npm ci --only=production

# Build Go binary
FROM golang:1.22.5-alpine AS builder-go

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy go workspace and modules
COPY go.work go.work
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY internal/ internal/
COPY runtime/ runtime/
COPY Makefile Makefile

# Copy node_modules from builder-node
COPY --from=builder-node /app/parser/node_modules /app/internal/transpiler/parser/node_modules

# Build the application
RUN go work sync && make build

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates nodejs npm

WORKDIR /app

# Copy the binary
COPY --from=builder-go /app/ts2go /usr/local/bin/ts2go

# Copy parser with dependencies
COPY --from=builder-node /app/parser/node_modules /app/internal/transpiler/parser/node_modules
COPY internal/transpiler/parser/*.js /app/internal/transpiler/parser/
COPY internal/transpiler/parser/package.json /app/internal/transpiler/parser/

# Copy runtime
COPY runtime/ /app/runtime/

# Copy docs
COPY README.md LICENSE* /app/

# Create directories for input/output
RUN mkdir -p /input /output

WORKDIR /workspace

ENTRYPOINT ["ts2go"]
CMD ["--help"]
