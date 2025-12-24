
export class ZanalyticsModel {
  apiKey = null

  constructor() {}

  initialize() {
    console.log("Zanlaytics Model Initialized")

    if (typeof window === "undefined") return
    window.Zanalytics = this
  }
}

(function () {
  const instance = new ZanalyticsModel()
  instance.initialize()

  return {
    instance
  }
})()
