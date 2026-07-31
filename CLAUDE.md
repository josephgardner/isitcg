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
2. Product data is JSON-serialized and encoded as versioned UTF-8 Base64URL in the URL hash
3. Hash routing (`#HASH` = results, `#edit/HASH` = prefill form) handles navigation
4. Every matching rule is retained, including multiple classifications for the same ingredient
5. Category panels are sorted by `Rank`; the headline uses the most severe match (`danger` > `warning` > `good`)
6. If there is no danger or warning but at least one ingredient is unmatched, the headline result is `unknown` rather than approved

### Key Files

- `docs/index.html` — SPA shell
- `docs/src/app.js` — hash routing and app bootstrap
- `docs/src/isitcg.js` — core logic: `analyze`, `encode`, `decode`, `normalize`, fuzzy matching
- `docs/src/render.js` — DOM rendering (`renderForm`, `renderResults`)
- `docs/static/css/site.css` — styles
- `docs/ingredientrules.json` — rules data (synced from root by CI)
- `ingredientrules.json` — source of truth (community edits via PR)

### Matching Logic

Uses equality (not substring) matching. `normalize()` strips `[bracket content]`, `(paren content)`, accents, and formatting while preserving Unicode letters and numbers, then lowercases. Ingredients are split on `/` for per-part matching. Each submitted ingredient is checked against every rule, so intentional multi-category matches are preserved rather than consumed by the first matching rule.

### Data Structures

**Rule** — name, description, result type (`danger`/`warning`/`good`; legacy `success` means `good`), rank, and ingredient list

**Results** — ProductName, overall Result (`danger`/`warning`/`good`/`unknown`), all matching rules sorted by rank, and Remainder (unmatched ingredients)

## Deployment

GitHub Pages serves `docs/`. CNAME: `www.isitcg.com`. Push to `main` → CI runs `npm test` → auto-deploys.
