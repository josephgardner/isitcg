const CACHE = 'isitcg-v1'

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
    caches.open(CACHE).then(cache => cache.addAll(PRECACHE))
  )
  self.skipWaiting()
})

self.addEventListener('activate', e => {
  e.waitUntil(
    caches.keys().then(keys =>
      Promise.all(keys.filter(k => k !== CACHE).map(k => caches.delete(k)))
    )
  )
  self.clients.claim()
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
          const fresh = fetch(e.request).then(response => {
            if (response.ok) cache.put(e.request, response.clone())
            return response
          })
          return cached || fresh
        })
      )
    )
    return
  }

  // Cache-first for everything else
  e.respondWith(
    caches.match(e.request).then(cached => {
      if (cached) return cached
      return fetch(e.request).then(response => {
        if (response.ok && e.request.url.startsWith(self.location.origin)) {
          const clone = response.clone()
          caches.open(CACHE).then(cache => cache.put(e.request, clone))
        }
        return response
      })
    })
  )
})
