import { apiClient } from './api'
import type { LoginRequest, RegisterRequest, AuthResponse } from '../types'

export class AuthService {
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/v1/auth/login', credentials)
    apiClient.setToken(response.token)
    return response
  }

  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/v1/auth/register', data)
    apiClient.setToken(response.token)
    return response
  }

  async logout(): Promise<void> {
    apiClient.setToken(null)
  }

  async getCurrentUser() {
    return apiClient.get('/api/v1/auth/me')
  }

  getToken(): string | null {
    if (typeof window === 'undefined') return null
    return localStorage.getItem('auth_token')
  }

  isAuthenticated(): boolean {
    return !!this.getToken()
  }
}

export const authService = new AuthService()
