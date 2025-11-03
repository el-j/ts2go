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

#[tauri::command]
async fn transpile_code(code: String, filename: String) -> Result<String, String> {
    // TODO: Call ts2go CLI to transpile code
    // For now, return a mock Go translation
    Ok(format!(
        "// Generated Go code from {}\npackage main\n\n// TODO: Implement actual transpilation\n// Original TypeScript:\n{}\n",
        filename,
        code.lines().map(|l| format!("// {}", l)).collect::<Vec<_>>().join("\n")
    ))
}

fn main() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .invoke_handler(tauri::generate_handler![
            analyze_project,
            transpile_project,
            get_project_files,
            transpile_code
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

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_analyze_project() {
        let result = analyze_project("/test/path".to_string()).await;
        assert!(result.is_ok());
        let message = result.unwrap();
        assert!(message.contains("Analyzing project at"));
        assert!(message.contains("/test/path"));
    }

    #[tokio::test]
    async fn test_transpile_project() {
        let result = transpile_project(
            "/test/input".to_string(),
            "/test/output".to_string(),
            true,
        )
        .await;
        assert!(result.is_ok());
        let message = result.unwrap();
        assert!(message.contains("Transpiling project"));
        assert!(message.contains("/test/input"));
        assert!(message.contains("/test/output"));
        assert!(message.contains("optimize: true"));
    }

    #[tokio::test]
    async fn test_get_project_files() {
        let result = get_project_files("/test/path".to_string()).await;
        assert!(result.is_ok());
        let files = result.unwrap();
        assert_eq!(files.len(), 1);
        assert_eq!(files[0], "/test/path/example.ts");
    }
}

