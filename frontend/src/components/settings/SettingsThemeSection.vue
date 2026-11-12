<template>
  <SettingSection :title="t('settings.theme')">
    <SettingItem :label="t('settings.theme')" :description="t('settings.themeSystem')">
      <BaseSelect :model-value="theme" :options="themeOptions" @update:model-value="(v: any) => emit('themeChange', v)" />
    </SettingItem>
  </SettingSection>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { BaseSelect } from '@/components/ui'
import { SettingItem, SettingSection } from '@/widgets/settings'

const props = defineProps<{
  theme: any
  themeOptions: { value: string; label: string }[]
}>()
const emit = defineEmits<{ themeChange: [string] }>()
const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";
</style>
