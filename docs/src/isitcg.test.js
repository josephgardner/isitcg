// Behavioral tests ported from internal/isitcg/testdata/rule-tests.yml
import { describe, test, expect } from 'vitest'
import { analyze, parts, matchAny, encode, decode } from './isitcg.js'

describe('parts()', () => {
  test('splits on comma with whitespace', () => {
    expect(parts('one, two, three')).toEqual(['one', 'two', 'three'])
  })

  test('trims whitespace', () => {
    expect(parts('one,            two,\n       three')).toEqual(['one', 'two', 'three'])
  })

  test('trims leading/trailing periods', () => {
    expect(parts('.one., .two., .three.')).toEqual(['one', 'two', 'three'])
  })

  test('preserves commas inside parentheses', () => {
    expect(parts('Water (Aqua, Eau), Glycerin')).toEqual(['Water (Aqua, Eau)', 'Glycerin'])
  })

  test('handles empty string', () => {
    expect(parts('')).toEqual([])
  })

  test('filters out empty parts from multiple commas', () => {
    expect(parts('  ,,,,   , ,   , \u0b99&*^%$       ,')).toEqual(['\u0b99&*^%$'])
  })
})

describe('analyze() — ported from rule-tests.yml', () => {
  test('Simplest test: one ingredient matches one rule', () => {
    const res = analyze('', 'one, two, three', [
      { Result: 'one-result', Ingredients: ['one'] },
    ])
    expect(res.result).toBe('good')
    expect(res.matches).toHaveLength(1)
    expect(res.matches[0].Result).toBe('one-result')
    expect(res.matches[0].Ingredients).toEqual(['one'])
    expect(res.remainder).toEqual(['two', 'three'])
  })

  test('split on comma: all ingredients match', () => {
    const res = analyze('', 'one, two, three', [
      { Result: 'comma-result', Ingredients: ['one', 'two', 'three'] },
    ])
    expect(res.result).toBe('good')
    expect(res.matches[0].Ingredients).toEqual(['one', 'two', 'three'])
    expect(res.remainder).toEqual([])
  })

  test('trims whitespace in ingredients', () => {
    const res = analyze('', 'one,            two,\n       three', [
      { Result: 'trim-result', Ingredients: ['one', 'two', 'three'] },
    ])
    expect(res.matches[0].Ingredients).toEqual(['one', 'two', 'three'])
    expect(res.remainder).toEqual([])
  })

  test('trims periods from ingredients', () => {
    const res = analyze('', '.one., .two., .three.', [
      { Result: 'comma-result', Ingredients: ['one', 'two', 'three'] },
    ])
    expect(res.matches[0].Ingredients).toEqual(['one', 'two', 'three'])
    expect(res.remainder).toEqual([])
  })

  test('matching ignores case', () => {
    const res = analyze('', 'ONE, Two, thrEE', [
      { Result: 'comma-result', Ingredients: ['one', 'two', 'three'] },
    ])
    // Matched ingredients are original (un-normalized) strings from input
    expect(res.matches[0].Ingredients).toEqual(['ONE', 'Two', 'thrEE'])
    expect(res.remainder).toEqual([])
  })

  test('results sorted by Rank ascending', () => {
    const res = analyze('', 'one, two, three, four', [
      { Result: 'good',   Rank: 2, Ingredients: ['four'] },
      { Result: 'danger', Rank: 3, Ingredients: ['one'] },
      { Result: 'danger', Rank: 1, Ingredients: ['two'] },
    ])
    expect(res.result).toBe('danger')
    expect(res.matches[0].Rank).toBe(1)  // two
    expect(res.matches[1].Rank).toBe(2)  // four
    expect(res.matches[2].Rank).toBe(3)  // one
    expect(res.remainder).toEqual(['three'])
  })

  test('no matching rules: everything is remainder', () => {
    const res = analyze('', 'one, two, three', [])
    expect(res.result).toBe('good')
    expect(res.matches).toEqual([])
    expect(res.remainder).toEqual(['one', 'two', 'three'])
  })

  test('match with slashes: exact slash ingredient matches', () => {
    const res = analyze('', 'one/two/three', [
      { Result: 'slash-result', Ingredients: ['one/two/three'] },
    ])
    expect(res.matches[0].Ingredients).toEqual(['one/two/three'])
    expect(res.remainder).toEqual([])
  })

  test('match slashes one part: ingredient matches via slash part', () => {
    const res = analyze('', 'one/two/three', [
      { Result: 'two-result', Ingredients: ['two'] },
    ])
    // "one/two/three" is matched because its slash part "two" matches
    expect(res.matches[0].Ingredients).toEqual(['one/two/three'])
    expect(res.remainder).toEqual([])
  })

  test('match slashes two parts: every matching rule is retained', () => {
    const res = analyze('', 'one/two/three', [
      { Result: 'one-result', Ingredients: ['one'] },
      { Result: 'two-result', Ingredients: ['two'] },
    ])
    expect(res.matches).toHaveLength(2)
    expect(res.matches[0].Result).toBe('one-result')
    expect(res.matches[0].Ingredients).toEqual(['one/two/three'])
    expect(res.matches[1].Result).toBe('two-result')
    expect(res.matches[1].Ingredients).toEqual(['one/two/three'])
    expect(res.remainder).toEqual([])
  })

  test('match more descriptive: equality not substring (benzyl alcohol != alcohol)', () => {
    const res = analyze('', 'alcohol, benzyl alcohol', [
      { Result: 'one', Ingredients: ['alcohol'] },
      { Result: 'two', Ingredients: ['Benzyl alcohol'] },
    ])
    expect(res.matches).toHaveLength(2)
    expect(res.matches[0].Result).toBe('one')
    expect(res.matches[0].Ingredients).toEqual(['alcohol'])
    expect(res.matches[1].Result).toBe('two')
    expect(res.matches[1].Ingredients).toEqual(['benzyl alcohol'])
    expect(res.remainder).toEqual([])
  })

  test('ignore formatting: special chars and brackets stripped for matching', () => {
    const res = analyze('', '*Foo, Fo - o, Fo (asdf) o, F [asdf] oo', [
      { Result: 'one', Ingredients: ['foo'] },
    ])
    expect(res.matches[0].Ingredients).toEqual(['*Foo', 'Fo - o', 'Fo (asdf) o', 'F [asdf] oo'])
    expect(res.remainder).toEqual([])
  })

  test('ignore invalid characters: unicode and special chars stripped', () => {
    const res = analyze('', 'F!@#\u2593\u263a$%^o+=\u045d\u0b99(asdf)-o, F   \u01ce\u1e50 [ \u0b99&*^%$asdf]~`oo  ', [
      { Result: 'one', Ingredients: ['foo'] },
    ])
    expect(res.matches[0].Ingredients).toHaveLength(2)
    expect(res.remainder).toEqual([])
  })

  test('remainder only: non-matching ingredient ends up in remainder', () => {
    const res = analyze('', '  ,,,,   , ,   , \u0b99&*^%$       ,', [
      { Result: 'one', Ingredients: ['foo'] },
    ])
    expect(res.matches).toEqual([])
    expect(res.remainder).toEqual(['\u0b99&*^%$'])
  })

  test('empty ingredients: empty result', () => {
    const res = analyze('', '', [
      { Result: 'one', Ingredients: ['foo'] },
    ])
    expect(res.matches).toEqual([])
    expect(res.remainder).toEqual([])
  })

  test('preserve commas in parentheses: Water (Aqua, Eau) is one ingredient', () => {
    const res = analyze('', 'Water (Aqua, Eau), Glycerin', [
      { Result: 'water-result', Ingredients: ['water'] },
      { Result: 'glycerin-result', Ingredients: ['glycerin'] },
    ])
    expect(res.matches).toHaveLength(2)
    expect(res.matches[0].Result).toBe('water-result')
    expect(res.matches[0].Ingredients).toEqual(['Water (Aqua, Eau)'])
    expect(res.matches[1].Result).toBe('glycerin-result')
    expect(res.matches[1].Ingredients).toEqual(['Glycerin'])
    expect(res.remainder).toEqual([])
  })

  test('danger escalates result, warning does not override danger', () => {
    const res = analyze('', 'a, b, c', [
      { Result: 'warning', Ingredients: ['a'] },
      { Result: 'danger',  Ingredients: ['b'] },
      { Result: 'warning', Ingredients: ['c'] },
    ])
    expect(res.result).toBe('danger')
  })

  test('all matching rules contribute to the verdict regardless of rule order', () => {
    const res = analyze('', 'Cyclopentasiloxane', [
      {
        Name: 'Silicones to Use with Caution',
        Result: 'warning',
        Rank: 5,
        Ingredients: ['Cyclopentasiloxane'],
      },
      {
        Name: 'Silicones to Avoid',
        Result: 'danger',
        Rank: 2,
        Ingredients: ['Cyclopentasiloxane'],
      },
    ])

    expect(res.result).toBe('danger')
    expect(res.matches.map(rule => rule.Name)).toEqual([
      'Silicones to Avoid',
      'Silicones to Use with Caution',
    ])
    expect(res.remainder).toEqual([])
  })

  test('warning escalates good result', () => {
    const res = analyze('', 'a, b', [
      { Result: 'warning', Ingredients: ['a'] },
      { Result: 'good',    Ingredients: ['b'] },
    ])
    expect(res.result).toBe('warning')
  })
})

describe('encode / decode', () => {
  test('round-trips product name and ingredients', () => {
    const hash = encode('Test Product', 'water, glycerin')
    const { n, i } = decode(hash)
    expect(n).toBe('Test Product')
    expect(i).toBe('water, glycerin')
  })

  test('handles empty product name', () => {
    const hash = encode('', 'sodium lauryl sulfate')
    const { n, i } = decode(hash)
    expect(n).toBe('')
    expect(i).toBe('sodium lauryl sulfate')
  })

  test('encoded string is URL-safe (no +, /, or =)', () => {
    const hash = encode('A Product', 'water, glycerin, citric acid')
    expect(hash).not.toMatch(/[+/=]/)
  })
})
