use actix_web::{post, web, App, HttpServer, Responder, HttpResponse};
use serde::Deserialize;

#[derive(Deserialize)]
struct Student {
    student : String,
    age: u32,
    faculty: String,
    discipline: u32,
}

#[post("/engineering")]
async fn engineering_service(student: web::Json<Student>) ->  impl Responder {
    println!(
        "Recibido: Facultad = {}, Disciplina = {}, Estudiante = {}, Edad = {}",
        student.faculty, student.discipline, student.student, student.age
    );
    
    HttpResponse::Ok().body("Recibido")
}

#[actix_web::main]
async fn main() -> std::io::Result<()>{
    println!("Servicio de Rust corriendo en el puerto 8081");
    HttpServer::new(|| {
        App::new()
            .service(engineering_service)
    })
    .bind("0.0.0.0:8081")?
    .run()
    .await
}