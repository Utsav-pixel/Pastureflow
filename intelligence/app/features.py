from collections import deque
import numpy as np

WINDOW_SIZE = 10

class FeatureStore:
    def __init__(self):
        self.buffer = {}
    
    def update(self,telemetry):
        zone = telemetry.zone_id
        if zone not in self.buffer:
            self.buffer[zone] = deque(maxlen=WINDOW_SIZE)
        self.buffer[zone].append(telemetry)

    def compute(self,zone_id):
        buf = self.buffer.get(zone_id)
        if not buf or len(buf)<WINDOW_SIZE:
            return None

        biomass = np.array([t.metrics.biomass for t in buf])
        soil = np.array([t.metrics.soil_moisture for t in buf])
        load = np.array([t.metrics.animal_load for t in buf])

        return {
            "biomass_mean": biomass.mean(),
            "soil_mean": soil.mean(),
            "biomass_trend": biomass[-1] - biomass[0],
            "animal_load_mean": load.mean(),
            "pressure_index": (load.mean()/(biomass.mean() + 0.1)),
        }