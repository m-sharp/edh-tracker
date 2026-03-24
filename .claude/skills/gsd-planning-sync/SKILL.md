---
name: gsd-planning-sync
description: Sync and commit .planning/ artifacts after any GSD operation. TRIGGER automatically after any GSD workflow completes — /gsd:execute-phase, /gsd:plan-phase, /gsd:quick, /gsd:debug, /gsd:verify-work, or any other /gsd:* command — to ensure .planning/ changes are committed and STATE.md/ROADMAP.md reflect the latest status. DO NOT TRIGGER for non-GSD operations or when the user is mid-workflow and explicitly continuing to the next step.
version: 1.0.0
---

# GSD Planning Sync

Run this after any GSD operation completes to ensure `.planning/` changes are committed
and `STATE.md` / `ROADMAP.md` reflect the current project state.

## Step 1 — Check for uncommitted .planning/ changes

```bash
git status --short -- .planning/
```

If the output is empty, the working tree is clean. **Skip to Step 3.**

## Step 2 — Validate and stage .planning/ changes

For each modified or untracked file under `.planning/`, verify it is in the correct state:

**STATE.md checks:**
- `status:` reflects the outcome of the operation that just ran (e.g., "Executing Phase X", "Phase X execution complete", "Planning Phase X")
- `stopped_at:` references the last completed plan or operation
- `progress.completed_plans` matches the actual count of plans with a SUMMARY.md on disk
- `Current Position` section shows the correct phase and plan number

**ROADMAP.md checks:**
- Any plan that now has a SUMMARY.md is marked `[x]` (not `[ ]`)
- Any plan that was just created is listed under its phase section
- `**Plans:** N plans` count matches the actual plan files on disk

**PLAN/SUMMARY/VERIFICATION files:**
- These are generally correct as written by the executor/verifier agent — commit them as-is unless there is an obvious error (e.g., wrong phase number in frontmatter)

If a file has stale or incorrect data, correct it before staging.

## Step 3 — Ensure STATE.md and ROADMAP.md are current

Even if Step 1 found no uncommitted changes, verify that both files reflect the latest operation:

**STATE.md must have:**
- `last_updated:` within the current session (not hours/days stale)
- `status:` that accurately describes where the project stands right now
- `stopped_at:` referencing the actual last completed artifact

**ROADMAP.md must have:**
- All completed plans (those with a SUMMARY.md on disk) marked `[x]`

To find the ground truth:
```bash
# Count completed plans per phase
for phase_dir in .planning/phases/*/; do
  phase=$(basename "$phase_dir")
  summaries=$(ls "$phase_dir"*-SUMMARY.md 2>/dev/null | wc -l)
  plans=$(ls "$phase_dir"*-PLAN.md 2>/dev/null | wc -l)
  echo "$phase: $summaries/$plans plans complete"
done
```

If either file is stale, update it directly using Edit before committing.

## Step 4 — Commit all pending .planning/ changes

Use the gsd-tools commit helper so the commit follows GSD conventions:

```bash
node "$HOME/.claude/get-shit-done/bin/gsd-tools.cjs" commit \
  "docs: sync .planning/ artifacts after GSD operation" \
  --files $(git ls-files --modified --others --exclude-standard -- .planning/ | tr '\n' ' ')
```

If gsd-tools is unavailable, fall back to a direct git commit:

```bash
git add .planning/
git commit --no-verify -m "docs: sync .planning/ artifacts after GSD operation

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
```

## Step 5 — Confirm

```bash
git status --short -- .planning/
```

Output should be empty. Report: "`.planning/` is clean — STATE.md and ROADMAP.md are up to date."

---

## What counts as "after a GSD operation"

Trigger this skill after:
- `/gsd:execute-phase` — after all waves complete (or after gap-closure run)
- `/gsd:plan-phase` — after PLAN.md files are written
- `/gsd:quick` — after the inline task completes
- `/gsd:debug` — after a debug session concludes
- `/gsd:verify-work` — after verification produces a UAT or VERIFICATION.md
- Any other `/gsd:*` command that produces or modifies files under `.planning/`

Do NOT trigger mid-workflow when the user is about to chain to the next GSD step immediately — wait until the chain is done or the user pauses.
