# Mobile App Module - Designer Documentation

## Mobile App Home Screen

```
┌────────────────────────────────┐
│  ☰                    🔔 (3)   │
│                                │
│  Good morning, John! 👋        │
│                                │
│  ┌──────────────────────────┐ │
│  │ 📅 Today: Jan 15, 2024   │ │
│  │ ⏰ 9:15 AM               │ │
│  │ [CHECK IN]               │ │
│  └──────────────────────────┘ │
│                                │
│  Quick Actions                 │
│  ┌────┐ ┌────┐ ┌────┐ ┌────┐ │
│  │ 🏖️  │ │ 💰  │ │ ✈️  │ │ 📊  │ │
│  │Leave│ │Exp │ │Trvl│ │Pay │ │
│  └────┘ └────┘ └────┘ └────┘ │
│                                │
│  Summary                       │
│  ┌──────────────────────────┐ │
│  │ Leave Balance: 12 days   │ │
│  │ Pending Approvals: 3     │ │
│  │ This Month: 18/20 days   │ │
│  └──────────────────────────┘ │
└────────────────────────────────┘
```

## Check-in/Check-out Screen

```
┌────────────────────────────────┐
│  ← Check In                    │
│                                │
│     ┌──────────────────┐       │
│     │                  │       │
│     │   [Selfie Area]  │       │
│     │                  │       │
│     └──────────────────┘       │
│                                │
│  📍 Location                   │
│  Office HQ - San Francisco     │
│  Distance: 15m from office     │
│  ✓ Within geofence             │
│                                │
│  ⏰ Current Time               │
│  9:15 AM                       │
│                                │
│  ┌──────────────────────────┐ │
│  │                          │ │
│  │     [CHECK IN NOW]       │ │
│  │                          │ │
│  └──────────────────────────┘ │
│                                │
│  Today's Schedule              │
│  • Team standup - 10:00 AM    │
│  • Client call - 2:00 PM      │
└────────────────────────────────┘
```

## Leave Application Form (Mobile)

```
┌────────────────────────────────┐
│  ← Apply for Leave             │
│                                │
│  Leave Type                    │
│  [Casual Leave      ▼]         │
│                                │
│  Duration                      │
│  ● Full Day  ○ Half Day        │
│                                │
│  From Date                     │
│  [Jan 20, 2024     📅]         │
│                                │
│  To Date                       │
│  [Jan 22, 2024     📅]         │
│                                │
│  Total: 3 days                 │
│  Available: 12 days            │
│                                │
│  Reason                        │
│  ┌──────────────────────────┐ │
│  │ Family vacation          │ │
│  └──────────────────────────┘ │
│                                │
│  [Submit Application]          │
└────────────────────────────────┘
```

## Approvals Screen

```
┌────────────────────────────────┐
│  ← Pending Approvals (3)       │
│                                │
│  ┌──────────────────────────┐ │
│  │ 🏖️ Leave Request          │ │
│  │ Jane Smith                │ │
│  │ Feb 5-7, 2024 (3 days)   │ │
│  │ Casual Leave              │ │
│  │                           │ │
│  │ [✓ Approve]  [✗ Reject]  │ │
│  └──────────────────────────┘ │
│                                │
│  ┌──────────────────────────┐ │
│  │ 💰 Expense Claim          │ │
│  │ Bob Wilson                │ │
│  │ $245.50 - Travel         │ │
│  │ 📎 3 receipts             │ │
│  │                           │ │
│  │ [View Details]            │ │
│  └──────────────────────────┘ │
│                                │
│  ┌──────────────────────────┐ │
│  │ ✈️ Travel Request         │ │
│  │ Alice Brown               │ │
│  │ NYC - Jan 25-28          │ │
│  │                           │ │
│  │ [Review]                  │ │
│  └──────────────────────────┘ │
└────────────────────────────────┘
```

## Salary Slip View (Mobile)

```
┌────────────────────────────────┐
│  ← Salary Slip                 │
│                                │
│  January 2024                  │
│                                │
│  Earnings                      │
│  Basic Salary      $3,000      │
│  HRA                 $600      │
│  Special Allow.      $400      │
│  ────────────────────────────  │
│  Gross Pay         $4,000      │
│                                │
│  Deductions                    │
│  Professional Tax    $200      │
│  Employee PF         $360      │
│  ────────────────────────────  │
│  Total Deductions    $560      │
│                                │
│  ════════════════════════════  │
│  NET PAY          $3,440       │
│  ════════════════════════════  │
│                                │
│  [Download PDF]  [Share]       │
└────────────────────────────────┘
```

## Design System
- Primary color: Company brand color
- Success: Green (#4CAF50)
- Warning: Orange (#FF9800)
- Error: Red (#F44336)
- Background: Light gray (#F5F5F5)
- Cards: White with shadow
- Icons: Material Design/SF Symbols
- Typography: System font (San Francisco/Roboto)

## Navigation
- Bottom tab bar: Home, Team, Approvals, Profile
- Hamburger menu for secondary features
- Swipe gestures for quick actions
- Pull to refresh on lists
