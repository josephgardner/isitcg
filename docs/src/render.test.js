import { describe, expect, test } from 'vitest'
import { renderResults } from './render.js'

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
