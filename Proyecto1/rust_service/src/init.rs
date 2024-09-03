pub fn start_module() -> std::process::Output {
    let output = std::process::Command::new("sudo")
        .arg("insmod")
        .arg("../module/sysinfo.ko")
        .output()
        .expect("failed to execute module");
    println!("Module started");
    output
}
