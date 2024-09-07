use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct SystemInfo {
    #[serde(rename = "Memory")]
    pub memory: Memory,
    #[serde(rename = "Processes")]
    pub processes: Vec<Process>,
}

#[derive(Debug,Serialize, Deserialize, PartialEq)]
pub struct Process {
    #[serde(rename = "pid")]
    pub pid: i32,
    #[serde(rename = "name")]
    pub name: String,
    #[serde(rename = "cmdLine")]
    pub cmd_line: String,
    #[serde(rename = "cpuUsage")]
    pub cpu_usage: f32,
    #[serde(rename = "memoryUsage")]
    pub memory_usage: f32,
    #[serde(rename = "vsz")]
    pub vsz: i32,
    #[serde(rename = "rss")]
    pub rss: i32,
}

#[derive(Debug, Serialize, Clone)]
pub struct LogProcess{
    pub pid: i32,
    pub container_id: String,
    pub name: String,
    pub vsz: i32,
    pub rss: i32,
    pub memory_usage: f32,
    pub cpu_usage: f32,
    pub action: String,
    pub timestamp: String,
}

#[derive(Deserialize,Serialize, Debug)]
pub struct Memory{
    pub total_ram: i32,
    pub free_ram: i32,
    pub used_ram: i32,
}

impl Process {
    pub fn get_container_id(&self) -> &str {
        let parts: Vec<&str> = self.cmd_line.split_whitespace().collect();
        for (i, part) in parts.iter().enumerate() {
            if *part == "-id" {
                if let Some(id) = parts.get(i + 1) {
                    return id;
                }
            }
        }
        "N/A"
    }
}
// Implementación de la interfaz Eq para el struct Process
// Se implementa para poder comparar dos instancias de Process
impl Eq for Process {}

impl Ord for Process {
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        self.cpu_usage
            .partial_cmp(&other.cpu_usage)
            .unwrap_or(std::cmp::Ordering::Equal)
            .then_with( || {
                self.memory_usage
                    .partial_cmp(&other.memory_usage)
                    .unwrap_or(std::cmp::Ordering::Equal)
                    .then_with(|| {
                        self.vsz
                            .partial_cmp(&other.vsz)
                            .unwrap_or(std::cmp::Ordering::Equal)
                            .then_with(|| {
                                self.rss
                                    .partial_cmp(&other.rss)
                                    .unwrap_or(std::cmp::Ordering::Equal)
                            })
                    })
            })
    }
}

impl PartialOrd for Process {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        Some(self.cmp(other))
    }
}