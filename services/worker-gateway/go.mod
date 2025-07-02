module github.com/informatik-mannheim/cmg-ss2025/services/worker-gateway

replace github.com/informatik-mannheim/cmg-ss2025/pkg/logging => ../../pkg/logging

go 1.24.3

toolchain go1.24.4

require github.com/informatik-mannheim/cmg-ss2025/pkg/logging v0.0.0-20250622162429-594d99bc7456

require github.com/informatik-mannheim/cmg-ss2025/pkg/auth v0.0.0-20250628182332-e35f88de24b9

require (
	github.com/MicahParks/keyfunc v1.9.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
)
