# jwt-otp-auth (Go + GORM + Redis + MySQL)

Authentication service with **OTP (Redis)** + **JWT** and user management.  
Implements **Clean/Hexagonal Architecture**, fully containerized with **Docker** (MySQL, Redis, phpMyAdmin).

## Features

- 🔑 OTP via Redis (hash + TTL + attempts + blocklist)
- 🔒 JWT HS256 with standard claims (`sub`, `phone`, `iat`, `exp`)
- 🗄 MySQL + GORM (AutoMigrate)
- ⚡ Rate limit & temporary blocking
- 📖 Swagger UI & Prometheus `/metrics`
- ⚙️ Config via ENV (Viper)
- 📝 Structured logging (Zerolog)

---

## Quick Start

### 1) Requirements

- Docker & Docker Compose
- (Optional) GNU `make`

### 2) Setup

```bash
cp infra/.env.example infra/.env
cp src/.env.example src/.env
```

Edit **infra/.env** → MySQL root password & ports.  
Edit **src/.env** → DB/Redis/JWT configs (must match infra).

### 3) Run

```bash
make up
# or
docker compose --env-file infra/.env -f infra/docker-compose.dev.yml up -d --build
```

### 4) Database Migration

```bash
make migrate
```

### 5) Health Checks

- API → http://localhost:8080/healthz
- DB → /health/db
- Redis → /health/redis
- Metrics → /metrics
- phpMyAdmin → http://localhost:8081 (Server: mysql)

### 6) Swagger

http://localhost:8080/swagger/index.html

---

## Folder Structure

```
infra/   → docker-compose, Dockerfiles, env
src/
  cmd/api/main.go
  cmd/migrate/main.go
  internal/
    adapters/ (http, db, cache, jwt)
    core/ (domain, ports, services)
    pkg/ (config, db, logger, util)
  docs/ (swagger)
  migrations/
Makefile
```

---

## Example API Calls

**1) Request OTP**

```bash
curl -X POST http://localhost:8080/v1/auth/otp/request   -H "Content-Type: application/json"   -d '{"phone":"09120000000"}'
# → 204 No Content (OTP printed in logs on dev)
```

**2) Verify OTP & Get Token**

```bash
curl -X POST http://localhost:8080/v1/auth/otp/verify   -H "Content-Type: application/json"   -d '{"phone":"09120000000","otp":"<OTP>"}'
```

**3) Current User**

```bash
curl http://localhost:8080/v1/users/me   -H "Authorization: Bearer <TOKEN>"
```

**4) List Users (search + pagination)**

```bash
curl "http://localhost:8080/v1/users?search=0912&page=1&per_page=20"   -H "Authorization: Bearer <TOKEN>"
```

---

## Makefile Commands

- `make up` → Start services
- `make down` → Stop services
- `make logs` → Tail logs
- `make ps` → Show containers
- `make migrate` → Run DB migrations
- `make swagger` → Regenerate Swagger docs
- `make restart` → Restart stack

---

## Notes

- Use a **strong, private `JWT_SECRET`** in production.
- Enable **HTTPS & HSTS** in production environments.
