// API Service Layer
import axios from 'axios'
import type {
  Process,
  Log,
  CreateProcessRequest,
  UpdateProcessStatusRequest,
  AddLogRequest,
} from '../types/api'

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Process API
export const processApi = {
  // List all processes
  list: async (status?: string): Promise<Process[]> => {
    const params = status ? { status } : {}
    const response = await api.get<Process[]>('/processes', { params })
    return response.data
  },

  // Create a new process
  create: async (data: CreateProcessRequest): Promise<Process> => {
    const response = await api.post<Process>('/processes', data)
    return response.data
  },

  // Get a single process
  get: async (id: number): Promise<Process> => {
    const response = await api.get<Process>(`/processes/${id}`)
    return response.data
  },

  // Update process status
  updateStatus: async (id: number, data: UpdateProcessStatusRequest): Promise<{ status: string }> => {
    const response = await api.put<{ status: string }>(`/processes/${id}/status`, data)
    return response.data
  },

  // Delete a process
  delete: async (id: number): Promise<{ status: string }> => {
    const response = await api.delete<{ status: string }>(`/processes/${id}`)
    return response.data
  },

  // Get process logs
  getLogs: async (id: number): Promise<Log[]> => {
    const response = await api.get<Log[]>(`/processes/${id}/logs`)
    return response.data
  },

  // Add a log entry
  addLog: async (id: number, data: AddLogRequest): Promise<Log> => {
    const response = await api.post<Log>(`/processes/${id}/logs`, data)
    return response.data
  },
}

// Search API
export const searchApi = {
  // Global search
  search: async (query: string): Promise<Process[]> => {
    const response = await api.get<Process[]>('/search', { params: { q: query } })
    return response.data
  },
}

export default api
