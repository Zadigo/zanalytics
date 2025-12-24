import type { BrowserStatistics, PerformanceData, PostData } from './types/index.ts'

export class ZanalyticsModel {
  apiUrl = 'http://127.0.0.1:9000/v1/analytics'
  apiKey = null
  private initialized = false

  constructor() {}

  async initialize(_e?: Event): Promise<void>{
    if (this.initialized) {
      console.log("⚠️ Zanalytics already initialized")
      return
    }

    if (typeof window === "undefined") return

    console.log("⭐️ Zanlaytics Model Initialized")
    this.initialized = true

    window.Zanalytics = this

    try {
      // Send POST request to server
      await this.post({
        timestamp: Date.now(),
        event: {
          name: 'load',
          url: window.location.href
        },
        urlDetails: {
          protocol: window.location.protocol,
          origin: window.location.origin,
          pathname: window.location.pathname,
          search: window.location.search
        }
      })
    } catch (error) {
      throw new Error("Error during initialization POST:" + (error as Error).message)
    }
  }

  /**
   * Helper method to post analytics data to the server.
   * @param data - The analytics data to post.
   */
  async post(data: PostData) {
    try {
      const headers: Record<string, string> = {
        'Content-Type': 'application/json'
      }

      // Add API key if available
      if (this.apiKey) {
        headers['Authorization'] = `Bearer ${this.apiKey}`
      }

      const response = await fetch(this.apiUrl, {
        method: 'POST',
        headers,
        body: JSON.stringify(data)
      })

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`)
      }

      console.log("✅ Analytics data posted successfully")
    } catch (error) {
      console.error(`Error posting analytics data:`, error)
      throw error
    }
  }

  /**
   * Loads performance data using the Performance API.
   * @returns A promise that resolves to an object containing performance data.
   */
  async loadPerformanceData<P extends PerformanceData = PerformanceData>(): Promise<P | undefined> {
    if (typeof window === "undefined") return

    const data = performance.getEntriesByType("navigation")[0]
    const returnData = {
      dnsTime: 0,
      tcpTime: 0,
      requestTime: 0,
      responseTime: 0,
      domProcessing: 0,
      totalLoadTime: 0,
      domReady: 0
    } as P

    if (typeof data !== "undefined") {
      try {
        // @ts-ignore DNS lookup time
        returnData.dnsTime = data.domainLookupEnd - data.domainLookupStart
        // @ts-ignore TCP connection time
        returnData.tcpTime = data.connectEnd - data.connectStart
        // @ts-ignore Request time
        returnData.requestTime = data.responseStart - data.requestStart
        // @ts-ignore Response time
        returnData.responseTime = data.responseEnd - data.responseStart
        // @ts-ignore DOM processing time
        returnData.domProcessing = data.domComplete - data.domInteractive
        // @ts-ignore Total load time
        returnData.totalLoadTime = data.loadEventEnd - data.fetchStart
        // @ts-ignore Time to DOM ready
        returnData.domReady = data.domContentLoadedEventEnd - data.fetchStart
        return returnData
      } catch (error) {
        console.error("Error parsing performance data:", error)
      }
    }
    return returnData
  }

  /**
   * Loads browser statistics such as user agent, language, platform, etc.
   * @returns A promise that resolves to an object containing browser statistics.
   */
  async loadBroswerStatistics<B extends BrowserStatistics = BrowserStatistics>(): Promise<B | undefined> {
    if (typeof window === "undefined") return

    return {
      userAgent: window.navigator.userAgent,
      language: window.navigator.language,
      platform: window.navigator.platform,
      cookiesEnabled: window.navigator.cookieEnabled,
      onLine: window.navigator.onLine,
      screenResolution: `${window.screen.width}x${window.screen.height}`,
      viewportSize: `${window.innerWidth}x${window.innerHeight}`,
      colorDepth: window.screen.colorDepth,
      pixelRatio: window.devicePixelRatio
    } as B
  }
}

(function () {
  const instance = new ZanalyticsModel()

  // Use DOMContentLoaded for faster initialization
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', async () => {
      await instance.initialize()
    })
  } else {
    // DOM already loaded
    instance.initialize()
  }

  // Handle page visibility changes (for back/forward cache)
  document.addEventListener('visibilitychange', async () => {
    if (document.visibilityState === 'visible') {
      console.log('Page became visible')
      // Post a visibility event instead of reinitializing
      await instance.post({
        timestamp: Date.now(),
        event: { 
          name: 'navigation', 
          state: 'visible'
        }
      } as PostData)
    }
  })

  // Handle page hide (for sending final analytics before unload)
  window.addEventListener('pagehide', () => {
    instance.post({
      timestamp: Date.now(),
      event: { 
        name: 'pagehide'
      }
    }).catch(err => console.error('Failed to send pagehide event:', err))
  })
  
  // window.onload = async (e: Event) => {
  //   console.log('On load:', e)
  //   await instance.initialize(e)
  // }

  // window.onpageshow = async (e: Event) => {
  //   console.log('On page show:', e)
  //   console.log(document.hasFocus())
  //   console.log(document.visibilityState)

  //   console.log('Performance Data:', await instance.loadPerformanceData())
  //   console.log('Browser Statistics:', await instance.loadBroswerStatistics())

  //   // // Scroll position using VueUse
  //   // const { x, y } = useScroll(window)
  //   // console.log(`Scroll Position - X: ${x.value}, Y: ${y.value}`)

  //   // // Mouse tracking
  //   // const { x: mouseX, y: mouseY } = useMouse()
  //   // console.log(`Mouse Position - X: ${mouseX.value}, Y: ${mouseY.value}`)
  // }

  // window.onbeforeunload = () => {
  //   // Do something
  // }

  // window.onpagehide = () => {
  //   // Do something
  // }
})()
