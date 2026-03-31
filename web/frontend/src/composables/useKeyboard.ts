// Keyboard Shortcuts Composable
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

type KeyboardHandler = (e: KeyboardEvent) => void

// Module-level handlers storage - shared across all composable instances
const handlers = new Map<string, KeyboardHandler[]>()
let isGlobalListenerSetup = false

export function useKeyboard() {
  // Register a handler for a specific key
  const register = (key: string, handler: KeyboardHandler) => {
    if (!handlers.has(key)) {
      handlers.set(key, [])
    }
    handlers.get(key)!.push(handler)
  }

  // Unregister a handler
  const unregister = (key: string, handler: KeyboardHandler) => {
    const keyHandlers = handlers.get(key)
    if (keyHandlers) {
      const index = keyHandlers.indexOf(handler)
      if (index > -1) {
        keyHandlers.splice(index, 1)
      }
    }
  }

  // Handle keyboard events
  const handleKeydown = (e: KeyboardEvent) => {
    // Ignore if user is typing in an input
    const target = e.target as HTMLElement
    if (
      target.tagName === 'INPUT' ||
      target.tagName === 'TEXTAREA' ||
      target.contentEditable === 'true'
    ) {
      // Allow Escape to close modals even when typing
      if (e.key === 'Escape') {
        const escapeHandlers = handlers.get('Escape')
        if (escapeHandlers) {
          escapeHandlers.forEach((h) => h(e))
        }
      }
      return
    }

    const keyHandlers = handlers.get(e.key)
    if (keyHandlers) {
      keyHandlers.forEach((h) => h(e))
    }
  }

  // Set up global listener only once
  if (!isGlobalListenerSetup) {
    window.addEventListener('keydown', handleKeydown)
    isGlobalListenerSetup = true
  }

  return {
    register,
    unregister,
  }
}

// Global keyboard shortcuts registration
export function useGlobalKeyboard(showHelp?: { value: boolean }) {
  const { register } = useKeyboard()
  const router = useRouter()

  onMounted(() => {
    // c - Create process
    register('c', (_e) => {
      _e.preventDefault()
      window.dispatchEvent(new CustomEvent('open-create-modal'))
    })

    // / - Focus search
    register('/', (_e) => {
      _e.preventDefault()
      router.push('/search')
    })

    // n - Navigate to process list
    register('n', (_e) => {
      router.push('/')
    })

    // ? - Show shortcuts help
    register('?', (_e) => {
      _e.preventDefault()
      if (showHelp) {
        showHelp.value = true
      } else {
        // Fallback to dispatching event for component to handle
        window.dispatchEvent(new CustomEvent('show-keyboard-help'))
      }
    })
  })
}
