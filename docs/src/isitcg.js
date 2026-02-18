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
  let remainder = parts(ingredientsStr)
  const matches = []
  let result = 'good'

  for (const rule of rules) {
    const hit = remainder.filter(i => matchAny(i, rule.Ingredients))
    if (hit.length === 0) continue

    matches.push({ ...rule, Ingredients: hit })
    remainder = remainder.filter(i => !hit.includes(i))

    if (rule.Result === 'danger') result = 'danger'
    else if (rule.Result === 'warning' && result === 'good') result = 'warning'
  }

  matches.sort((a, b) => (a.Rank || 0) - (b.Rank || 0))
  return { productName, result, matches, remainder }
}

// Slugify a rule name for use in URLs and filenames.
export function slugify(name) {
  return name.toLowerCase().replace(/[\s-]+/g, '-').replace(/[^a-z0-9-]/g, '')
}

// Encode product name + ingredients into a URL-safe base64 hash.
export function encode(name, ingredients) {
  return btoa(JSON.stringify({ n: name, i: ingredients }))
    .replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
}

// Decode a hash back to { n, i }.
export function decode(hash) {
  const b64 = hash.replace(/-/g, '+').replace(/_/g, '/')
  return JSON.parse(atob(b64))
}
