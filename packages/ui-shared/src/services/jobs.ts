import { apiClient } from './api'
import type { ApiProject as Project, CreateProjectRequest, UpdateProjectRequest } from '../types'

export interface Job {
  id: string
  project_id: string
  status: 'pending' | 'running' | 'completed' | 'failed'
  type: 'transpile' | 'build' | 'test'
  created_at: string
  updated_at: string
  completed_at?: string
  error?: string
  result?: any
}

export class JobService {
  async listJobs(projectId: string): Promise<Job[]> {
    const response = await apiClient.get<{ jobs: Job[] }>(`/api/v1/projects/${projectId}/jobs`)
    return response.jobs || []
  }

  async getJob(projectId: string, jobId: string): Promise<Job> {
    return apiClient.get<Job>(`/api/v1/projects/${projectId}/jobs/${jobId}`)
  }

  async createJob(projectId: string, type: Job['type'], payload?: any): Promise<Job> {
    return apiClient.post<Job>(`/api/v1/projects/${projectId}/jobs`, { type, payload })
  }

  async cancelJob(projectId: string, jobId: string): Promise<void> {
    await apiClient.delete(`/api/v1/projects/${projectId}/jobs/${jobId}`)
  }
}

export const jobService = new JobService()

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
