# isitcg.com — Is It Curly Girl?

**[isitcg.com](https://www.isitcg.com)** — paste your product's ingredient list, find out if it's CG-approved in seconds.

---

## What is isitcg.com?

isitcg.com is a free tool for the curly hair community. Paste in a product's ingredient list, and it checks each ingredient against a community-maintained database of 1,000+ ingredients — flagging sulfates, silicones, drying alcohols, and other ingredients the Curly Girl Method recommends avoiding. It's open source and hosted on GitHub Pages, so anyone can help keep the ingredient database accurate.

## How it works

1. You paste a product's ingredient list into the form
2. The app checks each ingredient against a community-maintained list of rules (stored in [`ingredientrules.json`](ingredientrules.json))
3. It shows you which ingredients matched, what category they fall into (danger / caution / good), and which ingredients weren't recognized

The ingredient database is the heart of the tool. If an ingredient is missing or miscategorized, anyone can fix it.

## How to contribute

You don't need to know how to code. You just need a free GitHub account.

### Option A — Report a problem (easiest)

If you spot a mistake — wrong category, missing ingredient, misleading description — [open an issue](https://github.com/josephgardner/isitcg/issues/new/choose). There are templates to guide you. A maintainer or volunteer admin will take it from there.

### Option B — Fix it yourself

Want to add or correct an ingredient directly? You can edit [`ingredientrules.json`](ingredientrules.json) right here on GitHub — no terminal required.

See **[CONTRIBUTING.md](CONTRIBUTING.md)** for a step-by-step walkthrough.

### Option C — Improve a glossary page

The site has glossary pages explaining each ingredient category. They're in [`docs/glossary/`](docs/glossary/) and most are stubs waiting to be filled in. Click any file, hit the pencil icon, and write something helpful.

## Become a volunteer admin

The maintainer of this repo is a developer, not a CG expert. We're looking for **knowledgeable community members** to help review ingredient suggestions and pull requests.

**What admins do:**
- Review pull requests that add or change ingredients ("is this ingredient correct? is it in the right category?")
- Triage issues (confirm reports, ask for sources, label things)
- No coding required — the automated tests catch JSON syntax errors automatically

**Interested?** [Open an issue](https://github.com/josephgardner/isitcg/issues/new?title=I+want+to+help+as+a+volunteer+admin) and introduce yourself. Tell us a bit about your CG knowledge and how much time you have.

See **[guides/admin-guide.md](guides/admin-guide.md)** for more about what the role involves.

## Project structure

```
ingredientrules.json      ← The ingredient database (this is what you'll edit most often)
docs/                     ← The website (auto-deployed to isitcg.com — don't edit manually)
docs/glossary/            ← Glossary pages, one per ingredient category
guides/                   ← Contributor guides (admin guide, format reference)
.github/                  ← Issue templates and CI configuration
```

**How deployment works:** When a pull request is merged to `main`, GitHub Actions runs the tests automatically. If they pass, the site deploys to isitcg.com within a minute or two. You never need to touch the deployment yourself.

## Links

- **Live site:** [isitcg.com](https://www.isitcg.com)
- **Contributor guide:** [CONTRIBUTING.md](CONTRIBUTING.md)
- **Ingredient format reference:** [guides/ingredient-format.md](guides/ingredient-format.md)
- **Admin guide:** [guides/admin-guide.md](guides/admin-guide.md)
- **Open an issue:** [github.com/josephgardner/isitcg/issues](https://github.com/josephgardner/isitcg/issues)
