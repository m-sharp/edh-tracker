package format

import (
	"context"
	"fmt"

	repos "github.com/m-sharp/edh-tracker/lib/repositories"
)

// TODO: There will never by more than a handful of functions.
// TODO: Create a flyweight pattern here where a var cachedFormats map[int]Entity map will be held in memory and consulted first by GetByID.
// TODO: If cachedFormats is empty or missing the target ID for GetByID, populate it via GetAll.
// TODO: GetAll should return from cachedFormats and only populate it on a daily basis.

func GetAll(formatRepo repos.FormatRepository) GetAllFunc {
	return func(ctx context.Context) ([]Entity, error) {
		models, err := formatRepo.GetAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get formats: %w", err)
		}

		entities := make([]Entity, 0, len(models))
		for _, m := range models {
			entities = append(entities, ToEntity(m))
		}

		return entities, nil
	}
}

func GetByID(formatRepo repos.FormatRepository) GetByIDFunc {
	return func(ctx context.Context, id int) (*Entity, error) {
		m, err := formatRepo.GetById(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get format %d: %w", id, err)
		}
		if m == nil {
			return nil, nil
		}
		e := ToEntity(*m)
		return &e, nil
	}
}
