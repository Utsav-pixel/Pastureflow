from pydantic import BaseModel
from datetime import datetime


class Metrics(BaseModel):
    biomass: float
    soil_moisture: float
    temperature: float
    animal_load: int

class Telemetry(BaseModel):
    ts: datetime
    zone_id: str
    metrics: Metrics
    quality: str
