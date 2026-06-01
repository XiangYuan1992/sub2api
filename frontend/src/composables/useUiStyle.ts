import { ref } from 'vue'

export type UiStyle = 'openai' | 'classic'

const STORAGE_KEY = 'ui-style'
const DEFAULT_STYLE: UiStyle = 'openai'

function readStyle(): UiStyle {
  return localStorage.getItem(STORAGE_KEY) === 'classic' ? 'classic' : DEFAULT_STYLE
}

const style = ref<UiStyle>(readStyle())

/** Apply the style to <html data-theme>. openai is the default (no attribute). */
export function applyUiStyle(next: UiStyle): void {
  if (next === 'classic') {
    document.documentElement.dataset.theme = 'classic'
  } else {
    delete document.documentElement.dataset.theme
  }
}

export function useUiStyle() {
  function setStyle(next: UiStyle): void {
    style.value = next
    localStorage.setItem(STORAGE_KEY, next)
    applyUiStyle(next)
  }
  function toggleStyle(): void {
    setStyle(style.value === 'openai' ? 'classic' : 'openai')
  }
  return { style, setStyle, toggleStyle }
}
