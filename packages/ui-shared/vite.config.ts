import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { fileURLToPath, URL } from 'node:url'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const isTauri = mode === 'tauri'
  const isWeb = mode === 'web'

  return {
    plugins: [
      vue(),
      tailwindcss()
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    define: {
      'import.meta.env.VITE_PLATFORM': JSON.stringify(isTauri ? 'tauri' : 'web'),
    },
    build: {
      outDir: isTauri ? 'dist-tauri' : 'dist-web',
      // For web, we want a standard SPA build
      // For Tauri, we keep the same but it will be consumed by Tauri
      target: isWeb ? 'es2015' : 'esnext',
    },
    // Tauri expects a fixed port to communicate with the frontend
    server: isTauri ? {
      port: 1420,
      strictPort: true,
      watch: {
        ignored: ['**/src-tauri/**', '**/desktop/tauri-backend/**']
      }
    } : {
      port: 3000,
    }
  }
})
