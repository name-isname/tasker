// Processes Store
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { processApi } from '../services/api'
import type { Process, ProcessStatus } from '../types/api'

export const useProcessesStore = defineStore('processes', () => {
  // State
  const processes = ref<Process[]>([])
  const loading = ref<boolean>(false)
  const error = ref<string | null>(null)
  const statusFilter = ref<ProcessStatus | 'all'>('all')

  // Fetch all processes
  const fetchProcesses = async () => {
    loading.value = true
    error.value = null
    try {
      const status = statusFilter.value === 'all' ? undefined : statusFilter.value
      processes.value = await processApi.list(status)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch processes'
    } finally {
      loading.value = false
    }
  }

  // Create a new process
  const createProcess = async (data: { title: string; description?: string; priority?: string }) => {
    loading.value = true
    error.value = null
    try {
      await processApi.create({
        title: data.title,
        description: data.description,
        priority: (data.priority as any) || 'medium',
      })
      await fetchProcesses()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to create process'
      throw e
    } finally {
      loading.value = false
    }
  }

  // Update process status
  const updateStatus = async (id: number, status: ProcessStatus, reason?: string) => {
    loading.value = true
    error.value = null
    try {
      await processApi.updateStatus(id, { status, reason })
      await fetchProcesses()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to update status'
      throw e
    } finally {
      loading.value = false
    }
  }

  // Delete a process
  const deleteProcess = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      await processApi.delete(id)
      await fetchProcesses()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to delete process'
      throw e
    } finally {
      loading.value = false
    }
  }

  // Set status filter
  const setStatusFilter = (filter: ProcessStatus | 'all') => {
    statusFilter.value = filter
    fetchProcesses()
  }

  return {
    processes,
    loading,
    error,
    statusFilter,
    fetchProcesses,
    createProcess,
    updateStatus,
    deleteProcess,
    setStatusFilter,
  }
})
