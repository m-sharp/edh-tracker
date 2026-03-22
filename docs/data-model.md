# Data Model

## Entity Relationships

```
Player ──→ many Decks
       ──→ one User (optional)
       ──→ many PlayerPods ──→ Pod

Pod ──→ many Games
    ──→ many PlayerPodRoles ──→ Player (with manager/member role)
    ──→ many PodInvites

Format ──→ many Decks
       ──→ many Games

Deck ──→ one Player
     ──→ one Format
     ──→ one DeckCommander (optional, via deck_commander join table)
     ──→ many GameResults

DeckCommander ──→ one Commander
              ──→ one Commander (partner_commander_id, nullable)

Game ──→ one Pod
     ──→ one Format
     ──→ many GameResults

GameResult ──→ one Game
           ──→ one Deck
```

## Tables

| Table | Key Columns |
|---|---|
| `player` | id, name, deleted_at |
| `deck` | id, player_id, name, format_id, retired, deleted_at |
| `game` | id, description, pod_id, format_id, deleted_at |
| `game_result` | id, game_id, deck_id, place, kill_count, deleted_at |
| `format` | id, name — seeded with `commander` and `other` |
| `commander` | id, name (unique card name) |
| `deck_commander` | id, deck_id, commander_id, partner_commander_id (nullable) |
| `pod` | id, name, deleted_at |
| `player_pod` | id, pod_id, player_id — legacy join table; superseded by `player_pod_role` |
| `player_pod_role` | id, pod_id, player_id, role (ENUM: manager/member), deleted_at — UNIQUE(pod_id, player_id) |
| `pod_invite` | id, pod_id, invite_code (unique), created_by_player_id, expires_at (nullable), used_count, deleted_at |
| `user` | id, player_id, role_id, oauth_provider, oauth_subject, email, display_name, avatar_url |
| `user_role` | id, name — seeded with `admin` and `player` |

All tables use soft deletes (`deleted_at` nullable DATETIME) and track `created_at`/`updated_at`, except `format`, `commander`, `user`, and `user_role` which do not have `deleted_at`.

**Cascading deletes** (Migration 19): all foreign keys on `deck_commander`, `player_pod`, `player_pod_role`, `pod_invite`, and `user` use `ON DELETE CASCADE`, so deleting a parent record automatically removes dependent rows.

## Model/Entity Split

Each domain has two struct types:

- **`<domain>.Model`** (`lib/repositories/<domain>/model.go`) — raw DB row struct with `db:""` tags; embeds `base.ModelBase` (ID, CreatedAt, UpdatedAt)
- **`<domain>.Entity`** (`lib/business/<domain>/entity.go`) — enriched business object returned to callers; includes computed/joined fields like `Stats`, `PodIDs`, commander names

Conversion is done by `ToEntity(model, ...)` functions in the business layer.

## Key Entity Fields

- **`player.Entity`** — `ID`, `Name`, `Stats stats.Stats`, `PodIDs []int`
- **`deck.Entity`** — `ID`, `PlayerName`, `Name`, `FormatName`, `Retired`, `Commanders *DeckCommanderEntry`
- **`game.Entity`** — `ID`, `Description`, `PodID`, `FormatID`, `Results []gameResult.Entity`
- **`gameResult.Entity`** — `DeckName`, `CommanderName *string`, `PartnerCommanderName *string`, `Points int`
- **`stats.Stats`** — `Record map[int]int` (place → count), `Games`, `Kills`, `Points` — computed from aggregates, never stored

## Points Formula

`points = kills + place_bonus` where place bonuses are **3/2/1/0** for 1st–4th place.

## Format Validation

When creating a game result, the deck's `format_id` must match the game's `format_id`, **except** when the format name is `"other"` (which skips the check). Enforced in the business layer (`lib/business/game/`).

## Index Convention

Any column used in a WHERE clause that is not a primary key or foreign key must have an explicit index. Add indexes in a new migration whenever adding queries that filter on non-PK/FK columns.
