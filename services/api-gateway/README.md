# API gateway

The gateway uses Gin to accept client HTTP/JSON requests on port `8080` and calls auth-service and property-service over gRPC with mutual TLS.

## Local run

Start the databases from the repository root:

```powershell
docker compose up -d rt-auth-db rt-propety-db
```

Generate the shared development CA and gateway client certificate once from `services/auth-service` with `go run ./cmd/devcert`. Generate the property server certificate from `services/property-service` with `go run ./cmd/devcert`, then start both gRPC services from their own directories with `go run ./cmd/grpc`.

Set `ACCESS_SECRET` in the gateway process to the same value used by auth-service. In PowerShell, set it from your local auth `.env` value, then start the gateway from `services/api-gateway`:

```powershell
$env:ACCESS_SECRET = "<same value as services/auth-service/.env>"
go run ./cmd/gateway
```

The gateway expects the generated certificates at `../auth-service/certs/` when started from this directory. Set `GATEWAY_AUTH_CA_FILE`, `GATEWAY_MTLS_CERT_FILE`, and `GATEWAY_MTLS_KEY_FILE` to override those paths. Auth defaults to `localhost:8081`; property defaults to `localhost:8082`. `GATEWAY_HTTP_ADDR` defaults to `:8080`.

## Postman requests

- `GET http://localhost:8080/healthz`
- `POST http://localhost:8080/auth/register`
- `POST http://localhost:8080/auth/login`
- `GET http://localhost:8080/property/complex` (Bearer access token required)
- `POST http://localhost:8080/property/complex` (Bearer access token required)

For auth POST requests, select **Body > raw > JSON** and send:

```json
{
  "email": "person@example.com",
  "password": "your-password"
}
```

For property requests, set Postman's Authorization type to **Bearer Token** and use the access token returned by `/auth/login`. A create-complex body can contain `name`, `avatar`, `address`, `employee_user_ids`, and `services`; user ID and role are taken from the JWT claims, not the body.
