/**
 * Get file icon based on file extension
 */
export function getFileIcon(fileName: string): string {
  const ext = fileName.split('.').pop()?.toLowerCase()
  
  const iconMap: Record<string, string> = {
    // TypeScript/JavaScript
    'ts': '📘',
    'tsx': '⚛️',
    'js': '📙',
    'jsx': '⚛️',
    
    // Go
    'go': '🐹',
    'mod': '📦',
    
    // Markup/Data
    'json': '📋',
    'md': '📝',
    'html': '🌐',
    'xml': '📄',
    'yaml': '⚙️',
    'yml': '⚙️',
    'toml': '⚙️',
    
    // Styles
    'css': '🎨',
    'scss': '🎨',
    'sass': '🎨',
    'less': '🎨',
    
    // Config
    'config': '⚙️',
    'conf': '⚙️',
    'ini': '⚙️',
    'env': '🔒',
    
    // Build/Package
    'lock': '🔒',
    'sum': '✓',
    
    // Vue
    'vue': '💚',
    
    // Images
    'png': '🖼️',
    'jpg': '🖼️',
    'jpeg': '🖼️',
    'gif': '🖼️',
    'svg': '🖼️',
    'ico': '🖼️',
    
    // Documents
    'pdf': '📕',
    'doc': '📘',
    'docx': '📘',
    'txt': '📄',
    
    // Archives
    'zip': '📦',
    'tar': '📦',
    'gz': '📦',
    'rar': '📦',
    
    // Shell/Scripts
    'sh': '🐚',
    'bash': '🐚',
    'zsh': '🐚',
    'fish': '🐚',
    'ps1': '💻'
  }
  
  return iconMap[ext || ''] || '📄'
}

/**
 * Get file type label
 */
export function getFileType(fileName: string): string {
  const ext = fileName.split('.').pop()?.toLowerCase()
  
  const typeMap: Record<string, string> = {
    'ts': 'TypeScript',
    'tsx': 'TypeScript React',
    'js': 'JavaScript',
    'jsx': 'JavaScript React',
    'go': 'Go',
    'json': 'JSON',
    'md': 'Markdown',
    'vue': 'Vue',
    'css': 'CSS',
    'html': 'HTML',
    'yml': 'YAML',
    'yaml': 'YAML'
  }
  
  return typeMap[ext || ''] || 'File'
}

/**
 * Check if file is a TypeScript file
 */
export function isTypeScriptFile(fileName: string): boolean {
  const ext = fileName.split('.').pop()?.toLowerCase()
  return ext === 'ts' || ext === 'tsx'
}

/**
 * Check if file is a Go file
 */
export function isGoFile(fileName: string): boolean {
  const ext = fileName.split('.').pop()?.toLowerCase()
  return ext === 'go'
}

/**
 * Check if file is a config file
 */
export function isConfigFile(fileName: string): boolean {
  const ext = fileName.split('.').pop()?.toLowerCase()
  const configExts = ['json', 'yaml', 'yml', 'toml', 'ini', 'conf', 'config']
  return configExts.includes(ext || '')
}

/**
 * Get file size in human-readable format
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return `${Math.round(bytes / Math.pow(k, i) * 100) / 100} ${sizes[i]}`
}

/**
 * Sanitize filename (remove invalid characters)
 */
export function sanitizeFilename(filename: string): string {
  // Remove invalid characters
  return filename.replace(/[<>:"/\\|?*\x00-\x1f]/g, '_')
}

/**
 * Validate filename
 */
export function isValidFilename(filename: string): boolean {
  if (!filename || filename.trim() === '') return false
  if (filename.length > 255) return false
  if (/[<>:"/\\|?*\x00-\x1f]/.test(filename)) return false
  if (filename === '.' || filename === '..') return false
  
  return true
}

/**
 * Get relative path from project root
 */
export function getRelativePath(fullPath: string, projectRoot: string): string {
  if (fullPath.startsWith(projectRoot)) {
    return fullPath.slice(projectRoot.length + 1)
  }
  return fullPath
}
