package sim

import (
	"math"
	"math/rand/v2"
	"time"
)

const (
	MaxValue = 1.0
	MinValue = 0.0
)

// clamp.keep the value between min and max
func clamp(value float64) float64 {
	if value < MinValue {
		return MinValue
	}
	if value > MaxValue {
		return MaxValue
	}
	return value
}

func ComputeTemperature(baseTemp float64, t time.Time) float64 {
	hour := float64(t.Hour()) + float64(t.Minute())/60.0
	radian := (hour / 24.0) * 2 * math.Pi

	diurnal := 6 * math.Sin(radian-math.Pi/2) // peak afternoon
	noise := rand.NormFloat64() * 0.5         // random noise
	return baseTemp + diurnal + noise
}

func UpdateSoilMoisture(
	current float64,
	temperature float64,
	evapRate float64,
) (next float64, rain bool) {

	evapLoss := temperature * evapRate
	next = current - evapLoss

	// Rain simulation: 3% chance of rain each step
	if rand.Float64() < 0.03 {
		rainAmount := 0.1 + rand.Float64()*0.2 // 0.1 to 0.3 range
		next += rainAmount
		rain = true
	}

	return clamp(next), rain
}

func UpdateBiomass(
	current float64,
	regenRate float64,
	degradeRate float64,
	animalLoad int,
	soilMoisture float64,
	temperature float64,
) float64 {

	regrowth := 0.5

	if soilMoisture > 0.3 {
		regrowth = regenRate // Grass only growths if moisture is sufficient

		//Heat stress
		if temperature > 35 {
			regrowth *= 0.5
		}
	}
	grazingImpact := float64(animalLoad) * degradeRate

	return clamp(current + regrowth - grazingImpact)
}

func ApplySensorNoise(value float64, stdDev float64) float64 {
	return clamp(value + rand.NormFloat64()*stdDev)
}

type Quality string

const (
	QualityOK      Quality = "OK"
	QualityNoisy   Quality = "NOISY"
	QualityPartial Quality = "PARTIAL"
	QualityCorrupt Quality = "CORRUPT"
)

func DetermineQuality() Quality {
	r := rand.Float64()
	switch {
	case r < 0.01:
		return QualityCorrupt
	case r < 0.03:
		return QualityPartial
	case r < 0.08:
		return QualityNoisy
	default:
		return QualityOK
	}
}
