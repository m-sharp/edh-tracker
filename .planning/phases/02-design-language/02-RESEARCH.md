# Phase 2: Design Language - Research

**Researched:** 2026-03-23
**Domain:** MUI v5 theming, React TypeScript frontend, dark design system
**Confidence:** HIGH

## Summary

Phase 2 introduces a MUI dark theme to an existing React/TypeScript app that currently has no custom theme. All decisions are fully locked in CONTEXT.md and the approved UI-SPEC — the planner does not need to make any visual or stack choices. The work is purely implementation: create `theme.ts`, wire it into `index.tsx`, update `index.html` for the Google Font, fix the one hardcoded color in `root.tsx`, update one Tabs prop in `pod.tsx`, and verify the Pod view at 375px.

The existing codebase is well-suited for this work. MUI v5 with `@emotion` is already installed and all required components are already in use. No new dependencies are needed. The `CssBaseline` already in `index.tsx` is compatible with dark themes. The DataGrid in `pod.tsx` adapts automatically when `palette.mode: 'dark'` is set — no custom DataGrid styling required.

The primary risk is that the `gradientBackground` CSS class is currently applied to `<body>` in `index.html` — this must be removed there AND the class definition retired from `styles.css`. Missing either removal would leave a gradient layered over the flat dark background, potentially making text unreadable.

**Primary recommendation:** Implement `theme.ts` first, wire it, then verify Pod view at 375px. Fix `gradientBackground` removal in both files atomically to avoid partial state.

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

- **D-01:** Dark theme — MUI `mode: 'dark'` throughout
- **D-02:** Background: `#0f1117` (near-black)
- **D-03:** Surface/Paper: `#1a1a2e` (dark navy — used for cards, elevated surfaces)
- **D-04:** Primary text: `#e8edf5` (off-white); secondary text at MUI's default dark-mode opacity
- **D-05:** Primary/accent: `#c9a227` (legendary gold), hover `#e0b83a`
- **D-06:** Error/danger: MUI default `error.main` — no custom color needed
- **D-07:** AppBar: solid dark color that coheres with the surface palette (no gradient)
- **D-08:** Display font (h1–h4): Josefin Sans from Google Fonts — geometric sans-serif, all-caps friendly
- **D-09:** Body text: Roboto (MUI default, unchanged)
- **D-10:** Numeric stats/code: monospace (MUI default)
- **D-11:** Load Josefin Sans via `<link>` in `app/public/index.html` (Google Fonts CDN), reference it in `createTheme` typography
- **D-12:** Create a dedicated theme file `app/src/theme.ts` exporting a `createTheme(...)` result
- **D-13:** Wrap the app in `<ThemeProvider theme={theme}>` in `index.tsx` (alongside existing `CssBaseline`)
- **D-14:** `CssBaseline` stays — it handles body background reset correctly with dark themes
- **D-15:** Use MUI DataGrid's built-in dark mode support — no custom row styling
- **D-16:** Pod view (`app/src/routes/pod.tsx`) is the validation target for Phase 2
- **D-17:** Theme applied globally via ThemeProvider — other views get the dark palette automatically
- **D-18:** Phase 2 delivers: theme file + Google Font import + ThemeProvider wiring + Pod view verified on mobile
- **D-19:** Individual view-level restyling beyond Pod view is Phase 3 work

### Claude's Discretion

- Exact AppBar background color (must be dark, cohesive with `#1a1a2e` surface — e.g., same color or slightly darker)
- Exact MUI spacing scale (`theme.spacing` default of 8px is likely fine)
- Component-level overrides in `theme.components` (button border-radius, chip variants) — choose what fits the dark/gold aesthetic

### Deferred Ideas (OUT OF SCOPE)

- Applying refined styles to all views individually (Login, Player, Deck, Game) — Phase 3
- Dark/light mode toggle for users — out of scope for launch
- Custom component variants beyond what Phase 2 validates — Phase 3+

</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| DSNG-01 | App has a defined visual design language (color palette, typography, spacing tokens) implemented consistently across all views | theme.ts with createTheme covers palette + typography + spacing; ThemeProvider applies globally |
| DSNG-02 | All views are usable on mobile (phone-sized viewport) — layout adapts, touch targets are adequate, text is readable without zooming | Pod view verified at 375px: Tabs scrollable, DataGrid height, 44px touch targets, 16px min body text |
| DSNG-03 | MUI components used properly and consistently — no ad-hoc styling where MUI has a clear pattern | Replace hardcoded `bgcolor: "#f0f5fa"` with `bgcolor: "background.default"` from theme; tab overrides via theme.components |

</phase_requirements>

---

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| @mui/material | v5.15.2 (installed) | createTheme, ThemeProvider, CssBaseline, all UI components | Already installed; no new dep |
| @mui/x-data-grid | v6.18.6 (installed) | DataGrid with dark mode auto-support | Already installed; adapts to palette.mode automatically |
| @emotion/react | v11.11.3 (installed) | CSS-in-JS engine for MUI v5 | Required peer dep; already present |
| @emotion/styled | v11.11.0 (installed) | Styled component support for MUI v5 | Required peer dep; already present |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| Google Fonts CDN | N/A | Josefin Sans font delivery | Added via `<link>` in index.html only |

**No new npm dependencies required for this phase.** All required packages are already installed.

**Installation:** None needed.

---

## Architecture Patterns

### File Structure Changes

```
app/
├── public/
│   └── index.html          # Add Josefin Sans <link>; remove gradientBackground from <body>
└── src/
    ├── theme.ts             # NEW — exports createTheme result
    ├── index.tsx            # ADD ThemeProvider + import theme
    ├── styles.css           # RETIRE gradientBackground class
    └── routes/
        ├── root.tsx         # REPLACE bgcolor: "#f0f5fa" with bgcolor: "background.default"
        └── pod.tsx          # ADD variant="scrollable" scrollButtons="auto" to Tabs
                             # RENAME "Save" button label to "Save Pod Name"
```

### Pattern 1: createTheme in a dedicated file

**What:** All theme configuration lives in `app/src/theme.ts`. The file exports a single `createTheme(...)` result. `index.tsx` imports and uses it.

**When to use:** Always — isolates theme from component tree, allows future updates without touching index.tsx.

```typescript
// app/src/theme.ts
import { createTheme } from "@mui/material/styles";

const theme = createTheme({
  palette: {
    mode: "dark",
    background: {
      default: "#0f1117",
      paper: "#1a1a2e",
    },
    primary: {
      main: "#c9a227",
      light: "#e0b83a",
    },
    text: {
      primary: "#e8edf5",
      // secondary: MUI dark-mode default opacity applied to text.primary
    },
    // error: MUI default (no override needed per D-06)
  },
  typography: {
    fontFamily: '"Roboto", sans-serif', // body default
    h1: { fontFamily: '"Josefin Sans", "Roboto", sans-serif' },
    h2: { fontFamily: '"Josefin Sans", "Roboto", sans-serif' },
    h3: { fontFamily: '"Josefin Sans", "Roboto", sans-serif' },
    h4: { fontFamily: '"Josefin Sans", "Roboto", sans-serif', fontSize: "28px", fontWeight: 700, lineHeight: 1.2 },
    h6: { fontFamily: '"Josefin Sans", "Roboto", sans-serif', fontSize: "20px", fontWeight: 700, lineHeight: 1.2 },
    body1: { fontSize: "16px", fontWeight: 400, lineHeight: 1.5 },
    body2: { fontSize: "14px", fontWeight: 400, lineHeight: 1.4 },
  },
  components: {
    MuiButton: {
      styleOverrides: { root: { borderRadius: "6px" } },
      defaultProps: { disableElevation: true },
    },
    MuiAppBar: {
      styleOverrides: { root: { backgroundColor: "#1a1a2e" } },
    },
    MuiChip: {
      styleOverrides: { root: { fontFamily: '"Roboto", sans-serif' } },
    },
    MuiTab: {
      styleOverrides: { root: { textTransform: "none" } },
    },
    // MuiDataGrid: no override needed — adapts automatically
  },
});

export default theme;
```

### Pattern 2: ThemeProvider wiring in index.tsx

**What:** Wrap the existing `RouterProvider` with `ThemeProvider`. `CssBaseline` must remain inside (or alongside within) `ThemeProvider` to pick up dark mode body background reset.

```typescript
// app/src/index.tsx — relevant changes only
import { ThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import theme from "./theme";

// Render tree becomes:
<StrictMode>
  <ThemeProvider theme={theme}>
    <CssBaseline enableColorScheme />
    <AuthProvider>
      <RouterProvider router={router} />
    </AuthProvider>
  </ThemeProvider>
</StrictMode>
```

**Important:** `CssBaseline` must be INSIDE `ThemeProvider` to inherit the dark background. Currently `CssBaseline` is outside any ThemeProvider — this is the critical wiring change.

### Pattern 3: Google Fonts CDN in index.html

**What:** Add Josefin Sans preconnect + stylesheet link. Remove the `gradientBackground` class from `<body>`.

```html
<!-- app/public/index.html — add after existing Roboto link -->
<link rel="preconnect" href="https://fonts.googleapis.com" />
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Josefin+Sans:wght@400;700&display=swap" />

<!-- REMOVE class="gradientBackground" from <body> -->
<body>
  <div id="root"></div>
</body>
```

Note: The Roboto preconnect links are already present. Only the Josefin Sans link needs to be added. The `gradientBackground` class attribute must be removed from `<body>` here, and the class definition retired from `styles.css`.

### Pattern 4: Mobile Tabs fix in pod.tsx

**What:** Add `variant="scrollable" scrollButtons="auto"` to the `<Tabs>` component so all 4 tabs fit at 375px width.

```typescript
// pod.tsx — existing Tabs component
<Tabs
  value={tab}
  onChange={(_, v) => setTab(v)}
  sx={{ mb: 2 }}
  variant="scrollable"
  scrollButtons="auto"
>
```

### Anti-Patterns to Avoid

- **CssBaseline outside ThemeProvider:** `CssBaseline` must be inside `ThemeProvider` — otherwise body background will not reset to `#0f1117`; the gradient or browser default will show through.
- **Hardcoded hex colors in sx props:** After this phase, `bgcolor: "#f0f5fa"` in root.tsx is the only known hardcoded background color. Replace with `bgcolor: "background.default"` from theme. Do not introduce new hardcoded colors.
- **Leaving gradientBackground on `<body>`:** The class is applied in `index.html` AND defined in `styles.css`. Both must be addressed. Removing only the CSS class definition while leaving the attribute on `<body>` has no effect (class just has no rules). Removing only the attribute while leaving the definition in CSS is safe but incomplete — retire the definition too.
- **Overriding DataGrid colors manually:** DataGrid adapts automatically. Do not add custom `sx` overrides for DataGrid row/header background colors.
- **Using `npx tsc` for TypeScript verification:** Use `./node_modules/.bin/tsc --noEmit` from `app/` per the frontend-verify skill.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Dark mode body background | Custom CSS class on `<body>` | `CssBaseline` inside `ThemeProvider` | CssBaseline injects the `background.default` color onto the body element automatically |
| Font loading | Manual JS font loading | Google Fonts `<link>` in index.html | CDN handles font file delivery and caching; `font-display: swap` via `display=swap` query param |
| DataGrid dark theme | Custom row/cell `sx` props | `palette.mode: 'dark'` in createTheme | DataGrid reads palette mode and applies its own dark styles automatically |
| Component hover states | Custom `&:hover` in sx | `palette.primary.light` for hover token | MUI uses `palette.primary.light` and `palette.primary.dark` for hover derivation on contained buttons |

---

## Existing Code State (Audit)

This section documents the pre-existing state that each task must transform. HIGH confidence — all verified by direct file reading.

### Files to create
| File | Action |
|------|--------|
| `app/src/theme.ts` | Create new — does not exist |

### Files to modify
| File | Current State | Required Change |
|------|--------------|-----------------|
| `app/src/index.tsx` | `CssBaseline` at root with no `ThemeProvider` | Add `ThemeProvider` wrapper; move `CssBaseline` inside it |
| `app/public/index.html` | `<body class="gradientBackground">` with Roboto Google Font link | Add Josefin Sans link; remove `class="gradientBackground"` from `<body>` |
| `app/src/styles.css` | `.gradientBackground` class defined with radial gradients | Remove/comment out `gradientBackground` class definition |
| `app/src/routes/root.tsx` | `Container` has `bgcolor: "#f0f5fa"` in `sx` prop | Replace with `bgcolor: "background.default"` (or remove — CssBaseline handles body) |
| `app/src/routes/pod.tsx` | `<Tabs>` has no `variant` or `scrollButtons` props; "Save" button label | Add `variant="scrollable" scrollButtons="auto"`; rename "Save" to "Save Pod Name" |

### No-change files
| File | Reason |
|------|--------|
| `app/src/routes/pod.tsx` DataGrid | No changes needed — auto-adapts to dark palette |
| All other route files | Phase 2 scope is global theme only; individual view styling is Phase 3 |

---

## Common Pitfalls

### Pitfall 1: CssBaseline outside ThemeProvider

**What goes wrong:** Body background stays at browser default (white or system default) instead of `#0f1117`. App renders dark components on a white/grey background.

**Why it happens:** `CssBaseline` reads the MUI theme context to inject `background.default` onto `body`. If placed outside `ThemeProvider`, there is no theme context and the default (light) styles apply.

**How to avoid:** Ensure the render tree is `ThemeProvider > CssBaseline > ... rest of app`. Verify by checking that `<body style="background-color: #0f1117">` appears in DevTools after wiring.

**Warning signs:** White flash visible on initial load; body element has no `background-color` style in DevTools.

### Pitfall 2: gradientBackground class still active

**What goes wrong:** The radial gradient defined in `styles.css` remains on `<body>` even after the flat dark background is applied, potentially causing visual artifacts or discoloration.

**Why it happens:** The class is applied as `class="gradientBackground"` directly in `index.html`. Deleting the CSS rules without removing the attribute from the HTML tag is harmless but the class definition is still present (or vice versa).

**How to avoid:** Address both locations atomically: remove the `class="gradientBackground"` attribute from `<body>` in `index.html` AND retire the CSS class definition from `styles.css`.

**Warning signs:** Dark background appears but has subtle radial gradient banding visible against a flat dark background.

### Pitfall 3: Typography font not loading

**What goes wrong:** Josefin Sans is referenced in `theme.ts` but not loaded — headings fall back to Roboto. The visual difference may be subtle on first glance.

**Why it happens:** Google Fonts requires a `<link>` in the HTML `<head>`. The CDN link must specify the exact weights used (400 and 700 per UI-SPEC).

**How to avoid:** Add the link with `wght@400;700` in the URL, not just a single weight. Verify by inspecting an h4 element in DevTools — computed font-family should show `Josefin Sans`.

**Warning signs:** Pod name heading looks indistinguishable from body text font.

### Pitfall 4: AppBar still renders in primary.main gold

**What goes wrong:** Without the `MuiAppBar` component override, MUI's default behavior sets the AppBar background to `palette.primary.main` — which becomes `#c9a227` (gold). This makes the top bar look garish and is not the design intent.

**Why it happens:** MUI AppBar defaults to `color="primary"` which maps to `palette.primary.main` as background. The theme override `MuiAppBar.styleOverrides.root.backgroundColor: "#1a1a2e"` suppresses this default.

**How to avoid:** Ensure the `MuiAppBar` component override is present in `theme.components`. Also note that `root.tsx` currently sets no explicit `color` prop on `<AppBar>` — it relies on the component override being present.

**Warning signs:** AppBar renders gold instead of dark navy.

### Pitfall 5: TypeScript import path for ThemeProvider

**What goes wrong:** `ThemeProvider` imported from `@mui/material` (not `@mui/material/styles`) can cause duplicate theme context issues or React reconciliation warnings in some CRA builds.

**Why it happens:** MUI re-exports from both paths but the `@mui/material/styles` path is the canonical one for theme infrastructure.

**How to avoid:** Import `ThemeProvider` and `createTheme` from `"@mui/material/styles"`, not from `"@mui/material"`. Both work but the styles sub-path is conventional and avoids potential confusion.

---

## Code Examples

### Minimal working theme.ts structure (Source: UI-SPEC + CONTEXT.md decisions)

```typescript
import { createTheme } from "@mui/material/styles";

const theme = createTheme({
  palette: {
    mode: "dark",
    background: { default: "#0f1117", paper: "#1a1a2e" },
    primary: { main: "#c9a227", light: "#e0b83a" },
    text: { primary: "#e8edf5" },
  },
  typography: {
    h4: { fontFamily: '"Josefin Sans", "Roboto", sans-serif', fontSize: "28px", fontWeight: 700 },
    h6: { fontFamily: '"Josefin Sans", "Roboto", sans-serif', fontSize: "20px", fontWeight: 700 },
    // h1–h3 also get Josefin Sans; body1/body2 stay Roboto
  },
  components: {
    MuiButton: { styleOverrides: { root: { borderRadius: "6px" } }, defaultProps: { disableElevation: true } },
    MuiAppBar: { styleOverrides: { root: { backgroundColor: "#1a1a2e" } } },
    MuiChip: { styleOverrides: { root: { fontFamily: '"Roboto", sans-serif' } } },
    MuiTab: { styleOverrides: { root: { textTransform: "none" } } },
  },
});

export default theme;
```

### index.tsx render tree after wiring

```typescript
import { ThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import theme from "./theme";

createRoot(document.getElementById("root") as HTMLElement).render(
  <StrictMode>
    <ThemeProvider theme={theme}>
      <CssBaseline enableColorScheme />
      <AuthProvider>
        <RouterProvider router={router} />
      </AuthProvider>
    </ThemeProvider>
  </StrictMode>
);
```

### root.tsx Container fix

```typescript
// Before:
<Container id="detail" component="main" sx={{ p: 3, width: "90%", bgcolor: "#f0f5fa", mt: 12, mb: 5 }} maxWidth="xl">

// After:
<Container id="detail" component="main" sx={{ p: 3, width: "90%", bgcolor: "background.default", mt: 12, mb: 5 }} maxWidth="xl">
// OR: remove bgcolor entirely — CssBaseline sets body background; Container inherits
```

### pod.tsx Tabs fix for mobile

```typescript
// Before:
<Tabs value={tab} onChange={(_, v) => setTab(v)} sx={{ mb: 2 }}>

// After:
<Tabs value={tab} onChange={(_, v) => setTab(v)} sx={{ mb: 2 }} variant="scrollable" scrollButtons="auto">
```

### pod.tsx Settings tab button label fix

```typescript
// Before:
<Button variant="contained" onClick={handleSaveName}>Save</Button>

// After:
<Button variant="contained" onClick={handleSaveName}>Save Pod Name</Button>
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| MUI v4 `makeStyles` / JSS | MUI v5 `sx` prop + `createTheme` | MUI v5 (2021) | Project already uses v5 — no migration needed |
| ThemeProvider from `@material-ui/core` | ThemeProvider from `@mui/material/styles` | MUI v5 | Import path changed; use `/styles` sub-path |
| DataGrid manual dark styling | DataGrid auto dark mode via `palette.mode` | x-data-grid v5+ | No custom DataGrid overrides needed |

**Deprecated/outdated:**
- `makeStyles` hook: replaced by `sx` prop and `styled()` in MUI v5 — do not introduce it
- `withStyles` HOC: same — do not use

---

## Open Questions

None. All design decisions are locked in CONTEXT.md and fully specified in the approved UI-SPEC. No gaps require resolution before planning.

---

## Environment Availability

Step 2.6: This phase is purely frontend code/config changes. All dependencies (`@mui/material`, `@emotion`, `@mui/x-data-grid`) are already installed in `app/node_modules/`. Google Fonts is a CDN link with no package installation required.

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| @mui/material | ThemeProvider, createTheme | Yes | v5.15.2 | — |
| @emotion/react | MUI v5 peer dep | Yes | v11.11.3 | — |
| @emotion/styled | MUI v5 peer dep | Yes | v11.11.0 | — |
| @mui/x-data-grid | DataGrid dark mode | Yes | v6.18.6 | — |
| Google Fonts CDN | Josefin Sans | CDN (no install) | — | System font fallback (Roboto already loaded) |
| tsc (TypeScript check) | frontend-verify skill | Yes | v4.9.5 via node_modules | — |

**No missing dependencies.**

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | TypeScript compiler (tsc) — no Jest/RTL tests for frontend per REQUIREMENTS.md out-of-scope |
| Config file | `app/tsconfig.json` |
| Quick run command | `cd /mnt/d/msharp/Documents/projects/edh-tracker/app && ./node_modules/.bin/tsc --noEmit 2>&1` |
| Full suite command | Same (no separate full suite for frontend) |

REQUIREMENTS.md explicitly marks "Frontend automated tests" as out of scope: "High effort for existing code; post-launch investment." This means all validation for Phase 2 is manual or TypeScript-compile-only.

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | Exists? |
|--------|----------|-----------|-------------------|---------|
| DSNG-01 | theme.ts exports createTheme with correct palette/typography/spacing | TypeScript compile | `./node_modules/.bin/tsc --noEmit` | Wave 0 — file doesn't exist yet |
| DSNG-01 | ThemeProvider wraps app correctly | TypeScript compile | `./node_modules/.bin/tsc --noEmit` | After index.tsx edit |
| DSNG-02 | Pod view usable at 375px viewport | Manual visual check | Browser DevTools — set to 375px iPhone SE | Manual only |
| DSNG-02 | Touch targets ≥ 44px | Manual visual check | Browser DevTools inspect element heights | Manual only |
| DSNG-03 | No hardcoded bgcolor after root.tsx fix | Code review + compile | `./node_modules/.bin/tsc --noEmit` | After root.tsx edit |

### Sampling Rate

- **Per task commit:** `cd /mnt/d/msharp/Documents/projects/edh-tracker/app && ./node_modules/.bin/tsc --noEmit 2>&1`
- **Per wave merge:** Same TypeScript check + manual Pod view at 375px
- **Phase gate:** TypeScript clean + manual mobile verification before `/gsd:verify-work`

### Wave 0 Gaps

- [ ] `app/src/theme.ts` — must be created in Wave 1 before any other task can import from it; this is the dependency anchor for the entire phase

No missing test files — no automated frontend tests exist or are required per project scope.

---

## Project Constraints (from CLAUDE.md)

These directives from CLAUDE.md must be honored by all planning and implementation:

- **Tech stack locked:** React + MUI v5 + React Router v6 — no framework changes
- **TypeScript verify command:** Use `./node_modules/.bin/tsc --noEmit` from `app/` — NEVER `npm run build` or `npx tsc`
- **Frontend component conventions:** All HTTP calls via `app/src/http.ts`; all TypeScript interfaces in `app/src/types.ts`; MUI components throughout; `credentials: "include"` on all fetch calls
- **JSON field names:** snake_case throughout — no camelCase in API or TypeScript interfaces
- **No breaking changes** to existing game/player/deck data in the database (frontend-only phase — not applicable here)
- **GSD workflow enforcement:** All file changes go through GSD execute-phase — no direct repo edits outside workflow
- **Frontend component default exports:** Components are default exports from `app/src/routes/<name>.tsx`; component functions return `ReactElement`

---

## Sources

### Primary (HIGH confidence)

- Direct file reads: `app/src/index.tsx`, `app/src/routes/root.tsx`, `app/src/routes/pod.tsx`, `app/public/index.html`, `app/src/styles.css` — current state of all files to be modified
- `app/package.json` — confirmed all MUI packages already installed at specified versions
- `.planning/phases/02-design-language/02-CONTEXT.md` — locked decisions D-01 through D-19
- `.planning/phases/02-design-language/02-UI-SPEC.md` — approved design contract (checker sign-off pending but approved per frontmatter)
- `CLAUDE.md` — project conventions and constraints
- `.claude/skills/frontend-verify/SKILL.md` — verified TypeScript check command

### Secondary (MEDIUM confidence)

- MUI v5 ThemeProvider/CssBaseline interaction pattern — knowledge confirmed by project's existing use of `CssBaseline enableColorScheme` in index.tsx (if it were broken, the existing dev sessions would have caught it)
- `@mui/x-data-grid` v6 automatic dark mode — confirmed by D-15 decision (explicitly decided after verification during discuss phase)

### Tertiary (LOW confidence)

- None — all critical claims verified by direct file reading or locked decisions

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all packages confirmed installed; versions from package.json
- Architecture: HIGH — all file contents read directly; no guessing about current state
- Pitfalls: HIGH — derived from direct inspection of existing code (CssBaseline placement, gradientBackground dual-location issue, AppBar default color behavior)
- Mobile contract: HIGH — fully specified in UI-SPEC; no ambiguity

**Research date:** 2026-03-23
**Valid until:** 2026-04-22 (stable MUI v5 API; 30-day window appropriate)
