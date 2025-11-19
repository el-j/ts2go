import { ref, onUnmounted } from 'vue'
import { usePlatform } from './usePlatform'
import { listen, type UnlistenFn } from '@tauri-apps/api/event'

export interface WebSocketMessage {
  type: string
  data: any
}

export function useWebSocket(url?: string) {
  const { isTauri, isWeb } = usePlatform()
  const ws = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const error = ref<string | null>(null)
  const messages = ref<WebSocketMessage[]>([])
  let unlisten: UnlistenFn | null = null

  // Web WebSocket connection
  function connectWeb(wsUrl: string) {
    try {
      ws.value = new WebSocket(wsUrl)

      ws.value.onopen = () => {
        isConnected.value = true
        error.value = null
      }

      ws.value.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)
          messages.value.push(message)
        } catch (e) {
          console.error('Failed to parse WebSocket message:', e)
        }
      }

      ws.value.onerror = (event) => {
        error.value = 'WebSocket error occurred'
        console.error('WebSocket error:', event)
      }

      ws.value.onclose = () => {
        isConnected.value = false
      }
    } catch (e: any) {
      error.value = e.message
    }
  }

  // Tauri event listener (alternative to WebSocket)
  async function connectTauri() {
    try {
      unlisten = await listen('backend-event', (event) => {
        messages.value.push(event.payload as WebSocketMessage)
      })
      isConnected.value = true
    } catch (e: any) {
      error.value = e.message
    }
  }

  function connect(wsUrl?: string) {
    if (isWeb) {
      const finalUrl = wsUrl || `ws://localhost:8080/ws`
      connectWeb(finalUrl)
    } else if (isTauri) {
      connectTauri()
    }
  }

  function disconnect() {
    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
    if (unlisten) {
      unlisten()
      unlisten = null
    }
    isConnected.value = false
  }

  function send(message: any) {
    if (isWeb && ws.value && isConnected.value) {
      ws.value.send(JSON.stringify(message))
    }
    // For Tauri, messages are typically sent via invoke commands
  }

  function clearMessages() {
    messages.value = []
  }

  onUnmounted(() => {
    disconnect()
  })

  // Auto-connect if URL provided
  if (url) {
    connect(url)
  }

  return {
    isConnected,
    error,
    messages,
    connect,
    disconnect,
    send,
    clearMessages,
  }
}
