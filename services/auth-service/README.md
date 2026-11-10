# Auth service

The auth service exposes its internal API over gRPC with mutual TLS. It has no HTTP listener.

## Required environment

- `AUTH_GRPC_PORT`: TCP port for the gRPC listener.
- `AUTH_TLS_CERT_FILE`: PEM server certificate.
- `AUTH_TLS_KEY_FILE`: PEM private key matching the server certificate.
- `AUTH_CLIENT_CA_FILE`: PEM CA certificate bundle trusted to issue client certificates.
- `AUTH_ALLOWED_CLIENTS`: comma-separated exact DNS SAN identities allowed to call auth, for example `api-gateway.rentix.internal`.
- `DB_URL`, `ACCESS_SECRET`, and `REFRESH_SECRET`: database connection and token signing secrets.

The server refuses to start when any TLS setting or the caller allowlist is missing. Client certificates must chain to the configured CA and contain at least one DNS SAN listed in `AUTH_ALLOWED_CLIENTS`. Common names are not used as identities.

## Local development certificates

Create a local-only CA, auth server certificate, and gateway client certificate with:

```powershell
go run ./cmd/devcert
```

The generated `certs/` directory is ignored by Git. Its private keys are for local development only and must not be used in deployed environments. The gateway client identity is `api-gateway.rentix.internal`.

Run from this directory with `go run ./cmd/grpc` or build and run `grpc.exe`.
