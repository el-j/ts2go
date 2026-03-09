<template>
  <div class="diagnostics-panel">
    <div class="diagnostics-header">
      <h3 class="text-lg font-semibold">Problems</h3>
      <div class="diagnostics-stats">
        <span v-if="errorCount > 0" class="error-count">
          ❌ {{ errorCount }}
        </span>
        <span v-if="warningCount > 0" class="warning-count">
          ⚠️ {{ warningCount }}
        </span>
        <span v-if="infoCount > 0" class="info-count">
          ℹ️ {{ infoCount }}
        </span>
        <span v-if="totalCount === 0" class="no-problems">
          ✓ No problems
        </span>
      </div>
      <button @click="clearDiagnostics" class="clear-btn" title="Clear all problems">
        Clear
      </button>
    </div>
    
    <div class="diagnostics-list" v-if="totalCount > 0">
      <div
        v-for="(diagnostic, index) in diagnostics"
        :key="index"
        :class="['diagnostic-item', `severity-${diagnostic.severity}`]"
        @click="jumpToLocation(diagnostic)"
      >
        <div class="diagnostic-icon">
          {{ getSeverityIcon(diagnostic.severity) }}
        </div>
        <div class="diagnostic-content">
          <div class="diagnostic-message">
            {{ diagnostic.message }}
          </div>
          <div class="diagnostic-location">
            {{ diagnostic.file || 'current file' }} ({{ diagnostic.line }}:{{ diagnostic.column }})
            <span v-if="diagnostic.code" class="diagnostic-code">
              [{{ diagnostic.code }}]
            </span>
          </div>
        </div>
      </div>
    </div>
    
    <div v-else class="no-diagnostics">
      <div class="text-center text-gray-500 py-8">
        <div class="text-4xl mb-2">✓</div>
        <div>No problems detected</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useDiagnosticsStore } from '../stores/diagnostics';
import { getSeverityIcon } from '../utils/diagnostics';

const diagnosticsStore = useDiagnosticsStore();

const diagnostics = computed(() => diagnosticsStore.allDiagnostics);
const errorCount = computed(() => diagnosticsStore.errorCount);
const warningCount = computed(() => diagnosticsStore.warningCount);
const infoCount = computed(() => diagnosticsStore.infoCount);
const totalCount = computed(() => diagnosticsStore.totalCount);

function clearDiagnostics() {
  diagnosticsStore.clearAll();
}

function jumpToLocation(diagnostic: any) {
  // Emit event to scroll editor to the error location
  diagnosticsStore.setActiveDiagnostic(diagnostic);
}
</script>

<style scoped>
.diagnostics-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #1e1e1e;
  color: #d4d4d4;
}

.diagnostics-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid #3e3e42;
  background: #252526;
}

.diagnostics-header h3 {
  flex: 1;
  margin: 0;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: uppercase;
}

.diagnostics-stats {
  display: flex;
  gap: 1rem;
  font-size: 0.875rem;
}

.error-count {
  color: #f48771;
}

.warning-count {
  color: #cca700;
}

.info-count {
  color: #75beff;
}

.no-problems {
  color: #89d185;
}

.clear-btn {
  padding: 0.25rem 0.75rem;
  background: #0e639c;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 0.75rem;
  cursor: pointer;
  transition: background 0.2s;
}

.clear-btn:hover {
  background: #1177bb;
}

.diagnostics-list {
  flex: 1;
  overflow-y: auto;
  padding: 0.5rem;
}

.diagnostic-item {
  display: flex;
  gap: 0.75rem;
  padding: 0.75rem;
  margin-bottom: 0.5rem;
  background: #2d2d30;
  border-radius: 4px;
  border-left: 3px solid;
  cursor: pointer;
  transition: background 0.2s;
}

.diagnostic-item:hover {
  background: #37373d;
}

.diagnostic-item.severity-error {
  border-left-color: #f48771;
}

.diagnostic-item.severity-warning {
  border-left-color: #cca700;
}

.diagnostic-item.severity-info {
  border-left-color: #75beff;
}

.diagnostic-icon {
  font-size: 1.25rem;
  flex-shrink: 0;
}

.diagnostic-content {
  flex: 1;
  min-width: 0;
}

.diagnostic-message {
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
  word-wrap: break-word;
}

.diagnostic-location {
  font-size: 0.75rem;
  color: #888;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.diagnostic-code {
  font-family: 'Courier New', monospace;
  color: #888;
}

.no-diagnostics {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
