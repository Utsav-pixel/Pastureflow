import csv
import os
from datetime import datetime

FILE = "training_data.csv"

def log_training_sample(zone_id, features, risk_label):
    file_exists = os.path.isfile(FILE)

    with open(FILE, mode="a", newline="") as f:
        writer = csv.writer(f)

        if not file_exists:
            writer.writerow([
                "timestamp",
                "zone_id",
                "biomass_mean",
                "biomass_trend",
                "soil_mean",
                "animal_load_mean",
                "pressure_index",
                "risk_label"
            ])

        writer.writerow([
            datetime.utcnow().isoformat(),
            zone_id,
            features["biomass_mean"],
            features["biomass_trend"],
            features["soil_mean"],
            features["animal_load_mean"],
            features["pressure_index"],
            risk_label
        ])