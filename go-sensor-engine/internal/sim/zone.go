package sim

import (
	"math/rand"
	"time"
)

type PastureZone struct {
	ZoneID       string
	AreaHectares float64

	// Dynamic state
	Biomass      float64 // 0.0 - 1.0
	SoilMoisture float64 // 0.0 - 1.0
	Temperature  float64 // Celsius
	AnimalLoad   int

	// Environment param
	RegenRate   float64 // grass regrowth rate per tick
	DegradeRate float64 // grazing impact per animal
	EvapRate    float64 // moisture loss per temp unit
}

type Telemetry struct {
	ZoneID    string
	Timestamp time.Time
	Metrics   Metrics
	Quality   Quality
}

type Metrics struct {
	Biomass      float64
	SoilMoisture float64
	Temperature  float64
	AnimalLoad   int
}

func (z *PastureZone) Tick(now time.Time) Telemetry {
	z.Temperature = ComputeTemperature(25.0, now)

	z.SoilMoisture, _ = UpdateSoilMoisture(z.SoilMoisture, z.Temperature, z.EvapRate)
	z.Biomass = UpdateBiomass(z.Biomass, z.RegenRate, z.DegradeRate, z.AnimalLoad, z.SoilMoisture, z.Temperature)

	quality := DetermineQuality()

	biomass := ApplySensorNoise(z.Biomass, 0.05)

	soilMoisture := ApplySensorNoise(z.SoilMoisture, 0.05)

	temperature := z.Temperature + rand.NormFloat64()*0.4

	return Telemetry{
		ZoneID:    z.ZoneID,
		Timestamp: now,
		Metrics: Metrics{
			Biomass:      biomass,
			SoilMoisture: soilMoisture,
			Temperature:  temperature,
			AnimalLoad:   z.AnimalLoad,
		},
		Quality: quality,
	}
}
