# Usa la imagen base de Python 3.9 slim
FROM python:3.9-slim

# Establece el directorio de trabajo
WORKDIR /code

# Copia el archivo de requisitos al contenedor
COPY ./requirements.txt /code/requirements.txt

# Actualiza pip y luego instala las dependencias
RUN pip install --no-cache-dir --upgrade pip \
    && pip install --no-cache-dir --upgrade -r /code/requirements.txt

# Copia el código fuente al contenedor
COPY ./ /code/

# Comando para ejecutar la aplicación FastAPI
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
