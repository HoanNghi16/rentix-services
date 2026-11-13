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

## Complex-created events

Set `KAFKA_BROKERS` to a comma-separated broker list (defaults to `localhost:9092`). Before starting property-service, create the topics from the repository root with `./kafka-topics/create-topics.ps1`. Creating a complex writes the complex, employee assignments, and a `complex.created` outbox event in one PostgreSQL transaction. A background publisher sends pending events to Kafka and retries failures with capped exponential backoff. Delivery is at least once, so consumers should deduplicate events if they require exactly-once effects.
