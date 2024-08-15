#!/usr/bin/env bash

# Función para generar nombres aleatorios
gen_random_name(){
    # Generar un nombre aleatorio de 8 caracteres
    echo "container_$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1)"
}
# Generar 10 contenedores con nombres aleatorios
num_containers=10

# ciclo para crear los contenedores
for i in $(seq 1 $num_containers); do
    # Generar un nombre aleatorio
    name=$(gen_random_name)
    # Crear el contenedor
    docker run -d --name $name alpine:latest sleep 120
    echo "Container $name created"
done

echo "All containers created"