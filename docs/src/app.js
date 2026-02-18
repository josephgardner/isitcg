import { analyze, encode, decode } from './isitcg.js'
import { renderForm, renderResults } from './render.js'

const container = document.getElementById('app')
let rules = []

async function init() {
  try {
    const res = await fetch('ingredientrules.json')
    const data = await res.json()
    rules = data.Rules
  } catch (e) {
    container.innerHTML = '<p style="padding:24px;color:#ed878d">Failed to load ingredient rules. Try refreshing the page.</p>'
    return
  }
  route()
  window.addEventListener('hashchange', route)
}

function route() {
  const hash = location.hash.slice(1) // strip leading #

  if (!hash) {
    renderForm(container)
    attachForm()
    return
  }

  if (hash.startsWith('edit/')) {
    try {
      const { n, i } = decode(hash.slice(5))
      renderForm(container, { name: n, ingredients: i })
      attachForm()
    } catch (_) {
      renderForm(container)
      attachForm()
    }
    return
  }

  try {
    const { n, i } = decode(hash)
    const results = analyze(n, i, rules)
    renderResults(container, results, hash)
    ga4('results_viewed', { result_type: results.result })
  } catch (_) {
    renderForm(container)
    attachForm()
  }
}

function attachForm() {
  const form = document.getElementById('ingredient-form')
  if (!form) return
  form.addEventListener('submit', e => {
    e.preventDefault()
    const name = document.getElementById('productname').value.trim()
    const ingredients = document.getElementById('ingredients').value.trim()
    if (!ingredients) return
    ga4('form_submit', { has_product_name: Boolean(name) })
    location.hash = encode(name, ingredients)
  })
}

function ga4(event, params) {
  if (typeof gtag === 'function') gtag('event', event, params)
}

init()
