<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [content: string]
}>()

const content = ref('')
const textarea = ref<HTMLTextAreaElement>()

const submit = () => {
  const trimmed = content.value.trim()
  if (!trimmed) return

  emit('submit', trimmed)
  content.value = ''
}

const focus = () => {
  textarea.value?.focus()
}

defineExpose({
  focus,
})
</script>

<template>
  <div class="add-log-form">
    <textarea
      ref="textarea"
      v-model="content"
      rows="3"
      placeholder="Add a progress log..."
      class="log-textarea"
      @keydown.ctrl.enter="submit"
      @keydown.meta.enter="submit"
    ></textarea>
    <div class="form-footer">
      <span class="hint">Ctrl+Enter to submit</span>
      <button
        type="button"
        class="btn-submit"
        :disabled="!content.trim() || loading"
        @click="submit"
      >
        Add Log
      </button>
    </div>
  </div>
</template>

<style scoped>
.add-log-form {
  padding: 16px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background: var(--code-bg);
}

.log-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: 4px;
  background: var(--bg);
  color: var(--text-h);
  font-size: 14px;
  font-family: inherit;
  resize: vertical;
  min-height: 60px;
  transition: border-color 0.15s;
}

.log-textarea:focus {
  outline: none;
  border-color: var(--accent);
}

.log-textarea::placeholder {
  color: var(--text);
  opacity: 0.5;
}

.form-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
}

.hint {
  font-size: 12px;
  color: var(--text);
  opacity: 0.7;
}

.btn-submit {
  padding: 6px 14px;
  background: var(--accent);
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: filter 0.15s;
}

.btn-submit:hover:not(:disabled) {
  filter: brightness(1.1);
}

.btn-submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
