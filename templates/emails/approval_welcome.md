# 🎉 Welcome to {{.AppName}}, {{.Name}}!

**Congratulations!** Your application has been **approved** and you are now officially part of the **{{.AppName}}** community!

## 🚀 Get Started - Connect Your Telegram Account

To complete your onboarding and access the community, please connect your Telegram account:

**👉 [Connect to Telegram Bot]({{.BotDeepLink}})**

This will:
- Link your profile with your Telegram account
- Give you access to community groups
- Enable notifications and updates
- Activate your full member dashboard

---

## 📋 Community Rules & Guidelines

Welcome to our community! Please review these important guidelines:

### **Core Values**
- **Respect:** Treat all members with respect and kindness
- **Collaboration:** Help others and share knowledge freely  
- **Growth:** Support each other's personal and professional development
- **Quality:** Contribute meaningful content and discussions

### **Communication Guidelines**
- Keep discussions relevant and constructive
- No spam, self-promotion, or off-topic content
- Use appropriate channels for different topics
- Be patient with new members

### **Community Behavior**
- Professional conduct is expected at all times
- Disagreements should be handled respectfully
- Report any issues to administrators
- Respect privacy and confidentiality

---

## 🎯 Onboarding Guide

### **Step 1: Complete Telegram Setup**
1. Click the [Telegram connection link]({{.BotDeepLink}}) above
2. Start a conversation with **@{{.BotUsername}}**
3. Follow the bot instructions to verify your account
4. You'll be automatically added to relevant community groups

### **Step 2: Access Your Dashboard**
- Visit your member dashboard: [{{.AppName}} Dashboard]({{.DashboardURL}})
- Complete your profile if needed
- Explore available community features

### **Step 3: Introduce Yourself**
- Join the welcome channel on Telegram
- Share a brief introduction with the community
- Let others know about your interests: {{range .Interests}}**{{.}}**, {{end}}

### **Step 4: Explore & Engage**
- Browse community discussions
- Join groups related to your interests
- Participate in events and activities

---

## 🔗 Quick Links

- **📱 Telegram Bot:** [@{{.BotUsername}}](https://t.me/{{.BotUsername}})
- **🏠 Dashboard:** [Member Dashboard]({{.DashboardURL}})
- **❓ Support:** {{.AdminEmail}}
- **📚 Community Guidelines:** [Available in Dashboard]({{.DashboardURL}})

---

## 🤝 Your Profile Summary

We're excited to have you in our community! Here's a reminder of your profile:

- **Location:** {{.City}}, {{.Location}}
- **Profession:** {{.JobField}}  
- **Interests:** {{range .Interests}}{{.}}, {{end}}
- **Member Since:** {{.ApprovalDate}}

---

**Ready to join the conversation?**

**[🚀 Connect Telegram Now]({{.BotDeepLink}})**

Welcome aboard, {{.Name}}! We can't wait to see what you'll contribute to our community.

Best regards,  
**{{.AppName}} Team**

---
*Approved on {{.ApprovalDate}} | Request ID: {{.RequestID}}*  
*This welcome email was sent automatically from {{.AppName}}*