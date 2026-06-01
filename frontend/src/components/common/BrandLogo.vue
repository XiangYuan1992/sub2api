<template>
  <img
    v-if="src"
    :src="src"
    :alt="alt"
    class="h-full w-full object-contain"
  />
  <svg
    v-else
    class="h-full w-full"
    viewBox="0 0 512 512"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    :aria-label="alt"
    role="img"
  >
    <title>{{ alt }}</title>
    <g
      stroke="currentColor"
      :stroke-width="strokeWidth"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <rect x="88" y="72" width="336" height="368" rx="88" />
      <path d="M164 154H348" />
      <path d="M196 220H316" />
      <path d="M196 220V330" />
      <path d="M316 220V330" />
      <path d="M170 330H342" />
      <path d="M170 330L256 256L342 330" />
    </g>
  </svg>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    src?: string
    alt?: string
    strokeWidth?: number
  }>(),
  {
    src: '',
    alt: '壹站 Logo',
    strokeWidth: 24,
  }
)
</script>
