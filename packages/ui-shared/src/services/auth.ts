import { apiClient } from './api'

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name?: string
}

export interface AuthResponse {
  token: string
  user: {
    id: string
    email: string
    name?: string
  }
}

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
