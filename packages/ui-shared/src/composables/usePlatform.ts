export type Platform = 'tauri' | 'web'

export function usePlatform() {
  const isTauri = typeof window !== 'undefined' && !!(window as any).__TAURI__
  const isWeb = !isTauri
  const platform: Platform = isTauri ? 'tauri' : 'web'

  return {
    isTauri,
    isWeb,
    platform
  }
}
