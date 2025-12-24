export type Nullable<T> = T | null

export interface ZanalyticsModel {
  apiKey: Nullable<string>
  initialize: () => void
}

declare global {
  interface Window {
    Zanalytics: ZanalyticsModel
  }
}
