# Kafka topics

`topics.json` is the registry of event topic names. The property names are the event identifiers used in code, and the values are the Kafka topic names.

Start the single-node development broker from the repository root:

```powershell
docker compose up -d rt-kafka
```

After the broker is healthy, create any missing registry topics:

```powershell
.\kafka-topics\create-topics.ps1
```

The script is safe to rerun. By default it creates topics with three partitions and replication factor one; use `-Partitions` to choose another partition count. Services on the `rentix-network` Docker network should connect to `rt-kafka:29092`. Host applications should connect to `localhost:9092`.

This Compose broker is a single-node local development setup. It uses plaintext Kafka listeners and replication factor one, so it is not a production cluster configuration.
