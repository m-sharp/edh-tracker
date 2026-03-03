package deck

import (
	"context"
	"fmt"

	"github.com/m-sharp/edh-tracker/lib/business/commander"
	"github.com/m-sharp/edh-tracker/lib/business/format"
	"github.com/m-sharp/edh-tracker/lib/business/player"
	repos "github.com/m-sharp/edh-tracker/lib/repositories"
)

func GetCommanderEntry(
	deckCmdrRepo repos.DeckCommanderRepository,
	getCommanderName commander.GetCommanderNameFunc,
) GetCommanderEntryFunc {
	return func(ctx context.Context, deckID int) (*CommanderInfo, error) {
		dcm, err := deckCmdrRepo.GetByDeckId(ctx, deckID)
		if err != nil {
			return nil, fmt.Errorf("failed to get commander for deck %d: %w", deckID, err)
		}
		if dcm == nil {
			return nil, nil
		}

		cmdName, err := getCommanderName(ctx, dcm.CommanderID)
		if err != nil {
			return nil, err
		}

		entry := &CommanderInfo{
			CommanderID:   dcm.CommanderID,
			CommanderName: cmdName,
		}

		if dcm.PartnerCommanderID != nil {
			partnerName, err := getCommanderName(ctx, *dcm.PartnerCommanderID)
			if err != nil {
				return nil, err
			}
			entry.PartnerCommanderID = dcm.PartnerCommanderID
			entry.PartnerCommanderName = &partnerName
		}

		return entry, nil
	}
}

func GetDeckName(deckRepo repos.DeckRepository) GetDeckNameFunc {
	return func(ctx context.Context, deckID int) (string, error) {
		d, err := deckRepo.GetById(ctx, deckID)
		if err != nil {
			return "", fmt.Errorf("failed to look up deck %d: %w", deckID, err)
		}
		if d == nil {
			return "", fmt.Errorf("deck %d not found", deckID)
		}
		return d.Name, nil
	}
}

func GetAll(
	deckRepo repos.DeckRepository,
	gameResultRepo repos.GameResultRepository,
	getPlayerName player.GetPlayerNameFunc,
	getFormat format.GetByIDFunc,
	getCommanderEntry GetCommanderEntryFunc,
) GetAllFunc {
	return func(ctx context.Context) ([]EntityWithStats, error) {
		decks, err := deckRepo.GetAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get decks: %w", err)
		}

		playerCache := map[int]string{}

		result := make([]EntityWithStats, 0, len(decks))
		for _, d := range decks {
			f, err := getFormat(ctx, d.FormatID)
			if err != nil {
				return nil, fmt.Errorf("failed to look up format %d: %w", d.FormatID, err)
			}
			if f == nil {
				return nil, fmt.Errorf("format %d not found", d.FormatID)
			}

			playerName, err := cachedPlayerName(ctx, d.PlayerID, playerCache, getPlayerName)
			if err != nil {
				return nil, err
			}

			commanders, err := getCommanderEntry(ctx, d.ID)
			if err != nil {
				return nil, err
			}

			entity := ToEntity(d, playerName, f.Name, commanders)

			agg, err := gameResultRepo.GetStatsForDeck(ctx, d.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get stats for deck %d: %w", d.ID, err)
			}

			result = append(result, ToEntityWithStats(entity, agg))
		}

		return result, nil
	}
}

func GetAllForPlayer(
	deckRepo repos.DeckRepository,
	gameResultRepo repos.GameResultRepository,
	getPlayerName player.GetPlayerNameFunc,
	getFormat format.GetByIDFunc,
	getCommanderEntry GetCommanderEntryFunc,
) GetAllForPlayerFunc {
	return func(ctx context.Context, playerID int) ([]EntityWithStats, error) {
		decks, err := deckRepo.GetAllForPlayer(ctx, playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get decks for player %d: %w", playerID, err)
		}

		playerName, err := getPlayerName(ctx, playerID)
		if err != nil {
			return nil, err
		}

		result := make([]EntityWithStats, 0, len(decks))
		for _, d := range decks {
			f, err := getFormat(ctx, d.FormatID)
			if err != nil {
				return nil, fmt.Errorf("failed to look up format %d: %w", d.FormatID, err)
			}
			if f == nil {
				return nil, fmt.Errorf("format %d not found", d.FormatID)
			}

			commanders, err := getCommanderEntry(ctx, d.ID)
			if err != nil {
				return nil, err
			}

			entity := ToEntity(d, playerName, f.Name, commanders)

			agg, err := gameResultRepo.GetStatsForDeck(ctx, d.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get stats for deck %d: %w", d.ID, err)
			}

			result = append(result, ToEntityWithStats(entity, agg))
		}

		return result, nil
	}
}

func GetByID(
	deckRepo repos.DeckRepository,
	gameResultRepo repos.GameResultRepository,
	getPlayerName player.GetPlayerNameFunc,
	getFormat format.GetByIDFunc,
	getCommanderEntry GetCommanderEntryFunc,
) GetByIDFunc {
	return func(ctx context.Context, deckID int) (*EntityWithStats, error) {
		d, err := deckRepo.GetById(ctx, deckID)
		if err != nil {
			return nil, fmt.Errorf("failed to get deck %d: %w", deckID, err)
		}
		if d == nil {
			return nil, nil
		}

		f, err := getFormat(ctx, d.FormatID)
		if err != nil {
			return nil, fmt.Errorf("failed to look up format %d: %w", d.FormatID, err)
		}
		if f == nil {
			return nil, fmt.Errorf("format %d not found", d.FormatID)
		}

		playerName, err := getPlayerName(ctx, d.PlayerID)
		if err != nil {
			return nil, err
		}

		commanders, err := getCommanderEntry(ctx, d.ID)
		if err != nil {
			return nil, err
		}

		entity := ToEntity(*d, playerName, f.Name, commanders)

		agg, err := gameResultRepo.GetStatsForDeck(ctx, deckID)
		if err != nil {
			return nil, fmt.Errorf("failed to get stats for deck %d: %w", deckID, err)
		}

		e := ToEntityWithStats(entity, agg)
		return &e, nil
	}
}

func Create(
	deckRepo repos.DeckRepository,
	deckCmdrRepo repos.DeckCommanderRepository,
	getFormat format.GetByIDFunc,
) CreateFunc {
	return func(ctx context.Context, playerID int, name string, formatID int, commanderID *int, partnerCommanderID *int) (int, error) {
		f, err := getFormat(ctx, formatID)
		if err != nil {
			return 0, fmt.Errorf("failed to look up format %d: %w", formatID, err)
		}
		if f == nil {
			return 0, fmt.Errorf("format %d not found", formatID)
		}

		if f.Name == "commander" && commanderID == nil {
			return 0, fmt.Errorf("commander_id is required for commander format decks")
		}

		deckID, err := deckRepo.Add(ctx, playerID, name, formatID)
		if err != nil {
			return 0, fmt.Errorf("failed to create deck: %w", err)
		}

		if commanderID != nil {
			if _, err := deckCmdrRepo.Add(ctx, deckID, *commanderID, partnerCommanderID); err != nil {
				return 0, fmt.Errorf("failed to create deck commander association: %w", err)
			}
		}

		return deckID, nil
	}
}

func Retire(deckRepo repos.DeckRepository) RetireFunc {
	return func(ctx context.Context, deckID int) error {
		return deckRepo.Retire(ctx, deckID)
	}
}

func cachedPlayerName(ctx context.Context, playerID int, cache map[int]string, getPlayerName player.GetPlayerNameFunc) (string, error) {
	if name, ok := cache[playerID]; ok {
		return name, nil
	}
	name, err := getPlayerName(ctx, playerID)
	if err != nil {
		return "", err
	}
	cache[playerID] = name
	return name, nil
}
