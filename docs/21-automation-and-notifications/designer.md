# Automation and Notifications Module - Designer Documentation

## Notification Center

```
┌─ Notifications ─────────────────────┐
│ ┌─────────────────────────────────┐ │
│ │ 🏖️ Leave Approved               │ │
│ │    Your leave for Jan 20-22     │ │
│ │    has been approved            │ │
│ │    2 hours ago                  │ │
│ └─────────────────────────────────┘ │
│                                      │
│ ┌─────────────────────────────────┐ │
│ │ 💰 Salary Slip Ready            │ │
│ │    January 2024 salary slip     │ │
│ │    is now available             │ │
│ │    [Download]                   │ │
│ │    5 hours ago                  │ │
│ └─────────────────────────────────┘ │
│                                      │
│ ┌─────────────────────────────────┐ │
│ │ ⏰ Missing Attendance            │ │
│ │    Please mark attendance for   │ │
│ │    January 14, 2024             │ │
│ │    [Mark Now]                   │ │
│ │    1 day ago                    │ │
│ └─────────────────────────────────┘ │
│                                      │
│ [Mark all as read]                   │
└──────────────────────────────────────┘
```

## Email Notification Settings

```
┌─ Notification Preferences ──────────┐
│                                      │
│ Leave Management                     │
│ ☑ Leave applications                │
│ ☑ Leave approvals                   │
│ ☑ Leave balance warnings            │
│ ☐ Team leave calendar               │
│                                      │
│ Attendance                           │
│ ☑ Missing attendance alerts         │
│ ☐ Weekly attendance summary         │
│ ☐ Late arrival notifications        │
│                                      │
│ Payroll                              │
│ ☑ Salary slip available             │
│ ☑ Tax documents                     │
│ ☐ Payroll announcements             │
│                                      │
│ Approvals                            │
│ ☑ Pending approval reminders        │
│ ☑ Daily digest (9:00 AM)            │
│ ☐ Real-time alerts                  │
│                                      │
│ Delivery Method                      │
│ ☑ Email                             │
│ ☑ In-app notifications              │
│ ☑ Mobile push (if app installed)    │
│ ☐ SMS (critical only)               │
│                                      │
│ [Save Preferences]                   │
└──────────────────────────────────────┘
```

## Automated Email Template

```
┌──────────────────────────────────────┐
│ From: HRMS <noreply@company.com>     │
│ To: john.doe@company.com             │
│ Subject: Leave Application Approved  │
│                                      │
│ ┌────────────────────────────────┐  │
│ │  [Company Logo]                │  │
│ │                                │  │
│ │  Dear John Doe,                │  │
│ │                                │  │
│ │  Your leave application has    │  │
│ │  been approved.                │  │
│ │                                │  │
│ │  Leave Details                 │  │
│ │  ────────────────────────────  │  │
│ │  Type: Casual Leave            │  │
│ │  From: January 20, 2024        │  │
│ │  To: January 22, 2024          │  │
│ │  Total Days: 3                 │  │
│ │                                │  │
│ │  Approved by: Jane Smith       │  │
│ │  (Manager)                     │  │
│ │                                │  │
│ │  Remaining Balance: 9 days     │  │
│ │                                │  │
│ │  [View in HRMS]                │  │
│ │                                │  │
│ │  ────────────────────────────  │  │
│ │  This is an automated email    │  │
│ └────────────────────────────────┘  │
└──────────────────────────────────────┘
```

## Scheduled Reports Dashboard

```
┌─ Scheduled Reports ─────────────────┐
│                                      │
│ ┌─────────────────────────────────┐ │
│ │ Weekly Attendance Report        │ │
│ │ Frequency: Every Monday 9:00 AM │ │
│ │ Recipients: All Managers        │ │
│ │ Format: Excel                   │ │
│ │ Status: ✓ Active                │ │
│ │ Last Run: Jan 15, 2024          │ │
│ │ [Edit] [Pause] [Run Now]        │ │
│ └─────────────────────────────────┘ │
│                                      │
│ ┌─────────────────────────────────┐ │
│ │ Monthly Payroll Summary         │ │
│ │ Frequency: 1st of month         │ │
│ │ Recipients: Finance Team        │ │
│ │ Format: PDF                     │ │
│ │ Status: ✓ Active                │ │
│ │ Next Run: Feb 1, 2024           │ │
│ │ [Edit] [Pause]                  │ │
│ └─────────────────────────────────┘ │
│                                      │
│ [Create New Schedule]                │
└──────────────────────────────────────┘
```

## Workflow Automation Rules

```
┌─ Automation Rules ──────────────────┐
│                                      │
│ ┌─────────────────────────────────┐ │
│ │ Rule: Auto-assign Leave Approver│ │
│ │ Trigger: On Leave Application   │ │
│ │          submission             │ │
│ │ Action: Assign to reporting mgr │ │
│ │ Status: ✓ Enabled               │ │
│ │ [Edit]                          │ │
│ └─────────────────────────────────┘ │
│                                      │
│ ┌─────────────────────────────────┐ │
│ │ Rule: Escalate Pending Approvals│ │
│ │ Trigger: After 3 days pending   │ │
│ │ Action: Notify senior manager   │ │
│ │ Status: ✓ Enabled               │ │
│ │ [Edit]                          │ │
│ └─────────────────────────────────┘ │
│                                      │
│ [Add New Rule]                       │
└──────────────────────────────────────┘
```

## Color Scheme
- New notification: Blue dot
- Unread: Bold text
- Read: Gray text
- Action required: Orange badge
- Success notification: Green icon
- Error notification: Red icon
