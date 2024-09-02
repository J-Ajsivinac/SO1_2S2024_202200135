mod process;
mod utils;
mod analyzer;
mod request;

use utils::{read_proc_file, parser_proc_to_struct};
use process::{SystemInfo};
use analyzer::{analyzer};
use tokio;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
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

    Ok(())
}
