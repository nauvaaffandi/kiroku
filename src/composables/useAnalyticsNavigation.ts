import { ref, computed, watch, onMounted, nextTick, type Ref } from 'vue'
import { useRoute } from 'vue-router'

export function useAnalyticsNavigation(
  subNavRef: Ref<HTMLElement | null>
) {
  const route = useRoute()
  const isAnalyticsOpen = ref(false)
  const analyticsHeight = ref(0)
  const isAnalyticsActive = computed(() => route.path.startsWith('/analytics'))

  isAnalyticsOpen.value = isAnalyticsActive.value

  const updateHeight = () => {
    if (subNavRef.value) {
      analyticsHeight.value = subNavRef.value.scrollHeight
    }
  }

  watch(isAnalyticsActive, (val) => {
    if (val) {
      isAnalyticsOpen.value = true
    }
  })

  watch(isAnalyticsOpen, () => {
    nextTick(() => {
      updateHeight()
    })
  })

  onMounted(() => {
    if (isAnalyticsActive.value) {
      isAnalyticsOpen.value = true
    }
    nextTick(() => {
      updateHeight()
    })
  })

  const toggleAnalytics = () => {
    isAnalyticsOpen.value = !isAnalyticsOpen.value
  }

  return {
    isAnalyticsOpen,
    analyticsHeight,
    isAnalyticsActive,
    toggleAnalytics
  }
}