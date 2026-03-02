package format

import (
	"context"
	"fmt"
	"sync"

	repos "github.com/m-sharp/edh-tracker/lib/repositories"
)

var cache struct {
	sync.RWMutex
	m map[int]Entity
}

func GetAll(formatRepo repos.FormatRepository) GetAllFunc {
	return func(ctx context.Context) ([]Entity, error) {
		if err := ensureCache(ctx, formatRepo); err != nil {
			return nil, err
		}
		cache.RLock()
		defer cache.RUnlock()
		entities := make([]Entity, 0, len(cache.m))
		for _, e := range cache.m {
			entities = append(entities, e)
		}
		return entities, nil
	}
}

func GetByID(formatRepo repos.FormatRepository) GetByIDFunc {
	return func(ctx context.Context, id int) (*Entity, error) {
		if err := ensureCache(ctx, formatRepo); err != nil {
			return nil, fmt.Errorf("failed to get format %d: %w", id, err)
		}
		cache.RLock()
		defer cache.RUnlock()
		if e, ok := cache.m[id]; ok {
			return &e, nil
		}
		return nil, nil
	}
}

func ensureCache(ctx context.Context, formatRepo repos.FormatRepository) error {
	cache.RLock()
	populated := cache.m != nil
	cache.RUnlock()
	if populated {
		return nil
	}

	cache.Lock()
	defer cache.Unlock()
	if cache.m != nil {
		return nil
	}

	models, err := formatRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get formats: %w", err)
	}

	cache.m = make(map[int]Entity, len(models))
	for _, m := range models {
		e := ToEntity(m)
		cache.m[e.ID] = e
	}
	return nil
}
