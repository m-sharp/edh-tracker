# Phase 3: Frontend Structure - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-23
**Phase:** 03-frontend-structure
**Areas discussed:** File org depth, Shared tab API, UI-SPEC step, Tooltip scope

---

## File Organization Depth

**Q: How granular should the split be within each view directory?**
Options: Move + extract tabs | Move only
Selected: Move + extract tabs — per-tab files extracted (DecksTab.tsx, PlayersTab.tsx, etc.)

**Q: Where should HomeView live after the restructure?**
Options: home/index.tsx | Stay in index.tsx
Selected: home/index.tsx — consistent with per-view pattern

**Q: Where should shared non-route utilities live?**
Options: Keep in app/src/ | Move to app/src/components/
Selected: Move to app/src/components/ — matches.tsx, stats.tsx, common.ts all move there along with new shared components

---

## Shared Tab API

**Q: What shape should the shared tab component take?**
Options: Config component (<TabbedLayout>) | useTabs hook
Selected: Config component — `<TabbedLayout queryKey="podTab" tabs={[{id, label, content, hidden}]} />`

**Q: What query string key naming convention?**
Options: Per-view keys | Single global key 'tab'
Selected: Per-view keys — `podTab`, `playerTab`, `deckTab`
User note: Distinct keys required; query string must not bloat — use replace:true, only write on explicit tab switch

**Q: Should TabbedLayout also own a shared loading/error state?**
Options: Yes — loading prop | No — each tab handles its own
Selected: Yes — loading prop with CircularProgress centered

**Q: Where does TabbedLayout live?**
Options: app/src/components/TabbedLayout.tsx | app/src/components/tabs/index.tsx
Selected: app/src/components/TabbedLayout.tsx

**Q: When TabbedLayout shows loading, what should it render?**
Options: CircularProgress centered | Skeleton rows | Claude's discretion
Selected: CircularProgress centered

**Q: Should TabbedLayout handle conditionally-hidden tabs?**
Options: Just 'hidden' prop | Also a 'disabled' state
Selected: Just 'hidden' — hidden: boolean filters out the tab

**Q: Should TabbedLayout use scrollable tabs on mobile?**
Options: Scrollable (variant="scrollable" scrollButtons="auto") | Claude's discretion
Selected: Scrollable — carries forward from existing pod.tsx

**Q: Should active tab be validated against available tabs on mount?**
Options: Yes — clamp to valid range | No — MUI handles it
Selected: Yes — clamp to first valid tab if out of range

**Q: Should TabbedLayout have an error prop?**
Options: Yes — error prop | No — loading only
Selected: No — loading only; errors handled at route level or inline

**Q: What query string key names per view?**
Options: podTab / playerTab / deckTab | pod_tab / player_tab / deck_tab
Selected: podTab / playerTab / deckTab — camelCase

**Q: Store as string label or numeric index?**
Options: String label ('decks', 'players') | Numeric index (0, 1, 2)
Selected: String label — robust to tab reordering

**Q: Auto-slug from label or explicit 'id' field?**
Options: Explicit 'id' field | Auto-slug from label
Selected: Explicit 'id' field — changing display label doesn't break URL

**Q: What happens when query string key is absent?**
Options: Show first non-hidden tab | Write default to URL immediately
Selected: Show first non-hidden tab — no URL write on initial render

**Q: History push or replace for tab switches?**
Options: Replace | Push
Selected: Replace — tab switches don't create history entries

---

## UI-SPEC Step

**Q: Run /gsd:ui-phase first or fold audit into planning?**
Options: Run /gsd:ui-phase first | Fold into planning
Selected: Run /gsd:ui-phase first — thorough documented audit before implementation

**Q: Which views should the UI-SPEC cover?**
Options: All 5 (Login, Home, Player, Deck, Game) | All 6 including Pod view
Selected: All 6 including Pod view — re-audit Pod to catch anything missed in Phase 2

**Q: What should UI-SPEC focus on?**
Options: Layout + spacing + typography | Full visual audit including interactions
Selected: Full visual audit — layout, spacing, typography, hover states, focus states, dialogs, touch targets

**Q: Should UI-SPEC specify new shared component patterns or visual only?**
Options: Visual issues only | Both visual and structural
Selected: Visual issues only — structural decisions in CONTEXT.md

**Q: Should NewGameView be included in the audit?**
Options: No — Phase 4 redesigns it | Yes — audit it too
Selected: No — Phase 4 redesigns it entirely, wasted effort

**Q: Should JoinView be in the audit?**
Options: Yes — include it | No — it's small
Selected: Yes — part of the user-facing flow, especially for new users

**Q: Issue format?**
Options: Structured issue list | Prose description
Selected: Structured issue list — numbered issues with type tag, current → desired behavior

**Q: Should researcher read Phase 2 artifacts?**
Options: Yes — mandatory | CONTEXT.md is enough
Selected: Yes — mandatory: 02-CONTEXT.md and 02-UI-SPEC.md

**Q: Mobile viewport scope?**
Options: Mobile-first 375px primary | Both 375px and 1280px | Desktop only
Selected: Mobile-first 375px primary

**Q: Include positive patterns to preserve?**
Options: Issues only | Issues + patterns to preserve
Selected: Issues only

---

## Tooltip Scope

**Q: Known use cases beyond commander disambiguation?**
Options: Commander disambiguation only | I have specific places in mind
Selected: Commander disambiguation only — UI-SPEC may surface more

**Q: What should the components look like?**
Options: InfoIcon + IconButton wrapper | Just one component
Selected: InfoIcon + IconButton wrapper — two variants: TooltipIcon (info, not tappable) and TooltipIconButton (tappable action)

**Q: Mobile behavior?**
User note: TooltipIcon tap toggles tooltip. TooltipIconButton tap = click. Mobile principle: if purpose can't be understood from icon alone, use text button.

**Q: Where should components live?**
Options: app/src/components/TooltipIcon.tsx | Separate files
Selected: app/src/components/TooltipIcon.tsx — both in one file

**Q: Should TooltipIconButton accept onClick prop (generic)?**
Options: Generic — accepts onClick | Commander-specific
Selected: Generic — accepts title + onClick + icon

**Q: Should TooltipIcon accept custom icon prop?**
Options: Yes — custom icon prop | No — always InfoOutlined
Selected: Yes — optional icon?: ReactElement, defaults to InfoOutlined
