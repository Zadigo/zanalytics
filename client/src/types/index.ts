export type Nullable<T> = T | null

export type BaseEventNames = 'load'
  | 'pagehide'
  | 'pageview'
  | 'click'
  | 'scroll'
  | 'navigation'
  | 'input'
  | 'submit'
  | 'custom'
  | 'error'

export type AnalyticsBaseEvent = {
  name: BaseEventNames
} & Partial<{
  url: string
}>

export interface ZanalyticsModel {
  apiUrl: string
  apiKey: Nullable<string>
  initialize: (e?: Event) => Promise<void>
  post: (data: PostData) => Promise<void>
}

export interface InitialPostData {
  url: string
  title: string
  referrer: string
  screenWidth: number
  screenHeight: number
}

interface UrlDetails {
  protocol: string
  origin: string
  pathname: string
  search: string
}

export type PostData = Partial<InitialPostData & { urlDetails: UrlDetails }> & {
  event: AnalyticsBaseEvent
  timestamp: number
}

declare global {
  interface Window {
    Zanalytics: ZanalyticsModel
  }
}

export interface PerformanceData {
  dnsTime: number
  tcpTime: number
  requestTime: number
  responseTime: number
  domProcessing: number 
  totalLoadTime: number
  domReady: number  
}

export type BrowserStatistics = {
  userAgent: string
  language: string
  platform: string
  cookiesEnabled: boolean
  onLine: boolean
  screenResolution: `${number}x${number}`
  viewportSize: `${number}x${number}`
  colorDepth: number
  pixelRatio: number
}
