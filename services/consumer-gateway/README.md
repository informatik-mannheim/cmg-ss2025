# Consumer Gateway Service

The consumer-gateway service is a forwarding service that routes all client requests to the appropriate service. It does not store, evaluate or authenticate any data itself. 

> The implementation is in an early stage. Many things are still missing. Use with care.

## Usage

```bash
curl -X PUT -d '{ "Id": "34", "IntProp" : 23, "StringProp": "test" }' localhost:8080/entity
curl localhost:8080/entity/34
```
