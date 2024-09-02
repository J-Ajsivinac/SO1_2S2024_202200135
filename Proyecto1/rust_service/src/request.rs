use reqwest::Client;
use crate::process::{LogProcess};
use std::error::Error;

const API_URL: &str = "http://localhost:8000";

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
        let response = res.text().await?;
        println!("Procesos enviados correctamente {}", response);
    } else {
        println!("Error al enviar procesos {} ",res.status());
    }
    Ok(())
}
