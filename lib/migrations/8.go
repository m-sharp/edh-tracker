package migrations

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	dropGames    = `DELETE * FROM game;`
	insertResult = `INSERT INTO game_result (game_id, deck_id, place, kill_count) VALUES (
		(SELECT id FROM game where description = ?),
		(SELECT id FROM deck where commander = ? AND player_id = (SELECT id FROM player where name = ?)),
		?,
		?
	);`
)

type resultInfo struct {
	Player    string
	Commander string
	Place     int
	Kills     int
}

// ToDo: This should really live in a seed file rather than a migration
var (
	gameSeeds = []string{
		`INSERT INTO game (description, ctime) VALUES ('Game 1', '2023-04-06');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 2', '2023-04-06');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 3', '2023-04-06');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 4', '2023-04-15');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 5', '2023-04-15');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 6', '2023-04-15');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 7', '2023-04-15');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 8', '2023-04-18');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 9', '2023-04-18');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 10', '2023-04-18');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 11', '2023-04-27');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 12', '2023-04-27');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 13', '2023-04-27');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 14', '2023-05-10');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 15', '2023-05-10');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 16', '2023-05-18');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 17', '2023-05-18');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 18', '2023-05-18');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 19', '2023-05-24');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 20', '2023-05-24');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 21', '2023-05-24');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 22', '2023-06-14');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 23', '2023-06-14');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 24', '2023-07-11');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 25', '2023-09-28');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 26', '2023-09-28');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 27', '2023-09-28');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 28', '2023-11-02');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 29', '2023-11-02');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 30', '2023-11-02');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 31', '2023-11-08');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 32', '2023-11-08');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 33', '2023-11-08');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 34', '2023-11-29');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 35', '2023-11-29');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 36', '2023-11-29');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 37', '2023-12-12');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 38', '2023-12-12');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 39', '2023-12-12');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 40', '2023-12-30');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 41', '2023-12-30');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 42', '2023-12-30');`,
		`INSERT INTO game (description, ctime) VALUES ('Game 43', '2023-12-30');`,
	}

	gameResultSeeds = map[string][]resultInfo{
		"Game 1": {
			{
				Player:    "Peter",
				Commander: "Thalisse, Reverent Medium",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Tom",
				Commander: "Zethi, Arcane Blademaster",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Go-Shintai of Life's Origin",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Jhoira, Weatherlight Captain",
				Place:     1,
				Kills:     3,
			},
		},
		"Game 2": {
			{
				Player:    "Peter",
				Commander: "Atraxa, Praetor's Voice",
				Place:     1,
				Kills:     1,
			},
			{
				Player:    "Tom",
				Commander: "Old Stickfingers",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Rakdos, Lord of Riots",
				Place:     2,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Sakashima of a Thousand Faces",
				Place:     3,
				Kills:     1,
			},
		},
		"Game 3": {
			{
				Player:    "Peter",
				Commander: "Tazri, Beacon of Unity",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Tom",
				Commander: "Grothama, All-Devouring",
				Place:     2,
				Kills:     1,
			},
			{
				Player:    "Mike",
				Commander: "Tovolar, Dire Overlord",
				Place:     1,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Nekusar, the Mindrazer",
				Place:     3,
				Kills:     1,
			},
		},
		"Game 4": {
			{
				Player:    "Peter",
				Commander: "Zimone and Dina",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Tom",
				Commander: "Dina, Soul Steeper",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Tergrid, God of Fright",
				Place:     1,
				Kills:     1,
			},
			{
				Player:    "Mike",
				Commander: "Go-Shintai of Life's Origin",
				Place:     2,
				Kills:     0,
			},
		},
		"Game 5": {
			{
				Player:    "Peter",
				Commander: "Ruxa, Patient Professor",
				Place:     1,
				Kills:     3,
			},
			{
				Player:    "Tom",
				Commander: "Sephara, Sky's Blade",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Sakashima of a Thousand Faces",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Korvold, Fae-Cursed King",
				Place:     2,
				Kills:     0,
			},
		},
		"Game 6": {
			{
				Player:    "Peter",
				Commander: "The Scorpion God",
				Place:     3,
				Kills:     1,
			},
			{
				Player:    "Tom",
				Commander: "Vadrik, Astral Archmage",
				Place:     1,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Jhoira, Weatherlight Captain",
				Place:     4,
				Kills:     1,
			},
			{
				Player:    "Mike",
				Commander: "Korvold, Fae-Cursed King",
				Place:     2,
				Kills:     0,
			},
		},
		"Game 7": {
			{
				Player:    "Peter",
				Commander: "Tazri, Beacon of Unity",
				Place:     1,
				Kills:     2,
			},
			{
				Player:    "Tom",
				Commander: "Grazilaxx, Illithid Scholar",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Miirym, Sentinel Wyrm",
				Place:     2,
				Kills:     0,
			},
		},
		"Game 8": {
			{
				Player:    "Tom",
				Commander: "Firesong and Sunspeaker",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Zaxara, the Exemplary",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Tergrid, God of Fright",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Ruxa, Patient Professor",
				Place:     1,
				Kills:     3,
			},
		},
		"Game 9": {
			{
				Player:    "Tom",
				Commander: "Atla Palani, Nest Tender",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Isshin, Two Heavens as One",
				Place:     1,
				Kills:     2,
			},
			{
				Player:    "Dillon",
				Commander: "Magus Lucea Kane",
				Place:     2,
				Kills:     1,
			},
			{
				Player:    "Peter",
				Commander: "Zimone and Dina",
				Place:     4,
				Kills:     0,
			},
		},
		"Game 10": {
			{
				Player:    "Tom",
				Commander: "Sephara, Sky's Blade",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Tovolar, Dire Overlord",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Sakashima of a Thousand Faces",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Thalisse, Reverent Medium",
				Place:     1,
				Kills:     3,
			},
		},
		"Game 11": {
			{
				Player:    "Tom",
				Commander: "Vadrik, Astral Archmage",
				Place:     3,
				Kills:     1,
			},
			{
				Player:    "Mike",
				Commander: "Miirym, Sentinel Wyrm",
				Place:     1,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Kelsien, the Plague",
				Place:     3,
				Kills:     1,
			},
			{
				Player:    "Peter",
				Commander: "The Ur-Dragon",
				Place:     4,
				Kills:     0,
			},
		},
		"Game 12": {
			{
				Player:    "Tom",
				Commander: "Yurlok of Scorch Thrash",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Kathril, Aspect Warper",
				Place:     1,
				Kills:     3,
			},
			{
				Player:    "Dillon",
				Commander: "Kelsien, the Plague",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "The Ur-Dragon",
				Place:     4,
				Kills:     0,
			},
		},
		"Game 13": {
			{
				Player:    "Tom",
				Commander: "Sephara, Sky's Blade",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Rakdos, Lord of Riots",
				Place:     3,
				Kills:     2,
			},
			{
				Player:    "Dillon",
				Commander: "Tergrid, God of Fright",
				Place:     1,
				Kills:     1,
			},
			{
				Player:    "Peter",
				Commander: "Umbris, Fear Manifest",
				Place:     4,
				Kills:     0,
			},
		},
		"Game 14": {
			{
				Player:    "Tom",
				Commander: "Sephara, Sky's Blade",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Raffine, Scheming Seer",
				Place:     1,
				Kills:     2,
			},
			{
				Player:    "Dillon",
				Commander: "Kelsien, the Plague",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Ganax, Astral Hunter",
				Place:     3,
				Kills:     1,
			},
		},
		"Game 15": {
			{
				Player:    "Tom",
				Commander: "Firesong and Sunspeaker",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Kathril, Aspect Warper",
				Place:     2,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Arcades, the Stategist",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Zimone and Dina",
				Place:     1,
				Kills:     2,
			},
		},
		"Game 16": {
			{
				Player:    "Tom",
				Commander: "Ruric Thar, the Unbowed",
				Place:     2,
				Kills:     1,
			},
			{
				Player:    "Mike",
				Commander: "Korvold, Fae-Cursed King",
				Place:     1,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Tergrid, God of Fright",
				Place:     3,
				Kills:     1,
			},
			{
				Player:    "Peter",
				Commander: "The Scorpion God",
				Place:     4,
				Kills:     0,
			},
		},
		"Game 17": {
			{
				Player:    "Tom",
				Commander: "Grazilaxx, Illithid Scholar",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Raffine, Scheming Seer",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Sakashima of a Thousand Faces",
				Place:     2,
				Kills:     2,
			},
			{
				Player:    "Peter",
				Commander: "The Scorpion God",
				Place:     1,
				Kills:     1,
			},
		},
		"Game 18": {
			{
				Player:    "Tom",
				Commander: "Ruric Thar, the Unbowed",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Raffine, Scheming Seer",
				Place:     1,
				Kills:     3,
			},
			{
				Player:    "Dillon",
				Commander: "Jhoira, Weatherlight Captain",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Umbris, Fear Manifest",
				Place:     3,
				Kills:     0,
			},
		},
		"Game 19": {
			{
				Player:    "Tom",
				Commander: "Ruric Thar, the Unbowed",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Isshin, Two Heavens as One",
				Place:     1,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Sythis, Harvest's Hand",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Umbris, Fear Manifest",
				Place:     2,
				Kills:     2,
			},
		},
		"Game 20": {
			{
				Player:    "Tom",
				Commander: "Kros, Defense Contractor",
				Place:     1,
				Kills:     2,
			},
			{
				Player:    "Mike",
				Commander: "Raffine, Scheming Seer",
				Place:     2,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Sythis, Harvest's Hand",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Ruxa, Patient Professor",
				Place:     2,
				Kills:     0,
			},
		},
		"Game 21": {
			{
				Player:    "Tom",
				Commander: "Dina, Soul Steeper",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Syr Gwyn, Hero of Ashvale",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Tergrid, God of Fright",
				Place:     1,
				Kills:     3,
			},
			{
				Player:    "Peter",
				Commander: "Ruxa, Patient Professor",
				Place:     3,
				Kills:     0,
			},
		},
		"Game 22": {
			{
				Player:    "Tom",
				Commander: "Plargg and Nassari",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Admiral Beckett Brass",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Sythis, Harvest's Hand",
				Place:     1,
				Kills:     2,
			},
			{
				Player:    "Peter",
				Commander: "Hofri Ghostforge",
				Place:     2,
				Kills:     1,
			},
		},
		"Game 23": {
			{
				Player:    "Tom",
				Commander: "Kros, Defense Contractor",
				Place:     3,
				Kills:     1,
			},
			{
				Player:    "Mike",
				Commander: "Admiral Beckett Brass",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Kelsien, the Plague",
				Place:     1,
				Kills:     2,
			},
			{
				Player:    "Peter",
				Commander: "Umbris, Fear Manifest",
				Place:     4,
				Kills:     0,
			},
		},
		"Game 24": {
			{
				Player:    "Tom",
				Commander: "Plargg and Nassari",
				Place:     1,
				Kills:     3,
			},
			{
				Player:    "Mike",
				Commander: "Admiral Beckett Brass",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Sythis, Harvest's Hand",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Ruxa, Patient Professor",
				Place:     3,
				Kills:     1,
			},
		},
		"Game 25": {
			{
				Player:    "Tom",
				Commander: "Kazuul, Tyrant of the Cliffs",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Kathril, Aspect Warper",
				Place:     1,
				Kills:     3,
			},
			{
				Player:    "Dillon",
				Commander: "Sakashima of a Thousand Faces",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Sauron, the Dark Lord",
				Place:     2,
				Kills:     0,
			},
		},
		"Game 26": {
			{
				Player:    "Tom",
				Commander: "Azusa, Lost but Seeking",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Giada, Font of Hope",
				Place:     1,
				Kills:     3,
			},
			{
				Player:    "Dillon",
				Commander: "Tergrid, God of Fright",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: `Henzie 'Toolbox' Torre`,
				Place:     2,
				Kills:     0,
			},
		},
		"Game 27": {
			{
				Player:    "Tom",
				Commander: "Kazuul, Tyrant of the Cliffs",
				Place:     3,
				Kills:     1,
			},
			{
				Player:    "Mike",
				Commander: "Raffine, Scheming Seer",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Jhoira, Weatherlight Captain",
				Place:     1,
				Kills:     1,
			},
			{
				Player:    "Peter",
				Commander: "Eriette of the Charmed Apple",
				Place:     2,
				Kills:     1,
			},
		},
		"Game 28": {
			{
				Player:    "Tom",
				Commander: "Kros, Defense Contractor",
				Place:     1,
				Kills:     2,
			},
			{
				Player:    "Mike",
				Commander: "Tovolar, Dire Overlord",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Oona, Queen of the Fae",
				Place:     2,
				Kills:     1,
			},
			{
				Player:    "Peter",
				Commander: "Lord of the Nazgûl",
				Place:     4,
				Kills:     0,
			},
		},
		"Game 29": {
			{
				Player:    "Tom",
				Commander: "Azusa, Lost but Seeking",
				Place:     2,
				Kills:     1,
			},
			{
				Player:    "Mike",
				Commander: "Miirym, Sentinel Wyrm",
				Place:     3,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Oona, Queen of the Fae",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Eriette of the Charmed Apple",
				Place:     1,
				Kills:     1,
			},
		},
		"Game 30": {
			{
				Player:    "Tom",
				Commander: "Ruric Thar, the Unbowed",
				Place:     3,
				Kills:     1,
			},
			{
				Player:    "Mike",
				Commander: "Giada, Font of Hope",
				Place:     1,
				Kills:     2,
			},
			{
				Player:    "Dillon",
				Commander: "Sythis, Harvest's Hand",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Eriette of the Charmed Apple",
				Place:     2,
				Kills:     0,
			},
		},
		"Game 31": {
			{
				Player:    "Tom",
				Commander: "Kazuul, Tyrant of the Cliffs",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Zaxara, the Exemplary",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Satoru Umezawa",
				Place:     1,
				Kills:     3,
			},
			{
				Player:    "Peter",
				Commander: "Sauron, the Dark Lord",
				Place:     4,
				Kills:     0,
			},
		},
		"Game 32": {
			{
				Player:    "Tom",
				Commander: "Kazuul, Tyrant of the Cliffs",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Kathril, Aspect Warper",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Oona, Queen of the Fae",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Sauron, the Dark Lord",
				Place:     1,
				Kills:     3,
			},
		},
		"Game 33": {
			{
				Player:    "Tom",
				Commander: "Vadrik, Astral Archmage",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Giada, Font of Hope",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Satoru Umezawa",
				Place:     3,
				Kills:     1,
			},
			{
				Player:    "Peter",
				Commander: "Lord of the Nazgûl",
				Place:     1,
				Kills:     2,
			},
		},
		"Game 34": {
			{
				Player:    "Tom",
				Commander: "Gyome, Master Chef",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Giada, Font of Hope",
				Place:     3,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Satoru Umezawa",
				Place:     1,
				Kills:     1,
			},
			{
				Player:    "Peter",
				Commander: "The Ur-Dragon",
				Place:     2,
				Kills:     1,
			},
		},
		"Game 35": {
			{
				Player:    "Tom",
				Commander: "Feather, the Redeemed",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Aragorn, King of Gondor",
				Place:     3,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Oona, Queen of the Fae",
				Place:     2,
				Kills:     1,
			},
			{
				Player:    "Peter",
				Commander: "Tazri, Beacon of Unity",
				Place:     1,
				Kills:     1,
			},
		},
		"Game 36": {
			{
				Player:    "Tom",
				Commander: "Lurrus of the Dream-Den",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Kathril, Aspect Warper",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Sakashima of a Thousand Faces",
				Place:     1,
				Kills:     3,
			},
			{
				Player:    "Peter",
				Commander: `Henzie 'Toolbox' Torre`,
				Place:     2,
				Kills:     0,
			},
		},
		"Game 37": {
			{
				Player:    "Tom",
				Commander: "Feather, the Redeemed",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Rakdos, Lord of Riots",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Arcades, the Stategist",
				Place:     1,
				Kills:     2,
			},
			{
				Player:    "Peter",
				Commander: "Umbris, Fear Manifest",
				Place:     3,
				Kills:     1,
			},
		},
		"Game 38": {
			{
				Player:    "Tom",
				Commander: "Azusa, Lost but Seeking",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Aragorn, King of Gondor",
				Place:     1,
				Kills:     3,
			},
			{
				Player:    "Dillon",
				Commander: "Oona, Queen of the Fae",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Thalisse, Reverent Medium",
				Place:     2,
				Kills:     0,
			},
		},
		"Game 39": {
			{
				Player:    "Tom",
				Commander: "Ruric Thar, the Unbowed",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Admiral Beckett Brass",
				Place:     1,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Satoru Umezawa",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Ruxa, Patient Professor",
				Place:     2,
				Kills:     2,
			},
		},
		"Game 40": {
			{
				Player:    "Tom",
				Commander: "Feather, the Redeemed",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Tovolar, Dire Overlord",
				Place:     1,
				Kills:     2,
			},
			{
				Player:    "Dillon",
				Commander: "Magar of the Magic Strings",
				Place:     2,
				Kills:     1,
			},
			{
				Player:    "Peter",
				Commander: "Wyll, Blade of Frontiers",
				Place:     4,
				Kills:     0,
			},
		},
		"Game 41": {
			{
				Player:    "Tom",
				Commander: "Kazuul, Tyrant of the Cliffs",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Thantis, the Warweaver",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Satoru Umezawa",
				Place:     3,
				Kills:     1,
			},
			{
				Player:    "Peter",
				Commander: "Umbris, Fear Manifest",
				Place:     1,
				Kills:     2,
			},
		},
		"Game 42": {
			{
				Player:    "Tom",
				Commander: "Kros, Defense Contractor",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Thantis, the Warweaver",
				Place:     2,
				Kills:     0,
			},
			{
				Player:    "Dillon",
				Commander: "Magar of the Magic Strings",
				Place:     1,
				Kills:     3,
			},
			{
				Player:    "Peter",
				Commander: "Kumena, Tyrant of Orazca",
				Place:     2,
				Kills:     0,
			},
		},
		"Game 43": {
			{
				Player:    "Tom",
				Commander: "Azusa, Lost but Seeking",
				Place:     3,
				Kills:     0,
			},
			{
				Player:    "Mike",
				Commander: "Raffine, Scheming Seer",
				Place:     2,
				Kills:     1,
			},
			{
				Player:    "Dillon",
				Commander: "Sakashima of a Thousand Faces",
				Place:     4,
				Kills:     0,
			},
			{
				Player:    "Peter",
				Commander: "Sauron, the Dark Lord",
				Place:     1,
				Kills:     2,
			},
		},
	}
)

type Migration8 struct{}

func (m *Migration8) Upgrade(ctx context.Context, client *lib.DBClient) error {
	for _, gameSeed := range gameSeeds {
		if _, err := client.Db.ExecContext(ctx, gameSeed); err != nil {
			return lib.NewDBError(gameSeed, err)
		}
	}

	for game, resultInfos := range gameResultSeeds {
		for _, result := range resultInfos {
			if _, err := client.Db.ExecContext(
				ctx,
				insertResult,
				game,
				result.Commander,
				result.Player,
				result.Place,
				result.Kills,
			); err != nil {
				return lib.NewDBError(insertResult, err)
			}
		}
	}
	return nil
}

func (m *Migration8) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, dropGames); err != nil {
		return lib.NewDBError(dropGames, err)
	}
	return nil
}
