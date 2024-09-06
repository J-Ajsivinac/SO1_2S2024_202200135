import os
import json
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
from typing import List
from fastapi import FastAPI # type: ignore
from models.models import LogProcess
from fastapi.responses import JSONResponse

def gen_heatmap(data, type):
    type_name = 'CPU' if type == 'cpu_usage' else 'Memoria'
    data['container_id'] = data['container_id'].str[:12]
    df_heatmap = data.groupby(['timestamp', 'container_id'])[type].mean().unstack()
    plt.style.use('dark_background')
    print(len(df_heatmap))
    plt.figure(figsize=(2*len(df_heatmap.columns), 2.2*len(df_heatmap)))
    plt.tight_layout()
    sns.heatmap(df_heatmap, cmap=sns.cubehelix_palette(as_cmap=True), annot=True, fmt=".2f", linewidths=0.5, square=False)
    plt.title(f'Heatmap de uso de {type_name}')
    plt.xlabel('Container ID')
    plt.ylabel('Tiempo')
    # plt.grid(True)

    # Guardamos la gráfica en la carpeta imgs

    # verificamos si la carpeta imgs existe
    if not os.path.exists('./imgs'):
        # si no existe, la creamos
        os.makedirs('./imgs')

    plt.savefig(f'./imgs/{type}_graph.png', dpi=300)
    plt.close()

def get_bar(data, type):
    type_name = 'CPU' if type == 'cpu_usage' else 'Memoria'
    data_grouped = data.groupby('timestamp').agg({f'{type}': 'sum'}).reset_index()

    plt.figure(figsize=(7.5*len(data_grouped.columns), 2.3*len(data_grouped)))
    plt.tight_layout()
    plt.barh(data_grouped['timestamp'], data_grouped[type], color='#a7aff1')

    plt.title(f'Suma de {type_name} por tiempo')
    plt.xlabel(f'% {type_name}')
    plt.ylabel('Tiempo')

    # Guardamos la gráfica en la carpeta imgs

    # verificamos si la carpeta imgs existe
    if not os.path.exists('./imgs'):
        # si no existe, la creamos
        os.makedirs('./imgs')
        
    plt.savefig(f'./imgs/{type_name}_graph.png', dpi=300)
    plt.close()

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

    return {"received": existing_logs}

@app.get("/graph")
def get_graph():
    logs_file = 'logs/logs.json'
    
    if os.path.exists(logs_file):
        with open(logs_file, 'r') as file:
            existing_logs = json.load(file)
    else:
        existing_logs = []
        response = {"message": "No logs found"}
        return JSONResponse(content=response, status_code=404)
    
    if existing_logs == []:
        response = {"message": "No logs found"}
        return JSONResponse(content=response, status_code=404)
    
    df = pd.DataFrame(existing_logs)
    df["timestamp"] = pd.to_datetime(df["timestamp"]).dt.strftime('%Y-%m-%d %H:%M:%S')
    
    # Heatmap de uso de CPU
    gen_heatmap(df, 'cpu_usage')
    # Heatmap de uso de memoria
    gen_heatmap(df, 'memory_usage')

    # Gráfica de barras de uso de CPU
    get_bar(df, 'cpu_usage')
    # Gráfica de barras de uso de memoria
    get_bar(df, 'memory_usage')

    return {"message": "Graph created"}