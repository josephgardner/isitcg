import { analyze, encode, decode } from './isitcg.js'
import { renderForm, renderResults, renderGlossary, renderGlossaryError } from './render.js'

const container = document.getElementById('app')
let rules = []
let routeVersion = 0
let lastResultsHash = ''

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
  const version = ++routeVersion
  const hash = location.hash.slice(1) // strip leading #

  if (!hash) {
    lastResultsHash = ''
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

  if (hash.startsWith('glossary/')) {
    const slug = hash.slice(9)
    loadGlossary(slug, version)
    return
  }

  try {
    const { n, i } = decode(hash)
    const results = analyze(n, i, rules)
    lastResultsHash = hash
    renderResults(container, results, hash)
    ga4('results_viewed', { result_type: results.result })
  } catch (_) {
    renderForm(container)
    attachForm()
  }
}

async function loadGlossary(slug, version) {
  try {
    const res = await fetch(`glossary/${slug}.md`)
    if (!res.ok) throw new Error()
    const md = await res.text()
    if (version !== routeVersion) return
    renderGlossary(container, parseMarkdown(md), slug, lastResultsHash)
  } catch (_) {
    if (version !== routeVersion) return
    renderGlossaryError(container, lastResultsHash)
  }
}

// Minimal markdown parser — handles headings, paragraphs, lists, bold, italic, links, code.
function parseMarkdown(md) {
  function escHtml(s) {
    return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
  }

  function inline(s) {
    return escHtml(s)
      .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
      .replace(/__(.+?)__/g, '<strong>$1</strong>')
      .replace(/\*(.+?)\*/g, '<em>$1</em>')
      .replace(/_(.+?)_/g, '<em>$1</em>')
      .replace(/`([^`]+)`/g, '<code>$1</code>')
      .replace(/\[([^\]]+)\]\(([^)]+)\)/g, (_, text, url) =>
        `<a href="${escHtml(url)}" target="_blank" rel="noopener">${text}</a>`)
  }

  const blocks = []
  let listItems = []

  function flushList() {
    if (listItems.length) {
      blocks.push(`<ul>${listItems.join('')}</ul>`)
      listItems = []
    }
  }

  for (const line of md.split('\n')) {
    const t = line.trim()
    if      (t.startsWith('### ')) { flushList(); blocks.push(`<h3>${inline(t.slice(4))}</h3>`) }
    else if (t.startsWith('## '))  { flushList(); blocks.push(`<h2>${inline(t.slice(3))}</h2>`) }
    else if (t.startsWith('# '))   { flushList(); blocks.push(`<h1>${inline(t.slice(2))}</h1>`) }
    else if (t.startsWith('> '))   { flushList(); blocks.push(`<blockquote>${inline(t.slice(2))}</blockquote>`) }
    else if (t.match(/^[-*] /))    { listItems.push(`<li>${inline(t.slice(2))}</li>`) }
    else if (t === '')             { flushList() }
    else                           { flushList(); blocks.push(`<p>${inline(t)}</p>`) }
  }
  flushList()

  return blocks.join('\n')
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

if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => navigator.serviceWorker.register('/sw.js'))
}

init()
