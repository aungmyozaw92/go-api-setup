# 🐳 Docker Setup Guide

This guide covers how to run the Go API with Docker and MySQL using Docker Compose.

## 📋 Quick Start

### **1. Production Setup**

```bash
# Build and run in production mode
make docker-run

# Or manually:
docker-compose up -d
```

### **2. Development Setup**

```bash
# Run in development mode with hot reload
make docker-dev

# Or manually:
docker-compose -f docker-compose.dev.yml up
```

## 🏗️ Architecture

### **Production (docker-compose.yml)**

- **Go API**: Optimized binary in scratch container
- **MySQL 8.0**: Persistent data with named volumes
- **Network**: Isolated bridge network
- **Health Checks**: MySQL health monitoring

### **Development (docker-compose.dev.yml)**

- **Go API**: Live code reload with volume mounts
- **MySQL 8.0**: Same as production
- **Adminer**: Database admin interface (port 8081)
- **Hot Reload**: Code changes trigger automatic restart

## 🔧 Configuration

### **Environment Variables**

**Production:**

```env
DB_HOST=mysql                    # Docker service name
DB_PORT=3306
DB_USER=api_user
DB_PASSWORD=api_password
DB_NAME=go_api_setup
SERVER_PORT=8080
JWT_SECRET=your-production-secret
APP_ENV=production
```

**Development:**

```env
DB_HOST=mysql                    # Docker service name
DB_PORT=3306
DB_USER=api_user
DB_PASSWORD=api_password
DB_NAME=go_api_setup
SERVER_PORT=8080
JWT_SECRET=dev-secret-key-not-for-production
APP_ENV=development
```

### **Ports**

- **8080**: Go API Server
- **3306**: MySQL Database
- **8081**: Adminer (dev only)

## 🚀 Usage Examples

### **Start Services**

```bash
# Production (detached)
make docker-run

# Development (with logs)
make docker-dev

# Check status
docker-compose ps
```

### **View Logs**

```bash
# All services
make docker-logs

# API only
make docker-logs-api

# MySQL only
docker-compose logs -f mysql
```

### **Database Management**

```bash
# Access Adminer (development)
open http://localhost:8081
# Server: mysql
# Username: api_user
# Password: api_password
# Database: go_api_setup

# Direct MySQL access
docker-compose exec mysql mysql -u api_user -p go_api_setup
```

### **API Testing**

```bash
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password123"}'
```

## 🛠️ Development Workflow

### **1. Code Changes (Hot Reload)**

```bash
# Start development environment
make docker-dev

# Edit code - server automatically restarts
# View logs to see restart
make docker-logs-api
```

### **2. Database Changes**

```bash
# Access database shell
docker-compose exec mysql mysql -u api_user -p go_api_setup

# Or use Adminer
open http://localhost:8081
```

### **3. Debugging**

```bash
# Enter API container
make docker-shell

# Check Go version
docker-compose exec api go version

# View environment variables
docker-compose exec api env | grep DB_
```

## 📊 Container Details

### **Go API Container**

```yaml
# Production build stages:
# 1. golang:1.21-alpine (builder)
# 2. scratch (runtime)

# Key features:
- Multi-stage build for minimal size
- Non-root user for security
- Static binary compilation
- CA certificates included
```

### **MySQL Container**

```yaml
# Features:
- Persistent data volume
- Health checks
- Initialization script
- Custom user/database creation
```

## 🔍 Troubleshooting

### **Common Issues**

**1. API can't connect to MySQL**

```bash
# Check MySQL is healthy
docker-compose ps

# Check network connectivity
docker-compose exec api ping mysql

# Verify environment variables
docker-compose exec api env | grep DB_
```

**2. Permission denied errors**

```bash
# Check file permissions
ls -la Dockerfile docker-compose.yml

# Rebuild with no cache
docker-compose build --no-cache
```

**3. Port conflicts**

```bash
# Check what's using the port
lsof -i :8080
netstat -tulpn | grep 8080

# Change port in docker-compose.yml
ports:
  - "8081:8080"  # host:container
```

**4. Database initialization issues**

```bash
# Remove volumes and restart
make docker-clean
make docker-run

# Check MySQL logs
docker-compose logs mysql
```

### **Performance Issues**

**1. Slow startup**

```bash
# MySQL needs time to initialize
# Wait for health check to pass
docker-compose ps

# Watch logs for "ready for connections"
docker-compose logs -f mysql
```

**2. Build times**

```bash
# Use development mode for faster iteration
make docker-dev

# Build with cache
docker-compose build

# Build without cache (if needed)
docker-compose build --no-cache
```

## 🧹 Cleanup

### **Stop Services**

```bash
# Stop containers
make docker-stop

# Remove containers and volumes
make docker-clean

# Full cleanup (including images)
docker system prune -a
```

### **Reset Database**

```bash
# Remove only database volume
docker volume rm go-api-setup_mysql_data

# Restart to recreate
make docker-run
```

## 🚀 Production Deployment

### **Environment Setup**

```bash
# Create production environment file
cp env.template .env.production

# Edit with production values
vim .env.production

# Load environment
export $(cat .env.production | xargs)
```

### **Security Considerations**

```yaml
# Change default passwords
MYSQL_ROOT_PASSWORD: use-strong-password
MYSQL_PASSWORD: use-strong-password
JWT_SECRET: use-long-random-string
# Use secrets management
# - Docker secrets
# - Kubernetes secrets
# - HashiCorp Vault
```

### **Monitoring**

```bash
# Container health
docker-compose ps

# Resource usage
docker stats

# Logs
docker-compose logs --tail=100
```

## 📚 Additional Resources

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [MySQL Docker Hub](https://hub.docker.com/_/mysql)
- [Go Docker Best Practices](https://docs.docker.com/language/golang/)

## 🤝 Tips

1. **Use `.dockerignore`** to reduce build context size
2. **Multi-stage builds** keep production images small
3. **Health checks** ensure services are ready
4. **Named volumes** persist data between restarts
5. **Development volumes** enable hot reload
6. **Network isolation** improves security

---

**🎯 Ready to containerize your Go API!**
