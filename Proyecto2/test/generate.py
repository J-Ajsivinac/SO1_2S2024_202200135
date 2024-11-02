import json
import random

def generate_student_data(num_students):
    faculties = ["Ingenieria", "Agronomia"]
    first_names = ["John", "Jane", "Alex", "Emily", "Chris", "Anna"]
    last_names = ["Doe", "Smith", "Johnson", "Brown", "Davis", "Martinez"]
    
    students = []

    for i in range(num_students):
        student = {
            "student": f"{random.choice(first_names)} {random.choice(last_names)}",
            "age": random.randint(20, 30),
            "faculty": random.choice(faculties),
            "discipline": random.randint(1, 3)
        }
        students.append(student)
    
    return students

# Genera el JSON con el número deseado de estudiantes
num_students = 10000  # Cambia este valor al tamaño deseado
students_data = generate_student_data(num_students)

# Guarda el resultado en un archivo JSON
with open("test/students.json", "w") as json_file:
    json.dump(students_data, json_file, indent=4)

print("JSON generado y guardado en 'students.json'")
