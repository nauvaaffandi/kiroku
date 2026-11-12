<template>
  <SettingSection :title="t('settings.telegramTitle')">
    <SettingItem :label="t('settings.telegram')" :description="isTelegramLinked ? t('settings.telegramLinked') : t('settings.telegramDesc')">
      <template #default>
        <BaseButton :variant="isTelegramLinked ? 'secondary' : 'primary'" @click="emit('openRegister')" :loading="telegramLoading">{{ isTelegramLinked ? 'Linked (Re-link)' : t('settings.telegramRegister') }}</BaseButton>
      </template>
    </SettingItem>
  </SettingSection>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales'
import { BaseButton } from '@/components/ui'
import { SettingItem, SettingSection } from '@/widgets/settings'

const props = defineProps<{ isTelegramLinked: boolean; telegramLoading: boolean }>()
const emit = defineEmits<{ openRegister: [] }>()
const { t } = useI18n()
</script>

<style scoped>
@reference "tailwindcss";
</style>
