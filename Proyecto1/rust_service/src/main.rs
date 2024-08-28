mod process;
mod utils;

use utils::{read_proc_file, parser_proc_to_struct};
use process::{SystemInfo};

fn main() {
    let system_info: Result<SystemInfo, _>;
    let json_str = read_proc_file("sysinfo").unwrap();
    system_info = parser_proc_to_struct(&json_str);

    match system_info{
        Ok(info) => {
            for process in info.processes{
                println!("PID: {}, Name: {}, CPU Usage: {}, Memory Usage: {}, container_id: {}", process.pid, process.name, process.cpu_usage, process.memory_usage, process.get_container_id());
            }
        },
        Err(e) => {
            println!("Error: {}", e);
        }
    }
}
