def decide(features):
    if features is None:
        return None
    if features["biomass_mean"] < 5.0:
        return {
            "decision": "ROTATE_PASTURE",
            "risk": "HIGH"
        }

    if features["biomass_trend"] < -1.0:
        return {
            "decision": "ALERT_DEGRADING",
            "risk": "MEDIUM"
        }

    return {
        "decision": "NO_ACTION",
        "risk": "LOW"
    }


def classify_risk(features):
    if features is None:
        return None

    score = 0
    biomass = features["biomass_mean"]
    trend = features["biomass_trend"]
    pressure = features["pressure_index"]

    # Biomass level (0–1 scale)
    if biomass < 0.25:
        score += 2
    elif biomass < 0.45:
        score += 1

    # Biomass trend (small deltas)
    if trend < -0.02:
        score += 2
    elif trend < -0.01:
        score += 1

    # Grazing pressure
    if pressure > 15:
        score += 2
    elif pressure > 8:
        score += 1

    if score >= 5:
        return "HIGH"
    elif score >= 3:
        return "MEDIUM"
    return "LOW"