# Phase 2: Design Language - Context

**Gathered:** 2026-03-23
**Status:** Ready for planning

<domain>
## Phase Boundary

Define a MUI theme (color palette, typography scale, spacing tokens) that all subsequent UI work builds against. Apply the design system to the Pod view as the representative view to validate it. Other views inherit the global MUI theme automatically but are not individually restyled in this phase.

</domain>

<decisions>
## Implementation Decisions

### Visual Tone
- **D-01:** Dark theme — MUI `mode: 'dark'` throughout
- **D-02:** Background: `#0f1117` (near-black)
- **D-03:** Surface/Paper: `#1a1a2e` (dark navy — used for cards, elevated surfaces)
- **D-04:** Primary text: `#e8edf5` (off-white); secondary text at MUI's default dark-mode opacity

### Color Palette
- **D-05:** Primary/accent: `#c9a227` (legendary gold), hover `#e0b83a`
- **D-06:** Error/danger: MUI default `error.main` — no custom color needed
- **D-07:** AppBar: Claude's discretion — solid dark color that coheres with the surface palette (no gradient)

### Typography
- **D-08:** Display font (h1–h4): **Josefin Sans** from Google Fonts — geometric sans-serif, all-caps friendly, clean/modern
- **D-09:** Body text: Roboto (MUI default, unchanged)
- **D-10:** Numeric stats/code: monospace (MUI default)
- **D-11:** Load Josefin Sans via `<link>` in `app/public/index.html` (Google Fonts CDN), reference it in `createTheme` typography

### MUI Theme Setup
- **D-12:** Create a dedicated theme file (e.g., `app/src/theme.ts`) exporting a `createTheme(...)` result
- **D-13:** Wrap the app in `<ThemeProvider theme={theme}>` in `index.tsx` (alongside existing `CssBaseline`)
- **D-14:** `CssBaseline` stays — it handles body background reset correctly with dark themes

### DataGrid
- **D-15:** Use MUI DataGrid's built-in dark mode support — no custom row styling. DataGrid adapts automatically when the theme `mode` is `'dark'`.

### Representative View
- **D-16:** Pod view (`app/src/routes/pod.tsx`) is the validation target for Phase 2
  - Validates: AppBar dark treatment, Tabs on dark bg, DataGrid on dark theme, gold accent buttons, mobile layout at ~375px
- **D-17:** The theme is applied globally via ThemeProvider — other views get the dark palette automatically even if not individually refined

### Scope
- **D-18:** Phase 2 delivers: theme file + Google Font import + ThemeProvider wiring + Pod view verified on mobile
- **D-19:** Individual view-level restyling beyond Pod view is Phase 3 (Frontend Structure) work

### Claude's Discretion
- Exact AppBar background color (must be dark, cohesive with `#1a1a2e` surface — e.g., same color or slightly darker)
- Exact MUI spacing scale (`theme.spacing` default of 8px is likely fine)
- Component-level overrides in `theme.components` (e.g., button border-radius, chip variants) — choose what fits the dark/gold aesthetic

</decisions>

<specifics>
## Specific Ideas

- Gold accent should evoke the "legendary card border" from MTG — #c9a227 is the anchor; avoid making it look yellow or washed out on dark backgrounds
- Josefin Sans works well in all-caps for headings (pod names, player names, section titles)
- The existing `gradientBackground` CSS class (dark blue radial gradients) may be reused or removed in favor of the flat dark background — Claude's discretion

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Phase Requirements
- `.planning/REQUIREMENTS.md` §Design — DSNG-01, DSNG-02, DSNG-03 (color/typography/spacing tokens, mobile usability, MUI consistency)
- `.planning/ROADMAP.md` §Phase 2 — success criteria and "UI hint: yes" note

### Existing Frontend
- `app/src/index.tsx` — ThemeProvider insertion point (alongside existing CssBaseline + AuthProvider)
- `app/src/routes/root.tsx` — AppBar implementation; hardcoded `bgcolor: "#f0f5fa"` to replace
- `app/src/routes/pod.tsx` — Representative view; contains Tabs, DataGrid, Buttons, Dialogs
- `app/src/styles.css` — Contains `gradientBackground` class with existing dark-blue radial gradients

### MUI Theme
- No external spec — theme decisions are fully captured in this CONTEXT.md

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `AppBar` in `root.tsx`: Already uses MUI AppBar — just needs theme color to apply; no structural change needed
- `Tabs` in `pod.tsx` and `player.tsx`: Standard MUI Tabs — will inherit dark theme automatically
- `DataGrid` in pod/player views: `@mui/x-data-grid` v6 supports dark mode via theme `mode`

### Established Patterns
- No custom MUI theme exists — raw defaults throughout; this phase introduces `createTheme` for the first time
- One hardcoded color: `bgcolor: "#f0f5fa"` in `root.tsx` Container — replace with `background.default` from theme
- `gradientBackground` CSS class: dark-blue radial gradients — may be retired in favor of flat `#0f1117` background

### Integration Points
- `app/src/index.tsx`: Add `<ThemeProvider>` wrapper and import theme; also add Google Fonts `<link>` in `app/public/index.html`
- The existing `CssBaseline enableColorScheme` is compatible with MUI dark theme — keep it

</code_context>

<deferred>
## Deferred Ideas

- Applying refined styles to all views individually (Login, Player, Deck, Game) — Phase 3 (Frontend Structure)
- Dark/light mode toggle for users — out of scope for launch
- Custom component variants beyond what Phase 2 validates — Phase 3+

</deferred>

---

*Phase: 02-design-language*
*Context gathered: 2026-03-23*
