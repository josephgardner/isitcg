# Ingredient Format Reference

This document explains the structure of [`ingredientrules.json`](../ingredientrules.json) in detail. It's for contributors who want to understand the data model fully before making changes.

---

## File structure

`ingredientrules.json` is a single JSON object with one key, `"Rules"`, whose value is an array of rule objects:

```json
{
  "Rules": [
    { ... },
    { ... }
  ]
}
```

There are 24 rules total, one for each ingredient category.

---

## Anatomy of a rule

Here is a complete, annotated example:

```json
{
  "Name": "Sulfates to Avoid",
  "Description": "These ingredients are sulfates that should be avoided.",
  "Result": "danger",
  "Rank": 1,
  "BlogUrl": "https://isitcg.lisagardnerdesign.com/say-no-to-sulfates-and-silicones-for-juicy-curls/",
  "Ingredients": [
    "sodium lauryl sulfate",
    "sodium laureth sulfate",
    "ammonium lauryl sulfate"
  ]
}
```

### `Name`
The category name displayed in results. This is also used to link to the corresponding glossary page.

### `Description`
A short description shown below the category name in results.

### `Result`
Controls the color and icon shown in results. Valid values:

| Value | Meaning | Display |
|-------|---------|---------|
| `danger` | Avoid this ingredient | Red / warning |
| `warning` | Use with caution | Orange / caution |
| `good` | CG-friendly | Green / check |
| `success` | Same as `good` (older alias, still works) | Green / check |

### `Rank`
Determines the order in which matched categories appear in results. **Lower number = higher up in the results.**

Danger categories are ranked 1–4. The full ranking:

| Rank | Category |
|------|---------|
| 1 | Sulfates to Avoid |
| 2 | Silicones to Avoid |
| 3 | Waxes and Oils to Avoid |
| 4 | Drying Alcohols to Avoid |
| 5–23 | Everything else |

Don't change ranks unless you have a good reason — the ordering is intentional.

### `BlogUrl`
Optional. A link to a blog post or resource about this category. Used in glossary pages.

### `Ingredients`
The list of ingredient strings to match. This is what you'll edit most often.

---

## How matching works

Understanding how ingredients are matched prevents common mistakes.

### The `normalize()` function

Before comparing, every ingredient string is passed through a normalization step:

1. Remove anything in `[square brackets]` (e.g., `[preservative]` → gone)
2. Remove anything in `(parentheses)` (e.g., `(fragrance)` → gone)
3. Fold accented characters to their base form (`Óleo` → `Oleo`)
4. Remove formatting characters — spaces, hyphens, slashes, periods, etc. Unicode letters and numbers are preserved.
5. Lowercase everything

So `"Sodium Lauryl Sulfate"` becomes `sodiumlaurylsulfate`.

An ingredient from a product label is normalized the same way before matching.

**Matching is exact after normalization** — not substring. An ingredient must normalize to exactly the same string as a rule entry to match.

### Case doesn't matter at all

All of these normalize to `sodiumlaurylsulfate` and are treated as **identical**:

- `"Sodium Lauryl Sulfate"`
- `"sodium lauryl sulfate"`
- `"SODIUM LAURYL SULFATE"`
- `"sodium lauryl Sulfate"`

Adding the same ingredient in multiple cases is redundant. Pick one form (Title Case is conventional for multi-word names) and use it once.

### Spelling variants ARE different

This is the most important distinction:

| Ingredient | Normalized form |
|-----------|----------------|
| `"Sodium Lauryl Sulfate"` | `sodiumlaurylsulfate` |
| `"Sodium Laurel Sulfate"` | `sodiumlaurelsulfate` |

These are genuinely different strings. Both spellings appear on real product labels. You need **separate entries** for each spelling variation.

Similarly:
- `"sulfate"` vs `"sulphate"` — different (both needed)
- `"Dimethicone"` vs `"dimethiconol"` — different
- `"peg-40"` vs `"peg40"` — after normalization, these ARE the same (hyphens are stripped)

### Slashes

If a product lists an ingredient as `"water/aqua"`, the app splits on `/` and checks each part separately. So `"water"` and `"aqua"` would be checked independently. You don't need to add slash-combined variants.

### How matches become a result

Each submitted ingredient is checked against every rule. A single ingredient may therefore appear in more than one category panel when it has multiple relevant classifications.

The headline product result is determined in this order:

1. `danger` if any ingredient matches a danger rule
2. `warning` if there is no danger match and any ingredient matches a warning rule
3. `unknown` if there is no danger or warning match and any submitted ingredient is unmatched
4. `good` if every submitted ingredient is recognized and none matches a danger or warning rule

`success` is treated the same as `good`. Unknown ingredients are displayed separately and are never assumed to be safe. A danger or warning result remains the headline when unknown ingredients are also present, and the unknown ingredients are still listed for the user.

Matched category panels are sorted by `Rank`; rank controls presentation, not whether a rule participates in matching.

---

## Common pitfalls

### Don't add partial names

The matching is not substring-based. Adding `"alcohol"` to the danger list will **not** match `"benzyl alcohol"` — it would only match a product that lists exactly `"alcohol"` (after normalization: `alcohol`).

If you want to match `"benzyl alcohol"`, add `"benzyl alcohol"` as its own entry.

### Check for duplicates before adding

Search `ingredientrules.json` for your ingredient before adding it. It may already exist:
- In the same rule (case-only or formatting-only duplicate — redundant, remove it)
- In a different rule (possibly an intentional second classification)

An ingredient may appear in multiple rules when each category accurately describes it. Results will show every applicable category, and the most severe match contributes to the headline verdict. When adding a second classification, include a source or explanation so reviewers can distinguish an intentional overlap from a mistake.

### Don't add overly generic terms

Adding `"extract"` or `"oil"` would match nearly everything. Be specific.

### Test your change

After submitting a pull request, check the automated test results. A failing test usually means a JSON syntax error (missing comma, extra comma, mismatched bracket). GitHub's editor will sometimes highlight these for you.

---

## The 24 rules at a glance

| Name | Result | Rank |
|------|--------|------|
| Sulfates to Avoid | danger | 1 |
| Silicones to Avoid | danger | 2 |
| Waxes and Oils to Avoid | danger | 3 |
| Drying Alcohols to Avoid | danger | 4 |
| Silicones to Use with Caution | warning | 5 |
| Proteins | warning | 6 |
| Humectants | warning | 7 |
| Polyquats | warning | 8 |
| Oils to use with Caution | warning | 9 |
| Conditioning Agents | warning | 10 |
| Parabens | warning | 11 |
| Itchy M's | warning | 12 |
| Chelating Agents | warning | 13 |
| Clarifying Ingredients | warning | 14 |
| Cleansing Agents | warning | 15 |
| Gentle Cleansers | good | 16 |
| Film Forming Humectants | good | 17 |
| Water Soluble Silicones | good | 18 |
| Moisturizing Alcohols | good | 19 |
| Waxes | good | 20 |
| Botanical Ingredients - Extracts and Oils | good | 21 |
| Preservatives | good | 22 |
| Verified - Okay | good | 22 |
| Verified - Caution | warning | 23 |
