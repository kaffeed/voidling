# Event Notification Role Feature - Complete

## Implementation Summary

Successfully implemented the ability for administrators to set a role that gets pinged when events (BOTW, SOTW, Mass) are created.

## Changes Made

### 1. Database Schema
**File:** `migrations/00007_add_event_notification_role.sql`
- Added `event_notification_role_id INTEGER` column to `guild_config` table
- Includes proper up/down migration support

### 2. Database Queries
**File:** `queries/guild_config.sql`
- Added `UpdateEventNotificationRole` query to update the notification role
- Updated `UpsertGuildConfig` to include the new column

### 3. Command Implementation
**File:** `internal/commands/config.go`
- Added `HandleSetEventNotificationRole()` function
  - Admin-only access check
  - Role ID parsing and validation
  - Database update with proper error handling
  - Success confirmation message
- Updated `HandleShowConfig()` to display the event notification role
  - Added role display with mention format `<@&roleID>`
  - Shows "Not configured" if role not set
  - Added default timezone display as well

### 4. Bot Command Registration
**File:** `internal/bot/bot.go`
- Added `/config set-event-notification-role` subcommand definition
  - Role parameter (required)
  - Clear description of functionality
- Added routing in `handleConfigCommand()` to new handler

### 5. Event Creation Updates
**File:** `internal/commands/trackable.go`
- Modified `createTrackableEvent()` to ping notification role
  - Fetches guild config before sending announcement
  - Includes role mention in content if configured
  - Applied to both BOTW and SOTW events

**File:** `internal/commands/schedulable.go`
- Modified `HandleMassEvent()` to ping notification role
  - Fetches guild config before sending confirmation
  - Includes role mention in content if configured

### 6. Documentation
**File:** `README.md`
- Added `/config set-event-notification-role` to Admin Commands section

## Usage

### Setting the Notification Role
```
/config set-event-notification-role role:@EventNotifications
```
- Requires Administrator permission
- Role will be mentioned when BOTW, SOTW, or Mass events are created
- Confirmation message shows the configured role

### Viewing Current Configuration
```
/config show
```
Displays:
- Coordinator Role
- Competition Code Channel
- Event Notification Role (NEW)
- Default Timezone

### Event Creation Behavior
When events are created, the configured role will be mentioned:
- `/botw wildy boss:Callisto` → mentions @EventNotifications in announcement
- `/sotw start skill:Mining` → mentions @EventNotifications in announcement
- `/mass activity:Nex ...` → mentions @EventNotifications in announcement

## Technical Details

### Database Schema
```sql
ALTER TABLE guild_config ADD COLUMN event_notification_role_id INTEGER;
```

### sqlc Generated Code
- `GuildConfig` struct now includes `EventNotificationRoleID sql.NullInt64`
- `UpdateEventNotificationRoleParams` struct for query parameters

### Error Handling
- Gracefully handles missing guild config (creates if needed)
- Validates role exists before storing
- Admin permission check prevents unauthorized access
- Handles database errors with user-friendly messages

## Build Status
✅ **Build Successful** - All changes compile without errors

## Testing Checklist
- [x] Build passes
- [x] SQL queries generated correctly
- [x] Command registered in bot
- [x] Admin permission check in place
- [x] Role mention format correct
- [ ] Manual testing required:
  - Set notification role
  - Create BOTW event (verify role ping)
  - Create SOTW event (verify role ping)
  - Create Mass event (verify role ping)
  - View config (verify role displays)

## Migration Path
1. Run database migration: `make migrate-up` (automatic on bot startup)
2. Restart bot to register new command
3. Admin runs `/config set-event-notification-role role:@YourRole`
4. Role will be mentioned on all future event creations

## Future Enhancements (Not Implemented)
- Separate roles per event type (BOTW role, SOTW role, Mass role)
- Option to remove/clear the notification role
- Notification role for event completions/winners
- Audit log for role changes
