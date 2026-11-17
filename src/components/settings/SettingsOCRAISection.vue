<template>
  <SettingSection :title="t('settings.ocrTitle')">
    <SettingItem :label="t('settings.ocrToggle')">
      <input
        type="checkbox"
        v-model="ocrEnabled"
        id="ocr-toggle"
        :aria-label="t('settings.ocrToggle')"
        class="ocr-toggle"
      />
    </SettingItem>
    <SettingItem :label="t('settings.ocrApiKey')">
      <BaseInput
        v-model="ocrApiKey"
        type="password"
        :placeholder="t('settings.ocrApiKey')"
        :disabled="!ocrEnabled"
        :label="t('settings.ocrApiKey')"
        id="ocr-api-key"
      />
    </SettingItem>
  </SettingSection>

  <SettingSection :title="t('settings.aiTitle')">
    <SettingItem :label="t('settings.aiToggle')">
      <input
        type="checkbox"
        v-model="aiProviderEnabled"
        id="ai-toggle"
        :aria-label="t('settings.aiToggle')"
        class="ai-toggle"
      />
    </SettingItem>
    <SettingItem :label="t('settings.aiProviderUrl')">
      <BaseInput
        v-model="aiProviderUrl"
        type="url"
        :placeholder="t('settings.aiProviderUrl')"
        :disabled="!aiProviderEnabled"
        :label="t('settings.aiProviderUrl')"
        id="ai-provider-url"
      />
    </SettingItem>
    <SettingItem :label="t('settings.aiModel')">
      <BaseInput
        v-model="aiModel"
        type="text"
        :placeholder="t('settings.aiModel')"
        :disabled="!aiProviderEnabled"
        :label="t('settings.aiModel')"
        id="ai-model"
      />
    </SettingItem>
    <SettingItem :label="t('settings.aiApiKey')">
      <BaseInput
        v-model="aiApiKey"
        type="password"
        :placeholder="t('settings.aiApiKey')"
        :disabled="!aiProviderEnabled"
        :label="t('settings.aiApiKey')"
        id="ai-api-key"
      />
    </SettingItem>
  </SettingSection>

  <div class="ocr-ai-warning">{{ t('settings.aiLocalOnlyWarning') }}</div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from '@/locales'
import { useSettingsStore } from '@/stores/settings'
import { BaseInput } from '@/components/ui'
import { SettingItem, SettingSection } from '@/widgets/settings'

const store = useSettingsStore()
const { t } = useI18n()

const ocrEnabled = computed({
  get: () => store.isOcrEnabled,
  set: (value: boolean) => {
    console.log('[OCR] toggle changed')
    store.setSetting('ocr_enabled', value)
  }
})

const aiProviderEnabled = computed({
  get: () => store.isAiProviderEnabled,
  set: (value: boolean) => {
    console.log('[AI] toggle changed')
    store.setSetting('ai_provider_enabled', value)
  }
})

const ocrApiKey = computed({
  get: () => store.settings.ocr_api_key ?? '',
  set: (value: string) => {
    console.log('[OCR] api key updated')
    store.setSetting('ocr_api_key', value)
  }
})

const aiProviderUrl = computed({
  get: () => store.settings.ai_provider_url ?? '',
  set: (value: string) => {
    console.log('[AI] provider url updated')
    store.setSetting('ai_provider_url', value)
  }
})

const aiModel = computed({
  get: () => store.settings.ai_model ?? '',
  set: (value: string) => {
    console.log('[AI] model updated')
    store.setSetting('ai_model', value)
  }
})

const aiApiKey = computed({
  get: () => store.settings.ai_api_key ?? '',
  set: (value: string) => {
    console.log('[AI] api key updated')
    store.setSetting('ai_api_key', value)
  }
})
</script>

<style scoped>
@reference "tailwindcss";

.ocr-toggle,
.ai-toggle {
  @apply w-[44px] h-[44px] rounded-full accent-[var(--accent)];
}

.ocr-ai-warning {
  @apply mt-6 p-4 rounded-lg border text-sm bg-[var(--bg-hover)] border-[var(--border-primary)] text-[var(--text-muted)];
}
</style>
