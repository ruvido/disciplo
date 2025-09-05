# New Membership Request - {{.AppName}}

A new user has applied to join **{{.AppName}}**!

## Applicant Details

**Name:** {{.Name}}  
**Email:** {{.Email}}  
**Location:** {{.City}}, {{.Location}}  
**Job Field:** {{.JobField}}  
**Date of Birth:** {{.DateOfBirth}}  

**Interests:**  
{{range .Interests}}- {{.}}  
{{end}}

**Why they want to join:**  
> {{.WhyJoin}}

**Profile Picture:** [View Picture]({{.ProfilePictureURL}})

---

## Quick Actions

**✅ [Approve Application]({{.ApprovalURL}})**

🔍 [Review Full Application]({{.ReviewURL}})  
📊 [Admin Dashboard]({{.DashboardURL}})

---

**Request Details:**  
- **Request ID:** {{.RequestID}}  
- **Submitted:** {{.SubmissionDate}}  
- **Status:** {{.Status}}

*This email was sent automatically from {{.AppName}}*