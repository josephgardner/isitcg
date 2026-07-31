// Normalize an ingredient or rule ingredient for matching:
// lowercase, remove [bracket content], remove (paren content), remove non-word chars.
function normalize(s) {
  return s.toLowerCase()
    .replace(/\[.*?\]/g, '')
    .replace(/\(.*?\)/g, '')
    .replace(/\W/g, '')
}

// Split an ingredient string by commas, respecting parentheses nesting.
// Trims whitespace and leading/trailing periods from each part.
export function parts(ingredientsStr) {
  const result = []
  let current = ''
  let depth = 0
  for (const ch of ingredientsStr) {
    if (ch === '(') { depth++; current += ch }
    else if (ch === ')') { depth--; current += ch }
    else if (ch === ',' && depth === 0) {
      const part = current.trim().replace(/^\.+|\.+$/g, '')
      if (part.length > 0) result.push(part)
      current = ''
    } else {
      current += ch
    }
  }
  const part = current.trim().replace(/^\.+|\.+$/g, '')
  if (part.length > 0) result.push(part)
  return result
}

// Returns true if ingredient matches any of the candidates.
// Matching is done by normalizing both sides and comparing for equality.
// Also splits ingredient on "/" and checks each part separately.
export function matchAny(ingredient, candidates) {
  const norm = normalize(ingredient)
  const slashParts = ingredient.split('/')
  return candidates.some(c => {
    const nc = normalize(c)
    if (norm === nc) return true
    return slashParts.some(p => normalize(p) === nc)
  })
}

// Analyze ingredients against rules.
// Returns { productName, result, matches, remainder }.
// result is 'danger', 'warning', or 'good'.
// matches is sorted by Rank ascending.
// remainder contains unmatched ingredients.
export function analyze(productName, ingredientsStr, rules) {
  const ingredients = parts(ingredientsStr)
  const matches = []
  const matchedIngredients = new Set()
  let result = 'good'

  for (const rule of rules) {
    const hit = ingredients.filter(i => matchAny(i, rule.Ingredients))
    if (hit.length === 0) continue

    matches.push({ ...rule, Ingredients: hit })
    hit.forEach(ingredient => matchedIngredients.add(ingredient))

    if (rule.Result === 'danger') result = 'danger'
    else if (rule.Result === 'warning' && result === 'good') result = 'warning'
  }

  matches.sort((a, b) => (a.Rank || 0) - (b.Rank || 0))
  const remainder = ingredients.filter(ingredient => !matchedIngredients.has(ingredient))
  return { productName, result, matches, remainder }
}

// Slugify a rule name for use in URLs and filenames.
export function slugify(name) {
  return name.toLowerCase().replace(/[\s-]+/g, '-').replace(/[^a-z0-9-]/g, '')
}

const UTF8_HASH_PREFIX = 'v1.'

function toBase64Url(binary) {
  return btoa(binary)
    .replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
}

function fromBase64Url(encoded) {
  const b64 = encoded.replace(/-/g, '+').replace(/_/g, '/')
  const padding = '='.repeat((4 - (b64.length % 4)) % 4)
  return atob(b64 + padding)
}

// Encode product name + ingredients as versioned UTF-8 Base64URL.
export function encode(name, ingredients) {
  const bytes = new TextEncoder().encode(JSON.stringify({ n: name, i: ingredients }))
  let binary = ''
  for (const byte of bytes) binary += String.fromCharCode(byte)
  return UTF8_HASH_PREFIX + toBase64Url(binary)
}

// Decode a hash back to { n, i }, retaining support for unversioned legacy hashes.
export function decode(hash) {
  if (!hash.startsWith(UTF8_HASH_PREFIX)) {
    return JSON.parse(fromBase64Url(hash))
  }

  const binary = fromBase64Url(hash.slice(UTF8_HASH_PREFIX.length))
  const bytes = Uint8Array.from(binary, ch => ch.charCodeAt(0))
  return JSON.parse(new TextDecoder('utf-8', { fatal: true }).decode(bytes))
}
