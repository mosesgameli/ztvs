// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

use std::io::{self, BufRead};
use serde_json::{json, Value};

fn main() {
    let stdin = io::stdin();
    for line in stdin.lock().lines() {
        let line = line.expect("Failed to read line");
        let request: Value = match serde_json::from_str(&line) {
            Ok(v) => v,
            Err(_) => continue,
        };

        let method = request.get("method").and_then(|m| m.as_str()).unwrap_or("");
        let id = request.get("id").unwrap();

        if method == "handshake" {
            let response = json!({
                "jsonrpc": "2.0",
                "id": id,
                "result": {
                    "name": "rust-fixture",
                    "version": "1.0.0",
                    "api_version": 1,
                    "checks_supported": ["rust_test_check"]
                }
            });
            println!("{}", response.to_string());
        } else if method == "run_check" {
            let params = request.get("params").unwrap();
            if params.get("check_id").and_then(|c| c.as_str()) == Some("rust_test_check") {
                let response = json!({
                    "jsonrpc": "2.0",
                    "id": id,
                    "result": {
                        "status": "pass",
                        "finding": {
                            "id": "F-RS-001",
                            "severity": "info",
                            "title": "Rust Fixture Executed",
                            "description": "This is a successful result from the raw Rust polyglot fixture."
                        }
                    }
                });
                println!("{}", response.to_string());
            }
        }
    }
}
