import json
from locust import HttpUser, task, between
class StudentTraffic(HttpUser):
    wait_time = between(1, 5)
    def on_start(self):
        # Cargar el archivo JSON al iniciar Locust
        with open('estudiantes.json', 'r') as file:
            self.students_data = json.load(file)
    @task
    def send_traffic(self):
        # Iterar sobre los datos cargados desde el archivo JSON
        for student in self.students_data:
            self.client.post("http://34.16.87.214.nip.io/agronomy", json=student)
            print(f"Enviando datos: {student}")