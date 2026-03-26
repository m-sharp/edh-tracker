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
- Pod -> Decks tab should sort the data grid by record by default
- Pod -> Players view should show player records and points within the pod
- Current pagination is great for performance, but it makes operations like seeing all the decks for a pod sorted by record difficult.
    - E.g., decks on the second page of results won't be sorted properly.
- The new pod creation mechanism is currently hidden withing player settings
    - See TODO in @app/src/routes/player.tsx at line 175
    - This functionality should live in the pod page
    - Need to make sure users aren't confusingly orphaned when they first start and have no pods
    - See TODO in @app/src/index.tsx line 42
- Various simple functional TODOs:
    - @app/src/routes/pod.tsx line 63
    - @app/src/routes/deck.tsx line 43
    - @app/src/routes/error.tsx line 9 - make hard error page match base styling of everywhere else
- The new Commander GetAll repository method needs tests - see TODO at line 175 in @lib/repositories/commander/repo.go
- Rebuild CLAUDE context on frontend patterns
