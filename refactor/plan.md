# Refactor Plan - TopezEventBot C# → voidling Go

**Session Started:** 2025-10-01
**Type:** Complete Language Rewrite (C# → Go)
**Original Codebase:** TopezEventBot (C#/.NET Discord.Net)
**Target Codebase:** voidling (Go/discordgo)

---

## Initial State Analysis

### Current Architecture (C# TopezEventBot)

**Tech Stack:**
- **Framework:** .NET 8.0 Worker Service
- **Discord Library:** Discord.Net (v3.12.0)
- **Database:** SQLite with Entity Framework Core
- **Task Scheduling:** Coravel (v5.0.2)
- **HTTP Client:** Built-in HttpClient for RuneScape Hiscore API

**Core Components:**

1. **Database Layer (TopezEventBot.Data)**
   - `TopezContext` - EF Core DbContext
   - Entities:
     - `AccountLink` - Links Discord users to RuneScape names
     - `TrackableEvent` - Events that track progress (Boss/Skill of the Week)
     - `TrackableEventParticipation` - User participation in trackable events
     - `SchedulableEvent` - Scheduled events (Mass events, Wildy Wednesday)
     - `SchedulableEventParticipation` - User registration for scheduled events
     - `GuildWarningChannel` - Guild warning system configuration
     - `Warning` - User warnings
   - Many-to-many relationships via join entities

2. **Discord Modules (Command Handlers)**
   - `BossOfTheWeekModule` - Boss of the Week competitions (5 boss categories: wildy, group, quest, slayer, world)
   - `SkillOfTheWeekModule` - Skill of the Week competitions
   - `MassEventModule` - Schedule mass events
   - `WildernessEventModule` - Wildy Wednesday events
   - `RegisterRunescapeNameModule` - Link/unlink RuneScape accounts
   - `WarningModule` - User warning system
   - `AdminModule` - Admin commands
   - Base classes: `TrackableEventModuleBase`, `SchedulableEventModuleBase`

3. **Scheduled Tasks (Invocables)**
   - `CheckForScheduledEventNotification` - Notifies users 30 minutes before scheduled events
   - `FetchEventProgressInvocable` - Periodically fetches player progress from RuneScape API
   - `RemindPeopleToLinkUsernameInvocable` - Reminds users to link accounts
   - Runs via Coravel scheduler (every minute for notifications)

4. **External API Integration**
   - `IRunescapeHiscoreHttpClient` - Fetches player stats from OSRS Hiscores
   - Parses CSV-like format from RuneScape API
   - Supports boss KC and skill XP tracking

5. **Features**
   - Slash commands with role-based permissions
   - Interactive components (buttons, modals)
   - Thread creation for event discussions
   - Ephemeral messages for privacy
   - Event registration with validation
   - Leaderboards and winner announcements
   - Progress tracking with starting/ending points
   - Scheduled event reminders via DM

### Target Architecture (Go voidling)

**Tech Stack:**
- **Language:** Go 1.24.1
- **Discord Library:** discordgo (bwmarrin/discordgo)
- **Database:** SQLite with sqlc for type-safe queries
- **Migrations:** goose for schema management
- **Task Scheduling:** Custom scheduler or cron-like library (robfig/cron/v3)
- **HTTP Client:** net/http for RuneScape API
- **Configuration:** context7 pattern for structured context management

---

## Problem Areas

1. **Database Layer Migration**
   - EF Core entities → Go structs
   - LINQ queries → sqlc generated code
   - Migrations need to be recreated with goose
   - Many-to-many relationships require explicit join tables

2. **Discord Interaction Model**
   - Discord.Net's interaction framework → discordgo event handlers
   - Slash command registration differences
   - Component interaction handling (buttons, modals)
   - Message builder patterns different between libraries

3. **Scheduled Tasks**
   - Coravel scheduler → Go scheduler (cron or custom)
   - Background task lifecycle management
   - Concurrent processing patterns

4. **HTTP Client & API Integration**
   - C# HttpClient → Go net/http
   - CSV parsing for RuneScape API
   - Player data models

5. **Error Handling & Logging**
   - ILogger → structured logging in Go (logrus, zap, or slog)
   - Different error handling idioms

---

## Dependencies

**External:**
- discordgo for Discord bot functionality
- sqlc for database query generation
- goose for migrations
- robfig/cron for scheduling (or alternative)
- Standard library for HTTP, JSON, CSV parsing

**Internal:**
- Database context and connection management
- Discord session management
- Command router/handler system
- Permission/role checking middleware

---

## Test Coverage

**Original (C#):**
- No visible test files in the scanned codebase
- Will need to create tests from scratch in Go

**Target (Go):**
- Unit tests for database operations
- Integration tests for Discord commands
- HTTP client mocks for RuneScape API
- Scheduler tests

---

## Refactoring Tasks

### Phase 1: Project Setup & Infrastructure
**Priority:** Critical | **Risk:** Low

- [x] Initialize Go module (already done: go.mod exists)
- [ ] Set up project structure following Go conventions
  ```
  voidling/
  ├── cmd/voidling/         # Main application
  ├── internal/
  │   ├── bot/              # Discord bot logic
  │   ├── commands/         # Command handlers
  │   ├── scheduler/        # Background tasks
  │   ├── database/         # DB queries (sqlc)
  │   ├── models/           # Domain models
  │   └── runescape/        # OSRS API client
  ├── migrations/           # goose migrations
  ├── queries/              # sqlc query definitions
  └── config/               # Configuration files
  ```
- [ ] Add required dependencies to go.mod
- [ ] Create configuration system (environment variables, config file)
- [ ] Set up logging infrastructure

### Phase 2: Database Layer
**Priority:** Critical | **Risk:** Medium

- [ ] Define database schema in SQL (for goose)
- [ ] Create initial migration files
  - `001_create_account_links.sql`
  - `002_create_trackable_events.sql`
  - `003_create_schedulable_events.sql`
  - `004_create_warnings.sql`
- [ ] Write sqlc queries for:
  - Account link CRUD operations
  - Event creation and management
  - Participation registration
  - Leaderboard queries
  - Warning system
- [ ] Generate Go code with sqlc
- [ ] Create database connection manager
- [ ] Implement transaction support

### Phase 3: RuneScape API Client
**Priority:** High | **Risk:** Low

- [ ] Create HTTP client for OSRS Hiscore API
- [ ] Implement player data fetching
- [ ] Parse CSV response format
- [ ] Define Go structs for:
  - Player stats
  - Skills (23 skills)
  - Bosses (categorized: wildy, group, quest, slayer, world)
- [ ] Add rate limiting/retry logic
- [ ] Error handling for player not found

### Phase 4: Discord Bot Core
**Priority:** Critical | **Risk:** Medium

- [ ] Initialize discordgo session
- [ ] Implement bot startup and authentication
- [ ] Register slash commands with Discord API
- [ ] Create command router/dispatcher
- [ ] Implement role-based permission checks
- [ ] Set up interaction handlers:
  - Slash command handler
  - Button interaction handler
  - Modal submit handler
- [ ] Implement message builders for embeds
- [ ] Create thread management utilities

### Phase 5: Command Modules - Registration System
**Priority:** High | **Risk:** Low

- [ ] `/link-rsn` command
  - Modal for username input
  - Fetch player from OSRS API
  - Display player embed for confirmation
  - Confirmation button handler
  - Store account link in database
  - Handle existing linked accounts (deactivate previous)
- [ ] `/unlink-rsn` command
  - Deactivate current linked account
  - Confirmation message

### Phase 6: Command Modules - Trackable Events
**Priority:** High | **Risk:** Medium

- [ ] Boss of the Week (`/botw`) group commands:
  - `/botw wildy [boss] [isActive]`
  - `/botw group [boss] [isActive]`
  - `/botw quest [boss] [isActive]`
  - `/botw slayer [boss] [isActive]`
  - `/botw world [boss] [isActive]`
  - Create event in database
  - Create thread for event
  - Post event embed with register/list buttons
  - Button handler: Register for event
  - Button handler: List participants
- [ ] Skill of the Week (`/sotw`) commands:
  - Similar structure to BOTW
  - 23 skill choices
- [ ] `/finish` command (for both BOTW/SOTW):
  - Mark event as inactive
  - Fetch final progress for all participants
  - Calculate winners (top 3)
  - Post winner announcement with embed
  - Store final results
- [ ] `/leaderboard` command:
  - Query all-time wins
  - Display formatted leaderboard

### Phase 7: Command Modules - Schedulable Events
**Priority:** High | **Risk:** Low

- [ ] Mass Event (`/mass`) commands:
  - `/mass schedule [activity] [location] [time]`
  - Create scheduled event
  - Post event embed with buttons
  - Register button handler
  - List participants button handler
  - Validate time format
- [ ] Wildy Wednesday commands:
  - Similar structure to mass events
  - Specific for wilderness activities

### Phase 8: Command Modules - Admin/Warnings
**Priority:** Medium | **Risk:** Low

- [ ] Warning system commands:
  - Issue warning with modal
  - Store in database
  - Notify user
  - Query warnings for user
- [ ] Admin commands:
  - Configure guild settings
  - Manage warning channels

### Phase 9: Scheduled Tasks
**Priority:** High | **Risk:** Medium

- [ ] Set up cron scheduler
- [ ] Scheduled event notifications task:
  - Run every minute
  - Query events in 30-minute window
  - Send DMs to unnotified participants
  - Mark as notified
- [ ] Progress tracking task:
  - Run periodically (every minute/hour)
  - Fetch current progress for active event participants
  - Store progress snapshots
  - Handle API failures gracefully
- [ ] Reminder task:
  - Remind users to link accounts

### Phase 10: Utilities & Helpers
**Priority:** Medium | **Risk:** Low

- [ ] Embed builders:
  - Event embeds (BOTW, SOTW, Mass, Wildy)
  - Player info embed
  - Winner announcement embed
  - Reminder embed
- [ ] Role checking middleware
- [ ] Discord ID utilities
- [ ] Time/date formatters
- [ ] CSV parser for OSRS API
- [ ] Context7 integration for structured context

### Phase 11: Configuration & Deployment
**Priority:** Medium | **Risk:** Low

- [ ] Environment variable support:
  - `DISCORD_TOKEN`
  - `DATABASE_PATH`
  - `LOG_LEVEL`
- [ ] Database migration runner in main()
- [ ] Graceful shutdown handling
- [ ] Health check endpoint (optional)
- [ ] Docker support (Dockerfile)

### Phase 12: Testing
**Priority:** Low | **Risk:** Low

- [ ] Unit tests for database queries
- [ ] Tests for RuneScape API client
- [ ] Mock Discord interactions for command tests
- [ ] Integration tests for full flows
- [ ] Scheduler tests

### Phase 13: Documentation
**Priority:** Low | **Risk:** Low

- [ ] README.md with setup instructions
- [ ] Architecture documentation
- [ ] Command usage guide
- [ ] Development guide
- [ ] Migration guide from old bot

---

## Validation Checklist

After implementation, verify:

- [ ] All slash commands registered and working
- [ ] Database migrations run successfully
- [ ] Account linking flow complete
- [ ] BOTW registration and finish flow works
- [ ] SOTW registration and finish flow works
- [ ] Mass event scheduling works
- [ ] Wildy Wednesday scheduling works
- [ ] Scheduled notifications sent correctly
- [ ] Progress tracking updates database
- [ ] Leaderboards display correctly
- [ ] Role permissions enforced
- [ ] Ephemeral messages work
- [ ] Threads created properly
- [ ] Embeds render correctly
- [ ] All buttons/modals functional
- [ ] Error handling works
- [ ] Logging is comprehensive
- [ ] Configuration loaded correctly
- [ ] No race conditions in concurrent code
- [ ] Database transactions work
- [ ] OSRS API integration functional

---

## De-Para Mapping (C# → Go)

| C# Component | Go Component | Status |
|-------------|-------------|--------|
| TopezContext (EF Core) | database.Queries (sqlc) | Pending |
| AccountLink entity | AccountLink struct + sqlc queries | Pending |
| TrackableEvent entity | TrackableEvent struct + sqlc queries | Pending |
| SchedulableEvent entity | SchedulableEvent struct + sqlc queries | Pending |
| BossOfTheWeekModule | commands/botw.go | Pending |
| SkillOfTheWeekModule | commands/sotw.go | Pending |
| MassEventModule | commands/mass.go | Pending |
| RegisterRunescapeNameModule | commands/register.go | Pending |
| WarningModule | commands/warnings.go | Pending |
| TrackableEventModuleBase | commands/trackable_base.go | Pending |
| SchedulableEventModuleBase | commands/schedulable_base.go | Pending |
| IRunescapeHiscoreHttpClient | runescape/client.go | Pending |
| CheckForScheduledEventNotification | scheduler/notifications.go | Pending |
| FetchEventProgressInvocable | scheduler/progress.go | Pending |
| Coravel scheduler | robfig/cron or custom | Pending |
| Program.cs (startup) | cmd/voidling/main.go | Pending |
| Discord.Net | discordgo | Pending |
| Entity Framework Core | sqlc + goose | Pending |

---

## Risk Assessment

**High Risk:**
- Discord interaction model differences (need to study discordgo patterns)
- Concurrent task scheduling and database access (need proper locking)
- Many-to-many relationship queries in sqlc (more verbose than LINQ)

**Medium Risk:**
- Database migration data preservation (if migrating existing data)
- RuneScape API parsing edge cases
- Error handling patterns across goroutines

**Low Risk:**
- Project structure setup
- Configuration management
- Logging implementation
- Most CRUD operations

---

## Migration Strategy

**Approach:** Greenfield rewrite (not incremental migration)

1. Build new Go bot from scratch
2. Run side-by-side with C# bot during testing
3. Migrate database schema (export/import if needed)
4. Switch Discord bot token to new bot
5. Decommission C# bot

**Rollback Plan:**
- Keep C# bot deployable
- Database backup before migration
- Can revert to C# bot if critical issues found

---

## Progress Tracking

- **Total Tasks:** 13 phases, ~80 individual tasks
- **Completed:** 1 (Go module initialized)
- **In Progress:** 0
- **Pending:** All phases

---

## Notes

- Use context7 pattern as requested (need to research this pattern)
- Follow Go best practices: error handling, package structure, naming
- Keep goroutine usage safe with proper synchronization
- Use structured logging for better debugging
- Add metrics/observability if needed for production
- Consider using dependency injection pattern for testability
