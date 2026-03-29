<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { searchApi } from '../services/api'
import type { Process } from '../types/api'
import ProcessCard from '../components/ProcessCard.vue'

const router = useRouter()

const query = ref('')
const results = ref<Process[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const searchInput = ref<HTMLInputElement>()

const search = async () => {
  if (!query.value.trim()) {
    results.value = []
    return
  }

  loading.value = true
  error.value = null
  try {
    results.value = await searchApi.search(query.value.trim())
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Search failed'
    results.value = []
  } finally {
    loading.value = false
  }
}

const searchTimeout = ref<ReturnType<typeof setTimeout>>()

watch(query, () => {
  // Debounce search
  clearTimeout(searchTimeout.value)
  searchTimeout.value = setTimeout(() => {
    if (query.value.trim()) {
      search()
    }
  }, 300)
})

const onKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    query.value = ''
    results.value = []
    router.push('/')
  }
}

onMounted(async () => {
  await nextTick()
  searchInput.value?.focus()
})

onUnmounted(() => {
  clearTimeout(searchTimeout.value)
})

const openProcess = (processId: number) => {
  router.push(`/process/${processId}`)
}
</script>

<template>
  <div class="search-view">
    <div class="search-header">
      <div class="search-input-wrapper">
        <span class="search-icon">🔍</span>
        <input
          ref="searchInput"
          v-model="query"
          type="text"
          placeholder="Search processes and logs... (press Escape to clear)"
          class="search-input"
          @keydown="onKeydown"
        />
      </div>
      <button
        v-if="query"
        class="clear-btn"
        @click="query = ''; results = []"
      >
        Clear
      </button>
    </div>

    <div class="search-content">
      <div v-if="loading" class="loading">
        Searching...
      </div>

      <div v-else-if="error" class="error">
        {{ error }}
      </div>

      <div v-else-if="query && !loading && results.length === 0" class="empty">
        <p>No results found for "{{ query }}"</p>
      </div>

      <div v-else-if="!query" class="empty">
        <p>Start typing to search...</p>
      </div>

      <div v-else class="results-header">
        <span class="results-count">{{ results.length }} result{{ results.length !== 1 ? 's' : '' }}</span>
      </div>

      <div v-if="results.length > 0" class="results-list">
        <ProcessCard
          v-for="process in results"
          :key="process.id"
          :process="process"
          @click="openProcess(process.id)"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.search-view {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.search-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
  background: var(--code-bg);
}

.search-input-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 4px;
  transition: border-color 0.15s;
}

.search-input-wrapper:focus-within {
  border-color: var(--accent);
}

.search-icon {
  font-size: 16px;
  opacity: 0.5;
}

.search-input {
  flex: 1;
  border: none;
  background: transparent;
  color: var(--text-h);
  font-size: 14px;
  font-family: inherit;
}

.search-input:focus {
  outline: none;
}

.search-input::placeholder {
  color: var(--text);
  opacity: 0.5;
}

.clear-btn {
  padding: 6px 12px;
  background: transparent;
  color: var(--text);
  border: 1px solid var(--border);
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: background 0.15s;
}

.clear-btn:hover {
  background: var(--bg);
}

.search-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
}

.loading,
.error,
.empty {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px 20px;
  color: var(--text);
  font-size: 14px;
}

.error {
  color: rgb(239, 68, 68);
}

.results-header {
  padding: 0 0 12px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 16px;
}

.results-count {
  font-size: 13px;
  color: var(--text);
  opacity: 0.7;
}

.results-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 12px;
}
</style>
