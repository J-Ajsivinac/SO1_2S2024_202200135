use crate::process::{SystemInfo,Process,LogProcess};
// use crate process::{};

pub fn analyzer( system_info:  SystemInfo) {
    let mut processes_list: Vec<Process> = system_info.processes;
    
    processes_list.sort();
    let mut log_proc_list: Vec<LogProcess> = Vec::new();
    // filtrar contenedores de alto y bajo consumo (alto consumo > 1%)

    let mut high_usage_containers: Vec<Process> = Vec::new();
    let mut low_usage_containers: Vec<Process> = Vec::new();

    for process in processes_list {
        if process.cpu_usage > 1.3 || process.memory_usage > 2.0 {
            high_usage_containers.push(process);
        } else {
            low_usage_containers.push(process);
        }
    }

    println!("Contenedores de alto consumo");
    for process in high_usage_containers {
        println!("container_id: {}, CPU Usage: {}, Memory Usage: {}", process.get_container_id(), process.cpu_usage, process.memory_usage);
    }

    println!("------------------------------");
    println!("Contenedores de bajo consumo");
    for process in low_usage_containers {
        println!("container_id: {}, CPU Usage: {}, Memory Usage: {}", process.get_container_id(), process.cpu_usage, process.memory_usage);
    }
    
    if low_usage_containers.len > 3{
        for process in low_usage_containers.skip(3){
            let log_process = LogProcess{
                pid: process.pid,
                container_id: process.get_container_id().to_string(),
                name: process.name.to_string(),
                cpu_usage: process.cpu_usage,
                memory_usage: process.memory_usage,
            };
            log_proc_list.push(log_process.clone());
            let _output = stop_container(&process.container_id);
        }
    }

    if high_usage_containers.len() > 2{
        for process in high_usage_containers.iter().take(high_usage_containers.len() - 2){
            let log_process = LogProcess{
                pid: process.pid,
                container_id: process.get_container_id().to_string(),
                name: process.name.to_string(),
                cpu_usage: process.cpu_usage,
                memory_usage: process.memory_usage,
            };
            log_proc_list.push(log_process.clone());

            let _output = stop_container(&process.container_id);
        }
    }

    
}

pub fn stop_container(id: &str) -> std::process::Output {
    let output = std::process::Command::new("docker")
        .arg("stop")
        .arg(id)
        .output()
        .expect("failed to execute process");
    println!("Contenendor {} matado", id);
    return output;
}