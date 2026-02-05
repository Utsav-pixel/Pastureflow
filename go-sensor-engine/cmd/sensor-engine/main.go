package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Utsav-pixel/go-sensor-engine/internal/publisher"
	"github.com/Utsav-pixel/go-sensor-engine/internal/sim"
)

func main() {

	tickerInterval, zones := ReadPastureZones()
	telemetry_chan := make(chan sim.Telemetry, len(zones[:4]))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := range zones[:3] {
		zone := &zones[i]
		go runZoneScheduler(ctx, zone, tickerInterval, telemetry_chan)
	}
	go consumeTelemetry(ctx, telemetry_chan)
	// Use tickerInterval to avoid unused variable warning
	<-ctx.Done()
	fmt.Println("Context cancelled or expired")
	fmt.Printf("Ticker interval: %v\n", tickerInterval)

}

func runZoneScheduler(ctx context.Context, zone *sim.PastureZone, tickerInterval time.Duration, out chan<- sim.Telemetry) {
	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			produced := zone.Tick(time.Now())
			fmt.Printf("Produced telemetry for zone %s at %v-------%v\n", produced.ZoneID, produced.Timestamp, produced.Metrics.Temperature)
			out <- produced
		case <-ctx.Done():
			return
		}
	}
}

func ReadPastureZones() (tickerInterval time.Duration, zones []sim.PastureZone) {
	// Open the config file
	file, err := os.Open("config/pastures.json")
	if err != nil {
		fmt.Printf("Error opening config file: %v\n", err)
		return time.Second, []sim.PastureZone{}
	}
	defer file.Close()

	// Read file using bufio
	scanner := bufio.NewScanner(file)
	var jsonContent string
	for scanner.Scan() {
		jsonContent += scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return time.Second, []sim.PastureZone{}
	}

	// Parse JSON
	var config struct {
		TickMs int `json:"tick_ms"`
		Zones  []struct {
			ID                  string  `json:"id"`
			AreaHectares        float64 `json:"area_hectares"`
			Lat                 float64 `json:"lat"`
			Lon                 float64 `json:"lon"`
			InitialBiomass      float64 `json:"initial_biomass"`
			InitialSoilMoisture float64 `json:"initial_soil_moisture"`
			InitialAnimalLoad   int     `json:"initial_animal_load"`
			RegenRate           float64 `json:"regen_rate"`
			DegradeRate         float64 `json:"degrade_rate"`
			EvapRate            float64 `json:"evap_rate"`
		} `json:"zones"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &config); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return time.Second, []sim.PastureZone{}
	}

	// Convert to PastureZone structs
	for _, zoneConfig := range config.Zones {
		zone := sim.PastureZone{
			ZoneID:       zoneConfig.ID,
			AreaHectares: zoneConfig.AreaHectares,
			Biomass:      zoneConfig.InitialBiomass,
			SoilMoisture: zoneConfig.InitialSoilMoisture,
			AnimalLoad:   zoneConfig.InitialAnimalLoad,
			RegenRate:    zoneConfig.RegenRate,
			DegradeRate:  zoneConfig.DegradeRate,
			EvapRate:     zoneConfig.EvapRate,
		}
		zones = append(zones, zone)
	}

	tickerInterval = time.Duration(config.TickMs) * time.Millisecond
	return tickerInterval, zones
}

func consumeTelemetry(ctx context.Context, in <-chan sim.Telemetry) {
	// pub := publisher.NewHTTPPublisher("https://webhook.site/dc4fff2a-f1a9-4fab-96f4-1fb0f6f97182")
	pub := publisher.NewKafkaPublisher([]string{"localhost:9092"}, "pasture.telemetry.v1")
	defer pub.Close()
	for {
		select {
		case telemetry := <-in:

			// formatPrint(telemetry)
			// pub.Publish(ctx, telemetry)
			if err := pub.Publish(ctx, telemetry); err != nil {
				fmt.Printf("Error publishing telemetry: %v\n", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
func formatPrint(telemetry sim.Telemetry) {
	buf := make([]byte, 0, 256)
	buf = append(buf, "ts="...)
	buf = telemetry.Timestamp.AppendFormat(buf, time.RFC3339Nano)
	buf = append(buf, " zone="...)
	buf = append(buf, []byte(telemetry.ZoneID)...)
	buf = append(buf, " biomass="...)
	buf = strconv.AppendFloat(buf, telemetry.Metrics.Biomass, 'f', 3, 64)
	buf = append(buf, " soil="...)
	buf = strconv.AppendFloat(buf, telemetry.Metrics.SoilMoisture, 'f', 3, 64)
	buf = append(buf, " temp_c="...)
	buf = strconv.AppendFloat(buf, telemetry.Metrics.Temperature, 'f', 2, 64)
	buf = append(buf, " animals="...)
	buf = strconv.AppendInt(buf, int64(telemetry.Metrics.AnimalLoad), 10)
	buf = append(buf, " quality="...)
	buf = append(buf, []byte(telemetry.Quality)...)
	buf = append(buf, '\n')
	_, _ = os.Stdout.Write(buf)
}
