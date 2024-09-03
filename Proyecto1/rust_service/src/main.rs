mod process;
mod utils;
mod analyzer;
mod request;
mod init;

use utils::{read_proc_file, parser_proc_to_struct};
use process::{SystemInfo};
use analyzer::{analyzer};
use std::env;
use init::{start_module,start_cronjob};
use tokio;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let cwd = env::current_dir()?;
    let parent = cwd.parent().unwrap().display().to_string();
    let path = format!("{}/script/cronjob.sh",parent);
    start_module();
    // println!("Path: {:?}",path);
    start_cronjob(&path);
    std::thread::sleep(std::time::Duration::from_secs(10));
    // println!("-> {:?}",temp);
    loop {
        let system_info: Result<SystemInfo, _>;
        let json_str = read_proc_file("sysinfo_202200135").unwrap();
        system_info = parser_proc_to_struct(&json_str);

        match system_info{
            Ok(info) => {
                analyzer(info).await?;
            },
            Err(e) => {
                println!("Error: {}", e);
            }
        }
        std::thread::sleep(std::time::Duration::from_secs(10));
    }
}
