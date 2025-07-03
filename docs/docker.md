# 🐳 Docker Setup Guide

This guide covers how to run the Go API with Docker and MySQL.

## 📋 Quick Start

### Production Setup

```bash
make docker-run
```

### Development Setup

```bash
make docker-dev
```

## 🏗️ Architecture

- **Go API**: Optimized container with MySQL connection
- **MySQL 8.0**: Persistent database with health checks
- **Network**: Isolated bridge network

## 🔧 Services

### docker-compose.yml (Production)

- **api**: Go application (port 8080)
- **mysql**: MySQL database (port 3306)

### docker-compose.dev.yml (Development)

- **api**: Go with hot reload
- **mysql**: Same as production
- **adminer**: Database admin UI (port 8081)

## 🚀 Usage

```bash
# Start production
make docker-run

# Start development
make docker-dev

# View logs
make docker-logs

# Stop services
make docker-stop

# Cleanup
make docker-clean
```

## 🔍 Troubleshooting

**Connection Issues:**

```bash
docker-compose ps        # Check status
docker-compose logs mysql # Check MySQL logs
```

**Reset Database:**

```bash
make docker-clean
make docker-run
```
