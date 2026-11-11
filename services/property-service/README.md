# Property service

Property-service exposes its complex operations over gRPC with mutual TLS on `PROPERTY_GRPC_PORT` (local default: `8082`). The API gateway is the allowed client. It passes the authenticated user ID and role from the verified access token.

## Local run

Generate the shared development CA and gateway client certificate from `services/auth-service` if they do not already exist:

```powershell
go run ./cmd/devcert
```

From this directory, issue the property server certificate from that development CA and start the service:

```powershell
go run ./cmd/devcert
go run ./cmd/grpc
```

The local `.env` expects the server certificate under `certs/` and trusts the development CA under `../auth-service/certs/`. The property gRPC server requires an authorized client certificate whose DNS SAN matches `PROPERTY_ALLOWED_CLIENTS`.
