---
phase: quick
plan: 260327-sli
type: execute
wave: 1
depends_on: []
files_modified:
  - .planning/PROJECT.md
  - .planning/REQUIREMENTS.md
  - .planning/ROADMAP.md
autonomous: true
requirements: []
must_haves:
  truths:
    - "Every requirement marked Complete in REQUIREMENTS.md traceability table has its corresponding PROJECT.md bullet checked"
    - "PROJECT.md footer reflects Phase 05 completion date"
    - "ROADMAP.md Phase 1 and Phase 2 checkboxes are marked complete"
    - "PROJECT.md Key Decisions table reflects implemented outcomes"
  artifacts:
    - path: ".planning/PROJECT.md"
      provides: "Updated active requirements checklist and key decisions"
    - path: ".planning/REQUIREMENTS.md"
      provides: "Updated last-updated footer"
    - path: ".planning/ROADMAP.md"
      provides: "Phase 1 and 2 completion checkboxes"
  key_links: []
---

<objective>
Update planning documents (PROJECT.md, REQUIREMENTS.md, ROADMAP.md) to accurately reflect the completion status of Phases 1-5.

Purpose: Planning docs have stale unchecked items for work completed in Phases 1-5. This creates confusion about what is actually done vs pending.
Output: Accurate planning documents with all completed work properly checked off.
</objective>

<execution_context>
@$HOME/.claude/get-shit-done/workflows/execute-plan.md
@$HOME/.claude/get-shit-done/templates/summary.md
</execution_context>

<context>
@.planning/PROJECT.md
@.planning/REQUIREMENTS.md
@.planning/ROADMAP.md
</context>

<tasks>

<task type="auto">
  <name>Task 1: Update PROJECT.md to reflect Phase 1-5 completions</name>
  <files>.planning/PROJECT.md</files>
  <action>
Cross-reference REQUIREMENTS.md traceability table (Complete status) against PROJECT.md Active bullets. Make these specific changes:

**Active section - check off completed items:**

1. "Define and apply an overarching visual design language" line 27 — check it (DSNG-01 complete, Phase 2)
2. "Game model change" section (lines 39-41) — all three bullets should be checked:
   - "Remove player requirement from game entry" (GAME-01 complete, Phase 4)
   - "Deck picker in game form displays owner name" (GAME-02 complete, Phase 4)
   - "Remove/hide player field from game creation and result forms" (GAME-01 complete, Phase 4)
3. "New game form complete redesign" line 47 — check it (GAME-03 complete, Phase 4)
4. "Tooltip on deck commander update" line 49 — check it (DECK-02 complete, Phase 5)
5. "Investigate and define retired deck behavior" line 50 — check it (DECK-03 complete, Phase 5)
6. Backend correctness section (lines 56-62) — check ALL seven bullets:
   - Pod-membership check on POST /api/game (SEC-01, Phase 1)
   - player_id-from-context ownership check in DeckCreate (SEC-02, Phase 1)
   - Wrap game creation in DB transaction (SEC-03, Phase 1)
   - Fix PromotePlayer/KickPlayer returning 403 for all errors (INFRA-02, Phase 1)
   - Add startup check rejecting JWT secrets shorter than 32 bytes (AUTH-02, Phase 1)
   - Validate used_count against max on pod invite join (SEC-04, Phase 1)
   - Add input length validation on string fields (SEC-05, Phase 1)
7. Performance section (lines 65-66) — check both bullets:
   - Batch deck stats queries (PERF-01, Phase 1)
   - Require at least one filter on GET /api/decks (PERF-02, Phase 1)

**DO NOT check these (still pending per REQUIREMENTS.md):**
- "Mobile-friendly layout and interaction patterns" (DSNG-02 still Pending)
- "Rebuild CLAUDE.md context section" (not tracked)
- "401 interceptor in http.ts" (AUTH-01, Phase 6 Pending)
- Test coverage items (TEST-01 through TEST-04, Phase 6 Pending)
- Production readiness items (INFRA-01, INFRA-03, Phase 7 Pending)

**Key Decisions table:** Update outcomes for implemented decisions:
- "Games track decks only, not players" — change "Pending" to "Implemented (Phase 4)"
- "Deck picker displays owner name" — change "Pending" to "Implemented (Phase 4)"
- "Frontend design language to be defined before implementation" — change "Pending" to "Implemented (Phase 2)"
- "Soft launch before full polish" — leave as Pending (not yet launched)

**Footer:** Change "Last updated: 2026-03-24 after Phase 03" to "Last updated: 2026-03-27 after Phase 05 (pod-deck-ux) completion"

Add "Validated in Phase N" annotations to newly checked items, matching the style of existing checked items (e.g., "— Validated in Phase 1: Backend Hardening").
  </action>
  <verify>
    <automated>grep -c "\- \[x\]" .planning/PROJECT.md</automated>
  </verify>
  <done>All items completed in Phases 1-5 are checked in PROJECT.md Active section. Key Decisions outcomes updated. Footer reflects Phase 05. Count of checked items increases from ~10 to ~22.</done>
</task>

<task type="auto">
  <name>Task 2: Update ROADMAP.md Phase 1/2 checkboxes and REQUIREMENTS.md footer</name>
  <files>.planning/ROADMAP.md, .planning/REQUIREMENTS.md</files>
  <action>
**ROADMAP.md:** Lines 15-16, the Phase 1 and Phase 2 entries still show `[ ]` while Phases 3-5 show `[x]`. Update:
- `- [ ] **Phase 1: Backend Hardening**` to `- [x] **Phase 1: Backend Hardening** - Close authorization, transaction, validation, and performance gaps in the API (completed 2026-03-23)`
- `- [ ] **Phase 2: Design Language**` to `- [x] **Phase 2: Design Language** - Define and apply the visual design system before any UI work begins (completed 2026-03-23)`

Note: Use approximate completion dates based on Phase 3 starting around 2026-03-23. If the exact dates are unclear, use the same format as Phase 3 but without a date (just mark `[x]`).

**REQUIREMENTS.md:** Update the footer from "Last updated: 2026-03-22 after roadmap creation" to "Last updated: 2026-03-27 after Phase 05 completion". No content changes needed — the checkboxes and traceability table are already accurate.
  </action>
  <verify>
    <automated>grep -E "^\- \[.\] \*\*Phase [12]:" .planning/ROADMAP.md && grep "Last updated" .planning/REQUIREMENTS.md</automated>
  </verify>
  <done>ROADMAP.md Phase 1 and Phase 2 lines show [x]. REQUIREMENTS.md footer reflects current date.</done>
</task>

</tasks>

<verification>
After both tasks:
1. `grep -c "\- \[ \]" .planning/PROJECT.md` — unchecked count should decrease (from ~18 to ~7)
2. `grep -c "\- \[x\]" .planning/PROJECT.md` — checked count should increase (from ~10 to ~22)
3. `grep "Phase 05" .planning/PROJECT.md` — footer mentions Phase 05
4. `grep -E "^\- \[x\] \*\*Phase [12]:" .planning/ROADMAP.md` — Phase 1 and 2 checked
5. No changes to STATE.md
</verification>

<success_criteria>
- Every requirement with "Complete" status in REQUIREMENTS.md traceability table has its corresponding PROJECT.md bullet checked off
- No still-pending requirements are incorrectly checked
- ROADMAP.md Phase 1 and 2 show [x] matching the Progress table
- All three doc footers reflect current date/phase
- Key Decisions table outcomes reflect implementation status
</success_criteria>

<output>
After completion, create `.planning/quick/260327-sli-update-planning-docs-to-reflect-phase-1-/260327-sli-SUMMARY.md`
</output>
