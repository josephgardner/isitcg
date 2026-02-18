function esc(s) {
  if (!s) return ''
  return String(s)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

const WIKI_BASE = 'https://github.com/josephgardner/isitcg/wiki/'

function wikiUrl(ruleName) {
  return WIKI_BASE + encodeURIComponent(ruleName.replace(/[\s-]+/g, '-'))
}

const EXAMPLES = [
  { label: 'CG approved',           result: 'good',    hash: 'eyJuIjoiIiwiaSI6IldhdGVyLCBDZXR5bCBBbGNvaG9sLCBDb2NhbWlkb3Byb3B5bCBCZXRhaW5lLCBQYW50aGVub2wsIENpdHJpYyBBY2lkIn0' },
  { label: 'Approved with caution',  result: 'warning', hash: 'eyJuIjoiIiwiaSI6IldhdGVyLCBTaGVhIEJ1dHRlciwgQ2V0eWwgQWxjb2hvbCwgR2x5Y2VyaW4sIEJlaGVudHJpbW9uaXVtIE1ldGhvc3VsZmF0ZSwgRnJhZ3JhbmNlLCBDaXRyaWMgQWNpZCJ9' },
  { label: 'Not CG approved',        result: 'danger',  hash: 'eyJuIjoiIiwiaSI6IldhdGVyIChBcXVhKSwgU29kaXVtIExhdXJ5bCBTdWxmYXRlLCBDb2NhbWlkb3Byb3B5bCBCZXRhaW5lLCBEaW1ldGhpY29uZSwgR2x5Y2VyaW4sIENldHlsIEFsY29ob2wsIENpdHJpYyBBY2lkLCBGcmFncmFuY2UifQ' },
]

export function renderForm(container, prefilled = {}) {
  const exampleLinks = EXAMPLES.map(e =>
    `<a href="#${e.hash}" class="example-link example-${e.result}">${e.label}</a>`
  ).join('')

  container.innerHTML = `
    <div class="titlebar">
      <h3>Enter ingredients</h3>
      <p>1. Enter the product name (optional).<br>
         2. Paste the ingredient list, comma-separated.<br>
         3. Click <strong>Submit</strong> to check if it's CG approved.</p>
      <div class="examples-row">
        <span class="examples-label">See an example:</span>
        ${exampleLinks}
      </div>
    </div>
    <div class="form-wrap">
      <form id="ingredient-form">
        <input type="text" id="productname" placeholder="Product name (optional)" value="${esc(prefilled.name || '')}">
        <textarea id="ingredients" placeholder="Ingredients (Water, Glycerin, Citric Acid, ...)">${esc(prefilled.ingredients || '')}</textarea>
        <button type="submit" class="btn btn-secondary">Submit</button>
      </form>
    </div>
  `
}

function resultLabel(r) {
  if (r === 'danger')  return 'Not CG Approved'
  if (r === 'warning') return 'CG Approved with caution'
  return 'CG Approved'
}

function resultMessage(r) {
  if (r === 'danger')  return 'Uh oh! This product contains ingredients that aren\'t CG approved.'
  if (r === 'warning') return 'This product is CG approved, but contains some ingredients that may not work for everyone.'
  return 'Great news — this product is CG approved!'
}

export function renderResults(container, { productName, result, matches, remainder }, hash) {
  const verdictLabel = resultLabel(result)
  const verdictMessage = resultMessage(result)

  const editBtn = hash
    ? `<a href="#edit/${esc(hash)}" class="btn btn-ghost">Edit</a>`
    : ''

  const matchPanels = matches.map(m => `
    <div class="panel panel-${esc(m.Result)}">
      <div class="panel-heading">${resultLabel(m.Result)}</div>
      <div class="panel-body">
        <h4>${esc(m.Name || '')}</h4>
        <p>${esc(m.Description || '')} <a href="${wikiUrl(m.Name)}" target="_blank" rel="noopener">(More info)</a></p>
      </div>
      <ul class="list-group">
        ${m.Ingredients.map(i => `<li>${esc(i)}</li>`).join('')}
      </ul>
    </div>
  `).join('')

  const remainderPanel = remainder.length > 0 ? `
    <div class="panel panel-primary">
      <div class="panel-heading">Probably OK</div>
      <div class="panel-body">
        <h4>Remaining Ingredients</h4>
        <p>These weren't matched to any known rule — likely fine, or not yet in our database.</p>
      </div>
      <ul class="list-group">
        ${remainder.map(i => `<li>${esc(i)}</li>`).join('')}
      </ul>
    </div>
  ` : ''

  container.innerHTML = `
    <div class="titlebar result-titlebar">
      ${productName ? `<h2 class="result-product">${esc(productName)}</h2>` : ''}
      <div class="result-verdict verdict-${result}">
        <span class="verdict-dot"></span>
        <span>${verdictLabel}</span>
      </div>
      <p class="result-message">${verdictMessage}</p>
      <div class="actions">
        ${editBtn}
        <a href="#" class="btn btn-ghost">Start Over</a>
      </div>
    </div>
    ${matchPanels}
    ${remainderPanel}
    <div class="community-callout">
      <strong>Wrong result or missing ingredient?</strong> This analyzer is entirely community-driven.<br>
      <a href="https://github.com/josephgardner/isitcg/blob/main/ingredientrules.json" target="_blank" rel="noopener">Edit ingredientrules.json on GitHub</a>
      &nbsp;&middot;&nbsp;
      <a href="https://github.com/josephgardner/isitcg/issues" target="_blank" rel="noopener">Open an issue</a>
    </div>
  `
}
