# User Management Service

A minimal user management microservice that supports secure user registration, login, and authentication via Auth0. Includes support for role-based users (`consumer`, `provider`, `job scheduler`) and persistence via JSON storage.

---

## ✨ Features

* Manual registration of users by admin only (via shared admin secret)
* Login using generated secret
* JWT issued via Auth0 (Machine-to-Machine)
* Secure Argon2id hashing for secrets
* Role enforcement: only one `job scheduler` allowed
* Notifier integration for logging registrations/logins
* Persistent user storage via `users.json`

---

## 📦 Project Structure

```
services/user-management/
├── main.go                  # Entry point
├── model/model.go           # User + role types
├── core/service.go          # Business logic + persistence
├── notifier/notifier.go     # Event logger
├── adapters/handler-http/   # HTTP handlers (auth.go)
├── users.json               # Stored users (auto-created)
├── .env                     # Auth0 and admin config
```

---

## 🔧 .env Configuration

Create a `.env` file with the following:

```env
AUTH0_CLIENT_ID=your_auth0_client_id
AUTH0_CLIENT_SECRET=your_auth0_client_secret
AUTH0_TOKEN_URL=https://your-domain.eu.auth0.com/oauth/token
JWT_AUDIENCE=https://user-management.local
JWT_ISSUER=https://your-domain.eu.auth0.com/
JWKS_URL=https://your-domain.eu.auth0.com/.well-known/jwks.json
YOUR_SECRET=your_super_secret_admin_token
```

---

## 🚀 Run the Service

```bash
cd services/user-management
go run main.go
```

---

## 🔐 Register a User (Admin Only)

```bash
curl -X POST http://localhost:8081/auth/register \
  -H "Content-Type: application/json" \
  -H "X-Admin-Secret: your_super_secret_admin_token" \
  -d '{ "role": "consumer" }'
```

Returns:

```json
{
  "id": "uuid",
  "secret": "generated-secret"
}
```

---

## 🔑 Login

```bash
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{ "secret": "generated-secret" }'
```

Returns:

```json
{
  "token": "<auth0-jwt>"
}
```

---

## 📁 Persistent Storage

* Users are stored in `users.json` after each registration.
* Secrets are hashed with Argon2id.

---

## 🔍 TODO / Extensions

* Add `/auth/me` endpoint (token introspection)
* Database storage (PostgreSQL, SQLite)
* Admin panel (CLI or Web)
* Docker support

---