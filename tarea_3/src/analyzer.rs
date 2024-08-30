use crate::process::{SystemInfo,Process};
// use crate process::{};

pub fn analyzer( system_info:  SystemInfo) {
    let mut processes_list: Vec<Process> = system_info.processes;
    
    processes_list.sort();

    // filtrar contenedores de alto y bajo consumo (alto consumo > 1%)

    let mut high_usage_containers: Vec<Process> = Vec::new();
    let mut low_usage_containers: Vec<Process> = Vec::new();

    for process in processes_list {
        if process.cpu_usage > 1.0 || process.memory_usage > 2.0 {
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
    

    // TODO: ENVIAR LOGS AL CONTENEDOR REGISTRO

    // Hacemos un print de los contenedores que matamos.
    // println!("Contenedores matados");
    // for process in log_proc_list {
    //     println!("PID: {}, Name: {}, Container ID: {}, Memory Usage: {}, CPU Usage: {} ", process.pid, process.name, process.container_id,  process.memory_usage, process.cpu_usage);
    // }

    // println!("------------------------------");

    
}