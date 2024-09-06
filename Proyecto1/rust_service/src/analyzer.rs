use crate::process::{SystemInfo,Process,LogProcess};
use crate::request::send_process; 
use std::error::Error;
use chrono::{DateTime, Utc, Local};
use std::process::{Stdio};

pub async fn analyzer(system_info: SystemInfo, id_logs:&str) -> Result<(), Box<dyn Error>> {
    let mut processes_list: Vec<Process> = system_info.processes;
    
    processes_list.sort();
    let mut log_proc_list: Vec<LogProcess> = Vec::new();
    // filtrar contenedores de alto y bajo consumo (alto consumo > 1%)
    let now_utc: DateTime<Utc> = Utc::now();
    let now_gt = now_utc.with_timezone(&Local::now().timezone());
    let formatted_date = now_gt.to_rfc3339();
    
    let (highest_list, lowest_list): (Vec<Process>, Vec<Process>) = processes_list
    .into_iter()
    .partition(|process| process.cpu_usage > 0.6 || process.memory_usage > 2.0);

    println!("Contenedores de alto consumo");
    for process in &highest_list {
        println!("container_id: {}, CPU Usage: {}, Memory Usage: {}", process.get_container_id(), process.cpu_usage, process.memory_usage);
    }

    println!("------------------------------");
    println!("Contenedores de bajo consumo");
    for process in &lowest_list {
        println!("container_id: {}, CPU Usage: {}, Memory Usage: {}", process.get_container_id(), process.cpu_usage, process.memory_usage);
    }
    
    if lowest_list.len() > 3 {
        for process in lowest_list.iter().skip(3) {
            let log_process = LogProcess {
                pid: process.pid,
                container_id: process.get_container_id().to_string(),
                name: process.name.clone(),
                vsz: process.vsz,
                rss: process.rss,
                memory_usage: process.memory_usage,
                cpu_usage: process.cpu_usage,
                action: "stop".to_string(),
                timestamp: formatted_date.to_string()
            };
    
            log_proc_list.push(log_process.clone());

            if !process.get_container_id().to_string().starts_with(&id_logs) {
                // Matamos el contenedor.
                let _output = stop_container(&process.get_container_id());
            }

        }
    } 


    if highest_list.len() > 2 {
        // Iteramos sobre los procesos en la lista de alto consumo.
        for process in highest_list.iter().take(highest_list.len() - 2) {
            let log_process = LogProcess {
                pid: process.pid,
                container_id: process.get_container_id().to_string(),
                name: process.name.clone(),
                vsz: process.vsz,
                rss: process.rss,
                memory_usage: process.memory_usage,
                cpu_usage: process.cpu_usage,
                action: "stop".to_string(),
                timestamp: formatted_date.to_string()
            };
    
            log_proc_list.push(log_process.clone());

            if !log_process.container_id.starts_with(&id_logs) {
                // Matamos el contenedor.
                let _output = stop_container(&process.get_container_id());
            }
            // Matamos el contenedor.
            // let _output = kill_container(&process.get_container_id());

        }
    }

    println!("Contenedores matados");
    for process in &log_proc_list {
        println!("PID: {}, Name: {}, Container ID: {}, Memory Usage: {}, CPU Usage: {} ", process.pid, process.name, process.container_id,  process.memory_usage, process.cpu_usage);
    }

    println!("------------------------------");

    let end_url: &str = "logs";
    println!("Enviando procesos al servidor");
    send_process(log_proc_list,end_url).await?;
    Ok(())  
}

pub fn stop_container(id: &str) -> std::process::Child {
    println!("Matando contenedor {}", id);
    let output = std::process::Command::new("docker")
        .arg("stop")
        .arg(id)
        .stdin(Stdio::null())
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .spawn()
        .expect("failed to execute process");
    println!("Contenendor {} matado", id);
    return output;
}