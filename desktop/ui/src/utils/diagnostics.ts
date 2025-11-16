/**
 * Diagnostic severity levels
 */
export type DiagnosticSeverity = 'error' | 'warning' | 'info';

/**
 * Represents a diagnostic message (error, warning, or info)
 */
export interface DiagnosticMessage {
  severity: DiagnosticSeverity;
  message: string;
  line: number;
  column: number;
  code?: string;
  source?: string;
}

/**
 * Editor diagnostics for a specific file
 */
export interface EditorDiagnostics {
  file: string;
  diagnostics: DiagnosticMessage[];
}

/**
 * Get an icon for a diagnostic severity level
 * @param severity - The diagnostic severity
 * @returns An emoji icon representing the severity
 */
export function getSeverityIcon(severity: DiagnosticSeverity): string {
  switch (severity) {
    case 'error':
      return '❌';
    case 'warning':
      return '⚠️';
    case 'info':
      return 'ℹ️';
    default:
      return '•';
  }
}

/**
 * Get a color for a diagnostic severity level
 * @param severity - The diagnostic severity
 * @returns A CSS color string
 */
export function getSeverityColor(severity: DiagnosticSeverity): string {
  switch (severity) {
    case 'error':
      return '#f48771';
    case 'warning':
      return '#cca700';
    case 'info':
      return '#75beff';
    default:
      return '#888';
  }
}

/**
 * Parse diagnostics from compiler output
 * @param output - The compiler output string
 * @returns An array of diagnostic messages
 */
export function parseDiagnostics(output: string): DiagnosticMessage[] {
  const diagnostics: DiagnosticMessage[] = [];
  const lines = output.split('\n');

  for (const line of lines) {
    // Try to match common diagnostic patterns
    // Example: "file.ts(10,5): error TS2304: Cannot find name 'foo'."
    const tsMatch = line.match(/^(.+)\((\d+),(\d+)\):\s+(error|warning|info)\s+(\w+):\s+(.+)$/);
    if (tsMatch) {
      diagnostics.push({
        severity: tsMatch[4] as DiagnosticSeverity,
        message: tsMatch[6],
        line: parseInt(tsMatch[2], 10),
        column: parseInt(tsMatch[3], 10),
        code: tsMatch[5],
      });
      continue;
    }

    // Example: "file.go:10:5: error message"
    const goMatch = line.match(/^(.+):(\d+):(\d+):\s+(.+)$/);
    if (goMatch) {
      const message = goMatch[4];
      let severity: DiagnosticSeverity = 'error';
      
      if (message.toLowerCase().includes('warning')) {
        severity = 'warning';
      } else if (message.toLowerCase().includes('info') || message.toLowerCase().includes('note')) {
        severity = 'info';
      }

      diagnostics.push({
        severity,
        message,
        line: parseInt(goMatch[2], 10),
        column: parseInt(goMatch[3], 10),
      });
      continue;
    }

    // Generic error pattern
    if (line.toLowerCase().includes('error') || line.toLowerCase().includes('failed')) {
      diagnostics.push({
        severity: 'error',
        message: line,
        line: 1,
        column: 1,
      });
    }
  }

  return diagnostics;
}

/**
 * Group diagnostics by severity
 * @param diagnostics - Array of diagnostics
 * @returns Object with diagnostics grouped by severity
 */
export function groupBySeverity(diagnostics: DiagnosticMessage[]) {
  return {
    errors: diagnostics.filter(d => d.severity === 'error'),
    warnings: diagnostics.filter(d => d.severity === 'warning'),
    info: diagnostics.filter(d => d.severity === 'info'),
  };
}

/**
 * Check if there are any errors in the diagnostics
 * @param diagnostics - Array of diagnostics
 * @returns True if there are any errors
 */
export function hasErrors(diagnostics: DiagnosticMessage[]): boolean {
  return diagnostics.some(d => d.severity === 'error');
}

/**
 * Format a diagnostic message for display
 * @param diagnostic - The diagnostic to format
 * @returns A formatted string
 */
export function formatDiagnostic(diagnostic: DiagnosticMessage): string {
  const icon = getSeverityIcon(diagnostic.severity);
  const code = diagnostic.code ? `[${diagnostic.code}] ` : '';
  return `${icon} Line ${diagnostic.line}:${diagnostic.column} - ${code}${diagnostic.message}`;
}
