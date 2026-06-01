<template>
  <header class="relative z-20 px-6 py-4">
    <nav class="mx-auto flex max-w-6xl items-center justify-between">
      <!-- Logo + site name -->
      <router-link to="/home" class="flex items-center gap-2">
        <div class="h-10 w-10 overflow-hidden rounded-xl text-gray-900 shadow-md dark:text-white">
          <BrandLogo :src="siteLogo" alt="壹站 Logo" :stroke-width="24" />
        </div>
        <span class="text-lg font-bold text-gray-900 dark:text-white">{{ siteName }}</span>
      </router-link>

      <div class="flex items-center gap-2 sm:gap-3">
        <!-- Page links -->
        <router-link
          v-for="link in links"
          :key="link.to"
          :to="link.to"
          :class="[
            'rounded-lg px-3 py-2 text-sm font-medium transition-colors',
            isActive(link.to)
              ? 'text-primary-600 dark:text-primary-400'
              : 'text-gray-600 hover:text-primary-600 dark:text-dark-300 dark:hover:text-primary-400'
          ]"
        >
          {{ link.label }}
        </router-link>

        <LocaleSwitcher />
        <StyleSwitcher />
        <button
          @click="toggleTheme"
          class="rounded-lg p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
          :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
        >
          <Icon v-if="isDark" name="sun" size="md" />
          <Icon v-else name="moon" size="md" />
        </button>
        <router-link
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="inline-flex items-center rounded-full bg-gray-900 px-3 py-1 text-xs font-medium text-white hover:bg-gray-800 dark:bg-gray-800"
        >
          {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
        </router-link>
      </div>
    </nav>
  </header>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import BrandLogo from '@/components/common/BrandLogo.vue'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import StyleSwitcher from '@/components/common/StyleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()
const route = useRoute()
const authStore = useAuthStore()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const pricingEnabled = computed(() => appStore.cachedPublicSettings?.pricing_page_enabled === true)
const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))

const links = computed(() => {
  const items = [{ to: '/home', label: t('home.navHome') }]
  if (pricingEnabled.value) items.push({ to: '/pricing', label: t('pricing.nav') })
  items.push({ to: '/docs', label: t('home.docs') })
  return items
})

function isActive(to: string): boolean {
  return route.path === to
}

const isDark = ref(document.documentElement.classList.contains('dark'))
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}
</script>
