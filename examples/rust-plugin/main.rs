use std::io::{self, BufRead};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
struct Request {
    jsonrpc: String,
    id: String,
    method: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    params: Option<serde_json::Value>,
}

#[derive(Serialize)]
struct Response<T> {
    jsonrpc: String,
    id: String,
    result: T,
}

#[derive(Serialize)]
struct HandshakeResult {
    name: String,
    version: String,
    api_version: u32,
    checks_supported: Vec<String>,
}

#[derive(Serialize)]
struct RunCheckResult {
    status: String,
    finding: Finding,
}

#[derive(Serialize)]
struct Finding {
    id: String,
    check_id: String,
    severity: String,
    title: String,
    description: String,
    evidence: serde_json::Value,
    remediation: String,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let stdin = io::stdin();
    for line in stdin.lock().lines() {
        let line = match line {
            Ok(l) => l,
            Err(_) => break,
        };

        let req: Request = match serde_json::from_str(&line) {
            Ok(r) => r,
            Err(e) => {
                eprintln!("Error: {}", e);
                continue;
            }
        };

        match req.method.as_str() {
            "handshake" => {
                let resp = Response {
                    jsonrpc: "2.0".to_string(),
                    id: req.id,
                    result: HandshakeResult {
                        name: "example-rust".to_string(),
                        version: "1.0.0".to_string(),
                        api_version: 1,
                        checks_supported: vec!["rust_check".to_string()],
                    },
                };
                println!("{}", serde_json::to_string(&resp)?);
            }
            "run_check" => {
                let resp = Response {
                    jsonrpc: "2.0".to_string(),
                    id: req.id,
                    result: RunCheckResult {
                        status: "pass".to_string(),
                        finding: Finding {
                            id: "F-RS-001".to_string(),
                            check_id: "rust_check".to_string(),
                            severity: "info".to_string(),
                            title: "Rust Plugin Running".to_string(),
                            description: "Rust-based plugin is communicating correctly.".to_string(),
                            evidence: serde_json::json!({"language": "rust"}),
                            remediation: "None".to_string(),
                        },
                    },
                };
                println!("{}", serde_json::to_string(&resp)?);
            }
            _ => {
                eprintln!("Unknown method: {}", req.method);
            }
        }
    }
    Ok(())
}
