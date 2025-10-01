package commands

import "github.com/bwmarrin/discordgo"

// Boss and skill choices for BOTW/SOTW events
// These map to Wise Old Man API metric names

// WildyBoss represents wilderness bosses
type WildyBoss string

const (
	WildyKingBlackDragon    WildyBoss = "king_black_dragon"
	WildyScorpia            WildyBoss = "scorpia"
	WildyArtio              WildyBoss = "artio"
	WildyCallisto           WildyBoss = "callisto"
	WildyCalvarion          WildyBoss = "calvarion"
	WildyChaosElemental     WildyBoss = "chaos_elemental"
	WildyChaosFanatic       WildyBoss = "chaos_fanatic"
	WildyCrazyArchaeologist WildyBoss = "crazy_archaeologist"
	WildySpindel            WildyBoss = "spindel"
	WildyVenenatis          WildyBoss = "venenatis"
	WildyVetion             WildyBoss = "vetion"
)

// GroupBoss represents group bosses
type GroupBoss string

const (
	GroupCorporealBeast    GroupBoss = "corporeal_beast"
	GroupNex               GroupBoss = "nex"
	GroupNightmare         GroupBoss = "nightmare"
	GroupCommanderZilyana  GroupBoss = "commander_zilyana"
	GroupKrilTsutsaroth    GroupBoss = "kril_tsutsaroth"
	GroupGeneralGraardor   GroupBoss = "general_graardor"
	GroupKreeArra          GroupBoss = "kreearra"
)

// QuestBoss represents quest/solo bosses
type QuestBoss string

const (
	QuestDukeSucellus      QuestBoss = "duke_sucellus"
	QuestLeviathan         QuestBoss = "the_leviathan"
	QuestWhisperer         QuestBoss = "the_whisperer"
	QuestVardorvis         QuestBoss = "vardorvis"
	QuestPhantomMuspah     QuestBoss = "phantom_muspah"
	QuestGauntlet          QuestBoss = "the_gauntlet"
	QuestCorruptedGauntlet QuestBoss = "the_corrupted_gauntlet"
	QuestVorkath           QuestBoss = "vorkath"
	QuestZalcano           QuestBoss = "zalcano"
)

// SlayerBoss represents slayer bosses
type SlayerBoss string

const (
	SlayerGrotesqueGuardians      SlayerBoss = "grotesque_guardians"
	SlayerAbyssalSire             SlayerBoss = "abyssal_sire"
	SlayerAlchemicalHydra         SlayerBoss = "alchemical_hydra"
	SlayerThermonuclearSmokeDevil SlayerBoss = "thermonuclear_smoke_devil"
	SlayerKraken                  SlayerBoss = "kraken"
	SlayerCerberus                SlayerBoss = "cerberus"
)

// WorldBoss represents world bosses
type WorldBoss string

const (
	WorldBarrows              WorldBoss = "barrows_chests"
	WorldGiantMole            WorldBoss = "giant_mole"
	WorldDerangedArchaeologist WorldBoss = "deranged_archaeologist"
	WorldDagannothPrime       WorldBoss = "dagannoth_prime"
	WorldDagannothRex         WorldBoss = "dagannoth_rex"
	WorldDagannothSupreme     WorldBoss = "dagannoth_supreme"
	WorldSarachnis            WorldBoss = "sarachnis"
	WorldKalphiteQueen        WorldBoss = "kalphite_queen"
	WorldSkotizo              WorldBoss = "skotizo"
)

// Skill represents non-combat skills for SOTW
type Skill string

const (
	SkillPrayer       Skill = "prayer"
	SkillCooking      Skill = "cooking"
	SkillWoodcutting  Skill = "woodcutting"
	SkillFletching    Skill = "fletching"
	SkillFishing      Skill = "fishing"
	SkillFiremaking   Skill = "firemaking"
	SkillCrafting     Skill = "crafting"
	SkillSmithing     Skill = "smithing"
	SkillMining       Skill = "mining"
	SkillHerblore     Skill = "herblore"
	SkillAgility      Skill = "agility"
	SkillThieving     Skill = "thieving"
	SkillSlayer       Skill = "slayer"
	SkillFarming      Skill = "farming"
	SkillRunecraft    Skill = "runecraft"
	SkillHunter       Skill = "hunter"
	SkillConstruction Skill = "construction"
)

// WildyBossChoices returns Discord choices for wilderness bosses
func WildyBossChoices() []*discordgo.ApplicationCommandOptionChoice {
	return []*discordgo.ApplicationCommandOptionChoice{
		{Name: "King Black Dragon", Value: string(WildyKingBlackDragon)},
		{Name: "Scorpia", Value: string(WildyScorpia)},
		{Name: "Artio", Value: string(WildyArtio)},
		{Name: "Callisto", Value: string(WildyCallisto)},
		{Name: "Calvarion", Value: string(WildyCalvarion)},
		{Name: "Chaos Elemental", Value: string(WildyChaosElemental)},
		{Name: "Chaos Fanatic", Value: string(WildyChaosFanatic)},
		{Name: "Crazy Archaeologist", Value: string(WildyCrazyArchaeologist)},
		{Name: "Spindel", Value: string(WildySpindel)},
		{Name: "Venenatis", Value: string(WildyVenenatis)},
		{Name: "Vet'ion", Value: string(WildyVetion)},
	}
}

// GroupBossChoices returns Discord choices for group bosses
func GroupBossChoices() []*discordgo.ApplicationCommandOptionChoice {
	return []*discordgo.ApplicationCommandOptionChoice{
		{Name: "Corporeal Beast", Value: string(GroupCorporealBeast)},
		{Name: "Nex", Value: string(GroupNex)},
		{Name: "Nightmare", Value: string(GroupNightmare)},
		{Name: "Commander Zilyana (Saradomin)", Value: string(GroupCommanderZilyana)},
		{Name: "K'ril Tsutsaroth (Zamorak)", Value: string(GroupKrilTsutsaroth)},
		{Name: "General Graardor (Bandos)", Value: string(GroupGeneralGraardor)},
		{Name: "Kree'arra (Armadyl)", Value: string(GroupKreeArra)},
	}
}

// QuestBossChoices returns Discord choices for quest bosses
func QuestBossChoices() []*discordgo.ApplicationCommandOptionChoice {
	return []*discordgo.ApplicationCommandOptionChoice{
		{Name: "Duke Sucellus", Value: string(QuestDukeSucellus)},
		{Name: "The Leviathan", Value: string(QuestLeviathan)},
		{Name: "The Whisperer", Value: string(QuestWhisperer)},
		{Name: "Vardorvis", Value: string(QuestVardorvis)},
		{Name: "Phantom Muspah", Value: string(QuestPhantomMuspah)},
		{Name: "The Gauntlet", Value: string(QuestGauntlet)},
		{Name: "The Corrupted Gauntlet", Value: string(QuestCorruptedGauntlet)},
		{Name: "Vorkath", Value: string(QuestVorkath)},
		{Name: "Zalcano", Value: string(QuestZalcano)},
	}
}

// SlayerBossChoices returns Discord choices for slayer bosses
func SlayerBossChoices() []*discordgo.ApplicationCommandOptionChoice {
	return []*discordgo.ApplicationCommandOptionChoice{
		{Name: "Grotesque Guardians", Value: string(SlayerGrotesqueGuardians)},
		{Name: "Abyssal Sire", Value: string(SlayerAbyssalSire)},
		{Name: "Alchemical Hydra", Value: string(SlayerAlchemicalHydra)},
		{Name: "Thermonuclear Smoke Devil", Value: string(SlayerThermonuclearSmokeDevil)},
		{Name: "Kraken", Value: string(SlayerKraken)},
		{Name: "Cerberus", Value: string(SlayerCerberus)},
	}
}

// WorldBossChoices returns Discord choices for world bosses
func WorldBossChoices() []*discordgo.ApplicationCommandOptionChoice {
	return []*discordgo.ApplicationCommandOptionChoice{
		{Name: "Barrows Chests", Value: string(WorldBarrows)},
		{Name: "Giant Mole", Value: string(WorldGiantMole)},
		{Name: "Deranged Archaeologist", Value: string(WorldDerangedArchaeologist)},
		{Name: "Dagannoth Prime", Value: string(WorldDagannothPrime)},
		{Name: "Dagannoth Rex", Value: string(WorldDagannothRex)},
		{Name: "Dagannoth Supreme", Value: string(WorldDagannothSupreme)},
		{Name: "Sarachnis", Value: string(WorldSarachnis)},
		{Name: "Kalphite Queen", Value: string(WorldKalphiteQueen)},
		{Name: "Skotizo", Value: string(WorldSkotizo)},
	}
}

// SkillChoices returns Discord choices for skills
func SkillChoices() []*discordgo.ApplicationCommandOptionChoice {
	return []*discordgo.ApplicationCommandOptionChoice{
		{Name: "Prayer", Value: string(SkillPrayer)},
		{Name: "Cooking", Value: string(SkillCooking)},
		{Name: "Woodcutting", Value: string(SkillWoodcutting)},
		{Name: "Fletching", Value: string(SkillFletching)},
		{Name: "Fishing", Value: string(SkillFishing)},
		{Name: "Firemaking", Value: string(SkillFiremaking)},
		{Name: "Crafting", Value: string(SkillCrafting)},
		{Name: "Smithing", Value: string(SkillSmithing)},
		{Name: "Mining", Value: string(SkillMining)},
		{Name: "Herblore", Value: string(SkillHerblore)},
		{Name: "Agility", Value: string(SkillAgility)},
		{Name: "Thieving", Value: string(SkillThieving)},
		{Name: "Slayer", Value: string(SkillSlayer)},
		{Name: "Farming", Value: string(SkillFarming)},
		{Name: "Runecraft", Value: string(SkillRunecraft)},
		{Name: "Hunter", Value: string(SkillHunter)},
		{Name: "Construction", Value: string(SkillConstruction)},
	}
}

// FormatActivityName converts snake_case to Title Case for display
func FormatActivityName(activity string) string {
	// Simple conversion for display
	result := ""
	capitalize := true
	for _, c := range activity {
		if c == '_' {
			result += " "
			capitalize = true
		} else if capitalize {
			result += string(c - 32) // Convert to uppercase
			capitalize = false
		} else {
			result += string(c)
		}
	}
	return result
}
