<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { useProcessesStore } from '../stores/processes'

const props = defineProps<{
  show?: boolean
}>()

const emit = defineEmits<{
  close: []
  'update:show': [value: boolean]
}>()

const processesStore = useProcessesStore()

const title = ref('')
const description = ref('')
const priority = ref<'low' | 'medium' | 'high'>('medium')
const titleInput = ref<HTMLInputElement>()

const reset = () => {
  title.value = ''
  description.value = ''
  priority.value = 'medium'
}

const submit = async () => {
  if (!title.value.trim()) return

  try {
    await processesStore.createProcess({
      title: title.value.trim(),
      description: description.value.trim() || undefined,
      priority: priority.value,
    })
    reset()
    emit('close')
    emit('update:show', false)
  } catch (e) {
    console.error('Failed to create process:', e)
  }
}

const onBackdropClick = (e: MouseEvent) => {
  if (e.target === e.currentTarget) {
    emit('close')
    emit('update:show', false)
  }
}

const onKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    e.preventDefault()
    emit('close')
    emit('update:show', false)
    return
  }
  if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) {
    e.preventDefault()
    submit()
  }
}

// Focus title input when modal opens
watch(() => props.show, async (show) => {
  if (show) {
    reset()
    await nextTick()
    titleInput.value?.focus()
  }
})

defineExpose({
  focusTitle: () => titleInput.value?.focus(),
})
</script>

<template>
  <Transition name="fade">
    <div v-if="show" class="modal-backdrop" @click="onBackdropClick">
      <div class="modal-content" @keydown="onKeydown">
        <div class="modal-header">
          <h2>Create Process</h2>
          <button class="close-btn" @click="emit('close')">×</button>
        </div>

        <form @submit.prevent="submit" class="modal-form">
          <div class="form-field">
            <label for="title">Title *</label>
            <input
              id="title"
              ref="titleInput"
              v-model="title"
              type="text"
              required
              placeholder="Process title..."
              class="form-input"
            />
          </div>

          <div class="form-field">
            <label for="description">Description</label>
            <textarea
              id="description"
              v-model="description"
              rows="4"
              placeholder="Optional description..."
              class="form-textarea"
            ></textarea>
          </div>

          <div class="form-field">
            <label>Priority</label>
            <div class="priority-options">
              <label class="priority-option">
                <input
                  v-model="priority"
                  type="radio"
                  value="low"
                />
                <span>Low</span>
              </label>
              <label class="priority-option">
                <input
                  v-model="priority"
                  type="radio"
                  value="medium"
                />
                <span>Medium</span>
              </label>
              <label class="priority-option">
                <input
                  v-model="priority"
                  type="radio"
                  value="high"
                />
                <span>High</span>
              </label>
            </div>
          </div>

          <div class="form-actions">
            <button type="button" class="btn-secondary" @click="emit('close')">
              Cancel
            </button>
            <button type="submit" class="btn-primary" :disabled="!title.trim()">
              Create Process
            </button>
          </div>
        </form>

        <div class="modal-footer">
          <span class="hint"><kbd>Ctrl+Enter</kbd> to submit · <kbd>Escape</kbd> to close</span>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dark .modal-backdrop {
  background: rgba(0, 0, 0, 0.7);
}

.modal-content {
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
}

.modal-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
  color: var(--text-h);
}

.close-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  color: var(--text);
  font-size: 24px;
  line-height: 1;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.15s;
}

.close-btn:hover {
  background: var(--code-bg);
}

.modal-form {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
}

.form-field {
  margin-bottom: 16px;
}

.form-field label {
  display: block;
  margin-bottom: 6px;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-h);
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--border);
  border-radius: 4px;
  background: var(--bg);
  color: var(--text-h);
  font-size: 14px;
  font-family: inherit;
  transition: border-color 0.15s;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: var(--accent);
}

.form-input::placeholder,
.form-textarea::placeholder {
  color: var(--text);
  opacity: 0.5;
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.priority-options {
  display: flex;
  gap: 16px;
}

.priority-option {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-size: 14px;
  color: var(--text);
}

.priority-option input[type="radio"] {
  accent-color: var(--accent);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.btn-primary,
.btn-secondary {
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}

.btn-primary {
  background: var(--accent);
  color: white;
  border: none;
}

.btn-primary:hover:not(:disabled) {
  filter: brightness(1.1);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: transparent;
  color: var(--text-h);
  border: 1px solid var(--border);
}

.btn-secondary:hover {
  background: var(--code-bg);
}

.modal-footer {
  padding: 12px 20px;
  border-top: 1px solid var(--border);
  background: var(--code-bg);
}

.hint {
  font-size: 12px;
  color: var(--text);
}

.hint kbd {
  display: inline;
  padding: 2px 6px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 3px;
  font-family: var(--mono);
  font-size: 11px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
