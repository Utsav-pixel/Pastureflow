package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/Utsav-pixel/gosense"
)

var (
	configFile = flag.String("config", "config/pastures.json", "Path to pasture configuration file")
	duration   = flag.Duration("duration", 30*time.Second, "How long to run the engine")
	publisher  = flag.String("publisher", "kafka", "Publisher type (kafka, http, console)")
	brokers    = flag.String("brokers", "localhost:9092", "Kafka broker addresses")
	topic      = flag.String("topic", "pasture.telemetry.v1", "Kafka topic name")
	endpoint   = flag.String("endpoint", "https://api.example.com/telemetry", "HTTP endpoint URL")
	verbose    = flag.Bool("verbose", false, "Enable verbose logging")

	// NEW: Configuration controls
	productionRate = flag.Duration("rate", 100*time.Millisecond, "Data generation interval (e.g., 100ms, 1s)")
	batchSize      = flag.Int("batch", 100, "Number of items per batch")
	batchTimeout   = flag.Duration("timeout", 500*time.Millisecond, "How long to wait before publishing a batch")
	maxWorkers     = flag.Int("workers", 3, "Number of concurrent workers")
	profile        = flag.String("profile", "default", "Performance profile (default, high-throughput, low-latency)")
)

type PastureTelemetry struct {
	ZoneID       string    `json:"zone_id"`
	Timestamp    time.Time `json:"ts"`
	Biomass      float64   `json:"biomass"`
	SoilMoisture float64   `json:"soil_moisture"`
	Temperature  float64   `json:"temperature"`
	AnimalLoad   int       `json:"animal_load"`
	Quality      string    `json:"quality"`
}

type PastureConfig struct {
	ID                  string  `json:"id"`
	AreaHectares        float64 `json:"area_hectares"`
	InitialBiomass      float64 `json:"initial_biomass"`
	InitialSoilMoisture float64 `json:"initial_soil_moisture"`
	InitialAnimalLoad   int     `json:"initial_animal_load"`
	RegenRate           float64 `json:"regen_rate"`
	DegradeRate         float64 `json:"degrade_rate"`
	EvapRate            float64 `json:"evap_rate"`
}

// PastureSensorFunction generates realistic pasture telemetry
type PastureSensorFunction struct {
	configs []PastureConfig
	zoneIdx int
}

func NewPastureSensorFunction(configs []PastureConfig) *PastureSensorFunction {
	return &PastureSensorFunction{configs: configs, zoneIdx: 0}
}

func (p *PastureSensorFunction) Generate(input float64, timestamp time.Time) PastureTelemetry {
	config := p.configs[p.zoneIdx]
	p.zoneIdx = (p.zoneIdx + 1) % len(p.configs)

	hour := float64(timestamp.Hour())
	baseTemp := 25.0 + 5.0*math.Sin((hour/24.0)*2*math.Pi-math.Pi/2)
	temperature := baseTemp + (input-0.5)*10.0

	evapLoss := temperature * config.EvapRate * 0.01
	soilMoisture := math.Max(0, config.InitialSoilMoisture-evapLoss+0.05)

	degradation := float64(config.InitialAnimalLoad) * config.DegradeRate * 0.001
	regeneration := config.RegenRate * soilMoisture * 0.01
	biomass := math.Max(0, math.Min(1, config.InitialBiomass-degradation+regeneration+0.02))

	return PastureTelemetry{
		ZoneID:       config.ID,
		Timestamp:    timestamp,
		Biomass:      biomass,
		SoilMoisture: soilMoisture,
		Temperature:  temperature,
		AnimalLoad:   config.InitialAnimalLoad,
		Quality:      "OK",
	}
}

func main() {
	flag.Parse()

	if *verbose {
		log.Printf("Starting Pasture Sensor Engine - Using Gosense v0.2.3")
		log.Printf("Config file: %s", *configFile)
		log.Printf("Duration: %v", *duration)
		log.Printf("Publisher: %s", *publisher)
		log.Printf("Performance Profile: %s", *profile)
		log.Printf("Production Rate: %v", *productionRate)
		log.Printf("Batch Size: %d", *batchSize)
		log.Printf("Batch Timeout: %v", *batchTimeout)
		log.Printf("Max Workers: %d", *maxWorkers)
	}

	configs, err := loadPastureConfigs(*configFile)
	if err != nil {
		log.Fatalf("Failed to load pasture configs: %v", err)
	}

	// Create publisher based on type
	pub, err := createPublisher(*publisher)
	if err != nil {
		log.Fatalf("Failed to create publisher: %v", err)
	}
	defer pub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), *duration)
	defer cancel()

	log.Printf("Starting gosense engine for %v...", *duration)
	startTime := time.Now()

	// Use gosense components directly
	seeder := gosense.NewTimeSeeder(1.0, 0.1, 0.5)
	sensorFunc := gosense.NewFunction(NewPastureSensorFunction(configs).Generate)

	// NEW: Create custom configuration based on profile or custom settings
	var engineConfig gosense.Config
	switch *profile {
	case "high-throughput":
		engineConfig = gosense.HighThroughputConfig()
		if *verbose {
			log.Printf("Using High-Throughput Profile")
		}
	case "low-latency":
		engineConfig = gosense.LowLatencyConfig()
		if *verbose {
			log.Printf("Using Low-Latency Profile")
		}
	default:
		// Custom configuration from command line flags
		engineConfig = gosense.Config{
			ProductionRate: *productionRate,
			BatchSize:      *batchSize,
			BatchTimeout:   *batchTimeout,
			MaxWorkers:     *maxWorkers,
		}
		if *verbose {
			log.Printf("Using Custom Configuration")
		}
	}

	gosenseEngine := gosense.NewEngine(engineConfig, seeder, sensorFunc, pub)

	if err := gosenseEngine.Start(ctx); err != nil {
		log.Fatalf("Engine error: %v", err)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\n=== Gosense v0.2.3 Engine Statistics ===\n")
	fmt.Printf("Total Duration: %v\n", elapsed)
	fmt.Printf("Production Rate: %v\n", engineConfig.ProductionRate)
	fmt.Printf("Batch Size: %d\n", engineConfig.BatchSize)
	fmt.Printf("Batch Timeout: %v\n", engineConfig.BatchTimeout)
	fmt.Printf("Max Workers: %d\n", engineConfig.MaxWorkers)
	fmt.Printf("Engine completed successfully\n")
	fmt.Printf("========================\n")
}

func loadPastureConfigs(filename string) ([]PastureConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var wrapper struct {
		Zones []PastureConfig `json:"zones"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return wrapper.Zones, nil
}

func createPublisher(pubType string) (gosense.Publisher[PastureTelemetry], error) {
	switch pubType {
	case "kafka":
		return gosense.NewGenericKafkaPublisher[PastureTelemetry]([]string{*brokers}, *topic), nil
	case "http":
		return gosense.NewGenericHTTPPublisher[PastureTelemetry](*endpoint), nil
	case "grpc":
		return gosense.NewGenericGRPCPublisher[PastureTelemetry](*endpoint)
	case "console":
		return gosense.NewConsolePublisher[PastureTelemetry](), nil
	default:
		return nil, fmt.Errorf("unsupported publisher type: %s", pubType)
	}
}
