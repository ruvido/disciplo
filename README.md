# Disciplo - Community Platform MVP

**🔒 Security-First Community Platform** with web onboarding and Telegram integration.

A modern community management platform that bridges web-based administration with Telegram group messaging through an intelligent bot-as-gatekeeper architecture.

## 🚀 Quick Start

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd disciplo
   cp .env.example .env  # Create your .env file
   ```

2. **Configure environment** (edit `.env`):
   ```bash
   # Required: Telegram Bot
   BOT_TOKEN=your_bot_token
   BOT_USERNAME=your_bot_username
   
   # Required: Admin Account
   ADMIN_EMAIL=your@email.com
   ADMIN_PASSWORD=your_secure_password
   ADMIN_NAME=Your Name
   
   # Required: Production Host
   HOST=https://yourdomain.com
   PORT=8080
   ```

3. **Start development**:
   ```bash
   make dev
   ```

4. **Access your platform**:
   - **Dashboard**: http://localhost:8080
   - **PocketBase Admin**: http://localhost:8080/_/
   - **Connect Telegram**: Check your email for the invitation link

## 🛠️ Development Commands

```bash
make dev          # Start development server with hot reload
make build        # Build production executable
make run          # Run production build
make kill         # Stop all running processes
make clean        # Clean build artifacts and database
```

## 🏗️ Architecture

### Technology Stack
- **Backend**: Go 1.24+ with PocketBase embedded database
- **Frontend**: Vanilla JavaScript with modern CSS
- **Bot**: Telegram Bot API integration
- **Email**: SMTP via PocketBase
- **Auth**: PocketBase authentication system

### Project Structure
```
disciplo/
├── src/                   # Source code
│   ├── main.go           # Application entry point
│   ├── bot/              # Telegram bot logic
│   ├── collections/      # PocketBase schemas
│   ├── config/           # Environment configuration
│   ├── email/            # Email templates & sending
│   ├── web/              # Web routes & handlers
│   └── static/           # Templates and assets
├── build/                # Production build output
├── Makefile              # Build automation
└── .env                  # Configuration (not committed)
```

## 🔐 Security Features

- **Environment-based configuration** (no hardcoded secrets)
- **PocketBase authentication** with secure session management
- **Input validation** and sanitization throughout
- **HTTPS enforcement** in production
- **Telegram webhook security** with proper validation
- **Password security** with state-of-the-art change procedures

## 🌟 Features

### Phase 0 - MVP (Current)
- ✅ **Web Admin Dashboard** with authentication
- ✅ **User Profile Management** with Telegram connection
- ✅ **Community Management** (general/local/special groups)
- ✅ **Member Management** with approval workflow
- ✅ **Telegram Bot Integration** with inline keyboards
- ✅ **Email System** with customizable templates
- ✅ **Auto-setup** of database collections and admin user

### Planned Phases
- **Phase 1**: Member onboarding workflow
- **Phase 2**: Advanced group management
- **Phase 3**: Bot group administration features

## 📊 Database Schema

### Collections
- **`users`** - Member profiles with group membership and verification status
- **`communities`** - Community metadata with type classification (default/local/special)
- **`requests`** - Pending member requests with admin approval workflow

### Key Fields
- `verified` - Boolean flag for Telegram connection status
- `telegram_id` - Telegram user ID for bot integration
- `admin` - Admin access permissions
- `group_admin` - Community administration rights
- `status` - Member approval status (pending/accepted)

## 🔗 Integration Guide

### Telegram Bot Setup
1. Create bot with [@BotFather](https://t.me/BotFather)
2. Get bot token and username
3. Add to `.env` file
4. Bot automatically connects on startup

### SMTP Configuration
1. Configure SMTP settings in `.env`
2. SMTP gets auto-configured in PocketBase
3. Email templates are customizable in `build/pb_public/email_templates/`

## 🛡️ Production Deployment

### Environment Requirements
```bash
# Production .env template
BOT_TOKEN=prod_bot_token_here
BOT_USERNAME=your_prod_bot
ADMIN_EMAIL=admin@yourdomain.com
ADMIN_PASSWORD=secure_password_here
ADMIN_NAME=Admin Name
HOST=https://yourdomain.com
PORT=8080
SMTP_HOST=your.smtp.server
SMTP_PORT=587
SMTP_USER=smtp_username
SMTP_PASS=smtp_password
SMTP_FROM="Your Community <noreply@yourdomain.com>"
```

### Deployment Steps
1. Set production environment variables
2. Run `make build` to create production executable
3. Deploy `build/` directory to your server
4. Run `./disciplo` from the `build/` directory
5. Configure reverse proxy (nginx/caddy) for HTTPS
6. Set up monitoring and logging

## 📝 Contributing

- Follow Go best practices and naming conventions
- Use PocketBase SDK methods (never raw SQL)
- Test all functionality before committing
- Maintain security-first approach
- Keep code minimal and robust

## 📞 Support

- **Issues**: Create GitHub issues for bugs and feature requests
- **Documentation**: All configuration options are in CLAUDE.md
- **Security**: Report security issues privately to maintainers

---

**Built with ❤️ for secure community management**