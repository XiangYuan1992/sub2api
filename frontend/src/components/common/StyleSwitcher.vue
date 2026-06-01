<template>
  <button
    type="button"
    @click="toggleStyle"
    class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
    :title="style === 'openai' ? t('theme.switchToClassic') : t('theme.switchToOpenai')"
  >
    <Icon name="sparkles" size="md" />
  </button>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { useUiStyle } from '@/composables/useUiStyle'

const { t } = useI18n()
const { style, toggleStyle } = useUiStyle()
</script>
