import { afterAll, describe, expect, test, vi } from 'vitest'

describe('hash routing', () => {
  test('a stale glossary response cannot overwrite a newer route', async () => {
    const container = { innerHTML: '' }
    const listeners = {}
    let resolveGlossary
    const glossaryResponse = new Promise(resolve => { resolveGlossary = resolve })

    vi.stubGlobal('document', {
      getElementById: id => id === 'app' ? container : null,
    })
    vi.stubGlobal('window', {
      addEventListener: (type, listener) => { listeners[type] = listener },
    })
    vi.stubGlobal('navigator', {})
    vi.stubGlobal('location', { hash: '#glossary/first' })
    vi.stubGlobal('fetch', vi.fn(url => {
      if (url === 'ingredientrules.json') {
        return Promise.resolve({ json: async () => ({ Rules: [] }) })
      }
      if (url === 'glossary/first.md') return glossaryResponse
      throw new Error(`Unexpected URL: ${url}`)
    }))

    await import('./app.js')
    await vi.waitFor(() => {
      expect(fetch).toHaveBeenCalledWith('glossary/first.md')
    })

    location.hash = ''
    listeners.hashchange()
    expect(container.innerHTML).toContain('ingredient-form')

    resolveGlossary({
      ok: true,
      text: async () => '# Stale Glossary',
    })
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(container.innerHTML).toContain('ingredient-form')
    expect(container.innerHTML).not.toContain('Stale Glossary')
  })
})

afterAll(() => {
  vi.unstubAllGlobals()
})
