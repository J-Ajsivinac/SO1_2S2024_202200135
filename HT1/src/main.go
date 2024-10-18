package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"src/schemas"
)

func agronomyHandler(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "No se pudo leer la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var student schemas.Student
	err = json.Unmarshal(body, &student)
	if err != nil {
		http.Error(w, "No se pudo decodificar el JSON", http.StatusBadRequest)
		return
	}
	fmt.Printf("Recibido: Facultad = %s, Disciplina = %d, Estudiante = %s, Edad = %d", student.Faculty, student.Discipline, student.Student, student.Age)
	fmt.Println()
}
func handleRequests() {
	http.HandleFunc("/agronomy", agronomyHandler)
	http.ListenAndServe(":8080", nil)
}
func main() {
	fmt.Println("Servidor iniciado en http://localhost:8080")
	handleRequests()
	
}
