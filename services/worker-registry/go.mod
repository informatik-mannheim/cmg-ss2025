module github.com/informatik-mannheim/cmg-ss2025/services/worker-registry

go 1.24.1

toolchain go1.24.2

require github.com/gorilla/mux v1.8.1

require github.com/google/uuid v1.6.0

require (
	github.com/informatik-mannheim/cmg-ss2025/pkg/logging v0.0.0-20250605183446-f825cbc16886
	github.com/lib/pq v1.10.9
)
