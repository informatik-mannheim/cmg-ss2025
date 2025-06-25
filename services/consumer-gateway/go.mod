module github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway

replace github.com/informatik-mannheim/cmg-ss2025/pkg/auth => ../../pkg/auth


go 1.24.3

require (
	github.com/gorilla/mux v1.8.1
	github.com/informatik-mannheim/cmg-ss2025/pkg/auth v0.0.0-00010101000000-000000000000
)

require (
	github.com/MicahParks/keyfunc v1.9.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
)
