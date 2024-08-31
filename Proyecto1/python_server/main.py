import datetime
from fastapi import FastAPI # type: ignore
import os
import json
from typing import List
from models.models import LogProcess
import matplotlib.pyplot as plt

app = FastAPI()


@app.get("/")
def read_root():
    return {"Hello": "test"}


@app.post("/logs")
def get_logs(logs_proc: List[LogProcess]):
    logs_file = 'logs/logs.json'
    
    # Checamos si existe el archivo logs.json
    if os.path.exists(logs_file):
        # Leemos el archivo logs.json
        with open(logs_file, 'r') as file:
            existing_logs = json.load(file)
    else:
        # Sino existe, creamos una lista vacía
        existing_logs = []

    # Agregamos los nuevos logs a la lista existente
    new_logs = [log.dict() for log in logs_proc]
    existing_logs.extend(new_logs)

    # Escribimos la lista de logs en el archivo logs.json
    with open(logs_file, 'w') as file:
        json.dump(existing_logs, file, indent=4)

    return {"received": True}

@app.get("/graph")
def get_graph():
    logs_file = 'logs/logs.json'
    
    # Checamos si existe el archivo logs.json
    if os.path.exists(logs_file):
        # Leemos el archivo logs.json
        with open(logs_file, 'r') as file:
            existing_logs = json.load(file)
    else:
        # Sino existe, creamos una lista vacía
        existing_logs = []

    # Obtenemos los time_stamps de los logs (solo los que no se repiten)
    time_stamps = list(set([log["timestamp"] for log in existing_logs]))
    
    # agrupar los logs del mismo time_stamp
    grouped_logs = {}
    for time_stamp in time_stamps:
        grouped_logs[time_stamp] = [log for log in existing_logs if log["timestamp"] == time_stamp]

    # memory_usage = [log["timestamp"] for log in grouped_logs]

    # obtjener el procentaje de uso por fechas
    memory_usage = []
    for time_stamp in time_stamps:
        logs = grouped_logs[time_stamp]
        total_memory = sum([log["memory_usage"] for log in logs])
        memory_usage.append(total_memory)

    # plt.figure(figsize=(10, 5))
    # plt.plot(time_stamps, memory_usage, marker='o', linestyle='-', color='b')
    # plt.xlabel('Timestamp')
    # plt.ylabel('Memory Usage (%)')
    # plt.title('Memory Usage Over Time')
    # plt.grid(True)
    # plt.xticks(rotation=45)
    # plt.tight_layout()

    output_path = 'home/ajsivinac/Documentos/memory_usage_graph.png'
    plt.savefig(output_path)

# Cerrar la figura para liberar memoria
    plt.close()
    return {"received": memory_usage}