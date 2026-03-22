- Need to handle hooking up seeded users
- Handle user oath token expiration mid-session:
    - A btw session responded with this when asked out handling a user's oauth token expiring mid session:
      ```
      Mid-session expiry: If the token expires while the user is actively using the app, useAuth() won't re-check. The
      user state remains set from the initial load. The user will only get bounced to /login when they make an API call
      that returns 401 — but that's in the individual fetch functions in http.ts, which currently just throw new
      Error(...). Nothing catches that throw and redirects to /login.
  
      So for now: graceful on load, silent failure mid-session. That's acceptable for Phase 3 — the mid-session case would
      be addressed later by adding a 401 interceptor in http.ts (e.g., call logout() from context on any 401 response).
      Not a gap in 3E itself.
      ```
- Refreshing any page in browser results in a blank white screen
- When loading the HomeView, the default `No pods yet` text is always shown until loading data lands. See TODO in @app/src/index.tsx line 43
- Pod -> Decks tab should sort the data grid by record by default
- Pod -> Players view should show player records and points within the pod
- Current pagination is great for performance, but it makes operations like seeing all the decks for a pod sorted by record difficult.
    - E.g., decks on the second page of results won't be sorted properly.
- Possible 403 error trying to retire a deck - `forbidden: deck 3 does not belong to caller`
- Frontend file restructure:
    - The following view files have gotten very large with multiple components:
        - @app/src/routes/pod.tsx
        - @app/src/routes/player.tsx
        - @app/src/routes/deck.tsx
        - @app/src/routes/game.tsx
    - Refactor into a structure like:
        - @app/src/routes/<ModelName> ->
            - view.tsx
            - players.tsx
            - settings.tsx
            - etc, file name depending on components it holds
- There should be a general component for rendering the tabs used in various views:
    - The current tab needs to be managed via a query string param - that way it persists during browser navigation
    - There should be some general error and loading handling for tab content so that each content component does not need to implement its own Skeleton and Error read out
    - The relevant views that need to be refactored are:
        - @app/src/routes/pod.tsx
        - @app/src/routes/player.tsx
        - @app/src/routes/deck.tsx
- There should be general components for rendering tooltip enabled icons and tooltip enabled icon buttons
    - @app/src/routes/pod.tsx should use tooltip icon buttons for promoting or removing pod members
    - @app/src/routes/pod.tsx should use tooltip icon buttons for saving pod name change and copying invite links
    - @app/src/routes/game.tsx has a number of buttons
- The new pod creation mechanism is currently hidden withing player settings
    - See TODO in @app/src/routes/player.tsx at line 175
    - This functionality should live in the pod page
    - Need to make sure users aren't confusingly orphaned when they first start and have no pods
    - See TODO in @app/src/index.tsx line 42
- Various simple functional TODOs:
    - @app/src/routes/pod.tsx line 63
    - @app/src/routes/pod.tsx line 165 - title case pod roles server side the same way format names are handled
    - @app/src/routes/deck.tsx line 43
    - @app/src/index.tsx line 25
    - @app/src/routes/join.tsx line 67
    - @app/src/routes/error.tsx line 9 - make hard error page match base styling of everywhere else
- Investigate and confirm retirement behavior, see TODO at line 175 in @app/src/routes/deck.tsx
- The new Commander GetAll repository method needs tests - see TODO at line 175 in @lib/repositories/commander/repo.go
- Rebuild CLAUDE context on frontend patterns
- GameResults and players
    - See TODO at line 207 in @app/src/routes/game.tsx
    - When adding/updating game results (also when creating new games), I don't think the player actually matters
    - The deck plays in the game, technically a single player could own any or all of the decks being played
    - Need to discuss options and desired behavior
- The new game form still looks terrible - @app/src/routes/new.tsx
- The record renderer in @app/sc/stats.tsx is hardcoded to only account for 4 person games.
  - This limitation is no longer in place server side
  - Frontend should account for any number of places within the record display
  - Will need some max truncate though, need to discuss options
- Need mobile styling help everywhere
- 