# Carbon Intensity Provider Service

This service provides an API to store and retrieve carbon intensity data, following the Ports & Adapters (Hexagonal) architecture.

> **WARNING**
> The implementation is in an early stage. Some functionality may be missing or subject to change.

## Usage

```bash
curl -X PUT -d '{ "Id": "34", "IntProp" : 23, "StringProp": "test" }' localhost:8080/carbon-intensity-provider
curl localhost:8080/carbon-intensity-provider/34
```
