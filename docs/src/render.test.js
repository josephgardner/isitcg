import { describe, expect, test } from 'vitest'
import { renderForm, renderGlossary, renderResults } from './render.js'

describe('renderForm()', () => {
  test('provides programmatic labels for both form fields', () => {
    const container = { innerHTML: '' }

    renderForm(container)

    expect(container.innerHTML).toContain('<label for="productname"')
    expect(container.innerHTML).toContain('<label for="ingredients"')
  })
})

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
