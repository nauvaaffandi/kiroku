// ============================================================================
// Kiroku Custom Service Worker
// Handles push notifications + notification clicks.
// Workbox runtime caching is injected at build time by VitePWA injectManifest.
// ============================================================================

importScripts('https://storage.googleapis.com/workbox-cdn/releases/7.4.1/workbox-sw.js')

const { precaching, routing, strategies } = workbox

// Register precache (injected by VitePWA injectManifest at build time)
precaching.precacheAndRoute(self.__WB_MANIFEST)

// Workbox runtime caching (preserved from original vite.config.ts)
routing.registerRoute(
  /^https?.*\/api\/v1\//,
  new strategies.NetworkFirst({
    cacheName: 'api-cache',
    plugins: []
  })
)

routing.registerRoute(
  /^https?.*\/shared\/.*/,
  new strategies.NetworkFirst({
    cacheName: 'shared-cache',
    plugins: []
  })
)

const CACHE_NAMES = {
  API: 'api-cache',
  SHARED: 'shared-cache',
}

// Push notification handler
self.addEventListener('push', (event) => {
  let data = {}
  try {
    data = event.data ? event.data.json() : {}
  } catch {
    data = {}
  }

  const title = data.title || 'Kiroku'
  const body = data.body || ''
  const url = data.url || '/'

  const options = {
    body,
    icon: '/icon.svg',
    badge: '/icon.svg',
    tag: 'kiroku-push',
    requireInteraction: false,
    data: {
      url,
      timestamp: Date.now(),
    },
  }

  event.waitUntil(self.registration.showNotification(title, options))
})

// Notification click handler
self.addEventListener('notificationclick', (event) => {
  event.notification.close()

  const url = event.notification.data?.url || '/'

  event.waitUntil(
    self.clients.matchAll({ type: 'window', includeUncontrolled: true }).then((clients) => {
      // If an app window is already open, focus it and navigate
      for (const client of clients) {
        if (client.url.includes(self.location.origin)) {
          return client.focus().then(() => {
            client.postMessage({ type: 'NAVIGATE_TO', url })
          })
        }
      }
      // Otherwise open a new window
      return self.clients.openWindow(url)
    })
  )
})

// Handle subscription changes (expired/invalid)
self.addEventListener('pushsubscriptionchange', (event) => {
  console.log('[SW] pushsubscriptionchange: subscription expired')
  // The app layer should re-subscribe on next load
})

// Optional: handle messages from the app (e.g., navigate after notification click)
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting()
  }
})
