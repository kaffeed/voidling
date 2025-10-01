# WOM Competition Integration Plan

## Overview
Replace custom event tracking with Wise Old Man competition system.

## Benefits
- ✅ WOM handles all progress tracking automatically
- ✅ Built-in leaderboards and standings
- ✅ Historical snapshots maintained by WOM
- ✅ No need for our scheduler/progress tracking
- ✅ Web interface at wiseoldman.net for participants
- ✅ Better data accuracy and reliability

## Architecture Changes

### Current System (To Remove)
- `trackable_events` table
- `trackable_event_participations` table
- `trackable_event_progress` table
- Custom progress tracking scheduler
- Manual winner calculation

### New System (WOM-Based)
- Store WOM competition ID and verification code
- Let WOM handle progress tracking
- Fetch standings from WOM API when finishing
- Simple mapping table: Discord event → WOM competition

## Database Schema Changes

**New table: `wom_competitions`**
```sql
CREATE TABLE wom_competitions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wom_competition_id INTEGER NOT NULL UNIQUE,
    verification_code TEXT NOT NULL,
    discord_thread_id TEXT NOT NULL,
    metric TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('BOSS_OF_THE_WEEK', 'SKILL_OF_THE_WEEK')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## WOM API Integration

### Create Competition
**Endpoint:** `POST /v2/competitions`
**Body:**
```json
{
  "title": "BOTW - Callisto",
  "metric": "callisto",
  "startsAt": "2025-10-01T00:00:00.000Z",
  "endsAt": "2025-10-08T00:00:00.000Z",
  "participants": []
}
```
**Returns:** Competition ID + verification code

### Add Participant
**Endpoint:** `POST /v2/competitions/{id}/participants`
**Body:**
```json
{
  "verificationCode": "123-456-789",
  "participants": ["username"]
}
```

### Get Competition Details (Standings)
**Endpoint:** `GET /v2/competitions/{id}`
**Returns:** Full competition with participant progress

### Delete Competition (if needed)
**Endpoint:** `DELETE /v2/competitions/{id}`
**Body:**
```json
{
  "verificationCode": "123-456-789"
}
```

## Implementation Steps

### 1. Update WOM Client
Add methods:
- `CreateCompetition(title, metric, startsAt, endsAt)`
- `AddParticipant(competitionID, username, verificationCode)`
- `GetCompetitionDetails(competitionID)`
- `DeleteCompetition(competitionID, verificationCode)` (optional)

### 2. Database Migration
- Create `wom_competitions` table
- Mark old tables as deprecated (don't drop yet)

### 3. Update Commands
**StartEvent:**
- Calculate start/end times (1 week from now)
- Create WOM competition
- Store competition ID + verification code
- Create Discord thread
- Post registration buttons

**RegisterForEvent:**
- Fetch WOM competition ID from DB
- Add participant to WOM competition
- Post confirmation in thread

**FinishEvent:**
- Fetch competition details from WOM
- Extract top 3 participants from standings
- Announce winners with embeds
- Competition remains on WOM for historical records

### 4. Remove Scheduler
- Delete progress tracking task
- Remove `internal/scheduler/` package
- Update main.go to not start scheduler

## Data Flow Comparison

### Old Flow
1. User runs `/botw wildy callisto`
2. Bot creates DB record, thread, embed
3. User clicks "Register"
4. Bot fetches current KC from WOM
5. Bot stores starting KC in DB
6. **Every 15 min:** Bot fetches current KC, stores snapshot
7. User runs `/botw finish`
8. Bot fetches final KC from WOM
9. Bot calculates gains (end - start)
10. Bot announces winners

### New Flow
1. User runs `/botw wildy callisto`
2. Bot creates WOM competition (7 days), gets ID + code
3. Bot stores WOM ID + code in DB
4. Bot creates Discord thread, embed
5. User clicks "Register"
6. Bot adds user to WOM competition
7. **WOM handles all tracking automatically**
8. User runs `/botw finish`
9. Bot fetches standings from WOM
10. Bot announces winners (WOM already calculated gains)

## Migration Strategy

### Phase 1: Keep Both Systems
- Add WOM integration alongside existing system
- Test WOM integration thoroughly
- Compare results between systems

### Phase 2: Switch to WOM
- Update commands to use WOM
- Keep old data for historical reference
- Don't delete old tables yet

### Phase 3: Cleanup (Later)
- Remove old tracking code
- Archive old data
- Drop deprecated tables

## Edge Cases

### User Not on WOM
- WOM automatically creates player on first add
- No special handling needed

### Competition Already Ended
- Use WOM's historical data
- Standings preserved forever on WOM

### Verification Code Lost
- Store in DB, never show to users
- Bot always has access

## Benefits Summary

**For Users:**
- View competition on wiseoldman.net
- See real-time progress updates
- Historical records of all competitions
- Better mobile experience

**For Bot:**
- Less code to maintain
- More reliable tracking
- No background jobs needed
- Leverage WOM's infrastructure

**For Developers:**
- Standard API integration
- Well-documented endpoints
- Active community support
- Continuous improvements by WOM team
