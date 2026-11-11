param(
    [ValidateRange(1, 1000)]
    [int]$Partitions = 3
)

$ErrorActionPreference = "Stop"
$repoRoot = Split-Path -Parent $PSScriptRoot
$registryPath = Join-Path $PSScriptRoot "topics.json"
$topicRegistry = Get-Content -LiteralPath $registryPath -Raw | ConvertFrom-Json

Push-Location $repoRoot
try {
    foreach ($entry in $topicRegistry.PSObject.Properties) {
        $topicName = [string]$entry.Value
        if ([string]::IsNullOrWhiteSpace($topicName)) {
            throw "Topic registry entry '$($entry.Name)' has an empty topic name."
        }

        Write-Host "Registering Kafka topic '$topicName'..."
        $dockerArgs = @(
            "compose", "exec", "-T", "rt-kafka",
            "/opt/kafka/bin/kafka-topics.sh",
            "--bootstrap-server", "localhost:9092",
            "--create", "--if-not-exists",
            "--topic", $topicName,
            "--partitions", $Partitions.ToString(),
            "--replication-factor", "1"
        )
        & docker @dockerArgs
        if ($LASTEXITCODE -ne 0) {
            throw "Failed to create Kafka topic '$topicName'."
        }
    }
}
finally {
    Pop-Location
}
