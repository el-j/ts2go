import { apiClient } from './api'
import type { ApiProject as Project, CreateProjectRequest, UpdateProjectRequest } from '../types'

export class ProjectService {
  async listProjects(): Promise<Project[]> {
    const response = await apiClient.get<{ projects: Project[] }>('/api/v1/projects')
    return response.projects || []
  }

  async getProject(id: string): Promise<Project> {
    return apiClient.get<Project>(`/api/v1/projects/${id}`)
  }

  async createProject(data: CreateProjectRequest): Promise<Project> {
    return apiClient.post<Project>('/api/v1/projects', data)
  }

  async updateProject(id: string, data: UpdateProjectRequest): Promise<Project> {
    return apiClient.put<Project>(`/api/v1/projects/${id}`, data)
  }

  async deleteProject(id: string): Promise<void> {
    await apiClient.delete(`/api/v1/projects/${id}`)
  }
}

export const projectService = new ProjectService()
