# Testing Implementation Plan - Voidling Bot

**Created:** 2025-10-03
**Status:** Starting
**Session ID:** implement-testing-2025-10-03

## Source Analysis

**Source Type:** New Testing Infrastructure
**Target:** Go Discord Bot (voidling) - Complete test coverage

**Current State:**
- ❌ **0 test files** in the entire project
- ✅ testify dependency already installed (v1.11.0)
- ✅ Makefile has test targets defined
- ❌ No test infrastructure setup
- ❌ No mocking framework for Discord/WOM APIs

**Packages Requiring Tests:**
1. `internal/wiseoldman` - WOM API client (2 files)
2. `internal/commands` - Command handlers (7 files)
3. `internal/embeds` - Embed builders (2 files)
4. `internal/bot` - Core bot logic (2 files)
5. `internal/database` - Database layer (9 files, mostly generated)
6. `internal/models` - Data models (1 file)
7. `internal/timezone` - Timezone utilities (1 file)
8. `config` - Configuration (1 file)

## Testing Strategy

### Test Categories

**1. Unit Tests** (Priority: High)
- Test individual functions in isolation
- Mock external dependencies (Discord, WOM, Database)
- Target: 70%+ coverage for business logic

**2. Integration Tests** (Priority: Medium)
- Test component interactions
- Use in-memory SQLite for database tests
- Mock only Discord/WOM APIs

**3. Table-Driven Tests** (Priority: High)
- Use Go's table-driven pattern for comprehensive coverage
- Test edge cases, error paths, and happy paths

## Implementation Tasks

### ✅ Phase 1: Setup Testing Infrastructure (30 min)
- [x] Analyze existing code structure
- [ ] Install additional test dependencies
  - [ ] Add `github.com/stretchr/testify/mock` for mocking
  - [ ] Add `github.com/stretchr/testify/suite` for test suites
- [ ] Create test helper utilities
  - [ ] `internal/testutil/discord_mock.go` - Mock Discord session
  - [ ] `internal/testutil/database.go` - In-memory DB setup
  - [ ] `internal/testutil/fixtures.go` - Test data fixtures
- [ ] Update Makefile with enhanced test targets
  - [ ] `make test-unit` - Unit tests only
  - [ ] `make test-integration` - Integration tests
  - [ ] `make test-coverage` - Generate coverage report
  - [ ] `make test-watch` - Watch mode for TDD

### Phase 2: WOM Client Tests (45 min)
- [ ] Create `internal/wiseoldman/client_test.go`
  - [ ] Test `GetPlayer()` success case
  - [ ] Test `GetPlayer()` player not found (404)
  - [ ] Test `GetPlayer()` network errors
  - [ ] Test `UpdatePlayer()` success
  - [ ] Test `CreateCompetition()` success
  - [ ] Test `CreateCompetition()` validation errors
  - [ ] Test `AddParticipantsToCompetition()` batch operations
  - [ ] Mock HTTP client for all tests
- [ ] Create `internal/wiseoldman/models_test.go`
  - [ ] Test `Player.GetSkill()` for valid skills
  - [ ] Test `Player.GetSkill()` for invalid skills
  - [ ] Test `Player.GetBoss()` for valid bosses
  - [ ] Test JSON unmarshaling

### Phase 3: Command Handler Tests (90 min)
- [ ] Create `internal/commands/register_test.go`
  - [ ] Test `HandleLinkRSN()` modal display
  - [ ] Test `HandleLinkRSNModal()` with valid username
  - [ ] Test `HandleLinkRSNModal()` with empty username
  - [ ] Test `HandleConfirmRSN()` new link creation
  - [ ] Test `HandleConfirmRSN()` link reactivation
  - [ ] Test `HandleConfirmRSN()` duplicate active link
  - [ ] Test `HandleConfirmRSN()` nickname update success/failure
  - [ ] Test `HandleUnlinkRSN()` success
  - [ ] Test `HandleUnlinkRSN()` no linked account
  - [ ] Test helper methods (getUserAndGuildIDs, parseDiscordID, etc.)

- [ ] Create `internal/commands/trackable_test.go`
  - [ ] Test `StartEvent()` BOTW creation
  - [ ] Test `StartEvent()` SOTW creation
  - [ ] Test `StartEvent()` thread creation failure
  - [ ] Test `StartEvent()` WOM competition creation failure
  - [ ] Test `RegisterForEvent()` first-time registration
  - [ ] Test `RegisterForEvent()` duplicate registration
  - [ ] Test `RegisterForEvent()` no linked account
  - [ ] Test `ListParticipants()` with participants
  - [ ] Test `ListParticipants()` empty list
  - [ ] Test finish logic (winner calculation)

- [ ] Create `internal/commands/schedulable_test.go`
  - [ ] Test `HandleMassEvent()` creation
  - [ ] Test timezone parsing
  - [ ] Test participation button handling

- [ ] Create `internal/commands/config_test.go`
  - [ ] Test role configuration commands
  - [ ] Test channel configuration commands
  - [ ] Test permission checks

### Phase 4: Embed Builder Tests (30 min)
- [ ] Create `internal/embeds/embeds_test.go`
  - [ ] Test `PlayerInfo()` with valid player
  - [ ] Test `PlayerInfo()` with nil player
  - [ ] Test `BossOfTheWeek()` embed generation
  - [ ] Test `SkillOfTheWeek()` embed generation
  - [ ] Test `EventWinners()` with 3 winners
  - [ ] Test `EventWinners()` with fewer than 3
  - [ ] Test `MassEvent()` timestamp formatting
  - [ ] Test `ErrorEmbed()` and `SuccessEmbed()`
  - [ ] Test `formatNumber()` with various inputs
  - [ ] Test `formatActivityName()` snake_case conversion

- [ ] Create `internal/embeds/boss_info_test.go`
  - [ ] Test `GetBossInfo()` for known bosses
  - [ ] Test `GetBossInfo()` for unknown bosses

### Phase 5: Bot Core Tests (60 min)
- [ ] Create `internal/bot/bot_test.go`
  - [ ] Test `New()` initialization
  - [ ] Test `Start()` success
  - [ ] Test `Stop()` cleanup
  - [ ] Test command registration
  - [ ] Test interaction routing (commands, buttons, modals)
  - [ ] Test `handleComponentInteraction()` parsing
  - [ ] Test `handleModalSubmit()` routing
  - [ ] Test `handleGuildMemberAdd()` greeting DM

- [ ] Create `internal/bot/permissions_test.go`
  - [ ] Test `HasPermission()` for admin
  - [ ] Test `HasPermission()` for coordinator
  - [ ] Test `HasPermission()` for regular user
  - [ ] Test role hierarchy

### Phase 6: Database Layer Tests (45 min)
- [ ] Create `internal/database/database_test.go`
  - [ ] Set up in-memory SQLite for tests
  - [ ] Test account link CRUD operations
  - [ ] Test trackable event queries
  - [ ] Test schedulable event queries
  - [ ] Test WOM competition queries
  - [ ] Test transaction handling
  - [ ] Test concurrent access (Go routines)

### Phase 7: Utility Tests (20 min)
- [ ] Create `internal/timezone/timezone_test.go`
  - [ ] Test `SearchTimezones()` with various queries
  - [ ] Test timezone validation

- [ ] Create `internal/models/events_test.go`
  - [ ] Test event type constants
  - [ ] Test hiscore field validation

- [ ] Create `config/config_test.go`
  - [ ] Test `Load()` with env vars
  - [ ] Test `Load()` with .env file
  - [ ] Test `Load()` with defaults
  - [ ] Test missing required fields

### Phase 8: Integration Tests (60 min)
- [ ] Create `tests/integration/registration_test.go`
  - [ ] Test full link-rsn flow (modal → WOM → DB → nickname)
  - [ ] Test full unlink-rsn flow

- [ ] Create `tests/integration/botw_test.go`
  - [ ] Test full BOTW event lifecycle
  - [ ] Test registration → tracking → finish → winners

- [ ] Create `tests/integration/sotw_test.go`
  - [ ] Test full SOTW event lifecycle

- [ ] Create `tests/integration/mass_test.go`
  - [ ] Test mass event creation and participation

### Phase 9: Test Infrastructure Enhancement (30 min)
- [ ] Create mock implementations
  - [ ] `internal/testutil/mock_discord.go` - Discord session mock
  - [ ] `internal/testutil/mock_wom.go` - WOM client mock
  - [ ] `internal/testutil/mock_db.go` - Database mock (if needed)

- [ ] Create test fixtures
  - [ ] `testdata/wom_responses/` - Sample WOM API responses
  - [ ] `testdata/discord_interactions/` - Sample Discord payloads
  - [ ] `testdata/database/` - Sample database states

- [ ] Add test helpers
  - [ ] `CreateTestPlayer()` - Generate test player data
  - [ ] `CreateTestInteraction()` - Generate Discord interactions
  - [ ] `AssertEmbedEquals()` - Compare embeds for tests

### Phase 10: Coverage & CI (30 min)
- [ ] Update Makefile
  - [ ] Add `test-coverage-html` target (opens in browser)
  - [ ] Add `test-race` for race condition detection
  - [ ] Add `test-bench` for benchmarks

- [ ] Create GitHub Actions workflow (optional)
  - [ ] `.github/workflows/test.yml`
  - [ ] Run tests on push/PR
  - [ ] Upload coverage to codecov

- [ ] Add coverage badge to README
- [ ] Document testing approach

## Test Patterns & Best Practices

### Table-Driven Tests Example
```go
func TestFormatNumber(t *testing.T) {
    tests := []struct {
        name     string
        input    int64
        expected string
    }{
        {"small number", 123, "123"},
        {"thousand", 1000, "1,000"},
        {"million", 1000000, "1,000,000"},
        {"zero", 0, "0"},
        {"negative", -1234, "-1,234"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := formatNumber(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Mock Discord Session
```go
type MockDiscordSession struct {
    mock.Mock
}

func (m *MockDiscordSession) InteractionRespond(i *discordgo.Interaction, resp *discordgo.InteractionResponse) error {
    args := m.Called(i, resp)
    return args.Error(0)
}
```

### In-Memory Database Setup
```go
func setupTestDB(t *testing.T) (*sql.DB, *database.Queries) {
    db, err := sql.Open("sqlite3", ":memory:")
    require.NoError(t, err)

    // Run migrations
    goose.SetDialect("sqlite3")
    err = goose.Up(db, "./migrations")
    require.NoError(t, err)

    queries := database.New(db)
    return db, queries
}
```

## Validation Checklist

- [ ] All packages have corresponding test files
- [ ] Unit test coverage > 70%
- [ ] Integration tests cover critical flows
- [ ] All tests pass (`make test`)
- [ ] No race conditions (`make test-race`)
- [ ] Coverage report generated
- [ ] Test documentation complete
- [ ] CI/CD integration (optional)

## Success Metrics

**Coverage Targets:**
- `internal/commands`: 80%+
- `internal/wiseoldman`: 90%+ (critical API integration)
- `internal/embeds`: 70%+
- `internal/bot`: 75%+
- `config`: 85%+
- Overall project: 75%+

**Quality Metrics:**
- ✅ All edge cases covered
- ✅ Error paths tested
- ✅ Mocks properly isolate units
- ✅ Tests run in < 5 seconds (unit tests)
- ✅ Integration tests in < 30 seconds
- ✅ No flaky tests

## Risk Mitigation

**Challenges:**
1. **Discord API Mocking** - Complex interaction types
   - Solution: Create comprehensive mock with common patterns

2. **Database Testing** - Generated sqlc code
   - Solution: Focus on testing queries, not generated code

3. **Async Operations** - Go routines in finish command
   - Solution: Use sync primitives properly in tests

4. **External API Deps** - WOM API changes
   - Solution: Use recorded responses (vcr pattern)

**Rollback Strategy:**
- Tests are additive - no risk to existing code
- Can implement incrementally
- Each test file is independent

## Next Steps After Testing

1. Add benchmark tests for performance-critical paths
2. Implement mutation testing with `go-mutesting`
3. Add property-based testing with `gopter`
4. Create E2E tests with real Discord bot in test server
5. Add load testing for concurrent event handling

## Notes

- Use `testify/suite` for setup/teardown in complex tests
- Mock external APIs (Discord, WOM) to ensure fast, reliable tests
- Use `testing.Short()` to skip slow tests in dev workflow
- Keep test data in `testdata/` directory
- Follow naming: `TestFunctionName` for tests, `BenchmarkFunctionName` for benchmarks
- Use subtests (`t.Run`) for better organization and parallel execution
