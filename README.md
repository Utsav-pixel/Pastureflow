# Pastureflow 🌱

**Pastureflow** is an event-driven pasture intelligence system designed to simulate real-world agricultural telemetry and feed it into downstream machine learning decision models.

The project models pasture zones as dynamic systems (biomass, soil moisture, animal load, environment), streams their state as telemetry events through Kafka, and provides a clean foundation for building ML-driven decision intelligence on top.

This is **not a dashboard project**.  
It is an **ML-first data and decision pipeline**, built the way real-world ag-tech and IoT-adjacent systems are designed.

---

## Why Pastureflow exists

Modern ag-tech systems don't start with ML models — they start with **reliable signals**.

Pastureflow focuses on:
- generating realistic, structured telemetry
- enforcing strong event contracts
- decoupling data production from intelligence
- enabling replayable, observable ML pipelines

The result is a system that mirrors how production-grade ML infrastructure works in practice.

---

## High-level architecture
```
┌──────────────────┐
│ Go Sensor Engine │
│ (Pasture Zones)  │
└────────┬─────────┘
│ Telemetry (JSON)
▼
┌──────────────────┐
│ Kafka            │
│ pasture.*        │
└────────┬─────────┘
│ Stream
▼
┌──────────────────┐
│ Intelligence     │
│ (FastAPI + ML)   │
└──────────────────┘
```
- **Go** is used for deterministic, high-performance simulation
- **Kafka** acts as the system backbone (durability, ordering, replay)
- **Python / FastAPI** consumes telemetry for ML-based decisions (next phase)

---

## Repository structure
```
pastureflow/
├── go-sensor-engine/          # Go-based pasture & sensor simulation
│   ├── internal/
│   │   ├── sim/            # Domain models & simulation logic
│   │   └── publisher/      # Kafka publisher
│   ├── cmd/
│   │   └── sensor-engine/  # Main application entry point
│   └── config/             # Configuration files
│
├── infra/                    # Infrastructure definitions
│   └── kafka/               # Docker-based Kafka infrastructure
│
├── docs/                    # Deep-dive documentation
│   ├── docker-kafka-operations.md  # Kafka setup & operations
│
└── README.md
```
Each module is intentionally small, explicit, and independently evolvable.
<!-- │   ├── simulation.md       # Pasture & sensor modeling (planned) -->
<!-- │   └── intelligence.md     # ML & decision layer (planned) -->
---

## Core components

### 1. Sensor Engine (Go)
The sensor engine simulates pasture zones as stateful systems evolving over time.

Each **PastureZone**:
- tracks biomass, soil moisture, temperature, and animal load
- evolves deterministically on each tick
- emits structured `Telemetry` events

Concurrency is handled using goroutines and channels, mirroring real sensor fleets.

📄 Detailed design: [`docs/simulation.md`](docs/simulation.md) (planned)

---

### 2. Event Backbone (Kafka)
Kafka is used as the **single source of truth for telemetry events**.

Design principles:
- append-only event log
- deterministic partitioning by `zone_id` 
- JSON-encoded, language-agnostic contracts
- replayable streams for ML experimentation

Primary topic:
- `pasture.telemetry.v1`

📄 Kafka setup & operations: [`docs/docker-kafka-operations.md`](docs/docker-kafka-operations.md)

---

### 3. Intelligence Layer (FastAPI + ML)
The intelligence layer consumes Kafka telemetry streams and turns signals into decisions.

Planned capabilities:
- feature extraction from telemetry windows
- overgrazing risk detection
- pasture rotation recommendations
- alerting and decision outputs
- model evaluation and iteration via replay

This layer is intentionally decoupled from data production.

📄 Design notes: [`docs/intelligence.md`](docs/intelligence.md) (planned)

---

## Running the project (local)

### Start Kafka
```bash
cd infra/kafka
docker compose up -d
```

### Create topic
```bash
docker exec -it kafka-kafka-1 kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic pasture.telemetry.v1 \
  --partitions 3 \
  --replication-factor 1
```

### Run sensor engine
```bash
cd go-sensor-engine
go run cmd/sensor-engine/main.go
```


## Roadmap

- [x] Pasture simulation engine (Go)
- [x] Kafka-backed telemetry streaming
- [ ] FastAPI Kafka consumer
- [ ] Feature extraction layer
- [ ] Baseline decision models
- [ ] ML-based optimization & forecasting
- [ ] Feedback loop into pasture control signals

---

## Who this project is for

- ML Infrastructure Engineers
- Applied ML / Decision Systems Engineers
- Backend Engineers working on streaming systems
- Anyone interested in how ML systems actually start

---

## License

MIT
