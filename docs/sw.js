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
