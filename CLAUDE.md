# Disciplo - Community Platform MVP

Community platform with web onboarding, built with Go, PocketBase, and Telegram Bot API.

## Vision

Create a community platform that uses a webapp access control, organization, and content management -> messaging with Telegram. The community is formed with members and telegram private groups through a bot-as-gatekeeper approach.

## Similar projects
- mymembers (telegram)
- https://www.ternarydev.com/ (discord)
- patreon
- memberful

## Phase 0: Webapp admin dashboard

- PocketBase integration for data persistence (always use pocketbase go sdk api if possible)
- .env file is used for setup
- Auto-create PocketBase collections with proper schemas
- Auto-create first admin from .env
- SMTP email integration via PocketBase (config from .env)
- Organized and modular codebase
- Bot template system externalized
- Basic bot structure with modular architecture
- Once automatic setup is done send an email to admin to link telegram id (deeplink), bot automatically starts and gives an inline button to go dashboard HOST
- Admin dashboard with:
    * profile page (with connect to telegram button -> token deeplink in /start to connect admin with telegram_id)
    * groups page (will be connected later by the bot)
    * members page (all members -> phase 2 for the onboarding)
- pocketbase collections (list of minimal fields):
- users:
    * name
    * email
    * password
    * created (date)
    * json with fields defined later
    * admin (bool; only admin=true can access admin dashboard)
    * group_admin (-> record groups; array of groups where member is admin)
    * group_admin_since (date)
    * groups (-> record grousp; array of groups where member belongs to)
    * status (pending; accepted)
    * verified (bool, verified=true when telegram_id is filled)
    * telegram_id
    * telegram_name (telegram handle)
- communities:
    * name
    * description
    * created (date)
    * telegram_id (this is sent from bot, once bot becomes admin, do it in phase 3)
    * type: (default | local | special) (only one default group can be present, check this when turning type=default)

**Overview PocketBase Collections:**
- `users` - Member profiles with group membership and status
- `communities` - Community metadata with general/local group flags
- `requests` - Pending member requests with admin approval workflow

## Group Structure

- **General Group**: Super group where all active members belong
- **Local Groups**: Exclusive membership (one at a time per user)
- **Local-admin**: One per local group, rotates every 2-3 months
- Special Group: request access to admin

## Architecture Decisions

**Web-First Design**: User onboarding and management happens primarily through web interface, not bot interactions

## Technology Stack

- **Backend**: Go 1.24.6
- **Database**: PocketBase (auto-create collections)
- **Bot**: go-telegram-bot-api/v5
- **Web Frontend**: Vanilla JS or Alpine.js (no React)
- **Email**: Pocketbase SMTP integration for invitations

## Development Commands

```bash
# Development mode with verbose logging
make dev

# Kill application if it's still running
make kill

# Build for production
make build

# Run the built application
make run

# Clean build artifacts and database
make clean
```

**Setup Steps:**
1. Copy `.env` file and update with your credentials
2. Run `make dev` to start in development mode
3. Access PocketBase admin at `http://localhost:8080/_/` 
4. Setup pocketbase SMTP from .env
5. Default admin user will be created automatically

## Project Structure

```
disciplo/
├── src/                   # ALL SOURCE CODE
│   ├── main.go           # Application entry point
│   ├── bot/              # Telegram bot commands
│   ├── collections/      # PocketBase collection schemas
│   ├── config/           # Environment configuration loader
│   ├── email/            # PocketBase email integration
│   ├── pb/               # PocketBase database utilities
│   ├── utils/            # Token generation utilities
│   ├── web/              # Web registration server
│   └── static/           # Source templates and CSS
│       ├── templates/    # HTML templates
│       └── styles.css    # CSS styles
├── build/                # COMPLETE RUNTIME ENVIRONMENT
│   ├── disciplo          # Compiled executable
│   ├── pb_public/        # Static files (copied from src/static)
│   ├── pb_data/          # Database (created at runtime)
│   └── .env              # Configuration (copied from root)
├── Makefile              # Build and development commands
├── .env                  # Environment configuration
├── go.mod                # Go module definition
└── go.sum                # Go dependencies
```

**Key Benefits:**
- **Clean separation**: All source in `src/`, all runtime in `build/`
- **Fresh start**: Delete `build/` and rebuild completely anytime
- **Self-contained**: `build/` directory contains everything needed to run
- **Source preservation**: Original templates stay in source, get copied to build

## Email and Message Templates

### Email Templates System

Email templates are externalized and customizable:

**Location**: `src/static/email_templates/`
- `admin_invitation.html` - Welcome email sent to admin on startup
- `admin_notification.html` - New member request notification
- `user_onboarding.html` - Welcome email sent to approved users

**Features**:
- Go template syntax with data interpolation
- Fallback to hardcoded templates if file not found
- Templates copied to `build/pb_public/email_templates/` during build
- Admins can customize email content in `build/pb_public/email_templates/` without code changes or restarting the server
- design grade interface, professional, clean, minimalistic
- no colors, or very few colors
- mobile first

**Template Variables Available**:
- `{{.Subject}}` - Email subject line
- `{{.UserName}}` - Member name
- `{{.AdminEmail}}` - Admin email address
- `{{.BotUsername}}` - Telegram bot username
- `{{.Token}}` - Registration token
- `{{.ApprovalLink}}` - Admin approval link


## Development Standards & References

### **Required Reading Before Code Changes**
1. **PocketBase Official Docs**: https://pocketbase.io/docs/ 
   - Must use official methods for email, file handling, database operations
   - Check docs before implementing custom solutions
2. **Telegram Bot API**: https://core.telegram.org/bots/api
   - Follow official patterns for bot commands and updates
3. **Webapp components**: https://ui.shadcn.com/docs/components 
    - use state of the art modern components
    - design grade interface, professional, clean, minimalistic
    - no colors, or very few colors
    - mobile first

### **Configuration Standards**
- All configurable values must be in .env file
- No hardcoded values (bot username, admin name, ports, etc.)
- Respect PORT environment variable if provided
- Use PocketBase official methods over custom implementations

### **Code Quality Standards**
- Proper error handling for all database operations
- Use PocketBase official email sending methods
- Follow Go best practices and naming conventions
- Test all components before marking as complete
- Test all bot commands before marking as complete
- Less is more: use robust, minimal code, don't create complex or redundant structures

- dont ask for permission to execute cli commands within the project
- use env variables from env dont hardcode!
- killa sempre i processi che lanci, quando finisci il comando