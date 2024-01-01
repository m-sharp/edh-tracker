package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	dropDecks = `DELETE * FROM deck;`
)

var (
	deckSeeds = []string{
		// Mike Decks
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Go-Shintai of Life's Origin", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Rakdos, Lord of Riots", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Tovolar, Dire Overlord", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Korvold, Fae-Cursed King", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Miirym, Sentinel Wyrm", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Zaxara, the Exemplary", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Isshin, Two Heavens as One", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Kathril, Aspect Warper", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Raffine, Scheming Seer", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Syr Gwyn, Hero of Ashvale", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Admiral Beckett Brass", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Giada, Font of Hope", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Aragorn, King of Gondor", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Mike"), "Thantis, the Warweaver", now());`,
		// Tom Decks
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Zethi, Arcane Blademaster", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Old Stickfingers", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Grothama, All-Devouring", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Dina, Soul Steeper", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Sephara, Sky's Blade", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Vadrik, Astral Archmage", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Grazilaxx, Illithid Scholar", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Firesong and Sunspeaker", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Atla Palani, Nest Tender", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Yurlok of Scorch Thrash", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Ruric Thar, the Unbowed", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Kros, Defense Contractor", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Plargg and Nassari", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Kazuul, Tyrant of the Cliffs", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Azusa, Lost but Seeking", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Feather, the Redeemed", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Lurrus of the Dream-Den", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Tom"), "Gyome, Master Chef", now());`,
		// Dillon Decks
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Dillon"), "Jhoira, Weatherlight Captain", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Dillon"), "Sakashima of a Thousand Faces", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Dillon"), "Nekusar, the Mindrazer", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Dillon"), "Tergrid, God of Fright", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Dillon"), "Magus Lucea Kane", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Dillon"), "Kelsien, the Plague", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Dillon"), "Arcades, the Stategist", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Dillon"), "Sythis, Harvest's Hand", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Dillon"), "Oona, Queen of the Fae", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Dillon"), "Satoru Umezawa", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Dillon"), "Magar of the Magic Strings", now());`,
		// Peter Decks
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Thalisse, Reverent Medium", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Atraxa, Praetor's Voice", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Tazri, Beacon of Unity", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Zimone and Dina", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Ruxa, Patient Professor", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "The Scorpion God", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "The Ur-Dragon", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Umbris, Fear Manifest", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Ganax, Astral Hunter", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Hofri Ghostforge", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Sauron, the Dark Lord", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Henzie 'Toolbox' Torre", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Eriette of the Charmed Apple", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Lord of the Nazgûl", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Kumena, Tyrant of Orazca", now());`,
		`INSERT INTO deck (player_id, commander, ctime) VALUES ((SELECT id FROM player where name = "Peter"), "Wyll, Blade of Frontiers", now());`,
	}
)

type Migration5 struct{}

func (m *Migration5) Upgrade(ctx context.Context, client *lib.DBClient) error {
	for _, deckSeed := range deckSeeds {
		if _, err := client.Db.ExecContext(ctx, deckSeed); err != nil {
			return lib.NewDBError(deckSeed, err)
		}
	}
	return nil
}

func (m *Migration5) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, dropDecks); err != nil {
		return lib.NewDBError(dropDecks, err)
	}
	return nil
}
