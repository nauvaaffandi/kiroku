<template>
  <section class="not-found">
    <div class="character-wrap">
      <div class="character">{{ t('notFound.character') }}</div>
      <div class="float-deco float-deco--1"></div>
      <div class="float-deco float-deco--2"></div>
      <div class="float-deco float-deco--3"></div>
    </div>

    <h1 class="not-found-title">{{ t('notFound.title') }}</h1>

    <p class="not-found-message">{{ t('notFound.message') }}</p>
    <p class="not-found-subtitle">{{ t('notFound.subtitle') }}</p>

    <button class="not-found-btn" @click="goHome">
      <span>{{ t('notFound.home') }}</span>
      <span class="btn-arrow">&#x2192;</span>
    </button>
  </section>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { useRouter } from 'vue-router'

const { t } = useI18n()
const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
@reference "tailwindcss";

.not-found {
  @apply m-auto flex max-w-md flex-col items-center justify-center gap-6 px-4 py-12 min-h-[60vh] text-center;
}

.character-wrap {
  @apply relative inline-block;
}

.character {
  @apply text-6xl;
  animation: float 4s ease-in-out infinite;
}

.float-deco {
  @apply absolute inline-block text-2xl;
  animation: float 5s ease-in-out infinite;
}

.float-deco--1 {
  @apply top-1/4 left-1/4 text-yellow-400;
  animation-delay: 0s;
  content: '✨';
}

.float-deco--2 {
  @apply bottom-1/3 right-1/3 text-emerald-400;
  animation-delay: 1.5s;
  content: '🌙';
}

.float-deco--3 {
  @apply top-1/2 left-2/3 text-rose-400;
  animation-delay: 3s;
  content: '⭐';
}

.not-found-title {
  @apply text-6xl font-bold tracking-tight text-[var(--text-primary)] mb-2;
  letter-spacing: -0.05em;
}

.not-found-message {
  @apply text-xl font-semibold text-[var(--text-primary)];
}

.not-found-subtitle {
  @apply text-base text-[var(--text-secondary)] max-w-sm mt-2 leading-relaxed;
}

.not-found-btn {
  @apply inline-flex items-center justify-center gap-2 px-6 py-3 bg-[var(--bg-accent)] text-[var(--text-accent)] font-semibold rounded-[var(--radius-lg)] transition-all duration-200 hover:scale-[1.02] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--accent)] focus-visible:ring-offset-2 mt-6;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-12px);
  }
}

@media (max-width: 480px) {
  .not-found-title {
    @apply text-4xl;
  }

  .character {
    @apply text-4xl;
  }
}
</style>
