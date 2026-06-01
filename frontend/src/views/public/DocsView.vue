<template>
  <div class="relative flex min-h-screen flex-col bg-white dark:bg-dark-950">
    <PublicNav />

    <main class="relative z-10 flex-1 px-6 py-10">
      <div class="mx-auto flex max-w-6xl flex-col gap-8 lg:flex-row">
        <!-- Sidebar -->
        <aside class="lg:w-56 lg:shrink-0">
          <nav class="sticky top-6 flex flex-row flex-wrap gap-1 lg:flex-col">
            <button
              v-for="s in sections"
              :key="s.id"
              @click="scrollTo(s.id)"
              :class="[
                'rounded-lg px-3 py-2 text-left text-sm font-medium transition-colors',
                activeId === s.id
                  ? 'bg-gray-100 text-gray-900 dark:bg-dark-800 dark:text-white'
                  : 'text-gray-500 hover:text-gray-800 dark:text-dark-400 dark:hover:text-dark-200'
              ]"
            >
              {{ s.title }}
            </button>
          </nav>
        </aside>

        <!-- Content -->
        <div class="min-w-0 flex-1">
          <h1 class="mb-2 text-3xl font-semibold tracking-tight text-gray-900 dark:text-white">{{ t('docs.title') }}</h1>
          <p class="mb-8 text-gray-600 dark:text-dark-300">{{ t('docs.subtitle') }}</p>

          <!-- Quick start -->
          <section id="quick-start" class="mb-12 scroll-mt-6">
            <h2 class="mb-3 text-xl font-semibold text-gray-900 dark:text-white">{{ t('docs.quickStart.title') }}</h2>
            <ol class="ml-5 list-decimal space-y-2 text-sm leading-relaxed text-gray-600 dark:text-dark-300">
              <li>{{ t('docs.quickStart.step1') }}</li>
              <li>{{ t('docs.quickStart.step2') }}</li>
              <li>{{ t('docs.quickStart.step3') }}</li>
            </ol>
            <div class="mt-4 flex flex-wrap items-center gap-2 text-sm">
              <span class="text-gray-500 dark:text-dark-400">{{ t('docs.baseUrlLabel') }}:</span>
              <code class="rounded bg-gray-100 px-2 py-1 font-mono text-gray-800 dark:bg-dark-800 dark:text-dark-100">{{ baseUrl }}</code>
            </div>
          </section>

          <!-- Tool sections -->
          <section v-for="tool in tools" :id="tool.id" :key="tool.id" class="mb-12 scroll-mt-6">
            <h2 class="mb-1 text-xl font-semibold text-gray-900 dark:text-white">{{ tool.title }}</h2>
            <p class="mb-4 text-sm leading-relaxed text-gray-600 dark:text-dark-300">{{ tool.desc }}</p>

            <template v-for="(step, i) in tool.steps" :key="i">
              <h3 class="mb-2 mt-5 text-sm font-semibold text-gray-800 dark:text-dark-100">{{ step.heading }}</h3>
              <p v-if="step.text" class="mb-2 text-sm leading-relaxed text-gray-600 dark:text-dark-300">{{ step.text }}</p>
              <CodeBlock v-if="step.code" :code="step.code" />
            </template>
          </section>

          <!-- FAQ -->
          <section id="faq" class="mb-12 scroll-mt-6">
            <h2 class="mb-3 text-xl font-semibold text-gray-900 dark:text-white">{{ t('docs.faq.title') }}</h2>
            <div class="space-y-4">
              <div v-for="i in 3" :key="i">
                <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t(`docs.faq.q${i}`) }}</p>
                <p class="mt-1 text-sm leading-relaxed text-gray-600 dark:text-dark-300">{{ t(`docs.faq.a${i}`) }}</p>
              </div>
            </div>
          </section>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, h, defineComponent } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import PublicNav from '@/components/common/PublicNav.vue'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()
const appStore = useAppStore()

const baseUrl = computed(
  () => appStore.cachedPublicSettings?.api_base_url || appStore.apiBaseUrl || window.location.origin
)

interface Step { heading: string; text?: string; code?: string }
interface Tool { id: string; title: string; desc: string; steps: Step[] }

const tools = computed<Tool[]>(() => {
  const b = baseUrl.value
  return [
    {
      id: 'claude-code',
      title: 'Claude Code',
      desc: t('docs.claude.desc'),
      steps: [
        { heading: t('docs.common.install'), code: 'npm install -g @anthropic-ai/claude-code' },
        { heading: t('docs.claude.skipLogin'), text: t('docs.claude.skipLoginText'), code: '// ~/.claude.json\n{\n  "hasCompletedOnboarding": true\n}' },
        {
          heading: t('docs.common.config'),
          text: t('docs.claude.configText'),
          code: `// ~/.claude/settings.json\n{\n  "env": {\n    "ANTHROPIC_AUTH_TOKEN": "YOUR_API_KEY",\n    "ANTHROPIC_BASE_URL": "${b}",\n    "ANTHROPIC_MODEL": "claude-sonnet-4-5",\n    "ANTHROPIC_DEFAULT_HAIKU_MODEL": "claude-haiku-4-5",\n    "ANTHROPIC_DEFAULT_SONNET_MODEL": "claude-sonnet-4-5",\n    "ANTHROPIC_DEFAULT_OPUS_MODEL": "claude-opus-4-5"\n  }\n}`
        },
        { heading: t('docs.common.verify'), text: t('docs.claude.verifyText'), code: 'claude "你好"' }
      ]
    },
    {
      id: 'codex',
      title: 'Codex',
      desc: t('docs.codex.desc'),
      steps: [
        { heading: t('docs.common.install'), code: 'npm install -g @openai/codex' },
        {
          heading: t('docs.common.config'),
          text: t('docs.codex.configText'),
          code: `# ~/.codex/config.toml\nmodel_provider = "one_code"\nmodel = "gpt-5-codex"\n\n[model_providers.one_code]\nname = "one_code"\nbase_url = "${b}/v1"\nenv_key = "OPENAI_API_KEY"\nwire_api = "responses"`
        },
        { heading: t('docs.codex.envHeading'), text: t('docs.codex.envText'), code: 'export OPENAI_API_KEY="YOUR_API_KEY"' },
        { heading: t('docs.common.verify'), text: t('docs.codex.verifyText'), code: 'codex' }
      ]
    },
    {
      id: 'openclaw',
      title: 'OpenClaw',
      desc: t('docs.openclaw.desc'),
      steps: [
        { heading: t('docs.common.install'), text: t('docs.openclaw.installText'), code: 'npm install -g openclaw@latest' },
        {
          heading: t('docs.common.config'),
          text: t('docs.openclaw.configText'),
          code: `// ~/.openclaw/openclaw.json\n{\n  "models": {\n    "mode": "merge",\n    "providers": {\n      "one-code": {\n        "baseUrl": "${b}",\n        "apiKey": "YOUR_API_KEY",\n        "api": "anthropic-messages",\n        "models": [\n          { "id": "claude-sonnet-4-5", "name": "claude-sonnet-4-5" }\n        ]\n      }\n    }\n  },\n  "agents": {\n    "defaults": {\n      "model": { "primary": "one-code/claude-sonnet-4-5" }\n    }\n  },\n  "gateway": { "mode": "local", "auth": { "mode": "none" } }\n}`
        },
        { heading: t('docs.common.verify'), text: t('docs.openclaw.verifyText'), code: 'openclaw gateway restart' }
      ]
    },
    {
      id: 'hermes',
      title: 'Hermes Agent',
      desc: t('docs.hermes.desc'),
      steps: [
        { heading: t('docs.common.install'), code: 'curl -fsSL https://raw.githubusercontent.com/NousResearch/hermes-agent/main/scripts/install.sh | bash' },
        {
          heading: t('docs.common.config'),
          text: t('docs.hermes.configText'),
          code: `hermes config set model.provider custom\nhermes config set model.base_url ${b}\nhermes config set model.api_mode anthropic_messages\nhermes config set model.api_key YOUR_API_KEY\nhermes config set model.default claude-sonnet-4-5`
        },
        { heading: t('docs.common.verify'), text: t('docs.hermes.verifyText'), code: 'hermes chat -q "你好"' }
      ]
    }
  ]
})

const sections = computed(() => [
  { id: 'quick-start', title: t('docs.quickStart.title') },
  ...tools.value.map((tool) => ({ id: tool.id, title: tool.title })),
  { id: 'faq', title: t('docs.faq.title') }
])

const activeId = ref('quick-start')
function scrollTo(id: string) {
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}
function onScroll() {
  for (const s of sections.value) {
    const el = document.getElementById(s.id)
    if (el && el.getBoundingClientRect().top <= 120) activeId.value = s.id
  }
}

onMounted(() => {
  if (!appStore.publicSettingsLoaded) appStore.fetchPublicSettings()
  window.addEventListener('scroll', onScroll, { passive: true })
})
onUnmounted(() => window.removeEventListener('scroll', onScroll))

const CodeBlock = defineComponent({
  props: { code: { type: String, required: true } },
  setup(props) {
    const copied = ref(false)
    const copy = () => {
      navigator.clipboard?.writeText(props.code)
      copied.value = true
      setTimeout(() => (copied.value = false), 1500)
    }
    return () =>
      h('div', { class: 'relative mb-2 rounded-lg bg-gray-900 p-4 dark:bg-dark-800' }, [
        h('button', {
          class: 'absolute right-2 top-2 rounded-md p-1.5 text-gray-400 hover:bg-white/10 hover:text-white',
          onClick: copy,
          title: t('pricing.copy')
        }, [h(Icon, { name: copied.value ? 'check' : 'copy', size: 'xs' })]),
        h('pre', { class: 'overflow-x-auto text-xs leading-relaxed text-gray-100' }, [
          h('code', { class: 'whitespace-pre font-mono' }, props.code)
        ])
      ])
  }
})
</script>
