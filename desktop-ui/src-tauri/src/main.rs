// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use tauri::Manager;

// Tauri commands
#[tauri::command]
async fn analyze_project(path: String) -> Result<String, String> {
    // TODO: Implement project analysis
    Ok(format!("Analyzing project at: {}", path))
}

#[tauri::command]
async fn transpile_project(
    path: String,
    output: String,
    optimize: bool,
) -> Result<String, String> {
    // TODO: Implement transpilation
    Ok(format!(
        "Transpiling project from {} to {} (optimize: {})",
        path, output, optimize
    ))
}

#[tauri::command]
async fn get_project_files(path: String) -> Result<Vec<String>, String> {
    // TODO: Implement file listing
    Ok(vec![format!("{}/example.ts", path)])
}

fn main() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .invoke_handler(tauri::generate_handler![
            analyze_project,
            transpile_project,
            get_project_files
        ])
        .setup(|app| {
            #[cfg(debug_assertions)]
            {
                let window = app.get_webview_window("main").unwrap();
                window.open_devtools();
            }
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
