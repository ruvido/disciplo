# Registration Received - {{.AppName}}

Hello **{{.Name}}**,

Thank you for applying to join **{{.AppName}}**! We have received your application and it is currently under review.

## Your Application Summary

**Personal Information:**
- **Name:** {{.Name}}
- **Email:** {{.Email}}  
- **Location:** {{.City}}, {{.Location}}
- **Job Field:** {{.JobField}}
- **Date of Birth:** {{.DateOfBirth}}

**Your Interests:**
{{range .Interests}}- {{.}}  
{{end}}

**Why you want to join:**
> {{.WhyJoin}}

**Profile Picture:** ✅ Uploaded successfully

---

## Next Steps

1. **Review Process:** Our team will review your application within **2-3 business days**
2. **Approval Notification:** You will receive an email confirmation once your application is approved
3. **Telegram Connection:** After approval, you'll receive instructions to connect your Telegram account
4. **Dashboard Access:** You'll be able to access your member dashboard and community features

## Timeline Expectations

- **Review:** 2-3 business days
- **Decision:** You'll be notified by email
- **Onboarding:** Immediate after approval

## Questions?

If you have any questions about your application or the community, feel free to contact us at **{{.AdminEmail}}**.

---

**Application Details:**
- **Request ID:** {{.RequestID}}
- **Submitted:** {{.SubmissionDate}}
- **Status:** {{.Status}}

Welcome to the {{.AppName}} community journey!

Best regards,  
**{{.AppName}} Team**

*This is an automated confirmation email from {{.AppName}}*