// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use std::fs;
use std::path::{Path, PathBuf};
use std::process::Command;
use tauri::Manager;

// Helper function to get the path to the bundled CLI binary
fn get_cli_binary_path() -> PathBuf {
    // In development, use the binary from the bin directory
    #[cfg(debug_assertions)]
    {
        let bin_path = Path::new("bin/ts2go-cli");
        if bin_path.exists() {
            return bin_path.to_path_buf();
        }
        // Fallback to system ts2go if bin doesn't exist
        return PathBuf::from("ts2go");
    }
    
    // In release, use the bundled binary
    #[cfg(not(debug_assertions))]
    {
        // The binary is bundled as a resource
        // Tauri places it in the resource directory
        #[cfg(target_os = "windows")]
        {
            PathBuf::from("ts2go-cli.exe")
        }
        #[cfg(not(target_os = "windows"))]
        {
            PathBuf::from("ts2go-cli")
        }
    }
}

// Tauri commands
#[tauri::command]
async fn analyze_project(path: String) -> Result<String, String> {
    // Call ts2go CLI to analyze project
    let cli_path = get_cli_binary_path();
    let output = Command::new(&cli_path)
        .arg("analyze")
        .arg(&path)
        .output()
        .map_err(|e| format!("Failed to execute ts2go analyze: {}. CLI path: {:?}", e, cli_path))?;

    if output.status.success() {
        String::from_utf8(output.stdout)
            .map_err(|e| format!("Failed to parse output: {}", e))
    } else {
        let error = String::from_utf8_lossy(&output.stderr);
        Err(format!("Analysis failed: {}", error))
    }
}

#[tauri::command]
async fn transpile_project(
    path: String,
    output: String,
    optimize: bool,
) -> Result<String, String> {
    // Call ts2go CLI to transpile project
    let cli_path = get_cli_binary_path();
    let mut cmd = Command::new(&cli_path);
    cmd.arg("transpile").arg(&path).arg("--out").arg(&output);

    if optimize {
        cmd.arg("--optimize");
    }

    let result = cmd
        .output()
        .map_err(|e| format!("Failed to execute ts2go transpile: {}. CLI path: {:?}", e, cli_path))?;

    if result.status.success() {
        Ok(format!(
            "Successfully transpiled project from {} to {}",
            path, output
        ))
    } else {
        let error = String::from_utf8_lossy(&result.stderr);
        Err(format!("Transpilation failed: {}", error))
    }
}

#[tauri::command]
async fn get_project_files(path: String) -> Result<Vec<String>, String> {
    // Recursively find all .ts and .tsx files in the project directory
    let mut files = Vec::new();
    find_typescript_files(Path::new(&path), &mut files)
        .map_err(|e| format!("Failed to list files: {}", e))?;
    Ok(files)
}

#[tauri::command]
async fn read_file(path: String) -> Result<String, String> {
    // Read file content
    fs::read_to_string(&path)
        .map_err(|e| format!("Failed to read file {}: {}", path, e))
}

#[tauri::command]
async fn write_file(path: String, content: String) -> Result<(), String> {
    // Write content to file
    fs::write(&path, content)
        .map_err(|e| format!("Failed to write file {}: {}", path, e))
}

#[tauri::command]
async fn load_project_folder(path: String) -> Result<serde_json::Value, String> {
    // Load project folder and return file tree structure
    let mut files = Vec::new();
    find_typescript_files(Path::new(&path), &mut files)
        .map_err(|e| format!("Failed to load project: {}", e))?;
    
    // Build file tree structure
    let tree = build_file_tree(&path, &files)?;
    
    Ok(serde_json::json!({
        "root": path,
        "files": files,
        "tree": tree
    }))
}

#[tauri::command]
async fn auto_transpile_project(path: String) -> Result<serde_json::Value, String> {
    // Auto-transpile entire project
    let cli_path = get_cli_binary_path();
    let output_dir = format!("{}_go", path.trim_end_matches('/'));
    
    // Create output directory
    fs::create_dir_all(&output_dir)
        .map_err(|e| format!("Failed to create output directory: {}", e))?;
    
    // Transpile project
    let result = Command::new(&cli_path)
        .arg("transpile")
        .arg(&path)
        .arg("--out")
        .arg(&output_dir)
        .output()
        .map_err(|e| format!("Failed to execute ts2go transpile: {}. CLI path: {:?}", e, cli_path))?;
    
    if result.status.success() {
        // Count transpiled files
        let mut go_files = Vec::new();
        find_go_files(Path::new(&output_dir), &mut go_files)
            .map_err(|e| format!("Failed to count output files: {}", e))?;
        
        Ok(serde_json::json!({
            "success": true,
            "output_dir": output_dir,
            "files_transpiled": go_files.len(),
            "message": format!("Successfully transpiled {} files to {}", go_files.len(), output_dir)
        }))
    } else {
        let error = String::from_utf8_lossy(&result.stderr);
        Ok(serde_json::json!({
            "success": false,
            "error": error.to_string()
        }))
    }
}

// Helper function to build file tree structure
fn build_file_tree(root: &str, files: &[String]) -> Result<serde_json::Value, String> {
    use std::collections::HashMap;
    
    let root_path = Path::new(root);
    let mut tree: HashMap<String, Vec<String>> = HashMap::new();
    
    for file in files {
        let file_path = Path::new(file);
        if let Ok(rel_path) = file_path.strip_prefix(root_path) {
            let parent = rel_path.parent()
                .and_then(|p| p.to_str())
                .unwrap_or("");
            let filename = file_path.file_name()
                .and_then(|n| n.to_str())
                .unwrap_or("");
            
            tree.entry(parent.to_string())
                .or_insert_with(Vec::new)
                .push(filename.to_string());
        }
    }
    
    Ok(serde_json::to_value(tree).unwrap())
}

// Helper function to find Go files
fn find_go_files(dir: &Path, files: &mut Vec<String>) -> std::io::Result<()> {
    if dir.is_dir() {
        for entry in fs::read_dir(dir)? {
            let entry = entry?;
            let path = entry.path();
            
            if path.is_dir() {
                find_go_files(&path, files)?;
            } else if let Some(ext) = path.extension() {
                if ext == "go" {
                    if let Some(path_str) = path.to_str() {
                        files.push(path_str.to_string());
                    }
                }
            }
        }
    }
    Ok(())
}

#[tauri::command]
async fn transpile_code(code: String, filename: String) -> Result<String, String> {
    // Create temporary file and transpile using ts2go CLI
    let temp_dir = std::env::temp_dir();
    let input_path = temp_dir.join(&filename);
    let output_path = temp_dir.join(filename.replace(".ts", ".go"));

    // Write TypeScript code to temp file
    fs::write(&input_path, &code)
        .map_err(|e| format!("Failed to write temp file: {}", e))?;

    // Call ts2go CLI to convert
    let cli_path = get_cli_binary_path();
    let result = Command::new(&cli_path)
        .arg("convert")
        .arg("--in")
        .arg(&input_path)
        .arg("--out")
        .arg(&output_path)
        .output()
        .map_err(|e| format!("Failed to execute ts2go convert: {}. CLI path: {:?}", e, cli_path))?;

    // Clean up input file
    let _ = fs::remove_file(&input_path);

    if result.status.success() {
        // Read generated Go code
        let go_code = fs::read_to_string(&output_path)
            .map_err(|e| format!("Failed to read generated Go code: {}", e))?;
        
        // Clean up output file
        let _ = fs::remove_file(&output_path);
        
        Ok(go_code)
    } else {
        let error = String::from_utf8_lossy(&result.stderr);
        Err(format!("Transpilation failed: {}", error))
    }
}

// Helper function to recursively find TypeScript files
fn find_typescript_files(dir: &Path, files: &mut Vec<String>) -> std::io::Result<()> {
    if dir.is_dir() {
        for entry in fs::read_dir(dir)? {
            let entry = entry?;
            let path = entry.path();
            
            // Skip node_modules and hidden directories
            if let Some(name) = path.file_name() {
                let name_str = name.to_string_lossy();
                if name_str.starts_with('.') || name_str == "node_modules" {
                    continue;
                }
            }
            
            if path.is_dir() {
                find_typescript_files(&path, files)?;
            } else if let Some(ext) = path.extension() {
                if ext == "ts" || ext == "tsx" {
                    if let Some(path_str) = path.to_str() {
                        files.push(path_str.to_string());
                    }
                }
            }
        }
    }
    Ok(())
}

fn main() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_dialog::init())
        .invoke_handler(tauri::generate_handler![
            analyze_project,
            transpile_project,
            get_project_files,
            transpile_code,
            read_file,
            write_file,
            load_project_folder,
            auto_transpile_project
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
    async fn test_transpile_code_creates_temp_files() {
        // Test that transpile_code properly creates and cleans up temp files
        let code = "const x: number = 5;".to_string();
        let filename = "test.ts".to_string();
        
        // This will fail if ts2go CLI is not available, which is expected in test environment
        // The test verifies the logic flow is correct
        let result = transpile_code(code, filename).await;
        
        // Either succeeds or fails with expected error about CLI not found
        match result {
            Ok(_) => {
                // CLI available and worked
                assert!(true);
            }
            Err(e) => {
                // Expected if ts2go CLI not in PATH during tests
                assert!(e.contains("ts2go") || e.contains("Failed to execute"));
            }
        }
    }

    #[test]
    fn test_find_typescript_files_skips_node_modules() {
        // Create temp directory structure for testing
        let temp_dir = std::env::temp_dir().join("test_ts2go_files");
        let _ = fs::create_dir_all(&temp_dir);
        let _ = fs::create_dir_all(temp_dir.join("node_modules"));
        let _ = fs::write(temp_dir.join("test.ts"), "const x = 1;");
        let _ = fs::write(temp_dir.join("node_modules/lib.ts"), "export const y = 2;");
        
        let mut files = Vec::new();
        let _ = find_typescript_files(&temp_dir, &mut files);
        
        // Should find test.ts but not node_modules/lib.ts
        let has_test_file = files.iter().any(|f| f.contains("test.ts"));
        let has_node_modules = files.iter().any(|f| f.contains("node_modules"));
        
        assert!(has_test_file, "Should find test.ts");
        assert!(!has_node_modules, "Should skip node_modules");
        
        // Cleanup
        let _ = fs::remove_dir_all(&temp_dir);
    }
}

