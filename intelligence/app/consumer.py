import json 
from .models import Telemetry
from kafka import KafkaConsumer


def start_consumer(handler):
    consumer = KafkaConsumer(
        'pasture.telemetry.v1',
        bootstrap_servers=['localhost:9092'],
        value_deserializer=lambda x: json.loads(x.decode('utf-8')),
        auto_offset_reset='latest',
        enable_auto_commit=True,
    )
    
    counter_message = 0 
    for message in consumer:
        counter_message += 1
        # print(f"Raw JSON: {message.value}")  # Debug: Show raw JSON
        try:
            telemetry = Telemetry(**message.value)
            handler(telemetry)
        except Exception as e:
            print(f"Error parsing telemetry: {e}")
            print(f"JSON data: {message.value}")
    
    print(f"Total messages consumed: {counter_message}")
