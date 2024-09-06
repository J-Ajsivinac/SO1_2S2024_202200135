mod process;
mod utils;
mod analyzer;
mod request;
mod init;

use utils::{read_proc_file, parser_proc_to_struct};
use process::{SystemInfo};
use analyzer::{analyzer};
use std::env;
use init::{start_module,start_cronjob,start_logs_server,get_logs_id,stop_cronjob};
use tokio;
use crate::request::gen_graph; 

// use std::env;
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;
use ctrlc;
use tokio::runtime::Runtime;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let cwd = env::current_dir()?;
    let parent = cwd.parent().unwrap().display().to_string();

    let path_cronjob = format!("{}/script/cronjob.sh",parent);
    let path_logs = format!("{}/python_server/docker-compose.yml",&parent);

    start_cronjob(&path_cronjob);

    println!("Iniciando servicio de logs");
    start_logs_server(&path_logs);
    std::thread::sleep(std::time::Duration::from_secs(10));
    let id_container_logs = get_logs_id("python_server_python_service");
    println!("ID del contenedor de logs: {}",&id_container_logs);
    
    start_module();

    let running = Arc::new(AtomicBool::new(true));
    let r = running.clone();

    // Crear un runtime de Tokio para ejecutar código asíncrono en un contexto síncrono
    let rt = Runtime::new().unwrap();

    ctrlc::set_handler(move || {
        println!("Ctrl+C detected, generating graph...");
        

        // Ejecutar la función asíncrona en el runtime de Tokio
        rt.block_on(async {
            if let Err(e) = gen_graph("graph").await {
                println!("Error generating graph: {}", e);
            }
        });
        stop_cronjob();
        r.store(false, Ordering::SeqCst);
        std::thread::sleep(std::time::Duration::from_secs(2));

    }).expect("Error setting Ctrl-C handler");
    // start_module();
    // start_cronjob(&path);

    // Bucle principal
    while running.load(Ordering::SeqCst) {
        let system_info: Result<SystemInfo, _>;
        let json_str = read_proc_file("sysinfo_202200135").unwrap();
        system_info = parser_proc_to_struct(&json_str);

        match system_info {
            Ok(info) => {
                analyzer(info,&id_container_logs).await?;
            }
            Err(e) => {
                println!("Error: {}", e);
            }
        }
        std::thread::sleep(std::time::Duration::from_secs(20));
    }

    println!("Exiting main loop.");
    Ok(())
}