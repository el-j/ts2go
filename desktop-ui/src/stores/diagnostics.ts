import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { DiagnosticMessage } from '../utils/diagnostics';

export const useDiagnosticsStore = defineStore('diagnostics', () => {
  // State
  const diagnosticsByFile = ref<Map<string, DiagnosticMessage[]>>(new Map());
  const activeDiagnostic = ref<DiagnosticMessage | null>(null);

  // Getters
  const allDiagnostics = computed(() => {
    const all: (DiagnosticMessage & { file?: string })[] = [];
    diagnosticsByFile.value.forEach((diagnostics, file) => {
      diagnostics.forEach(diagnostic => {
        all.push({ ...diagnostic, file });
      });
    });
    return all.sort((a, b) => {
      // Sort by severity: error > warning > info
      const severityOrder = { error: 0, warning: 1, info: 2 };
      const severityDiff = severityOrder[a.severity] - severityOrder[b.severity];
      if (severityDiff !== 0) return severityDiff;
      // Then by line number
      return a.line - b.line;
    });
  });

  const errorCount = computed(() => {
    return allDiagnostics.value.filter(d => d.severity === 'error').length;
  });

  const warningCount = computed(() => {
    return allDiagnostics.value.filter(d => d.severity === 'warning').length;
  });

  const infoCount = computed(() => {
    return allDiagnostics.value.filter(d => d.severity === 'info').length;
  });

  const totalCount = computed(() => {
    return allDiagnostics.value.length;
  });

  // Actions
  function setDiagnostics(file: string, diagnostics: DiagnosticMessage[]) {
    if (diagnostics.length === 0) {
      diagnosticsByFile.value.delete(file);
    } else {
      diagnosticsByFile.value.set(file, diagnostics);
    }
  }

  function addDiagnostic(file: string, diagnostic: DiagnosticMessage) {
    const existing = diagnosticsByFile.value.get(file) || [];
    diagnosticsByFile.value.set(file, [...existing, diagnostic]);
  }

  function clearFile(file: string) {
    diagnosticsByFile.value.delete(file);
  }

  function clearAll() {
    diagnosticsByFile.value.clear();
    activeDiagnostic.value = null;
  }

  function getDiagnosticsForFile(file: string): DiagnosticMessage[] {
    return diagnosticsByFile.value.get(file) || [];
  }

  function setActiveDiagnostic(diagnostic: DiagnosticMessage | null) {
    activeDiagnostic.value = diagnostic;
  }

  return {
    // State
    diagnosticsByFile,
    activeDiagnostic,
    // Getters
    allDiagnostics,
    errorCount,
    warningCount,
    infoCount,
    totalCount,
    // Actions
    setDiagnostics,
    addDiagnostic,
    clearFile,
    clearAll,
    getDiagnosticsForFile,
    setActiveDiagnostic
  };
});
