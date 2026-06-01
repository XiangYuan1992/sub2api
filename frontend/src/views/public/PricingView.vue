<template>
  <div class="relative flex min-h-screen flex-col overflow-hidden bg-white dark:bg-dark-950">
    <PublicNav />

    <main class="relative z-10 flex-1 px-6 py-10">
      <div class="mx-auto max-w-6xl">
        <h1 class="mb-2 text-4xl font-semibold tracking-tight text-gray-900 dark:text-white">{{ t('pricing.title') }}</h1>
        <p class="mb-8 text-gray-600 dark:text-dark-300">{{ t('pricing.subtitle') }}</p>

        <div v-if="loading" class="py-20 text-center">
          <Icon name="refresh" size="xl" class="inline-block animate-spin text-gray-400" />
        </div>

        <div v-else-if="!data?.enabled || platforms.length === 0" class="rounded-lg border border-gray-200 p-12 text-center text-gray-500 dark:border-dark-700 dark:text-dark-400">
          {{ data && !data.enabled ? t('pricing.disabled') : t('pricing.empty') }}
        </div>

        <template v-else>
          <div class="mb-6 flex gap-6 border-b border-gray-200 dark:border-dark-700">
            <button
              v-for="p in platforms"
              :key="p.platform"
              @click="activePlatform = p.platform"
              :class="[
                'flex items-center gap-2 border-b-2 px-1 pb-3 text-sm font-medium transition-colors -mb-px',
                activePlatform === p.platform
                  ? 'border-gray-900 text-gray-900 dark:border-white dark:text-white'
                  : 'border-transparent text-gray-500 hover:text-gray-800 dark:text-dark-400 dark:hover:text-dark-200'
              ]"
            >
              {{ platformLabel(p.platform) }}
            </button>
          </div>

          <div class="mb-6 flex flex-wrap items-center gap-x-6 gap-y-2 rounded-lg border border-gray-200 px-4 py-3 text-sm dark:border-dark-700">
            <span class="font-medium text-gray-900 dark:text-white">{{ t('pricing.rulesTitle') }}</span>
            <span class="text-gray-500 dark:text-dark-400">{{ t('pricing.rulesRate', { rate: cnyRate }) }}</span>
            <span class="text-gray-500 dark:text-dark-400">{{ t('pricing.rulesSiteRate', { rate: siteRate.toFixed(2) }) }}</span>
            <span class="text-gray-500 dark:text-dark-400">{{ t('pricing.rulesFormula') }}</span>
          </div>

          <div>
            <div class="mb-4 flex items-center justify-between">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('pricing.listTitle') }}</h2>
              <div class="flex rounded-lg border border-gray-200 p-0.5 dark:border-dark-700">
                <button @click="showOfficial = false" :class="['rounded-md px-3 py-1 text-xs font-medium transition-colors', !showOfficial ? 'bg-gray-900 text-white dark:bg-white dark:text-gray-900' : 'text-gray-500 dark:text-dark-400']">{{ t('pricing.groupPrice') }}</button>
                <button @click="showOfficial = true" :class="['rounded-md px-3 py-1 text-xs font-medium transition-colors', showOfficial ? 'bg-gray-900 text-white dark:bg-white dark:text-gray-900' : 'text-gray-500 dark:text-dark-400']">{{ t('pricing.officialPrice') }}</button>
              </div>
            </div>

            <div class="mb-4 flex flex-wrap gap-3">
              <button
                v-for="g in currentGroups"
                :key="g.id"
                @click="activeGroupId = g.id"
                :class="[
                  'rounded-lg border px-5 py-3 text-left transition-colors',
                  activeGroupId === g.id
                    ? 'border-gray-900 dark:border-white'
                    : 'border-gray-200 hover:border-gray-300 dark:border-dark-700 dark:hover:border-dark-600'
                ]"
              >
                <div class="flex items-center gap-2">
                  <span class="font-semibold text-gray-900 dark:text-white">{{ g.name }}</span>
                  <span
                    v-if="groupSavingPercent(g) > 0"
                    class="rounded-full bg-gray-100 px-2 py-0.5 text-[10px] font-medium text-gray-600 dark:bg-dark-700 dark:text-dark-200"
                  >{{ discountLabelFromSaving(groupSavingPercent(g)) }}</span>
                </div>
                <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ g.rate_multiplier }}x</div>
              </button>
            </div>

            <div v-if="activeGroup?.description" class="mb-4 rounded-lg border border-gray-200 px-4 py-3 text-sm text-gray-600 dark:border-dark-700 dark:text-dark-300">
              <span class="font-medium text-gray-900 dark:text-white">{{ t('pricing.groupIntro') }}：</span>{{ activeGroup.description }}
            </div>

            <div class="overflow-x-auto">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-gray-100 text-xs font-medium uppercase tracking-wide text-gray-500 dark:border-dark-700 dark:text-dark-400">
                    <th class="px-4 py-3 text-left">{{ t('pricing.colModel') }}</th>
                    <th class="px-4 py-3 text-left">{{ t('pricing.colInput') }}</th>
                    <th class="px-4 py-3 text-left">{{ t('pricing.colOutput') }}</th>
                    <th class="px-4 py-3 text-left">{{ t('pricing.colCacheRead') }}</th>
                    <th v-if="!showOfficial" class="px-4 py-3 text-left">{{ t('pricing.colSaving') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="m in sortedModels" :key="m.name" class="border-b border-gray-50 last:border-0 dark:border-dark-800">
                    <td class="px-4 py-4">
                      <button class="inline-flex items-center gap-1.5 font-medium text-gray-900 hover:text-gray-600 dark:text-white dark:hover:text-dark-300" @click="copyModel(m.name)">
                        {{ m.name }}
                        <Icon name="copy" size="xs" class="text-gray-400" />
                      </button>
                    </td>
                    <td class="px-4 py-4"><PriceCell :perTokenUsd="m.input_price" /></td>
                    <td class="px-4 py-4"><PriceCell :perTokenUsd="m.output_price" /></td>
                    <td class="px-4 py-4"><PriceCell :perTokenUsd="m.cache_read_price" /></td>
                    <td v-if="!showOfficial" class="px-4 py-4">
                      <span v-if="activeGroup && savingPercent(rechargeMultiplier, activeGroup.rate_multiplier, cnyRate) > 0" class="text-xs font-medium text-green-600 dark:text-green-400">
                        {{ t('pricing.saving', { percent: savingPercent(rechargeMultiplier, activeGroup.rate_multiplier, cnyRate) }) }}
                      </span>
                      <span v-else class="text-gray-400">-</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </template>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, h, defineComponent } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import PublicNav from '@/components/common/PublicNav.vue'
import Icon from '@/components/icons/Icon.vue'
import { getPricing, type PricingResponse, type PricingGroup } from '@/api/pricing'
import {
  officialCny,
  officialUsd,
  groupCny,
  siteCnyRate,
  discountLabelFromSaving,
  savingPercent,
  sortModelsByVersionDesc
} from './pricingCalc'

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

const loading = ref(true)
const data = ref<PricingResponse | null>(null)
const activePlatform = ref('')
const activeGroupId = ref<number | null>(null)
const showOfficial = ref(false)

const cnyRate = computed(() => data.value?.cny_rate ?? 7)
const rechargeMultiplier = computed(() => data.value?.recharge_multiplier ?? 1)
const siteRate = computed(() => siteCnyRate(rechargeMultiplier.value))
const platforms = computed(() => data.value?.platforms ?? [])
const currentGroups = computed<PricingGroup[]>(
  () => platforms.value.find((p) => p.platform === activePlatform.value)?.groups ?? []
)
const activeGroup = computed<PricingGroup | undefined>(
  () => currentGroups.value.find((g) => g.id === activeGroupId.value)
)
const sortedModels = computed(() => sortModelsByVersionDesc(activeGroup.value?.models ?? []))

const platformLabels: Record<string, string> = {
  anthropic: 'Claude Code',
  openai: 'Codex',
  gemini: 'Gemini'
}
function platformLabel(p: string): string {
  return platformLabels[p] || p
}

function groupSavingPercent(g: PricingGroup): number {
  return savingPercent(rechargeMultiplier.value, g.rate_multiplier, cnyRate.value)
}

const PriceCell = defineComponent({
  props: { perTokenUsd: { type: Number as () => number | null, default: null } },
  setup(props) {
    return () => {
      const rate = cnyRate.value
      const mult = activeGroup.value?.rate_multiplier ?? 1
      const fmtCny = (v: number | null) => (v == null ? '-' : `¥${v.toFixed(2)}`)
      const fmtUsd = (v: number | null) => (v == null ? '' : `$${v.toFixed(2)}`)
      const perM = h('span', { class: 'ml-1 text-xs text-gray-400' }, t('pricing.perM'))

      if (showOfficial.value) {
        const cny = officialCny(props.perTokenUsd, rate)
        const usd = officialUsd(props.perTokenUsd)
        return h('div', {}, [
          h('div', { class: 'font-medium text-gray-900 dark:text-white' }, [fmtCny(cny), perM]),
          usd != null
            ? h('div', { class: 'text-xs text-gray-400' }, fmtUsd(usd))
            : null
        ])
      }

      const gpCny = groupCny(props.perTokenUsd, rechargeMultiplier.value, mult)
      const officialUsdVal = officialUsd(props.perTokenUsd)
      const officialCnyVal = officialCny(props.perTokenUsd, rate)
      return h('div', {}, [
        h('div', { class: 'font-semibold text-primary-600 dark:text-primary-400' }, [fmtCny(gpCny), perM]),
        officialUsdVal != null
          ? h('div', { class: 'text-xs text-gray-400 line-through' }, `${t('pricing.official')} ${fmtUsd(officialUsdVal)} ${fmtCny(officialCnyVal)}`)
          : null
      ])
    }
  }
})

function copyModel(name: string) {
  navigator.clipboard?.writeText(name)
}
function syncDefaults() {
  if (platforms.value.length > 0) {
    activePlatform.value = platforms.value[0].platform
    activeGroupId.value = currentGroups.value[0]?.id ?? null
  }
}

watch(activePlatform, () => {
  activeGroupId.value = currentGroups.value[0]?.id ?? null
})

onMounted(async () => {
  if (!appStore.publicSettingsLoaded) appStore.fetchPublicSettings()
  authStore.checkAuth()
  try {
    data.value = await getPricing()
    syncDefaults()
  } finally {
    loading.value = false
  }
})
</script>
