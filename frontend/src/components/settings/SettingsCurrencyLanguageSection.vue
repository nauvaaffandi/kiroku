<template>
  <SettingSection :title="t('settings.currency') + ' & ' + t('settings.language')">
    <SettingItem :label="t('settings.currency')" :description="t('settings.currencyDesc')">
      <div class="flex flex-col gap-2 w-full">
        <BaseSelect :model-value="baseCurrency" :options="currencyOptions" @update:model-value="(v: any) => emit('baseCurrencyChange', v)" />
        <div class="flex gap-2 mt-2">
          <input :model-value="newCurrencyCode" @update:model-value="emit('update:newCurrencyCode', $event)" placeholder="Add Currency (e.g. IDR)" class="p-2 border border-[var(--border-primary)] rounded-lg bg-[var(--bg-primary)] text-[var(--text-primary)] text-sm flex-1" maxlength="5" />
          <BaseButton variant="secondary" @click="handleAddCurrency">Add</BaseButton>
        </div>
      </div>
    </SettingItem>
    <SettingItem :label="t('settings.language')" :description="t('settings.languageDesc')">
      <LanguageSwitcher />
    </SettingItem>
    <SettingItem :label="t('settings.date_format')" :description="t('settings.dateFormatDesc')">
      <BaseSelect :model-value="settings.date_format" :options="dateFormatOptions" @update:model-value="(v: any) => emit('dateFormatChange', v)" />
    </SettingItem>
    <SettingItem :label="t('settings.number_format')" :description="t('settings.numberFormatDesc')">
      <BaseSelect :model-value="settings.number_format" :options="numberFormatOptions" @update:model-value="(v: any) => emit('numberFormatChange', v)" />
    </SettingItem>
    <SettingItem :label="t('settings.week_start')" :description="t('settings.weekStartDesc')">
      <BaseSelect :model-value="weekStartStr" :options="weekStartOptions" @update:model-value="(v: any) => emit('weekStartChange', v)" />
    </SettingItem>
  </SettingSection>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { BaseSelect, BaseButton } from '@/components/ui'
import LanguageSwitcher from '@/components/forms/LanguageSwitcher.vue'
import { SettingItem, SettingSection } from '@/widgets/settings'

const props = defineProps<{
  settings: any
  currencyOptions: { value: string; label: string }[]
  dateFormatOptions: { value: string; label: string }[]
  numberFormatOptions: { value: string; label: string }[]
  weekStartOptions: { value: string; label: string }[]
  weekStartStr: string
  newCurrencyCode: string
  baseCurrency: string
}>()
const emit = defineEmits<{
  addCurrency: []
  'update:newCurrencyCode': [string]
  dateFormatChange: [string]
  numberFormatChange: [string]
  weekStartChange: [string]
  baseCurrencyChange: [string]
}>()
const { t } = useI18n()
const handleAddCurrency = () => emit('addCurrency')
</script>
