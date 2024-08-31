#!/bin/bash

# Lista de imágenes base
images=("cpu-high" "ram-high" "cpu-low" "ram-low")

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
  docker run -d --name "$container_name" "$base_image"
  
  echo "Contenedor $i generado: $container_name usando $base_image"
done
