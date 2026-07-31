import { readFile } from 'node:fs/promises'
import { describe, expect, test } from 'vitest'
import { isValidRulesData } from './isitcg.js'

const rootRules = JSON.parse(
  await readFile(new URL('../../ingredientrules.json', import.meta.url), 'utf8'),
)
const pagesRules = JSON.parse(
  await readFile(new URL('../ingredientrules.json', import.meta.url), 'utf8'),
)

describe('ingredient rules deployment data', () => {
  test('the source-of-truth rules have the expected shape', () => {
    expect(isValidRulesData(rootRules)).toBe(true)
    expect(new Set(rootRules.Rules.map(rule => rule.Name)).size)
      .toBe(rootRules.Rules.length)
  })

  test('the GitHub Pages rules match the source of truth', () => {
    expect(pagesRules).toEqual(rootRules)
  })
})
