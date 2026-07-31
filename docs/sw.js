const CACHE = 'isitcg-v2'

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

  const url = new URL(e.request.url)
  const isRules = url.pathname.endsWith('/ingredientrules.json')

  if (isRules) {
    // Stale-while-revalidate: serve cache instantly, refresh in background
    e.respondWith(
      caches.open(CACHE).then(cache =>
        cache.match(e.request).then(cached => {
          const fresh = fetch(e.request).then(async response => {
            if (response.ok) await cache.put(e.request, response.clone())
            return response
          })
          if (!cached) return fresh
          e.waitUntil(fresh.catch(() => undefined))
          return cached
        })
      )
    )
    return
  }

  // Network-first keeps the app current while retaining an offline fallback.
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
