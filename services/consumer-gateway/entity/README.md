# Consumer Gateway Service

The consumer-gateway service is a service that forwards all client requests to the appropriate service. Endpoints like `/jobs` or `/me` require a valid JWT in the `Authorization` header.

> The implementation is in an early stage. Many things are still missing. Use with care.

## Usage

```bash
curl -X PUT -d '{ "Id": "34", "IntProp" : 23, "StringProp": "test" }' localhost:8080/entity
curl localhost:8080/entity/34
```
