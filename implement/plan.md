# Implementation Plan - BOTW & SOTW Commands
**Created:** 2025-10-01
**Status:** In Progress

## Source Analysis

**Source Type:** Existing C# Implementation (TopezEventBot)
**Target:** Go Discord Bot (voidbound)

**Core Features to Implement:**
1. Boss of the Week (BOTW) - 5 categories of bosses
   - `/botw wildy` - Wilderness bosses (11 choices)
   - `/botw group` - Group bosses (7 choices)
   - `/botw quest` - Quest bosses (9 choices)
   - `/botw slayer` - Slayer bosses (6 choices)
   - `/botw world` - World bosses (9 choices)

2. Skill of the Week (SOTW) - 16 skills
   - `/sotw start` - Select from 16 non-combat skills

3. Shared functionality (base trackable events):
   - Event creation with Discord thread
   - Registration button ("Register" + "List participants")
   - `/finish` command to end event and calculate winners
   - Progress tracking (starting point → ending point)
   - Winner announcement with embeds

**Dependencies:**
- Wise Old Man API (already integrated)
- Discord thread creation
- Database operations (sqlc queries already exist)
- Embeds (already created but need event-specific ones)

**Complexity:** Medium-High
- ~5 files to create
- ~1000 lines of code
- Database integration complex (join queries)
- Button/modal interactions

## Target Integration

**Integration Points:**
1. `internal/bot/bot.go` - Register new command handlers
2. `internal/commands/` - New package for trackable events
3. `internal/database/` - Use existing sqlc queries
4. `internal/embeds/embeds.go` - Already has BOTW/SOTW embeds
5. `internal/wiseoldman/` - Player data fetching

**Affected Files:**
- **Create:**
  - `internal/commands/trackable.go` - Base handler for shared logic
  - `internal/commands/botw.go` - BOTW command handlers
  - `internal/commands/sotw.go` - SOTW command handlers
  - `internal/commands/choices.go` - Boss/skill choice definitions

- **Modify:**
  - `internal/bot/bot.go` - Add command registration + button routing
  - `internal/bot/commands.go` - Add new slash commands

**Pattern Matching:**
- Follow existing `internal/commands/register.go` structure
- Use same database transaction pattern
- Reuse button/component interaction routing
- Match error handling style (ephemeral messages)

## Implementation Tasks

### ✅ Phase 1: Setup (5 min)
- [x] Create `implement/` directory
- [x] Write implementation plan
- [x] Analyze C# source code
- [x] Map database queries to Go implementation

### Phase 2: Data Structures (15 min)
- [ ] Create `internal/commands/choices.go`
  - [ ] Boss choice enums (WildyBosses, GroupBosses, QuestBosses, SlayerBosses, WorldBosses)
  - [ ] Skill choice enum (16 non-combat skills)
  - [ ] Mapping functions to HiscoreField/WOM metric names

### Phase 3: Base Trackable Event Handler (30 min)
- [ ] Create `internal/commands/trackable.go`
  - [ ] `TrackableCommands` struct with DB, WOMClient
  - [ ] `StartEvent()` - Creates event + thread + embed with buttons
  - [ ] `RegisterForEvent()` - Handles registration button click
  - [ ] `ListParticipants()` - Shows participant list
  - [ ] `FinishEvent()` - Fetches end points, calculates winners, announces
  - [ ] Thread creation helper
  - [ ] WOM data extraction helpers

### Phase 4: BOTW Commands (20 min)
- [ ] Create `internal/commands/botw.go`
  - [ ] `HandleBOTWWildy()` - `/botw wildy <boss>`
  - [ ] `HandleBOTWGroup()` - `/botw group <boss>`
  - [ ] `HandleBOTWQuest()` - `/botw quest <boss>`
  - [ ] `HandleBOTWSlayer()` - `/botw slayer <boss>`
  - [ ] `HandleBOTWWorld()` - `/botw world <boss>`
  - [ ] All call `StartEvent()` with boss name

### Phase 5: SOTW Commands (10 min)
- [ ] Create `internal/commands/sotw.go`
  - [ ] `HandleSOTWStart()` - `/sotw start <skill>`
  - [ ] Calls `StartEvent()` with skill name

### Phase 6: Button Handlers (15 min)
- [ ] Add button routing in `internal/bot/bot.go`
  - [ ] `register-for-botw:<eventId>,<threadId>`
  - [ ] `register-for-sotw:<eventId>,<threadId>`
  - [ ] `list-participants-botw:<eventId>`
  - [ ] `list-participants-sotw:<eventId>`
  - [ ] Route to trackable handler methods

### Phase 7: Command Registration (15 min)
- [ ] Update `internal/bot/bot.go`
  - [ ] Create `trackableCmds` field
  - [ ] Initialize in `New()`
  - [ ] Register BOTW subcommands (wildy, group, quest, slayer, world)
  - [ ] Register SOTW subcommand (start)
  - [ ] Register finish command
  - [ ] Add handlers to handler map

### Phase 8: Testing & Validation (30 min)
- [ ] Build and verify compilation
- [ ] Test `/botw wildy <boss>` command
  - [ ] Verify thread creation
  - [ ] Check embed display
  - [ ] Test registration button
  - [ ] Verify database insertion
- [ ] Test `/sotw start <skill>` command
- [ ] Test `/finish` command
  - [ ] Verify WOM data fetch
  - [ ] Check winner calculation
  - [ ] Validate embed display
- [ ] Test edge cases
  - [ ] No participants
  - [ ] Player not found in WOM
  - [ ] Duplicate registration
  - [ ] Finish with no active event

## Validation Checklist

- [ ] All BOTW commands implemented (5 categories)
- [ ] SOTW command implemented
- [ ] Registration button works for both event types
- [ ] List participants button works
- [ ] Finish command calculates winners correctly
- [ ] Thread creation works
- [ ] Embeds display properly
- [ ] Database operations succeed
- [ ] WOM API integration works
- [ ] No broken functionality
- [ ] Build successful

## Risk Mitigation

**Potential Issues:**
1. Discord thread creation API differences (discordgo vs Discord.Net)
2. WOM API metric names may differ from C# implementation
3. Complex database queries with joins
4. Button interaction state management
5. Parallel processing of finish command (Go concurrency patterns)

**Rollback Strategy:**
- Git checkpoints after each phase
- Keep C# implementation as reference
- Test incrementally to catch issues early

## C# → Go Pattern Mapping

### Event Creation Flow
**C#:**
```csharp
var eventId = await db.CreateEvent(_eventType, activity, isActive);
var threadId = await NewThreadInCurrentChannelAsync(activity, _eventType);
componentBuilder.AddRow(...)
await FollowupAsync(embed: ..., components: componentBuilder.Build());
```

**Go:**
```go
eventID, err := r.DB.CreateTrackableEvent(ctx, ...)
threadID, err := createThread(s, i.ChannelID, eventName)
components := []discordgo.MessageComponent{...}
s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
    Embeds: []*discordgo.MessageEmbed{embed},
    Components: components,
})
```

### Registration Flow
**C#:**
```csharp
var player = await _rsClient.LoadPlayer(linkedAccount.RunescapeName);
@event.EventParticipations.Add(new TrackableEventParticipation() {
    StartingPoint = player.Bosses[@event.Activity].KillCount,
});
await db.SaveChangesAsync();
```

**Go:**
```go
player, err := r.WOMClient.GetPlayer(ctx, linkedAccount.RunescapeName)
startingPoint := getStartingPoint(player, eventActivity, eventType)
_, err = r.DB.CreateTrackableParticipation(ctx, ...)
```

### Finish Event Flow
**C#:**
```csharp
await Parallel.ForEachAsync(participants, async (participation, token) => {
    var player = await _rsClient.LoadPlayer(...);
    participation.EndPoint = player.Bosses[activity].KillCount;
});
db.UpdateRange(participants);
var result = participants.OrderByDescending(x => x.Progress).Take(3);
```

**Go:**
```go
var wg sync.WaitGroup
for _, p := range participants {
    wg.Add(1)
    go func(p Participation) {
        defer wg.Done()
        player, _ := r.WOMClient.GetPlayer(ctx, p.RunescapeName)
        endPoint := getEndPoint(player, eventActivity, eventType)
        r.DB.UpdateTrackableParticipationEndPoint(ctx, ...)
    }(p)
}
wg.Wait()
winners, _ := r.DB.GetEventWinners(ctx, eventID)
```

## Notes

- WOM API uses snake_case for metrics (e.g., "abyssal_sire", "king_black_dragon")
- Need to map C# HiscoreField names to WOM metric names
- Thread archive duration: 1 week (Discord limit)
- Event embeds already created in `internal/embeds/embeds.go`
- SQL queries already exist in `queries/trackable_events.sql`
- Context handling: use `context.Background()` for now, improve later

## Next Steps After Implementation

1. Add permission middleware (RequireRole equivalent)
2. Implement scheduled progress tracking
3. Add leaderboard command
4. Implement Mass/Wildy Wednesday events
5. Add tests for trackable events
