use std::path::PathBuf;
use std::process::{Child, Command, Stdio};
use std::sync::Mutex;

static CHILDREN: Mutex<Vec<Child>> = Mutex::new(Vec::new());

/// Запускает Go API и Rust engine, если заданы пути (для сборки «всё в одном»).
pub fn spawn_sidecars() {
    let api = std::env::var("FIN_API_BIN").ok();
    let engine = std::env::var("FIN_ENGINE_BIN").ok();
    if api.is_none() && engine.is_none() {
        return;
    }

    let mut children = CHILDREN.lock().unwrap();

    if let Some(bin) = api {
        if let Some(child) = spawn_bin(&bin, &[]) {
            eprintln!("started API: {bin}");
            children.push(child);
        }
    }

    if let Some(bin) = engine {
        if let Some(child) = spawn_bin(&bin, &[]) {
            eprintln!("started engine: {bin}");
            children.push(child);
        }
    }
}

fn spawn_bin(path: &str, _args: &[&str]) -> Option<Child> {
    let p = PathBuf::from(path);
    if !p.exists() {
        eprintln!("sidecar not found: {path}");
        return None;
    }
    Command::new(p)
        .stdin(Stdio::null())
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .spawn()
        .ok()
}
