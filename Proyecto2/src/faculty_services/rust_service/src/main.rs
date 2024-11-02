use studentgrpc::student_client::StudentClient;
use actix_web::{web, App, HttpServer, HttpResponse, Responder};
use studentgrpc::StudentRequest;
use serde::{Deserialize, Serialize};
use std::time::Duration;

pub mod studentgrpc {
    tonic::include_proto!("student_grpc");
}

#[derive(Deserialize, Serialize, Debug, Clone)]
struct StudentData {
    student: String,
    age: i32,
    faculty: String,
    discipline: i32,
}

// Función para procesar la solicitud gRPC
async fn process_grpc_request(student: StudentData, host: String) -> Result<String, String> {
    let mut client = StudentClient::connect(host)
        .await
        .map_err(|e| format!("Failed to connect to gRPC server: {}", e))?;

    let request = tonic::Request::new(StudentRequest {
        student: student.student,
        age: student.age,
        faculty: student.faculty,
        discipline: student.discipline,
    });

    client.get_student_req(request)
        .await
        .map(|response| format!("Student: {:?}", response))
        .map_err(|e| format!("gRPC call failed: {}", e))
}

async fn handle_student(student: web::Json<StudentData>) -> impl Responder {
    let host = match student.discipline {
        1 => "http://swimming-service:50051",
        2 => "http://athletics-service:50051",
        3 => "http://boxing-service:50051",
        _ => return HttpResponse::BadRequest().body("Invalid discipline"),
    };

    let student_data = student.into_inner();
    let host_string = host.to_string();

    // Usar tokio::spawn para manejar la solicitud de forma asíncrona
    match tokio::time::timeout(
        Duration::from_secs(10),
        process_grpc_request(student_data, host_string)
    ).await {
        Ok(Ok(response)) => {
            println!("RESPONSE={}", response);
            HttpResponse::Ok().json(response)
        },
        Ok(Err(error)) => {
            println!("ERROR={}", error);
            HttpResponse::InternalServerError().body(error)
        },
        Err(_) => {
            println!("ERROR=Timeout");
            HttpResponse::GatewayTimeout().body("Request timed out")
        }
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Starting server at port 8080");
    
    HttpServer::new(|| {
        App::new()
            .wrap(actix_web::middleware::Logger::default())
            .route("/engineering", web::post().to(handle_student))
    })
    .workers(4)
    .keep_alive(Duration::from_secs(30))
    .client_request_timeout(Duration::from_secs(30))
    .bind("0.0.0.0:8080")?
    .run()
    .await
}