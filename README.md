# voidling

A Discord bot for managing RuneScape (OSRS) clan events, built in Go. Complete rewrite of [TopezEventBot](https://github.com/kaffeed/TopezEventBot) with improved performance, better maintainability, and modern patterns.

## Features

### ✅ Implemented

- **Account Linking** (`/link-rsn`, `/unlink-rsn`)
  - Link Discord accounts to RuneScape usernames
  - Player verification with Wise Old Man API
  - Interactive confirmation flow with player stats embed

- **Boss of the Week** (`/botw`)
  - Weekly boss kill count competitions across 5 categories:
    - Wilderness bosses (Callisto, Vet'ion, Venenatis, etc.)
    - Group bosses (CoX, ToB, ToA, Corp, Nex, etc.)
    - Quest bosses (Galvek, Vanstrom, Glough, etc.)
    - Slayer bosses (Cerberus, Hydra, Sire, etc.)
    - World bosses (Phantom Muspah, DT2 bosses, etc.)
  - Automatic tracking via Wise Old Man competitions
  - Thread-based participation with buttons
  - Winner announcements with medals (🥇🥈🥉)

- **Skill of the Week** (`/sotw`)
  - Weekly skill experience competitions (all 23 OSRS skills)
  - Automatic XP gain tracking via Wise Old Man
  - Thread-based participation

- **Mass Events** (`/mass`)
  - Schedule clan mass events with boss dropdown
  - OSRS Wiki images for activities
  - Discord timestamp formatting with timezone support
  - User and server-specific timezone preferences

- **Server Configuration** (`/config`)
  - Set coordinator role for event management
  - Configure competition code notification channel
  - Set default server timezone
  - Personal timezone preferences

### 📋 Planned

- Scheduled notifications (30min before events)
- Progress tracking background job
- Leaderboard history tracking
- Admin/warning system
- Comprehensive test coverage

## Tech Stack

- **Language**: Go 1.24.1
- **Discord Library**: [discordgo](https://github.com/bwmarrin/discordgo) v0.29.0
- **Database**: SQLite with [sqlc](https://sqlc.dev/) (type-safe queries) and [goose](https://github.com/pressly/goose) (migrations)
- **External API**: [Wise Old Man API](https://docs.wiseoldman.net/) for player tracking
- **Scheduler**: [robfig/cron](https://github.com/robfig/cron) v3 for background tasks (planned)

## Project Structure

```
voidling/
├── cmd/voidling/         # Application entry point with migration runner
├── internal/
│   ├── bot/              # Discord bot core with interaction routing
│   ├── commands/         # Command handlers (register, trackable, schedulable, config)
│   ├── embeds/           # Discord embed builders for all event types
│   ├── database/         # Generated sqlc code (type-safe queries)
│   ├── models/           # Domain models (events, players, hiscores)
│   ├── wiseoldman/       # Wise Old Man API client
│   └── timezone/         # Timezone utilities and autocomplete
├── migrations/           # Database migrations (goose) - 6 migrations
├── queries/              # SQL query definitions for sqlc
├── config/               # Configuration management (.env support)
├── refactor/             # Development documentation and session notes
└── Makefile              # Build automation (25+ targets)
```

## Quick Start

### Prerequisites

- Go 1.24.1 or higher
- SQLite3
- Discord bot token ([Create one here](https://discord.com/developers/applications))

### Installation

```bash
# Clone repository
git clone https://github.com/kaffeed/voidling.git
cd voidling

# Install dev tools and dependencies
make init

# Configure environment
cp .env.example .env
# Edit .env and add your DISCORD_TOKEN

# Build and run
make run
```

The database migrations run automatically on startup. Your bot is now ready!

### Environment Configuration

Create a `.env` file (or set environment variables):

```bash
# Required
DISCORD_TOKEN=your_bot_token_here

# Optional
DATABASE_PATH=~/.voidling/voidling.db     # Default database location
LOG_LEVEL=info                             # debug|info|warn|error
DISCORD_GUILD_ID=123456789                 # For fast command registration during dev
```

## Development

### Useful Make Targets

```bash
make build          # Build for current OS
make run            # Build and run
make test           # Run tests
make fmt            # Format code
make vet            # Run go vet
make check          # Run fmt + vet + test

make sqlc-generate  # Regenerate database queries
make migrate-up     # Run migrations
make migrate-down   # Rollback migration
make migrate-status # Show migration status

make build-all      # Build for Windows, Linux, macOS
make release        # Create optimized release builds
make help           # Show all available targets
```

### Database

SQLite database with automatic migrations on startup. Schema managed via goose, queries via sqlc.

**Schema includes:**
- Account links (Discord ↔ RuneScape)
- Trackable events (BOTW/SOTW competitions)
- Schedulable events (Mass events)
- Guild configuration (roles, channels, timezones)
- User timezone preferences

**Regenerate queries after modifying `queries/*.sql`:**
```bash
make sqlc-generate
```

**Create new migration:**
```bash
make migrate-create NAME=add_new_feature
```

## Architecture

### Key Components

**Bot Core** (`internal/bot/`)
- Discord session management
- Interaction routing (commands, buttons, modals, autocomplete)
- Permission checking (coordinator role, admin role)
- Handler registration

**Commands** (`internal/commands/`)
- `register.go` - Account linking (`/link-rsn`, `/unlink-rsn`)
- `trackable.go` - Base logic for BOTW/SOTW events
- `botw.go` - Boss of the Week command handlers
- `sotw.go` - Skill of the Week command handlers
- `schedulable.go` - Mass event scheduling
- `config.go` - Server configuration commands
- `choices.go` - Boss and skill dropdown data

**Embeds** (`internal/embeds/`)
- PlayerInfo - Player stats with WOM data
- BossOfTheWeek / SkillOfTheWeek - Event announcements
- EventWinners - Winner displays with medals
- MassEvent - Mass event scheduling with timestamps
- Error/Success - Consistent messaging

**Wise Old Man Client** (`internal/wiseoldman/`)
- HTTP client for WOM API
- Player data fetching
- Competition creation and management
- Participant tracking

**Database Layer** (`internal/database/`)
- Type-safe queries generated by sqlc
- Transaction support
- Proper error handling

### Design Patterns

**Custom ID Routing**: `"action:data"` format for buttons/modals
```go
confirm-rsn:username
register-for-botw:womCompetitionID,threadID
```

**Modal Flow**: Command → Modal → Processing → Followup
```go
/link-rsn → show modal → process submission → show player embed → button confirmation
```

**Transaction Pattern**: Proper rollback/commit
```go
tx, _ := db.Begin()
defer tx.Rollback()
qtx := queries.WithTx(tx)
// ... operations ...
tx.Commit()
```

## Commands Reference

### User Commands
- `/link-rsn` - Link your RuneScape account
- `/unlink-rsn` - Unlink your account
- `/config set-my-timezone` - Set your timezone preference

### Coordinator Commands (requires Coordinator role)
- `/botw wildy|group|quest|slayer|world` - Start BOTW competition
- `/botw finish` - Finish current BOTW and announce winners
- `/sotw start` - Start SOTW competition
- `/sotw finish` - Finish current SOTW and announce winners
- `/mass` - Schedule a mass event

### Admin Commands (requires Administrator permission)
- `/config set-coordinator-role` - Set coordinator role
- `/config set-competition-code-channel` - Set WOM code channel
- `/config set-default-timezone` - Set server default timezone
- `/config show` - Show current configuration

## Migration from TopezEventBot

Complete rewrite of [TopezEventBot](https://github.com/kaffeed/TopezEventBot) (C#/.NET) in Go.

**Key Improvements:**
- 🚀 Better performance and lower resource usage
- 📦 Single binary deployment (no runtime dependencies)
- 🔒 Type-safe database queries (sqlc vs EF Core)
- 🎯 Cleaner migration management (goose vs EF)
- 🌐 Modern API integration (Wise Old Man vs OSRS Hiscore)
- ⚡ Faster command registration and interaction handling

**Data Migration:**
Database schema is similar but not directly compatible. Manual migration required.

## Contributing

Issues and pull requests welcome! This is an active project.

**Development priorities:**
1. Scheduled notifications implementation
2. Progress tracking background job
3. Test coverage
4. CI/CD setup

## License

MIT License

## Links

- **Original Project**: [TopezEventBot (C#)](https://github.com/kaffeed/TopezEventBot)
- **Wise Old Man**: [API Documentation](https://docs.wiseoldman.net/)
- **OSRS Wiki**: [Boss Information](https://oldschool.runescape.wiki/)

---

**Status**: Active development • Production ready for core features • ~5,300 lines of Go code
