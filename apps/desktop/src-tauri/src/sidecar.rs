use std::path::{Path, PathBuf};
use std::process::{Child, Command, Stdio};
use std::sync::Mutex;

static CHILDREN: Mutex<Vec<Child>> = Mutex::new(Vec::new());

/// Запускает Go API и Rust engine, если заданы пути (для сборки «всё в одном»).
pub fn spawn_sidecars() {
    let api = resolve_bin(
        "FIN_API_BIN",
        &["fin-api.exe", "fin-api-x86_64-pc-windows-msvc.exe"],
    );
    let engine = resolve_bin(
        "FIN_ENGINE_BIN",
        &["fin-engine.exe", "fin-engine-x86_64-pc-windows-msvc.exe"],
    );
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
        .env("API_ADDR", "127.0.0.1:8080")
        .env("ENGINE_HTTP_URL", "http://127.0.0.1:50052")
        .env("ENGINE_HTTP_ADDR", "127.0.0.1:50052")
        .env("API_TOKEN", std::env::var("VITE_API_TOKEN").unwrap_or_else(|_| "dev-token".to_string()))
        .stdin(Stdio::null())
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .spawn()
        .ok()
}

fn resolve_bin(env_name: &str, candidates: &[&str]) -> Option<String> {
    if let Ok(v) = std::env::var(env_name) {
        if Path::new(&v).exists() {
            return Some(v);
        }
    }
    let exe_dir = std::env::current_exe()
        .ok()
        .and_then(|p| p.parent().map(|d| d.to_path_buf()))?;
    for c in candidates {
        let direct = exe_dir.join(c);
        if direct.exists() {
            return Some(direct.to_string_lossy().to_string());
        }
        let bundled = exe_dir.join("binaries").join(c);
        if bundled.exists() {
            return Some(bundled.to_string_lossy().to_string());
        }
    }
    None
}
