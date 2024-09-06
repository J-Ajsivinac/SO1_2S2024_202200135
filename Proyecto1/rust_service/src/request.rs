use reqwest::Client;
use crate::process::{LogProcess};
use std::error::Error;
use serde::Deserialize;

const API_URL: &str = "http://localhost:8000";

#[derive(Deserialize)]
struct ErrorMessage {
    message: String,
}

pub async fn send_process(process_list: Vec<LogProcess>,url:&str)-> Result<(), Box<dyn Error>>{
    let json_body = serde_json::to_string(&process_list)?;
    let client = Client::new();
    let full_url = format!("{}/{}",API_URL,url);
    // let temp = format!("{}/{}",API_URL,url);
    let res = client.post(full_url)
        .body(json_body)
        .send()
        .await?;

    if res.status().is_success(){
        // let response = res.text().await?;
        println!("Procesos enviados correctamente");
    } else {
        println!("Error al enviar procesos {} ",res.status());
    }
    Ok(())
}

pub async fn gen_graph(url:&str)-> Result<(), Box<dyn Error>>{
    let client = Client::new();
    let full_url = format!("{}/{}",API_URL,url);
    let res = client.get(full_url)
        .send()
        .await?;

    if res.status().is_success(){
        // let response = res.text().await?;
        println!("Grafica generada correctamente");
    } else {
        // {"message":"Error al generar la grafica"}
        let response = res.text().await?;
        let error_message: ErrorMessage = serde_json::from_str(&response)?;
        println!("Error al generar la grafica: {} ",error_message.message);
    }
    Ok(())
}
