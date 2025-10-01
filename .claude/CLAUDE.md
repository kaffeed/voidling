# Voidbound Discord Bot - Development Memory

## Project Overview

**Project Name:** voidbound
**Type:** Discord bot rewrite from C# to Go
**Original Project:** TopezEventBot (C#/.NET)
**Primary Language:** Go 1.24.1
**Purpose:** RuneScape clan event management bot

### Tech Stack
- **Discord Library:** discordgo (bwmarrin/discordgo)
- **Database:** SQLite3 with sqlc (type-safe queries) and goose (migrations)
- **Scheduler:** robfig/cron/v3
- **HTTP Client:** Standard library for OSRS Hiscore API
- **Pattern:** context7 (requested by user)

---

## Session History

### Session 1: 2025-10-01 - Initial Project Setup & Foundation

**Duration:** Full session
**Status:** Infrastructure Complete ✅

#### Accomplished

**1. Project Analysis & Planning**
- Analyzed entire TopezEventBot C# codebase
- Identified 7 core modules and 14 command groups
- Documented features: BOTW (5 categories), SOTW (23 skills), Mass events, Wildy Wednesday
- Created comprehensive 13-phase refactoring plan (80+ tasks)
- Established de-para mapping from C# to Go components
- Session state tracking in `refactor/plan.md` and `refactor/state.json`

**2. Project Structure Setup**
```
voidbound/
├── cmd/voidbound/main.go          # Application entry point
├── config/config.go               # Environment-based configuration
├── internal/
│   ├── bot/bot.go                 # Discord bot core with discordgo
│   ├── database/                  # sqlc generated (7 files)
│   ├── models/hiscore.go          # OSRS domain models
│   └── runescape/client.go        # Hiscore API client
├── migrations/
│   └── 00001_initial_schema.sql   # Initial database schema
├── queries/                        # SQL query definitions (4 files)
│   ├── account_links.sql
│   ├── trackable_events.sql
│   ├── schedulable_events.sql
│   └── warnings.sql
├── refactor/                       # Refactoring documentation
│   ├── plan.md                    # 13-phase implementation plan
│   └── state.json                 # Session state tracking
├── .env.example                    # Configuration template
├── .gitignore                      # Git ignore rules
├── Makefile                        # Build automation
├── README.md                       # Project documentation
├── go.mod & go.sum                # Go dependencies
└── sqlc.yaml                      # sqlc configuration
```

**3. Database Layer (Complete)**
- Designed 8-table schema:
  - `account_links` - Discord ↔ RuneScape username mapping
  - `trackable_events` - BOTW/SOTW competitions
  - `trackable_event_participations` - User participation with progress
  - `trackable_event_progress` - Progress snapshots over time
  - `schedulable_events` - Mass events, Wildy Wednesday
  - `schedulable_event_participations` - Event registrations
  - `guild_warning_channels` - Warning system config
  - `warnings` - User warnings
- Created goose migration with proper indexes and constraints
- Wrote 60+ SQL queries for all CRUD operations
- Generated type-safe Go code with sqlc (7 files, ~1000 LOC)

**4. Core Components Implemented**

**RuneScape API Client** (`internal/runescape/client.go`):
- HTTP client with 10s timeout
- CSV parsing for OSRS Hiscore API
- Handles 24 skills + boss kill counts
- Error handling for player not found
- Maps to structured domain models

**Discord Bot Core** (`internal/bot/bot.go`):
- Session management with discordgo
- Interaction handler routing (commands, buttons, modals)
- Command registration system
- Graceful startup/shutdown
- Stub implementations for `/link-rsn` and `/unlink-rsn`
- Modal handling infrastructure

**Configuration System** (`config/config.go`):
- Environment variable loading
- Defaults for DATABASE_PATH (~/.voidbound/)
- Optional guild ID for dev command registration
- Log level configuration

**Main Application** (`cmd/voidbound/main.go`):
- Database initialization
- Automatic migrations on startup
- Bot lifecycle management
- Signal handling (SIGTERM, SIGINT)
- Error propagation

**5. Build System & Tooling**

**Makefile** - 25+ targets:
- Build: `build`, `build-windows`, `build-linux`, `build-darwin`, `build-all`, `release`
- Development: `run`, `dev`, `fmt`, `vet`, `lint`, `check`
- Testing: `test`, `coverage`
- Database: `sqlc-generate`, `migrate-up`, `migrate-down`, `migrate-status`, `migrate-create`
- Setup: `init`, `install-tools`, `download`, `tidy`
- Docker: `docker-build`, `docker-run`
- Help: `help` (self-documenting)

**6. Dependencies Installed**
```go
require (
    github.com/bwmarrin/discordgo v0.29.0
    github.com/mattn/go-sqlite3 v1.14.32
    github.com/pressly/goose/v3 v3.25.0
    github.com/robfig/cron/v3 v3.0.1
    // + transitive dependencies
)
```

**7. Build Validation**
- ✅ Project compiles successfully
- ✅ Binary size: 25MB (voidbound.exe)
- ✅ All imports resolved
- ✅ sqlc code generation working
- ✅ No build errors

**8. Documentation**
- Comprehensive README.md with setup instructions
- Refactoring plan with 13 phases documented
- Session state JSON for continuity
- Makefile self-documenting with help target

#### Technical Decisions Made

**1. Database Design**
- **Decision:** Use INTEGER for Discord IDs instead of TEXT
- **Rationale:** Better indexing performance, Discord snowflakes fit in int64
- **Constraint:** Unique constraint on active account links per user

**2. Project Structure**
- **Decision:** Follow standard Go project layout (cmd/, internal/)
- **Rationale:** Community conventions, clear public vs internal API
- **Note:** All bot logic in internal/ (not meant to be imported)

**3. sqlc vs ORM**
- **Decision:** Use sqlc instead of GORM/ent
- **Rationale:** Type safety, performance, explicit queries, no runtime reflection
- **Tradeoff:** More verbose for complex queries, but clearer intent

**4. Migration Strategy**
- **Decision:** Greenfield rewrite, not incremental migration
- **Rationale:** Clean slate, modern patterns, different language paradigms
- **Plan:** Run side-by-side during testing, then switch bot token

**5. Command Architecture**
- **Decision:** Handler map with custom ID routing
- **Rationale:** Flexible routing, supports buttons/modals/commands uniformly
- **Pattern:** "action:data" or "action:data1,data2" in custom IDs

**6. Error Handling**
- **Decision:** Explicit error returns, no panics in handlers
- **Rationale:** Graceful degradation, Discord ephemeral error messages
- **Pattern:** Log errors, respond to user with friendly message

#### Known Issues & Technical Debt

**1. Boss Parsing Incomplete**
- **Issue:** OSRS API boss order not fully mapped in runescape/client.go
- **Impact:** Boss KC tracking needs completion
- **TODO:** Add complete boss name mapping based on API response order
- **Priority:** High (required for BOTW functionality)

**2. Context7 Pattern Not Implemented**
- **Issue:** User requested context7 pattern, not yet integrated
- **Impact:** Context management may need refactoring
- **TODO:** Research context7 pattern and integrate
- **Priority:** Medium (architectural preference)

**3. Command Handlers Stubbed**
- **Issue:** Most command handlers are placeholders
- **Impact:** Bot runs but commands don't do anything
- **TODO:** Implement full command logic (Phase 5-8)
- **Priority:** High (core functionality)

**4. No Scheduled Tasks**
- **Issue:** Cron jobs for notifications/progress tracking not implemented
- **Impact:** Missing background processing
- **TODO:** Add scheduler setup in main.go (Phase 9)
- **Priority:** High (key feature)

**5. No Role Permission Checks**
- **Issue:** Permission middleware not implemented
- **Impact:** Anyone can run admin commands
- **TODO:** Add role checking decorator/middleware
- **Priority:** Critical (security)

**6. No Embeds**
- **Issue:** Discord embeds for events/winners not created
- **Impact:** Responses are plain text
- **TODO:** Create embed builders (Phase 10)
- **Priority:** Medium (UX enhancement)

**7. No Tests**
- **Issue:** Zero test coverage
- **Impact:** No validation of functionality
- **TODO:** Add unit and integration tests (Phase 12)
- **Priority:** Medium (quality assurance)

**8. Build Timeout Issue**
- **Issue:** Windows cross-compilation times out (>2 minutes)
- **Impact:** Makefile build-windows target slow
- **Workaround:** Use direct `go build` or increase timeout
- **Priority:** Low (dev environment specific)

#### Pending Work

**Next Session Priorities:**

**Phase 5: Registration System** (High Priority)
- [ ] Complete `/link-rsn` modal submission handler
- [ ] Fetch player from OSRS API with validation
- [ ] Display player embed with confirmation buttons
- [ ] Handle confirmation/rejection interactions
- [ ] Store account link with proper deactivation logic
- [ ] Implement `/unlink-rsn` full flow
- [ ] Add error handling for API failures

**Phase 6: Trackable Events - BOTW** (High Priority)
- [ ] Implement `/botw wildy` command group
- [ ] Implement `/botw group` command group
- [ ] Implement `/botw quest` command group
- [ ] Implement `/botw slayer` command group
- [ ] Implement `/botw world` command group
- [ ] Create event in database with thread
- [ ] Build registration button handler
- [ ] Build list participants button handler
- [ ] Implement `/finish` command with winner calculation

**Phase 6: Trackable Events - SOTW** (High Priority)
- [ ] Implement `/sotw` command with 23 skill choices
- [ ] Reuse trackable event base logic
- [ ] Skill-specific progress tracking

**Phase 9: Scheduled Tasks** (High Priority)
- [ ] Set up cron scheduler in main.go
- [ ] Implement notification task (every minute)
- [ ] Implement progress tracking task
- [ ] Handle concurrent database access safely
- [ ] Add error recovery for failed tasks

**Phase 10: Utilities** (Medium Priority)
- [ ] Create embed builders for all event types
- [ ] Implement role permission middleware
- [ ] Add Discord ID utilities
- [ ] Time/date formatters for event scheduling

**Later Phases:**
- Phase 7: Schedulable Events (Mass, Wildy Wednesday)
- Phase 8: Admin/Warnings
- Phase 11: Configuration & Deployment
- Phase 12: Testing
- Phase 13: Documentation

#### Files Created This Session

**Source Code (10 files):**
1. `cmd/voidbound/main.go` - 96 lines - Application entry point with migrations
2. `config/config.go` - 47 lines - Configuration loading from env vars
3. `internal/bot/bot.go` - 177 lines - Discord bot core with interaction routing
4. `internal/models/hiscore.go` - 97 lines - OSRS domain models (skills, bosses)
5. `internal/runescape/client.go` - 168 lines - OSRS Hiscore API client with CSV parsing
6. `internal/database/*.go` - 7 files - sqlc generated (auto-generated, ~1000 LOC)

**Database (5 files):**
7. `migrations/00001_initial_schema.sql` - 113 lines - Complete schema with indexes
8. `queries/account_links.sql` - 32 lines - 9 queries for account management
9. `queries/trackable_events.sql` - 82 lines - 21 queries for BOTW/SOTW
10. `queries/schedulable_events.sql` - 54 lines - 11 queries for scheduled events
11. `queries/warnings.sql` - 35 lines - 7 queries for warning system

**Configuration (3 files):**
12. `sqlc.yaml` - 13 lines - sqlc configuration
13. `go.mod` - 17 lines - Go module with dependencies
14. `.env.example` - 11 lines - Environment variable template

**Documentation (4 files):**
15. `README.md` - 210 lines - Comprehensive project documentation
16. `refactor/plan.md` - 680 lines - 13-phase refactoring plan with de-para mapping
17. `refactor/state.json` - 62 lines - Session state tracking
18. `.gitignore` - 35 lines - Git ignore rules

**Build System (1 file):**
19. `Makefile` - 270 lines - 25+ automated targets with documentation

**Total:** 19 files created, ~2,900 lines of code/docs/config

#### Key Code Patterns Established

**1. Database Query Pattern (sqlc)**
```go
// Query definition in queries/*.sql
-- name: GetAccountLinkByDiscordID :one
SELECT * FROM account_links
WHERE discord_member_id = ? AND is_active = 1
LIMIT 1;

// Generated Go code usage
link, err := queries.GetAccountLinkByDiscordID(ctx, discordID)
```

**2. Discord Interaction Pattern**
```go
// Route by interaction type
switch i.Type {
case discordgo.InteractionApplicationCommand:
    handler := b.handlers[i.ApplicationCommandData().Name]
case discordgo.InteractionMessageComponent:
    b.handleComponentInteraction(s, i)
case discordgo.InteractionModalSubmit:
    b.handleModalSubmit(s, i)
}
```

**3. Configuration Pattern**
```go
// Environment variables with defaults
cfg, err := config.Load()
// Defaults: ~/.voidbound/voidbound.db, log level "info"
```

**4. Error Handling Pattern**
```go
// Explicit error returns, no panics
if err != nil {
    return fmt.Errorf("context: %w", err)
}
```

#### Integration Points

**External Services:**
1. **Discord API** - Via discordgo, requires DISCORD_TOKEN
2. **OSRS Hiscore API** - HTTP GET to secure.runescape.com, no auth required
3. **SQLite Database** - Local file, no external service

**Internal Components:**
- Bot ← Config (token, paths)
- Bot ← Database (queries via sqlc)
- Bot ← RSClient (player data)
- Commands → Database (state persistence)
- Commands → RSClient (progress tracking)
- Scheduler → Database (event queries)
- Scheduler → Discord (notifications)

#### Environment Variables Required

```bash
DISCORD_TOKEN=required          # Bot authentication token
DATABASE_PATH=optional          # Default: ~/.voidbound/voidbound.db
LOG_LEVEL=optional             # Default: info (debug|info|warn|error)
DISCORD_GUILD_ID=optional      # For dev: fast command registration
```

#### Build & Run Instructions

```bash
# First time setup
make init                      # Install tools, download deps, generate code

# Development
make run                       # Build and run bot
make dev                       # Live reload (requires air)

# Building
make build                     # Current OS
make build-windows            # Windows binary
make build-all                # All platforms

# Database
make migrate-up               # Run migrations
make sqlc-generate           # Regenerate queries
make migrate-create NAME=xyz # New migration

# Quality
make check                    # fmt + vet + test
make test                     # Run tests
make coverage                # Coverage report
```

#### Architecture Decisions for Next Session

**1. Embed Builder Strategy**
- Create package `internal/embeds/` with builders for each event type
- Use discordgo.MessageEmbed structs
- Include color coding by event type
- Standardize footer with bot branding

**2. Permission Middleware**
- Create `internal/bot/middleware.go` with role checking
- Use Discord role IDs from config or hardcoded
- Return early with ephemeral error if unauthorized
- Apply to all Coordinator/Admin commands

**3. Scheduler Architecture**
- Use robfig/cron in main.go after bot starts
- Create `internal/scheduler/` package with task implementations
- Each task gets DB connection and Discord session
- Handle errors gracefully, log failures, continue running

**4. Command Module Organization**
- Create `internal/commands/` package
- One file per command group (botw.go, sotw.go, mass.go, register.go)
- Register handlers in bot.go from command packages
- Share common logic via internal/commands/base.go

**5. Context7 Integration Research**
- Research context7 pattern (not standard Go pattern)
- Evaluate if structured context fits our use case
- Consider using context.Context with typed keys
- Decision pending research

#### Session Metrics

- **Files Created:** 19
- **Lines of Code:** ~1,200 (excluding generated)
- **Lines of Config/SQL:** ~400
- **Lines of Documentation:** ~1,300
- **Phases Completed:** 4 of 13
- **Features Implemented:** 0 (infrastructure only)
- **Build Status:** ✅ Successful
- **Test Coverage:** 0% (no tests yet)

#### Handoff Notes for Next Developer

**What's Working:**
- ✅ Project builds successfully
- ✅ Database schema complete and migrated
- ✅ Discord bot connects (needs token)
- ✅ Configuration system functional
- ✅ OSRS API client ready
- ✅ Build automation complete

**What Needs Work:**
- ⚠️ All command handlers are stubs
- ⚠️ No scheduled tasks running
- ⚠️ No permission checks
- ⚠️ Boss parsing incomplete
- ⚠️ No embed builders
- ⚠️ No tests

**Quick Start:**
1. Copy `.env.example` to `.env` and add Discord token
2. Run `make init` to set up tools
3. Run `make run` to start bot
4. Bot will connect but commands won't do anything yet

**To Continue Development:**
- Start with Phase 5 (Registration System) - it's the foundation
- Then Phase 6 (BOTW/SOTW) - core functionality
- Then Phase 9 (Scheduled Tasks) - background processing
- Reference `refactor/plan.md` for detailed task breakdown

**Critical Next Steps:**
1. Complete boss name mapping in runescape/client.go
2. Implement full /link-rsn flow with confirmation
3. Add role permission middleware
4. Build event embeds
5. Set up cron scheduler

---

## Important Patterns & Conventions

**Code Style:**
- Follow standard Go conventions (gofmt, golint)
- Error messages: lowercase, no punctuation
- Package comments on every package
- Exported function comments required

**Database Patterns:**
- Use transactions for multi-step operations
- Always defer stmt.Close() and rows.Close()
- Use context.Context for cancellation
- Integer primary keys, not UUIDs

**Discord Patterns:**
- Ephemeral messages for user-specific responses
- Public messages for events/announcements
- Custom IDs format: "action:data1,data2"
- Always respond to interactions (avoid timeout)

**Error Handling:**
- Return errors, don't panic in handlers
- Wrap errors with context: fmt.Errorf("context: %w", err)
- Log errors before responding to user
- User-friendly messages, log technical details

---

## Git Information

**Repository:** Not initialized yet
**Branch:** N/A
**Recent Commits:** None

**Uncommitted Changes:** All files (19 new files)

**Recommended First Commit:**
```bash
git init
git add .
git commit -m "Initial voidbound project setup

- Go project structure with cmd/internal layout
- Database schema with goose migrations
- sqlc configuration and queries for all entities
- RuneScape Hiscore API client implementation
- Discord bot core with discordgo
- Configuration system with environment variables
- Makefile with 25+ build/dev/test targets
- Comprehensive documentation and refactoring plan

Bot compiles and connects but commands need implementation.
Foundation complete for Phase 5+ development."
```

---

## References & Resources

**Original Project:**
- Location: `C:\Users\s.schubert\source\personal\TopezEventBot`
- Language: C#/.NET 8.0
- Framework: Discord.Net v3.12.0
- Database: SQLite + Entity Framework Core
- Scheduler: Coravel v5.0.2

**Documentation:**
- Refactoring Plan: `refactor/plan.md`
- Project README: `README.md`
- OSRS Hiscore API: https://secure.runescape.com/m=hiscore_oldschool/
- discordgo Docs: https://github.com/bwmarrin/discordgo
- sqlc Docs: https://docs.sqlc.dev/

**Key Dependencies:**
- github.com/bwmarrin/discordgo@v0.29.0
- github.com/pressly/goose/v3@v3.25.0
- github.com/robfig/cron/v3@v3.0.1
- github.com/mattn/go-sqlite3@v1.14.32

---

## Notes for Future Sessions

**Before Starting:**
- Check `refactor/state.json` for current phase
- Review pending tasks in this document
- Verify build still works: `make build`
- Check for dependency updates: `go list -u -m all`

**During Development:**
- Update `refactor/state.json` after completing phases
- Add tests as you implement features
- Run `make check` before committing
- Update README.md with new features

**When Stuck:**
- Reference original C# code in TopezEventBot
- Check `refactor/plan.md` for detailed task breakdown
- Discord bot examples: https://github.com/bwmarrin/discordgo/tree/master/examples
- sqlc examples: https://github.com/sqlc-dev/sqlc/tree/main/examples

**Performance Considerations:**
- OSRS API has rate limits (unknown, test cautiously)
- SQLite handles concurrency via serialization (single writer)
- Discord rate limits: 50 commands per guild per hour
- Batch DB operations where possible

---

## End of Session 1

**Status:** Foundation Complete ✅
**Next Session:** Begin Phase 5 (Registration System Implementation)
**Estimated Completion:** 40% infrastructure, 0% features (10-12 more sessions estimated)

---

## Session 2: 2025-10-01 - Phase 5 Complete + Enhancements

**Duration:** Full session
**Status:** Phase 5 Complete ✅, Enhancements Added

### Accomplished

**Phase 5: Registration System - COMPLETE**

**1. Full `/link-rsn` Implementation**
- Created modal-based username input
- OSRS Hiscore API integration with error handling
- Player stats embed with thumbnail and formatted numbers
- Confirmation button flow ("That's me!" / "Not me")
- Database transaction with proper link management
- Auto-deactivation of previous links
- Reactivation of existing inactive links
- Comprehensive error handling and user feedback

**2. Full `/unlink-rsn` Implementation**
- Fetches and deactivates active account link
- Validates user has linked account
- Friendly error messages
- Transaction-safe operations

**3. Embed System Created** (`internal/embeds/embeds.go` - 367 lines)

Created 9 embed types:
- `PlayerInfo()` - Player stats with level, XP, rank, combat stats
- `BossOfTheWeek()` - BOTW event announcements
- `SkillOfTheWeek()` - SOTW event announcements
- `EventWinners()` - Winner displays with 🥇🥈🥉 medals
- `MassEvent()` - Mass event scheduling with Discord timestamps
- `WildyWednesday()` - Wilderness event scheduling
- `ScheduledEventReminder()` - DM reminders for events
- `ErrorEmbed()` - Consistent error messaging
- `SuccessEmbed()` - Consistent success messaging

Features:
- Color-coded by event type (purple BOTW, turquoise SOTW, etc.)
- Number formatting with commas (1,234,567)
- Discord timestamp formatting (`<t:unix:F>`)
- OSRS logo thumbnails
- Footer text for context

**4. Registration Command Module** (`internal/commands/register.go` - 321 lines)

Implemented:
- `RegisterCommands` struct with DB, OSRS client
- `HandleLinkRSN()` - Shows modal
- `HandleLinkRSNModal()` - Processes submission, fetches player, shows embed
- `HandleConfirmRSN()` - Stores link with transaction
- `HandleCancelRSN()` - Cancels linking process
- `HandleUnlinkRSN()` - Deactivates link

Technical details:
- Discord ID parsing (string → int64 with strconv)
- Database transactions (Begin/Commit/Rollback pattern)
- Component ID routing ("action:data" pattern)
- Ephemeral messages for privacy
- Case-insensitive username matching

**5. Bot Integration** (`internal/bot/bot.go` - Updated)

Changes:
- Added `registerCmds` field to Bot struct
- Updated `handleComponentInteraction()` with custom ID parser
- Routes "confirm-rsn:username" and "cancel-rsn:username" buttons
- Updated `handleModalSubmit()` to route to registration handlers
- Removed placeholder stub handlers

**6. Build System Enhancements**

Makefile improvements:
- ✅ `make build-windows` working (8.9MB optimized binary)
- Build flags: `-trimpath -ldflags "-s -w"` for smaller binaries
- 64% size reduction vs debug build (8.9MB vs 25MB)

**7. Configuration Enhancement**

Added godotenv support:
- Installed `github.com/joho/godotenv@v1.5.1`
- Updated `config/config.go` to load `.env` files
- Maintains backward compatibility with system env vars
- Priority: System env → .env file → defaults

**8. Documentation**

Created:
- `refactor/phase5_complete.md` - Comprehensive phase documentation (320 lines)
- Documented all features, technical decisions, and integration points
- Listed known limitations and next phase preview

### Files Created/Modified This Session

**New Files (3):**
1. `internal/embeds/embeds.go` - 367 lines - Embed builder system
2. `internal/commands/register.go` - 321 lines - Registration handlers
3. `refactor/phase5_complete.md` - 320 lines - Phase documentation

**Modified Files (2):**
1. `internal/bot/bot.go` - Updated for command integration and routing
2. `config/config.go` - Added godotenv support

**Total New Code:** ~690 lines (excluding documentation)

### Technical Decisions Made

**1. Wise Old Man API Migration (Pending)**
- **Decision:** User requested switch from OSRS Hiscore to Wise Old Man API
- **Rationale:** Better data structure, more reliable, includes computed stats (EHP/EHB)
- **Status:** Research started, implementation pending
- **Blocker:** User requested context7 pattern integration (non-standard pattern)
- **Action Required:** Research context7, implement WOM client, update models

**2. Context7 Pattern Integration (Pending)**
- **Request:** User wants context7 pattern for structured context
- **Status:** Research needed - not a standard Go pattern
- **Note:** User attempted MCP server connection for context7
- **Next Steps:**
  - Understand context7 requirements
  - Evaluate implementation approach
  - Integrate with WOM API client

**3. Discord ID Type Handling**
- **Challenge:** discordgo returns IDs as strings, DB uses int64
- **Solution:** `strconv.ParseInt(idStr, 10, 64)` with error handling
- **Applied:** All registration handlers

**4. Transaction Pattern Established**
```go
tx, err := r.DBSQL.Begin()
defer tx.Rollback()
qtx := r.DB.WithTx(tx)
// operations...
tx.Commit()
```

**5. Custom ID Routing Pattern**
- Format: `"action:data"`
- Parser: `strings.SplitN(customID, ":", 2)`
- Examples: "confirm-rsn:PlayerName", "cancel-rsn:PlayerName"

### Known Issues & Technical Debt

**1. Wise Old Man Migration Incomplete (HIGH PRIORITY)**
- **Issue:** User wants WOM API instead of OSRS Hiscore
- **Current:** Using OSRS Hiscore API
- **Impact:** Need to rewrite `internal/runescape/client.go`
- **TODO:**
  - Create WOM client with JSON parsing
  - Update models for WOM data structures
  - Test API endpoints
  - Handle player not found errors
- **Blocker:** Context7 pattern understanding needed

**2. Context7 Pattern Not Implemented (HIGH PRIORITY)**
- **Issue:** User requested context7 pattern for context management
- **Research:** Not a standard Go pattern, needs investigation
- **Impact:** Architectural decision affects all API clients
- **TODO:**
  - Research context7 pattern/library
  - Understand user's requirements
  - Implement or adapt pattern
  - Document usage

**3. Boss Parsing Still Incomplete**
- **Issue:** OSRS API boss mapping incomplete in current client
- **Impact:** Will need similar fix in WOM client
- **Priority:** Medium (will be addressed during WOM migration)

**4. No Role Permission Checks**
- **Issue:** Anyone can use commands
- **Impact:** Security risk for admin commands
- **Priority:** High (needed before BOTW/SOTW commands)
- **TODO:** Create middleware in next phase

**5. No Scheduled Tasks Yet**
- **Issue:** Background jobs not implemented
- **Impact:** No progress tracking or notifications
- **Priority:** High (Phase 9)

**6. No Tests**
- **Issue:** Zero test coverage
- **Impact:** No automated validation
- **Priority:** Medium (Phase 12)

### Pending Work for Next Session

**IMMEDIATE PRIORITIES:**

**1. Research & Implement Context7 Pattern**
- Understand what context7 means in this context
- Check if it's a specific library or pattern
- Implement or adapt for Go
- Document the pattern

**2. Wise Old Man API Integration**
- Create new client: `internal/wiseoldman/client.go`
- Update models in `internal/models/`
- WOM API endpoints:
  - `GET /v2/players/{username}` - Fetch player data
  - `POST /v2/players/{username}` - Update player (optional)
- Parse JSON response (easier than CSV!)
- Handle errors (404 = player not found)
- Update registration commands to use WOM

**3. Update Models for WOM**
```go
type WOMPlayer struct {
    ID          int64
    Username    string
    DisplayName string
    Type        string
    Build       string
    LatestSnapshot WOMSnapshot
}

type WOMSnapshot struct {
    Data struct {
        Skills map[string]WOMSkill
        Bosses map[string]WOMBoss
    }
}
```

**LATER PRIORITIES:**

**4. Role Permission Middleware** (before Phase 6)
- Create `internal/bot/middleware.go`
- Check Discord role IDs
- Decorator pattern for handlers
- Apply to admin commands

**5. Continue Phase 6: BOTW/SOTW** (after WOM migration)
- Implement `/botw` command groups
- Implement `/sotw` command
- Event registration buttons
- Progress tracking

### Wise Old Man API Details

**Base URL:** `https://api.wiseoldman.net/v2`

**Key Endpoints:**
- `GET /players/{username}` - Fetch player
- `POST /players/{username}` - Update/track player
- `GET /players/{id}/snapshots` - Historical data

**Response Structure (Tested):**
```json
{
  "id": 1057,
  "username": "lynx titan",
  "displayName": "Lynx titan",
  "latestSnapshot": {
    "data": {
      "skills": {
        "overall": {"experience": 4600000000, "rank": 1, "level": 2277},
        "attack": {"experience": 200000000, "rank": 15, "level": 99}
      },
      "bosses": {
        "tztok_jad": {"kills": 186, "rank": 362}
      }
    }
  }
}
```

**Benefits over OSRS Hiscore:**
- JSON (easier to parse)
- Computed stats (EHP/EHB)
- Player IDs (stable identifier)
- More reliable uptime
- Historical snapshots
- Better error responses

### Code Patterns Established

**1. Embed Builder Pattern**
```go
embed := embeds.PlayerInfo(player)
// Returns *discordgo.MessageEmbed
```

**2. Modal Flow Pattern**
```go
// 1. Show modal
InteractionResponseModal

// 2. Process modal submission
InteractionResponseDeferredChannelMessageWithSource

// 3. Followup with result
s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{...})
```

**3. Button Confirmation Pattern**
```go
// Show embed with buttons
components := []discordgo.MessageComponent{
    discordgo.ActionsRow{
        Components: []discordgo.MessageComponent{
            discordgo.Button{CustomID: "confirm-rsn:username", ...},
            discordgo.Button{CustomID: "cancel-rsn:username", ...},
        },
    },
}
```

**4. Database Transaction Pattern**
```go
tx, err := r.DBSQL.Begin()
defer tx.Rollback()
qtx := r.DB.WithTx(tx)
// ... operations ...
if err := tx.Commit(); err != nil { ... }
```

### Session Metrics

- **Phases Completed:** 5 of 13
- **Features Working:** 2 commands (/link-rsn, /unlink-rsn)
- **Build Status:** ✅ Successful (8.9MB optimized)
- **Lines Added:** ~690 (code) + 320 (docs)
- **Files Created:** 3
- **Files Modified:** 2
- **Test Coverage:** 0% (no tests yet)

### Integration Points for WOM Migration

**Files to Update:**
1. `internal/models/hiscore.go` - Add WOM structs
2. `internal/wiseoldman/client.go` - NEW - WOM HTTP client
3. `internal/commands/register.go` - Update client reference
4. `internal/bot/bot.go` - Update client initialization

**Pattern to Follow:**
```go
// New client
type WOMClient struct {
    httpClient *http.Client
    baseURL    string
    ctx        context7.Context // User's requested pattern
}

// Fetch player
func (c *WOMClient) GetPlayer(ctx context.Context, username string) (*WOMPlayer, error) {
    // GET https://api.wiseoldman.net/v2/players/{username}
    // Parse JSON response
    // Return structured data
}
```

### Handoff Notes

**What's Working:**
- ✅ Full registration system (link/unlink)
- ✅ Embed system for all event types
- ✅ Database transactions
- ✅ Component routing
- ✅ Modal handling
- ✅ Build system with optimization
- ✅ .env file support

**What Needs Work:**
- ⚠️ Switch to Wise Old Man API (user requirement)
- ⚠️ Implement context7 pattern (user requirement)
- ⚠️ Update models for WOM data structures
- ⚠️ Add role permission checks
- ⚠️ Implement BOTW/SOTW commands
- ⚠️ Add scheduled tasks

**Blockers:**
- Context7 pattern needs research/understanding
- User attempted MCP server connection (unclear purpose)
- Need to understand user's expectations for context7

**Quick Start for Next Session:**
1. Research context7 - what does user expect?
2. Create WOM client with context7 integration
3. Update models for WOM data structures
4. Test WOM API integration
5. Update registration commands to use WOM
6. Continue with Phase 6 (BOTW/SOTW)

**Critical Next Steps:**
1. Clarify context7 requirements with user
2. Implement WOM API client
3. Test player lookup with WOM
4. Verify embed still renders correctly with WOM data
5. Add role permissions before admin commands

---

## Important Context

### Context7 Mystery
- User requested "Use context7" for WOM integration
- Not a standard Go pattern in common use
- User attempted MCP server connection: `claude mcp add --transport http context7 https://mcp.context7.com/mcp`
- May be referring to a specific service or pattern
- **Action Required:** Clarify with user what context7 means

### Wise Old Man API
- Well-documented REST API
- JSON responses (much easier than CSV)
- Free tier available
- Rate limiting unknown (test cautiously)
- Player data more comprehensive than OSRS API

### Development Flow
User prefers:
- Makefile for builds
- .env for configuration
- Concise responses
- Direct implementation over explanation

---

## End of Session 2

**Status:** Phase 5 Complete ✅, WOM Migration Pending
**Next Session:** Research context7, implement WOM client, continue Phase 6
**Estimated Completion:** 50% infrastructure, 15% features

**Blockers:** Context7 pattern understanding needed

This memory will persist across sessions for seamless development continuity.
