import { readFile } from 'node:fs/promises'
import { describe, expect, test } from 'vitest'

const rootRules = JSON.parse(
  await readFile(new URL('../../ingredientrules.json', import.meta.url), 'utf8'),
)
const pagesRules = JSON.parse(
  await readFile(new URL('../ingredientrules.json', import.meta.url), 'utf8'),
)

const supportedResults = new Set(['danger', 'warning', 'good', 'success'])

function normalizeIngredient(ingredient) {
  return ingredient.toLowerCase()
    .replace(/\[.*?\]/g, '')
    .replace(/\(.*?\)/g, '')
    .replace(/\W/g, '')
}

function expectValidRules(data) {
  expect(data).toHaveProperty('Rules')
  expect(Array.isArray(data.Rules)).toBe(true)
  expect(data.Rules.length).toBeGreaterThan(0)

  for (const rule of data.Rules) {
    expect(rule.Name).toEqual(expect.any(String))
    expect(rule.Name.length).toBeGreaterThan(0)
    expect(rule.Description).toEqual(expect.any(String))
    expect(supportedResults.has(rule.Result)).toBe(true)
    expect(Number.isFinite(rule.Rank)).toBe(true)
    expect(Array.isArray(rule.Ingredients)).toBe(true)
    expect(rule.Ingredients.length).toBeGreaterThan(0)
    expect(rule.Ingredients.every(ingredient =>
      typeof ingredient === 'string' && ingredient.trim().length > 0
    )).toBe(true)
    expect(rule.Ingredients.every(ingredient =>
      normalizeIngredient(ingredient).length > 0
    )).toBe(true)
  }
}

describe('ingredient rules deployment data', () => {
  test('the source-of-truth rules have the expected shape', () => {
    expectValidRules(rootRules)
    expect(new Set(rootRules.Rules.map(rule => rule.Name)).size)
      .toBe(rootRules.Rules.length)
  })

  test('the GitHub Pages rules match the source of truth', () => {
    expect(pagesRules).toEqual(rootRules)
  })
})
