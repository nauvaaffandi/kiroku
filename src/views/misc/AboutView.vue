<template>
  <div class="about-page">
    <header class="page-header">
      <h1>{{ t('about.title') }}</h1>
    </header>

    <div class="about-content">
      <section class="about-section">
        <div class="app-info">
          <div class="app-icon">
            <AppIcon />
          </div>
          <div class="app-details">
            <h2>{{ t('about.appName') }}</h2>
            <p class="app-version">{{ t('about.version', { version: appVersion }) }}</p>
            <p class="app-build">{{ t('about.buildDate', { date: buildDate }) }}</p>
            <p class="app-desc">{{ t('about.desc') }}</p>
          </div>
        </div>
      </section>

      <section class="about-section">
        <h3>{{ t('about.links') }}</h3>
        <div class="links-grid">
          <a
            href="https://github.com"
            target="_blank"
            rel="noopener"
            class="link-card"
            :aria-label="t('about.githubRepo')"
          >
            <GitBranch :size="20" stroke-width="1.75" aria-hidden="true" />
            <span>{{ t('about.githubRepo') }}</span>
          </a>
          <a
            href="https://github.com/issues"
            target="_blank"
            rel="noopener"
            class="link-card"
            :aria-label="t('about.reportIssues')"
          >
            <FileText :size="20" stroke-width="1.75" aria-hidden="true" />
            <span>{{ t('about.reportIssues') }}</span>
          </a>
          <a
            href="https://opensource.org/licenses/MIT"
            target="_blank"
            rel="noopener"
            class="link-card"
            :aria-label="t('about.mitLicense')"
          >
            <Shield :size="20" stroke-width="1.75" aria-hidden="true" />
            <span>{{ t('about.mitLicense') }}</span>
          </a>
          <a
            href="#"
            target="_blank"
            rel="noopener"
            class="link-card"
            :aria-label="t('about.privacyPolicy')"
          >
            <Lock :size="20" stroke-width="1.75" aria-hidden="true" />
            <span>{{ t('about.privacyPolicy') }}</span>
          </a>
        </div>
      </section>

      <section class="about-section">
        <h3>{{ t('about.systemInfo') }}</h3>
        <div class="sys-info">
          <div class="sys-item">
            <Code class="sys-icon" aria-hidden="true" />
            <span class="sys-label">{{ t('about.frontend') }}</span>
            <span class="sys-value">{{ t('about.stackFrontend') }}</span>
          </div>
          <div class="sys-item">
            <Layout class="sys-icon" aria-hidden="true" />
            <span class="sys-label">{{ t('about.uiFramework') }}</span>
            <span class="sys-value">{{ t('about.stackUi') }}</span>
          </div>
          <div class="sys-item">
            <Store class="sys-icon" aria-hidden="true" />
            <span class="sys-label">{{ t('about.stateManagement') }}</span>
            <span class="sys-value">{{ t('about.stackState') }}</span>
          </div>
          <div class="sys-item">
            <PieChart class="sys-icon" aria-hidden="true" />
            <span class="sys-label">{{ t('about.charts') }}</span>
            <span class="sys-value">{{ t('about.stackChart') }}</span>
          </div>
          <div class="sys-item">
            <Database class="sys-icon" aria-hidden="true" />
            <span class="sys-label">{{ t('about.backend') }}</span>
            <span class="sys-value">{{ t('about.stackBackend') }}</span>
          </div>
          <div class="sys-item">
            <HardDrive class="sys-icon" aria-hidden="true" />
            <span class="sys-label">{{ t('about.database') }}</span>
            <span class="sys-value">{{ t('about.stackDatabase') }}</span>
          </div>
          <div class="sys-item">
            <Shield class="sys-icon" aria-hidden="true" />
            <span class="sys-label">{{ t('about.auth') }}</span>
            <span class="sys-value">{{ t('about.stackAuth') }}</span>
          </div>
          <div class="sys-item">
            <Layers class="sys-icon" aria-hidden="true" />
            <span class="sys-label">{{ t('about.architecture') }}</span>
            <span class="sys-value">{{ t('about.stackArch') }}</span>
          </div>
        </div>
      </section>

      <section class="about-section">
        <h3>{{ t('about.acknowledgments') }}</h3>
        <p class="acknowledgment">{{ t('about.acknowledgmentsDesc') }}</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from '@/locales'
import AppIcon from '@/components/icons/AppIcon.vue'
import {
  GitBranch,
  FileText,
  Shield,
  Lock,
  Code,
  Layout,
  Store,
  PieChart,
  Database,
  HardDrive,
  Layers
} from '@lucide/vue'

const { t } = useI18n()

const appVersion = '1.0.0'
const buildDate = '2026-10-14'

onMounted(() => {
  console.log('[AboutView] mounted')
})
</script>

<style scoped>
@reference "tailwindcss";

.about-page {
  @apply min-h-screen flex flex-col w-full max-w-[480px] mx-auto;
}

.page-header {
  @apply px-6 py-6 flex items-center;
  border-bottom: 1px solid var(--border-primary);
  background-color: var(--bg-card);
}

.page-header h1 {
  @apply text-[20px] font-semibold m-0;
  color: var(--text-primary);
}

.about-content {
  @apply flex-1 p-6 flex flex-col gap-6 w-full;
}

.about-section {
  @apply rounded-[14px] border p-6 flex flex-col gap-4;
  background-color: var(--bg-card);
  border-color: var(--border-primary);
}

.about-section h3 {
  @apply text-[16px] font-semibold mb-3 pb-3 border-b m-0;
  border-color: var(--border-primary);
  color: var(--text-primary);
}

.app-info {
  @apply flex items-center gap-4;
}

.app-icon {
  @apply shrink-0;
}

.app-details h2 {
  @apply text-[20px] font-semibold m-[0_0_0.5rem_0];
  color: var(--text-primary);
}

.app-version {
  @apply text-[14px] m-[0_0_0.25rem_0];
  color: var(--text-secondary);
}

.app-build {
  @apply text-sm m-[0_0_0.5rem_0];
  color: var(--text-muted);
}

.app-desc {
  @apply text-[14px] m-0;
  color: var(--text-muted);
}

.links-grid {
  @apply grid grid-cols-[repeat(auto-fit,minmax(140px,1fr))] gap-3;
}

.link-card {
  @apply flex flex-col items-center gap-2 px-4 py-3 rounded-[10px] border text-[14px] font-medium transition-colors duration-120 ease-out;
  background-color: transparent;
  border-color: var(--border-primary);
  color: var(--text-secondary);
  text-decoration: none;
}

.link-card:hover {
  @apply bg-[var(--bg-hover)];
  color: var(--text-primary);
  border-color: var(--border-primary);
}

.link-card svg {
  color: var(--text-muted);
}

.link-card:hover svg {
  color: var(--text-primary);
}

.link-card:focus-visible {
  @apply ring-2 ring-offset-2 rounded-[10px];
  ring-offset-color: var(--bg-card);
  ring-color: var(--accent);
}

.sys-info {
  @apply grid grid-cols-[repeat(auto-fit,minmax(160px,1fr))] gap-3;
}

.sys-item {
  @apply flex flex-col gap-1 px-4 py-3 rounded-[6px] border;
  background-color: var(--bg-hover);
  border-color: var(--border-primary);
}

.sys-icon {
  @apply mt-[0.1rem];
  color: var(--text-muted);
}

.sys-label {
  @apply text-[0.75rem] font-semibold text-[var(--text-muted)];
}

.sys-value {
  @apply text-sm text-[var(--text-primary)];
}

.acknowledgment {
  @apply text-sm leading-6 text-[var(--text-muted)];
}
</style>
