// API Type Definitions

export type ProcessStatus = 'running' | 'blocked' | 'suspended' | 'terminated'
export type ProcessPriority = 'low' | 'medium' | 'high'
export type LogType = 'state_change' | 'progress'

export interface Process {
  id: number
  parent_id?: number | null
  title: string
  description: string
  status: ProcessStatus
  priority: ProcessPriority
  ranking: number
  created_at: string
  updated_at: string
  parent?: Process | null
  children?: Process[]
  logs?: Log[]
}

export interface Log {
  id: number
  process_id: number
  log_type: LogType
  content: string
  created_at: string
}

export interface CreateProcessRequest {
  title: string
  description?: string
  parent_id?: number | null
  priority?: ProcessPriority
}

export interface UpdateProcessStatusRequest {
  status: ProcessStatus
  reason?: string
}

export interface AddLogRequest {
  log_type: LogType
  content: string
}

export interface ApiError {
  error: string
}

// Theme types
export type ThemeMode = 'light' | 'dark' | 'system'
