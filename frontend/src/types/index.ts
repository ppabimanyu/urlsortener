export interface User {
  id: string
  name: string
  email: string
  created_at: string
}

export interface URLItem {
  id: string
  user_id: string
  original_url: string
  short_code: string
  title: string
  has_password?: boolean
  expires_at?: string | null
  click_limit?: number | null
  click_count: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface DateClickStat {
  date: string
  count: number
}

export interface GroupCountStat {
  name: string
  count: number
}

export interface URLClick {
  id: number
  url_id: string
  clicked_at: string
  ip_address: string
  user_agent: string
  device_type: string
  os: string
  browser: string
  referrer: string
  country: string
  country_code: string
  city: string
}

export interface OverviewAnalytics {
  total_links: number
  total_clicks: number
  active_links: number
  timeline: DateClickStat[]
  recent_clicks: URLClick[]
}

export interface URLAnalyticsDetail {
  url: URLItem
  timeline: DateClickStat[]
  devices: GroupCountStat[]
  os: GroupCountStat[]
  browsers: GroupCountStat[]
  referrers: GroupCountStat[]
  countries: GroupCountStat[]
  cities: GroupCountStat[]
}

export interface AuthResponse {
  message: string
  token: string
  user: User
}

export interface URLListResponse {
  data: URLItem[]
  total: number
  page: number
  limit: number
}
