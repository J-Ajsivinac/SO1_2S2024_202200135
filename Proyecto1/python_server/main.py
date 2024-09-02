import datetime
from fastapi import FastAPI # type: ignore
import os
import json
from typing import List
from models.models import LogProcess
import matplotlib.pyplot as plt
import pandas as pd

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

    # # Escribimos la lista de logs en el archivo logs.json
    # with open(logs_file, 'w') as file:
    #     json.dump(existing_logs, file, indent=4)

    return {"received": existing_logs}

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
    df = pd.DataFrame(existing_logs) 
    df['timestamp'] = pd.to_datetime(df['timestamp'])

    plt.figure(figsize=(10, 6))
    for container_id, group in df.groupby('container_id'):
        # group = group.sort_values('timestamp')
        # plt.plot(group['timestamp'], group['memory_usage'], label=container_id)
        plt.plot(group['timestamp'], group['cpu_usage'], marker='o', label=f'CPU - {container_id}')
    # output_path = 'home/ajsivinac/Documentos/memory_usage_graph.png'
    # plt.savefig(output_path)
    plt.xlabel('Timestamp')
    plt.ylabel('CPU Usage')
    plt.title('CPU Usage by Container')
    plt.legend(loc='upper left', bbox_to_anchor=(1, 1.06))
    plt.tight_layout(rect=[0, 0, 0.85, 1])  # Deja espacio a la derecha

    plt.grid(True)
    # output_path = 'home/ajsivinac/Documentos/memory_usage_graph.png'
    plt.savefig('cpu_usage_graph.png')

# Cerrar la figura para liberar memoria
    plt.close()
    return {"received": "True", "output_path": "output_path"}