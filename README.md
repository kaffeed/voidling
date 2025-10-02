# voidling

A Discord bot for managing RuneScape clan events, rewritten in Go from TopezEventBot (C#).

## Features

- **Account Linking**: Link Discord accounts to RuneScape usernames
- **Boss of the Week**: Weekly boss kill count competitions (5 categories: Wildy, Group, Quest, Slayer, World)
- **Skill of the Week**: Weekly skill experience competitions (23 OSRS skills)
- **Mass Events**: Schedule and manage clan mass events
- **Wildy Wednesday**: Wilderness events scheduling
- **Progress Tracking**: Automatic fetching of player progress via OSRS Hiscore API
- **Leaderboards**: Track winners and display all-time leaderboards
- **Scheduled Notifications**: DM reminders 30 minutes before events
- **Warning System**: Guild warning management for moderators

## Tech Stack

- **Language**: Go 1.24.1
- **Discord Library**: discordgo
- **Database**: SQLite with sqlc (type-safe queries) and goose (migrations)
- **Scheduler**: robfig/cron for background tasks
- **External API**: OSRS Hiscore API integration

## Project Structure

```
voidling/
├── cmd/voidling/         # Main application entry point
├── internal/
│   ├── bot/              # Discord bot core logic
│   ├── commands/         # Command handlers (to be implemented)
│   ├── scheduler/        # Background tasks (to be implemented)
│   ├── database/         # Generated sqlc code
│   ├── models/           # Domain models
│   └── runescape/        # OSRS API client
├── migrations/           # Database migrations (goose)
├── queries/              # SQL query definitions for sqlc
├── config/               # Configuration management
└── refactor/             # Refactoring plan and state
```

## Setup

### Prerequisites

- Go 1.24.1 or higher
- SQLite3
- sqlc (for regenerating database code)
- A Discord bot token

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd voidling
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

4. Edit `.env` and add your Discord bot token:
   ```
   DISCORD_TOKEN=your_bot_token_here
   DISCORD_GUILD_ID=your_guild_id_for_testing  # Optional
   ```

5. Build the bot:
   ```bash
   go build -o voidling ./cmd/voidling
   ```

6. Run the bot:
   ```bash
   ./voidling
   ```

## Database

The bot uses SQLite for data persistence. Migrations are run automatically on startup.

### Running Migrations Manually

```bash
goose -dir migrations sqlite3 voidling.db up
```

### Regenerating sqlc Code

If you modify SQL queries in `queries/`:

```bash
sqlc generate
```

## Development Status

### Completed
- ✅ Project structure setup
- ✅ Database schema and migrations
- ✅ sqlc query definitions
- ✅ RuneScape API client
- ✅ Discord bot core with discordgo
- ✅ Basic command infrastructure
- ✅ Configuration system

### In Progress
- 🚧 Command implementations (link-rsn, unlink-rsn, etc.)
- 🚧 Scheduled tasks (notifications, progress tracking)

### Planned
- [ ] Full command implementations:
  - [ ] `/botw` (Boss of the Week)
  - [ ] `/sotw` (Skill of the Week)
  - [ ] `/mass` (Mass events)
  - [ ] `/wildy` (Wildy Wednesday)
  - [ ] Warning system
  - [ ] Admin commands
- [ ] Embed builders
- [ ] Role-based permission checks
- [ ] Thread management
- [ ] Comprehensive testing
- [ ] Documentation

## Configuration

Environment variables:

- `DISCORD_TOKEN` (required): Your Discord bot token
- `DATABASE_PATH` (optional): Path to SQLite database (default: `~/.voidling/voidling.db`)
- `LOG_LEVEL` (optional): Logging level - debug, info, warn, error (default: `info`)
- `DISCORD_GUILD_ID` (optional): Guild ID for command registration during development

## Migration from TopezEventBot

This bot is a complete rewrite of TopezEventBot (C#/.NET) to Go. Key improvements:

- Type-safe database queries with sqlc
- Cleaner migration management with goose
- Better performance and lower resource usage
- Simpler deployment (single binary)
- Context7 pattern for structured context management

### Data Migration

If you have existing data in TopezEventBot, you'll need to export and import:

1. Export data from the old SQLite database
2. Transform schema to match new structure
3. Import into voidling database

(Detailed migration script to be added)

## Contributing

This is a personal project, but contributions are welcome!

## License

MIT License (or specify your license)

## Original Project

Based on TopezEventBot (C#/.NET) - rewritten in Go for better performance and maintainability.
