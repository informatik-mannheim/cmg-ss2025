### How this works:

A middleware is a component that intercepts incoming HTTP requests to perform validation before they reach the actual
code.

1. The middleware reads the contents from the Authorization Header-Field:
   `Authorization: Bearer eyxyz1234qwer...` and parses it with the implemented JWT-Library
2. It extracts the `kid` (Key ID) from the `JWT-Header` to find the corresponding entry in the JWKS. If the exact `kid`
   exists, a public key gets build.
3. It now has access to the hash that was sent and locked by the pirvate key.
4. The Middleware creates a hash (the algorithm is provided in the `JWT-Header`)containing `Header` and `Payload` and
   then compares createdHash == received Hash
5. If both values match, the token is valid and has not been manipulated. Otherwise, the request is rejected

### How to implement:

1. In your main.go:

```go
jwksUrl := os.Getenv("JWKS_URL") // URL saved in .env
if jwksUrl == "" {log.Fatalf("JWKS_URL variable not set")
}
err := auth.InitJWKS(jwksUrl)
if err != nil {
log.Fatalf("Failed to initialize JWKS: %v", err)
}
   ```

2. Now we create a helper function that all handler-functions pass through.

```go
func secure(h http.HandlerFunc) http.Handler {
return auth.AuthMiddleware(h)
}
```
and then register your endpoints like this:
```go
mux.Handle("/endpoint", secure(handler.HandleEndpointRequest))
```
