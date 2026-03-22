---
name: react-router
description: Use this skill to load React Router v6 patterns and conventions for the EDH Tracker frontend. TRIGGER when adding or modifying routes in app/src/index.tsx, creating files in app/src/routes/, adding route protection (RequireAuth, loader-based redirects), adding/changing loader or action functions, using hooks like useNavigate/useParams/useLocation/useFetcher/useNavigation/useActionData/useLoaderData, adding Link/NavLink/Form components, or wiring a new page into the router. DO NOT TRIGGER for backend Go routing (lib/routers/, api.go, Gorilla Mux), MUI styling changes without routing involvement, API client or auth context changes that don't affect navigation, or seed data / smoke test tasks.
version: 1.0.0
---

# React Router v6 Patterns — EDH Tracker Frontend

## Project Router Setup

The app uses React Router v6 **data router** (`createBrowserRouter` + `RouterProvider`).

- **Router config**: `app/src/index.tsx` — `createBrowserRouter([...])` defines all routes
- **Root layout**: `app/src/routes/root.tsx` — nav bar + `<Outlet />` renders child routes
- **Wrapper**: `<RouterProvider router={router} />` in `index.tsx` (no `<BrowserRouter>` component)
- All page routes are children of the root `/` route

Route objects support: `path`, `element`, `loader`, `action`, `errorElement`, `children`

## Current Route Tree

```
/ (Root layout — nav bar, Outlet)
├── /decks           loader: GetDecks
├── /deck/:deckId    loader: GetDeck
├── /games           loader: GetGames
├── /game/:gameId    loader: GetGame
├── /players         loader: GetPlayers
├── /player/:playerId  loader: GetPlayer
├── /new-game        loader: GetNewDeckInfo, action: createGame
└── [/login — file exists at routes/login.tsx but not yet registered in router]
```

## Adding a New Route

Add to the `children` array in `app/src/index.tsx`:

```tsx
{
  path: "/my-path",
  element: <MyView />,
  loader: myLoader,        // optional — runs before render
  action: myAction,        // optional — handles form/submit
  errorElement: <ErrorPage />,  // optional — inherits root's errorElement if omitted
}
```

## Loader / Action Signatures

```tsx
import { LoaderFunctionArgs, ActionFunctionArgs, redirect } from "react-router-dom";

// Loader — pre-fetch data before component renders
export async function myLoader({ params }: LoaderFunctionArgs) {
  return await fetchSomething(params.id);
}

// Action — handles form submission (POST/PUT/DELETE)
export async function myAction({ request }: ActionFunctionArgs) {
  const data = await request.json();  // or request.formData()
  const res = await postSomething(data);
  if (!res.ok) throw new Error("Failed");
  return redirect("/target");
}
```

## API Quick-Reference

| Need | Use |
|------|-----|
| Pre-fetch data before render | `loader` on route + `useLoaderData()` in component |
| Handle form/mutation | `action` on route + `<Form method="post">` or `useSubmit()` |
| Programmatic navigation | `useNavigate()` in component, or `redirect()` in loader/action |
| URL param (`:deckId`) | `useParams()` |
| Current location/pathname | `useLocation()` |
| Query string | `useSearchParams()` |
| Fetch without navigating | `useFetcher()` — e.g. like/unlike, inline edits |
| Navigation pending state | `useNavigation()` → `state: "idle" \| "loading" \| "submitting"` |
| Action response data | `useActionData()` |
| Route error display | `errorElement` on route + `useRouteError()` |
| Active nav link styling | `<NavLink>` — auto-applies `active` class |
| Protect a route | See route protection section below |

## Route Protection

Two approaches — pick based on context:

### 1. `RequireAuth` wrapper component (`app/src/routes/RequireAuth.tsx`)

Wrap the route's `element`:
```tsx
{
  path: "/protected",
  element: <RequireAuth><ProtectedView /></RequireAuth>,
  loader: protectedLoader,
}
```
Shows `<CircularProgress>` while auth loads, redirects to `/login` if no user. Uses `useAuth()` from `app/src/auth.tsx`.

### 2. Loader-based redirect (preferred for data router)

```tsx
import { redirect } from "react-router-dom";
import { GetMe } from "../http";

export async function protectedLoader({ params }: LoaderFunctionArgs) {
  const user = await GetMe();
  if (!user) throw redirect("/login");
  return fetchPageData(params.id);
}
```

Runs before the component mounts — cleaner separation from UI, no loading flicker.

## Key Files

| File | Purpose |
|------|---------|
| `app/src/index.tsx` | `createBrowserRouter` config — add routes here |
| `app/src/routes/root.tsx` | Layout: nav bar, `<Outlet />`, `useLocation` |
| `app/src/routes/RequireAuth.tsx` | Auth guard wrapper component |
| `app/src/routes/login.tsx` | Login page (needs to be registered in router) |
| `app/src/routes/error.tsx` | Error boundary (`useRouteError`) |
| `app/src/auth.tsx` | `AuthContext`, `AuthProvider`, `useAuth()` hook |

## Common Hooks — Import Path

```tsx
import {
  useLoaderData, useActionData, useNavigation,
  useNavigate, useParams, useLocation, useSearchParams,
  useFetcher, useRouteError,
  Link, NavLink, Form, Navigate,
  redirect,
} from "react-router-dom";
```

## `useFetcher` — Mutations Without Navigation

Use when you want to call an action or loader without navigating away:
```tsx
const fetcher = useFetcher();

// Submit programmatically
fetcher.submit({ id: "123" }, { method: "post", action: "/my-action" });

// Or use as a form
<fetcher.Form method="post" action="/my-action">
  <button type="submit">Like</button>
</fetcher.Form>

// Check state
if (fetcher.state === "submitting") { /* show spinner */ }
```

## `useNavigation` — Global Pending State

```tsx
const navigation = useNavigation();
// navigation.state: "idle" | "loading" | "submitting"
const isLoading = navigation.state !== "idle";
```
