package testHelpers

import (
	"context"
	"time"

	"gorm.io/gorm"

	repos "github.com/m-sharp/edh-tracker/lib/repositories"
	deckRepo "github.com/m-sharp/edh-tracker/lib/repositories/deck"
	deckCommanderRepo "github.com/m-sharp/edh-tracker/lib/repositories/deckCommander"
	formatRepo "github.com/m-sharp/edh-tracker/lib/repositories/format"
	gameRepo "github.com/m-sharp/edh-tracker/lib/repositories/game"
	gameResultRepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
	playerRepo "github.com/m-sharp/edh-tracker/lib/repositories/player"
	playerPodRoleRepo "github.com/m-sharp/edh-tracker/lib/repositories/playerPodRole"
	podRepo "github.com/m-sharp/edh-tracker/lib/repositories/pod"
	podInviteRepo "github.com/m-sharp/edh-tracker/lib/repositories/podInvite"
	userRepo "github.com/m-sharp/edh-tracker/lib/repositories/user"
)

// Compile-time interface checks.
var (
	_ repos.GameResultRepository    = (*MockGameResultRepo)(nil)
	_ repos.DeckRepository          = (*MockDeckRepo)(nil)
	_ repos.DeckCommanderRepository = (*MockDeckCommanderRepo)(nil)
	_ repos.PodRepository           = (*MockPodRepo)(nil)
	_ repos.PlayerRepository        = (*MockPlayerRepo)(nil)
	_ repos.PlayerPodRoleRepository = (*MockPlayerPodRoleRepo)(nil)
	_ repos.PodInviteRepository     = (*MockPodInviteRepo)(nil)
	_ repos.GameRepository          = (*MockGameRepo)(nil)
	_ repos.FormatRepository        = (*MockFormatRepo)(nil)
	_ repos.UserRepository          = (*MockUserRepo)(nil)
)

// MockGameResultRepo implements repos.GameResultRepository.
type MockGameResultRepo struct {
	GetByGameIDFn        func(ctx context.Context, gameID int) ([]gameResultRepo.Model, error)
	GetByIDFn            func(ctx context.Context, resultID int) (*gameResultRepo.Model, error)
	GetStatsForPlayerFn  func(ctx context.Context, playerID int) (*gameResultRepo.Aggregate, error)
	GetStatsForDeckFn    func(ctx context.Context, deckID int) (*gameResultRepo.Aggregate, error)
	GetStatsForDecksFn          func(ctx context.Context, deckIDs []int) (map[int]*gameResultRepo.Aggregate, error)
	GetStatsForPlayersInPodFn   func(ctx context.Context, podID int, playerIDs []int) (map[int]*gameResultRepo.Aggregate, error)
	BulkAddFn                   func(ctx context.Context, results []gameResultRepo.Model) error
	AddFn                func(ctx context.Context, model gameResultRepo.Model) (int, error)
	UpdateFn             func(ctx context.Context, resultID, place, killCount, deckID int) error
	SoftDeleteFn         func(ctx context.Context, id int) error
}

func (m *MockGameResultRepo) GetByGameID(ctx context.Context, gameID int) ([]gameResultRepo.Model, error) {
	if m.GetByGameIDFn != nil {
		return m.GetByGameIDFn(ctx, gameID)
	}
	panic("unexpected call to GetByGameID")
}
func (m *MockGameResultRepo) GetByID(ctx context.Context, resultID int) (*gameResultRepo.Model, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, resultID)
	}
	panic("unexpected call to GetByID")
}
func (m *MockGameResultRepo) GetStatsForPlayer(ctx context.Context, playerID int) (*gameResultRepo.Aggregate, error) {
	if m.GetStatsForPlayerFn != nil {
		return m.GetStatsForPlayerFn(ctx, playerID)
	}
	panic("unexpected call to GetStatsForPlayer")
}
func (m *MockGameResultRepo) GetStatsForDeck(ctx context.Context, deckID int) (*gameResultRepo.Aggregate, error) {
	if m.GetStatsForDeckFn != nil {
		return m.GetStatsForDeckFn(ctx, deckID)
	}
	panic("unexpected call to GetStatsForDeck")
}
func (m *MockGameResultRepo) GetStatsForDecks(ctx context.Context, deckIDs []int) (map[int]*gameResultRepo.Aggregate, error) {
	if m.GetStatsForDecksFn != nil {
		return m.GetStatsForDecksFn(ctx, deckIDs)
	}
	panic("unexpected call to GetStatsForDecks")
}
func (m *MockGameResultRepo) GetStatsForPlayersInPod(ctx context.Context, podID int, playerIDs []int) (map[int]*gameResultRepo.Aggregate, error) {
	if m.GetStatsForPlayersInPodFn != nil {
		return m.GetStatsForPlayersInPodFn(ctx, podID, playerIDs)
	}
	panic("unexpected call to GetStatsForPlayersInPod")
}
func (m *MockGameResultRepo) Add(ctx context.Context, model gameResultRepo.Model) (int, error) {
	if m.AddFn != nil {
		return m.AddFn(ctx, model)
	}
	panic("unexpected call to Add")
}
func (m *MockGameResultRepo) BulkAdd(ctx context.Context, results []gameResultRepo.Model) error {
	if m.BulkAddFn != nil {
		return m.BulkAddFn(ctx, results)
	}
	panic("unexpected call to BulkAdd")
}
func (m *MockGameResultRepo) Update(ctx context.Context, resultID, place, killCount, deckID int) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, resultID, place, killCount, deckID)
	}
	panic("unexpected call to Update")
}
func (m *MockGameResultRepo) SoftDelete(ctx context.Context, id int) error {
	if m.SoftDeleteFn != nil {
		return m.SoftDeleteFn(ctx, id)
	}
	panic("unexpected call to SoftDelete")
}
func (m *MockGameResultRepo) StartTX(_ *gorm.DB) {}
func (m *MockGameResultRepo) EndTX()              {}
func (m *MockGameResultRepo) DB() *gorm.DB        { return nil }

// MockDeckRepo implements repos.DeckRepository.
// GetByIdCalled is set to true whenever GetById is invoked (exported tracking flag).
type MockDeckRepo struct {
	GetByIdFn                 func(ctx context.Context, deckID int) (*deckRepo.Model, error)
	AddFn                     func(ctx context.Context, playerID int, name string, formatID int) (int, error)
	UpdateFn                  func(ctx context.Context, deckID int, fields deckRepo.UpdateFields) error
	SoftDeleteFn              func(ctx context.Context, id int) error
	GetAllFn                  func(ctx context.Context) ([]deckRepo.Model, error)
	GetAllForPlayerFn         func(ctx context.Context, playerID int) ([]deckRepo.Model, error)
	GetAllByPlayerIDsFn       func(ctx context.Context, playerIDs []int) ([]deckRepo.Model, error)
	GetAllByPodPaginatedFn    func(ctx context.Context, podID, limit, offset int) ([]deckRepo.Model, int, error)
	GetAllByPlayerPaginatedFn func(ctx context.Context, playerID, limit, offset int) ([]deckRepo.Model, int, error)
	GetByIDHydratedFn         func(ctx context.Context, deckID int) (*deckRepo.Model, error)
	GetByIdCalled             bool
}

func (m *MockDeckRepo) GetAll(ctx context.Context) ([]deckRepo.Model, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(ctx)
	}
	panic("unexpected call to GetAll")
}
func (m *MockDeckRepo) GetAllForPlayer(ctx context.Context, playerID int) ([]deckRepo.Model, error) {
	if m.GetAllForPlayerFn != nil {
		return m.GetAllForPlayerFn(ctx, playerID)
	}
	panic("unexpected call to GetAllForPlayer")
}
func (m *MockDeckRepo) GetAllByPlayerIDs(ctx context.Context, playerIDs []int) ([]deckRepo.Model, error) {
	if m.GetAllByPlayerIDsFn != nil {
		return m.GetAllByPlayerIDsFn(ctx, playerIDs)
	}
	panic("unexpected call to GetAllByPlayerIDs")
}
func (m *MockDeckRepo) GetAllByPodPaginated(ctx context.Context, podID, limit, offset int) ([]deckRepo.Model, int, error) {
	if m.GetAllByPodPaginatedFn != nil {
		return m.GetAllByPodPaginatedFn(ctx, podID, limit, offset)
	}
	panic("unexpected call to GetAllByPodPaginated")
}
func (m *MockDeckRepo) GetAllByPlayerPaginated(ctx context.Context, playerID, limit, offset int) ([]deckRepo.Model, int, error) {
	if m.GetAllByPlayerPaginatedFn != nil {
		return m.GetAllByPlayerPaginatedFn(ctx, playerID, limit, offset)
	}
	panic("unexpected call to GetAllByPlayerPaginated")
}
func (m *MockDeckRepo) GetById(ctx context.Context, deckID int) (*deckRepo.Model, error) {
	m.GetByIdCalled = true
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, deckID)
	}
	panic("unexpected call to GetById")
}
func (m *MockDeckRepo) GetByIDHydrated(ctx context.Context, deckID int) (*deckRepo.Model, error) {
	if m.GetByIDHydratedFn != nil {
		return m.GetByIDHydratedFn(ctx, deckID)
	}
	panic("unexpected call to GetByIDHydrated")
}
func (m *MockDeckRepo) Add(ctx context.Context, playerID int, name string, formatID int) (int, error) {
	if m.AddFn != nil {
		return m.AddFn(ctx, playerID, name, formatID)
	}
	panic("unexpected call to Add")
}
func (m *MockDeckRepo) BulkAdd(ctx context.Context, decks []deckRepo.Model) ([]deckRepo.Model, error) {
	panic("unexpected call to BulkAdd")
}
func (m *MockDeckRepo) Update(ctx context.Context, deckID int, fields deckRepo.UpdateFields) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, deckID, fields)
	}
	panic("unexpected call to Update")
}
func (m *MockDeckRepo) Retire(ctx context.Context, deckID int) error {
	panic("unexpected call to Retire")
}
func (m *MockDeckRepo) SoftDelete(ctx context.Context, id int) error {
	if m.SoftDeleteFn != nil {
		return m.SoftDeleteFn(ctx, id)
	}
	panic("unexpected call to SoftDelete")
}
func (m *MockDeckRepo) StartTX(_ *gorm.DB) {}
func (m *MockDeckRepo) EndTX()             {}
func (m *MockDeckRepo) DB() *gorm.DB       { return nil }

// MockDeckCommanderRepo implements repos.DeckCommanderRepository.
// AddCalled is set to true whenever Add is invoked (exported tracking flag).
type MockDeckCommanderRepo struct {
	GetByDeckIdFn    func(ctx context.Context, deckID int) (*deckCommanderRepo.Model, error)
	AddFn            func(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error)
	DeleteByDeckIDFn func(ctx context.Context, deckID int) error
	AddCalled        bool
}

func (m *MockDeckCommanderRepo) GetByDeckId(ctx context.Context, deckID int) (*deckCommanderRepo.Model, error) {
	if m.GetByDeckIdFn != nil {
		return m.GetByDeckIdFn(ctx, deckID)
	}
	panic("unexpected call to GetByDeckId")
}
func (m *MockDeckCommanderRepo) Add(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error) {
	m.AddCalled = true
	if m.AddFn != nil {
		return m.AddFn(ctx, deckID, commanderID, partnerCommanderID)
	}
	panic("unexpected call to Add")
}
func (m *MockDeckCommanderRepo) BulkAdd(ctx context.Context, entries []deckCommanderRepo.Model) error {
	panic("unexpected call to BulkAdd")
}
func (m *MockDeckCommanderRepo) DeleteByDeckID(ctx context.Context, deckID int) error {
	if m.DeleteByDeckIDFn != nil {
		return m.DeleteByDeckIDFn(ctx, deckID)
	}
	panic("unexpected call to DeleteByDeckID")
}
func (m *MockDeckCommanderRepo) StartTX(_ *gorm.DB) {}
func (m *MockDeckCommanderRepo) EndTX()             {}
func (m *MockDeckCommanderRepo) DB() *gorm.DB       { return nil }

// MockPodRepo implements repos.PodRepository.
type MockPodRepo struct {
	GetByIDFn          func(ctx context.Context, podID int) (*podRepo.Model, error)
	GetByPlayerIDFn    func(ctx context.Context, playerID int) ([]podRepo.Model, error)
	AddFn              func(ctx context.Context, name string) (int, error)
	AddPlayerToPodFn   func(ctx context.Context, podID, playerID int) error
	RemovePlayerFn     func(ctx context.Context, podID, playerID int) error
	SoftDeleteFn       func(ctx context.Context, podID int) error
	UpdateFn           func(ctx context.Context, podID int, name string) error
	GetIDsByPlayerIDFn func(ctx context.Context, playerID int) ([]int, error)
	GetPlayerIDsFn     func(ctx context.Context, podID int) ([]int, error)
}

func (m *MockPodRepo) GetAll(ctx context.Context) ([]podRepo.Model, error) {
	panic("unexpected call to GetAll")
}
func (m *MockPodRepo) GetByIDWithMembers(ctx context.Context, podID int) (*podRepo.Model, error) {
	panic("unexpected call to GetByIDWithMembers")
}
func (m *MockPodRepo) GetByID(ctx context.Context, podID int) (*podRepo.Model, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, podID)
	}
	panic("unexpected call to GetByID")
}
func (m *MockPodRepo) GetByPlayerID(ctx context.Context, playerID int) ([]podRepo.Model, error) {
	if m.GetByPlayerIDFn != nil {
		return m.GetByPlayerIDFn(ctx, playerID)
	}
	panic("unexpected call to GetByPlayerID")
}
func (m *MockPodRepo) GetByName(ctx context.Context, name string) (*podRepo.Model, error) {
	panic("unexpected call to GetByName")
}
func (m *MockPodRepo) GetIDsByPlayerID(ctx context.Context, playerID int) ([]int, error) {
	if m.GetIDsByPlayerIDFn != nil {
		return m.GetIDsByPlayerIDFn(ctx, playerID)
	}
	panic("unexpected call to GetIDsByPlayerID")
}
func (m *MockPodRepo) GetPlayerIDs(ctx context.Context, podID int) ([]int, error) {
	if m.GetPlayerIDsFn != nil {
		return m.GetPlayerIDsFn(ctx, podID)
	}
	panic("unexpected call to GetPlayerIDs")
}
func (m *MockPodRepo) Add(ctx context.Context, name string) (int, error) {
	if m.AddFn != nil {
		return m.AddFn(ctx, name)
	}
	panic("unexpected call to Add")
}
func (m *MockPodRepo) BulkAddPlayers(ctx context.Context, podID int, playerIDs []int) error {
	panic("unexpected call to BulkAddPlayers")
}
func (m *MockPodRepo) AddPlayerToPod(ctx context.Context, podID, playerID int) error {
	if m.AddPlayerToPodFn != nil {
		return m.AddPlayerToPodFn(ctx, podID, playerID)
	}
	panic("unexpected call to AddPlayerToPod")
}
func (m *MockPodRepo) SoftDelete(ctx context.Context, podID int) error {
	if m.SoftDeleteFn != nil {
		return m.SoftDeleteFn(ctx, podID)
	}
	panic("unexpected call to SoftDelete")
}
func (m *MockPodRepo) Update(ctx context.Context, podID int, name string) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, podID, name)
	}
	panic("unexpected call to Update")
}
func (m *MockPodRepo) RemovePlayer(ctx context.Context, podID, playerID int) error {
	if m.RemovePlayerFn != nil {
		return m.RemovePlayerFn(ctx, podID, playerID)
	}
	panic("unexpected call to RemovePlayer")
}
func (m *MockPodRepo) StartTX(_ *gorm.DB) {}
func (m *MockPodRepo) EndTX()             {}
func (m *MockPodRepo) DB() *gorm.DB       { return nil }

// MockPlayerRepo implements repos.PlayerRepository.
type MockPlayerRepo struct {
	GetAllFn    func(ctx context.Context) ([]playerRepo.Model, error)
	GetByIdFn   func(ctx context.Context, playerID int) (*playerRepo.Model, error)
	GetByNameFn func(ctx context.Context, name string) (*playerRepo.Model, error)
	UpdateFn    func(ctx context.Context, playerID int, name string) error
}

func (m *MockPlayerRepo) GetAll(ctx context.Context) ([]playerRepo.Model, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(ctx)
	}
	panic("unexpected call to GetAll")
}
func (m *MockPlayerRepo) GetById(ctx context.Context, playerID int) (*playerRepo.Model, error) {
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, playerID)
	}
	panic("unexpected call to GetById")
}
func (m *MockPlayerRepo) GetByName(ctx context.Context, name string) (*playerRepo.Model, error) {
	if m.GetByNameFn != nil {
		return m.GetByNameFn(ctx, name)
	}
	panic("unexpected call to GetByName")
}
func (m *MockPlayerRepo) GetByNames(ctx context.Context, names []string) ([]playerRepo.Model, error) {
	panic("unexpected call to GetByNames")
}
func (m *MockPlayerRepo) Add(ctx context.Context, name string) (int, error) {
	panic("unexpected call to Add")
}
func (m *MockPlayerRepo) BulkAdd(ctx context.Context, names []string) ([]playerRepo.Model, error) {
	panic("unexpected call to BulkAdd")
}
func (m *MockPlayerRepo) Update(ctx context.Context, playerID int, name string) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, playerID, name)
	}
	panic("unexpected call to Update")
}
func (m *MockPlayerRepo) SoftDelete(ctx context.Context, id int) error {
	panic("unexpected call to SoftDelete")
}
func (m *MockPlayerRepo) StartTX(_ *gorm.DB) {}
func (m *MockPlayerRepo) EndTX()             {}
func (m *MockPlayerRepo) DB() *gorm.DB       { return nil }

// MockPlayerPodRoleRepo implements repos.PlayerPodRoleRepository.
type MockPlayerPodRoleRepo struct {
	GetRoleFn             func(ctx context.Context, podID, playerID int) (*playerPodRoleRepo.Model, error)
	SetRoleFn             func(ctx context.Context, podID, playerID int, role string) error
	GetMembersWithRolesFn func(ctx context.Context, podID int) ([]playerPodRoleRepo.Model, error)
}

func (m *MockPlayerPodRoleRepo) GetRole(ctx context.Context, podID, playerID int) (*playerPodRoleRepo.Model, error) {
	if m.GetRoleFn != nil {
		return m.GetRoleFn(ctx, podID, playerID)
	}
	panic("unexpected call to GetRole")
}
func (m *MockPlayerPodRoleRepo) SetRole(ctx context.Context, podID, playerID int, role string) error {
	if m.SetRoleFn != nil {
		return m.SetRoleFn(ctx, podID, playerID, role)
	}
	panic("unexpected call to SetRole")
}
func (m *MockPlayerPodRoleRepo) GetMembersWithRoles(ctx context.Context, podID int) ([]playerPodRoleRepo.Model, error) {
	if m.GetMembersWithRolesFn != nil {
		return m.GetMembersWithRolesFn(ctx, podID)
	}
	panic("unexpected call to GetMembersWithRoles")
}
func (m *MockPlayerPodRoleRepo) BulkAdd(ctx context.Context, podID int, playerIDs []int, role string) error {
	panic("unexpected call to BulkAdd")
}
func (m *MockPlayerPodRoleRepo) StartTX(_ *gorm.DB) {}
func (m *MockPlayerPodRoleRepo) EndTX()             {}
func (m *MockPlayerPodRoleRepo) DB() *gorm.DB       { return nil }

// MockPodInviteRepo implements repos.PodInviteRepository.
type MockPodInviteRepo struct {
	GetByCodeFn          func(ctx context.Context, code string) (*podInviteRepo.Model, error)
	AddFn                func(ctx context.Context, podID, createdByPlayerID int, code string, expiresAt *time.Time) error
	IncrementUsedCountFn func(ctx context.Context, code string) error
}

func (m *MockPodInviteRepo) GetByCode(ctx context.Context, code string) (*podInviteRepo.Model, error) {
	if m.GetByCodeFn != nil {
		return m.GetByCodeFn(ctx, code)
	}
	panic("unexpected call to GetByCode")
}
func (m *MockPodInviteRepo) Add(ctx context.Context, podID, createdByPlayerID int, code string, expiresAt *time.Time) error {
	if m.AddFn != nil {
		return m.AddFn(ctx, podID, createdByPlayerID, code, expiresAt)
	}
	panic("unexpected call to Add")
}
func (m *MockPodInviteRepo) IncrementUsedCount(ctx context.Context, code string) error {
	if m.IncrementUsedCountFn != nil {
		return m.IncrementUsedCountFn(ctx, code)
	}
	panic("unexpected call to IncrementUsedCount")
}
func (m *MockPodInviteRepo) StartTX(_ *gorm.DB) {}
func (m *MockPodInviteRepo) EndTX()             {}
func (m *MockPodInviteRepo) DB() *gorm.DB       { return nil }

// MockGameRepo implements repos.GameRepository.
type MockGameRepo struct {
	GetAllByPodFn               func(ctx context.Context, podID int) ([]gameRepo.Model, error)
	GetAllByDeckFn              func(ctx context.Context, deckID int) ([]gameRepo.Model, error)
	GetAllByPlayerIDFn          func(ctx context.Context, playerID int) ([]gameRepo.Model, error)
	GetAllByPodPaginatedFn      func(ctx context.Context, podID, limit, offset int) ([]gameRepo.Model, int, error)
	GetAllByDeckPaginatedFn     func(ctx context.Context, deckID, limit, offset int) ([]gameRepo.Model, int, error)
	GetAllByPlayerIDPaginatedFn func(ctx context.Context, playerID, limit, offset int) ([]gameRepo.Model, int, error)
	GetByIDFn                   func(ctx context.Context, gameID int) (*gameRepo.Model, error)
	AddFn                       func(ctx context.Context, description string, podID, formatID int) (int, error)
	UpdateFn                    func(ctx context.Context, gameID int, description string) error
	SoftDeleteFn                func(ctx context.Context, id int) error
}

func (m *MockGameRepo) GetAllByPod(ctx context.Context, podID int) ([]gameRepo.Model, error) {
	if m.GetAllByPodFn != nil {
		return m.GetAllByPodFn(ctx, podID)
	}
	panic("unexpected call to GetAllByPod")
}
func (m *MockGameRepo) GetAllByDeck(ctx context.Context, deckID int) ([]gameRepo.Model, error) {
	if m.GetAllByDeckFn != nil {
		return m.GetAllByDeckFn(ctx, deckID)
	}
	panic("unexpected call to GetAllByDeck")
}
func (m *MockGameRepo) GetAllByPlayerID(ctx context.Context, playerID int) ([]gameRepo.Model, error) {
	if m.GetAllByPlayerIDFn != nil {
		return m.GetAllByPlayerIDFn(ctx, playerID)
	}
	panic("unexpected call to GetAllByPlayerID")
}
func (m *MockGameRepo) GetAllByPodPaginated(ctx context.Context, podID, limit, offset int) ([]gameRepo.Model, int, error) {
	if m.GetAllByPodPaginatedFn != nil {
		return m.GetAllByPodPaginatedFn(ctx, podID, limit, offset)
	}
	panic("unexpected call to GetAllByPodPaginated")
}
func (m *MockGameRepo) GetAllByDeckPaginated(ctx context.Context, deckID, limit, offset int) ([]gameRepo.Model, int, error) {
	if m.GetAllByDeckPaginatedFn != nil {
		return m.GetAllByDeckPaginatedFn(ctx, deckID, limit, offset)
	}
	panic("unexpected call to GetAllByDeckPaginated")
}
func (m *MockGameRepo) GetAllByPlayerIDPaginated(ctx context.Context, playerID, limit, offset int) ([]gameRepo.Model, int, error) {
	if m.GetAllByPlayerIDPaginatedFn != nil {
		return m.GetAllByPlayerIDPaginatedFn(ctx, playerID, limit, offset)
	}
	panic("unexpected call to GetAllByPlayerIDPaginated")
}
func (m *MockGameRepo) GetByID(ctx context.Context, gameID int) (*gameRepo.Model, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, gameID)
	}
	panic("unexpected call to GetByID")
}
func (m *MockGameRepo) Add(ctx context.Context, description string, podID, formatID int) (int, error) {
	if m.AddFn != nil {
		return m.AddFn(ctx, description, podID, formatID)
	}
	panic("unexpected call to Add")
}
func (m *MockGameRepo) BulkAdd(ctx context.Context, games []gameRepo.Model) ([]int, error) {
	panic("unexpected call to BulkAdd")
}
func (m *MockGameRepo) Update(ctx context.Context, gameID int, description string) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, gameID, description)
	}
	panic("unexpected call to Update")
}
func (m *MockGameRepo) SoftDelete(ctx context.Context, id int) error {
	if m.SoftDeleteFn != nil {
		return m.SoftDeleteFn(ctx, id)
	}
	panic("unexpected call to SoftDelete")
}
func (m *MockGameRepo) StartTX(_ *gorm.DB) {}
func (m *MockGameRepo) EndTX()             {}
func (m *MockGameRepo) DB() *gorm.DB       { return nil }

// MockFormatRepo implements repos.FormatRepository.
type MockFormatRepo struct {
	GetAllFn    func(ctx context.Context) ([]formatRepo.Model, error)
	GetByIdFn   func(ctx context.Context, id int) (*formatRepo.Model, error)
	GetByNameFn func(ctx context.Context, name string) (*formatRepo.Model, error)
}

func (m *MockFormatRepo) GetAll(ctx context.Context) ([]formatRepo.Model, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(ctx)
	}
	panic("unexpected call to GetAll")
}
func (m *MockFormatRepo) GetById(ctx context.Context, id int) (*formatRepo.Model, error) {
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, id)
	}
	panic("unexpected call to GetById")
}
func (m *MockFormatRepo) GetByName(ctx context.Context, name string) (*formatRepo.Model, error) {
	if m.GetByNameFn != nil {
		return m.GetByNameFn(ctx, name)
	}
	panic("unexpected call to GetByName")
}
func (m *MockFormatRepo) StartTX(_ *gorm.DB) {}
func (m *MockFormatRepo) EndTX()             {}
func (m *MockFormatRepo) DB() *gorm.DB       { return nil }

// MockUserRepo implements repos.UserRepository.
type MockUserRepo struct {
	GetByIDFn              func(ctx context.Context, id int) (*userRepo.Model, error)
	GetByPlayerIDFn        func(ctx context.Context, playerID int) (*userRepo.Model, error)
	GetByOAuthFn           func(ctx context.Context, provider, subject string) (*userRepo.Model, error)
	GetRoleByNameFn        func(ctx context.Context, name string) (*userRepo.RoleModel, error)
	AddFn                  func(ctx context.Context, playerID, roleID int) (int, error)
	AddWithOAuthFn         func(ctx context.Context, playerID, roleID int, provider, subject, email, displayName, avatarURL string) (int, error)
	BulkAddFn              func(ctx context.Context, playerIDs []int, roleID int) error
	SoftDeleteFn           func(ctx context.Context, id int) error
	GetByEmailFn           func(ctx context.Context, email string) (*userRepo.Model, error)
	UpdateOAuthFn          func(ctx context.Context, userID int, provider, subject, email, displayName, avatarURL string) error
	SetEmailFn             func(ctx context.Context, playerID int, email string) error
}

func (m *MockUserRepo) GetByID(ctx context.Context, id int) (*userRepo.Model, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	panic("unexpected call to GetByID")
}
func (m *MockUserRepo) GetByPlayerID(ctx context.Context, playerID int) (*userRepo.Model, error) {
	if m.GetByPlayerIDFn != nil {
		return m.GetByPlayerIDFn(ctx, playerID)
	}
	panic("unexpected call to GetByPlayerID")
}
func (m *MockUserRepo) GetByOAuth(ctx context.Context, provider, subject string) (*userRepo.Model, error) {
	if m.GetByOAuthFn != nil {
		return m.GetByOAuthFn(ctx, provider, subject)
	}
	panic("unexpected call to GetByOAuth")
}
func (m *MockUserRepo) GetRoleByName(ctx context.Context, name string) (*userRepo.RoleModel, error) {
	if m.GetRoleByNameFn != nil {
		return m.GetRoleByNameFn(ctx, name)
	}
	panic("unexpected call to GetRoleByName")
}
func (m *MockUserRepo) Add(ctx context.Context, playerID, roleID int) (int, error) {
	if m.AddFn != nil {
		return m.AddFn(ctx, playerID, roleID)
	}
	panic("unexpected call to Add")
}
func (m *MockUserRepo) AddWithOAuth(ctx context.Context, playerID, roleID int, provider, subject, email, displayName, avatarURL string) (int, error) {
	if m.AddWithOAuthFn != nil {
		return m.AddWithOAuthFn(ctx, playerID, roleID, provider, subject, email, displayName, avatarURL)
	}
	panic("unexpected call to AddWithOAuth")
}
func (m *MockUserRepo) BulkAdd(ctx context.Context, playerIDs []int, roleID int) error {
	if m.BulkAddFn != nil {
		return m.BulkAddFn(ctx, playerIDs, roleID)
	}
	panic("unexpected call to BulkAdd")
}
func (m *MockUserRepo) SoftDelete(ctx context.Context, id int) error {
	if m.SoftDeleteFn != nil {
		return m.SoftDeleteFn(ctx, id)
	}
	panic("unexpected call to SoftDelete")
}
func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*userRepo.Model, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(ctx, email)
	}
	panic("unexpected call to GetByEmail")
}
func (m *MockUserRepo) UpdateOAuth(ctx context.Context, userID int, provider, subject, email, displayName, avatarURL string) error {
	if m.UpdateOAuthFn != nil {
		return m.UpdateOAuthFn(ctx, userID, provider, subject, email, displayName, avatarURL)
	}
	panic("unexpected call to UpdateOAuth")
}
func (m *MockUserRepo) SetEmail(ctx context.Context, playerID int, email string) error {
	if m.SetEmailFn != nil {
		return m.SetEmailFn(ctx, playerID, email)
	}
	panic("unexpected call to SetEmail")
}
func (m *MockUserRepo) StartTX(_ *gorm.DB) {}
func (m *MockUserRepo) EndTX()             {}
func (m *MockUserRepo) DB() *gorm.DB       { return nil }
