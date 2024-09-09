import os
import json
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
from typing import List
from fastapi import FastAPI # type: ignore
from models.models import LogProcess
from fastapi.responses import JSONResponse
from collections import Counter

def get_line(data, type):
    # Contar la cantidad de procesos por timestamp
    counter = Counter(data['timestamp'])
    plt.style.use('dark_background')
    plt.tight_layout()
    # Convertir a DataFrame para facilidad de uso con matplotlib
    count_df = pd.DataFrame(counter.items(), columns=['timestamp', 'count'])

    # Ordenar por timestamp
    count_df = count_df.sort_values(by='timestamp')
    
    # Graficar
    plt.figure(figsize=(12*len(count_df.columns), 0.5*len(count_df)))
    plt.plot(count_df['timestamp'], count_df['count'], color='darkviolet', marker='o', linestyle='-', linewidth=2, label='Number of Processes')

    # Agregar área de degradado debajo de la línea
    plt.fill_between(count_df['timestamp'], count_df['count'], color='blueviolet', alpha=0.3)
    plt.xlabel('Tiempo')
    plt.ylabel('Numero de Procesos')
    plt.title('Numero de Procesos por Tiempo')
    # plt.grid()
    plt.grid(axis='both', linestyle='--', alpha=0.6)

    plt.xticks(rotation=70)
    

    if not os.path.exists('./imgs'):
        # Si no existe, la creamos
        os.makedirs('./imgs')

    plt.savefig(f'./imgs/{type}_graph.png', dpi=120, bbox_inches='tight')
    plt.close()

def get_bar(data, type):
    type_name = 'CPU' if type == 'cpu_usage' else 'Memoria'
    data_grouped = data.groupby('timestamp').agg({f'{type}': 'sum'}).reset_index()
   

    plt.figure(figsize=(7*len(data_grouped.columns), 2*len(data_grouped)))
    plt.tight_layout()
    plt.barh(data_grouped['timestamp'], data_grouped[type], color='mediumpurple', edgecolor='black', linewidth=1.2)

    plt.title(f'Suma de {type_name} por tiempo')
    plt.xlabel(f'Suma de %{type_name}')
    plt.ylabel('Tiempo')
    plt.grid(axis='x', linestyle='--', alpha=0.6)
    # Guardamos la gráfica en la carpeta imgs

    # verificamos si la carpeta imgs existe
    if not os.path.exists('./imgs'):
        # si no existe, la creamos
        os.makedirs('./imgs')
        
    plt.savefig(f'./imgs/{type_name}_graph.png', dpi=120)
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
    df["timestamp"] = pd.to_datetime(df["timestamp"])

    df = df.sort_values(by='timestamp')   

    df['timestamp'] = df['timestamp'].dt.strftime('%Y-%m-%d %H:%M:%S')

    get_line(df, 'process')

    # Gráfica de barras de uso de CPU
    get_bar(df, 'cpu_usage')
    # Gráfica de barras de uso de memoria
    get_bar(df, 'memory_usage')

    return {"message": "Graph created"}