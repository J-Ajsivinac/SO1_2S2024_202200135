#!/bin/bash

# Lista de imágenes base
images=("high_ram" "high_cpu" "low_1" "low_2")

# Función para generar un nombre único de contenedor
generate_container_name() {
  echo "container_$(date +%s%N)"
}

# Crear 10 contenedores
for i in {1..10}; do
  # Seleccionar una imagen base aleatoriamente
  base_image=${images[$RANDOM % ${#images[@]}]}
  
  # Generar un nombre de contenedor único
  container_name=$(generate_container_name)
  
  # Ejecutar el contenedor
#   docker run -d --name "$container_name" "$base_image"
  
  echo "Contenedor $i generado: $container_name usando $base_image"
done