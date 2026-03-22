---
name: frontend-verify
description: Verify that the edh-tracker React/TypeScript frontend compiles cleanly after making changes. TRIGGER when finishing any edit to files under app/src/ — new components, type changes, http.ts additions, auth changes, route restructuring — or whenever the user asks to "verify frontend changes", "check TypeScript", or "does the frontend compile". DO NOT TRIGGER for backend Go changes (use smoke-test instead), for running the dev server (npm start), or for building the production Docker image.
version: 1.0.0
---

# EDH Tracker Frontend TypeScript Check

Run this after any change to files under `app/src/` to verify TypeScript compiles cleanly.

## Command

Run from the project root:

```bash
cd /mnt/d/msharp/Documents/projects/edh-tracker/app && ./node_modules/.bin/tsc --noEmit 2>&1
```

- Uses the local `tsc` binary directly — avoids the slow startup overhead of `npx tsc`
- `--noEmit` means no output files are written, just type-checking
- Clean exit (no output, exit 0) = success
- Any TypeScript errors are printed to stdout

## Notes

- Do NOT use `npm run build` to verify — it works but takes 60–90+ seconds
- Do NOT use `npx tsc --noEmit` — npx overhead causes timeouts
- This check covers all files in `app/src/` including new files added this session
