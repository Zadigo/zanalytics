export type Nullable<T> = T | null

export type BaseEventNames = 'load'
  | 'pageview'
  | 'click'
  | 'scroll'
  | 'navigate'
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
  initialize: () => Promise<void>
  post: (data: PostData) => Promise<void>
}

export interface InitialPostData {
  url: string
  title: string
  referrer: string
  screenWidth: number
  screenHeight: number
}

export type PostData = Partial<InitialPostData> & {
  event: AnalyticsBaseEvent
  timestamp: number
}

declare global {
  interface Window {
    Zanalytics: ZanalyticsModel
  }
}
