from fastapi import FastAPI
from .models import Telemetry
from .features import FeatureStore
from .decisions import classify_risk
from .ml_model import predict_risk
from .consumer import start_consumer
import threading
from .training_logger import log_training_sample
app = FastAPI()
store = FeatureStore()

def handle_telemetry(t:Telemetry):
    store.update(t)
    features = store.compute(t.zone_id)
    decision = classify_risk(features)
    print(f"Zone {t.zone_id}: Features={features}, Decision={decision}")
    log_training_sample(t.zone_id, features, decision)
    return {"features": features, "decision": decision}

@app.get("/")
def read_root():
    return {"message": "Pastureflow Intelligence API", "status": "running"}

@app.get("/health")
def health_check():
    return {"status": "healthy", "zones": len(store.buffer)}

@app.get("/zones/{zone_id}/features")
def get_zone_features(zone_id: str):
    features = store.compute(zone_id)
    if features:
        return {"zone_id": zone_id, "features": features}
    return {"zone_id": zone_id, "error": "No data available"}

@app.get("/zones/{zone_id}/decision")
def get_zone_decision(zone_id: str):
    features = store.compute(zone_id)
    if features:
        decision = predict_risk(features)
        # log_training_sample(zone_id, features, decision)
        return {"zone_id": zone_id, "features": features, "decision": decision}
    return {"zone_id": zone_id, "error": "No data available"}

@app.on_event("startup")
def startup():
    print("Starting Kafka consumer...")
    thread = threading.Thread(target=start_consumer, args=(handle_telemetry,), daemon=True)
    thread.start()
    print("Kafka consumer started successfully")