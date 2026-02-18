# Admin Guide

Thank you for volunteering to help maintain isitcg.com. This guide explains what the role involves and how to handle the most common situations.

You don't need to know how to code. The automated tests handle JSON validation. Your job is to bring your CG knowledge.

---

## Your role

The maintainer of this repo is a developer, not a CG expert. When someone proposes adding or changing an ingredient, **you know better than they do** whether it belongs where they put it.

Your responsibilities:
- Review pull requests that add or change ingredients
- Triage issues (confirm them, label them, ask for more info, close the ones that don't make sense)
- Be the human expert in the loop

You do **not** need to:
- Write code
- Fix JSON syntax errors (the CI tests will flag those automatically)
- Approve every PR immediately — it's fine to take your time and ask questions

---

## Reviewing a pull request

When someone opens a pull request to change `ingredientrules.json`, ask yourself:

**1. Is the ingredient name correct?**
Check the spelling against a real product label or a reliable CG resource. Ingredient names on labels can vary — "Sodium Lauryl Sulfate" vs "Sodium Laurel Sulfate" are genuinely different spellings that appear on different products. Both may be worth including.

**2. Is it in the right category?**
Does this ingredient belong under "Sulfates to Avoid"? "Gentle Cleansers"? "Silicones to Avoid"? Use your CG knowledge here. The category list is in [ingredientrules.json](../ingredientrules.json).

**3. Is there a source?**
The PR template asks contributors to link to a source. It doesn't have to be a scientific study — a product label, a known CG blog, or a well-established community resource all count. If there's no source at all, it's reasonable to ask for one before approving.

**4. Is it a duplicate?**
Search the PR diff for the ingredient name. Also do a quick search of `ingredientrules.json` to see if the same ingredient (or a near-identical spelling) already exists in another category.

> **Case note:** Matching is case-insensitive — "Sodium Lauryl Sulfate" and "sodium lauryl sulfate" are treated identically by the tool. If a PR adds a case-only duplicate of something already in the database, the new entry is redundant and can be removed.

**5. Did the CI tests pass?**
You'll see a green checkmark or red X on the PR. If the tests failed, the JSON probably has a syntax error. You don't need to fix it — just leave a comment like "The tests are failing, there may be a formatting issue with the JSON. Can you take a look?"

---

## Approving vs. requesting changes

**Approve** when:
- The ingredient name looks correct
- The category makes sense
- There's a reasonable source or explanation
- The CI tests pass

**Request changes** when:
- The ingredient seems wrong or misplaced
- No source is provided and you want one
- There's a duplicate you'd like removed

**Close without merging** when:
- The change is clearly incorrect and the contributor hasn't responded to requests for clarification
- The PR is spam or irrelevant

To approve: click "Review changes" → "Approve" → "Submit review."
To request changes: click "Review changes" → "Request changes" → write your comments → "Submit review."

---

## Triaging issues

When someone opens an issue, you'll see it in the [issues list](https://github.com/josephgardner/isitcg/issues).

**Common issue types and how to handle them:**

| Issue type | What to do |
|-----------|-----------|
| "This ingredient is missing" | Confirm whether it should be added. If yes, label it `good first issue`. If you can add it yourself, do so. |
| "This ingredient is in the wrong category" | Verify with your CG knowledge. If correct, label it `needs-verification` until you can confirm, then act on it. |
| "The tool gave me a wrong result" | Investigate — this is often a missing or misspelled ingredient. |
| "I want to help" | Welcome them! Point them to [CONTRIBUTING.md](../CONTRIBUTING.md) and this guide. |
| Spam or off-topic | Close it. |

**Useful labels:**
- `good first issue` — simple, well-defined, low risk
- `needs-verification` — needs a CG expert to confirm before acting
- `help wanted` — open to anyone who wants to contribute

---

## The CI pipeline

Every pull request automatically runs tests before it can be merged. You'll see the status on the PR page.

**Green checkmark** — tests passed. The JSON is valid and the known test cases still work. Safe to review the content.

**Red X** — tests failed. This almost always means the JSON has a syntax error (missing comma, extra comma, unclosed bracket). You don't need to fix this. Just leave a comment pointing it out — the contributor can fix it by editing the file again on GitHub.

Tests run on pull requests and on pushes to `main`. When a PR is merged, the site automatically deploys to isitcg.com.
