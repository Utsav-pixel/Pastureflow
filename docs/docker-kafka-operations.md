# Docker and Kafka Operations Guide

This document provides instructions for managing Docker services and Kafka topics in the Pastureflow project.

## Docker Operations

### Starting Services
To start all services using Docker Compose:
```bash
docker-compose up -d
```

### Stopping Services
To stop all services:
```bash
docker-compose down
```

### Viewing Logs
To view logs for all services:
```bash
docker-compose logs -f
```

To view logs for a specific service:
```bash
docker-compose logs -f kafka
```

## Kafka Topic Management

### Creating Topics

To create a new Kafka topic, use the following command:

```bash
docker exec -it kafka-kafka-1 kafka-topics --bootstrap-server localhost:9092 --create --topic pasture.telemetry.v1 --partitions 3 --replication-factor 1
```

**Command Breakdown:**
- `docker exec -it kafka-kafka-1`: Execute command inside the Kafka container
- `kafka-topics`: Kafka topics management utility
- `--bootstrap-server localhost:9092`: Kafka broker address
- `--create`: Create a new topic
- `--topic pasture.telemetry.v1`: Name of the topic
- `--partitions 3`: Number of partitions for the topic
- `--replication-factor 1`: Number of replicas for each partition

### Listing Topics

To list all available topics and verify creation:

```bash
docker exec -it kafka-kafka-1 kafka-topics --bootstrap-server localhost:9092 --list
```

### Describing Topics

To get detailed information about a specific topic:

```bash
docker exec -it kafka-kafka-1 kafka-topics --bootstrap-server localhost:9092 --describe --topic pasture.telemetry.v1
```

### Deleting Topics

To delete a topic (use with caution):

```bash
docker exec -it kafka-kafka-1 kafka-topics --bootstrap-server localhost:9092 --delete --topic pasture.telemetry.v1
```

## Testing Topic Creation

### Complete Workflow Example

1. **Start the services:**
   ```bash
   docker-compose up -d
   ```

2. **Wait for Kafka to be ready** (usually 30-60 seconds):
   ```bash
   docker exec -it kafka-kafka-1 kafka-broker-api-versions --bootstrap-server localhost:9092
   ```

3. **Create the telemetry topic:**
   ```bash
   docker exec -it kafka-kafka-1 kafka-topics --bootstrap-server localhost:9092 --create --topic pasture.telemetry.v1 --partitions 3 --replication-factor 1
   ```

4. **Verify topic creation:**
   ```bash
   docker exec -it kafka-kafka-1 kafka-topics --bootstrap-server localhost:9092 --list
   ```

5. **Get topic details:**
   ```bash
   docker exec -it kafka-kafka-1 kafka-topics --bootstrap-server localhost:9092 --describe --topic pasture.telemetry.v1
   ```

6. **Consume messages from the topic:**
   ```bash
   docker exec -it kafka-kafka-1 kafka-console-consumer --bootstrap-server localhost:9092 --topic pasture.telemetry.v1 --from-beginning
   ```

## Troubleshooting

### Common Issues

1. **Kafka container not ready:**
   - Wait longer after `docker-compose up`
   - Check container status: `docker-compose ps`

2. **Connection refused:**
   - Verify Kafka container is running
   - Check if port 9092 is accessible

3. **Topic already exists:**
   - Use `--if-not-exists` flag or delete the topic first

### Useful Commands

- **Check container status:** `docker-compose ps`
- **View all containers:** `docker ps -a`
- **Enter container shell:** `docker exec -it kafka-kafka-1 bash`
- **View Kafka logs:** `docker-compose logs kafka`

## Topic Naming Convention

Topics in this project follow the pattern: `pasture.{domain}.{version}`

Examples:
- `pasture.telemetry.v1`
- `pasture.events.v1`
- `pasture.commands.v1`

## Additional Resources

- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
