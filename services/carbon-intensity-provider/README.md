# 🌍 Carbon Intensity Provider

A microservice that fetches and provides real-time carbon intensity data for various electricity zones using the [Electricity Maps API](https://www.electricitymap.org/). Zones are filtered based on available API tokens.

---

## ✨ Features

- Live or offline data mode
- Fetches zone metadata (unauthenticated)
- Retrieves carbon intensity data (authenticated by zone)
- REST API: `/carbon-intensity/zones` and `/carbon-intensity/{zone}`
- Logs and stores data using file-based persistence (`zones.json`, `zones_metadata.json`)
- Uses Go interfaces and clean architecture with adapters (handlers, providers, repo, notifier)

---

## 📂 Project Structure

```
services/carbon-intensity-provider/
├── main.go                        # Entry point
├── core/                          # Business logic
├── ports/                         # Interfaces (Repo, Notifier, Provider)
├── adapters/
│   ├── handler-http/             # HTTP handlers (API endpoints)
│   ├── notifier/                 # Logging notifier
│   ├── repo-in-memory/          # File-based repository
│   └── provider/
│       └── electricity-maps/    # API fetcher
├── zones.json                   # Stored carbon intensity data
├── zones_metadata.json          # Stored zone metadata
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## 🔧 .env Configuration

Set these environment variables:

```env
USE_LIVE=true
TOKEN_GB=your_token_here
TOKEN_FR=your_token_here
# Add other zones and tokens as needed
```

Tokens are zone-specific and loaded dynamically via `TOKEN_<ZONE>`.

---

## 🚀 Running Locally

### With Go:

```bash
USE_LIVE=true TOKEN_GB=your_token_here TOKEN_DE=your_token_here go run main.go
```

---

## 🐳 Docker Usage

### 🚧 Build and Run with Docker Compose

To run the service inside Docker with fresh builds and logs:

```bash
docker compose up --build
```

To run it in the background:

```bash
docker compose up --build -d
```

### 🧹 Stop and Remove Containers

When you're done and want to stop and clean up:

```bash
docker compose down
```

---

## 🌐 API Endpoints

- `GET /carbon-intensity/zones`: Returns list of available zones (filtered by tokens)
- `GET /carbon-intensity/{zone}`: Returns current carbon intensity data for a specific zone

---

## 📁 Data Storage

- `zones.json` stores the fetched carbon data.
- `zones_metadata.json` stores zone names (for display).
- Both are automatically updated during runtime.

---

## 🧪 Example CURL

```bash
curl http://localhost:8080/carbon-intensity/zones
curl http://localhost:8080/carbon-intensity/GB
```

---

## 📌 Notes

- The `/zones` endpoint is built from the live metadata fetched from Electricity Maps, filtered by zones you’ve provided tokens for.
- The app works without authentication for metadata but requires tokens for real-time carbon data.
