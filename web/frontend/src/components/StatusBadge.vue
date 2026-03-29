<script setup lang="ts">
import type { ProcessStatus } from '../types/api'

const props = defineProps<{
  status: ProcessStatus
}>()

const statusConfig: Record<ProcessStatus, { label: string; class: string }> = {
  running: { label: 'Running', class: 'status-running' },
  blocked: { label: 'Blocked', class: 'status-blocked' },
  suspended: { label: 'Suspended', class: 'status-suspended' },
  terminated: { label: 'Terminated', class: 'status-terminated' },
}

const config = statusConfig[props.status]
</script>

<template>
  <span :class="['status-badge', config.class]">
    {{ config.label }}
  </span>
</template>

<style scoped>
.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.status-running {
  background: rgba(34, 197, 94, 0.15);
  color: rgb(34, 197, 94);
}

.status-blocked {
  background: rgba(239, 68, 68, 0.15);
  color: rgb(239, 68, 68);
}

.status-suspended {
  background: rgba(234, 179, 8, 0.15);
  color: rgb(234, 179, 8);
}

.status-terminated {
  background: rgba(107, 114, 128, 0.15);
  color: rgb(107, 114, 128);
}

.dark .status-running {
  background: rgba(34, 197, 94, 0.2);
  color: rgb(74, 222, 128);
}

.dark .status-blocked {
  background: rgba(239, 68, 68, 0.2);
  color: rgb(248, 113, 113);
}

.dark .status-suspended {
  background: rgba(234, 179, 8, 0.2);
  color: rgb(250, 204, 21);
}

.dark .status-terminated {
  background: rgba(107, 114, 128, 0.2);
  color: rgb(156, 163, 175);
}
</style>
