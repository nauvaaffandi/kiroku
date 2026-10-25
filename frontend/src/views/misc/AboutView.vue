<template>
  <div class="about-page">
    <header class="page-header">
      <h1>{{ $t('about.title') }}</h1>
    </header>

    <div class="about-content">
      <section class="about-section">
        <div class="app-info">
          <div class="app-icon">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M12 2L2 7l10 5 10-5-10-5z"></path>
              <path d="M2 17l10 5 10-5"></path>
              <path d="M2 12l10 5 10-5"></path>
            </svg>
          </div>
          <div class="app-details">
            <h2>Kiroku</h2>
            <p class="app-version">{{ $t('about.version', { version: appVersion }) }}</p>
            <p class="app-build">{{ $t('about.buildDate', { date: buildDate }) }}</p>
            <p class="app-desc">{{ $t('about.desc') }}</p>
          </div>
        </div>
      </section>

      <section class="about-section">
        <h3>{{ $t('about.links') }}</h3>
        <div class="links-grid">
          <a href="https://github.com" target="_blank" rel="noopener" class="link-card">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z"/>
            </svg>
            <span>{{ $t('about.githubRepo') }}</span>
          </a>
          <a href="https://github.com/issues" target="_blank" rel="noopener" class="link-card">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <span>{{ $t('about.reportIssues') }}</span>
          </a>
          <a href="https://opensource.org/licenses/MIT" target="_blank" rel="noopener" class="link-card">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path>
            </svg>
            <span>{{ $t('about.mitLicense') }}</span>
          </a>
          <a href="#" target="_blank" rel="noopener" class="link-card">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
            </svg>
            <span>{{ $t('about.privacyPolicy') }}</span>
          </a>
        </div>
      </section>

      <section class="about-section">
        <h3>{{ $t('about.systemInfo') }}</h3>
        <div class="sys-info">
          <div class="sys-item">
            <span class="sys-label">{{ $t('about.frontend') }}</span>
            <span class="sys-value">Vue 3 + TypeScript + Vite</span>
          </div>
          <div class="sys-item">
            <span class="sys-label">{{ $t('about.uiFramework') }}</span>
            <span class="sys-value">Tailwind CSS v4</span>
          </div>
          <div class="sys-item">
            <span class="sys-label">{{ $t('about.stateManagement') }}</span>
            <span class="sys-value">Pinia</span>
          </div>
          <div class="sys-item">
            <span class="sys-label">{{ $t('about.charts') }}</span>
            <span class="sys-value">ECharts</span>
          </div>
          <div class="sys-item">
            <span class="sys-label">{{ $t('about.backend') }}</span>
            <span class="sys-value">Go + Chi Router</span>
          </div>
          <div class="sys-item">
            <span class="sys-label">{{ $t('about.database') }}</span>
            <span class="sys-value">SQLite (WAL mode)</span>
          </div>
          <div class="sys-item">
            <span class="sys-label">{{ $t('about.auth') }}</span>
            <span class="sys-value">JWT + TOTP 2FA</span>
          </div>
          <div class="sys-item">
            <span class="sys-label">{{ $t('about.architecture') }}</span>
            <span class="sys-value">Monorepo</span>
          </div>
        </div>
      </section>

      <section class="about-section">
        <h3>{{ $t('about.acknowledgments') }}</h3>
        <p class="acknowledgment">
          {{ $t('about.acknowledgmentsDesc') }}
        </p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from '@/locales'

const { t: $t } = useI18n()

const appVersion = '1.0.0'
const buildDate = '2026-10-14'

const getSystemInfo = async () => {
  // System info is static in this build
}

onMounted(() => {
  getSystemInfo()
})
</script>

<style scoped>
.about-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.page-header {
  padding: 1.5rem;
  background-color: var(--color-neutral-100);
  border-bottom: 1px solid var(--color-neutral-200);
}
.page-header h1 { font-size: 1.5rem; font-weight: 600; color: var(--color-neutral-900); }

.about-content {
  flex: 1;
  padding: 1.5rem;
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.about-section {
  background-color: var(--color-neutral-100);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 1.5rem;
}

.about-section h3 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-neutral-900);
  margin-bottom: 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--color-neutral-200);
}

.app-info {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.app-icon {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, var(--color-primary), var(--color-primary-light));
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.app-details h2 {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-neutral-900);
  margin: 0 0 0.25rem;
}

.app-version {
  font-size: 1rem;
  color: var(--color-neutral-600);
  margin: 0 0 0.25rem;
}

.app-build {
  font-size: 0.875rem;
  color: var(--color-neutral-500);
  margin: 0 0 0.5rem;
}

.app-desc {
  font-size: 0.875rem;
  color: var(--color-neutral-500);
  margin: 0;
}

.links-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 1rem;
}

.link-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 1.25rem 1rem;
  background-color: var(--color-neutral-50);
  border: 1px solid var(--color-neutral-200);
  border-radius: var(--radius);
  color: var(--color-neutral-700);
  text-decoration: none;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s ease;
}

.link-card:hover {
  background-color: var(--color-primary);
  border-color: var(--color-primary);
  color: white;
  transform: translateY(-2px);
  box-shadow: var(--shadow);
}

.link-card svg {
  color: var(--color-neutral-400);
  transition: color 0.2s ease;
}

.link-card:hover svg {
  color: white;
}

.sys-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 0.75rem;
}

.sys-item {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
  padding: 0.75rem 1rem;
  background-color: var(--color-neutral-50);
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-neutral-200);
}

.sys-label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-neutral-500);
}

.sys-value {
  font-size: 0.875rem;
  color: var(--color-neutral-900);
}

.acknowledgment {
  font-size: 0.875rem;
  color: var(--color-neutral-600);
  line-height: 1.6;
}
</style>