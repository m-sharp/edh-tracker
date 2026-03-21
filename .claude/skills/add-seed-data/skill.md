---
name: add-seed-data
description: Use this skill to add new game results to data/gameInfos.json. TRIGGER when the user provides game results to record (players, commanders, kills, places, date). DO NOT TRIGGER for questions about existing data, schema changes, seeder code changes, or anything not involving appending new game records.
version: 1.0.0
---

# Add Seed Data — gameInfos.json

Use this skill when the user provides new game results to append to `data/gameInfos.json`.

## Step 1 — Read the end of the file

Read the tail of `data/gameInfos.json` to see the current last entry and confirm where to append. 
The file is large; use `offset`/`limit` or `tail` rather than reading the whole thing.

## Step 2 — Check player and commander names

Before writing anything:

- **Player names**: grep existing data for any player name you haven't seen before to confirm spelling.
- **Commander names**: scan the user's input for likely typos. Common pitfalls:
  - Homophones or near-homophones (e.g. "Road" vs "Roar", "Warper" vs "Wraper")
  - Apostrophes and punctuation (e.g. "Praetor's" not "Praetors")
  - Articles and prepositions (e.g. "of the", "the", "a")
  - Confirm any name you're uncertain about before writing via the AskUserQuestion tool
- If you correct a typo, call it out explicitly in your response.

## Step 3 — Append the new entries

Edit `data/gameInfos.json` to replace the closing `]` with the new game objects followed by `]`.

### JSON structure

Each game is an object with:

```json
{
  "date": "YYYY-MM-DDT00:00:00Z",
  "results": [
    {
      "player": "PlayerName",
      "place": 1,
      "kills": 2,
      "name": "Commander Name"
    }
  ],
  "format": "commander"
}
```

### Field rules

| Field | Type | Notes |
|---|---|---|
| `date` | string | ISO 8601, always midnight UTC: `"YYYY-MM-DDT00:00:00Z"` |
| `player` | string | Must match existing player name exactly (case-sensitive) |
| `place` | int | Finishing position; **ties are valid** — multiple players can share the same place |
| `kills` | int | **Omit entirely when 0** (do not write `"kills": 0`) |
| `name` | string | Commander name exactly as it appears on the card |
| `format` | string | Always `"commander"` unless told otherwise |

### Result ordering

List results in the order the user provided them. No sorting is required.

## Step 4 — Confirm

After editing, briefly summarize what was added and call out any typos corrected or ambiguities noted.
