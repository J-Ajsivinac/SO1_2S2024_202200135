import os
import json
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
from typing import List
from fastapi import FastAPI # type: ignore
from models.models import LogProcess
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
    
    if os.path.exists(logs_file):
        with open(logs_file, 'r') as file:
            existing_logs = json.load(file)
    else:
        existing_logs = []

    df = pd.DataFrame(existing_logs)
    df["timestamp"] = pd.to_datetime(df["timestamp"]).dt.strftime('%Y-%m-%d %H:%M:%S')
    df_heatmap = df.groupby(['timestamp', 'container_id'])['cpu_usage'].mean().unstack()
    
    plt.style.use('dark_background')
    sns.heatmap(df_heatmap, cmap=sns.cubehelix_palette(as_cmap=True), annot=True, fmt=".2f", linewidth=0.5)
    plt.title('Heatmap de uso de CPU')
    plt.xlabel('Container ID')
    plt.ylabel('Tiempo')
    plt.grid(True)
    plt.savefig('cpu_usage_graph.png', dpi=300)
    plt.close()

    # Heatmap de uso de memoria
    df_memory_heatmap = df.groupby(['timestamp', 'container_id'])['memory_usage'].mean().unstack()
    plt.figure(figsize=(10, 6))
    plt.style.use('dark_background')
    sns.heatmap(df_memory_heatmap, cmap=sns.cubehelix_palette(as_cmap=True), annot=True, fmt=".2f", linewidths=0.5)
    plt.title('Heatmap de uso de Memoria')
    plt.xlabel('Container ID')
    plt.ylabel('Tiempo')
    plt.grid(True)
    plt.savefig('memory_usage_graph.png', dpi=300)
    plt.close()
    return {"message": "Graph created"}


def gen_image(data, type):
    df_heatmap = data.groupby(['timestamp', 'container_id'])[type].mean().unstack()
    plt.style.use('dark_background')
    sns.heatmap(df_heatmap, cmap=sns.cubehelix_palette(as_cmap=True), annot=True, fmt=".2f", linewidths=0.5)
    plt.title(f'Heatmap de uso de {type}')
    plt.xlabel('Container ID')
    plt.ylabel('Tiempo')
    plt.grid(True)
    plt.savefig(f'{type}_usage_graph.png', dpi=300)
    plt.close()