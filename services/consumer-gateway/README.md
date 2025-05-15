# Consumer Gateway Service

The consumer-gateway service is a forwarding service that routes all client requests to the appropriate service. It does not store, evaluate or authenticate any data itself. 

> The implementation is in an early stage. Many things are still missing. Use with care.

## Usage

```bash
curl -X POST -d '{ "image_id": "1234", "zone" : "GER", "params": "-a" }' localhost:8080/jobs

```
