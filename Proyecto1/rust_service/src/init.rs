use std::process::{Stdio};

pub fn start_module() -> std::process::Output {
    let output = std::process::Command::new("sudo")
        .arg("insmod")
        .arg("../module/sysinfo.ko")
        .output()
        .expect("failed to execute module");
    println!("Module started");
    output
}

pub fn start_cronjob(path: &str) -> std::process::Child {
    let output = std::process::Command::new("sh")
        .arg(path)
        .stdin(Stdio::null())
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .spawn()
        .expect("failed to execute cronjob");
    output
}

pub fn stop_cronjob() -> std::process::Output {
    let output = std::process::Command::new("crontab")
        .arg("-r")
        .output()
        .expect("failed to execute cronjob");
    println!("Cronjob stopped");
    output
}

pub fn start_logs_server(path: &str) -> std::process::Child {
    println!("{}",&path.to_string());
    let output = std::process::Command::new("docker-compose")
        .arg("-f")
        .arg(path)
        .arg("up")
        .arg("-d")
        .stdin(Stdio::null())
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .spawn()
        .expect("failed to execute module");

    output
}

pub fn get_logs_id(service_name: &str) -> String {

    let output = std::process::Command::new("docker")
        .arg("ps")
        .arg("--format")
        .arg("{{.ID}}")
        .arg("--filter")
        .arg(format!("ancestor={}", service_name.to_string()))
        .output()
        .expect("Failed to execute command");

        let output_str = std::str::from_utf8(&output.stdout).expect("Failed to convert output to string");

        // Limpia los espacios en blanco y guarda el resultado
        let container_id = output_str.trim();
        println!("Container ID: {}", &container_id.to_string());
    
        // Imprime el ID del contenedor (o úsalo según sea necesario)
        return container_id.to_string();
    }