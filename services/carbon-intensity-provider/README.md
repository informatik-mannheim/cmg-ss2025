
# Carbon Intensity Provider Service

This service provides an API to store and retrieve carbon intensity data, following the Ports & Adapters (Hexagonal) architecture.

> **WARNING**  
> This is a basic in-memory implementation for academic use. No persistence or external integrations are included.

---

## Architecture

- `model/`: Contains data structures like `CarbonIntensityData` and `Zone`
- `core/`: Implements logic with an in-memory map (`zoneMap`)
- `ports/`: Defines the `CarbonIntensityProvider` interface
- `api/`: Exposes HTTP routes
- `main.go`: Wires everything and starts the server

---

## API Overview

### `GET /carbon-intensity/zones`

Returns the list of available zones.

#### Example Command

```bash
curl http://localhost:8080/carbon-intensity/zones
```

#### Sample Response

```json
{
  "zones": [
    { "code": "DE", "name": "DE" },
    { "code": "FR", "name": "FR" }
  ]
}
```

---

### `GET /carbon-intensity/{zone}`

Returns the carbon intensity for a specific zone.

#### Example Command

```bash
curl http://localhost:8080/carbon-intensity/DE
```

#### Sample Response

```json
{
  "zone": "DE",
  "carbonIntensity": 140.5
}
```

If the zone is not found:

```text
Zone not found
```

Status code: `404 Not Found`

---
