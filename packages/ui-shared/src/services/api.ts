import { usePlatform } from '../composables/usePlatform'
import { invoke } from '@tauri-apps/api/core'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// HTTP Client for Web
class HttpClient {
  private baseURL: string
  private token: string | null = null

  constructor(baseURL: string) {
    this.baseURL = baseURL
    this.token = localStorage.getItem('auth_token')
  }

  setToken(token: string | null) {
    this.token = token
    if (token) {
      localStorage.setItem('auth_token', token)
    } else {
      localStorage.removeItem('auth_token')
    }
  }

  private getHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    }
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
    }
    return headers
  }

  async get<T>(endpoint: string): Promise<T> {
    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'GET',
      headers: this.getHeaders(),
    })
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }
    return response.json()
  }

  async post<T>(endpoint: string, data?: any): Promise<T> {
    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: data ? JSON.stringify(data) : undefined,
    })
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }
    return response.json()
  }

  async put<T>(endpoint: string, data?: any): Promise<T> {
    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: data ? JSON.stringify(data) : undefined,
    })
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }
    return response.json()
  }

  async delete<T>(endpoint: string): Promise<T> {
    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    })
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }
    return response.json()
  }
}

// Tauri IPC Client
class TauriClient {
  async get<T>(endpoint: string): Promise<T> {
    return invoke(endpoint)
  }

  async post<T>(endpoint: string, data?: any): Promise<T> {
    return invoke(endpoint, data ? { payload: data } : undefined)
  }

  async put<T>(endpoint: string, data?: any): Promise<T> {
    return invoke(endpoint, data ? { payload: data } : undefined)
  }

  async delete<T>(endpoint: string): Promise<T> {
    return invoke(endpoint)
  }
}

// Unified API Client
export class ApiClient {
  private httpClient: HttpClient
  private tauriClient: TauriClient
  private platform = usePlatform()

  constructor() {
    this.httpClient = new HttpClient(API_BASE_URL)
    this.tauriClient = new TauriClient()
  }

  setToken(token: string | null) {
    if (this.platform.isWeb) {
      this.httpClient.setToken(token)
    }
  }

  async get<T>(endpoint: string): Promise<T> {
    if (this.platform.isTauri) {
      return this.tauriClient.get<T>(endpoint)
    }
    return this.httpClient.get<T>(endpoint)
  }

  async post<T>(endpoint: string, data?: any): Promise<T> {
    if (this.platform.isTauri) {
      return this.tauriClient.post<T>(endpoint, data)
    }
    return this.httpClient.post<T>(endpoint, data)
  }

  async put<T>(endpoint: string, data?: any): Promise<T> {
    if (this.platform.isTauri) {
      return this.tauriClient.put<T>(endpoint, data)
    }
    return this.httpClient.put<T>(endpoint, data)
  }

  async delete<T>(endpoint: string): Promise<T> {
    if (this.platform.isTauri) {
      return this.tauriClient.delete<T>(endpoint)
    }
    return this.httpClient.delete<T>(endpoint)
  }
}

export const apiClient = new ApiClient()
