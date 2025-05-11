# User Management Service

A minimal user management microservice supporting secure login and token issuance via Auth0. Designed for machine-to-machine authentication with configurable roles: `consumer`, `provider`, `job scheduler`.

---

## ✨ Features

- Role-based Auth0 integration (via M2M)
- Login via `client_id.client_secret` → JWT
- Admin-only `/auth/register` using shared secret
- Notifier support for logging events
- Mock vs. Live Auth0 switching via `USE_LIVE`
- Fully tested (≥94% coverage)
- Docker & Makefile integration

---

## 📦 Project Structure

```
services/user-management/
├── main.go                 # Entry point
├── adapters/
│   ├── handler-http/       # HTTP routes for /auth/login and /auth/register
│   ├── auth/               # Auth0 adapter (mock + live)
│   └── notifier/           # Stdout notifier
├── ports/                  # Interfaces: Notifier, AuthProvider, Role
├── Dockerfile              # Container build config
├── docker-compose.yaml     # Runtime environment
├── Makefile                # Task automation
└── README.md               # This file
```

---

## 🔧 Configuration (via ENV)

Required environment variables:

```env
AUTH0_TOKEN_URL=https://dev-jqhwcu7xuwgdqi56.eu.auth0.com/oauth/token
JWT_AUDIENCE=https://green-load-shifting-platform/
ADMIN_SECRET_HASH=sha256-hash-of-secret
USE_LIVE=true
```

You can define these in your Docker Compose or terminal.

---

## 🛠 Run Locally

```bash
make run       # go run main.go
make test      # run all unit tests
make build     # build binary
```

---

## 🐳 Run with Docker Compose

```bash
make docker-up
make docker-down
```

Ensure the following are defined via environment or .env:
- `AUTH0_TOKEN_URL`
- `JWT_AUDIENCE`
- `ADMIN_SECRET_HASH`
- `USE_LIVE`

---

## 🧪 Testing

```bash
make test
```

Includes full coverage of:
- HTTP handlers
- Auth adapter
- Notifier logging

---

## 🔐 Register

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -H "X-Admin-Secret: your_admin_secret" \
  -d '{ "role": "consumer" }'
```

---

## 🔑 Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{ "secret": "client_id.client_secret" }'
```

---

## 📈 Makefile Targets

```bash
make run          # Start the service
make test         # Run tests with coverage
make fmt          # Format code
make docker-up    # Start service with Docker Compose
make docker-down  # Tear down
```

---

## 📌 Notes

- No persistent storage: login data comes from Auth0 only
- Job scheduler is a singleton role
- Notifier can be replaced with alternate implementations

---
