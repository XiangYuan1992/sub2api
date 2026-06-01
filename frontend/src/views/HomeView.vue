<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
    <iframe v-if="isHomeContentUrl" :src="homeContent.trim()" class="h-screen w-full border-0" allowfullscreen></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Default Home Page -->
  <div v-else class="relative flex min-h-screen flex-col bg-white dark:bg-dark-950">
    <PublicNav />

    <main class="relative z-10 flex-1">
      <!-- Hero -->
      <section class="px-6 py-16 lg:py-24">
        <div class="mx-auto flex max-w-6xl flex-col items-center gap-12 lg:flex-row lg:gap-16">
          <div class="flex-1 text-center lg:text-left">
            <span class="inline-block rounded-full bg-primary-50 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-primary-600 dark:bg-dark-800 dark:text-primary-400">
              {{ t('home.landing.eyebrow') }}
            </span>
            <h1 class="mt-5 text-4xl font-bold tracking-tight text-gray-900 dark:text-white md:text-5xl lg:text-6xl">
              {{ t('home.landing.heroTitle1') }}<br />{{ t('home.landing.heroTitle2') }}
            </h1>
            <p class="mt-6 max-w-xl text-base leading-relaxed text-gray-600 dark:text-dark-300 lg:text-lg">
              {{ t('home.landing.heroDesc') }}
            </p>
            <div class="mt-8 flex flex-wrap items-center justify-center gap-3 lg:justify-start">
              <router-link :to="isAuthenticated ? dashboardPath : '/login'" class="btn btn-primary px-7 py-3 text-base">
                {{ t('home.landing.ctaConsole') }}
                <Icon name="arrowRight" size="md" class="ml-2" :stroke-width="2" />
              </router-link>
              <router-link v-if="pricingEnabled" to="/pricing" class="btn btn-secondary px-7 py-3 text-base">
                {{ t('home.landing.ctaPricing') }}
              </router-link>
              <router-link to="/docs" class="btn btn-secondary px-7 py-3 text-base">
                {{ t('home.landing.ctaDocs') }}
              </router-link>
            </div>
          </div>

          <!-- Right: Model hub -->
          <div class="flex flex-1 justify-center lg:justify-end">
            <div class="relative h-72 w-72">
              <div class="absolute left-1/2 top-1/2 flex h-16 w-16 -translate-x-1/2 -translate-y-1/2 items-center justify-center rounded-2xl bg-gray-900 text-white shadow-lg dark:bg-neutral-100 dark:text-gray-900 dark:shadow-black/30">
                <Icon name="menu" size="lg" />
              </div>
              <div
                v-for="(node, i) in hubNodes"
                :key="node.label"
                class="absolute flex h-20 w-20 flex-col items-center justify-center gap-1 rounded-full border border-gray-200 bg-white text-center shadow-sm dark:border-dark-700 dark:bg-dark-900"
                :style="nodeStyle(i)"
              >
                <Icon :name="node.icon" size="md" class="text-primary-600 dark:text-primary-400" />
                <span class="text-[10px] font-medium text-gray-600 dark:text-dark-300">{{ node.label }}</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Ecosystem / Clients -->
      <section class="bg-gray-50 px-6 py-16 dark:bg-dark-900/40">
        <div class="mx-auto max-w-6xl">
          <p class="text-sm font-semibold uppercase tracking-wide text-primary-600 dark:text-primary-400">{{ t('home.landing.ecosystem') }}</p>
          <h2 class="mt-2 text-3xl font-bold text-gray-900 dark:text-white">{{ t('home.landing.clientsTitle') }}</h2>
          <div class="mt-8 grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
            <router-link
              v-for="c in clients"
              :key="c.name"
              to="/docs"
              class="group rounded-2xl border border-gray-200 bg-white p-6 transition-colors hover:border-gray-900 dark:border-dark-700 dark:bg-dark-900 dark:hover:border-white"
            >
              <div class="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-gray-900 text-white dark:bg-dark-700">
                <Icon :name="c.icon" size="md" class="text-white" />
              </div>
              <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">{{ c.name }}</h3>
              <p class="mb-4 text-sm leading-relaxed text-gray-500 dark:text-dark-400">{{ c.desc }}</p>
              <span class="text-xs font-semibold uppercase tracking-wide text-primary-600 dark:text-primary-400">{{ c.tag }}</span>
            </router-link>
          </div>
        </div>
      </section>

      <!-- Pricing -->
      <section class="px-6 py-16">
        <div class="mx-auto max-w-6xl text-center">
          <p class="text-sm font-semibold uppercase tracking-wide text-primary-600 dark:text-primary-400">{{ t('home.landing.pricingEyebrow') }}</p>
          <h2 class="mt-2 text-3xl font-bold text-gray-900 dark:text-white md:text-4xl">{{ t('home.landing.pricingTitle') }}</h2>
          <p class="mt-3 text-gray-600 dark:text-dark-300">{{ t('home.landing.pricingSubtitle') }}</p>
          <div class="mt-10 grid gap-6 text-left lg:grid-cols-3">
            <div
              v-for="plan in plans"
              :key="plan.name"
              :class="[
                'relative rounded-2xl border bg-white p-6 dark:bg-dark-900',
                plan.featured ? 'border-gray-900 dark:border-white' : 'border-gray-200 dark:border-dark-700'
              ]"
            >
              <span v-if="plan.featured" class="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-gray-900 px-3 py-0.5 text-[10px] font-semibold text-white dark:bg-white dark:text-gray-900">{{ t('home.landing.recommended') }}</span>
              <p class="text-xs font-semibold uppercase tracking-wide text-gray-400">{{ plan.tag }}</p>
              <h3 class="mt-2 text-xl font-bold text-gray-900 dark:text-white">{{ plan.name }}</h3>
              <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ plan.sub }}</p>
              <ul class="mt-5 space-y-2.5">
                <li v-for="f in plan.features" :key="f" class="flex items-start gap-2 text-sm text-gray-600 dark:text-dark-300">
                  <Icon name="check" size="sm" class="mt-0.5 shrink-0 text-primary-600 dark:text-primary-400" />
                  <span>{{ f }}</span>
                </li>
              </ul>
              <router-link :to="plan.to" :class="['btn mt-6 w-full justify-center py-2.5', plan.featured ? 'btn-primary' : 'btn-secondary']">
                {{ plan.cta }}
              </router-link>
            </div>
          </div>
        </div>
      </section>

      <!-- Value -->
      <section class="bg-gray-50 px-6 py-16 dark:bg-dark-900/40">
        <div class="mx-auto max-w-6xl">
          <p class="text-sm font-semibold uppercase tracking-wide text-primary-600 dark:text-primary-400">{{ t('home.landing.valueEyebrow') }}</p>
          <h2 class="mt-2 max-w-3xl text-3xl font-bold text-gray-900 dark:text-white">{{ t('home.landing.valueTitle') }}</h2>
          <div class="mt-8 grid gap-5 sm:grid-cols-2">
            <div v-for="v in values" :key="v.title" class="rounded-2xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-900">
              <div class="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-primary-50 text-primary-600 dark:bg-dark-800 dark:text-primary-400">
                <Icon :name="v.icon" size="md" />
              </div>
              <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">{{ v.title }}</h3>
              <p class="text-sm leading-relaxed text-gray-500 dark:text-dark-400">{{ v.desc }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- FAQ -->
      <section class="px-6 py-16">
        <div class="mx-auto max-w-3xl text-center">
          <p class="text-sm font-semibold uppercase tracking-wide text-primary-600 dark:text-primary-400">{{ t('home.landing.faqEyebrow') }}</p>
          <h2 class="mt-2 text-3xl font-bold text-gray-900 dark:text-white md:text-4xl">{{ t('home.landing.faqTitle') }}</h2>
          <div class="mt-8 space-y-3 text-left">
            <div v-for="(f, i) in faqs" :key="i" class="rounded-xl border border-gray-200 dark:border-dark-700">
              <button @click="openFaq = openFaq === i ? -1 : i" class="flex w-full items-center justify-between px-5 py-4 text-left text-sm font-medium text-gray-900 dark:text-white">
                {{ f.q }}
                <Icon :name="openFaq === i ? 'chevronUp' : 'plus'" size="sm" class="shrink-0 text-gray-400" />
              </button>
              <p v-if="openFaq === i" class="px-5 pb-4 text-sm leading-relaxed text-gray-600 dark:text-dark-300">{{ f.a }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Final CTA -->
      <section class="px-6 pb-20">
        <div class="mx-auto max-w-5xl rounded-3xl border border-gray-200 bg-white px-6 py-14 text-center shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white md:text-4xl">{{ t('home.landing.finalTitle') }}</h2>
          <p class="mx-auto mt-4 max-w-xl text-gray-600 dark:text-dark-300">{{ t('home.landing.finalDesc') }}</p>
          <router-link :to="isAuthenticated ? dashboardPath : '/login'" class="btn btn-primary mt-8 px-8 py-3 text-base">
            {{ t('home.landing.ctaConsole') }}
          </router-link>
        </div>
      </section>
    </main>

    <!-- Footer -->
    <footer class="border-t border-gray-200/60 px-6 py-8 dark:border-dark-800/60">
      <div class="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 sm:flex-row">
        <p class="text-sm text-gray-500 dark:text-dark-400">&copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}</p>
        <div class="flex items-center gap-5 text-sm">
          <router-link v-if="pricingEnabled" to="/pricing" class="text-gray-500 hover:text-gray-900 dark:text-dark-400 dark:hover:text-white">{{ t('pricing.nav') }}</router-link>
          <router-link to="/docs" class="text-gray-500 hover:text-gray-900 dark:text-dark-400 dark:hover:text-white">{{ t('home.docs') }}</router-link>
          <a :href="githubUrl" target="_blank" rel="noopener noreferrer" class="text-gray-500 hover:text-gray-900 dark:text-dark-400 dark:hover:text-white">GitHub</a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import PublicNav from '@/components/common/PublicNav.vue'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const pricingEnabled = computed(() => appStore.cachedPublicSettings?.pricing_page_enabled === true)
const isHomeContentUrl = computed(() => {
  const c = homeContent.value.trim()
  return c.startsWith('http://') || c.startsWith('https://')
})

const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))
const currentYear = computed(() => new Date().getFullYear())
const githubUrl = 'https://github.com/Wei-Shaw/sub2api'

const hubNodes = [
  { label: 'Claude', icon: 'sparkles' as const },
  { label: 'All Models', icon: 'cube' as const },
  { label: 'Gemini', icon: 'bolt' as const },
  { label: 'ChatGPT', icon: 'chat' as const }
]
function nodeStyle(i: number) {
  const angle = (i / hubNodes.length) * 2 * Math.PI - Math.PI / 4
  const r = 120
  const x = 50 + (Math.cos(angle) * r) / 2.88
  const y = 50 + (Math.sin(angle) * r) / 2.88
  return { left: `${x}%`, top: `${y}%`, transform: 'translate(-50%, -50%)' }
}

const clients = computed(() => [
  { name: 'OpenClaw', icon: 'chatBubble' as const, desc: t('home.clients.openclaw'), tag: t('home.clients.openclawTag') },
  { name: 'Claude Code', icon: 'terminal' as const, desc: t('home.clients.claude'), tag: t('home.clients.claudeTag') },
  { name: 'Codex', icon: 'terminal' as const, desc: t('home.clients.codex'), tag: t('home.clients.codexTag') },
  { name: 'Hermes Agent', icon: 'bolt' as const, desc: t('home.clients.hermes'), tag: t('home.clients.hermesTag') }
])

const plans = computed(() => {
  const target = pricingEnabled.value ? '/pricing' : (isAuthenticated.value ? dashboardPath.value : '/login')
  return [
    { name: t('home.landing.planPaygoName'), tag: t('home.landing.planPaygoPrice'), sub: t('home.landing.planPaygoDesc'),
      features: [t('home.landing.planPaygoF1'), t('home.landing.planPaygoF2'), t('home.landing.planPaygoF3')],
      cta: t('home.landing.planPaygoCta'), to: target, featured: false },
    { name: t('home.landing.planClaudeName'), tag: t('home.landing.planClaudePrice'), sub: t('home.landing.planClaudePrice'),
      features: [t('home.landing.planClaudeF1'), t('home.landing.planClaudeF2'), t('home.landing.planClaudeF3')],
      cta: t('home.landing.planStart'), to: target, featured: true },
    { name: t('home.landing.planGptName'), tag: t('home.landing.planGptPrice'), sub: t('home.landing.planGptPrice'),
      features: [t('home.landing.planGptF1'), t('home.landing.planGptF2'), t('home.landing.planGptF3')],
      cta: t('home.landing.planStart'), to: target, featured: true }
  ]
})

const values = computed(() => [
  { icon: 'cloud' as const, title: t('home.landing.value1Title'), desc: t('home.landing.value1Desc') },
  { icon: 'server' as const, title: t('home.landing.value2Title'), desc: t('home.landing.value2Desc') },
  { icon: 'beaker' as const, title: t('home.landing.value3Title'), desc: t('home.landing.value3Desc') },
  { icon: 'shield' as const, title: t('home.landing.value4Title'), desc: t('home.landing.value4Desc') }
])

const faqs = computed(() => [
  { q: t('home.landing.faq1Q'), a: t('home.landing.faq1A') },
  { q: t('home.landing.faq2Q'), a: t('home.landing.faq2A') },
  { q: t('home.landing.faq3Q'), a: t('home.landing.faq3A') }
])
const openFaq = ref(0)

function initTheme() {
  const saved = localStorage.getItem('theme')
  if (saved === 'dark' || (!saved && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    document.documentElement.classList.add('dark')
  }
}

onMounted(() => {
  initTheme()
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) appStore.fetchPublicSettings()
})
</script>
