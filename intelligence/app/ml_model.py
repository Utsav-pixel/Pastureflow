import joblib
import numpy as np

model = joblib.load("risk_model.pkl")

def predict_risk(features):
    vector = np.array([[
        features["biomass_mean"],
        features["biomass_trend"],
        features["soil_mean"],
        features["animal_load_mean"],
        features["pressure_index"]
    ]])

    return model.predict(vector)[0]