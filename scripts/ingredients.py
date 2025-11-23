#!/usr/bin/env python3
"""
Ingredient suggestions management.

Usage:
  python scripts/ingredients.py update  # Update dashboard issue
  python scripts/ingredients.py sync    # Sync PRs with checkboxes
"""

import json, os, re, subprocess, sys
import requests

API_URL = "https://analytics.apps.algoplay.com/api/ingredients/pending"
RULES_FILE = "ingredientrules.json"
ISSUE_TITLE = "Ingredient Suggestions Dashboard"

def gh(*args):
    return subprocess.run(["gh", *args], capture_output=True, text=True)

def git(*args):
    return subprocess.run(["git", *args], capture_output=True, text=True)

def slugify(s):
    return re.sub(r'[^a-z0-9-]', '-', s.lower()).strip('-')

def load_rules():
    return json.load(open(RULES_FILE))

def save_rules(data):
    json.dump(data, open(RULES_FILE, "w"), indent=2)

def existing_ingredients():
    return {ing.lower() for rule in load_rules()["Rules"] for ing in rule["Ingredients"]}

def find_issue():
    result = gh("issue", "list", "--search", f'"{ISSUE_TITLE}" in:title', "--json", "number,title")
    return next((i["number"] for i in json.loads(result.stdout) if i["title"] == ISSUE_TITLE), None)

def category_prs():
    result = gh("pr", "list", "--json", "number,headRefName", "--limit", "500")
    if result.returncode != 0: return {}
    return {pr["headRefName"].replace("ingredient-suggestions/", ""): pr["number"]
            for pr in json.loads(result.stdout)
            if pr["headRefName"].startswith("ingredient-suggestions/")}


def cmd_update():
    """Update the dashboard issue with pending ingredients."""
    api_key = os.environ.get("ANALYTICS_API_KEY")
    if not api_key: raise ValueError("ANALYTICS_API_KEY not set")

    data = requests.get(API_URL, headers={"X-API-Key": api_key}).json()
    existing = existing_ingredients()
    pending = [i for i in data.get("ingredients", []) if i["name"].lower() not in existing]
    prs = category_prs()

    # Build issue body
    lines = [
        "# Ingredient Suggestions\n",
        "Check a box to include that ingredient in a PR. **Only check ONE category per ingredient.**\n",
        "---\n"
    ]

    for item in sorted(pending, key=lambda x: -x.get("count", 0)):
        lines.append(f"### {item['name']}\n*Seen {item.get('count', 0)} times*\n")

        for sug in item.get("suggested_categories", []):
            cat = sug["category"]
            conf = int(sug.get("confidence", 0) * 100)
            pr_num = prs.get(slugify(cat))
            check = "[x]" if pr_num else "[ ]"
            link = f" [PR #{pr_num}](../../pull/{pr_num})" if pr_num else ""

            lines.append(f"- {check} **{cat}** ({conf}%){link}")
            lines.append(f"  - {sug.get('reason', '')}\n")

        lines.append("---\n")

    body = "\n".join(lines) if pending else "*No pending ingredients.*"

    issue_num = find_issue()
    if issue_num:
        gh("issue", "edit", str(issue_num), "--body", body)
        print(f"Updated issue #{issue_num}")
    else:
        result = gh("issue", "create", "--title", ISSUE_TITLE, "--body", body)
        print(f"Created: {result.stdout.strip()}")


def cmd_sync():
    """Sync PRs with checkbox state."""
    issue_num = find_issue()
    if not issue_num: return print("Dashboard issue not found")

    body = json.loads(gh("issue", "view", str(issue_num), "--json", "body").stdout)["body"]

    # Parse checked items grouped by category
    by_category = {}
    current_ing = None
    for line in body.split("\n"):
        if m := re.match(r"^### (.+)$", line):
            current_ing = m.group(1).strip()
        elif (m := re.match(r"^- \[x\] \*\*(.+?)\*\*", line)) and current_ing:
            by_category.setdefault(m.group(1).strip(), []).append(current_ing)

    print(f"Found {sum(len(v) for v in by_category.values())} ingredients in {len(by_category)} categories")

    git("config", "user.name", "github-actions[bot]")
    git("config", "user.email", "github-actions[bot]@users.noreply.github.com")

    prs = category_prs()

    # Create/update PRs
    for cat, ingredients in by_category.items():
        slug = slugify(cat)
        branch = f"ingredient-suggestions/{slug}"

        git("fetch", "origin", "main")
        git("checkout", "-B", branch, "origin/main")

        rules = load_rules()
        for rule in rules["Rules"]:
            if rule["Name"] == cat:
                rule["Ingredients"] = sorted(set(rule["Ingredients"] + ingredients), key=str.lower)
                break
        save_rules(rules)

        git("add", RULES_FILE)
        if git("diff", "--cached", "--quiet").returncode == 0:
            continue

        git("commit", "-m", f"Add ingredients to {cat}")
        git("push", "-f", "origin", branch)

        if slug not in prs:
            ing_list = "\n".join(f"- {i}" for i in sorted(ingredients))
            gh("pr", "create", "--title", f"Add ingredients: {cat}",
               "--body", f"Adding to **{cat}**:\n\n{ing_list}", "--base", "main", "--head", branch)
            print(f"Created PR for {cat}")
        else:
            print(f"Updated PR #{prs[slug]}")

    # Close PRs for unchecked categories
    for slug, pr_num in prs.items():
        if slug not in {slugify(c) for c in by_category}:
            gh("pr", "close", str(pr_num), "--delete-branch")
            print(f"Closed PR #{pr_num}")

    git("checkout", "main")


if __name__ == "__main__":
    cmd = sys.argv[1] if len(sys.argv) > 1 else None
    if cmd == "update": cmd_update()
    elif cmd == "sync": cmd_sync()
    else: print(__doc__)
