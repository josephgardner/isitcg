import { describe, expect, test } from 'vitest'
import { renderGlossary, renderResults } from './render.js'

describe('renderResults()', () => {
  test('describes an unknown result without claiming approval', () => {
    const container = { innerHTML: '' }

    renderResults(container, {
      productName: 'Mystery Product',
      result: 'unknown',
      matches: [],
      remainder: ['Mystery Ingredient'],
    })

    expect(container.innerHTML).toContain('Unable to Verify')
    expect(container.innerHTML).toContain('Unknown Ingredients')
    expect(container.innerHTML).not.toContain('Great news')
  })
})

describe('renderGlossary()', () => {
  test('links back to the result that opened the glossary', () => {
    const container = { innerHTML: '' }

    renderGlossary(container, '<h1>Silicones</h1>', 'silicones', 'result-hash')

    expect(container.innerHTML).toContain('href="#result-hash"')
  })
})
