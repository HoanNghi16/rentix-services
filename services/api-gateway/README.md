# API gateway

The gateway uses Gin to accept client HTTP/JSON requests on port `8080` and calls auth-service, property-service, and employee-service over gRPC with mutual TLS.

## Local run

Start the databases from the repository root:

```powershell
docker compose up -d rt-auth-db rt-propety-db rt-employee-db
```

Generate the shared development CA and gateway client certificate once from `services/auth-service` with `go run ./cmd/devcert`. Generate the property and employee server certificates from their respective service directories with `go run ./cmd/devcert`.

Start each gRPC service from its own directory with `go run ./cmd/grpc`. For employee-service, set `DB_URL` to the employee database connection string and configure its mTLS files before starting:

```powershell
$env:DB_URL = "postgres://rentix:<EMPLOYEE_DB_PASSWORD>@localhost:5445/rentix_employee?sslmode=disable"
$env:EMPLOYEE_TLS_CERT_FILE = "certs/employee-server.crt"
$env:EMPLOYEE_TLS_KEY_FILE = "certs/employee-server.key"
$env:EMPLOYEE_CLIENT_CA_FILE = "../auth-service/certs/ca.crt"
$env:EMPLOYEE_ALLOWED_CLIENTS = "api-gateway.rentix.internal"
go run ./cmd/grpc
```

Set `ACCESS_SECRET` in the gateway process to the same value used by auth-service. In PowerShell, set it from your local auth `.env` value, then start the gateway from `services/api-gateway`:

```powershell
$env:ACCESS_SECRET = "<same value as services/auth-service/.env>"
go run ./cmd/gateway
```

The gateway expects the generated certificates at `../auth-service/certs/` when started from this directory. Set `GATEWAY_AUTH_CA_FILE`, `GATEWAY_PROPERTY_CA_FILE`, `GATEWAY_EMPLOYEE_CA_FILE`, `GATEWAY_MTLS_CERT_FILE`, and `GATEWAY_MTLS_KEY_FILE` to override those paths. Auth defaults to `localhost:8081`; property to `localhost:8082`; employee to `localhost:8083`. `GATEWAY_HTTP_ADDR` defaults to `:8080`.

## Postman requests

- `GET http://localhost:8080/healthz`
- `POST http://localhost:8080/auth/register`
- `POST http://localhost:8080/auth/login`
- `GET http://localhost:8080/property/complex` (Bearer access token required)
- `POST http://localhost:8080/property/complex` (Bearer access token required)
- `GET http://localhost:8080/employee` (Bearer access token required; returns all employees belonging to the owner)

For auth POST requests, select **Body > raw > JSON** and send:

```json
{
  "email": "person@example.com",
  "password": "your-password"
}
```

For property requests, set Postman's Authorization type to **Bearer Token** and use the access token returned by `/auth/login`. A create-complex body can contain `name`, `avatar`, `address`, `employee_user_ids`, and `services`; user ID and role are taken from the JWT claims, not the body.
