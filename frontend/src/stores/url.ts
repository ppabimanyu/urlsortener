import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'
import type {
  URLItem,
  URLListResponse,
  OverviewAnalytics,
  URLAnalyticsDetail,
} from '@/types'

export const useUrlStore = defineStore('url', () => {
  const urls = ref<URLItem[]>([])
  const total = ref(0)
  const page = ref(1)
  const limit = ref(10)
  const search = ref('')
  const activeFilter = ref<boolean | null>(null)
  const loading = ref(false)

  const overview = ref<OverviewAnalytics | null>(null)
  const overviewLoading = ref(false)

  const selectedUrlAnalytics = ref<URLAnalyticsDetail | null>(null)
  const analyticsLoading = ref(false)

  async function fetchUrls(currentPage = 1, searchQuery = '', statusFilter: boolean | null = null) {
    loading.value = true
    try {
      page.value = currentPage
      search.value = searchQuery
      activeFilter.value = statusFilter

      const params: Record<string, any> = {
        page: currentPage,
        limit: limit.value,
      }
      if (searchQuery) params.search = searchQuery
      if (statusFilter !== null) params.is_active = statusFilter

      const res = await api.get<URLListResponse>('/urls', { params })
      urls.value = res.data.data || []
      total.value = res.data.total
    } catch (err) {
      console.error('Failed to fetch URLs:', err)
    } finally {
      loading.value = false
    }
  }

  async function createUrl(payload: {
    original_url: string
    custom_slug?: string
    title?: string
    password?: string
    expires_at?: string | null
    click_limit?: number | null
  }) {
    const res = await api.post<{ message: string; data: URLItem }>('/urls', payload)
    await fetchUrls(1, search.value, activeFilter.value)
    return res.data.data
  }

  async function updateUrl(
    id: string,
    payload: {
      original_url?: string
      custom_slug?: string
      title?: string
      password?: string
      expires_at?: string | null
      click_limit?: number | null
      is_active?: boolean
    }
  ) {
    const res = await api.put<{ message: string; data: URLItem }>(`/urls/${id}`, payload)
    const index = urls.value.findIndex((u) => u.id === id)
    if (index !== -1) {
      urls.value[index] = res.data.data
    }
    return res.data.data
  }

  async function deleteUrl(id: string) {
    await api.delete(`/urls/${id}`)
    urls.value = urls.value.filter((u) => u.id !== id)
    total.value = Math.max(0, total.value - 1)
  }

  async function toggleUrlStatus(id: string) {
    const index = urls.value.findIndex((u) => u.id === id)
    const prevStatus = index !== -1 ? urls.value[index].is_active : null
    if (index !== -1) {
      urls.value[index].is_active = !urls.value[index].is_active
    }
    try {
      const res = await api.patch<{ message: string; data: URLItem }>(`/urls/${id}/toggle`)
      if (index !== -1) {
        urls.value[index] = res.data.data
      }
      return res.data.data
    } catch (err) {
      if (index !== -1 && prevStatus !== null) {
        urls.value[index].is_active = prevStatus
      }
      throw err
    }
  }

  async function fetchOverview() {
    overviewLoading.value = true
    try {
      const res = await api.get<{ data: OverviewAnalytics }>('/analytics/overview')
      overview.value = res.data.data
    } catch (err) {
      console.error('Failed to fetch overview analytics:', err)
    } finally {
      overviewLoading.value = false
    }
  }

  async function fetchUrlAnalytics(id: string, days = 14) {
    analyticsLoading.value = true
    try {
      const res = await api.get<{ data: URLAnalyticsDetail }>(`/analytics/urls/${id}`, {
        params: { days },
      })
      selectedUrlAnalytics.value = res.data.data
    } catch (err) {
      console.error('Failed to fetch URL analytics:', err)
    } finally {
      analyticsLoading.value = false
    }
  }

  async function exportCsv(id: string, shortCode: string) {
    const res = await api.get(`/analytics/urls/${id}/export`, {
      responseType: 'blob',
    })
    const blob = new Blob([res.data], { type: 'text/csv' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `analytics_${shortCode}.csv`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  }

  return {
    urls,
    total,
    page,
    limit,
    search,
    activeFilter,
    loading,
    overview,
    overviewLoading,
    selectedUrlAnalytics,
    analyticsLoading,
    fetchUrls,
    createUrl,
    updateUrl,
    deleteUrl,
    toggleUrlStatus,
    fetchOverview,
    fetchUrlAnalytics,
    exportCsv,
  }
})
