import pandas as pd
from sklearn.ensemble import RandomForestClassifier
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report
import joblib

# Load data
df = pd.read_csv("training_data.csv")

X = df[[
    "biomass_mean",
    "biomass_trend",
    "soil_mean",
    "animal_load_mean",
    "pressure_index"
]]

y = df["risk_label"]

X_train, X_test, y_train, y_test = train_test_split(
    X, y, test_size=0.2, random_state=42
)

model = RandomForestClassifier(
    n_estimators=100,
    max_depth=5,
    random_state=42
)

model.fit(X_train, y_train)

preds = model.predict(X_test)
print(classification_report(y_test, preds))

joblib.dump(model, "risk_model.pkl")

print("Model saved to risk_model.pkl")