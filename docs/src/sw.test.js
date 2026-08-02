import { beforeAll, beforeEach, describe, expect, test, vi } from 'vitest'

const handlers = {}
const cache = {
  addAll: vi.fn(),
  match: vi.fn(),
  put: vi.fn(),
}

beforeAll(async () => {
  vi.stubGlobal('self', {
    location: { origin: 'https://www.isitcg.com' },
    clients: { claim: vi.fn() },
    skipWaiting: vi.fn(),
    addEventListener: (type, handler) => { handlers[type] = handler },
  })
  vi.stubGlobal('caches', {
    open: vi.fn().mockResolvedValue(cache),
    match: vi.fn(),
    keys: vi.fn().mockResolvedValue([]),
    delete: vi.fn(),
  })
  vi.stubGlobal('fetch', vi.fn())

  await import('../sw.js')
})

beforeEach(() => {
  vi.clearAllMocks()
})

function dispatchFetch(request) {
  let response
  handlers.fetch({
    request,
    respondWith: promise => { response = promise },
    waitUntil: vi.fn(),
  })
  return response
}

describe('service worker fetch policy', () => {
  test('ingredient rules prefer a fresh response over a cached response', async () => {
    const cached = { source: 'cache' }
    const fresh = { ok: true, source: 'network', clone: () => ({ source: 'clone' }) }
    caches.match.mockResolvedValue(cached)
    fetch.mockResolvedValue(fresh)

    const response = await dispatchFetch({
      method: 'GET',
      url: 'https://www.isitcg.com/ingredientrules.json',
    })

    expect(response).toBe(fresh)
    expect(cache.put).toHaveBeenCalled()
  })

  test('ingredient rules fall back to the cached response while offline', async () => {
    const cached = { source: 'cache' }
    caches.match.mockResolvedValue(cached)
    fetch.mockRejectedValue(new TypeError('offline'))

    const response = await dispatchFetch({
      method: 'GET',
      url: 'https://www.isitcg.com/ingredientrules.json',
    })

    expect(response).toBe(cached)
  })

  test('navigation prefers a fresh response over a cached response', async () => {
    const cached = { source: 'cache' }
    const fresh = { ok: true, source: 'network', clone: () => ({ source: 'clone' }) }
    caches.match.mockResolvedValue(cached)
    fetch.mockResolvedValue(fresh)

    const response = await dispatchFetch({
      method: 'GET',
      mode: 'navigate',
      url: 'https://www.isitcg.com/',
    })

    expect(response).toBe(fresh)
    expect(cache.put).toHaveBeenCalled()
  })

  test('navigation falls back to the cached response while offline', async () => {
    const cached = { source: 'cache' }
    caches.match.mockResolvedValue(cached)
    fetch.mockRejectedValue(new TypeError('offline'))

    const response = await dispatchFetch({
      method: 'GET',
      mode: 'navigate',
      url: 'https://www.isitcg.com/',
    })

    expect(response).toBe(cached)
  })
})
