mod process;
mod utils;
mod analyzer;

use utils::{read_proc_file, parser_proc_to_struct};
use process::{SystemInfo};

use analyzer::{analyzer};

fn main() {
    let system_info: Result<SystemInfo, _>;
    let json_str = read_proc_file("sysinfo_202200135").unwrap();
    system_info = parser_proc_to_struct(&json_str);

    match system_info{
        Ok(info) => {
            analyzer(info);
        },
        Err(e) => {
            println!("Error: {}", e);
        }
    }
}
