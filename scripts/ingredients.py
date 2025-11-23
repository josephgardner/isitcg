#!/usr/bin/env python3
"""
Ingredient suggestions management.

Usage:
  python scripts/ingredients.py update [--dry-run]  # Create/update PRs for pending ingredients
"""

import json, os, re, subprocess, sys
import requests

API_URL = "https://analytics.apps.algoplay.com/api/ingredients/pending"
RULES_FILE = "ingredientrules.json"

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

def category_prs():
    """Get existing PRs mapped by category slug."""
    result = gh("pr", "list", "--json", "number,headRefName", "--limit", "500")
    if result.returncode != 0: return {}
    return {pr["headRefName"].replace("ingredient-suggestions/", ""): pr["number"]
            for pr in json.loads(result.stdout)
            if pr["headRefName"].startswith("ingredient-suggestions/")}

def pr_last_author(pr_num):
    """Get the login of the last commit author on a PR."""
    result = gh("pr", "view", str(pr_num), "--json", "commits")
    if result.returncode != 0: return None
    commits = json.loads(result.stdout).get("commits", [])
    if not commits: return None
    authors = commits[-1].get("authors", [])
    return authors[0].get("login") if authors else None


def cmd_update(dry_run=False):
    """Create/update PRs for pending ingredients."""
    api_key = os.environ.get("ANALYTICS_API_KEY")
    if not api_key: raise ValueError("ANALYTICS_API_KEY not set")

    data = requests.get(API_URL, headers={"X-API-Key": api_key}).json()
    existing = existing_ingredients()
    pending = [i for i in data.get("ingredients", []) if i["name"].lower() not in existing]

    # Group ingredients by their top suggested category
    by_category = {}
    for item in pending:
        name = item['name']
        count = item.get('count', 0)
        suggestions = item.get("suggested_categories", [])
        if not suggestions:
            continue
        # Use first (highest confidence) suggestion for grouping
        top = suggestions[0]
        cat = top["category"]
        by_category.setdefault(cat, []).append({
            "name": name,
            "count": count,
            "suggestions": suggestions  # Keep all suggestions for PR body
        })

    if not by_category:
        print("No pending ingredients to process")
        return

    print(f"Found {sum(len(v) for v in by_category.values())} ingredients in {len(by_category)} categories")

    if dry_run:
        for cat, ingredients in sorted(by_category.items()):
            print(f"\n{cat}:")
            for ing in sorted(ingredients, key=lambda x: -x["count"]):
                top = ing["suggestions"][0]
                conf = int(top.get("confidence", 0) * 100)
                print(f"  - {ing['name']} ({ing['count']}x, {conf}%) - {top.get('reason', '')}")
                for alt in ing["suggestions"][1:]:
                    alt_conf = int(alt.get("confidence", 0) * 100)
                    print(f"      ⊕ {alt['category']} ({alt_conf}%) - {alt.get('reason', '')}")
        return

    git("config", "user.name", "github-actions[bot]")
    git("config", "user.email", "github-actions[bot]@users.noreply.github.com")

    prs = category_prs()

    # Create/update PRs for each category
    for cat, ingredients in by_category.items():
        slug = slugify(cat)
        branch = f"ingredient-suggestions/{slug}"
        pr_num = prs.get(slug)

        # Check if PR was manually edited
        if pr_num:
            last_author = pr_last_author(pr_num)
            if last_author and last_author != "github-actions[bot]":
                print(f"Skipping {cat} - manually edited by {last_author}")
                continue

        git("fetch", "origin", "main")
        git("checkout", "-B", branch, "origin/main")

        rules = load_rules()
        ingredient_names = [ing["name"] for ing in ingredients]
        for rule in rules["Rules"]:
            if rule["Name"] == cat:
                rule["Ingredients"] = sorted(set(rule["Ingredients"] + ingredient_names), key=str.lower)
                break
        save_rules(rules)

        git("add", RULES_FILE)
        if git("diff", "--cached", "--quiet").returncode == 0:
            print(f"No changes for {cat}")
            continue

        git("commit", "-m", f"Add ingredients to {cat}")
        git("push", "-f", "origin", branch)

        # Build PR body with ingredient details
        ing_lines = []
        for ing in sorted(ingredients, key=lambda x: -x["count"]):
            top = ing["suggestions"][0]
            conf = int(top.get("confidence", 0) * 100)
            ing_lines.append(f"- **{ing['name']}** ({ing['count']}x, {conf}%) - {top.get('reason', '')}")
            for alt in ing["suggestions"][1:]:
                alt_conf = int(alt.get("confidence", 0) * 100)
                ing_lines.append(f"  - ⊕ {alt['category']} ({alt_conf}%) - {alt.get('reason', '')}")
        ing_list = "\n".join(ing_lines)

        if not pr_num:
            result = gh("pr", "create",
                "--title", f"Add ingredients: {cat}",
                "--body", f"Adding to **{cat}**:\n\n{ing_list}",
                "--base", "main",
                "--head", branch)
            print(f"Created PR for {cat}: {result.stdout.strip()}")
        else:
            # Update existing PR body
            gh("pr", "edit", str(pr_num), "--body", f"Adding to **{cat}**:\n\n{ing_list}")
            print(f"Updated PR #{pr_num} for {cat}")

    git("checkout", "main")


if __name__ == "__main__":
    cmd = sys.argv[1] if len(sys.argv) > 1 else None
    dry_run = "--dry-run" in sys.argv
    if cmd == "update": cmd_update(dry_run)
    else: print(__doc__)
