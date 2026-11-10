# API gateway

The gateway accepts client HTTP/JSON requests on port `8080` and calls auth-service over gRPC with mutual TLS.

## Local run

Start the auth database from the repository root:

```powershell
docker compose up -d rt-auth-db
```

Generate local certificates once from `services/auth-service` with `go run ./cmd/devcert`, then start auth-service from that directory with `go run ./cmd/grpc` (or run its built executable). Start the gateway from `services/api-gateway`:

```powershell
go run ./cmd/gateway
```

The gateway expects the generated certificates at `../auth-service/certs/` when started from this directory. Set `GATEWAY_AUTH_CA_FILE`, `GATEWAY_MTLS_CERT_FILE`, and `GATEWAY_MTLS_KEY_FILE` to override those paths. `AUTH_GRPC_TARGET` defaults to `localhost:8081`; `AUTH_GRPC_SERVER_NAME` defaults to `localhost`; `GATEWAY_HTTP_ADDR` defaults to `:8080`.

## Postman requests

- `GET http://localhost:8080/healthz`
- `POST http://localhost:8080/auth/register`
- `POST http://localhost:8080/auth/login`

For both POST requests, select **Body → raw → JSON** and send:

```json
{
  "email": "person@example.com",
  "password": "your-password"
}
```
