import { ref, computed } from 'vue'

export function usePagination(initialPage = 1, initialLimit = 20) {
  const page = ref(initialPage)
  const limit = ref(initialLimit)
  const total = ref(0)

  const totalPages = computed(() => Math.ceil(total.value / limit.value))
  const hasNextPage = computed(() => page.value < totalPages.value)
  const hasPrevPage = computed(() => page.value > 1)

  const setPage = (newPage: number) => {
    page.value = Math.max(1, Math.min(newPage, totalPages.value))
  }

  const nextPage = () => {
    if (hasNextPage.value) page.value++
  }

  const prevPage = () => {
    if (hasPrevPage.value) page.value--
  }

  const updateTotal = (newTotal: number) => {
    total.value = newTotal
  }

  const reset = () => {
    page.value = initialPage
    total.value = 0
  }

  return {
    page,
    limit,
    total,
    totalPages,
    hasNextPage,
    hasPrevPage,
    setPage,
    nextPage,
    prevPage,
    updateTotal,
    reset
  }
}