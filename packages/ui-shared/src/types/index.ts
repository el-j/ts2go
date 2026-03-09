// Centralized type definitions for TS2Go Shared UI

// ============================================================================
// Platform Types
// ============================================================================

export type Platform = 'tauri' | 'web'

// ============================================================================
// Project Types
// ============================================================================

/** Project as stored in backend API (SaaS) */
export interface ApiProject {
  id: string
  name: string
  description?: string
  user_id: string
  created_at: string
  updated_at: string
}

/** Project as stored locally (Desktop) */
export interface LocalProject {
  id: string
  name: string
  path: string
  lastModified: Date
  lastOpened?: string
  isPinned?: boolean
  accessCount?: number
}

export interface CreateProjectRequest {
  name: string
  description?: string
}

export interface UpdateProjectRequest {
  name?: string
  description?: string
}

// ============================================================================
// Authentication Types
// ============================================================================

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name?: string
}

export interface User {
  id: string
  email: string
  name?: string
}

export interface AuthResponse {
  token: string
  user: User
}

// ============================================================================
// Transpilation Types
// ============================================================================

export interface TranspileResult {
  success: boolean
  output_dir?: string
  files_transpiled?: number
  files_formatted?: number
  format_warnings?: string[]
  message?: string
  error?: string
  duration?: number
  goCode?: string
}

export interface TranspilationState {
  projectPath: string
  outputDir: string
  filesTranspiled: number
  timestamp: string
  success: boolean
}

export type TranspilationStatus = 'pending' | 'transpiling' | 'success' | 'error' | 'warning'

export interface TranspilationInfo {
  status: TranspilationStatus
  message?: string
  goCode?: string
  error?: string
  timestamp: number
  filesTranspiled?: number
  outputDir?: string
}

export interface TranspilationProgress {
  total: number
  current: number
  currentFile: string
  status: TranspilationStatus
}

// ============================================================================
// Build Types
// ============================================================================

export interface BuildResult {
  success: boolean
  output?: string
  error?: string
  duration?: number
  exitCode?: number
}

export interface BuildRecord {
  id: string
  timestamp: string
  success: boolean
  duration: number
  filesTranspiled: number
  outputPath: string
  errors?: string[]
}

// ============================================================================
// Testing Types
// ============================================================================

export interface TestCase {
  name: string
  input: string
  expectedOutput?: string
  description?: string
}

export interface TestResult {
  name: string
  passed: boolean
  error?: string
  output?: string
  duration?: number
}

export interface TestResults {
  total: number
  passed: number
  failed: number
  skipped: number
  results: TestResult[]
}

// ============================================================================
// File & Editor Types
// ============================================================================

export interface EditorTab {
  id: string
  filePath: string
  fileName: string
  content: string
  language: string
  isDirty: boolean
}

export interface FileNode {
  name: string
  path: string
  type: 'file' | 'directory'
  children?: FileNode[]
  size?: number
  modified?: Date
}

export interface OpenFile {
  path: string
  name: string
  content: string
  isDirty: boolean
}

export interface SaveOptions {
  askForLocation?: boolean
  defaultPath?: string
  filters?: Array<{ name: string; extensions: string[] }>
}

// ============================================================================
// Diagnostic Types
// ============================================================================

export type DiagnosticSeverity = 'error' | 'warning' | 'info'

export interface DiagnosticMessage {
  severity: DiagnosticSeverity
  message: string
  line?: number
  column?: number
  source?: string
  code?: string
}

export interface EditorDiagnostics {
  errors: DiagnosticMessage[]
  warnings: DiagnosticMessage[]
  infos: DiagnosticMessage[]
}

// ============================================================================
// Logging Types
// ============================================================================

export interface LogEntry {
  timestamp: string
  level: 'debug' | 'info' | 'warn' | 'error'
  message: string
  context?: string
  metadata?: Record<string, any>
}

// ============================================================================
// Settings Types
// ============================================================================

export interface AppSettings {
  theme: 'light' | 'dark' | 'system'
  autoSave: boolean
  autoSaveDelay: number
  formatOnSave: boolean
  showLineNumbers: boolean
  fontSize: number
  tabSize: number
  backupEnabled: boolean
  backupLocation: string
  defaultOutputDir: string
  goPath?: string
  modulePath?: string
  excludePatterns: string[]
  includePatterns: string[]
}

// ============================================================================
// Artifact Types
// ============================================================================

export interface Artifact {
  id: string
  name: string
  path: string
  type: 'file' | 'directory'
  size: number
  created: Date
  modified: Date
}

// ============================================================================
// WebSocket Types
// ============================================================================

export interface WebSocketMessage {
  type: string
  data: any
}

// ============================================================================
// Example Types
// ============================================================================

export interface Example {
  id: string
  title: string
  description: string
  typescript: string
  go?: string
  tags?: string[]
  difficulty?: 'beginner' | 'intermediate' | 'advanced'
}

// ============================================================================
// Keyboard Shortcut Types
// ============================================================================

export interface Shortcut {
  key: string
  action: () => void
  description?: string
  ctrl?: boolean
  shift?: boolean
  alt?: boolean
  meta?: boolean
}

// ============================================================================
// Backend State Types
// ============================================================================

export interface BackendStateOptions {
  debounce?: number
  validate?: (data: any) => boolean
}
