import json
from locust import HttpUser, task, between

class StudentTraffic(HttpUser):
    wait_time = between(10, 15)

    def on_start(self):
        # Cargar el archivo JSON al iniciar Locust
        with open('students.json', 'r') as file:
            self.students_data = json.load(file)

    @task
    def send_traffic(self):
        # Iterar sobre los datos cargados desde el archivo JSON
        for student in self.students_data:
            if student['faculty'].lower() == 'ingenieria':
                self.client.post("/engineering", json=student)
                print(f"Enviando datos: {student}")
            elif student['faculty'].lower() == 'agronomia':
                self.client.post("/agronomy", json=student)
                print(f"Enviando datos: {student}")


# Ejecutar Locust
# locust -f test/locustfile.py --host=http://localhost:8080
# http://localhost:8089/

