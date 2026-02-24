# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run tests
npm test

# Run tests in watch mode
npm run test:watch
```

## Architecture

Vanilla JS SPA served via GitHub Pages from the `docs/` directory. No build step.

### Core Flow
1. User submits product name + comma-separated ingredients via form
2. Product data is JSON-serialized, DEFLATE-compressed, and Base64-encoded into a URL hash
3. Hash routing (`#HASH` = results, `#edit/HASH` = prefill form) handles navigation
4. Results show matched rules ranked by priority (danger > warning > good)

### Key Files

- `docs/index.html` — SPA shell
- `docs/src/app.js` — hash routing and app bootstrap
- `docs/src/isitcg.js` — core logic: `analyze`, `encode`, `decode`, `normalize`, fuzzy matching
- `docs/src/render.js` — DOM rendering (`renderForm`, `renderResults`)
- `docs/static/css/site.css` — styles
- `docs/ingredientrules.json` — rules data (synced from root by CI)
- `ingredientrules.json` — source of truth (community edits via PR)

### Matching Logic

Uses equality (not substring) matching. `normalize()` strips `[bracket content]`, `(paren content)`, and non-word chars, then lowercases. Ingredients are split on `/` for per-part matching.

### Data Structures

**Rule** — name, description, result type (`danger`/`warning`/`good`), rank, and ingredient list

**Results** — ProductName, overall Result, sorted MatchingRules, and Remainder (unmatched ingredients)

## Deployment

GitHub Pages serves `docs/`. CNAME: `www.isitcg.com`. Push to `main` → CI runs `npm test` → auto-deploys.
