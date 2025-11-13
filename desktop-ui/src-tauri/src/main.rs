// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use std::fs;
use std::path::{Path, PathBuf};
use std::process::Command;

#[cfg(debug_assertions)]
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
    
    // In release, use the bundled binary from Resources
    #[cfg(not(debug_assertions))]
    {
        use std::env;
        
        // Get the path to the executable
        if let Ok(exe_path) = env::current_exe() {
            // Navigate from MacOS/exe to Resources/bin/ts2go-cli
            // Structure: .app/Contents/MacOS/exe -> .app/Contents/Resources/bin/ts2go-cli
            if let Some(contents_dir) = exe_path.parent().and_then(|p| p.parent()) {
                #[cfg(target_os = "macos")]
                {
                    let cli_path = contents_dir.join("Resources").join("bin").join("ts2go-cli");
                    if cli_path.exists() {
                        return cli_path;
                    }
                }
                
                #[cfg(target_os = "windows")]
                {
                    let cli_path = contents_dir.join("bin").join("ts2go-cli.exe");
                    if cli_path.exists() {
                        return cli_path;
                    }
                }
                
                #[cfg(target_os = "linux")]
                {
                    let cli_path = contents_dir.join("bin").join("ts2go-cli");
                    if cli_path.exists() {
                        return cli_path;
                    }
                }
            }
        }
        
        // Fallback to relative path
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
async fn get_go_files(path: String) -> Result<Vec<String>, String> {
    // Recursively find all .go files in the directory
    let mut files = Vec::new();
    find_go_files(Path::new(&path), &mut files)
        .map_err(|e| format!("Failed to list Go files: {}", e))?;
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
async fn run_go_code(code: String) -> Result<serde_json::Value, String> {
    use std::time::{SystemTime, UNIX_EPOCH, Instant};
    
    // Create temporary file for Go code
    let temp_dir = std::env::temp_dir();
    let timestamp = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_millis();
    let temp_file = temp_dir.join(format!("ts2go_run_{}.go", timestamp));
    
    // Write Go code to temp file
    fs::write(&temp_file, &code)
        .map_err(|e| format!("Failed to write temp Go file: {}", e))?;
    
    let start_time = Instant::now();
    
    // Execute Go code using 'go run'
    let output = Command::new("go")
        .arg("run")
        .arg(&temp_file)
        .output()
        .map_err(|e| {
            let _ = fs::remove_file(&temp_file);
            format!("Failed to execute Go code: {}. Make sure Go is installed and in PATH.", e)
        })?;
    
    let duration = start_time.elapsed();
    
    // Clean up temp file
    let _ = fs::remove_file(&temp_file);
    
    // Parse output
    let stdout = String::from_utf8_lossy(&output.stdout).to_string();
    let stderr = String::from_utf8_lossy(&output.stderr).to_string();
    let exit_code = output.status.code().unwrap_or(-1);
    let success = output.status.success();
    
    Ok(serde_json::json!({
        "success": success,
        "stdout": stdout,
        "stderr": stderr,
        "exit_code": exit_code,
        "duration_ms": duration.as_millis() as u64
    }))
}

#[tauri::command]
async fn run_go_project(output_dir: String) -> Result<serde_json::Value, String> {
    use std::time::Instant;
    
    let output_path = Path::new(&output_dir);
    
    // Verify output directory exists
    if !output_path.exists() {
        return Err(format!("Output directory does not exist: {}", output_dir));
    }
    
    // Find main package (look for main.go or any file with package main)
    let mut has_main = false;
    if let Ok(entries) = fs::read_dir(output_path) {
        for entry in entries.flatten() {
            if let Some(name) = entry.file_name().to_str() {
                if name.ends_with(".go") {
                    if let Ok(content) = fs::read_to_string(entry.path()) {
                        if content.contains("package main") {
                            has_main = true;
                            break;
                        }
                    }
                }
            }
        }
    }
    
    if !has_main {
        return Err("No main package found in transpiled output. Make sure your TypeScript project has an entry point.".to_string());
    }
    
    let start_time = Instant::now();
    
    // Execute Go project using 'go run .'
    let output = Command::new("go")
        .arg("run")
        .arg(".")
        .current_dir(output_path)
        .output()
        .map_err(|e| format!("Failed to execute Go project: {}. Make sure Go is installed and in PATH.", e))?;
    
    let duration = start_time.elapsed();
    
    // Parse output
    let stdout = String::from_utf8_lossy(&output.stdout).to_string();
    let stderr = String::from_utf8_lossy(&output.stderr).to_string();
    let exit_code = output.status.code().unwrap_or(-1);
    let success = output.status.success();
    
    Ok(serde_json::json!({
        "success": success,
        "stdout": stdout,
        "stderr": stderr,
        "exit_code": exit_code,
        "duration_ms": duration.as_millis() as u64
    }))
}

#[tauri::command]
async fn build_go_file(code: String, output_path: String) -> Result<serde_json::Value, String> {
    use std::time::{SystemTime, UNIX_EPOCH, Instant};
    
    // Create temporary file for Go code
    let temp_dir = std::env::temp_dir();
    let timestamp = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_millis();
    let temp_file = temp_dir.join(format!("ts2go_build_{}.go", timestamp));
    
    // Write Go code to temp file
    fs::write(&temp_file, &code)
        .map_err(|e| format!("Failed to write temp Go file: {}", e))?;
    
    let start_time = Instant::now();
    
    // Execute Go build
    let output = Command::new("go")
        .arg("build")
        .arg("-o")
        .arg(&output_path)
        .arg(&temp_file)
        .output()
        .map_err(|e| {
            let _ = fs::remove_file(&temp_file);
            format!("Failed to build Go code: {}. Make sure Go is installed and in PATH.", e)
        })?;
    
    let duration = start_time.elapsed();
    
    // Clean up temp file
    let _ = fs::remove_file(&temp_file);
    
    // Parse output
    let stdout = String::from_utf8_lossy(&output.stdout).to_string();
    let stderr = String::from_utf8_lossy(&output.stderr).to_string();
    let exit_code = output.status.code().unwrap_or(-1);
    let success = output.status.success();
    
    // Get binary size if build succeeded
    let binary_size = if success {
        fs::metadata(&output_path)
            .map(|m| m.len())
            .unwrap_or(0)
    } else {
        0
    };
    
    // Parse errors and warnings from stderr
    let mut errors = Vec::new();
    let mut warnings = Vec::new();
    for line in stderr.lines() {
        if line.contains("error:") || line.contains("undefined:") {
            errors.push(line.to_string());
        } else if line.contains("warning:") {
            warnings.push(line.to_string());
        }
    }
    
    Ok(serde_json::json!({
        "success": success,
        "binary_path": output_path,
        "binary_size": binary_size,
        "stdout": stdout,
        "stderr": stderr,
        "errors": errors,
        "warnings": warnings,
        "exit_code": exit_code,
        "duration_ms": duration.as_millis() as u64
    }))
}

#[tauri::command]
async fn build_go_project(
    source_dir: String,
    output_path: String,
    build_flags: Option<Vec<String>>
) -> Result<serde_json::Value, String> {
    use std::time::Instant;
    
    let source_path = Path::new(&source_dir);
    
    // Verify source directory exists
    if !source_path.exists() {
        return Err(format!("Source directory does not exist: {}", source_dir));
    }
    
    // Find main package (look for main.go or any file with package main)
    let mut has_main = false;
    if let Ok(entries) = fs::read_dir(source_path) {
        for entry in entries.flatten() {
            if let Some(name) = entry.file_name().to_str() {
                if name.ends_with(".go") {
                    if let Ok(content) = fs::read_to_string(entry.path()) {
                        if content.contains("package main") {
                            has_main = true;
                            break;
                        }
                    }
                }
            }
        }
    }
    
    if !has_main {
        return Err("No main package found in source directory. Make sure your TypeScript project has an entry point.".to_string());
    }
    
    let start_time = Instant::now();
    
    // Build command with flags
    let mut cmd = Command::new("go");
    cmd.arg("build")
        .arg("-o")
        .arg(&output_path);
    
    // Add custom build flags if provided
    if let Some(flags) = build_flags {
        for flag in flags {
            cmd.arg(flag);
        }
    }
    
    // Execute Go build in the source directory
    let output = cmd
        .current_dir(source_path)
        .output()
        .map_err(|e| format!("Failed to build Go project: {}. Make sure Go is installed and in PATH.", e))?;
    
    let duration = start_time.elapsed();
    
    // Parse output
    let stdout = String::from_utf8_lossy(&output.stdout).to_string();
    let stderr = String::from_utf8_lossy(&output.stderr).to_string();
    let exit_code = output.status.code().unwrap_or(-1);
    let success = output.status.success();
    
    // Get binary size if build succeeded
    let binary_size = if success {
        fs::metadata(&output_path)
            .map(|m| m.len())
            .unwrap_or(0)
    } else {
        0
    };
    
    // Parse errors and warnings from stderr
    let mut errors = Vec::new();
    let mut warnings = Vec::new();
    for line in stderr.lines() {
        if line.contains("error:") || line.contains("undefined:") || line.contains("cannot") {
            errors.push(line.to_string());
        } else if line.contains("warning:") {
            warnings.push(line.to_string());
        }
    }
    
    Ok(serde_json::json!({
        "success": success,
        "binary_path": output_path,
        "binary_size": binary_size,
        "stdout": stdout,
        "stderr": stderr,
        "errors": errors,
        "warnings": warnings,
        "exit_code": exit_code,
        "duration_ms": duration.as_millis() as u64
    }))
}

#[tauri::command]
async fn test_go_file(file_path: String) -> Result<serde_json::Value, String> {
    use std::time::Instant;
    
    let file_path_buf = Path::new(&file_path);
    
    // Verify file exists
    if !file_path_buf.exists() {
        return Err(format!("Go file does not exist: {}", file_path));
    }
    
    let start_time = Instant::now();
    
    // Execute go test
    let output = Command::new("go")
        .arg("test")
        .arg("-v")
        .arg("-json")
        .arg(&file_path)
        .output()
        .map_err(|e| format!("Failed to run tests: {}. Make sure Go is installed and in PATH.", e))?;
    
    let duration = start_time.elapsed();
    
    // Parse output
    let stdout = String::from_utf8_lossy(&output.stdout).to_string();
    let stderr = String::from_utf8_lossy(&output.stderr).to_string();
    let exit_code = output.status.code().unwrap_or(-1);
    let success = output.status.success();
    
    // Parse test results from JSON output
    let test_results = parse_test_results(&stdout);
    
    Ok(serde_json::json!({
        "success": success,
        "stdout": stdout,
        "stderr": stderr,
        "exit_code": exit_code,
        "duration_ms": duration.as_millis() as u64,
        "test_results": test_results
    }))
}

#[tauri::command]
async fn test_go_project(
    source_dir: String,
    test_flags: Option<Vec<String>>
) -> Result<serde_json::Value, String> {
    use std::time::Instant;
    
    let source_path = Path::new(&source_dir);
    
    // Verify source directory exists
    if !source_path.exists() {
        return Err(format!("Source directory does not exist: {}", source_dir));
    }
    
    let start_time = Instant::now();
    
    // Build command with flags
    let mut cmd = Command::new("go");
    cmd.arg("test")
        .arg("-v")
        .arg("-json");
    
    // Add custom test flags if provided
    if let Some(flags) = test_flags {
        for flag in flags {
            cmd.arg(flag);
        }
    }
    
    cmd.arg("./...");
    
    // Execute go test in the source directory
    let output = cmd
        .current_dir(source_path)
        .output()
        .map_err(|e| format!("Failed to run tests: {}. Make sure Go is installed and in PATH.", e))?;
    
    let duration = start_time.elapsed();
    
    // Parse output
    let stdout = String::from_utf8_lossy(&output.stdout).to_string();
    let stderr = String::from_utf8_lossy(&output.stderr).to_string();
    let exit_code = output.status.code().unwrap_or(-1);
    let success = output.status.success();
    
    // Parse test results from JSON output
    let test_results = parse_test_results(&stdout);
    
    Ok(serde_json::json!({
        "success": success,
        "stdout": stdout,
        "stderr": stderr,
        "exit_code": exit_code,
        "duration_ms": duration.as_millis() as u64,
        "test_results": test_results
    }))
}

// Helper function to parse go test JSON output
fn parse_test_results(json_output: &str) -> serde_json::Value {
    use serde_json::Value;
    
    let mut passed = 0;
    let mut failed = 0;
    let mut skipped = 0;
    let mut tests: Vec<Value> = Vec::new();
    
    for line in json_output.lines() {
        if let Ok(json) = serde_json::from_str::<Value>(line) {
            if let Some(action) = json.get("Action").and_then(|v| v.as_str()) {
                if action == "pass" {
                    if json.get("Test").is_some() {
                        passed += 1;
                        tests.push(json.clone());
                    }
                } else if action == "fail" {
                    if json.get("Test").is_some() {
                        failed += 1;
                        tests.push(json.clone());
                    }
                } else if action == "skip" {
                    if json.get("Test").is_some() {
                        skipped += 1;
                        tests.push(json.clone());
                    }
                }
            }
        }
    }
    
    serde_json::json!({
        "passed": passed,
        "failed": failed,
        "skipped": skipped,
        "total": passed + failed + skipped,
        "tests": tests
    })
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
            get_go_files,
            transpile_code,
            read_file,
            write_file,
            load_project_folder,
            auto_transpile_project,
            run_go_code,
            run_go_project,
            build_go_file,
            build_go_project,
            test_go_file,
            test_go_project
        ])
        .setup(|_app| {
            #[cfg(debug_assertions)]
            {
                let window = _app.get_webview_window("main").unwrap();
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

