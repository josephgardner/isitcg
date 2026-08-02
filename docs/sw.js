const CACHE = 'isitcg-v3'

const PRECACHE = [
  '/',
  '/src/app.js',
  '/src/isitcg.js',
  '/src/render.js',
  '/static/css/site.css',
  '/ingredientrules.json',
  '/apple-touch-icon.png',
  '/android-chrome-192x192.png',
  '/android-chrome-512x512.png',
]

self.addEventListener('install', e => {
  e.waitUntil(
    caches.open(CACHE)
      .then(cache => cache.addAll(PRECACHE))
      .then(() => self.skipWaiting())
  )
})

self.addEventListener('activate', e => {
  e.waitUntil(
    caches.keys().then(keys =>
      Promise.all(keys.filter(k => k !== CACHE).map(k => caches.delete(k)))
    ).then(() => self.clients.claim())
  )
})

self.addEventListener('fetch', e => {
  if (e.request.method !== 'GET') return

  // Network-first keeps app code and rules current while retaining an offline fallback.
  e.respondWith(
    fetch(e.request)
      .then(async response => {
        if (response.ok && e.request.url.startsWith(self.location.origin)) {
          const cache = await caches.open(CACHE)
          await cache.put(e.request, response.clone())
        }
        return response
      })
      .catch(async error => {
        const cached = await caches.match(e.request)
        if (cached) return cached
        throw error
      })
  )
})
