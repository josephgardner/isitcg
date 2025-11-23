#!/usr/bin/env python3
"""
Ingredient suggestions management.

Usage:
  python scripts/ingredients.py update [--dry-run]  # Update dashboard issue
  python scripts/ingredients.py sync [--dry-run]    # Sync PRs with checkboxes
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


def cmd_update(dry_run=False):
    """Update the dashboard issue with pending ingredients."""
    api_key = os.environ.get("ANALYTICS_API_KEY")
    if not api_key: raise ValueError("ANALYTICS_API_KEY not set")

    data = requests.get(API_URL, headers={"X-API-Key": api_key}).json()
    existing = existing_ingredients()
    pending = [i for i in data.get("ingredients", []) if i["name"].lower() not in existing]
    prs = category_prs() if not dry_run else {}

    # Group ingredients by category and track which categories each ingredient appears in
    by_category = {}
    ing_categories = {}  # ingredient name -> list of categories
    for item in pending:
        name = item['name']
        count = item.get('count', 0)
        for sug in item.get("suggested_categories", []):
            cat = sug["category"]
            conf = int(sug.get("confidence", 0) * 100)
            reason = sug.get('reason', '')
            by_category.setdefault(cat, []).append({
                "name": name, "count": count, "conf": conf, "reason": reason
            })
            ing_categories.setdefault(name, []).append(cat)

    # Sort categories by total ingredient count
    sorted_cats = sorted(by_category.items(), key=lambda x: -sum(i["count"] for i in x[1]))

    # Build issue body
    lines = [
        "# Ingredient Suggestions\n",
        "Check ingredients to add to each category.\n",
    ]

    for cat, ingredients in sorted_cats:
        pr_num = prs.get(slugify(cat))
        pr_link = f" ([PR #{pr_num}](../../pull/{pr_num}))" if pr_num else ""
        lines.append(f"## {cat}{pr_link}\n")

        for ing in sorted(ingredients, key=lambda x: -x["count"]):
            check = "[x]" if pr_num else "[ ]"
            # Add links to other categories this ingredient appears in
            other_cats = [c for c in ing_categories[ing['name']] if c != cat]
            alt_links = " ⊕ " + ", ".join(f"[{c}]" for c in other_cats) if other_cats else ""
            lines.append(f"- {check} **{ing['name']}** ({ing['count']}x, {ing['conf']}%) - {ing['reason']}{alt_links}")

        lines.append("")

    # Add reference-style links at the end
    for cat, _ in sorted_cats:
        lines.append(f"[{cat}]: #{slugify(cat)}")

    body = "\n".join(lines) if pending else "*No pending ingredients.*"

    if dry_run:
        print(body)
        return

    issue_num = find_issue()
    if issue_num:
        gh("issue", "edit", str(issue_num), "--body", body)
        print(f"Updated issue #{issue_num}")
    else:
        result = gh("issue", "create", "--title", ISSUE_TITLE, "--body", body)
        print(f"Created: {result.stdout.strip()}")


def cmd_sync(dry_run=False):
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

    if dry_run:
        for cat, ingredients in by_category.items():
            print(f"\n{cat}:")
            for ing in sorted(ingredients):
                print(f"  - {ing}")
        return

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
    dry_run = "--dry-run" in sys.argv
    if cmd == "update": cmd_update(dry_run)
    elif cmd == "sync": cmd_sync(dry_run)
    else: print(__doc__)
