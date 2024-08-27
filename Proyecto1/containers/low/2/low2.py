import time

# Este script consume poca RAM creando una lista pequeña
small_list = [0] * 10  # Lista muy pequeña
print("List created with", len(small_list), "elements.")

# Mantener el programa en ejecución para que puedas observar el consumo
while True:
    time.sleep(10)