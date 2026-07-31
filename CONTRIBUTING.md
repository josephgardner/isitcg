# Contributing to isitcg.com

Thank you for wanting to help! This guide walks you through making a contribution entirely through GitHub's website — no terminal, no git commands required.

You'll need a free [GitHub account](https://github.com/join).

---

## Before you start

A good contribution:
- Is based on a reliable source (a product label, a well-known CG resource, a blog post, a study)
- Includes a note about where the information comes from
- Checks that the ingredient isn't already in the database under a slightly different spelling

The most common mistakes are:
- Adding an ingredient that's already there (search the file first!)
- Putting an ingredient in the wrong category
- Adding a generic word that would accidentally match too many things

When in doubt, [open an issue](https://github.com/josephgardner/isitcg/issues/new/choose) and describe what you found — someone can help figure out the right approach.

---

## How to suggest an ingredient

### Step 1 — Open `ingredientrules.json`

Go to [ingredientrules.json](https://github.com/josephgardner/isitcg/blob/main/ingredientrules.json) on GitHub.

### Step 2 — Search for the ingredient first

Press `Ctrl+F` (or `Cmd+F` on Mac) and search for the ingredient name. If it's already in the intended category, you're done. An ingredient can intentionally belong to more than one category when it has multiple relevant properties, but a second classification should have a source or explanation.

### Step 3 — Click the pencil icon

In the top-right corner of the file view, click the pencil (edit) icon. GitHub will create a personal copy of the file for you to edit.

### Step 4 — Find the right rule

The file is a list of rules. Each rule has a name like `"Sulfates to Avoid"` or `"Gentle Cleansers"`. Scroll to find the rule or rules where your ingredient belongs.

See the [ingredient format reference](guides/ingredient-format.md) if you're not sure which rule to use or how the file is structured.

### Step 5 — Add the ingredient

Inside the rule's `"Ingredients"` array, add your ingredient as a new line. Follow the same format as the surrounding entries — quoted string, comma at the end (except the last item in the list).

**Example:** adding `"sodium phytate"` to a list:

```json
"Ingredients": [
  "disodium EDTA",
  "sodium citrate",
  "sodium phytate"
]
```

**Important:** Matching is case-insensitive — `"Sodium Phytate"` and `"sodium phytate"` are treated as identical. Don't add the same ingredient in multiple cases within one rule. Pick whichever looks most natural (Title Case is conventional for multi-word names).

Different *spellings* do need separate entries — `"Sodium Lauryl Sulfate"` and `"Sodium Laurel Sulfate"` are genuinely different strings that appear on different product labels.

### Step 6 — Commit your change

Scroll to the bottom of the editor. You'll see a "Commit changes" section.

Write a short description of what you changed and why, for example:

> Add sodium phytate to Chelating Agents — it's a chelating agent used in many CG-friendly products

Click **"Propose changes"**. This creates a pull request (a proposed change) for maintainers to review.

---

## The ingredient rules format

Here's an annotated example of one rule:

```json
{
  "Name": "Sulfates to Avoid",
  "Description": "These ingredients are sulfates that should be avoided.",
  "Result": "danger",
  "Rank": 1,
  "Ingredients": [
    "sodium lauryl sulfate",
    "sodium laureth sulfate",
    "ammonium lauryl sulfate"
  ]
}
```

| Field | What it means |
|-------|--------------|
| `Name` | The category name shown in results |
| `Description` | Shown below the category name |
| `Result` | `danger`, `warning`, or `good` (also `success`, which is an older alias for `good`) |
| `Rank` | Lower number = shown higher in results. Danger categories are ranked 1–4. |
| `Ingredients` | The list of ingredient strings to match against |

The analyzer checks each submitted ingredient against every rule. If an ingredient belongs to multiple rules, it appears in every applicable category panel. The overall product verdict is the most severe result found: `danger`, then `warning`, then `good`. If no danger or warning is found but one or more submitted ingredients are unmatched, the product is shown as **Unable to Verify** rather than CG-approved.

For a deeper explanation of how matching works and common pitfalls, see [guides/ingredient-format.md](guides/ingredient-format.md).

---

## How to improve a glossary page

The site has a glossary page for each ingredient category. Most are stubs. You can flesh them out with descriptions, context, tips, or links to helpful resources.

### Step 1 — Find the right file

Browse to [`docs/glossary/`](https://github.com/josephgardner/isitcg/tree/main/docs/glossary) and click the file for the category you want to improve.

### Step 2 — Click the pencil icon

Same as before — click the pencil icon in the top right to start editing.

### Step 3 — Write something helpful

Add information that would help a curly hair enthusiast understand the category. For example:
- What these ingredients do
- Why they're considered CG or not-CG
- Common products that contain them
- Links to CG resources or studies

### Step 4 — Commit your change

Same process as above — scroll down, write a commit message, click "Propose changes."

---

## What happens next

1. Your pull request appears in the [PR list](https://github.com/josephgardner/isitcg/pulls)
2. GitHub automatically runs tests to check that the JSON is valid
3. A maintainer or volunteer admin will review the change — usually within a few days
4. If everything looks good, it gets merged and the site automatically updates within minutes

If the reviewer asks for changes, you'll get a notification and can make edits directly on the pull request page.
