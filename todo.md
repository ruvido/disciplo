# Disciplo Development TODO

## ✅ PHASE 0 - COMPLETED (MVP Web Dashboard)

### Core Infrastructure - ✅ DONE
- ✅ User verification status (verified=false until telegram connected)
- ✅ Telegram reconnection link on profile page
- ✅ Sign-out button implementation  
- ✅ Complete profile editing with password change
- ✅ Telegram bot inline button HOST fix
- ✅ Environment variables (no hardcoded ports/configs)
- ✅ PocketBase integration with auto-setup
- ✅ Admin dashboard with authentication
- ✅ Email template system
- ✅ Comprehensive routing system

### Security Hardening - ✅ DONE
- ✅ Input validation and sanitization
- ✅ Secure password change procedures
- ✅ Environment variable configuration
- ✅ PocketBase authentication integration
- ✅ Proper error handling throughout

## 🔄 PHASE 1 - MEMBER ONBOARDING WORKFLOW

### 1. Web Registration System
- [ ] Public registration page with form validation
- [ ] Member application review workflow for admins
- [ ] Email invitation system for approved members
- [ ] Member status tracking (pending → approved → active)
- [ ] Onboarding: 
    * Multi-step webform with modern components (same as shadcn): fields include -> name, email, password -> date of birth, location, job field (menu a tendina), interests (multiple choices) -> why you want to join (longer text), upload picture of yourself
    * data goes to requests (pb collection)
    * new record requests triggers email sending to EMAIL_REQUESTS (disciplo.toml config file, use for now ruvido@gmail.com)
    * admin receives email (new membership request)
    * user receives email (your submission has been received)
    * admin accepts requests -> email to user with deeplink to bot to connect profile with telegram-id
    * admin accepts requests changes status from pending to accepted, telegram linking success -> verified=true
    * only when status is accepted then user can login

### 2. Enhanced Bot Features  
- [ ] Member welcome flow with community rules
- [ ] Automated group invitation after approval
- [ ] Member verification and status commands
- [ ] Help system with command documentation

### 3. Admin Tools
- [ ] Bulk member management (approve/reject multiple)
- [ ] Member search and filtering
- [ ] Community analytics dashboard
- [ ] Email notification preferences

## 🏢 PHASE 2 - ADVANCED GROUP MANAGEMENT

### Group Operations
- [ ] Automatic Telegram group creation
- [ ] Local group assignment system
- [ ] Group admin rotation scheduling (2-3 months)
- [ ] Member transfer between groups

### Bot Group Management
- [ ] Group admin commands (/promote, /demote, /transfer)
- [ ] Member status sync between web and Telegram
- [ ] Group activity monitoring
- [ ] Automated group moderation rules

### Advanced Features
- [ ] Member engagement tracking
- [ ] Community metrics and reporting
- [ ] Integration with external services
- [ ] Mobile app considerations

## 🚀 IMMEDIATE NEXT PRIORITIES

1. **Member Registration Flow** (Phase 1 start)
2. **Bot Welcome System** enhancement  
3. **Admin Approval Workflow** implementation
4. **Email Integration** testing and optimization
