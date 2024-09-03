pub fn start_module() -> std::process::Output {
    let output = std::process::Command::new("sudo")
        .arg("insmod")
        .arg("../module/sysinfo.ko")
        .output()
        .expect("failed to execute module");
    println!("Module started");
    output
}

// pub fn start_cronjob() -> std::process::Output {
//     let output = std::process::Command::new("./cronjob.sh")
//         .output()
//         .expect("failed to execute cronjob");
//     println!("Cronjob started");
//     output
// }

// pub fn stop_cronjob() -> std::process::Output {
//     let output = std::process::Command::new("sudo")
//         .arg("crontab")
//         .arg("-r")
//         .output()
//         .expect("failed to execute cronjob");
//     println!("Cronjob stopped");
//     output
// }