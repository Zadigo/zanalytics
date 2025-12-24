import type { PostData } from "./types/index.ts"

export class ZanalyticsModel {
  apiUrl = 'http://127.0.0.1:9000/v1/analytics'
  apiKey = null

  constructor() {}

  async initialize(): Promise<void>{
    console.log("⭐️ Zanlaytics Model Initialized")

    if (typeof window === "undefined") return

    window.Zanalytics = this

    // Send POST request to server
    await this.post({
      timestamp: Date.now(),
      event: {
        name: 'load',
        url: window.location.href
      }
    })
  }

  async post(data: PostData) {
    const response = await fetch(this.apiUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    })

    if (!response.ok) {
      console.error(`Error posting analytics data: ${response.statusText}`)
    }

    console.log("✅ Analytics data posted successfully.")
  }
}

(function () {
  const instance = new ZanalyticsModel()
  
  window.onload = async () => {
    await instance.initialize()
  }

  window.onpageshow = () => {
    // Do something
    console.log(document.hasFocus())
    console.log(document.visibilityState)
  }

  window.onbeforeunload = () => {
    // Do something
  }

  window.onpagehide = () => {
    // Do something
  }

  return {
    instance
  }
})()
