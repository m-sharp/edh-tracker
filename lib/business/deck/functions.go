package deck

import (
	"context"
	"fmt"

	"github.com/m-sharp/edh-tracker/lib/business/commander"
	"github.com/m-sharp/edh-tracker/lib/business/format"
	repos "github.com/m-sharp/edh-tracker/lib/repositories"
	deckRepository "github.com/m-sharp/edh-tracker/lib/repositories/deck"
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

func GetPlayerIDForDeck(deckRepo repos.DeckRepository) GetPlayerIDForDeckFunc {
	return func(ctx context.Context, deckID int) (int, error) {
		d, err := deckRepo.GetById(ctx, deckID)
		if err != nil {
			return 0, fmt.Errorf("failed to look up deck %d: %w", deckID, err)
		}
		if d == nil {
			return 0, fmt.Errorf("deck %d not found", deckID)
		}
		return d.PlayerID, nil
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

func commanderInfoFromModel(d deckRepository.Model) *CommanderInfo {
	if d.Commander == nil {
		return nil
	}
	info := &CommanderInfo{
		CommanderID:   d.Commander.CommanderID,
		CommanderName: d.Commander.Commander.Name,
	}
	if d.Commander.PartnerCommanderID != nil {
		info.PartnerCommanderID = d.Commander.PartnerCommanderID
		name := d.Commander.PartnerCommander.Name
		info.PartnerCommanderName = &name
	}
	return info
}

func GetAll(
	deckRepo repos.DeckRepository,
	gameResultRepo repos.GameResultRepository,
) GetAllFunc {
	return func(ctx context.Context) ([]EntityWithStats, error) {
		decks, err := deckRepo.GetAllHydrated(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get decks: %w", err)
		}

		result := make([]EntityWithStats, 0, len(decks))
		for _, d := range decks {
			entity := ToEntity(d, d.Player.Name, d.Format.Name, commanderInfoFromModel(d))

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
) GetAllForPlayerFunc {
	return func(ctx context.Context, playerID int) ([]EntityWithStats, error) {
		decks, err := deckRepo.GetAllForPlayerHydrated(ctx, playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get decks for player %d: %w", playerID, err)
		}

		result := make([]EntityWithStats, 0, len(decks))
		for _, d := range decks {
			entity := ToEntity(d, d.Player.Name, d.Format.Name, commanderInfoFromModel(d))

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
) GetByIDFunc {
	return func(ctx context.Context, deckID int) (*EntityWithStats, error) {
		d, err := deckRepo.GetByIDHydrated(ctx, deckID)
		if err != nil {
			return nil, fmt.Errorf("failed to get deck %d: %w", deckID, err)
		}
		if d == nil {
			return nil, nil
		}

		entity := ToEntity(*d, d.Player.Name, d.Format.Name, commanderInfoFromModel(*d))

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

func GetAllByPod(
	deckRepo repos.DeckRepository,
	podRepo repos.PodRepository,
	gameResultRepo repos.GameResultRepository,
) GetAllByPodFunc {
	return func(ctx context.Context, podID int) ([]EntityWithStats, error) {
		playerIDs, err := podRepo.GetPlayerIDs(ctx, podID)
		if err != nil {
			return nil, fmt.Errorf("failed to get player IDs for pod %d: %w", podID, err)
		}
		if len(playerIDs) == 0 {
			return []EntityWithStats{}, nil
		}

		decks, err := deckRepo.GetAllByPlayerIDsHydrated(ctx, playerIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get decks for pod %d: %w", podID, err)
		}

		result := make([]EntityWithStats, 0, len(decks))
		for _, d := range decks {
			entity := ToEntity(d, d.Player.Name, d.Format.Name, commanderInfoFromModel(d))

			agg, err := gameResultRepo.GetStatsForDeck(ctx, d.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get stats for deck %d: %w", d.ID, err)
			}

			result = append(result, ToEntityWithStats(entity, agg))
		}

		return result, nil
	}
}

func assertCallerOwnsDeck(ctx context.Context, deckRepo repos.DeckRepository, deckID, callerPlayerID int) error {
	d, err := deckRepo.GetById(ctx, deckID)
	if err != nil {
		return fmt.Errorf("failed to fetch deck %d: %w", deckID, err)
	}
	if d == nil {
		return fmt.Errorf("deck %d not found", deckID)
	}
	if d.PlayerID != callerPlayerID {
		return fmt.Errorf("forbidden: deck %d does not belong to caller", deckID)
	}
	return nil
}

func Update(
	deckRepo repos.DeckRepository,
	deckCmdrRepo repos.DeckCommanderRepository,
) UpdateFunc {
	return func(ctx context.Context, deckID int, callerPlayerID int, fields UpdateFields) error {
		if err := assertCallerOwnsDeck(ctx, deckRepo, deckID, callerPlayerID); err != nil {
			return err
		}

		repoFields := deckRepository.UpdateFields{
			Name:     fields.Name,
			FormatID: fields.FormatID,
			Retired:  fields.Retired,
		}
		if err := deckRepo.Update(ctx, deckID, repoFields); err != nil {
			return fmt.Errorf("failed to update deck %d: %w", deckID, err)
		}

		if fields.CommanderID != nil {
			if err := deckCmdrRepo.DeleteByDeckID(ctx, deckID); err != nil {
				return fmt.Errorf("failed to clear commander for deck %d: %w", deckID, err)
			}
			if _, err := deckCmdrRepo.Add(ctx, deckID, *fields.CommanderID, fields.PartnerCommanderID); err != nil {
				return fmt.Errorf("failed to set commander for deck %d: %w", deckID, err)
			}
		}

		return nil
	}
}

func SoftDelete(deckRepo repos.DeckRepository) SoftDeleteFunc {
	return func(ctx context.Context, deckID int, callerPlayerID int) error {
		if err := assertCallerOwnsDeck(ctx, deckRepo, deckID, callerPlayerID); err != nil {
			return err
		}
		return deckRepo.SoftDelete(ctx, deckID)
	}
}

func Retire(deckRepo repos.DeckRepository) RetireFunc {
	return func(ctx context.Context, deckID int, callerPlayerID int) error {
		if err := assertCallerOwnsDeck(ctx, deckRepo, deckID, callerPlayerID); err != nil {
			return err
		}
		return deckRepo.Retire(ctx, deckID)
	}
}
