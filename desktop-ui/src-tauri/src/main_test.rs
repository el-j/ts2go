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
