/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL?: string
  readonly VITE_WS_URL?: string
  readonly VITE_PLATFORM?: string
  readonly VITE_SAAS_API_URL?: string
  readonly VITE_PROJECT_ID?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
