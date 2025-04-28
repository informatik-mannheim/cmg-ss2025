# carbon-intensity-provider Service

The carbon-intensity-provider service is an example service that demonstrates the folder structure of a microservice following the ports & adapters architecture.

> **WARNING**
> The implementation is in an early stage. Many things are still missing. Use with care.

## Usage

```bash
curl -X PUT -d '{ "Id": "34", "IntProp" : 23, "StringProp": "test" }' localhost:8080/carbon-intensity-provider
curl localhost:8080/carbon-intensity-provider/34
```
