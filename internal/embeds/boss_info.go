package embeds

// BossInfo contains description and wiki link for a boss.
type BossInfo struct {
	Description string
	WikiURL     string
}

// bossInformation maps boss names to their info.
var bossInformation = map[string]BossInfo{
	// Wildy Bosses
	"callisto": {
		Description: "A powerful bear found in the Wilderness, Callisto requires both combat prowess and awareness of PKers.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Callisto/Strategies",
	},
	"venenatis": {
		Description: "A giant spider lurking in the Wilderness east of the Bone Yard, known for dropping the Treasonous ring.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Venenatis/Strategies",
	},
	"vetion": {
		Description: "An undead skeletal champion found in the Wilderness, accompanied by his two skeletal hellhounds.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Vet%27ion/Strategies",
	},
	"scorpia": {
		Description: "A scorpion boss found in the Scorpion Pit, one of the easier Wilderness bosses.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Scorpia/Strategies",
	},
	"chaos_elemental": {
		Description: "A chaotic entity found west of the Rogue's Castle, known for unequipping player items.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Chaos_Elemental/Strategies",
	},
	"chaos_fanatic": {
		Description: "A mad mage found in the Wilderness, wielding powerful magic attacks.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Chaos_Fanatic/Strategies",
	},
	"crazy_archaeologist": {
		Description: "An eccentric archaeologist found near the Rogue's Castle, attacking with ranged and special attacks.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Crazy_Archaeologist/Strategies",
	},
	"king_black_dragon": {
		Description: "The original dragon boss, found in his lair beneath the Wilderness. A classic challenge for mid-level players.",
		WikiURL:     "https://oldschool.runescape.wiki/w/King_Black_Dragon/Strategies",
	},
	"artio": {
		Description: "A bear boss found in the Hunter's End, part of the Wilderness Boss Rework.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Artio/Strategies",
	},
	"calvarion": {
		Description: "The reanimated form of Vet'ion, part of the Wilderness Boss Rework.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Calvar%27ion/Strategies",
	},
	"spindle": {
		Description: "The awakened form of Venenatis, part of the Wilderness Boss Rework.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Spindel/Strategies",
	},

	// Group Bosses
	"chambers_of_xeric": {
		Description: "The first raid in OSRS, featuring randomized rooms and the final boss Great Olm. Rewards include the Twisted bow.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Chambers_of_Xeric/Strategies",
	},
	"chambers_of_xeric_challenge_mode": {
		Description: "A harder version of CoX with increased difficulty and exclusive rewards like the Twisted ancestral colour kit.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Chambers_of_Xeric:_Challenge_Mode/Strategies",
	},
	"theatre_of_blood": {
		Description: "The second raid, featuring five bosses and the final encounter with Verzik Vitur. Rewards include the Scythe of vitur.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Theatre_of_Blood/Strategies",
	},
	"theatre_of_blood_hard_mode": {
		Description: "An extremely challenging version of ToB with mechanics changes and exclusive rewards.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Theatre_of_Blood:_Hard_Mode/Strategies",
	},
	"tombs_of_amascut": {
		Description: "The third raid set in the Kharidian Desert, featuring the Wardens and paths with varying difficulty.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Tombs_of_Amascut/Strategies",
	},
	"tombs_of_amascut_expert_mode": {
		Description: "Expert mode ToA with invocation level 300+, significantly harder with better rewards.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Tombs_of_Amascut/Strategies#Expert_Mode",
	},
	"nex": {
		Description: "A powerful Zarosian general requiring a team and strong gear. Drops the Zaryte crossbow and Torva armour.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Nex/Strategies",
	},
	"corporeal_beast": {
		Description: "A massive spirit beast from the Spirit Realm, known for the Elysian spirit shield and requiring a team.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Corporeal_Beast/Strategies",
	},
	"nightmare": {
		Description: "A dark manifestation found beneath Slepe, fought in groups. Drops the Inquisitor's set and nightmare staff.",
		WikiURL:     "https://oldschool.runescape.wiki/w/The_Nightmare/Strategies",
	},

	// Quest Bosses
	"barrows_chests": {
		Description: "Six undead brothers guarding their tomb. A popular mid-game money maker with iconic armour rewards.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Barrows/Strategies",
	},
	"bryophyta": {
		Description: "A moss giant boss found in the Varrock Sewers, accessible with mossy keys.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Bryophyta/Strategies",
	},
	"obor": {
		Description: "A hill giant boss found in Edgeville Dungeon, accessible with giant keys.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Obor/Strategies",
	},
	"deranged_archaeologist": {
		Description: "A crazed archaeologist found on Fossil Island, similar to Crazy Archaeologist.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Deranged_Archaeologist/Strategies",
	},
	"hespori": {
		Description: "A demi-boss plant grown in the Farming Guild's Hespori patch, fought solo.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Hespori/Strategies",
	},
	"mimic": {
		Description: "A shapeshifting boss that mimics treasure chests, accessed via master clue scrolls.",
		WikiURL:     "https://oldschool.runescape.wiki/w/The_Mimic/Strategies",
	},
	"sarachnis": {
		Description: "A spider boss found in the Forthos Dungeon, dropping the Cudgel and requiring moderate combat stats.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Sarachnis/Strategies",
	},
	"skotizo": {
		Description: "A demon boss summoned in the Catacombs of Kourend using dark totems.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Skotizo/Strategies",
	},
	"gauntlet": {
		Description: "A solo minigame where you gather resources and fight the Crystalline Hunllef.",
		WikiURL:     "https://oldschool.runescape.wiki/w/The_Gauntlet/Strategies",
	},
	"corrupted_gauntlet": {
		Description: "A much harder version of the Gauntlet with better rewards including the Blade of saeldor.",
		WikiURL:     "https://oldschool.runescape.wiki/w/The_Gauntlet/Strategies#Corrupted_Gauntlet",
	},

	// Slayer Bosses
	"abyssal_sire": {
		Description: "A high-level Slayer boss requiring 85 Slayer, found in the Abyssal Nexus. Drops the Abyssal bludgeon pieces.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Abyssal_Sire/Strategies",
	},
	"alchemical_hydra": {
		Description: "A powerful Slayer boss requiring 95 Slayer, fought in the Karuulm Slayer Dungeon. Drops the Dragon hunter lance.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Alchemical_Hydra/Strategies",
	},
	"cerberus": {
		Description: "A three-headed hellhound requiring 91 Slayer, guarding the Taverley Dungeon. Drops the Primordial, Pegasian, and Eternal crystals.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Cerberus/Strategies",
	},
	"grotesque_guardians": {
		Description: "A gargoyle boss duo requiring 75 Slayer, fought atop the Slayer Tower. Drops the Granite hammer and Black tourmaline core.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Grotesque_Guardians/Strategies",
	},
	"kraken": {
		Description: "A water-based boss requiring 87 Slayer, found in the Kraken Cove. A relaxed AFK boss dropping the Trident of the seas.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Kraken/Strategies",
	},
	"thermonuclear_smoke_devil": {
		Description: "A smoke devil boss requiring 93 Slayer, found in the Smoke Devil Dungeon. Drops the Occult necklace.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Thermonuclear_smoke_devil/Strategies",
	},

	// World Bosses
	"commander_zilyana": {
		Description: "The Saradominist general in the God Wars Dungeon, requiring 70 Agility. Drops the Saradomin sword and Armadyl crossbow.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Commander_Zilyana/Strategies",
	},
	"dagannoth_prime": {
		Description: "One of the three Dagannoth Kings, weak to magic. Drops the Seercull and Archer ring.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Dagannoth_Prime/Strategies",
	},
	"dagannoth_rex": {
		Description: "One of the three Dagannoth Kings, weak to melee. Drops the Warrior ring and Dragon axe.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Dagannoth_Rex/Strategies",
	},
	"dagannoth_supreme": {
		Description: "One of the three Dagannoth Kings, weak to ranged. Drops the Archers ring.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Dagannoth_Supreme/Strategies",
	},
	"general_graardor": {
		Description: "The Bandosian general in the God Wars Dungeon, requiring 70 Strength. Drops Bandos armour and the Godsword shards.",
		WikiURL:     "https://oldschool.runescape.wiki/w/General_Graardor/Strategies",
	},
	"giant_mole": {
		Description: "A giant mole found beneath Falador Park, a popular mid-level boss. Drops the Mole skin for Falador diaries.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Giant_Mole/Strategies",
	},
	"kalphite_queen": {
		Description: "A two-phase insect boss in the Kalphite Lair, requiring desert gear. Drops the Dragon chainbody and Kalphite head.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Kalphite_Queen/Strategies",
	},
	"kreearra": {
		Description: "The Armadylean general in the God Wars Dungeon, requiring 70 Ranged. Drops Armadyl armour and Godsword shards.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Kree%27arra/Strategies",
	},
	"kril_tsutsaroth": {
		Description: "The Zamorakian general in the God Wars Dungeon, requiring 70 Hitpoints. Drops Subjugation armour and Godsword shards.",
		WikiURL:     "https://oldschool.runescape.wiki/w/K%27ril_Tsutsaroth/Strategies",
	},
	"phosani_nightmare": {
		Description: "A solo-only, significantly harder version of The Nightmare. Drops unique items like the Harmonised orb.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Phosani%27s_Nightmare/Strategies",
	},
	"vorkath": {
		Description: "An undead dragon fought after Dragon Slayer II, one of the best money makers in the game. Drops the Skeletal visage.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Vorkath/Strategies",
	},
	"zulrah": {
		Description: "A serpent boss on Zul-Andra, cycling through three forms. Drops the Toxic blowpipe and Magic fang.",
		WikiURL:     "https://oldschool.runescape.wiki/w/Zulrah/Strategies",
	},
}

// GetBossInfo returns the boss info for a given boss name.
func GetBossInfo(bossName string) (BossInfo, bool) {
	info, exists := bossInformation[bossName]
	return info, exists
}
