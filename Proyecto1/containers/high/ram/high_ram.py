import time

# Crear una lista muy grande que consuma mucha memoria
large_list = [0] * (10**8)  # Ajusta este número para consumir más RAM

# Mantener el programa en ejecución para que puedas observar el consumo
while True:
    time.sleep(10)